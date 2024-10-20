package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_middleware "github.com/reversersed/LitGO-backend/tree/main/api_gateway/pkg/middleware/mocks"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testSecretKey string = "secretStringKey"

var userId = primitive.NewObjectID().Hex()

func generateToken(exp time.Duration) string {
	signer, _ := jwt.NewSignerHS(jwt.HS256, []byte(testSecretKey))
	builder := jwt.NewBuilder(signer)

	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        userId,
			Audience:  []string{"user"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		Roles: []string{"user"},
		Login: "user",
		Email: "user@example.com",
	}
	token, _ := builder.Build(claims)
	return token.String()
}
func TestGetCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_middleware.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	jwt := NewJwtMiddleware(logger, testSecretKey)

	router := gin.Default()
	router.Use(ErrorHandler)
	router.Group("/success").Use(jwt.Middleware()).GET("/", func(ctx *gin.Context) {
		credential, err := jwt.GetCredentialsFromContext(ctx)
		assert.NoError(t, err)
		assert.Equal(t, userId, credential.Id)
		assert.Equal(t, "user", credential.Login)
		assert.Equal(t, []string{"user"}, credential.Roles)
		ctx.Status(http.StatusOK)
	})
	router.GET("/error", func(ctx *gin.Context) {
		_, err := jwt.GetCredentialsFromContext(ctx)
		assert.Error(t, err)
		ctx.Error(err)
	})

	r := httptest.NewRequest(http.MethodGet, "/success/", nil)
	r.AddCookie(&http.Cookie{
		Name:   TokenCookieName,
		Value:  generateToken(time.Second * 5),
		Path:   "/",
		MaxAge: 30,
	})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	r = httptest.NewRequest(http.MethodGet, "/error", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}
func TestMiddleware(t *testing.T) {
	table := []struct {
		name           string
		key            string
		request        func() *http.Request
		roles          []string
		mockBehaviour  func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer)
		exceptedStatus int
	}{
		{
			name: "successful authorization",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				r.AddCookie(&http.Cookie{
					Name:   RefreshCookieName,
					Value:  primitive.NewObjectID().Hex(),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any()).AnyTimes()
				logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			},
			exceptedStatus: http.StatusOK,
		},
		{
			name: "empty token cookie",
			key:  testSecretKey,
			request: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			exceptedStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong secret key",
			key:  "",
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())
			},
			exceptedStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong token cookie",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  "randomtoken",
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any())
			},
			exceptedStatus: http.StatusUnauthorized,
		},
		{
			name: "user successful role authorization",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any()).AnyTimes()
				logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			},
			roles:          []string{"user"},
			exceptedStatus: http.StatusOK,
		},
		{
			name: "user unsuccessful role authorization",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any()).AnyTimes()
				logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			},
			roles:          []string{"user", "admin"},
			exceptedStatus: http.StatusForbidden,
		},
		{
			name: "old token without refresh",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(-time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any()).AnyTimes()
			},
			exceptedStatus: http.StatusUnauthorized,
		},
		{
			name: "service error",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(-time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				r.AddCookie(&http.Cookie{
					Name:   RefreshCookieName,
					Value:  primitive.NewObjectID().Hex(),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any())
				userServer.EXPECT().UpdateToken(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unauthenticated, "error"))
			},
			exceptedStatus: http.StatusUnauthorized,
		},
		{
			name: "token successful updated",
			key:  testSecretKey,
			request: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.AddCookie(&http.Cookie{
					Name:   TokenCookieName,
					Value:  generateToken(-time.Second * 5),
					Path:   "/",
					MaxAge: 30,
				})
				r.AddCookie(&http.Cookie{
					Name:   RefreshCookieName,
					Value:  primitive.NewObjectID().Hex(),
					Path:   "/",
					MaxAge: 30,
				})
				return r
			},
			mockBehaviour: func(logger *mock_middleware.MockLogger, userServer *mock_middleware.MockUserServer) {
				logger.EXPECT().Info(gomock.Any())
				userServer.EXPECT().UpdateToken(gomock.Any(), gomock.Any()).Return(&users_pb.TokenReply{
					Token:        "sometoken",
					Refreshtoken: "sometoken",
				}, nil)
				logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
			},
			exceptedStatus: http.StatusOK,
		},
	}
	gin.SetMode(gin.ReleaseMode)
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := mock_middleware.NewMockLogger(ctrl)
			server := mock_middleware.NewMockUserServer(ctrl)

			if tt.mockBehaviour != nil {
				tt.mockBehaviour(logger, server)
			}
			jwt := NewJwtMiddleware(logger, tt.key)
			jwt.ApplyUserServer(server)

			router := gin.Default()
			router.Use(ErrorHandler)
			router.Use(jwt.Middleware(tt.roles...))
			router.GET("/", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			r := tt.request()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.exceptedStatus, w.Result().StatusCode)
		})
	}
}
