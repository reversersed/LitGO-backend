package app

import (
	"context"
	"fmt"
	"net"

	"github.com/reversersed/go-grpc/tree/main/api_user/internal/config"
	srv "github.com/reversersed/go-grpc/tree/main/api_user/internal/service"
	"github.com/reversersed/go-grpc/tree/main/api_user/internal/storage"
	freecache "github.com/reversersed/go-grpc/tree/main/api_user/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_user/pkg/logging"
	"github.com/reversersed/go-grpc/tree/main/api_user/pkg/mongo"
	"google.golang.org/grpc"
)

type service interface {
	Register(s grpc.ServiceRegistrar)
}
type cache interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte, expiration int) error
	Delete(key []byte) (affected bool)

	EntryCount() int64
}
type app struct {
	config  *config.Config
	logger  *logging.Logger
	service service
	cache   cache
}

func New() (*app, error) {
	logger, err := logging.GetLogger()
	if err != nil {
		return nil, err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	cache := freecache.NewFreeCache(104857600)
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
		service: srv.NewServer(cfg.Server.JwtSecret, logger, cache, storage),
	}

	return app, nil
}

func (a *app) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.Server.Port))
	if err != nil {
		a.logger.Errorf("failed to listen: %v", err)
		return err
	}
	server := grpc.NewServer()
	a.service.Register(server)

	a.logger.Infof("starting listening %s:%d...", a.config.Server.Host, a.config.Server.Port)
	if err := server.Serve(listener); err != nil {
		a.logger.Errorf("failed to start server: %v", err)
		return err
	}
	return nil
}
