package storage

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/reversersed/LitGO-backend-pkg/mongo"
	mock_storage "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/storage/mocks"
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
	os.Exit(m.Run())
}
func TestFindBook(t *testing.T) {
	ctx := context.Background()
	dba, err := mongo.NewClient(context.Background(), cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	book, err := storage.CreateBook(ctx, &Book{Name: "Книга о книгопечатании", MonthPurchases: 20, Published: 1, Description: "Описание книги", Picture: "picture.png", Filepath: "book.epub", Genre: primitive.NewObjectID(), Authors: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}})
	assert.NoError(t, err)
	book2, err := storage.CreateBook(ctx, &Book{Name: "Треш какой-то", MonthPurchases: 60, Published: 0, Description: "Описание книги", Picture: "picture.png", Filepath: "book.epub", Genre: primitive.NewObjectID(), Authors: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}})
	assert.NoError(t, err)

	sugg, err := storage.Find(ctx, "(Книга)|(книгопечатании)", 1, 0, 0.0, Popular)

	assert.NoError(t, err)
	assert.Len(t, sugg, 1)
	assert.Equal(t, book, sugg[0])

	_, err = storage.Find(ctx, "(Книга)|(книгопечатании)", 1, 0, 2.0, Popular)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = no books found")

	_, err = storage.Find(ctx, "(КнигиНеСуществует)", 1, 0, 0.0, Newest)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = no books found")

	sugg, err = storage.Find(ctx, "(.*?)", 2, 0, 0.0, Popular)
	assert.NoError(t, err)
	assert.Len(t, sugg, 2)
	assert.Equal(t, book2, sugg[0])
	assert.Equal(t, book, sugg[1])

	sugg, err = storage.Find(ctx, "(.*?)", 2, 0, 0.0, Newest)
	assert.NoError(t, err)
	assert.Len(t, sugg, 2)
	assert.Equal(t, book2, sugg[1])
	assert.Equal(t, book, sugg[0])

	_, err = storage.Find(ctx, "(.*?)", 2, 0, 0.0, SortType("not-existing-sort-type"))
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = not known sort type: not-existing-sort-type")
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
func TestGetBookByGenre(t *testing.T) {
	ctx := context.Background()
	dba, err := mongo.NewClient(context.Background(), cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	genreIds := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	models := []*Book{{Name: "Эрагон 2.0", Rating: 2.0, Genre: genreIds[0]}, {Name: "Эрагон 4.0", Rating: 4.0, Genre: genreIds[0]}, {Name: "Эрагон 4.0", Rating: 4.0, Genre: genreIds[1]}}

	for i := 0; i < len(models); i++ {
		book, err := storage.CreateBook(ctx, models[i])
		assert.NoError(t, err)

		models[i] = book
	}

	t.Run("search all books", func(t *testing.T) {
		books, err := storage.GetBookByGenre(ctx, genreIds, Popular, false, 5, 0)
		if assert.NoError(t, err) {
			assert.Equal(t, models, books)
		}
	})
	t.Run("search only high rating", func(t *testing.T) {
		books, err := storage.GetBookByGenre(ctx, genreIds, Popular, true, 5, 0)
		if assert.NoError(t, err) {
			assert.Equal(t, models[1:], books)
		}
	})
	t.Run("search only one genre", func(t *testing.T) {
		books, err := storage.GetBookByGenre(ctx, []primitive.ObjectID{genreIds[1]}, Newest, false, 5, 0)
		if assert.NoError(t, err) {
			assert.Equal(t, []*Book{models[2]}, books)
		}
	})
	t.Run("empty array passed", func(t *testing.T) {
		_, err := storage.GetBookByGenre(ctx, []primitive.ObjectID{}, Popular, false, 5, 0)
		assert.Error(t, err)
	})
	t.Run("invalid sort type passed", func(t *testing.T) {
		_, err := storage.GetBookByGenre(ctx, []primitive.ObjectID{primitive.NewObjectID()}, SortType(""), false, 5, 0)
		assert.Error(t, err)
	})
	t.Run("not found error", func(t *testing.T) {
		books, err := storage.GetBookByGenre(ctx, []primitive.ObjectID{primitive.NewObjectID()}, Newest, false, 5, 0)
		assert.Error(t, err)
		assert.Nil(t, books)
	})
	t.Run("not found by page", func(t *testing.T) {
		books, err := storage.GetBookByGenre(ctx, []primitive.ObjectID{primitive.NewObjectID()}, Newest, false, 10, 5)
		assert.Error(t, err)
		assert.Nil(t, books)
	})
}
func TestGetBookList(t *testing.T) {
	ctx := context.Background()
	dba, err := mongo.NewClient(context.Background(), cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	genreId := primitive.NewObjectID()

	books := []Book{{Name: "book1", Genre: genreId}, {Name: "book2", Genre: genreId}, {Name: "book3", Genre: genreId}}
	for _, b := range books {
		_, err := storage.CreateBook(ctx, &b)
		assert.NoError(t, err)
	}

	mocked_books, err := storage.GetBookByGenre(ctx, []primitive.ObjectID{genreId}, Popular, false, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, mocked_books, len(books))

	_, e := storage.GetBookList(ctx, []primitive.ObjectID{primitive.NewObjectID()}, []string{})
	assert.EqualError(t, e, "rpc error: code = NotFound desc = no books found")

	storage = NewStorage(dba, cfg.Base, logger)

	a, e := storage.GetBookList(ctx, []primitive.ObjectID{mocked_books[0].Id, mocked_books[1].Id}, []string{mocked_books[2].TranslitName})
	assert.NoError(t, e)

	assert.Equal(t, mocked_books, a)

	_, e = storage.GetBookList(ctx, []primitive.ObjectID{}, []string{})
	assert.EqualError(t, e, "rpc error: code = InvalidArgument desc = no id or translit name argument presented")

}
