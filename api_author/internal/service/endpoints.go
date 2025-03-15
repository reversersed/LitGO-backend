package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/reversersed/LitGO-backend-pkg/copier"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *authorServer) GetAuthors(ctx context.Context, r *authors_pb.GetAuthorsRequest) (*authors_pb.GetAuthorsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if r == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}
	s.logger.Infof("received get authors with request id=%v, name=%v", r.GetId(), r.GetTranslit())
	var id []primitive.ObjectID
	for _, i := range r.GetId() {
		v, err := primitive.ObjectIDFromHex(i)
		if err != nil {
			status, _ := status.New(codes.InvalidArgument, "unable to convert value to type").WithDetails(&shared_pb.ErrorDetail{
				Field:       "id",
				Struct:      "authors_pb.GetAuthorsRequest",
				Description: fmt.Sprintf("unable to convert %s to primitive id type", i),
				Actualvalue: strings.Join(r.GetId(), ","),
			})
			return nil, status.Err()
		}
		id = append(id, v)
	}
	authors, err := s.storage.GetAuthors(ctx, id, r.GetTranslit())
	if err != nil {
		s.logger.Errorf("got error from storage: %v, id data: %v, translit data: %v", err, id, r.GetTranslit())
		return nil, err
	}
	s.logger.Infof("found %d authors: %v", len(authors), authors)
	authorModels := make([]*authors_pb.AuthorModel, len(authors))
	if err := copier.Copy(&authorModels, &authors, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authors_pb.GetAuthorsResponse{
		Authors: authorModels,
	}, nil
}

func (s *authorServer) FindAuthors(ctx context.Context, r *authors_pb.FindAuthorsRequest) (*authors_pb.GetAuthorsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if r == nil {
		return nil, status.Error(codes.InvalidArgument, "received nil request")
	}
	if err := s.validator.StructValidation(r); err != nil {
		return nil, err
	}
	var pattern string

	words := strings.Fields(r.GetQuery())
	for _, word := range words {
		pattern += fmt.Sprintf("(%s)|", regexp.QuoteMeta(word))
	}
	pattern = strings.Trim(pattern, "|")

	s.logger.Infof("received authors suggestion request %s with pattern %s", r.GetQuery(), pattern)
	authors, err := s.storage.Find(ctx, pattern, int(r.GetLimit()), int(r.GetPage()), r.GetRating())
	if err != nil {
		return nil, err
	}
	authorModels := make([]*authors_pb.AuthorModel, len(authors))
	if err := copier.Copy(&authorModels, &authors, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authors_pb.GetAuthorsResponse{
		Authors: authorModels,
	}, nil
}
