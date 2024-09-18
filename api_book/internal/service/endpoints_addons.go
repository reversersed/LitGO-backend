package service

import (
	"context"
	"fmt"

	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	model "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *bookServer) bookMapper(ctx context.Context, src *model.Book) (*books_pb.BookModel, error) {
	var book books_pb.BookModel
	if err := copier.Copy(&book, src, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !src.Genre.IsZero() {
		response, err := s.genreService.GetTree(ctx, &genres_pb.GetOneOfRequest{Query: src.Genre.Hex()})
		if err != nil {
			return nil, err
		}
		book.Category = new(books_pb.CategoryModel)
		if err := copier.Copy(book.GetCategory(), response.GetCategory(), copier.WithPrimitiveToStringConverter); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		book.Genre = new(books_pb.GenreModel)
		for _, v := range response.GetCategory().GetGenres() {
			if v.GetId() == src.Genre.Hex() {
				if err := copier.Copy(book.GetGenre(), v, copier.WithPrimitiveToStringConverter); err != nil {
					return nil, status.Error(codes.Internal, err.Error())
				}
				break
			}
		}
		if book.GetGenre() == nil {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("genre %s not found", src.Genre.Hex()))
		}
	} else {
		book.Genre = nil
		book.Category = nil
	}
	if len(src.Authors) > 0 {
		authorsId := make([]string, len(src.Authors))
		for i, v := range src.Authors {
			authorsId[i] = v.Hex()
		}
		authorResponse, err := s.authorService.GetAuthors(ctx, &authors_pb.GetAuthorsRequest{Id: authorsId})
		if err != nil {
			return nil, err
		}
		if err := copier.Copy(&book.Authors, authorResponse.GetAuthors(), copier.WithPrimitiveToStringConverter); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		book.Authors = nil
	}

	return &book, nil
}
