package service

import (
	"context"

	model "github.com/reversersed/go-grpc/tree/main/api_genre/internal/storage"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/proto/genres"
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
	GetAll(context.Context) ([]*model.Category, error)
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
	genres_pb.UnimplementedGenreServer
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
	genres_pb.RegisterGenreServer(s, u)
}