package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors/mock"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	mock_genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres/mock"
	mock_service "github.com/reversersed/go-grpc/tree/main/api_book/internal/service/mocks"
	model "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetBookSuggestion(t *testing.T) {
	books := []*model.Book{
		{
			Id:           primitive.NewObjectID(),
			Name:         "name1",
			TranslitName: "isname",
		},
		{
			Id:           primitive.NewObjectID(),
			Name:         "name2",
			TranslitName: "isname?",
		},
		{
			Id:           primitive.NewObjectID(),
			Name:         "named author",
			TranslitName: "translit-name-21421",
		},
	}
	bookModel := []*books_pb.BookModel{
		{
			Id:           books[0].Id.Hex(),
			Name:         "name1",
			Translitname: "isname",
		},
		{
			Id:           books[1].Id.Hex(),
			Name:         "name2",
			Translitname: "isname?",
		},
		{
			Id:           books[2].Id.Hex(),
			Name:         "named author",
			Translitname: "translit-name-21421",
		},
	}
	table := []struct {
		Name             string
		Request          *books_pb.GetSuggestionRequest
		ExceptedError    string
		ExceptedResponse *books_pb.GetBooksResponse
		MockBehaviour    func(*mock_service.Mockcache, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
	}{
		{
			Name:          "validation error",
			Request:       &books_pb.GetSuggestionRequest{},
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
			Request: &books_pb.GetSuggestionRequest{Query: "Проверка правильности разбиения", Limit: 5},
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().GetSuggestions(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", int64(5)).Return(books, nil)
			},
			ExceptedResponse: &books_pb.GetBooksResponse{Books: bookModel},
		},
		{
			Name:    "storage error",
			Request: &books_pb.GetSuggestionRequest{Query: "Проверка правильности разбиения", Limit: 5},
			MockBehaviour: func(m1 *mock_service.Mockcache, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().GetSuggestions(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", int64(5)).Return(nil, status.Error(codes.NotFound, "authors not found"))
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
			genreService := mock_genres_pb.NewMockGenreClient(ctrl)
			authorService := mock_authors_pb.NewMockAuthorClient(ctrl)

			service := NewServer(logger, cache, storage, validator, genreService, authorService)
			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, logger, storage, validator)
			}

			response, err := service.GetBookSuggestions(context.Background(), v.Request)
			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.ExceptedResponse, response)
		})
	}
}
