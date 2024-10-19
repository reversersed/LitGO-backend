package app

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/config"
)

type handler interface {
	RegisterRouter(router *gin.Engine)
	Close() error
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
type app struct {
	router   *gin.Engine
	config   *config.Config
	logger   logger
	handlers []handler
}
