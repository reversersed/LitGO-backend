package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	model "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage"
	"github.com/reversersed/LitGO-backend/tree/main/api_book/pkg/copier"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
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

	s.logger.Infof("received book find request %s, built pattern: %s", req.GetQuery(), pattern)
	response, err := s.storage.Find(ctx, pattern, int(req.GetLimit()), int(req.GetPage()))
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
