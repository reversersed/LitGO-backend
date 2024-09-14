package storage

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
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
	sync.RWMutex
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
		return nil, status.Error(codes.NotFound, "no authors found")
	}
	return books, nil
}
