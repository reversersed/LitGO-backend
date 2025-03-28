package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	mock_service "github.com/reversersed/LitGO-backend/tree/main/api_user/internal/service/mocks"
	model "github.com/reversersed/LitGO-backend/tree/main/api_user/internal/storage"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetUser(t *testing.T) {
	user := &model.User{Id: primitive.NewObjectID(), Login: "login", Email: "email", Roles: []string{"user"}}
	table := []struct {
		Name          string
		MockBehaviour func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		Request       *users_pb.UserRequest
		Response      *users_pb.UserModel
		ExceptedError string
	}{
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong request",
		},
		{
			Name: "error from every method",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindById(gomock.Any(), "id").Return(nil, status.Error(codes.NotFound, "user does not exist"))
				m3.EXPECT().FindByLogin(gomock.Any(), "login").Return(nil, status.Error(codes.NotFound, "user does not exist"))
				m3.EXPECT().FindByEmail(gomock.Any(), "email").Return(nil, status.Error(codes.NotFound, "user does not exist"))
			},
			Request: &users_pb.UserRequest{
				Id:    "id",
				Login: "login",
				Email: "email",
			},
			ExceptedError: "rpc error: code = NotFound desc = user does not exist",
		},
		{
			Name: "successful from id",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindById(gomock.Any(), "id").Return(user, nil)
			},
			Request: &users_pb.UserRequest{
				Id:    "id",
				Login: "login",
				Email: "email",
			},
			Response: &users_pb.UserModel{
				Id:    user.Id.Hex(),
				Login: "login",
				Email: "email",
				Roles: []string{"user"},
			},
		},
		{
			Name: "successful from id",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindById(gomock.Any(), "id").Return(nil, status.Error(codes.NotFound, "user does not exist"))
				m3.EXPECT().FindByLogin(gomock.Any(), "login").Return(user, nil)
			},
			Request: &users_pb.UserRequest{
				Id:    "id",
				Login: "login",
				Email: "email",
			},
			Response: &users_pb.UserModel{
				Id:    user.Id.Hex(),
				Login: "login",
				Email: "email",
				Roles: []string{"user"},
			},
		},
		{
			Name: "successful from id",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindById(gomock.Any(), "id").Return(nil, status.Error(codes.NotFound, "user does not exist"))
				m3.EXPECT().FindByLogin(gomock.Any(), "login").Return(nil, status.Error(codes.NotFound, "user does not exist"))
				m3.EXPECT().FindByEmail(gomock.Any(), "email").Return(user, nil)
			},
			Request: &users_pb.UserRequest{
				Id:    "id",
				Login: "login",
				Email: "email",
			},
			Response: &users_pb.UserModel{
				Id:    user.Id.Hex(),
				Login: "login",
				Email: "email",
				Roles: []string{"user"},
			},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}
			server := NewServer("secretKey", logger, cache, storage, validator)
			response, err := server.GetUser(context.Background(), v.Request)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.Response, response)
		})
	}
}
func TestLogin(t *testing.T) {
	pass, _ := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.MinCost)
	table := []struct {
		Name          string
		MockBehaviour func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		Request       *users_pb.LoginRequest
		Response      *users_pb.LoginResponse
		ExceptedError string
	}{
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong request",
		},
		{
			Name: "services error",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(nil, errors.New(""))
				m3.EXPECT().FindByEmail(gomock.Any(), "user").Return(nil, errors.New(""))
			},
			Request:       &users_pb.LoginRequest{Login: "user", Password: "123"},
			ExceptedError: "rpc error: code = InvalidArgument desc = invalid password or login",
		},
		{
			Name: "invalid password",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(&model.User{Login: "user", Password: []byte("321")}, nil)
			},
			Request:       &users_pb.LoginRequest{Login: "user", Password: "123"},
			ExceptedError: "rpc error: code = InvalidArgument desc = invalid password or login",
		},
		{
			Name: "success",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(&model.User{Login: "user", Password: pass, Roles: []string{"user"}}, nil)
				m2.EXPECT().Info(gomock.Any())
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any())
			},
			Request:  &users_pb.LoginRequest{Login: "user", Password: "123"},
			Response: &users_pb.LoginResponse{Login: "user", Roles: []string{"user"}},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}
			server := NewServer("secretKey", logger, cache, storage, validator)
			response, err := server.Login(grpc.NewContextWithServerTransportStream(context.Background(), &mock_service.MockServerTransportStream{}), v.Request)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			if v.Response != nil {
				if assert.NotNil(t, response) {
					assert.Equal(t, v.Response.Login, response.Login)
					assert.Equal(t, v.Response.Roles, response.Roles)
					assert.NotEmpty(t, response.Token)
					assert.NotEmpty(t, response.Refreshtoken)
				}
			} else {
				assert.Nil(t, response)
			}
		})
	}
}
func TestUpdateToken(t *testing.T) {
	table := []struct {
		Name          string
		MockBehaviour func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		Request       *users_pb.TokenRequest
		Response      bool
		ExceptedError string
	}{
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong request",
		},
		{
			Name: "empty cache",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m2.EXPECT().Warn(gomock.Any())
				m1.EXPECT().Delete([]byte("request"))
				m1.EXPECT().Get([]byte("request")).Return([]byte{}, errors.New(""))
			},
			Request:       &users_pb.TokenRequest{Refreshtoken: "request"},
			ExceptedError: "rpc error: code = Unauthenticated desc = refresh token not found",
		},
		{
			Name: "successful",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m2.EXPECT().Info(gomock.Any())
				m1.EXPECT().Delete([]byte("request"))
				user := &model.User{Id: primitive.NewObjectID(), Login: "user", Roles: []string{"user"}, Email: "email"}
				bytes, _ := json.Marshal(user)
				m1.EXPECT().Get([]byte("request")).Return(bytes, nil)
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any())
			},
			Request:  &users_pb.TokenRequest{Refreshtoken: "request"},
			Response: true,
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}
			server := NewServer("secretKey", logger, cache, storage, validator)
			response, err := server.UpdateToken(context.Background(), v.Request)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			if v.Response && assert.NotNil(t, response) {
				assert.NotEmpty(t, response.Refreshtoken)
				assert.NotEmpty(t, response.Token)
				assert.NotEqual(t, response.Refreshtoken, response.Token)
			}
		})
	}
}
func TestRegisterUser(t *testing.T) {
	table := []struct {
		Name          string
		MockBehaviour func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		Request       *users_pb.RegistrationRequest
		Response      *users_pb.LoginResponse
		ExceptedError string
	}{
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong request"))
			},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong request",
		},
		{
			Name: "login collision",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(nil, nil)
				m2.EXPECT().Warnf(gomock.Any(), "user")
			},
			Request:       &users_pb.RegistrationRequest{Login: "user"},
			ExceptedError: "rpc error: code = AlreadyExists desc = user with provided login already exist",
		},
		{
			Name: "email collision",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(nil, errors.New(""))
				m3.EXPECT().FindByEmail(gomock.Any(), "email").Return(nil, nil)
				m2.EXPECT().Warnf(gomock.Any(), "user", "email")
			},
			Request:       &users_pb.RegistrationRequest{Login: "user", Email: "email"},
			ExceptedError: "rpc error: code = AlreadyExists desc = email already taken",
		},
		{
			Name: "registration error",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(nil, errors.New(""))
				m3.EXPECT().FindByEmail(gomock.Any(), "email").Return(nil, errors.New(""))
				m3.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(primitive.ObjectID{}, status.Error(codes.Internal, "db not connected"))
				m2.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any())
			},
			Request:       &users_pb.RegistrationRequest{Login: "user", Email: "email"},
			ExceptedError: "rpc error: code = Internal desc = db not connected",
		},
		{
			Name: "successful registration",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().FindByLogin(gomock.Any(), "user").Return(nil, errors.New(""))
				m3.EXPECT().FindByEmail(gomock.Any(), "email").Return(nil, errors.New(""))
				m3.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(primitive.NewObjectID(), nil)
				m2.EXPECT().Info(gomock.Any())
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any())
			},
			Request:  &users_pb.RegistrationRequest{Login: "user", Email: "email"},
			Response: &users_pb.LoginResponse{Login: "user", Roles: []string{"user"}},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}
			server := NewServer("secretKey", logger, cache, storage, validator)
			response, err := server.RegisterUser(grpc.NewContextWithServerTransportStream(context.Background(), &mock_service.MockServerTransportStream{}), v.Request)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			if v.Response != nil && assert.NotNil(t, response) {
				assert.Equal(t, v.Response.Login, response.Login)
				assert.Equal(t, v.Response.Roles, response.Roles)
				assert.NotEmpty(t, response.Refreshtoken)
				assert.NotEmpty(t, response.Token)
			}
		})
	}
}

/*Template
func TestGetUser(t *testing.T) {
	table := []struct {
		Name          string
		MockBehaviour func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		Request       *users_pb.
		Response      *users_pb.
		ExceptedError string
	}{}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)

			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}
			server := NewServer("secretKey", logger, cache, storage, validator)
			response, err := server.(context.Background(), v.Request)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.Response, response)
		})
	}
}
*/
