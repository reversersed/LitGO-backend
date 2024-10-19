package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	mocks "github.com/reversersed/LitGO-backend/tree/main/api_genre/internal/service/mocks"
	model "github.com/reversersed/LitGO-backend/tree/main/api_genre/internal/storage"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
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
			m_rabbit := mocks.NewMockrabbitservice(ctrl)

			server := NewServer(m_logger, m_cache, m_storage, m_validator, m_rabbit)
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
func TestGetCategoryTree(t *testing.T) {
	model := &model.Category{

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
	}
	response := &genres_pb.CategoryModel{
		Id:           model.Id.Hex(),
		Name:         "test",
		TranslitName: "test-213",
		Genres: []*genres_pb.GenreModel{
			{
				Id:           model.Genres[0].Id.Hex(),
				Name:         "genre",
				TranslitName: "test-genre",
				BookCount:    2,
			},
			{
				Id:           model.Genres[1].Id.Hex(),
				Name:         "genre2",
				TranslitName: "test-genre2",
				BookCount:    22,
			},
		},
	}

	table := []struct {
		Name             string
		MockBehaviour    func(*mocks.Mockcache, *mocks.Mocklogger, *mocks.Mockstorage, *mocks.Mockvalidator)
		Query            string
		ExceptedError    string
		ExceptedResponse *genres_pb.CategoryResponse
	}{
		{
			Name: "get by genre translit name",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("tree_"+model.Genres[0].TranslitName)).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("tree_"+model.Genres[0].TranslitName), gomock.Any(), gomock.Any())
			},
			Query:            model.Genres[0].TranslitName,
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.CategoryResponse{Category: response},
		},
		{
			Name: "get by genre id",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("tree_"+model.Genres[0].Id.Hex())).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("tree_"+model.Genres[0].Id.Hex()), gomock.Any(), gomock.Any())
			},
			Query:            model.Genres[0].Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.CategoryResponse{Category: response},
		},
		{
			Name: "get by category translit name",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("tree_"+model.TranslitName)).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("tree_"+model.TranslitName), gomock.Any(), gomock.Any())
			},
			Query:            model.TranslitName,
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.CategoryResponse{Category: response},
		},
		{
			Name: "get by category id",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("tree_"+model.Id.Hex())).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("tree_"+model.Id.Hex()), gomock.Any(), gomock.Any())
			},
			Query:            model.Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.CategoryResponse{Category: response},
		},
		{
			Name: "get from cache",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				bytes, _ := json.Marshal(&response)
				m1.EXPECT().Get([]byte("tree_"+model.Id.Hex())).Return(bytes, nil)
			},
			Query:            model.Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.CategoryResponse{Category: response},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m_cache := mocks.NewMockcache(ctrl)
			m_logger := mocks.NewMocklogger(ctrl)
			m_storage := mocks.NewMockstorage(ctrl)
			m_validator := mocks.NewMockvalidator(ctrl)
			m_rabbit := mocks.NewMockrabbitservice(ctrl)

			server := NewServer(m_logger, m_cache, m_storage, m_validator, m_rabbit)
			v.MockBehaviour(m_cache, m_logger, m_storage, m_validator)

			response, err := server.GetTree(context.Background(), &genres_pb.GetOneOfRequest{Query: v.Query})
			if len(v.ExceptedError) > 0 {
				assert.EqualError(t, err, v.ExceptedError)
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, v.ExceptedResponse, response)
		})
	}
}
func TestGetOneOf(t *testing.T) {
	model := &model.Category{

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
	}
	response := &genres_pb.CategoryModel{
		Id:           model.Id.Hex(),
		Name:         "test",
		TranslitName: "test-213",
		Genres: []*genres_pb.GenreModel{
			{
				Id:           model.Genres[0].Id.Hex(),
				Name:         "genre",
				TranslitName: "test-genre",
				BookCount:    2,
			},
			{
				Id:           model.Genres[1].Id.Hex(),
				Name:         "genre2",
				TranslitName: "test-genre2",
				BookCount:    22,
			},
		},
	}

	table := []struct {
		Name             string
		MockBehaviour    func(*mocks.Mockcache, *mocks.Mocklogger, *mocks.Mockstorage, *mocks.Mockvalidator)
		Query            string
		ExceptedError    string
		ExceptedResponse *genres_pb.GetCategoryOrGenreResponse
	}{
		{
			Name: "get by genre translit name",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("one_"+model.Genres[0].TranslitName)).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("one_"+model.Genres[0].TranslitName), gomock.Any(), gomock.Any())
			},
			Query:            model.Genres[0].TranslitName,
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: response.Genres[0]}},
		},
		{
			Name: "get by genre id",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("one_"+model.Genres[0].Id.Hex())).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("one_"+model.Genres[0].Id.Hex()), gomock.Any(), gomock.Any())
			},
			Query:            model.Genres[0].Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: response.Genres[0]}},
		},
		{
			Name: "get category by id",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("one_"+model.Id.Hex())).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("one_"+model.Id.Hex()), gomock.Any(), gomock.Any())
			},
			Query:            model.Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: response}},
		},
		{
			Name: "get category by translit",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("one_"+model.TranslitName)).Return([]byte{}, nil)
				m3.EXPECT().FindCategoryTree(gomock.Any(), gomock.Any()).Return(model, nil)
				m1.EXPECT().Set([]byte("one_"+model.TranslitName), gomock.Any(), gomock.Any())
			},
			Query:            model.TranslitName,
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: response}},
		},
		{
			Name: "get genre from cache by translit",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				bytes, _ := json.Marshal(&response.GetGenres()[0])
				m1.EXPECT().Get([]byte("one_"+model.Genres[0].TranslitName)).Return(bytes, nil)
			},
			Query:            model.Genres[0].TranslitName,
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: response.Genres[0]}},
		},
		{
			Name: "get genre from cache by id",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				bytes, _ := json.Marshal(&response.GetGenres()[0])
				m1.EXPECT().Get([]byte("one_"+model.Genres[0].Id.Hex())).Return(bytes, nil)
			},
			Query:            model.Genres[0].Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Genre{Genre: response.Genres[0]}},
		},
		{
			Name: "get category from cache by translit",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				bytes, _ := json.Marshal(&response)
				m1.EXPECT().Get([]byte("one_"+model.TranslitName)).Return(bytes, nil)
			},
			Query:            model.TranslitName,
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: response}},
		},
		{
			Name: "get genre from cache by id",
			MockBehaviour: func(m1 *mocks.Mockcache, m2 *mocks.Mocklogger, m3 *mocks.Mockstorage, m4 *mocks.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				bytes, _ := json.Marshal(&response)
				m1.EXPECT().Get([]byte("one_"+model.Id.Hex())).Return(bytes, nil)
			},
			Query:            model.Id.Hex(),
			ExceptedError:    "",
			ExceptedResponse: &genres_pb.GetCategoryOrGenreResponse{Model: &genres_pb.GetCategoryOrGenreResponse_Category{Category: response}},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m_cache := mocks.NewMockcache(ctrl)
			m_logger := mocks.NewMocklogger(ctrl)
			m_storage := mocks.NewMockstorage(ctrl)
			m_validator := mocks.NewMockvalidator(ctrl)
			m_rabbit := mocks.NewMockrabbitservice(ctrl)

			server := NewServer(m_logger, m_cache, m_storage, m_validator, m_rabbit)
			v.MockBehaviour(m_cache, m_logger, m_storage, m_validator)

			response, err := server.GetOneOf(context.Background(), &genres_pb.GetOneOfRequest{Query: v.Query})
			if len(v.ExceptedError) > 0 {
				assert.EqualError(t, err, v.ExceptedError)
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, v.ExceptedResponse, response)
		})
	}
}
