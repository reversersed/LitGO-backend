package review

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
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
		auth := general.Group("").Use(h.jwt.Middleware())
		{
			auth.POST("/book/:id", h.CreateBookReview)
		}
		_ = general.Group("").Use(h.jwt.Middleware("admin"))
		{

		}
		general.GET("/book/:id", h.GetBookReviews)
	}
	h.logger.Info("review handler has been registered")
}
