package handlers

import "github.com/gin-gonic/gin"

type Logger interface {
	Infof(format string, args ...interface{})
	Info(...interface{})
}
type jwtMiddleware interface {
	Middleware(...string) gin.HandlerFunc
}
