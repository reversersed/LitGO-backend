package app

import (
	"io"

	"github.com/reversersed/LitGO-backend/tree/main/api_file/internal/config"
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
type logger interface {
	Info(...any)
	Error(...any)
	Warn(...any)
	Fatal(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
	Warnf(string, ...any)
	Fatalf(string, ...any)
}
type RabbitListenerService any

type app struct {
	config                *config.Config
	logger                logger
	service               service
	cache                 cache
	RabbitListenerService RabbitListenerService
	closers               []io.Closer
}
