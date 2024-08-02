package service

import (
	"context"

	model "github.com/reversersed/go-grpc/tree/main/api_author/internal/storage"
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_author/pkg/proto/authors"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	GetAuthors(ctx context.Context, id []primitive.ObjectID, translit []string) ([]*model.Author, error)
}
type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type authorServer struct {
	cache     cache
	logger    logger
	storage   storage
	validator validator
	authors_pb.UnimplementedAuthorServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator) *authorServer {
	return &authorServer{
		storage:   storage,
		logger:    logger,
		cache:     cache,
		validator: validator,
	}
}
func (u *authorServer) Register(s grpc.ServiceRegistrar) {
	authors_pb.RegisterAuthorServer(s, u)
}
