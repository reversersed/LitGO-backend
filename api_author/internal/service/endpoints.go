package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/reversersed/go-grpc/tree/main/api_author/pkg/copier"
	shared_pb "github.com/reversersed/go-grpc/tree/main/api_author/pkg/proto"
	authors_pb "github.com/reversersed/go-grpc/tree/main/api_author/pkg/proto/authors"
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
	var id []primitive.ObjectID
	for _, i := range r.Id {
		v, err := primitive.ObjectIDFromHex(i)
		if err != nil {
			status, _ := status.New(codes.InvalidArgument, "unable to convert value to type").WithDetails(&shared_pb.ErrorDetail{
				Field:       "id",
				Struct:      "authors_pb.GetAuthorsRequest",
				Description: fmt.Sprintf("unable to convert %s to primitive id type", i),
				Actualvalue: strings.Join(r.Id, ","),
			})
			return nil, status.Err()
		}
		id = append(id, v)
	}
	authors, err := s.storage.GetAuthors(ctx, id, r.Translit)
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
