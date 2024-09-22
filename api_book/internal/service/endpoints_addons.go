package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	s.logger.Infof("received mapping book %v, genre zero: %v", src, src.Genre.IsZero())
	if !src.Genre.IsZero() {
		var response *genres_pb.CategoryResponse
		genre, err := s.cache.Get([]byte("category_" + src.Genre.Hex()))
		if jsonErr := json.Unmarshal(genre, &response); err != nil || jsonErr != nil {
			response, err = s.genreService.GetTree(ctx, &genres_pb.GetOneOfRequest{Query: src.Genre.Hex()})
			if err != nil {
				return nil, status.Error(codes.NotFound, "genre not found: "+err.Error())
			}
			responseJson, err := json.Marshal(response)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			if err := s.cache.Set([]byte("category_"+response.GetCategory().GetId()), responseJson, int(time.Hour*2)); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
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
		if len(book.GetGenre().Id) == 0 {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("genre %s not found", src.Genre.Hex()))
		}
		s.logger.Infof("mapped genres: %v", book.GetGenre())
	} else {
		book.Genre = nil
		book.Category = nil
	}
	s.logger.Infof("authors: %v (len=%d)", src.Authors, len(src.Authors))
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
			return nil, status.Error(codes.Internal, "authors not found: "+err.Error())
		}
		s.logger.Infof("mapped authors: %v", book.GetAuthors())
	} else {
		book.Authors = nil
	}

	return &book, nil
}
