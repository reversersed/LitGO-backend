package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mdigger/translit"
	shared_pb "github.com/reversersed/go-grpc/tree/main/api_genre/pkg/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type logger interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Warnf(string, ...interface{})
	Warn(...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
}
type db struct {
	sync.RWMutex
	logger     logger
	collection *mongo.Collection
}

func generateIntegerFromObjectId(id primitive.ObjectID) int {
	lastBytes := id[len(id)-3:]
	return int(lastBytes[0])<<16 | int(lastBytes[1])<<8 | int(lastBytes[2])
}
func NewStorage(storage *mongo.Database, collection string, logger logger) *db {
	db := &db{
		collection: storage.Collection(collection),
		logger:     logger,
	}
	defer db.seedGenres()
	return db
}
func (d *db) seedGenres() {
	d.Lock()
	defer d.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documents, _ := d.collection.CountDocuments(ctx, bson.D{})
	if documents > 0 {
		d.logger.Infof("there are %d documents in database, seed canceled", documents)
		return
	}
	categories := []struct {
		Name   string
		Genres []string
	}{
		{
			Name: "Бизнес-книги",
			Genres: []string{
				"Менеджмент",
				"Работа с клиентами",
				"Переговоры",
				"Ораторское искусство / риторика",
			},
		},
		{
			Name: "Знания и навыки",
			Genres: []string{
				"Научно-популярная литература",
				"Учебная и научная литература",
				"Компьютерная литература",
				"Культура и искусство",
				"Саморазвитие / личностный рост",
				"Эзотерика",
				"Словари, справочники",
			},
		},
		{
			Name: "Хобби, досуг",
			Genres: []string{
				"Отдых / туризм",
				"Хобби / увлечения",
				"Охота",
				"Мода и стиль",
				"Автомобили и ПДД",
				"Сад и огород",
				"Прикладная литература",
				"Развлечения",
				"Йога",
				"Кулинария",
			},
		},
		{
			Name: "Легкое чтение",
			Genres: []string{
				"Детективы",
				"Фантастика",
				"Фэнтези",
				"Любовные романы",
				"Ужасы / мистика",
				"Боевики, остросюжетная литература",
				"Юмористическая литература",
				"Приключения",
				"Легкая проза",
			},
		},
	}
	wg := sync.WaitGroup{}
	for _, c := range categories {
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
	rxSpaces := regexp.MustCompile(`\s+`)
	reg := regexp.MustCompile(`[^\p{L}\s]`)

	genreName = strings.TrimSpace(genreName)
	genre := &Genre{
		Id:   primitive.NewObjectID(),
		Name: genreName,
	}
	genre.TranslitName = fmt.Sprintf("%s-%d", strings.ReplaceAll(strings.TrimSpace(rxSpaces.ReplaceAllString(translit.Ru(reg.ReplaceAllString(strings.ToLower(genreName), "")), " ")), " ", "-"), generateIntegerFromObjectId(genre.Id))

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
	rxSpaces := regexp.MustCompile(`\s+`)
	reg := regexp.MustCompile(`[^\p{L}\s]`)

	categoryName = strings.TrimSpace(categoryName)
	category := &Category{
		Id:     primitive.NewObjectID(),
		Name:   categoryName,
		Genres: []*Genre{},
	}
	category.TranslitName = fmt.Sprintf("%s-%d", strings.ReplaceAll(strings.TrimSpace(rxSpaces.ReplaceAllString(translit.Ru(reg.ReplaceAllString(strings.ToLower(categoryName), "")), " ")), " ", "-"), generateIntegerFromObjectId(category.Id))

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
			Description: fmt.Sprintf("wanted id: %s", category.Id.Hex()),
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
	return response, nil
}
