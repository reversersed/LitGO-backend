package service

import (
	"context"

	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	model "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
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
	GetSuggestions(context.Context, string, int64) ([]*model.Book, error)
	CreateBook(context.Context, *model.Book) (*model.Book, error)
}

type cache interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, int) error
	Delete([]byte) bool
}
type bookServer struct {
	cache         cache
	logger        logger
	storage       storage
	validator     validator
	authorService authors_pb.AuthorClient
	genreService  genres_pb.GenreClient
	books_pb.UnimplementedBookServer
}

func NewServer(logger logger, cache cache, storage storage, validator validator, genreService genres_pb.GenreClient, authorService authors_pb.AuthorClient) *bookServer {
	return &bookServer{
		storage:       storage,
		logger:        logger,
		cache:         cache,
		validator:     validator,
		genreService:  genreService,
		authorService: authorService,
	}
}
func (u *bookServer) Register(s grpc.ServiceRegistrar) {
	books_pb.RegisterBookServer(s, u)
}
