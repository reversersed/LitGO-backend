package user

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
)

type userHandler struct {
	client users_pb.UserClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func NewUserHandler(client users_pb.UserClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) (*userHandler, error) {
	return &userHandler{
		client: client,
		logger: logger,
		jwt:    jwtMiddleware,
	}, nil
}
func (h *userHandler) Close() error {
	return nil
}
func (h *userHandler) RegisterRouter(router *gin.Engine) {
	general := router.Group("/api/v1/users")
	{
		authorized := general.Group("/").Use(h.jwt.Middleware())
		{
			authorized.GET("/auth", h.UserAuthenticate)
		}
		general.GET("/", h.UserSearch)
		general.POST("/login", h.UserLogin)
		general.POST("/signin", h.UserRegister)
	}
	h.logger.Info("user handler has been registered")
}
