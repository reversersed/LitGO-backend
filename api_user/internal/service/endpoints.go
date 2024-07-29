package service

import (
	"context"
	"time"

	model "github.com/reversersed/go-grpc/tree/main/api_user/internal/storage"
	shared_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func (u *userServer) GetUser(c context.Context, r *users_pb.UserRequest) (*users_pb.UserModel, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := u.validator.StructValidation(r); err != nil {
		return nil, err
	}
	var user *model.User
	var err error
	details := make([]protoadapt.MessageV1, 0)

	if len(r.Id) > 0 {
		user, err = u.storage.FindById(c, r.Id)
		if err != nil {
			user = nil
			details = append(details, &shared_pb.ErrorDetail{Field: "Id", Struct: "users_pb.UserRequest", Description: err.Error(), Actualvalue: r.Id})
		}
	}
	if len(r.Login) > 0 && user == nil {
		user, err = u.storage.FindByLogin(c, r.Login)
		if err != nil {
			user = nil
			details = append(details, &shared_pb.ErrorDetail{Field: "Login", Struct: "users_pb.UserRequest", Description: err.Error(), Actualvalue: r.Login})
		}
	}
	if len(r.Email) > 0 && user == nil {
		user, err = u.storage.FindByEmail(c, r.Email)
		if err != nil {
			user = nil
			details = append(details, &shared_pb.ErrorDetail{Field: "Email", Struct: "users_pb.UserRequest", Description: err.Error(), Actualvalue: r.Email})
		}
	}

	if user == nil {
		err, _ := status.New(codes.NotFound, "user does not exist").WithDetails(details...)
		return nil, err.Err()
	}
	return &users_pb.UserModel{
		Id:    user.Id.Hex(),
		Login: user.Login,
		Email: user.Email,
		Roles: user.Roles,
	}, nil
}
func (u *userServer) UpdateToken(c context.Context, r *users_pb.TokenRequest) (*users_pb.TokenReply, error) {
	if err := u.validator.StructValidation(r); err != nil {
		return nil, err
	}

	token, refresh, err := u.UpdateRefreshToken(r.Refreshtoken)
	if err != nil {
		return nil, err
	}
	return &users_pb.TokenReply{
		Token:        token,
		Refreshtoken: refresh,
	}, nil
}
func (u *userServer) Login(c context.Context, r *users_pb.LoginRequest) (*users_pb.LoginResponse, error) {
	c, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := u.validator.StructValidation(r); err != nil {
		return nil, err
	}

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
func (u *userServer) RegisterUser(c context.Context, usr *users_pb.RegistrationRequest) (*users_pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := u.validator.StructValidation(usr); err != nil {
		return nil, err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Login:          usr.Login,
		Password:       pass,
		Roles:          []string{"user"},
		Email:          usr.Email,
		EmailConfirmed: false,
	}
	cntx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_, err = u.storage.FindByLogin(cntx, user.Login)
	if err == nil {
		u.logger.Warnf("user %s couldn't register because of login collision", user.Login)
		return nil, status.Error(codes.AlreadyExists, "user with provided login already exist")
	}
	_, err = u.storage.FindByEmail(cntx, user.Email)
	if err == nil {
		u.logger.Warnf("user %s couldn't register because of email (%s) collision", user.Login, user.Email)
		return nil, status.Error(codes.AlreadyExists, "email already taken")
	}
	result, err := u.storage.CreateUser(cntx, &user)
	if err != nil {
		u.logger.Errorf("couldn't add user %s to database: %v", user.Login, err)
		return nil, err
	}
	user.Id, err = primitive.ObjectIDFromHex(result)
	if err != nil {
		u.logger.Errorf("can't create user %s id: %v", user.Login, err)
		return nil, status.Error(codes.Internal, "can't create user id")
	}
	token, refresh, err := u.GenerateAccessToken(&user)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &users_pb.LoginResponse{
		Login:        user.Login,
		Roles:        user.Roles,
		Token:        token,
		Refreshtoken: refresh,
	}, nil
}
