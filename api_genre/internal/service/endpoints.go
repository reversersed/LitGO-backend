package service

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/proto/genres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *genreServer) GetAll(ctx context.Context, _ *genres_pb.Empty) (*genres_pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var categories []*genres_pb.CategoryModel
	if err := copier.Copy(&categories, &response); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(categories) == 0 {
		return nil, status.Error(codes.NotFound, "there is no genres in database")
	}
	return &genres_pb.GetAllResponse{
		Categories: categories,
	}, nil
}
