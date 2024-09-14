package book

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
	books_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/books"
)

type handler struct {
	client books_pb.BookClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func New(client books_pb.BookClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) *handler {
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
	general := router.Group("/api/v1/books")
	{
		_ = general.Group("/").Use(h.jwt.Middleware())
		{

		}
		general.GET("/suggest", h.GetBooksSuggestion)
	}
	h.logger.Info("book handler has been registered")
}
