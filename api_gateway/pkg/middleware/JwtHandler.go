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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Key string

const (
	tokenCookieName   string = "authTokenCookie"
	refreshCookieName string = "refreshTokenCookie"
	UserIdKey         Key    = "userAuthId"
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

func (j *jwtMiddleware) Middleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		headertoken, err := c.Cookie(tokenCookieName)
		if err != nil {
			c.Error(status.Error(codes.Unauthenticated, "user has no token cookie"))
			return
		}
		key := []byte(j.secret)
		verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
		if err != nil {
			j.logger.Errorf("error creating verifier for key. key length = %d, error = %v", len(key), err)
			c.Error(status.Errorf(codes.Unauthenticated, "error creating verifier for key"))
			return
		}
		j.logger.Info("parsing and verifying token...")
		token, err := jwt.ParseAndVerifyString(headertoken, verifier)
		if err != nil {
			c.Error(status.Errorf(codes.Unauthenticated, err.Error()))
			return
		}

		var claims claims
		if err := json.Unmarshal(token.RawClaims(), &claims); err != nil {
			c.Error(status.Errorf(codes.Unauthenticated, err.Error()))
			return
		}
		if !claims.IsValidAt(time.Now()) {
			refreshCookie, err := c.Cookie(refreshCookieName)
			if err != nil {
				c.SetCookie(tokenCookieName, "", -1, "/", "/", true, true)
				c.SetCookie(refreshCookieName, "", -1, "/", "/", true, true)
				c.Error(status.Errorf(codes.Unauthenticated, err.Error()))
				return
			}
			token, refresh, err := j.UpdateRefreshToken(refreshCookie)
			if err != nil {
				c.SetCookie(tokenCookieName, "", -1, "/", "/", true, true)
				c.SetCookie(refreshCookieName, "", -1, "/", "/", true, true)
				c.Error(status.Errorf(codes.Unauthenticated, err.Error()))
				return
			}
			c.SetCookie(tokenCookieName, token, (int)((31*24*time.Hour)/time.Second), "/", "/", true, true)
			c.SetCookie(refreshCookieName, refresh, (int)((31*24*time.Hour)/time.Second), "/", "/", true, true)
		}
		if len(roles) > 0 {
			var errorRoles []string
			for _, role := range roles {
				if len(role) > 0 && !slices.Contains(claims.Roles, role) {
					errorRoles = append(errorRoles, fmt.Sprintf("user has no %s right", role))
				}
			}
			if len(errorRoles) > 0 {
				c.Error(status.Errorf(codes.PermissionDenied, strings.Join(errorRoles, ", ")))
				return
			}
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), UserIdKey, claims.ID))
		c.Next()
	}
}
func (j *jwtMiddleware) UpdateRefreshToken(refreshToken string) (string, string, error) {
	defer j.cache.Delete([]byte(refreshToken))

	userBytes, err := j.cache.Get([]byte(refreshToken))
	if err != nil {
		j.logger.Warn(err)
		return "", "", status.Errorf(codes.NotFound, "refresh token not found")
	}
	var u UserTokenModel
	err = json.Unmarshal(userBytes, &u)
	if err != nil {
		j.logger.Error(err)
		return "", "", err
	}
	return j.GenerateAccessToken(&u)
}
func (j *jwtMiddleware) GenerateAccessToken(u *UserTokenModel) (string, string, error) {
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(j.secret))
	if err != nil {
		j.logger.Warn(err)
		return "", "", err
	}
	builder := jwt.NewBuilder(signer)

	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        u.Id,
			Audience:  u.Roles,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		Roles: u.Roles,
		Login: u.Login,
		Email: u.Email,
	}
	token, err := builder.Build(claims)
	if err != nil {
		j.logger.Warn(err)
		return "", "", err
	}

	j.logger.Info("creating refresh token...")
	refreshTokenUuid := primitive.NewObjectID().Hex()
	userBytes, _ := json.Marshal(u)
	j.cache.Set([]byte(refreshTokenUuid), userBytes, int((7*24*time.Hour)/time.Second))

	return token.String(), refreshTokenUuid, nil
}
func (j *jwtMiddleware) GetUserClaims(token string) (*UserTokenModel, error) {
	verifier, err := jwt.NewVerifierHS(jwt.HS256, []byte(j.secret))
	if err != nil {
		return nil, err
	}

	claimToken, err := jwt.ParseAndVerifyString(token, verifier)
	if err != nil {
		return nil, err
	}

	var claims claims
	if err := json.Unmarshal(claimToken.RawClaims(), &claims); err != nil {
		return nil, err
	}
	j.logger.Infof("user %s authorized with %v rights", claims.Login, claims.Roles)
	return &UserTokenModel{
		Login: claims.Login,
		Roles: claims.Roles,
	}, nil
}
