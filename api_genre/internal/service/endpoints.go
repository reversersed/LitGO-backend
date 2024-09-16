package service

import (
	"context"
	"encoding/json"
	"time"

	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	model "github.com/reversersed/go-grpc/tree/main/api_genre/internal/storage"
	"github.com/reversersed/go-grpc/tree/main/api_genre/pkg/copier"
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

		bytes, err := json.Marshal(&response)
		if err != nil {
			return nil, err
		}
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

func (s *genreServer) GetOneOf(ctx context.Context, req *genres_pb.GetOneOfRequest) (*genres_pb.GetCategoryOrGenreResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}

	if m, err := s.cache.Get([]byte("one_" + req.GetQuery())); len(m) > 0 && err == nil {
		var model genres_pb.CategoryModel
		err = json.Unmarshal(m, &model)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if model.GetGenres() == nil || len(model.GetGenres()) == 0 {
			var genreModel genres_pb.GenreModel
			err := json.Unmarshal(m, &genreModel)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: &genreModel}}, nil
		}
		return &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: &model}}, nil
	}

	cat, err := s.storage.FindCategoryTree(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}
	for _, v := range cat.Genres {
		if v.Id.Hex() == req.GetQuery() || v.TranslitName == req.GetQuery() {
			var model genres_pb.GenreModel
			err = copier.Copy(&model, &v, copier.WithPrimitiveToStringConverter)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}

			modelByte, err := json.Marshal(&model)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			err = s.cache.Set([]byte("one_"+req.GetQuery()), modelByte, int(3*time.Hour))
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: &model}}, nil
		}
	}

	var model genres_pb.CategoryModel
	err = copier.Copy(&model, &cat, copier.WithPrimitiveToStringConverter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	modelByte, err := json.Marshal(&model)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.cache.Set([]byte("one_"+req.GetQuery()), modelByte, int(3*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: &model}}, nil
}
func (s *genreServer) GetTree(ctx context.Context, req *genres_pb.GetOneOfRequest) (*genres_pb.CategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := s.validator.StructValidation(req); err != nil {
		return nil, err
	}

	if category, err := s.cache.Get([]byte("tree_" + req.GetQuery())); len(category) > 0 && err == nil {
		var model genres_pb.CategoryModel
		err := json.Unmarshal(category, &model)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		return &genres_pb.CategoryResponse{Category: &model}, nil
	}

	cat, err := s.storage.FindCategoryTree(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	var model genres_pb.CategoryModel
	err = copier.Copy(&model, &cat, copier.WithPrimitiveToStringConverter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	modelByte, err := json.Marshal(&model)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.cache.Set([]byte("tree_"+req.GetQuery()), modelByte, int(3*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &genres_pb.CategoryResponse{Category: &model}, nil
}
