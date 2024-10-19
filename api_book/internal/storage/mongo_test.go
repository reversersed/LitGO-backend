package storage

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_storage "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage/mocks"
	"github.com/reversersed/LitGO-backend/tree/main/api_book/pkg/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var cfg *mongo.DatabaseConfig

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		log.Println("\t--- Integration tests are not running in short mode")
		return
	}

	ctx := context.Background()
	var err error
	var mongoContainer testcontainers.Container
	for i := 0; i < 5; i++ {
		req := testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
			Env: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "root",
				"MONGO_INITDB_ROOT_PASSWORD": "root",
			},
		}
		mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err == nil {
			break
		}
		log.Printf("failed to create container: %v, retry %d/5", err, i+1)
		<-time.After(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not start mongo: %s", err)
	}
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop mongo: %s", err)
		}
	}()
	host, err := mongoContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	port, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		log.Fatal(err)
	}
	cfg = &mongo.DatabaseConfig{
		Host:     host,
		Port:     port.Int(),
		User:     "root",
		Password: "root",
		Base:     "testbase",
		AuthDb:   "admin",
	}
	exit := m.Run()
	os.Exit(exit)
}
func TestGetSuggestion(t *testing.T) {
	ctx := context.Background()
	dba, err := mongo.NewClient(context.Background(), cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	book, err := storage.CreateBook(ctx, &Book{Name: "Книга о книгопечатании", Description: "Описание книги", Picture: "picture.png", Filepath: "book.epub", Genre: primitive.NewObjectID(), Authors: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}})
	assert.NoError(t, err)
	sugg, err := storage.GetSuggestions(ctx, "(Книга)|(книгопечатании)", 1)

	assert.NoError(t, err)
	assert.Len(t, sugg, 1)
	assert.Equal(t, book, sugg[0])

	_, err = storage.GetSuggestions(ctx, "(КнигиНеСуществует)", 1)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = no books found")
}
func TestGetBook(t *testing.T) {
	ctx := context.Background()
	dba, err := mongo.NewClient(context.Background(), cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	book, err := storage.CreateBook(ctx, &Book{Name: "Книга о книгопечатании", Description: "Описание книги", Picture: "picture.png", Filepath: "book.epub", Genre: primitive.NewObjectID(), Authors: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}})
	assert.NoError(t, err)

	t.Run("get by translit", func(t *testing.T) {
		response, err := storage.GetBook(ctx, book.TranslitName)
		if assert.NoError(t, err) {
			assert.Equal(t, book, response)
		}
	})
	t.Run("get by id", func(t *testing.T) {
		response, err := storage.GetBook(ctx, book.Id.Hex())
		if assert.NoError(t, err) {
			assert.Equal(t, book, response, book.Id.Hex())
		}
	})
	t.Run("not found error by translit", func(t *testing.T) {
		_, err := storage.GetBook(ctx, "not-found-name")
		assert.EqualError(t, err, "rpc error: code = NotFound desc = mongo: no documents in result")
	})
	t.Run("not found error by id", func(t *testing.T) {
		_, err := storage.GetBook(ctx, primitive.NewObjectID().Hex())
		assert.EqualError(t, err, "rpc error: code = NotFound desc = mongo: no documents in result")
	})
}
