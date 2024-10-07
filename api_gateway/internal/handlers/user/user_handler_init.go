package user

import (
	"github.com/gin-gonic/gin"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
)

type handler struct {
	client users_pb.UserClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func New(client users_pb.UserClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) *handler {
	return &handler{
		client: client,
		logger: logger,
		jwt:    jwtMiddleware,
	}
}
func (h *handler) Close() error {
	return nil
}
func (h *handler) RegisterRouter(router *gin.Engine) {
	general := router.Group("/api/v1/users")
	{
		authorized := general.Group("").Use(h.jwt.Middleware())
		{
			authorized.GET("/auth", h.UserAuthenticate)
			authorized.POST("/logout", h.UserLogout)
		}
		general.GET("", h.UserSearch)
		general.POST("/login", h.UserLogin)
		general.POST("/signin", h.UserRegister)
	}
	h.logger.Info("user handler has been registered")
}
