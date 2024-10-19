package review

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mocks "github.com/reversersed/LitGO-backend/tree/main/api_gateway/internal/handlers/mocks"
	"github.com/reversersed/LitGO-backend/tree/main/api_gateway/pkg/middleware"
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	mock_reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews/mocks"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
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
		MockBehaviour  func(*mocks.MockLogger, *mocks.MockJwtMiddleware, *mock_reviews_pb.MockReviewClient)
		ExceptedStatus int
		ExceptedBody   string
	}{
		{
			Name:   "book reviews default request",
			Path:   "/api/v1/reviews/book/translit1",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mrc *mock_reviews_pb.MockReviewClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mrc.EXPECT().GetBookReviews(gomock.Any(), &reviews_pb.GetBookReviewsRequest{Id: "translit1", Page: 0, PageSize: 1}).Return(nil, status.Error(codes.InvalidArgument, "specify correct id"))
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"specify correct id\",\"details\":[]}",
		},
		{
			Name:   "book reviews request with parameters",
			Path:   "/api/v1/reviews/book/id?page=2&pagesize=5",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mrc *mock_reviews_pb.MockReviewClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mrc.EXPECT().GetBookReviews(gomock.Any(), &reviews_pb.GetBookReviewsRequest{Id: "id", Page: 2, PageSize: 5}).Return(&reviews_pb.GetBookReviewsResponse{Reviews: []*reviews_pb.ReviewModel{{Text: "review text", Creator: &reviews_pb.UserModel{Login: "user"}}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "[{\"text\":\"review text\",\"creator\":{\"login\":\"user\"}}]",
		},
		{
			Name:   "book reviews request with one parameter",
			Path:   "/api/v1/reviews/book/dsgfd34fdcvbw3r?pagesize=20",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mrc *mock_reviews_pb.MockReviewClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mrc.EXPECT().GetBookReviews(gomock.Any(), &reviews_pb.GetBookReviewsRequest{Id: "dsgfd34fdcvbw3r", Page: 0, PageSize: 20}).Return(&reviews_pb.GetBookReviewsResponse{Reviews: []*reviews_pb.ReviewModel{{Text: "review text", Creator: &reviews_pb.UserModel{Login: "user"}}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "[{\"text\":\"review text\",\"creator\":{\"login\":\"user\"}}]",
		},
		{
			Name:   "create book review request without user",
			Path:   "/api/v1/reviews/book/id",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byt, _ := json.Marshal(&reviews_pb.CreateBookReviewRequest{})
				return bytes.NewBuffer(byt)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mrc *mock_reviews_pb.MockReviewClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) {
					c.Next()
				})
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(nil, status.Error(codes.Unauthenticated, "no user found"))
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
			},
			ExceptedStatus: http.StatusUnauthorized,
			ExceptedBody:   "{\"code\":16,\"type\":\"Unauthenticated\",\"message\":\"no user found\",\"details\":[]}",
		},
		{
			Name:   "create book review request error from service",
			Path:   "/api/v1/reviews/book/id",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byt, _ := json.Marshal(&reviews_pb.CreateBookReviewRequest{
					Text:   "text",
					Rating: 2.0,
				})
				return bytes.NewBuffer(byt)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mrc *mock_reviews_pb.MockReviewClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) {
					c.Next()
				})
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{Id: "userid"}, nil)
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mrc.EXPECT().CreateBookReview(gomock.Any(), &reviews_pb.CreateBookReviewRequest{CreatorId: "userid", ModelId: "id", Text: "text", Rating: 2.0}).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"user not found\",\"details\":[]}",
		},
		{
			Name:   "create book review request success",
			Path:   "/api/v1/reviews/book/id",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byt, _ := json.Marshal(&reviews_pb.CreateBookReviewRequest{
					Text:   "text",
					Rating: 2.0,
				})
				return bytes.NewBuffer(byt)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mrc *mock_reviews_pb.MockReviewClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) {
					c.Next()
				})
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{Id: "userid"}, nil)
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mrc.EXPECT().CreateBookReview(gomock.Any(), &reviews_pb.CreateBookReviewRequest{CreatorId: "userid", ModelId: "id", Text: "text", Rating: 2.0}).Return(&reviews_pb.CreateBookReviewResponse{Review: &reviews_pb.ReviewModel{Text: "text", Rating: 2.0, UserAction: reviews_pb.UserActionEnum_like}}, nil)
			},
			ExceptedStatus: http.StatusCreated,
			ExceptedBody:   "{\"text\":\"text\",\"rating\":2,\"userAction\":1}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_reviews_pb.NewMockReviewClient(ctrl)
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
