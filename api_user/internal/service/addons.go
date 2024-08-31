package service

import (
	"encoding/json"
	"time"

	"github.com/cristalhq/jwt/v3"
	model "github.com/reversersed/go-grpc/tree/main/api_user/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type claims struct {
	jwt.RegisteredClaims
	Login string
	Roles []string
	Email string
}

func (j *userServer) UpdateRefreshToken(refreshToken string) (string, string, error) {
	defer j.cache.Delete([]byte(refreshToken))

	userBytes, err := j.cache.Get([]byte(refreshToken))
	if err != nil {
		j.logger.Warn(err)
		return "", "", status.Errorf(codes.Unauthenticated, "refresh token not found")
	}
	var u model.User
	err = json.Unmarshal(userBytes, &u)
	if err != nil {
		j.logger.Error(err)
		return "", "", err
	}
	return j.GenerateAccessToken(&u)
}
func (j *userServer) GenerateAccessToken(u *model.User) (string, string, error) {
	signer, err := jwt.NewSignerHS(jwt.HS256, []byte(j.jwtSecret))
	if err != nil {
		j.logger.Warn(err)
		return "", "", err
	}
	builder := jwt.NewBuilder(signer)

	claims := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        u.Id.Hex(),
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
	if err := j.cache.Set([]byte(refreshTokenUuid), userBytes, int((7*24*time.Hour)/time.Second)); err != nil {
		return "", "", err
	}

	return token.String(), refreshTokenUuid, nil
}
