package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
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
	response, err := s.storage.GetSuggestions(ctx, pattern, req.GetLimit())
	if err != nil {
		return nil, err
	}
	data := make([]*books_pb.BookModel, len(response))
	if err := copier.Copy(&data, &response, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, err
	}

	return &books_pb.GetBooksResponse{Books: data}, nil
}
