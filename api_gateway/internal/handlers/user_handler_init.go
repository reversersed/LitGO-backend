package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/config"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userHandler struct {
	connection *grpc.ClientConn
	Client     users_pb.UserClient
	jwt        jwtMiddleware
	logger     Logger
}

func NewUserHandler(config *config.UrlConfig, logger Logger, jwtMiddleware jwtMiddleware) (*userHandler, error) {
	con, err := grpc.NewClient(config.UserServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := users_pb.NewUserClient(con)

	return &userHandler{
		connection: con,
		Client:     client,
		logger:     logger,
		jwt:        jwtMiddleware,
	}, nil
}
func (h *userHandler) Close() error {
	if err := h.connection.Close(); err != nil {
		return err
	}
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
