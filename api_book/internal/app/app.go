package app

import (
	"context"
	"fmt"
	"net"

	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	"github.com/reversersed/go-grpc/tree/main/api_book/internal/config"
	srv "github.com/reversersed/go-grpc/tree/main/api_book/internal/service"
	"github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
	freecache "github.com/reversersed/go-grpc/tree/main/api_book/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/logging/logrus"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/mongo"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/shutdown"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/validator"
	"google.golang.org/grpc"
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
	app := &app{
		logger:  logger,
		config:  cfg,
		cache:   cache,
		service: srv.NewServer(logger, cache, storage, validator),
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
	return mongo.Close()
}
