package service

import (
	"context"

	model "github.com/reversersed/go-grpc/tree/main/api_user/internal/storage"
	users_pb "github.com/reversersed/go-grpc/tree/main/api_user/pkg/proto/users"
	"google.golang.org/grpc"
)

type logger interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
}
type storage interface {
	FindById(context.Context, string) (*model.User, error)
	FindByLogin(context.Context, string) (*model.User, error)
	FindByEmail(context.Context, string) (*model.User, error)
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
	users_pb.UnimplementedUserServer
}

func NewServer(secret string, logger logger, cache cache, storage storage) *userServer {
	return &userServer{
		jwtSecret: secret,
		storage:   storage,
		logger:    logger,
		cache:     cache,
	}
}
func (u *userServer) Register(s grpc.ServiceRegistrar) {
	users_pb.RegisterUserServer(s, u)
}
