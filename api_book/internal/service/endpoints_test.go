package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mock_service "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/service/mocks"
	model "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage"
	authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors"
	mock_authors_pb "github.com/reversersed/LitGO-proto/gen/go/authors/mocks"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	mock_genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFindBook(t *testing.T) {
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
			Category:     category,
			Genre:        category.GetGenres()[0],
			Authors:      []*authors_pb.AuthorModel{author},
		},
		{
			Id:           books[1].Id.Hex(),
			Name:         "name2",
			Translitname: "isname?",
			Category:     category,
			Genre:        category.GetGenres()[0],
			Authors:      []*authors_pb.AuthorModel{author},
		},
		{
			Id:           books[2].Id.Hex(),
			Name:         "named author",
			Translitname: "translit-name-21421",
			Category:     category,
			Genre:        category.GetGenres()[0],
			Authors:      []*authors_pb.AuthorModel{author},
		},
		{
			Id:           books[3].Id.Hex(),
			Name:         "book without genre and author",
			Translitname: "book-24142",
		},
	}
	table := []struct {
		Name             string
		Request          *books_pb.FindBookRequest
		ExceptedError    string
		ExceptedResponse *books_pb.FindBookResponse
		MockBehaviour    func(*mock_service.Mockcache, *mock_authors_pb.MockAuthorClient, *mock_genres_pb.MockGenreClient, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
	}{
		{
			Name:          "validation error",
			Request:       &books_pb.FindBookRequest{},
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
			Request: &books_pb.FindBookRequest{Query: "Проверка правильности разбиения", Limit: 5, Page: 1, Sorttype: "Popular"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().Find(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", 5, 1, float32(0.0), model.Popular).Return(books, nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mgc.EXPECT().GetTree(gomock.Any(), gomock.Any()).Return(&genres_pb.CategoryResponse{Category: category}, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{author}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.FindBookResponse{Books: bookModel},
		},
		{
			Name:    "successful with category from cache",
			Request: &books_pb.FindBookRequest{Query: "Проверка правильности разбиения", Limit: 5, Sorttype: "Popular"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().Find(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", 5, 0, float32(0.0), model.Popular).Return(books, nil)
				json, _ := json.Marshal(&genres_pb.CategoryResponse{Category: category})
				m1.EXPECT().Get(gomock.Any()).Return(json, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{author}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.FindBookResponse{Books: bookModel},
		},
		{
			Name:    "successful with no query request",
			Request: &books_pb.FindBookRequest{Limit: 5, Sorttype: "Popular"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m3.EXPECT().Find(gomock.Any(), "(.*?)", 5, 0, float32(0.0), model.Popular).Return(books, nil)
				json, _ := json.Marshal(&genres_pb.CategoryResponse{Category: category})
				m1.EXPECT().Get(gomock.Any()).Return(json, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{author}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.FindBookResponse{Books: bookModel},
		},
		{
			Name:    "storage error",
			Request: &books_pb.FindBookRequest{Query: "Проверка правильности разбиения", Limit: 5, Rating: 2.0, Sorttype: "Popular"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				m3.EXPECT().Find(gomock.Any(), "(Проверка)|(правильности)|(разбиения)", 5, 0, float32(2.0), model.Popular).Return(nil, status.Error(codes.NotFound, "books not found"))
			},
			ExceptedError: "rpc error: code = NotFound desc = books not found",
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

			response, err := service.FindBook(context.Background(), v.Request)
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

/*
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
			Category:    &genres_pb.CategoryModel{Id: category.GetId(), Name: "category"},
			Genre:       &genres_pb.GenreModel{Id: category.GetGenres()[0].GetId(), Name: "genre"},
			Authors:     []*authors_pb.AuthorModel{{Id: authors[0].GetId(), Name: "Author"}, {Id: authors[1].GetId(), Name: "Author2"}},
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
*/
func TestGetBook(t *testing.T) {
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
		Category:    category,
		Genre:       category.GetGenres()[0],
		Authors:     authors,
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
		MockBehaviour    func(*mock_service.Mockcache, *mock_authors_pb.MockAuthorClient, *mock_genres_pb.MockGenreClient, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		ExceptedError    string
		ExceptedResponse *books_pb.GetBookResponse
		Request          *books_pb.GetBookRequest
	}{
		{
			Name:          "nil request",
			ExceptedError: "rpc error: code = InvalidArgument desc = received nil request",
		},
		{
			Name:          "empty request",
			ExceptedError: "rpc error: code = InvalidArgument desc = validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "validation error"))
			},
			Request: &books_pb.GetBookRequest{},
		},
		{
			Name:    "get from cache",
			Request: &books_pb.GetBookRequest{Query: "bookId"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				bytes, _ := json.Marshal(model)
				m1.EXPECT().Get([]byte("book_bookId")).Return(bytes, nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mgc.EXPECT().GetTree(gomock.Any(), gomock.Any()).Return(&genres_pb.CategoryResponse{Category: category}, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{authors[0], authors[1]}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.GetBookResponse{Book: book},
		},
		{
			Name:    "get from database",
			Request: &books_pb.GetBookRequest{Query: "bookId"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("book_bookId")).Return([]byte{}, errors.New(""))
				m3.EXPECT().GetBook(gomock.Any(), "bookId").Return(model, nil)
				m1.EXPECT().Set([]byte("book_bookId"), gomock.Any(), gomock.Any()).Return(nil)
				m1.EXPECT().Get(gomock.Any()).Return([]byte{}, errors.New("")).AnyTimes()
				m1.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
				mgc.EXPECT().GetTree(gomock.Any(), gomock.Any()).Return(&genres_pb.CategoryResponse{Category: category}, nil).AnyTimes()
				mac.EXPECT().GetAuthors(gomock.Any(), gomock.Any()).Return(&authors_pb.GetAuthorsResponse{Authors: []*authors_pb.AuthorModel{authors[0], authors[1]}}, nil).AnyTimes()
			},
			ExceptedResponse: &books_pb.GetBookResponse{Book: book},
		},
		{
			Name:    "error from storage",
			Request: &books_pb.GetBookRequest{Query: "bookId"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("book_bookId")).Return([]byte{}, errors.New(""))
				m3.EXPECT().GetBook(gomock.Any(), "bookId").Return(nil, status.Error(codes.NotFound, "not found book"))
			},
			ExceptedError: "rpc error: code = NotFound desc = not found book",
		},
		{
			Name:    "error from cache",
			Request: &books_pb.GetBookRequest{Query: "bookId"},
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				m1.EXPECT().Get([]byte("book_bookId")).Return([]byte{}, errors.New(""))
				m3.EXPECT().GetBook(gomock.Any(), "bookId").Return(model, nil)
				m1.EXPECT().Set([]byte("book_bookId"), gomock.Any(), gomock.Any()).Return(errors.New("not enough space"))
			},
			ExceptedError: "rpc error: code = Internal desc = error saving book to cache: not enough space",
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

			response, err := service.GetBook(context.Background(), v.Request)
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
func TestGetBookByGenre(t *testing.T) {
	category := &genres_pb.CategoryModel{Id: primitive.NewObjectID().Hex(), Translitname: "category-1", Genres: []*genres_pb.GenreModel{
		{
			Id:           primitive.NewObjectID().Hex(),
			Translitname: "genre-1",
		},
		{
			Id:           primitive.NewObjectID().Hex(),
			Translitname: "genre-2",
		}, {
			Id:           primitive.NewObjectID().Hex(),
			Translitname: "genre-3",
		},
	}}
	book := &model.Book{Name: "book"}
	table := []struct {
		Name             string
		MockBehaviour    func(*mock_service.Mockcache, *mock_authors_pb.MockAuthorClient, *mock_genres_pb.MockGenreClient, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		ExceptedError    string
		ExceptedResponse *books_pb.GetBookByGenreResponse
		Request          *books_pb.GetBookByGenreRequest
	}{
		{
			Name:          "nil request",
			ExceptedError: "rpc error: code = InvalidArgument desc = received nil request",
			Request:       nil,
		},
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "validation error"))
			},
			ExceptedError: "rpc error: code = InvalidArgument desc = validation error",
			Request:       &books_pb.GetBookByGenreRequest{},
		},
		{
			Name: "success by single genre",
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				id, _ := primitive.ObjectIDFromHex(category.GetGenres()[0].GetId())

				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				mgc.EXPECT().GetTree(gomock.Any(), &genres_pb.GetOneOfRequest{Query: category.GetGenres()[0].GetId()}).Return(&genres_pb.CategoryResponse{Category: category}, nil)
				m3.EXPECT().GetBookByGenre(gomock.Any(), []primitive.ObjectID{id}, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.Book{book}, nil)
			},
			Request:          &books_pb.GetBookByGenreRequest{Query: category.GetGenres()[0].GetId(), Limit: 1, Page: 0, Sorttype: "Newest"},
			ExceptedResponse: &books_pb.GetBookByGenreResponse{Books: []*books_pb.BookModel{{Id: string(primitive.NilObjectID.Hex()), Name: "book"}}},
		},
		{
			Name: "success by whole category",
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				genres := make([]primitive.ObjectID, 0)
				for _, v := range category.GetGenres() {
					id, _ := primitive.ObjectIDFromHex(v.GetId())
					genres = append(genres, id)
				}

				m4.EXPECT().StructValidation(gomock.Any()).Return(nil)
				mgc.EXPECT().GetTree(gomock.Any(), &genres_pb.GetOneOfRequest{Query: category.Id}).Return(&genres_pb.CategoryResponse{Category: category}, nil)
				m3.EXPECT().GetBookByGenre(gomock.Any(), genres, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.Book{book}, nil)
			},
			Request:          &books_pb.GetBookByGenreRequest{Query: category.GetId(), Limit: 1, Page: 0, Sorttype: "Newest"},
			ExceptedResponse: &books_pb.GetBookByGenreResponse{Books: []*books_pb.BookModel{{Id: string(primitive.NilObjectID.Hex()), Name: "book"}}},
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

			response, err := service.GetBookByGenre(context.Background(), v.Request)
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

// TODO write tests
func TestGetBookList(t *testing.T) {
	/*category := &genres_pb.CategoryModel{Id: primitive.NewObjectID().Hex(), Translitname: "category-1", Genres: []*genres_pb.GenreModel{
		{
			Id:           primitive.NewObjectID().Hex(),
			Translitname: "genre-1",
		},
		{
			Id:           primitive.NewObjectID().Hex(),
			Translitname: "genre-2",
		}, {
			Id:           primitive.NewObjectID().Hex(),
			Translitname: "genre-3",
		},
	}}
	/book := &model.Book{Name: "book"}*/
	table := []struct {
		Name             string
		MockBehaviour    func(*mock_service.Mockcache, *mock_authors_pb.MockAuthorClient, *mock_genres_pb.MockGenreClient, *mock_service.Mocklogger, *mock_service.Mockstorage, *mock_service.Mockvalidator)
		ExceptedError    string
		ExceptedResponse *books_pb.GetBookListResponse
		Request          *books_pb.GetBookListRequest
	}{
		{
			Name:          "nil request",
			ExceptedError: "rpc error: code = InvalidArgument desc = received nil request",
			Request:       nil,
		},
		{
			Name: "validation error",
			MockBehaviour: func(m1 *mock_service.Mockcache, mac *mock_authors_pb.MockAuthorClient, mgc *mock_genres_pb.MockGenreClient, m2 *mock_service.Mocklogger, m3 *mock_service.Mockstorage, m4 *mock_service.Mockvalidator) {
				m4.EXPECT().StructValidation(gomock.Any()).Return(status.Error(codes.InvalidArgument, "validation error"))
			},
			ExceptedError: "rpc error: code = InvalidArgument desc = validation error",
			Request:       &books_pb.GetBookListRequest{},
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

			response, err := service.GetBookList(context.Background(), v.Request)
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
