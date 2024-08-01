package service

import (
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_author/pkg/proto/authors"
	"google.golang.org/grpc"
)

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
}
type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type genreServer struct {
	cache     cache
	logger    logger
	storage   storage
	validator validator
	authors_pb.UnimplementedGenreServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator) *genreServer {
	return &genreServer{
		storage:   storage,
		logger:    logger,
		cache:     cache,
		validator: validator,
	}
}
func (u *genreServer) Register(s grpc.ServiceRegistrar) {
	authors_pb.RegisterGenreServer(s, u)
}
