package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_service "github.com/reversersed/LitGO-backend/tree/main/api_author/internal/service/mocks"
	model "github.com/reversersed/LitGO-backend/tree/main/api_author/internal/storage"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAuthors(t *testing.T) {
	authors := []*model.Author{
		{
			Id:             primitive.NewObjectID(),
			Name:           "name1",
			TranslitName:   "isname",
			About:          "about this author",
			ProfilePicture: "url",
			Rating:         5.0,
		},
		{
			Id:             primitive.NewObjectID(),
			Name:           "name2",
			TranslitName:   "isname?",
			About:          "about this another author",
			ProfilePicture: "url to pic",
			Rating:         2.3,
		},
		{
			Id:             primitive.NewObjectID(),
			Name:           "named author",
			TranslitName:   "translit-name-21421",
			About:          "about this author again",
			ProfilePicture: "urls...",
			Rating:         0,
		},
	}
	authorModel := []*authors_pb.AuthorModel{
		{
			Id:             authors[0].Id.Hex(),
			Name:           "name1",
			Translitname:   "isname",
			About:          "about this author",
			Profilepicture: "url",
			Rating:         5.0,
		},
		{
			Id:             authors[1].Id.Hex(),
			Name:           "name2",
			Translitname:   "isname?",
			About:          "about this another author",
			Profilepicture: "url to pic",
			Rating:         2.3,
		},
		{
			Id:             authors[2].Id.Hex(),
			Name:           "named author",
			Translitname:   "translit-name-21421",
			About:          "about this author again",
			Profilepicture: "urls...",
			Rating:         0,
		},
	}
	table := []struct {
		Name             string
		Request          *authors_pb.GetAuthorsRequest
		ExceptedError    string
		ExceptedResponse *authors_pb.GetAuthorsResponse
		MockBehaviour    func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
	}{
		{
			Name:          "validation error",
			Request:       &authors_pb.GetAuthorsRequest{},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong arguments number",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong arguments number"))
			},
		},
		{
			Name:          "wrong id",
			Request:       &authors_pb.GetAuthorsRequest{Id: []string{primitive.NewObjectID().Hex(), "not id"}},
			ExceptedError: "rpc error: code = InvalidArgument desc = unable to convert value to type",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
			},
		},
		{
			Name:          "nil request",
			Request:       nil,
			ExceptedError: "rpc error: code = InvalidArgument desc = received nil request",
		},
		{
			Name: "successful response",
			Request: &authors_pb.GetAuthorsRequest{
				Id:       []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()},
				Translit: []string{"name1,name2,name3"},
			},
			ExceptedResponse: &authors_pb.GetAuthorsResponse{Authors: authorModel},
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().GetAuthors(gomock.Any(), gomock.Any(), []string{"name1,name2,name3"}).Return(authors, nil)
			},
		},
		{
			Name:          "storage error",
			Request:       &authors_pb.GetAuthorsRequest{},
			ExceptedError: "rpc error: code = NotFound desc = authors not found",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().GetAuthors(gomock.Any(), nil, nil).Return(nil, status.Error(codes.NotFound, "authors not found"))
			},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

			service := NewServer(logger, cache, storage, validator)
			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}

			response, err := service.GetAuthors(context.Background(), v.Request)
			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.ExceptedResponse, response)
		})
	}
}
func TestGetAuthorSuggestion(t *testing.T) {
	authors := []*model.Author{
		{
			Id:             primitive.NewObjectID(),
			Name:           "name1",
			TranslitName:   "isname",
			About:          "about this author",
			ProfilePicture: "url",
			Rating:         5.0,
		},
		{
			Id:             primitive.NewObjectID(),
			Name:           "name2",
			TranslitName:   "isname?",
			About:          "about this another author",
			ProfilePicture: "url to pic",
			Rating:         2.3,
		},
		{
			Id:             primitive.NewObjectID(),
			Name:           "named author",
			TranslitName:   "translit-name-21421",
			About:          "about this author again",
			ProfilePicture: "urls...",
			Rating:         0,
		},
	}
	authorModel := []*authors_pb.AuthorModel{
		{
			Id:             authors[0].Id.Hex(),
			Name:           "name1",
			Translitname:   "isname",
			About:          "about this author",
			Profilepicture: "url",
			Rating:         5.0,
		},
		{
			Id:             authors[1].Id.Hex(),
			Name:           "name2",
			Translitname:   "isname?",
			About:          "about this another author",
			Profilepicture: "url to pic",
			Rating:         2.3,
		},
		{
			Id:             authors[2].Id.Hex(),
			Name:           "named author",
			Translitname:   "translit-name-21421",
			About:          "about this author again",
			Profilepicture: "urls...",
			Rating:         0,
		},
	}
	table := []struct {
		Name             string
		Request          *authors_pb.FindAuthorsRequest
		ExceptedError    string
		ExceptedResponse *authors_pb.GetAuthorsResponse
		MockBehaviour    func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
	}{
		{
			Name:          "validation error",
			Request:       &authors_pb.FindAuthorsRequest{},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong arguments number",
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong arguments number"))
			},
		},
		{
			Name:          "nil request",
			Request:       nil,
			ExceptedError: "rpc error: code = InvalidArgument desc = received nil request",
		},
		{
			Name:    "successful",
			Request: &authors_pb.FindAuthorsRequest{Query: "Проверка правильности разбиения", Limit: 5},
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().Find(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", 5, 0, float32(0.0)).Return(authors, nil)
			},
			ExceptedResponse: &authors_pb.GetAuthorsResponse{Authors: authorModel},
		},
		{
			Name:    "storage error",
			Request: &authors_pb.FindAuthorsRequest{Query: "Проверка правильности разбиения", Limit: 5, Rating: 2.0},
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().Find(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", 5, 0, float32(2.0)).Return(nil, status.Error(codes.NotFound, "authors not found"))
			},
			ExceptedError: "rpc error: code = NotFound desc = authors not found",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			logger := mock_service.NewMocklogger(ctrl)
			cache := mock_service.NewMockcache(ctrl)
			storage := mock_service.NewMockstorage(ctrl)
			validator := mock_service.NewMockvalidator(ctrl)
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

			service := NewServer(logger, cache, storage, validator)
			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}

			response, err := service.FindAuthors(context.Background(), v.Request)
			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.ExceptedResponse, response)
		})
	}
}
