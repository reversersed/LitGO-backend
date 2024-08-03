package handlers

import (
	"github.com/gin-gonic/gin"
	shared_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto"
)

//go:generate mockgen -source=general.go -destination=mocks/general.go

type Logger interface {
	Infof(format string, args ...interface{})
	Info(...interface{})
}
type JwtMiddleware interface {
	GetCredentialsFromContext(c *gin.Context) (*shared_pb.UserCredentials, error)
	Middleware(...string) gin.HandlerFunc
}
