package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	mock_authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors/mocks"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	mock_genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres/mocks"
	mock_service "github.com/reversersed/go-grpc/tree/main/api_book/internal/service/mocks"
	model "github.com/reversersed/go-grpc/tree/main/api_book/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetBookSuggestion(t *testing.T) {
	category := &genres_pb.CategoryModel{
		Id:   primitive.NewObjectID().Hex(),
		Name: "category",
		Genres: []*genres_pb.GenreModel{
			{
				Id:   primitive.NewObjectID().Hex(),
				Name: "genre",
			},
		},
	}
	author := &authors_pb.AuthorModel{
		Id:   primitive.NewObjectID().Hex(),
		Name: "Author",
	}
	authorId, _ := primitive.ObjectIDFromHex(author.GetId())
	genreId, _ := primitive.ObjectIDFromHex(category.GetGenres()[0].GetId())
	books := []*model.Book{
		{
			Id:           primitive.NewObjectID(),
			Name:         "name1",
			TranslitName: "isname",
			Genre:        genreId,
			Authors:      []primitive.ObjectID{authorId},
		},
		{
			Id:           primitive.NewObjectID(),
			Name:         "name2",
			TranslitName: "isname?",
			Genre:        genreId,
			Authors:      []primitive.ObjectID{authorId},
		},
		{
			Id:           primitive.NewObjectID(),
			Name:         "named author",
			TranslitName: "translit-name-21421",
			Genre:        genreId,
			Authors:      []primitive.ObjectID{authorId},
		},
		{
			Id:           primitive.NewObjectID(),
			Name:         "book without genre and author",
			TranslitName: "book-24142",
		},
	}
	bookModel := []*books_pb.BookModel{
		{
			Id:           books[0].Id.Hex(),
			Name:         "name1",
			Translitname: "isname",
			Category:     &books_pb.CategoryModel{Name: category.GetName(), Id: category.GetId()},
			Genre:        &books_pb.GenreModel{Name: category.GetGenres()[0].GetName(), Id: category.GetGenres()[0].GetId()},
			Authors:      []*books_pb.AuthorModel{{Name: author.Name, Id: author.Id}},
		},
		{
			Id:           books[1].Id.Hex(),
			Name:         "name2",
			Translitname: "isname?",
			Category:     &books_pb.CategoryModel{Name: category.GetName(), Id: category.GetId()},
			Genre:        &books_pb.GenreModel{Name: category.GetGenres()[0].GetName(), Id: category.GetGenres()[0].GetId()},
			Authors:      []*books_pb.AuthorModel{{Name: author.Name, Id: author.Id}},
		},
		{
			Id:           books[2].Id.Hex(),
			Name:         "named author",
			Translitname: "translit-name-21421",
			Category:     &books_pb.CategoryModel{Name: category.GetName(), Id: category.GetId()},
			Genre:        &books_pb.GenreModel{Name: category.GetGenres()[0].GetName(), Id: category.GetGenres()[0].GetId()},
			Authors:      []*books_pb.AuthorModel{{Name: author.Name, Id: author.Id}},
		},
		{
			Id:           books[3].Id.Hex(),
			Name:         "book without genre and author",
			Translitname: "book-24142",
		},
	}
	table := []struct {
		Name             string
		Request          *books_pb.GetSuggestionRequest
		ExceptedError    string
		ExceptedResponse *books_pb.GetBooksResponse
		MockBehaviour    func(*mock_service.Mockcache, *mock_authors_pb.MockAuthorClient, *mock_genres_pb.MockGenreClient, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
	}{
		{
			Name:          "validation error",
			Request:       &books_pb.GetSuggestionRequest{},
			ExceptedError: "rpc error: code = InvalidArgument desc = wrong arguments number",
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "wrong arguments number"))
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
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
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().GetSuggestions(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", int64(5)).Return(books, nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mgc.EXPECT().GetTree(gomock.Any(), gomock.Any()).Return(&genres_pb.CategoryResponse{Category: category}, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{author}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.GetBooksResponse{Books: bookModel},
		},
		{
			Name:    "successful with category from cache",
			Request: &books_pb.GetSuggestionRequest{Query: "Проверка правильности разбиения", Limit: 5},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().GetSuggestions(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", int64(5)).Return(books, nil)
				json, _ := json.Marshal(&genres_pb.CategoryResponse{Category: category})
				m1.EXPECT().Get(gomock.Any()).Return(json, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{author}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.GetBooksResponse{Books: bookModel},
		},
		{
			Name:    "storage error",
			Request: &books_pb.GetSuggestionRequest{Query: "Проверка правильности разбиения", Limit: 5},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
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
			rabbit := mock_service.NewMockrabbitservice(ctrl)
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

			service := NewServer(logger, cache, storage, validator, genreService, authorService, rabbit)
			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, authorService, genreService, logger, storage, validator)
			}

			response, err := service.GetBookSuggestions(context.Background(), v.Request)
			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			if v.ExceptedResponse != nil && assert.NotNil(t, response) {
				assert.Equal(t, v.ExceptedResponse, response)
			} else if v.ExceptedResponse == nil {
				assert.Nil(t, response)
			}
		})
	}
}
func TestCreateBook(t *testing.T) {
	category := &genres_pb.CategoryModel{
		Id:   primitive.NewObjectID().Hex(),
		Name: "category",
		Genres: []*genres_pb.GenreModel{
			{
				Id:   primitive.NewObjectID().Hex(),
				Name: "genre",
			},
		},
	}
	authors := []*authors_pb.AuthorModel{{
		Id:   primitive.NewObjectID().Hex(),
		Name: "Author",
	}, {
		Id:   primitive.NewObjectID().Hex(),
		Name: "Author2",
	}}
	book := &books_pb.BookModel{
		Id:          primitive.NilObjectID.Hex(),
		Name:        "book name",
		Description: "book Description",
		Picture:     "picture.jpg",
		Filepath:    "book.epub",
		Category:    &books_pb.CategoryModel{Id: category.GetId(), Name: "category"},
		Genre:       &books_pb.GenreModel{Id: category.GetGenres()[0].GetId(), Name: "genre"},
		Authors:     []*books_pb.AuthorModel{{Id: authors[0].GetId(), Name: "Author"}, {Id: authors[1].GetId(), Name: "Author2"}},
	}
	authorId, _ := primitive.ObjectIDFromHex(authors[0].GetId())
	authorId2, _ := primitive.ObjectIDFromHex(authors[1].GetId())
	genreId, _ := primitive.ObjectIDFromHex(category.GetGenres()[0].GetId())
	model := &model.Book{
		Id:          primitive.NilObjectID,
		Name:        "book name",
		Description: "book Description",
		Picture:     "picture.jpg",
		Filepath:    "book.epub",
		Genre:       genreId,
		Authors:     []primitive.ObjectID{authorId, authorId2},
	}
	table := []struct {
		Name             string
		Request          *books_pb.CreateBookRequest
		ExceptedError    string
		ExceptedResponse *books_pb.CreateBookResponse
		MockBehaviour    func(*mock_service.Mockcache, *mock_authors_pb.MockAuthorClient, *mock_genres_pb.MockGenreClient, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
	}{
		{
			Name: "successful",
			Request: &books_pb.CreateBookRequest{
				Name:        "book name",
				Description: "book Description",
				Picture:     "picture.jpg",
				Filepath:    "book.epub",
				Genre:       category.GetGenres()[0].GetId(),
				Authors:     []string{authors[0].GetId(), authors[1].GetId()},
			},
			ExceptedResponse: &books_pb.CreateBookResponse{Book: book},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				m3.EXPECT().CreateBook(gomock.Any(), model).Return(model, nil)
				mgc.EXPECT().GetTree(gomock.Any(), gomock.Any()).Return(&genres_pb.CategoryResponse{Category: category}, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{authors[0], authors[1]}}, nil).AnyTimes()
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
			genreService := mock_genres_pb.NewMockGenreClient(ctrl)
			authorService := mock_authors_pb.NewMockAuthorClient(ctrl)
			rabbit := mock_service.NewMockrabbitservice(ctrl)
			rabbit.EXPECT().SendBookCreatedMessage(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Info(gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

			service := NewServer(logger, cache, storage, validator, genreService, authorService, rabbit)
			if v.MockBehaviour != nil {
				v.MockBehaviour(cache, authorService, genreService, logger, storage, validator)
			}

			response, err := service.CreateBook(context.Background(), v.Request)
			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			if v.ExceptedResponse != nil && assert.NotNil(t, response) {
				assert.Equal(t, v.ExceptedResponse, response)
			} else if v.ExceptedResponse == nil {
				assert.Nil(t, response)
			}
		})
	}
}
