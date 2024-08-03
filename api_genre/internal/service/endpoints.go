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

	if cats, err := s.cache.Get([]byte("all_categories")); len(cats) > 0 && err != nil {
		json.Unmarshal(cats, &response)
	} else {
		response, err = s.storage.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		bytes, _ := json.Marshal(&response)
		s.cache.Set([]byte("all_categories"), bytes, int(time.Hour*6))
	}

	var categories []*genres_pb.CategoryModel
	if err := copier.Copy(&categories, &response, copier.WithPrimitiveToStringConverter); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(categories) == 0 {
		return nil, status.Error(codes.NotFound, "there is no genres in database")
	}
	return &genres_pb.GetAllResponse{
		Categories: categories,
	}, nil
}
