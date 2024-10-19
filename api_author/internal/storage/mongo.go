package storage

import (
	"context"
	"strings"
	"time"

	"github.com/reversersed/LitGO-backend/tree/main/api_author/pkg/mongo"
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
	defer db.seedAuthors()
	return db
}
func (d *db) seedAuthors() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if count, _ := d.collection.CountDocuments(ctx, bson.D{}); count > 0 {
		d.logger.Infof("there are %d authors in base, seeding canceled", count)
		return
	}

	for _, a := range mocked_authors {
		if _, err := d.CreateAuthor(ctx, a); err != nil {
			d.logger.Fatalf("error seeding author %v: %v", a, err)
		}
	}
	d.logger.Infof("seeded %d authors", len(mocked_authors))
}
func (d *db) CreateAuthor(ctx context.Context, author *Author) (*Author, error) {
	author.Id = primitive.NewObjectID()
	author.TranslitName = mongo.GenerateTranslitName(author.Name, author.Id)

	result, err := d.collection.InsertOne(ctx, author)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok || id != author.Id {
		return nil, status.Error(codes.Internal, "error retrieving author id")
	}

	return author, nil
}
func (d *db) GetAuthors(ctx context.Context, id []primitive.ObjectID, translit []string) ([]*Author, error) {
	authors := make([]*Author, 0)

	if len(id) > 0 {
		result, err := d.collection.Find(ctx, bson.M{"_id": bson.D{{Key: "$in", Value: id}}})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		var temp []*Author
		err = result.All(ctx, &temp)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		authors = append(authors, temp...)
	}
	if len(translit) > 0 {
		result, err := d.collection.Find(ctx, bson.M{"translit": bson.D{{Key: "$in", Value: translit}}})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		var temp []*Author
		err = result.All(ctx, &temp)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		authors = append(authors, temp...)
	}
	if len(id) == 0 && len(translit) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no id or translit name argument presented")
	}
	if len(authors) == 0 {
		var str string
		for _, i := range id {
			str += i.Hex() + ","
		}
		status, _ := status.New(codes.NotFound, "no authors found").WithDetails(&shared_pb.ErrorDetail{
			Field:       "id",
			Struct:      "authors_pb.GetAuthorsRequest",
			Actualvalue: strings.Trim(str, ","),
		}, &shared_pb.ErrorDetail{
			Field:       "translit",
			Struct:      "authors_pb.GetAuthorsRequest",
			Actualvalue: strings.Join(translit, ","),
		})
		return nil, status.Err()
	}

	return authors, nil
}
func (d *db) GetSuggestions(ctx context.Context, regex string, limit int64) ([]*Author, error) {
	response, err := d.collection.Find(ctx, bson.M{"name": bson.M{"$regex": regex, "$options": "i"}}, &options.FindOptions{Limit: &limit})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	authors := make([]*Author, 0)
	err = response.All(ctx, &authors)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(authors) == 0 {
		return nil, status.Error(codes.NotFound, "no authors found")
	}
	return authors, nil
}
