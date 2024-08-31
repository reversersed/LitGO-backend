package service

import (
	"context"
	"encoding/json"
	"time"

	model "github.com/reversersed/go-grpc/tree/main/api_genre/internal/storage"
	"github.com/reversersed/go-grpc/tree/main/api_genre/pkg/copier"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/proto/genres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *genreServer) GetAll(ctx context.Context, _ *genres_pb.Empty) (*genres_pb.GetAllResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var response []*model.Category

	if cats, err := s.cache.Get([]byte("all_categories")); len(cats) > 0 && err == nil {
		if err := json.Unmarshal(cats, &response); err != nil {
			return nil, err
		}
	} else {
		response, err = s.storage.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		bytes, _ := json.Marshal(&response)
		if err := s.cache.Set([]byte("all_categories"), bytes, int(time.Hour*6)); err != nil {
			return nil, err
		}
	}

	var categories []*genres_pb.CategoryModel
	if err := copier.Copy(&categories, &response, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &genres_pb.GetAllResponse{
		Categories: categories,
	}, nil
}
