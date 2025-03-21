package storage

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/reversersed/LitGO-backend-pkg/mongo"
	shared_pb "github.com/reversersed/LitGO-proto/gen/go/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
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
	defer db.seedGenres()
	return db
}
func (d *db) seedGenres() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documents, _ := d.collection.CountDocuments(ctx, bson.D{})
	if documents > 0 {
		d.logger.Infof("there are %d documents in database, seed canceled", documents)
		return
	}

	wg := sync.WaitGroup{}
	for _, c := range mocked_categories {
		wg.Add(1)
		go func(c struct {
			Name   string
			Genres []string
		}, wg *sync.WaitGroup) {
			category, err := d.InsertCategory(ctx, c.Name)
			if err != nil {
				d.logger.Fatalf("error inserting category: %v", err)
			}
			d.logger.Infof("category %v inserted", category)
			group := sync.WaitGroup{}
			for _, g := range c.Genres {
				group.Add(1)
				go func(g string, group *sync.WaitGroup) {
					genre, err := d.InsertGenre(ctx, category.Id, g)
					if err != nil {
						d.logger.Fatalf("error inserting genre: %v", err)
					}
					d.logger.Infof("genre %v inserted", genre)
					group.Done()
				}(g, &group)
			}
			group.Wait()
			wg.Done()
		}(c, &wg)
	}
	wg.Wait()
	d.logger.Info("categories seeded")
}
func (d *db) InsertGenre(ctx context.Context, category primitive.ObjectID, genreName string) (*Genre, error) {
	genreName = strings.TrimSpace(genreName)
	genre := &Genre{
		Id:   primitive.NewObjectID(),
		Name: genreName,
	}
	genre.TranslitName = mongo.GenerateTranslitName(genre.Name, genre.Id)

	insertRequest := bson.M{"$push": bson.M{"genres": genre}}
	result, err := d.collection.UpdateByID(ctx, category, insertRequest)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if result.ModifiedCount == 0 {
		return nil, status.Error(codes.Internal, "no documents was modified")
	}

	return genre, nil
}
func (d *db) InsertCategory(ctx context.Context, categoryName string) (*Category, error) {
	categoryName = strings.TrimSpace(categoryName)
	category := &Category{
		Id:     primitive.NewObjectID(),
		Name:   categoryName,
		Genres: []*Genre{},
	}
	category.TranslitName = mongo.GenerateTranslitName(categoryName, category.Id)

	result, err := d.collection.InsertOne(ctx, category)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, "error inserting id")
	} else if id != category.Id {
		var detail protoadapt.MessageV1 = &shared_pb.ErrorDetail{
			Field:       "Id",
			Struct:      "Category",
			Description: ("wanted id: " + category.Id.Hex()),
			Actualvalue: id.Hex(),
		}
		status, _ := status.New(codes.Internal, "error retrieving id").WithDetails(detail)
		return nil, status.Err()
	}
	return category, nil
}
func (d *db) GetAll(ctx context.Context) ([]*Category, error) {
	cursor, err := d.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var response []*Category

	err = cursor.All(ctx, &response)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(response) == 0 {
		return nil, status.Error(codes.NotFound, "there is no genres in database")
	}
	return response, nil
}
func (d *db) FindCategoryTree(ctx context.Context, query string) (*Category, error) {
	id, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		id = primitive.NilObjectID
	}
	filter := bson.M{
		"$or": []bson.M{{
			"genres": bson.M{
				"$elemMatch": bson.M{
					"$or": []bson.M{
						{"_id": id},
						{"translit": query},
					},
				},
			}},
			{"$or": []bson.M{{"_id": id}, {"translit": query}}},
		},
	}

	var result Category
	if err := d.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongodb.ErrNoDocuments) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &result, nil
}

func (d *db) IncreateBookCount(ctx context.Context, genre primitive.ObjectID) error {
	result, err := d.collection.UpdateOne(ctx, bson.M{"genres._id": genre}, bson.M{"$inc": bson.M{"genres.$.bookCount": 1}})
	if err != nil {
		return status.Error(codes.NotFound, err.Error())
	}
	if result.MatchedCount == 0 {
		return status.Error(codes.NotFound, "no genre found")
	}
	if result.ModifiedCount == 0 {
		return status.Error(codes.Aborted, "no genre modified")
	}
	return nil
}
