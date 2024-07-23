package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/config"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type jwtMiddleware interface {
	Middleware(roles ...string) gin.HandlerFunc
	GenerateAccessToken(u *middleware.UserTokenModel) (string, string, error)
	GetUserClaims(token string) (*middleware.UserTokenModel, error)
}
type userHandler struct {
	connection *grpc.ClientConn
	client     users_pb.UserClient
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
		client:     client,
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
	var userRoutes = []struct {
		method  string
		route   string
		handler []gin.HandlerFunc
	}{}
	for _, v := range userRoutes {
		router.Handle(v.method, v.route, v.handler...)
	}
	h.logger.Infof("user handler has been registered with %d routes", len(userRoutes))
}
