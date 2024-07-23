package service

import (
	"context"
	"time"

	users_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto/users"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *userServer) GetUserById(c context.Context, r *users_pb.UserIdRequest) (*users_pb.UserModel, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	user, err := u.storage.FindById(c, r.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &users_pb.UserModel{
		Id:    user.Id.Hex(),
		Login: user.Login,
		Email: user.Email,
		Roles: user.Roles,
	}, nil
}
func (u *userServer) UpdateToken(c context.Context, r *users_pb.TokenRequest) (*users_pb.TokenReply, error) {
	token, refresh, err := u.UpdateRefreshToken(r.Refreshtoken)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &users_pb.TokenReply{
		Token:        token,
		Refreshtoken: refresh,
	}, nil
}
func (u *userServer) Login(c context.Context, r *users_pb.LoginRequest) (*users_pb.LoginResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	model, err := u.storage.FindByLogin(c, r.Login)
	if err != nil {
		model, err = u.storage.FindByEmail(c, r.Login)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid password or login")
		}
	}

	if err = bcrypt.CompareHashAndPassword(model.Password, []byte(r.Password)); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid password or login")
	}
	token, refresh, err := u.GenerateAccessToken(model)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &users_pb.LoginResponse{
		Login:        model.Login,
		Roles:        model.Roles,
		Token:        token,
		Refreshtoken: refresh,
	}, nil
}
