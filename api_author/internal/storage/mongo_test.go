package storage

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	mock_storage "github.com/reversersed/go-grpc/tree/main/api_author/internal/storage/mocks"
	"github.com/reversersed/go-grpc/tree/main/api_author/pkg/mongo"
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
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017/tcp"),
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
		},
	}
	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
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
func TestGetAuthors(t *testing.T) {
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
	_, e := storage.GetAuthors(ctx, []primitive.ObjectID{mocked_authors[0].Id}, []string{})
	assert.EqualError(t, e, "rpc error: code = NotFound desc = no authors found")

	storage = NewStorage(dba, cfg.Base, logger)

	a, e := storage.GetAuthors(ctx, []primitive.ObjectID{mocked_authors[0].Id, mocked_authors[1].Id}, []string{mocked_authors[2].TranslitName})
	assert.NoError(t, e)

	assert.Equal(t, mocked_authors, a)

	_, e = storage.GetAuthors(ctx, []primitive.ObjectID{}, []string{})
	assert.EqualError(t, e, "rpc error: code = InvalidArgument desc = no id or translit name argument presented")
}
