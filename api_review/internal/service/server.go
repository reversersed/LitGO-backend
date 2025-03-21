package service

import (
	reviews_pb "github.com/reversersed/LitGO-proto/gen/go/reviews"
	users_pb "github.com/reversersed/LitGO-proto/gen/go/users"
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
type storage any

type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}

type reviewServer struct {
	cache       cache
	logger      logger
	storage     storage
	validator   validator
	userService users_pb.UserClient
	reviews_pb.UnimplementedReviewServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator, userService users_pb.UserClient) *reviewServer {
	return &reviewServer{
		storage:     storage,
		logger:      logger,
		cache:       cache,
		validator:   validator,
		userService: userService,
	}
}
func (u *reviewServer) Register(s grpc.ServiceRegistrar) {
	reviews_pb.RegisterReviewServer(s, u)
}
