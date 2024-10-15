package app

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/reversersed/go-grpc/tree/main/api_genre/internal/config"
	"github.com/reversersed/go-grpc/tree/main/api_genre/internal/rabbitmq"
	srv "github.com/reversersed/go-grpc/tree/main/api_genre/internal/service"
	"github.com/reversersed/go-grpc/tree/main/api_genre/internal/storage"
	freecache "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_genre/pkg/logging/logrus"
	"github.com/reversersed/go-grpc/tree/main/api_genre/pkg/mongo"
	rabbit "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/rabbitmq"
	"github.com/reversersed/go-grpc/tree/main/api_genre/pkg/shutdown"
	"github.com/reversersed/go-grpc/tree/main/api_genre/pkg/validator"
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

	rabbitMqConnection, err := rabbit.New(cfg.Rabbit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	rabbitService := rabbitmq.New(rabbitMqConnection.Connection, logger, storage)

	app := &app{
		logger:                logger,
		config:                cfg,
		cache:                 cache,
		RabbitListenerService: rabbitService,
		service:               srv.NewServer(logger, cache, storage, validator, rabbitService),
		closers:               []io.Closer{rabbitMqConnection},
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

	if err := a.RabbitListenerService.InitiateBookCreatedReceiver(); err != nil {
		return err
	}

	a.logger.Infof("starting listening %s:%d...", a.config.Server.Host, a.config.Server.Port)
	if err := server.Serve(listener); err != nil {
		a.logger.Errorf("failed to start server: %v", err)
		return err
	}
	return nil
}
func (a *app) Close() error {
	for _, c := range a.closers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return mongo.Close()
}
