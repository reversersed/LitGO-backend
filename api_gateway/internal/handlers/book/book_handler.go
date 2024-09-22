package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
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

// TODO write tests for this handler
// @Summary      Create new book
// @Description  Creates new book by request
// @Description  Request must be multipart/form data only
// @Tags         books
// @Produce      json
// @Accept 		 x-www-form-urlencoded
// @Param		 data   formData     books_pb.CreateBookRequest 		true 		"Request model. Must be multipart/form data"
// @Param		 File formData file true "epub format book file"
// @Param		 Picture formData file true "book cover picture"
// @Success      201  {array}   books_pb.BookModel 		"Book created"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Wrong request received"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Authors or genre not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Some internal error"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Security 	 ApiKeyAuth
// @Router       /books/ [post]
func (h *handler) CreateBook(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	var req books_pb.CreateBookRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	response, err := h.client.CreateBook(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response)
}
