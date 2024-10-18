package app

import (
	"context"
	"fmt"
	"io"
	"net"

	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"github.com/reversersed/go-grpc/tree/main/api_review/internal/config"
	srv "github.com/reversersed/go-grpc/tree/main/api_review/internal/service"
	"github.com/reversersed/go-grpc/tree/main/api_review/internal/storage"
	freecache "github.com/reversersed/go-grpc/tree/main/api_review/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_review/pkg/logging/logrus"
	"github.com/reversersed/go-grpc/tree/main/api_review/pkg/mongo"
	"github.com/reversersed/go-grpc/tree/main/api_review/pkg/shutdown"
	"github.com/reversersed/go-grpc/tree/main/api_review/pkg/validator"
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

	logger.Info("setting up user grpc client...")
	userConnection, err := grpc.NewClient(cfg.Server.UserServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	userClient := users_pb.NewUserClient(userConnection)

	app := &app{
		logger:  logger,
		config:  cfg,
		cache:   cache,
		service: srv.NewServer(logger, cache, storage, validator, userClient),
		closers: []io.Closer{},
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
	for _, c := range a.closers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return mongo.Close()
}
