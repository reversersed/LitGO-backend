package review

import (
	"net/http"

	"github.com/gin-gonic/gin"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// @Summary      Get book's reviews
// @Description  Find's all book's reviews with page and pagesize arguments
// @Tags         reviews
// @Produce      json
// @Param		 id 		path 	string 	true 	"Book ID or translit name"
// @Param		 page 		query 	int 	false 	"Page number, must be greater or equal than 0, optional parameter"
// @Param		 pagesize 	query 	int 	false 	"Page size, default 1, must be greater or equal than 1"
// @Success      200  {array}   reviews_pb.ReviewModel "Reviews array"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Invalid request"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "No reviews found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Internal error occurred"
// @Failure      501  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Route not implemented yet"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /reviews/book/{id} [get]
func (h *handler) GetBookReviews(c *gin.Context) {
	var request reviews_pb.GetBookReviewsRequest
	if err := c.BindQuery(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	if request.GetPageSize() < 1 {
		request.PageSize = 1
	}
	request.Id = c.Param("id")
	response, err := h.client.GetBookReviews(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response.GetReviews())
}

// @Summary      Creates review on book
// @Description  Adds review to specified book and updates it's rating
// @Tags         reviews
// @Produce      json
// @Param		 id 		path 	string 	true 	"Book ID or translit name"
// @Param 		 body 		body 	reviews_pb.CreateBookReviewRequest true "Request body"
// @Success      201  {object}   reviews_pb.ReviewModel "Created review"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Invalid request"
// @Failure      401  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "User not authorized"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Book or user not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Internal error occurred"
// @Failure      501  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Route not implemented yet"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} "Service does not responding (maybe crush)"
// @Router       /reviews/book/{id} [post]
func (h *handler) CreateBookReview(c *gin.Context) {
	var request reviews_pb.CreateBookReviewRequest
	if err := c.BindJSON(&request); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	if user, err := h.jwt.GetCredentialsFromContext(c); err != nil {
		c.Error(err)
		return
	} else {
		request.CreatorId = user.GetId()
	}

	request.ModelId = c.Param("id")
	response, err := h.client.CreateBookReview(c.Request.Context(), &request)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, response.GetReview())
}
