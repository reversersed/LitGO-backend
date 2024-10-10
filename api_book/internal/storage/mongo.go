package storage

import (
	"context"

	"github.com/reversersed/go-grpc/tree/main/api_book/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source=mongo.go -destination=mocks/mongo.go

type logger interface {
	Infof(string, ...any)
	Info(...any)
	Warnf(string, ...any)
	Warn(...any)
	Fatalf(string, ...any)
	Fatal(...any)
}
type db struct {
	logger     logger
	collection *mongodb.Collection
}

func NewStorage(storage *mongodb.Database, collection string, logger logger) *db {
	db := &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
	return db
}

// TODDO write test for create book
func (d *db) CreateBook(ctx context.Context, book *Book) (*Book, error) {
	book.Id = primitive.NewObjectID()
	book.TranslitName = mongo.GenerateTranslitName(book.Name, book.Id)

	result, err := d.collection.InsertOne(ctx, book)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok || id != book.Id {
		return nil, status.Error(codes.Internal, "error retrieving book id")
	}
	return book, nil
}
func (d *db) GetSuggestions(ctx context.Context, regex string, limit int64) ([]*Book, error) {
	response, err := d.collection.Find(ctx, bson.M{"name": bson.M{"$regex": regex, "$options": "i"}}, &options.FindOptions{Limit: &limit})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	books := make([]*Book, 0)
	err = response.All(ctx, &books)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(books) == 0 {
		return nil, status.Error(codes.NotFound, "no books found")
	}
	return books, nil
}
