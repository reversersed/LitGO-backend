package storage

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/golang/mock/gomock"
	mock_storage "github.com/reversersed/go-grpc/tree/main/api_user/internal/storage/mocks"
	"github.com/reversersed/go-grpc/tree/main/api_user/pkg/mongo"
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
			Name:         "user_mongo",
			Image:        "mongo",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
			Env: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "root",
				"MONGO_INITDB_ROOT_PASSWORD": "root",
			},
			SkipReaper: true,
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.PortBindings = nat.PortMap{"27017/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "54003"}}}
			},
		}
		mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
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
func TestFindUserById(t *testing.T) {
	ctx := context.Background()
	db, err := mongo.NewClient(context.Background(), cfg)
	defer db.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	table := []struct {
		Name          string
		Id            string
		ExceptedUser  *User
		ExceptedError string
	}{
		{
			Name:         "admin found",
			Id:           seedAdmin.Id.Hex(),
			ExceptedUser: seedAdmin,
		},
		{
			Name:          "wrong id type",
			Id:            "not an id",
			ExceptedError: "the provided hex string is not a valid ObjectID",
		},
		{
			Name:          "user does not exist",
			Id:            primitive.NewObjectID().Hex(),
			ExceptedError: "mongo: no documents in result",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			storage := NewStorage(db, cfg.Base, logger)
			user, err := storage.FindById(context.Background(), v.Id)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.ExceptedUser, user)
		})
	}
}
func TestFindUserByLogin(t *testing.T) {
	ctx := context.Background()
	db, err := mongo.NewClient(context.Background(), cfg)
	defer db.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	table := []struct {
		Name          string
		Login         string
		ExceptedUser  *User
		ExceptedError string
	}{
		{
			Name:         "admin found",
			Login:        seedAdmin.Login,
			ExceptedUser: seedAdmin,
		},
		{
			Name:          "user does not exist",
			Login:         "notExistingUser123!",
			ExceptedError: "mongo: no documents in result",
		},
		{
			Name:          "empty login",
			Login:         "",
			ExceptedError: "mongo: no documents in result",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			storage := NewStorage(db, cfg.Base, logger)
			user, err := storage.FindByLogin(context.Background(), v.Login)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.ExceptedUser, user)
		})
	}
}
func TestFindUserByEmail(t *testing.T) {
	ctx := context.Background()
	db, err := mongo.NewClient(context.Background(), cfg)
	defer db.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	table := []struct {
		Name          string
		Email         string
		ExceptedUser  *User
		ExceptedError string
	}{
		{
			Name:         "admin found",
			Email:        seedAdmin.Email,
			ExceptedUser: seedAdmin,
		},
		{
			Name:          "user does not exist",
			Email:         "not real email",
			ExceptedError: "mongo: no documents in result",
		},
		{
			Name:          "empty email",
			Email:         "",
			ExceptedError: "mongo: no documents in result",
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			storage := NewStorage(db, cfg.Base, logger)
			user, err := storage.FindByEmail(context.Background(), v.Email)

			if len(v.ExceptedError) == 0 {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, v.ExceptedError)
			}
			assert.Equal(t, v.ExceptedUser, user)
		})
	}
}
func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	db, err := mongo.NewClient(context.Background(), cfg)
	defer db.Client().Disconnect(ctx)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	logger := mock_storage.NewMocklogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	table := []struct {
		Name  string
		Model *User
	}{
		{
			Name: "created user",
			Model: &User{
				Login: "newUser",
				Email: "user@new.com",
				Roles: []string{"user"},
			},
		},
		{
			Name: "user with id",
			Model: &User{
				Id:    primitive.NewObjectID(),
				Login: "newUserWithId",
				Email: "user@id.com",
				Roles: []string{"user"},
			},
		},
	}

	for _, v := range table {
		t.Run(v.Name, func(t *testing.T) {
			storage := NewStorage(db, cfg.Base, logger)
			id, err := storage.CreateUser(context.Background(), v.Model)
			assert.NoError(t, err)

			result, err := storage.FindById(context.Background(), id.Hex())
			assert.NoError(t, err)
			v.Model.Id = result.Id

			assert.Equal(t, v.Model, result)
		})
	}
}
