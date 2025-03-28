package app

import (
	"fmt"
	"io"
	"net"

	freecache "github.com/reversersed/LitGO-backend-pkg/cache"
	"github.com/reversersed/LitGO-backend-pkg/logging/logrus"
	"github.com/reversersed/LitGO-backend-pkg/mongo"
	rabbit "github.com/reversersed/LitGO-backend-pkg/rabbitmq"
	"github.com/reversersed/LitGO-backend-pkg/shutdown"
	"github.com/reversersed/LitGO-backend-pkg/validator"
	"github.com/reversersed/LitGO-backend/tree/main/api_file/internal/config"
	"github.com/reversersed/LitGO-backend/tree/main/api_file/internal/rabbitmq"
	srv "github.com/reversersed/LitGO-backend/tree/main/api_file/internal/service"
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

	rabbitMqConnection, err := rabbit.New(cfg.Rabbit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	rabbitService := rabbitmq.New(rabbitMqConnection.Connection, logger)

	app := &app{
		logger:                logger,
		config:                cfg,
		cache:                 cache,
		RabbitListenerService: rabbitService,
		service:               srv.NewServer(logger, cache, validator, rabbitService, cfg.File),
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

	/*if err := a.RabbitListenerService.InitiateBookCreatedReceiver(); err != nil {
		return err
	}*/

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
