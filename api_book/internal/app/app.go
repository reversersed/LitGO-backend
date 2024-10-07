package app

import (
	"context"
	"fmt"
	"io"
	"net"

	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	"github.com/reversersed/go-grpc/tree/main/api_book/internal/config"
	rabbitService "github.com/reversersed/go-grpc/tree/main/api_book/internal/rabbitmq"
	srv "github.com/reversersed/go-grpc/tree/main/api_book/internal/service"
	"github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
	freecache "github.com/reversersed/go-grpc/tree/main/api_book/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/logging/logrus"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/mongo"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/rabbitmq"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/shutdown"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New() (*app, error) {
	logger, err := logrus.GetLogger()
	if err != nil {
		return nil, err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	cache := freecache.NewFreeCache(104857600)
	validator := validator.New()
	dbClient, err := mongo.NewClient(context.Background(), cfg.Database)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	storage := storage.NewStorage(dbClient, cfg.Database.Base, logger)

	genreConnection, err := grpc.NewClient(cfg.Server.GenreServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	genreClient := genres_pb.NewGenreClient(genreConnection)

	authorConnection, err := grpc.NewClient(cfg.Server.AuthorServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	authorClient := authors_pb.NewAuthorClient(authorConnection)

	rabbitMqServer, err := rabbitmq.New(cfg.Rabbit)
	if err != nil {
		return nil, err
	}

	app := &app{
		logger:      logger,
		config:      cfg,
		cache:       cache,
		service:     srv.NewServer(logger, cache, storage, validator, genreClient, authorClient, rabbitService.New(rabbitMqServer.Connection, logger)),
		connections: []*grpc.ClientConn{genreConnection, authorConnection},
		closers:     []io.Closer{rabbitMqServer},
	}

	return app, nil
}

func (a *app) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.Server.Port))
	if err != nil {
		a.logger.Errorf("failed to listen: %v", err)
		return err
	}
	go shutdown.Graceful(a)
	server := grpc.NewServer()
	a.service.Register(server)

	a.logger.Infof("starting listening %s:%d...", a.config.Server.Host, a.config.Server.Port)
	if err := server.Serve(listener); err != nil {
		a.logger.Errorf("failed to start server: %v", err)
		return err
	}
	return nil
}
func (a *app) Close() error {
	for _, c := range a.connections {
		if err := c.Close(); err != nil {
			return err
		}
	}
	for _, c := range a.closers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return mongo.Close()
}
