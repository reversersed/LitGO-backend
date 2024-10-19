package genre

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mocks "github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/mocks"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/pkg/middleware"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	mock_genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres/mocks"
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
		MockBehaviour  func(*mocks.MockLogger, *mocks.MockJwtMiddleware, *mock_genres_pb.MockGenreClient)
		ExceptedStatus int
		ExceptedBody   string
	}{
		{
			Name:   "Get all genres success",
			Path:   "/api/v1/genres/all",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetAll(gomock.Any(), &genres_pb.Empty{}).Return(&genres_pb.GetAllResponse{Categories: []*genres_pb.CategoryModel{{Name: "Category"}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "[{\"name\":\"Category\"}]",
		},
		{
			Name:   "Get genres error",
			Path:   "/api/v1/genres/all",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetAll(gomock.Any(), &genres_pb.Empty{}).Return(nil, status.Error(codes.NotFound, "no genres in database"))
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"no genres in database\",\"details\":[]}",
		},
		{
			Name:   "Get category tree",
			Path:   "/api/v1/genres/tree?query=query+request",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetTree(gomock.Any(), &genres_pb.GetOneOfRequest{Query: "query request"}).Return(&genres_pb.CategoryResponse{Category: &genres_pb.CategoryModel{Name: "category", Genres: []*genres_pb.GenreModel{{Name: "genre nested"}}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"name\":\"category\",\"genres\":[{\"name\":\"genre nested\"}]}",
		},
		{
			Name:   "Get category tree error",
			Path:   "/api/v1/genres/tree?query=query+request",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetTree(gomock.Any(), &genres_pb.GetOneOfRequest{Query: "query request"}).Return(nil, status.Error(codes.NotFound, "category not found"))
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"category not found\",\"details\":[]}",
		},
		{
			Name:   "Search category",
			Path:   "/api/v1/genres?query=query+request",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetOneOf(gomock.Any(), &genres_pb.GetOneOfRequest{Query: "query request"}).Return(&genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: &genres_pb.CategoryModel{Name: "category", Genres: []*genres_pb.GenreModel{{Name: "genre nested"}}}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"category\":{\"name\":\"category\",\"genres\":[{\"name\":\"genre nested\"}]}}",
		},
		{
			Name:   "Search genre",
			Path:   "/api/v1/genres?query=query+request",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetOneOf(gomock.Any(), &genres_pb.GetOneOfRequest{Query: "query request"}).Return(&genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: &genres_pb.GenreModel{Name: "category"}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"genre\":{\"name\":\"category\"}}",
		},
		{
			Name:   "Search category error",
			Path:   "/api/v1/genres?query=query+request",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_genres_pb.MockGenreClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetOneOf(gomock.Any(), &genres_pb.GetOneOfRequest{Query: "query request"}).Return(nil, status.Error(codes.NotFound, "category not found"))
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"category not found\",\"details\":[]}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_genres_pb.NewMockGenreClient(ctrl)
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
