package user

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
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	mock_users_pb "github.com/reversersed/LitGO-proto/gen/go/users/mocks"
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
		MockBehaviour  func(*mocks.MockLogger, *mocks.MockJwtMiddleware, *mock_users_pb.MockUserClient)
		ExtraCheck     func(*httptest.ResponseRecorder, *testing.T)
		ExceptedStatus int
		ExceptedBody   string
	}{
		{
			Name:   "User search success",
			Path:   "/api/v1/users?id=123&login=456&email=789",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{Id: "123", Login: "456", Email: "789"}).Return(&users_pb.UserModel{Login: "user"}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"login\":\"user\"}",
		},
		{
			Name:   "User search error",
			Path:   "/api/v1/users",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{}).Return(nil, status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong request\",\"details\":[]}",
		},
		{
			Name:   "User register success without remember me",
			Path:   "/api/v1/users/signin",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byts, _ := json.Marshal(&users_pb.RegistrationRequest{Login: "user", Password: "password", PasswordRepeat: "password", Email: "user@example.com"})
				return bytes.NewReader(byts)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().RegisterUser(gomock.Any(), &users_pb.RegistrationRequest{Login: "user", Password: "password", PasswordRepeat: "password", Email: "user@example.com"}).Return(&users_pb.LoginResponse{Login: "user", Roles: []string{"user"}, Token: "token", Refreshtoken: "rtoken"}, nil)
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Equal(t, c.MaxAge, 0)
						assert.Equal(t, c.Value, "token")
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.False(t, refresh, "excepted no refresh token in cookies")
			},
			ExceptedStatus: http.StatusCreated,
			ExceptedBody:   "{\"login\":\"user\",\"roles\":[\"user\"]}",
		},
		{
			Name:   "User register success with remember me",
			Path:   "/api/v1/users/signin",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byts, _ := json.Marshal(&users_pb.RegistrationRequest{Login: "user", Password: "password", PasswordRepeat: "password", Email: "user@example.com", RememberMe: true})
				return bytes.NewReader(byts)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().RegisterUser(gomock.Any(), &users_pb.RegistrationRequest{Login: "user", Password: "password", PasswordRepeat: "password", Email: "user@example.com", RememberMe: true}).Return(&users_pb.LoginResponse{Login: "user", Roles: []string{"user"}, Token: "token", Refreshtoken: "rtoken"}, nil)
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0)
						assert.Equal(t, c.Value, "token")
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0)
						assert.Equal(t, c.Value, "rtoken")
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.True(t, refresh, "excepted refresh token in cookies")
			},
			ExceptedStatus: http.StatusCreated,
			ExceptedBody:   "{\"login\":\"user\",\"roles\":[\"user\"]}",
		},
		{
			Name:   "User register error",
			Path:   "/api/v1/users/signin",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byts, _ := json.Marshal(&users_pb.RegistrationRequest{})
				return bytes.NewReader(byts)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().RegisterUser(gomock.Any(), &users_pb.RegistrationRequest{}).Return(nil, status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong request\",\"details\":[]}",
		},
		{
			Name:   "User login success without remember me",
			Path:   "/api/v1/users/login",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byts, _ := json.Marshal(&users_pb.LoginRequest{Login: "user", Password: "password", RememberMe: false})
				return bytes.NewReader(byts)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().Login(gomock.Any(), &users_pb.LoginRequest{Login: "user", Password: "password"}).Return(&users_pb.LoginResponse{Id: "id", Login: "user", Roles: []string{"user"}, Token: "token", Refreshtoken: "rtoken"}, nil)
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Equal(t, c.MaxAge, 0)
						assert.Equal(t, c.Value, "token")
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.False(t, refresh, "excepted no refresh token in cookies")
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"id\":\"id\",\"login\":\"user\",\"roles\":[\"user\"]}",
		},
		{
			Name:   "User login success with remember me",
			Path:   "/api/v1/users/login",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byts, _ := json.Marshal(&users_pb.LoginRequest{Login: "user", Password: "password", RememberMe: true})
				return bytes.NewReader(byts)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().Login(gomock.Any(), &users_pb.LoginRequest{Login: "user", Password: "password", RememberMe: true}).Return(&users_pb.LoginResponse{Login: "user", Roles: []string{"user"}, Token: "token", Refreshtoken: "rtoken"}, nil)
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0)
						assert.Equal(t, c.Value, "token")
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0)
						assert.Equal(t, c.Value, "rtoken")
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.True(t, refresh, "excepted refresh token in cookies")
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"login\":\"user\",\"roles\":[\"user\"]}",
		},
		{
			Name:   "User login failure",
			Path:   "/api/v1/users/login",
			Method: http.MethodPost,
			Body: func() io.Reader {
				byts, _ := json.Marshal(&users_pb.LoginRequest{})
				return bytes.NewReader(byts)
			},
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes()
				muc.EXPECT().Login(gomock.Any(), &users_pb.LoginRequest{}).Return(nil, status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong request\",\"details\":[]}",
		},
		{
			Name:   "User authenticate success",
			Path:   "/api/v1/users/auth",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{Id: "123", Login: "user", Roles: []string{"user"}}, nil)
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{Id: "123"}).Return(&users_pb.UserModel{Id: "123", Roles: []string{"user"}, Login: "user"}, nil)
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"id\":\"123\",\"login\":\"user\",\"roles\":[\"user\"]}",
		},
		{
			Name:   "User authenticate jwt failure",
			Path:   "/api/v1/users/auth",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(nil, status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedStatus: http.StatusBadRequest,
			ExceptedBody:   "{\"code\":3,\"type\":\"InvalidArgument\",\"message\":\"wrong request\",\"details\":[]}",
		},
		{
			Name:   "User authenticate service failure",
			Path:   "/api/v1/users/auth",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{}, nil)
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{}).Return(nil, status.Error(codes.NotFound, "user not found"))
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"user not found\",\"details\":[]}",
		},
		{
			Name:   "User authenticate removed roles",
			Path:   "/api/v1/users/auth",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{Id: "123", Login: "user", Roles: []string{"user", "admin"}}, nil)
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{Id: "123"}).Return(&users_pb.UserModel{Id: "123", Roles: []string{"user"}, Login: "user"}, nil)
				muc.EXPECT().UpdateToken(gomock.Any(), &users_pb.TokenRequest{Refreshtoken: "refresh cookie"}).Return(&users_pb.TokenReply{Token: "token", Refreshtoken: "newrefresh"}, nil)
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0, "token must has MaxAge greater than 0")
						assert.Equal(t, c.Value, "token")
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0, "token must has MaxAge greater than 0")
						assert.Equal(t, "newrefresh", c.Value)
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.True(t, refresh, "excepted refresh token in cookies")
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"id\":\"123\",\"login\":\"user\",\"roles\":[\"user\"]}",
		},
		{
			Name:   "User authenticate new roles",
			Path:   "/api/v1/users/auth",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{Id: "123", Login: "user", Roles: []string{"user"}}, nil)
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{Id: "123"}).Return(&users_pb.UserModel{Id: "123", Roles: []string{"user", "admin"}, Login: "user"}, nil)
				muc.EXPECT().UpdateToken(gomock.Any(), &users_pb.TokenRequest{Refreshtoken: "refresh cookie"}).Return(&users_pb.TokenReply{Token: "token", Refreshtoken: "newrefresh"}, nil)
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0, "token must has MaxAge greater than 0")
						assert.Equal(t, c.Value, "token")
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Greater(t, c.MaxAge, 0, "token must has MaxAge greater than 0")
						assert.Equal(t, "newrefresh", c.Value)
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.True(t, refresh, "excepted refresh token in cookies")
			},
			ExceptedStatus: http.StatusOK,
			ExceptedBody:   "{\"id\":\"123\",\"login\":\"user\",\"roles\":[\"user\",\"admin\"]}",
		},
		{
			Name:   "User authenticate error updating token",
			Path:   "/api/v1/users/auth",
			Method: http.MethodGet,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
				mjm.EXPECT().GetCredentialsFromContext(gomock.Any()).Return(&shared_pb.UserCredentials{Id: "123", Login: "user", Roles: []string{"user"}}, nil)
				muc.EXPECT().GetUser(gomock.Any(), &users_pb.UserRequest{Id: "123"}).Return(&users_pb.UserModel{Id: "123", Roles: []string{"user", "admin"}, Login: "user"}, nil)
				muc.EXPECT().UpdateToken(gomock.Any(), &users_pb.TokenRequest{Refreshtoken: "refresh cookie"}).Return(nil, status.Error(codes.NotFound, "refresh token not found"))
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Less(t, c.MaxAge, 0, "cookie must has MaxAge < 0 to be removed")
						assert.Equal(t, "", c.Value)
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Less(t, c.MaxAge, 0, "cookie must has MaxAge < 0 to be removed")
						assert.Equal(t, "", c.Value)
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.True(t, refresh, "excepted refresh token in cookies")
			},
			ExceptedStatus: http.StatusNotFound,
			ExceptedBody:   "{\"code\":5,\"type\":\"NotFound\",\"message\":\"refresh token not found\",\"details\":[]}",
		},
		{
			Name:   "User logout checking",
			Path:   "/api/v1/users/logout",
			Method: http.MethodPost,
			Body:   func() io.Reader { return nil },
			MockBehaviour: func(ml *mocks.MockLogger, mjm *mocks.MockJwtMiddleware, muc *mock_users_pb.MockUserClient) {
				ml.EXPECT().Info(gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
				ml.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mjm.EXPECT().Middleware(gomock.Any()).AnyTimes().Return(func(c *gin.Context) { c.Next() })
			},
			ExtraCheck: func(rr *httptest.ResponseRecorder, t *testing.T) {
				cookie := rr.Result().Cookies()
				token := false
				refresh := false
				for _, c := range cookie {
					if c.Name == middleware.TokenCookieName {
						if token == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Less(t, c.MaxAge, 0, "cookie must has MaxAge < 0 to be removed")
						assert.Equal(t, "", c.Value)
						token = true
					}
					if c.Name == middleware.RefreshCookieName {
						if refresh == true {
							assert.Fail(t, "found cookie duplicate: refresh token")
						}
						assert.Less(t, c.MaxAge, 0, "cookie must has MaxAge < 0 to be removed")
						assert.Equal(t, "", c.Value)
						refresh = true
					}
				}
				assert.True(t, token, "excepted token in cookies")
				assert.True(t, refresh, "excepted refresh token in cookies")
			},
			ExceptedStatus: http.StatusNoContent,
			ExceptedBody:   "",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			client := mock_users_pb.NewMockUserClient(ctrl)
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
			r.AddCookie(&http.Cookie{Name: middleware.RefreshCookieName, Value: "refresh cookie", MaxAge: 100})
			e.ServeHTTP(w, r)

			assert.Equal(t, v.ExceptedStatus, w.Result().StatusCode)
			b, _ := io.ReadAll(w.Result().Body)
			assert.Equal(t, v.ExceptedBody, string(b))

			if v.ExtraCheck != nil {
				v.ExtraCheck(w, t)
			}

			er := h.Close()
			assert.NoError(t, er)
		})
	}
}
