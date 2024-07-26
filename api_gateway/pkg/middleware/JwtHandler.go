package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/gin-gonic/gin"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_gateway/pkg/proto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=JwtHandler.go -destination=mocks/jwt_mw_mock.go

const (
	TokenCookieName     string = "authTokenCookie"
	RefreshCookieName   string = "refreshTokenCookie"
	UserIdKey           string = "userAuthId"
	UserCredentialLogin string = "userLoginCredential"
	UserCredentialRoles string = "userRolesCredential"
)

type Logger interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
}
type Cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type UserServer interface {
	UpdateToken(context.Context, *users_pb.TokenRequest, ...grpc.CallOption) (*users_pb.TokenReply, error)
}
type jwtMiddleware struct {
	secret string
	logger Logger
	cache  Cache
}
type claims struct {
	jwt.RegisteredClaims
	Login string   `json:"login"`
	Roles []string `json:"roles"`
	Email string   `json:"email"`
}
type UserTokenModel struct {
	Id    string   `json:"-"`
	Login string   `json:"login"`
	Roles []string `json:"roles"`
	Email string   `json:"-"`
}

func NewJwtMiddleware(logger Logger, cache Cache, secret string) *jwtMiddleware {
	return &jwtMiddleware{
		secret: secret,
		logger: logger,
		cache:  cache,
	}
}

func (j *jwtMiddleware) Middleware(server UserServer, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headertoken, err := c.Cookie(TokenCookieName)
		if err != nil {
			c.Error(status.Error(codes.Unauthenticated, "user has no token cookie"))
			c.Abort()
			return
		}
		key := []byte(j.secret)
		verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
		if err != nil {
			j.logger.Errorf("error creating verifier for key. key length = %d, error = %v", len(key), err)
			c.Error(status.Error(codes.Unauthenticated, "error creating verifier for key"))
			c.Abort()
			return
		}
		j.logger.Info("parsing and verifying token...")
		token, err := jwt.ParseAndVerifyString(headertoken, verifier)
		if err != nil {
			c.Error(status.Error(codes.Unauthenticated, err.Error()))
			c.Abort()
			return
		}

		var claims claims
		if err := json.Unmarshal(token.RawClaims(), &claims); err != nil {
			c.Error(status.Errorf(codes.Unauthenticated, err.Error()))
			c.Abort()
			return
		}
		if !claims.IsValidAt(time.Now()) {
			refreshCookie, err := c.Cookie(RefreshCookieName)
			if err != nil {
				c.SetCookie(TokenCookieName, "", -1, "/", "", true, true)
				c.SetCookie(RefreshCookieName, "", -1, "/", "", true, true)
				c.Error(status.Error(codes.Unauthenticated, err.Error()))
				c.Abort()
				return
			}
			tokenReply, err := server.UpdateToken(c.Request.Context(), &users_pb.TokenRequest{Refreshtoken: refreshCookie})
			if err != nil {
				c.SetCookie(TokenCookieName, "", -1, "/", "", true, true)
				c.SetCookie(RefreshCookieName, "", -1, "/", "", true, true)
				c.Error(err)
				c.Abort()
				return
			}
			c.SetCookie(TokenCookieName, tokenReply.Token, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
			c.SetCookie(RefreshCookieName, tokenReply.Refreshtoken, (int)((31*24*time.Hour)/time.Second), "/", "", true, true)
		}
		if len(roles) > 0 {
			var errorRoles []string
			for _, role := range roles {
				if len(role) > 0 && !slices.Contains(claims.Roles, role) {
					errorRoles = append(errorRoles, fmt.Sprintf("user has no %s right", role))
				}
			}
			if len(errorRoles) > 0 {
				c.Error(status.Error(codes.PermissionDenied, strings.Join(errorRoles, ", ")))
				c.Abort()
				return
			}
		}
		j.logger.Infof("user's %s token has been verified with %v rights", claims.Login, claims.Roles)
		c.Set(UserIdKey, claims.ID)
		c.Set(UserCredentialLogin, claims.Login)
		c.Set(UserCredentialRoles, claims.Roles)
		c.Next()
	}
}
