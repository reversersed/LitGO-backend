package service

import (
	"context"

	model "github.com/reversersed/go-grpc/tree/main/api_user/internal/storage"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto/users"
	"google.golang.org/grpc"
)

var dataValidationRules = []struct {
	Data  any
	Rules map[string]string
}{
	{
		Data: users_pb.RegistrationRequest{},
		Rules: map[string]string{
			"Login":          "required,min=4,max=16,onlyenglish",
			"Email":          "required,email",
			"Password":       "required,min=8,max=32,lowercase,uppercase,digitrequired,specialsymbol",
			"PasswordRepeat": "required,eqfield=password",
		},
	},
	{
		Data: users_pb.LoginRequest{},
		Rules: map[string]string{
			"Login":    "required",
			"Password": "required",
		},
	},
	{
		Data: users_pb.TokenRequest{},
		Rules: map[string]string{
			"RefreshToken": "required,jwt",
		},
	},
	{
		Data: users_pb.UserIdRequest{},
		Rules: map[string]string{
			"Id": "required,primitiveid",
		},
	},
}

type validator interface {
	StructValidation(any) error
	Register(any, map[string]string)
}
type logger interface {
	Infof(string, ...any)
	Info(...any)
	Errorf(string, ...any)
	Error(...any)
	Warnf(string, ...any)
	Warn(...any)
}
type storage interface {
	FindById(context.Context, string) (*model.User, error)
	FindByLogin(context.Context, string) (*model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
	CreateUser(ctx context.Context, model *model.User) (string, error)
}
type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type userServer struct {
	jwtSecret string
	cache     cache
	logger    logger
	storage   storage
	validator validator
	users_pb.UnimplementedUserServer
}

func NewServer(secret string, logger logger, cache cache, storage storage, validator validator) *userServer {
	for _, rule := range dataValidationRules {
		validator.Register(rule.Data, rule.Rules)
	}
	return &userServer{
		jwtSecret: secret,
		storage:   storage,
		logger:    logger,
		cache:     cache,
		validator: validator,
	}
}
func (u *userServer) Register(s grpc.ServiceRegistrar) {
	users_pb.RegisterUserServer(s, u)
}
