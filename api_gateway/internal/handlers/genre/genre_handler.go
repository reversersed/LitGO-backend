package genre

import (
	"net/http"

	"github.com/gin-gonic/gin"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/genres"
)

// @Summary      Get all genres
// @Description  Fetches all categories (with genres included)
// @Tags         genres
// @Produce      json
// @Success      200  {array}   genres_pb.CategoryModel "Genres fetched successfully"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "There's no genres in database"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Internal error occurred"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /genres/all [get]
func (h *handler) GetAll(c *gin.Context) {
	response, err := h.client.GetAll(c.Request.Context(), &genres_pb.Empty{})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.GetCategories())
}
