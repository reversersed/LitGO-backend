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
	mock_storage "github.com/reversersed/LitGO-backend/tree/main/api_genre/internal/storage/mocks"
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
func TestGetAll(t *testing.T) {
	ctx := context.Background()
	dba, err := mongo.NewClient(context.Background(), cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := &db{
		collection: dba.Collection(cfg.Base),
		logger:     logger,
	}
	_, err = storage.GetAll(ctx)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = there is no genres in database")

	storage = NewStorage(dba, cfg.Base, logger)

	cats, err := storage.GetAll(ctx)
	assert.NoError(t, err)

	assert.Len(t, cats, len(mocked_categories))
}
func TestFindCategoryTree(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dba, err := mongo.NewClient(ctx, cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	cats, err := storage.GetAll(ctx)
	assert.NoError(t, err)

	result, err := storage.FindCategoryTree(ctx, cats[0].Genres[0].TranslitName)
	if assert.NoError(t, err) {
		assert.Equal(t, cats[0], result)
	}

	result, err = storage.FindCategoryTree(ctx, cats[0].Genres[2].Id.Hex())
	if assert.NoError(t, err) {
		assert.Equal(t, cats[0], result)
	}

	result, err = storage.FindCategoryTree(ctx, cats[0].TranslitName)
	if assert.NoError(t, err) {
		assert.Equal(t, cats[0], result)
	}

	result, err = storage.FindCategoryTree(ctx, cats[1].Id.Hex())
	if assert.NoError(t, err) {
		assert.Equal(t, cats[1], result)
	}

	result, err = storage.FindCategoryTree(ctx, "not existing")
	assert.Nil(t, result)
	assert.Error(t, err)
}
func TestIncreaseBookCount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dba, err := mongo.NewClient(ctx, cfg)
	defer dba.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	storage := NewStorage(dba, cfg.Base, logger)

	genres, err := storage.GetAll(ctx)
	assert.NoError(t, err)

	count := genres[0].Genres[0].BookCount
	err = storage.IncreateBookCount(ctx, genres[0].Genres[0].Id)
	assert.NoError(t, err)

	genres, err = storage.GetAll(ctx)
	assert.NoError(t, err)

	assert.Equal(t, count+1, genres[0].Genres[0].BookCount)

	err = storage.IncreateBookCount(ctx, primitive.NewObjectID())
	assert.Error(t, err)
}
