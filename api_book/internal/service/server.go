package service

import (
	"context"

	model "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
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
	Find(context.Context, string, int, int, float32) ([]*model.Book, error)
	CreateBook(context.Context, *model.Book) (*model.Book, error)
	GetBook(context.Context, string) (*model.Book, error)
}

type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type rabbitservice interface {
	SendBookCreatedMessage(context.Context, *books_pb.BookModel) error
}
type bookServer struct {
	cache         cache
	logger        logger
	storage       storage
	validator     validator
	authorService authors_pb.AuthorClient
	genreService  genres_pb.GenreClient
	rabbitService rabbitservice
	books_pb.UnimplementedBookServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator, genreService genres_pb.GenreClient, authorService authors_pb.AuthorClient, rabbit rabbitservice) *bookServer {
	return &bookServer{
		storage:       storage,
		logger:        logger,
		cache:         cache,
		validator:     validator,
		genreService:  genreService,
		authorService: authorService,
		rabbitService: rabbit,
	}
}
func (u *bookServer) Register(s grpc.ServiceRegistrar) {
	books_pb.RegisterBookServer(s, u)
}
