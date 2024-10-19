package genre

import (
	"github.com/gin-gonic/gin"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
)

type handler struct {
	client genres_pb.GenreClient
	jwt    handlers.JwtMiddleware
	logger handlers.Logger
}

func New(client genres_pb.GenreClient, logger handlers.Logger, jwtMiddleware handlers.JwtMiddleware) *handler {
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
	general := router.Group("/api/v1/genres")
	{
		_ = general.Group("").Use(h.jwt.Middleware())
		{

		}
		general.GET("/all", h.GetAll)
		general.GET("/tree", h.GetGenreTree)
		general.GET("", h.GetOneOfGenre)
	}
	h.logger.Info("genre handler has been registered")
}
