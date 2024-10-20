package storage

import (
	"context"
	"time"

	"github.com/reversersed/LitGO-backend/tree/main/api_book/pkg/mongo"
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

func (d *db) CreateBook(ctx context.Context, book *Book) (*Book, error) {
	book.Id = primitive.NewObjectID()
	book.TranslitName = mongo.GenerateTranslitName(book.Name, book.Id)
	book.Published = time.Now().UTC().Unix()

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
func (d *db) Find(ctx context.Context, regex string, limit, page int, rating float32) ([]*Book, error) {
	lim := int64(limit)
	skip := int64(page * limit)
	response, err := d.collection.Find(ctx, bson.M{"$and": []bson.M{bson.M{"name": bson.M{"$regex": regex, "$options": "i"}}, bson.M{"rating": bson.M{"$gte": rating}}}}, &options.FindOptions{Limit: &lim, Skip: &skip})
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

// Query can be ID or translit name. Id has higher priority
func (d *db) GetBook(ctx context.Context, query string) (*Book, error) {
	id, err := primitive.ObjectIDFromHex(query)

	var response *mongodb.SingleResult

	if err == nil {
		response = d.collection.FindOne(ctx, bson.M{"_id": id})
	} else {
		response = d.collection.FindOne(ctx, bson.M{"translit": query})
	}
	if response == nil {
		return nil, status.Error(codes.NotFound, "response was nil")
	}
	if response.Err() != nil {
		return nil, status.Error(codes.NotFound, response.Err().Error())
	}

	var book Book
	if err = response.Decode(&book); err != nil {
		return nil, status.Error(codes.Internal, "error decoding response: "+err.Error())
	}
	return &book, nil
}
