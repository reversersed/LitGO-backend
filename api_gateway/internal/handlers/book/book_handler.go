package book

import (
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mdigger/translit"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// @Summary      Create new book
// @Description  Creates new book by request
// @Description  Request must be multipart/form data only
// @Tags         books
// @Produce      json
// @Accept 		 x-www-form-urlencoded
// @Param		 data   formData     books_pb.CreateBookRequest 		true 		"Request model. Must be multipart/form data"
// @Param		 Book formData file true "epub format book file"
// @Param		 Cover formData file true "book cover picture"
// @Success      201  {array}   books_pb.BookModel 		"Book created"
// @Failure      400  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Wrong request received"
// @Failure      404  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Authors or genre not found"
// @Failure      500  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Some internal error"
// @Failure      503  {object}  middleware.CustomError{details=[]shared_pb.ErrorDetail} 	"Service does not responding (maybe crush)"
// @Security 	 ApiKeyAuth
// @Router       /books [post]
func (h *handler) CreateBook(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}

	fileArray, exist := form.File["Book"]
	if !exist {
		c.Error(status.Error(codes.InvalidArgument, "book is required"))
		return
	}
	pictureArray, exist := form.File["Cover"]
	if !exist {
		c.Error(status.Error(codes.InvalidArgument, "cover is required"))
		return
	}
	if len(fileArray) != 1 || len(pictureArray) != 1 {
		stat, err := status.New(codes.InvalidArgument, "you can upload only one file per field").WithDetails(&shared_pb.ErrorDetail{
			Field:       "Book",
			Description: "length of file received array",
			Actualvalue: strconv.Itoa(len(fileArray)),
		}, &shared_pb.ErrorDetail{
			Field:       "Cover",
			Description: "length of picture received array",
			Actualvalue: strconv.Itoa(len(pictureArray)),
		})
		if err != nil {
			c.Error(status.Error(codes.InvalidArgument, "you can upload only one file per field"))
		} else {
			c.Error(stat.Err())
		}
		return
	}
	file := fileArray[0]
	picture := pictureArray[0]

	fileDetail := make([]*shared_pb.ErrorDetail, 0)
	a := strings.Split(file.Filename, ".")
	if a[len(a)-1] != "epub" {
		fileDetail = append(fileDetail, &shared_pb.ErrorDetail{
			Field:       "Book",
			Tag:         "format",
			Description: "file must be in .epub format",
			TagValue:    "epub",
			Actualvalue: a[len(a)-1],
		})
	}

	a = strings.Split(picture.Filename, ".")
	allowedPictureFormat := []string{"jpg", "jpeg", "png"}
	if !slices.Contains(allowedPictureFormat, a[len(a)-1]) {
		fileDetail = append(fileDetail, &shared_pb.ErrorDetail{
			Field:       "Cover",
			Tag:         "format",
			Description: "available formats: " + strings.Join(allowedPictureFormat, " | "),
			TagValue:    strings.Join(allowedPictureFormat, " | "),
			Actualvalue: a[len(a)-1],
		})
	}
	const MB = 1 << 20

	if picture.Size > 5*MB {
		fileDetail = append(fileDetail, &shared_pb.ErrorDetail{
			Field:       "Cover",
			Tag:         "size",
			Description: "picture must be less than 5 MB size",
			TagValue:    strconv.Itoa(5 * MB),
			Actualvalue: strconv.FormatInt(picture.Size, 10),
		})
	}
	if file.Size > 15*MB {
		fileDetail = append(fileDetail, &shared_pb.ErrorDetail{
			Field:       "Book",
			Tag:         "size",
			Description: "file must be less than 15 MB size",
			TagValue:    strconv.Itoa(15 * MB),
			Actualvalue: strconv.FormatInt(picture.Size, 10),
		})
	}
	if len(fileDetail) > 0 {
		stat := status.New(codes.InvalidArgument, "wrong file format")

		for _, d := range fileDetail {
			stat, err = stat.WithDetails(d)
			if err != nil {
				break
			}
		}
		if stat == nil {
			c.Error(status.Error(codes.InvalidArgument, "wrong file format"))
		} else {
			c.Error(stat.Err())
		}
		return
	}

	var req books_pb.CreateBookRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(status.Error(codes.InvalidArgument, err.Error()))
		return
	}
	rxSpaces := regexp.MustCompile(`\s+`)
	reg := regexp.MustCompile(`[^\p{L}\s]`)

	const MAX_NAME_LENGTH = 24 // max length of book name that goes to filename

	var fileName string = strings.Split(req.GetName(), ".")[0]
	if len(fileName) > MAX_NAME_LENGTH {
		fileName = fileName[:MAX_NAME_LENGTH]
	}
	req.Filepath = strings.ReplaceAll(strings.TrimSpace(rxSpaces.ReplaceAllString(translit.Ru(reg.ReplaceAllString(strings.ToLower(strings.ReplaceAll(fileName, "_", " ")), "")), " ")), " ", "_") + "_" + primitive.NewObjectID().Hex() + "." + strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
	req.Picture = strings.ReplaceAll(strings.TrimSpace(rxSpaces.ReplaceAllString(translit.Ru(reg.ReplaceAllString(strings.ToLower(strings.ReplaceAll(fileName, "_", " ")), "")), " ")), " ", "_") + "_" + primitive.NewObjectID().Hex() + "." + strings.Split(picture.Filename, ".")[len(strings.Split(picture.Filename, "."))-1]

	if len(req.GetAuthors()) == 1 {
		req.Authors = strings.Split(req.GetAuthors()[0], ",") // bcz curl from swagger sending arrays as values joined by ,
	}
	response, err := h.client.CreateBook(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	if err := c.SaveUploadedFile(file, "files/books/"+req.GetFilepath()); err != nil {
		c.Error(status.Error(codes.Internal, "error saving book: "+err.Error()))
		return
	}
	if err := c.SaveUploadedFile(picture, "files/book_covers/"+req.GetPicture()); err != nil {
		c.Error(status.Error(codes.Internal, "error saving cover: "+err.Error()))
		return
	}
	_ = os.Chmod("files/book_covers", os.FileMode(0007))
	_ = os.Chmod("files/books", os.FileMode(0007))
	c.JSON(http.StatusCreated, response)
}
