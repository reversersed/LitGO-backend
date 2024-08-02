package author

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/authors"
)

type handler struct {
	client authors_pb.AuthorClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func New(client authors_pb.AuthorClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) *handler {
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
	general := router.Group("/api/v1/authors")
	{
		_ = general.Group("/").Use(h.jwt.Middleware())
		{

		}
		general.GET("/", h.GetAuthors)
	}
	h.logger.Info("author handler has been registered")
}
