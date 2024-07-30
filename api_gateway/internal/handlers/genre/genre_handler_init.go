package genre

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/genres"
)

type genreHandler struct {
	client genres_pb.GenreClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func NewGenreHandler(client genres_pb.GenreClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) (*genreHandler, error) {
	return &genreHandler{
		client: client,
		logger: logger,
		jwt:    jwtMiddleware,
	}, nil
}
func (h *genreHandler) Close() error {
	return nil
}
func (h *genreHandler) RegisterRouter(router *gin.Engine) {
	general := router.Group("/api/v1/genres")
	{
		_ = general.Group("/").Use(h.jwt.Middleware())
		{

		}
		general.GET("/all", h.GetAll)
	}
	h.logger.Info("genre handler has been registered")
}
