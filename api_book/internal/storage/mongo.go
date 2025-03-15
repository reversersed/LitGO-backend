package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/reversersed/LitGO-backend-pkg/mongo"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
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
func (d *db) Find(ctx context.Context, regex string, limit, page int, rating float32, sort SortType) ([]*Book, error) {
	options := &options.FindOptions{}
	options.SetLimit(int64(limit))
	options.SetSkip(int64(page * limit))

	switch sort {
	case Popular:
		options.SetSort(bson.M{"monthpurchases": -1})
	case Newest:
		options.SetSort(bson.M{"published": -1})
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("not known sort type: %v", sort))
	}

	response, err := d.collection.Find(ctx, bson.M{"$and": []bson.M{{"name": bson.M{"$regex": regex, "$options": "i"}}, {"rating": bson.M{"$gte": rating}}}}, options)
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
func (d *db) GetBookByGenre(ctx context.Context, genreIds []primitive.ObjectID, sortType SortType, onlyHighRating bool, limit int, page int) ([]*Book, error) {
	options := options.Find()

	if len(genreIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "received empty array")
	}

	filter := bson.M{"genre": bson.M{"$in": genreIds}}

	if onlyHighRating {
		filter["rating"] = bson.M{"$gte": 4.0}
	}

	switch sortType {
	case Popular:
		options.SetSort(bson.M{"monthpurchases": -1})
	case Newest:
		options.SetSort(bson.M{"published": -1})
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("not known sort type: %v", sortType))
	}
	options.SetSkip(int64(page * limit))
	options.SetLimit(int64(limit))

	response, err := d.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer response.Close(ctx)

	var book []*Book = make([]*Book, 0)
	if err := response.All(ctx, &book); err != nil {
		return nil, status.Error(codes.Internal, "error decoding response: "+err.Error())
	}
	if len(book) == 0 {
		return nil, status.Error(codes.NotFound, "books not found")
	}
	return book, nil
}
func (d *db) GetBookList(ctx context.Context, id []primitive.ObjectID, translit []string) ([]*Book, error) {
	if len(id) == 0 && len(translit) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no id or translit name argument presented")
	}

	books := make([]*Book, 0)
	result, err := d.collection.Find(ctx, bson.M{"$or": bson.A{bson.M{"_id": bson.D{{Key: "$in", Value: id}}}, bson.M{"translit": bson.D{{Key: "$in", Value: translit}}}}})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = result.All(ctx, &books)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(books) == 0 {
		var str string
		for _, i := range id {
			str += i.Hex() + ","
		}
		status, _ := status.New(codes.NotFound, "no books found").WithDetails(&shared_pb.ErrorDetail{
			Field:       "id",
			Struct:      "books_pb.GetBookListRequest",
			Actualvalue: strings.Trim(str, ","),
		}, &shared_pb.ErrorDetail{
			Field:       "translit",
			Struct:      "books_pb.GetBookListRequest",
			Actualvalue: strings.Join(translit, ","),
		})
		return nil, status.Err()
	}

	return books, nil
}
