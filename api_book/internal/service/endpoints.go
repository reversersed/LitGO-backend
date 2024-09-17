package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/copier"
<<<<<<< HEAD
	books_pb "github.com/reversersed/go-grpc/tree/main/api_book/pkg/proto/books"
=======
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
>>>>>>> 560fcec (separated projects and proto files)
)

func (s *bookServer) GetBookSuggestions(ctx context.Context, req *books_pb.GetSuggestionRequest) (*books_pb.GetBooksResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

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
