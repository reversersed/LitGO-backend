package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	model "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *bookServer) GetBookSuggestions(ctx context.Context, req *books_pb.GetSuggestionRequest) (*books_pb.GetBooksResponse, error) {
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

	s.logger.Infof("received book suggestion request %s, built pattern: %s", req.GetQuery(), pattern)
	response, err := s.storage.GetSuggestions(ctx, pattern, req.GetLimit())
	if err != nil {
		return nil, err
	}
	s.logger.Infof("got %d books by pattern %s, limit it %d", len(response), pattern, req.GetLimit())
	data := make([]*books_pb.BookModel, len(response))
	for i, v := range response {
		model, err := s.bookMapper(ctx, v)
		if err != nil {
			return nil, err
		}
		data[i] = model
	}

	return &books_pb.GetBooksResponse{Books: data}, nil
}

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
	response, err := s.storage.CreateBook(ctx, &book)
	if err != nil {
		return nil, err
	}
	responseModel, err := s.bookMapper(ctx, response)
	if err != nil {
		return nil, err
	}
	s.logger.Infof("created book mapped to: %v", responseModel)
	return &books_pb.CreateBookResponse{Book: responseModel}, nil
}
