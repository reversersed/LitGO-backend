package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/reversersed/LitGO-backend-pkg/copier"
	model "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *bookServer) FindBook(ctx context.Context, req *books_pb.FindBookRequest) (*books_pb.FindBookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}
	var pattern string

	words := strings.Fields(req.GetQuery())
	for _, word := range words {
		pattern += fmt.Sprintf("(%s)|", regexp.QuoteMeta(word))
	}
	pattern = strings.Trim(pattern, "|")

	if len(req.GetQuery()) == 0 {
		pattern = "(.*?)"
	}

	s.logger.Infof("received book find request %s, built pattern: %s", req.GetQuery(), pattern)
	response, err := s.storage.Find(ctx, pattern, int(req.GetLimit()), int(req.GetPage()), req.GetRating(), model.SortType(req.GetSorttype()))
	if err != nil {
		return nil, err
	}
	s.logger.Infof("got %d books by pattern %s, limit %d, page %d", len(response), pattern, req.GetLimit(), req.GetPage())
	data := make([]*books_pb.BookModel, len(response))
	for i, v := range response {
		model, err := s.bookMapper(ctx, v)
		if err != nil {
			return nil, err
		}
		data[i] = model
	}

	return &books_pb.FindBookResponse{Books: data}, nil
}

// TODO rework method to new request
func (s *bookServer) CreateBook(ctx context.Context, req *books_pb.CreateBookRequest) (*books_pb.CreateBookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2500*time.Millisecond)
	defer cancel()

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}

	var book model.Book
	if err := copier.Copy(&book, req, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.logger.Infof("received create book request: %v", book)
	if _, err := s.bookMapper(ctx, &book); err != nil {
		return nil, err
	}

	response, err := s.storage.CreateBook(ctx, &book)
	if err != nil {
		return nil, err
	}
	responseModel, err := s.bookMapper(ctx, response)
	if err != nil {
		return nil, err
	}
	s.logger.Infof("created book mapped to: %v", responseModel)
	err = s.rabbitService.SendBookCreatedMessage(ctx, responseModel)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &books_pb.CreateBookResponse{Book: responseModel}, nil
}
func (s *bookServer) GetBook(ctx context.Context, req *books_pb.GetBookRequest) (*books_pb.GetBookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}
	book := new(model.Book)
	if bytes, err := s.cache.Get([]byte("book_" + req.GetQuery())); err == nil {
		if err = json.Unmarshal(bytes, book); err != nil {
			return nil, status.Error(codes.Internal, "error decoding book: "+err.Error())
		}
	} else {
		book, err = s.storage.GetBook(ctx, req.GetQuery())
		if err != nil {
			return nil, err
		}
		bytes, err := json.Marshal(book)
		if err != nil {
			return nil, status.Error(codes.Internal, "error encoding book: "+err.Error())
		}
		if err = s.cache.Set([]byte("book_"+req.GetQuery()), bytes, int(time.Hour)); err != nil {
			return nil, status.Error(codes.Internal, "error saving book to cache: "+err.Error())
		}
	}
	responseModel, err := s.bookMapper(ctx, book)
	if err != nil {
		return nil, err
	}
	return &books_pb.GetBookResponse{Book: responseModel}, nil
}

func (s *bookServer) GetBookByGenre(ctx context.Context, req *books_pb.GetBookByGenreRequest) (*books_pb.GetBookByGenreResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}

	tree, err := s.genreService.GetTree(ctx, &genres_pb.GetOneOfRequest{Query: req.GetQuery()})
	if err != nil {
		return nil, err
	}
	genre := make([]primitive.ObjectID, 0)

	for _, v := range tree.GetCategory().GetGenres() {
		if v.GetId() == req.GetQuery() || v.GetTranslitname() == req.GetQuery() || tree.GetCategory().GetId() == req.GetQuery() || tree.GetCategory().GetTranslitname() == req.GetQuery() {
			id, err := primitive.ObjectIDFromHex(v.GetId())
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			genre = append(genre, id)
			if tree.GetCategory().GetId() != req.GetQuery() && tree.GetCategory().GetTranslitname() != req.GetQuery() {
				break
			}
		}
	}

	response, err := s.storage.GetBookByGenre(ctx, genre, model.SortType(req.GetSorttype()), req.GetOnlyhighrating(), int(req.GetLimit()), int(req.GetPage()))
	if err != nil {
		return nil, err
	}

	data := make([]*books_pb.BookModel, len(response))
	for i, v := range response {
		model, err := s.bookMapper(ctx, v)
		if err != nil {
			return nil, err
		}
		data[i] = model
	}
	return &books_pb.GetBookByGenreResponse{Books: data}, nil
}

// TODO write tests
func (s *bookServer) GetBookList(ctx context.Context, req *books_pb.GetBookListRequest) (*books_pb.GetBookListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}

	var ids []primitive.ObjectID
	if err := copier.Copy(&ids, req.GetId(), copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response, err := s.storage.GetBookList(ctx, ids, req.GetTranslit())
	if err != nil {
		return nil, err
	}

	data := make([]*books_pb.BookModel, len(response))
	for i, v := range response {
		model, err := s.bookMapper(ctx, v)
		if err != nil {
			return nil, err
		}
		data[i] = model
	}
	return &books_pb.GetBookListResponse{Books: data}, nil
}
