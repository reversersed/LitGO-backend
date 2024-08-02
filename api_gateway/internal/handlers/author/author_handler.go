package author

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/authors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Find authors
// @Description  there can be multiple search parameters, id or translit, or both
// @Description  example: ?id=1&id=2&translit=author-21&id=3&translit=author-756342
// @Tags         authors
// @Produce      json
// @Param		 id      query     string 		false 		"Author Id, must be a primitive id hex"
// @Param		 translit   query     string 		false 		"Translit author name"
// @Success      200  {array}   authors_pb.GetAuthorsResponse.Authors 		"Authors"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Field was not in a correct format"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Authors not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Some internal error"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Router       /authors [get]
func (h *handler) GetAuthors(c *gin.Context) {

	var request authors_pb.GetAuthorsRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	// need to made
	reply, err := h.client.GetAuthors(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reply)
}
