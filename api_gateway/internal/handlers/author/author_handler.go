package author

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
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
// @Success      200  {array}   authors_pb.AuthorModel 		"Authors"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Field was not in a correct format"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Authors not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Some internal error"
// @Failure      501  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Route not implemented yet"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Router       /authors [get]
func (h *handler) GetAuthors(c *gin.Context) {
	var request authors_pb.GetAuthorsRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}

	reply, err := h.client.GetAuthors(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reply.GetAuthors())
}

// @Summary      Find authors by keywords
// @Description  find authors by provided phares, keys or names
// @Tags         authors
// @Produce      json
// @Param		 query      query     string 		true 		"Query with keywords"
// @Param		 limit   query     int 		false 		"limit authors to display. default = 5 if not specified, min = 1, max = 50"
// @Param		 page   query     int 		false 		"page number to find, must be greater or equal than 0"
// @Param		 rating	query	float32 false	"rating to find, must be 0 <= rating <= 5"
// @Success      200  {array}   authors_pb.AuthorModel 		"Authors"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Query was empty or not validated"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Authors not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Some internal error"
// @Failure      501  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Route not implemented yet"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Router       /authors/search [get]
func (h *handler) FindAuthors(c *gin.Context) {
	var request authors_pb.FindAuthorsRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	if request.GetLimit() == 0 {
		request.Limit = 5
	}

	reply, err := h.client.FindAuthors(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reply.GetAuthors())
}
