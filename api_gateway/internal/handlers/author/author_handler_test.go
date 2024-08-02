package author

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mocks "github.com/reversersed/go-grpc/tree/main/api_gateway/internal/handlers/mocks"
	"github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/middleware"
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/authors"
	mock_authors_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/mocks/authors"
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
		MockBehaviour  func(*mocks.MockLogger, *mocks.MockJwtMiddleware, *mock_authors_pb.MockAuthorClient)
		ExceptedStatus int
		ExceptedBody   string
	}{
		{
			Name:   "Get authors success",
			Path:   "/api/v1/authors/?id=123&id=321&translit=421&translit=23",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_authors_pb.MockAuthorClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), &authors_pb.GetAuthorsRequest{Id: []string{"123", "321"}, Translit: []string{"421", "23"}}).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{{Name: "Author"}}}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"authors\":[{\"name\":\"Author\"}]}",
		},
		{
			Name:   "Get authors error",
			Path:   "/api/v1/authors/",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, mac *mock_authors_pb.MockAuthorClient) {
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.InvalidArgument, "wrong number of arguments"))
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong number of arguments\",\"details\":[]}",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			client := mock_authors_pb.NewMockAuthorClient(ctrl)
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
