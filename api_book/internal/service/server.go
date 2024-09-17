package service

import (
	"context"

	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	model "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
<<<<<<< HEAD
	books_pb "github.com/reversersed/go-grpc/tree/main/api_book/pkg/proto/books"
=======
>>>>>>> 560fcec (separated projects and proto files)
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
	GetSuggestions(ctx context.Context, regex string, limit int64) ([]*model.Book, error)
}

type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type bookServer struct {
	cache     cache
	logger    logger
	storage   storage
	validator validator
	books_pb.UnimplementedBookServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator) *bookServer {
	return &bookServer{
		storage:   storage,
		logger:    logger,
		cache:     cache,
		validator: validator,
	}
}
func (u *bookServer) Register(s grpc.ServiceRegistrar) {
	books_pb.RegisterBookServer(s, u)
}
