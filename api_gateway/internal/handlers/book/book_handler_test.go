package book

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	mock_books_pb "github.com/reversersed/LitGO-proto/gen/go/books/mock"
	mocks "github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/mocks"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	"github.com/stretchr/testify/assert"
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
