package collection

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers"
	collections_pb "github.com/reversersed/LitGO-proto/gen/go/collections"
)

type handler struct {
	client collections_pb.CollectionClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func New(client collections_pb.CollectionClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) *handler {
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
	general := router.Group("/api/v1/collections")
	{
		_ = general.Group("").Use(h.jwt.Middleware())
		{

		}
	}
	h.logger.Info("collection handler has been registered")
}
