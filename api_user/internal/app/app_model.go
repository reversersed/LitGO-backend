package app

import (
	"github.com/reversersed/go-grpc/tree/main/api_user/internal/config"
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
	Info(...interface{})
	Error(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Fatalf(string, ...interface{})
}
type app struct {
	config  *config.Config
	logger  logger
	service service
	cache   cache
}
