package handlers

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=general.go -destination=mocks/general.go

type Logger interface {
	Infof(format string, args ...interface{})
	Info(...interface{})
}
type JwtMiddleware interface {
	Middleware(...string) gin.HandlerFunc
}
