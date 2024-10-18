package review

import (
	"github.com/gin-gonic/gin"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
)

type handler struct {
	client reviews_pb.ReviewClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func New(client reviews_pb.ReviewClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) *handler {
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
	general := router.Group("/api/v1/reviews")
	{
		_ = general.Group("").Use(h.jwt.Middleware())
		{

		}
		_ = general.Group("").Use(h.jwt.Middleware("admin"))
		{

		}

	}
	h.logger.Info("review handler has been registered")
}
