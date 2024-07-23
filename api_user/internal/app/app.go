package app

import (
	"fmt"
	"net"

	"github.com/reversersed/go-grpc/tree/main/api_user/internal/config"
	freecache "github.com/reversersed/go-grpc/tree/main/api_user/pkg/cache"
	"github.com/reversersed/go-grpc/tree/main/api_user/pkg/logging"
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
	app := &app{
		logger: logger,
		config: cfg,
		cache:  freecache.NewFreeCache(104857600),
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
