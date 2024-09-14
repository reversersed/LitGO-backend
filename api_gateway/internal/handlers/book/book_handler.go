package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	books_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/books"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Find books by keywords
// @Description  find books by provided phares, keys or names
// @Tags         books
// @Produce      json
// @Param		 query      query     string 		true 		"Query with keywords"
// @Param		 limit   query     int 		false 		"limit books to display. default = 5 if not specified, min = 1, max = 10"
// @Success      200  {array}   books_pb.BookModel 		"Books"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Query was empty"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Books not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Some internal error"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Router       /books/suggest [get]
func (h *handler) GetBooksSuggestion(c *gin.Context) {
	var request books_pb.GetSuggestionRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	if request.GetLimit() == 0 {
		request.Limit = 5
	}

	reply, err := h.client.GetBookSuggestions(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, reply.GetBooks())
}
