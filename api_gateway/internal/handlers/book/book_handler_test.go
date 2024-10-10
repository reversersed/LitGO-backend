package book

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	mock_books_pb "github.com/reversersed/LitGO-proto/gen/go/books/mock"
	mocks "github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/mocks"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlers(t *testing.T) {
	table := []struct {
		Name           string
		Path           string
		Method         string
		Body           func() io.Reader
		MockBehaviour  func(*mocks.MockLogger, *mocks.MockJwtMiddleware, *mock_books_pb.MockBookClient)
		ExceptedStatus int
		ExceptedBody   string
	}{
		{
			Name:   "suggestion query",
			Path:   "/api/v1/books/suggest",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mbc *mock_books_pb.MockBookClient) {
				mjm.EXPECT().Middleware(gomock.Any()).Return(func(c *gin.Context) { c.Next() }).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()

				mbc.EXPECT().GetBookSuggestions(gomock.Any(), gomock.Any(), gomock.Any()).Return(&books_pb.GetBooksResponse{Books: []*books_pb.BookModel{{Name: "book name"}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "[{\"name\":\"book name\"}]",
		},
		{
			Name:   "suggestion query error",
			Path:   "/api/v1/books/suggest",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mbc *mock_books_pb.MockBookClient) {
				mjm.EXPECT().Middleware(gomock.Any()).Return(func(c *gin.Context) { c.Next() }).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()

				mbc.EXPECT().GetBookSuggestions(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "books not found"))
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"books not found\",\"details\":[]}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_books_pb.NewMockBookClient(ctrl)
			logger := mocks.NewMockLogger(ctrl)
			jwt := mocks.NewMockJwtMiddleware(ctrl)
			v.MockBehaviour(logger, jwt, client)

			h := New(client, logger, jwt)
			gin.SetMode(gin.TestMode)

			e := gin.Default()
			e.Use(middleware.ErrorHandler)
			h.RegisterRouter(e)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(v.Method, v.Path, v.Body())
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			b, _ := io.ReadAll(w.Result().Body)
			assert.Equal(t, v.ExceptedBody, string(b))

			er := h.Close()
			assert.NoError(t, er)
		})
	}
}
func TestCreateBookHandler(t *testing.T) {
	table := []struct {
		Name           string
		BodyFunction   func(*multipart.Writer, *testing.T)
		MockBehaviour  func(*mock_books_pb.MockBookClient, *testing.T)
		ExceptedStatus int
		ExceptedBody   string
	}{
		{
			Name: "missing book file",
			BodyFunction: func(w *multipart.Writer, t *testing.T) {
				f, err := w.CreateFormField("Name")
				assert.NoError(t, err)
				f.Write([]byte("Книга"))
			},
			MockBehaviour:  func(mbc *mock_books_pb.MockBookClient, t *testing.T) {},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"book is required\",\"details\":[]}",
		},
		{
			Name: "missing cover file",
			BodyFunction: func(w *multipart.Writer, t *testing.T) {
				f, err := w.CreateFormFile("Book", "eragon")
				assert.NoError(t, err)
				f.Write([]byte("Книга"))
			},
			MockBehaviour:  func(mbc *mock_books_pb.MockBookClient, t *testing.T) {},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"cover is required\",\"details\":[]}",
		},
		{
			Name: "file duplicating",
			BodyFunction: func(w *multipart.Writer, t *testing.T) {
				f, err := w.CreateFormFile("Book", "eragon")
				assert.NoError(t, err)
				f.Write([]byte("Книга"))
				f.Write([]byte("Книга"))
				f, err = w.CreateFormFile("Book", "eragon")
				assert.NoError(t, err)
				f.Write([]byte("Книга"))
				f.Write([]byte("Книга"))
				f, err = w.CreateFormFile("Cover", "eragon")
				assert.NoError(t, err)
				f.Write([]byte("Картинка"))
				f, err = w.CreateFormFile("Cover", "eragon")
				assert.NoError(t, err)
				f.Write([]byte("Картинка"))
			},
			MockBehaviour:  func(mbc *mock_books_pb.MockBookClient, t *testing.T) {},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"you can upload only one file per field\",\"details\":[{\"field\":\"Book\",\"description\":\"length of file received array\",\"actualvalue\":\"2\"},{\"field\":\"Cover\",\"description\":\"length of picture received array\",\"actualvalue\":\"2\"}]}",
		},
		{
			Name: "wrong formats and size",
			BodyFunction: func(w *multipart.Writer, t *testing.T) {
				f, err := w.CreateFormFile("Book", "eragon")
				assert.NoError(t, err)
				twentyBuffer := make([]byte, 20*(1<<20))
				rand.Read(twentyBuffer)
				f.Write(twentyBuffer)
				f, err = w.CreateFormFile("Cover", "cover")
				assert.NoError(t, err)
				f.Write(twentyBuffer)
			},
			MockBehaviour:  func(mbc *mock_books_pb.MockBookClient, t *testing.T) {},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong file format\",\"details\":[{\"field\":\"Book\",\"tag\":\"format\",\"tagValue\":\"epub\",\"description\":\"file must be in .epub format\",\"actualvalue\":\"eragon\"},{\"field\":\"Cover\",\"tag\":\"format\",\"tagValue\":\"jpg | jpeg | png\",\"description\":\"available formats: jpg | jpeg | png\",\"actualvalue\":\"cover\"},{\"field\":\"Cover\",\"tag\":\"size\",\"tagValue\":\"5242880\",\"description\":\"picture must be less than 5 MB size\",\"actualvalue\":\"20971520\"},{\"field\":\"Book\",\"tag\":\"size\",\"tagValue\":\"15728640\",\"description\":\"file must be less than 15 MB size\",\"actualvalue\":\"20971520\"}]}",
		},
		{
			Name: "error from service",
			BodyFunction: func(w *multipart.Writer, t *testing.T) {
				f, err := w.CreateFormFile("Book", "eragon.epub")
				assert.NoError(t, err)
				f.Write([]byte("book"))

				f, err = w.CreateFormFile("Cover", "cover.png")
				assert.NoError(t, err)
				f.Write([]byte("cover"))

				f, err = w.CreateFormField("Name")
				assert.NoError(t, err)
				f.Write([]byte("Название,? длиннее! 24-х символов для проверки длины названия файлов."))

				f, err = w.CreateFormField("Description")
				assert.NoError(t, err)
				f.Write([]byte("Описание"))

				f, err = w.CreateFormField("Genre")
				assert.NoError(t, err)
				f.Write([]byte("genreid"))

				f, err = w.CreateFormField("Authors")
				assert.NoError(t, err)
				f.Write([]byte("authorid1,authorid2,authorid3"))
			},
			MockBehaviour: func(mbc *mock_books_pb.MockBookClient, t *testing.T) {
				mbc.EXPECT().CreateBook(gomock.Any(), gomock.AssignableToTypeOf(&books_pb.CreateBookRequest{})).DoAndReturn(func(ctx context.Context, req *books_pb.CreateBookRequest, calloption ...grpc.CallOption) (*books_pb.CreateBookResponse, error) {
					assert.EqualValues(t, []string{"authorid1", "authorid2", "authorid3"}, req.GetAuthors())
					return nil, status.Error(codes.InvalidArgument, "wrong id format")
				})
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong id format\",\"details\":[]}",
		},
		{
			Name: "succeeded request",
			BodyFunction: func(w *multipart.Writer, t *testing.T) {
				f, err := w.CreateFormFile("Book", "eragon.epub")
				assert.NoError(t, err)
				f.Write([]byte("book"))

				f, err = w.CreateFormFile("Cover", "cover.png")
				assert.NoError(t, err)
				f.Write([]byte("cover"))

				f, err = w.CreateFormField("Name")
				assert.NoError(t, err)
				f.Write([]byte("Название,? длиннее! 24-х символов для проверки длины названия файлов."))

				f, err = w.CreateFormField("Description")
				assert.NoError(t, err)
				f.Write([]byte("Описание"))

				f, err = w.CreateFormField("Genre")
				assert.NoError(t, err)
				f.Write([]byte("genreid"))

				f, err = w.CreateFormField("Authors")
				assert.NoError(t, err)
				f.Write([]byte("authorid1"))

				f, err = w.CreateFormField("Authors")
				assert.NoError(t, err)
				f.Write([]byte("authorid2"))

				f, err = w.CreateFormField("Authors")
				assert.NoError(t, err)
				f.Write([]byte("authorid3"))
			},
			MockBehaviour: func(mbc *mock_books_pb.MockBookClient, t *testing.T) {
				mbc.EXPECT().CreateBook(gomock.Any(), gomock.AssignableToTypeOf(&books_pb.CreateBookRequest{})).DoAndReturn(func(ctx context.Context, req *books_pb.CreateBookRequest, calloption ...grpc.CallOption) (*books_pb.CreateBookResponse, error) {
					assert.EqualValues(t, []string{"authorid1", "authorid2", "authorid3"}, req.GetAuthors())
					return &books_pb.CreateBookResponse{
						Book: &books_pb.BookModel{
							Id:           "book id",
							Name:         "book name",
							Translitname: "translit-name-321312",
							Description:  "description",
							Picture:      "cover.png",
							Filepath:     "file.epub",
							Authors:      []*books_pb.AuthorModel{{Name: "author1"}, {Name: "author2"}},
							Category:     &books_pb.CategoryModel{Name: "category", Translitname: "translit-category-23312"},
							Genre:        &books_pb.GenreModel{Name: "genre", Translitname: "translit-genre-424"},
						},
					}, nil
				})
			},
			ExceptedStatus: http.StatusCreated,
			ExceptedBody:   "{\"book\":{\"id\":\"book id\",\"name\":\"book name\",\"translitname\":\"translit-name-321312\",\"description\":\"description\",\"picture\":\"cover.png\",\"filepath\":\"file.epub\",\"authors\":[{\"name\":\"author1\"},{\"name\":\"author2\"}],\"category\":{\"name\":\"category\",\"translitname\":\"translit-category-23312\"},\"genre\":{\"name\":\"genre\",\"translitname\":\"translit-genre-424\"}}}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock_books_pb.NewMockBookClient(ctrl)
			logger := mocks.NewMockLogger(ctrl)
			jwt := mocks.NewMockJwtMiddleware(ctrl)
			v.MockBehaviour(client, t)
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			jwt.EXPECT().Middleware(gomock.Any()).Return(func(c *gin.Context) { c.Next() }).AnyTimes()

			h := New(client, logger, jwt)
			gin.SetMode(gin.TestMode)

			e := gin.Default()
			e.Use(middleware.ErrorHandler)
			h.RegisterRouter(e)

			var buf bytes.Buffer
			w := multipart.NewWriter(&buf)
			v.BodyFunction(w, t)
			w.Close()

			wr := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/api/v1/books", &buf)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", w.FormDataContentType())
			e.ServeHTTP(wr, req)

			if assert.Equal(t, v.ExceptedStatus, wr.Result().StatusCode) {
				body, err := io.ReadAll(wr.Result().Body)
				assert.NoError(t, err)
				assert.Equal(t, v.ExceptedBody, string(body))
			}

			if v.ExceptedStatus == http.StatusCreated {
				os.MkdirAll("./files", os.FileMode(777))
				err = os.RemoveAll("./files")
				assert.NoError(t, err)
			}
		})
	}
}
