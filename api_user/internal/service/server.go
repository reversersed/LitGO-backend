package service

import (
	"context"

	model "github.com/reversersed/go-grpc/tree/main/api_user/internal/storage"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto/users"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

//go:generate mockgen -source=server.go -destination=mocks/server.go

type validator interface {
	StructValidation(any) error
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
	CreateUser(ctx context.Context, model *model.User) (primitive.ObjectID, error)
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
