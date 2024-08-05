package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/reversersed/go-grpc/tree/main/api_genre/internal/service/mocks"
	model "github.com/reversersed/go-grpc/tree/main/api_genre/internal/storage"
	genres_pb "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/proto/genres"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAll(t *testing.T) {
	models := []*model.Category{
		{
			Id:           primitive.NewObjectID(),
			Name:         "test",
			TranslitName: "test-213",
			Genres: []*model.Genre{
				{
					Id:           primitive.NewObjectID(),
					Name:         "genre",
					TranslitName: "test-genre",
					BookCount:    2,
				},
				{
					Id:           primitive.NewObjectID(),
					Name:         "genre2",
					TranslitName: "test-genre2",
					BookCount:    22,
				},
			},
		},
		{
			Id:           primitive.NewObjectID(),
			Name:         "test2",
			TranslitName: "test2-213",
			Genres: []*model.Genre{
				{
					Id:           primitive.NewObjectID(),
					Name:         "genre4",
					TranslitName: "test-genre4",
					BookCount:    24,
				},
				{
					Id:           primitive.NewObjectID(),
					Name:         "genre24",
					TranslitName: "test-genre24",
					BookCount:    24,
				},
			},
		},
	}
	response := []*genres_pb.CategoryModel{
		{
			Id:           models[0].Id.Hex(),
			Name:         "test",
			TranslitName: "test-213",
			Genres: []*genres_pb.GenreModel{
				{
					Id:           models[0].Genres[0].Id.Hex(),
					Name:         "genre",
					TranslitName: "test-genre",
					BookCount:    2,
				},
				{
					Id:           models[0].Genres[1].Id.Hex(),
					Name:         "genre2",
					TranslitName: "test-genre2",
					BookCount:    22,
				},
			},
		},
		{
			Id:           models[1].Id.Hex(),
			Name:         "test2",
			TranslitName: "test2-213",
			Genres: []*genres_pb.GenreModel{
				{
					Id:           models[1].Genres[0].Id.Hex(),
					Name:         "genre4",
					TranslitName: "test-genre4",
					BookCount:    24,
				},
				{
					Id:           models[1].Genres[1].Id.Hex(),
					Name:         "genre24",
					TranslitName: "test-genre24",
					BookCount:    24,
				},
			},
		},
	}
	table := []struct {
		Name             string
		MockBehaviour    func(*mocks.Mockcache, *mocks.Mocklogger, *mocks.Mockstorage, *mocks.Mockvalidator)
		ExceptedError    string
		ExceptedResponse *genres_pb.GetAllResponse
	}{
		{
			Name: "successful response",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m1.EXPECT().Get([]byte("all_categories")).Return([]byte{}, nil)
				m3.EXPECT().GetAll(gomock.Any()).Return(models, nil)
				m1.EXPECT().Set([]byte("all_categories"), gomock.Any(), gomock.Any())
			},
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetAllResponse{Categories: response},
		},
		{
			Name: "successful from cache",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				bytes, _ := json.Marshal(models)
				m1.EXPECT().Get([]byte("all_categories")).Return(bytes, nil)
			},
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetAllResponse{Categories: response},
		},
		{
			Name: "storage error",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m1.EXPECT().Get([]byte("all_categories")).Return([]byte{}, nil)
				m3.EXPECT().GetAll(gomock.Any()).Return(nil, status.Error(codes.Internal, "no database connection"))
			},
			ExceptedError:    status.Error(codes.Internal, "no database connection").Error(),
			ExceptedResponse: nil,
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m_cache := mocks.NewMockcache(ctrl)
			m_logger := mocks.NewMocklogger(ctrl)
			m_storage := mocks.NewMockstorage(ctrl)
			m_validator := mocks.NewMockvalidator(ctrl)

			server := NewServer(m_logger, m_cache, m_storage, m_validator)
			v.MockBehaviour(m_cache, m_logger, m_storage, m_validator)

			response, err := server.GetAll(context.Background(), &genres_pb.Empty{})
			if len(v.ExceptedError) > 0 {
				assert.EqualError(t, err, v.ExceptedError)
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, v.ExceptedResponse, response)
		})
	}
}
