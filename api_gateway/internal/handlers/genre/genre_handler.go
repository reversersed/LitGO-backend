package genre

import (
	"net/http"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/genres"
<<<<<<< HEAD
=======
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
>>>>>>> 560fcec (separated projects and proto files)
=======
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
>>>>>>> 4e07393 (added genre routes with tests)
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

// @Summary      Get category tree
// @Description  Searches category tree based on category or genre id or translate name
// @Description  Query can be: category id, category translit name, genre id or genre translit name
// @Description  If genre id or name matches, it returns whole category that contains that genre
// @Tags         genres
// @Produce      json
// @Param        query query genres_pb.GetOneOfRequest true "Request body"
// @Success      200  {object}   genres_pb.CategoryModel "Category"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Received wrong query"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Category not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Internal error occurred"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /genres/tree [get]
func (h *handler) GetGenreTree(c *gin.Context) {
	var request genres_pb.GetOneOfRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}

	response, err := h.client.GetTree(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.GetCategory())
}

// @Summary      Get category or genre
// @Description  Searches category or genre based on id or translit name
// @Description  Query can be: category id, category translit name, genre id or genre translit name
// @Description  If category found, it returns whole category with nested genre. Otherwise it returns a single genre
// @Tags         genres
// @Produce      json
// @Param        query query genres_pb.GetOneOfRequest true "Request body"
// @Success      200  {object}   genre.GetOneOfGenre.HandleResponse "Response body. Only one field will be presented"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Received wrong query"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Category or genre not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Internal error occurred"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /genres [get]
func (h *handler) GetOneOfGenre(c *gin.Context) {
	var request genres_pb.GetOneOfRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}

	response, err := h.client.GetOneOf(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	type HandleResponse struct {
		Category *genres_pb.CategoryModel `json:"category,omitempty"`
		Genre    *genres_pb.GenreModel    `json:"genre,omitempty"`
	}
	data := HandleResponse{
		Category: response.GetCategory(),
		Genre:    response.GetGenre(),
	}
	c.JSON(http.StatusOK, data)
}
