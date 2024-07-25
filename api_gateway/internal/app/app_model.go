package app

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/config"
)

type handler interface {
	RegisterRouter(router *gin.Engine)
	Close() error
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
	router   *gin.Engine
	config   *config.Config
	logger   logger
	cache    cache
	handlers []handler
}