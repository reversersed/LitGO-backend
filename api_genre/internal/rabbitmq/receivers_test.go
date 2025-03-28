package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/reversersed/LitGO-backend-pkg/rabbitmq"
	mock_rabbitmq "github.com/reversersed/LitGO-backend/tree/main/api_genre/internal/rabbitmq/mocks"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	genres_pb "github.com/reversersed/LitGO-proto/gen/go/genres"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

var conn *amqp091.Connection

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		log.Println("\t--- Integration tests are not running in short mode")
		return
	}
	var err error
	var cont testcontainers.Container
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		req := testcontainers.ContainerRequest{
			Image:        "rabbitmq:3.10.7-management",
			ExposedPorts: []string{"5672/tcp"},
			SkipReaper:   true,
			WaitingFor:   wait.ForLog("Server startup complete"),
			Env:          map[string]string{"RABBITMQ_DEFAULT_USER": "user", "RABBITMQ_DEFAULT_PASS": "password"},
		}
		cont, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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
		log.Fatal(err)
	}
	defer func() {
		err := cont.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
	port, err := cont.MappedPort(ctx, "5672")
	if err != nil {
		log.Fatal(err)
	}
	cfg := &rabbitmq.RabbitConfig{Rabbit_Port: port.Port(), Rabbit_Pass: "password", Rabbit_User: "user"}
	cfg.Rabbit_Host, err = cont.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client, err := rabbitmq.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	conn = client.Connection
	defer func() {
		client.Close()
	}()

	os.Exit(m.Run())
}
func TestBookCreatedReceiver(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_rabbitmq.NewMocklogger(ctrl)
	storage := mock_rabbitmq.NewMockstorage(ctrl)
	service := New(conn, logger, storage)
	defer service.Close()

	logger.EXPECT().Infof("Waiting for created books...")
	err := service.InitiateBookCreatedReceiver()
	assert.NoError(t, err)

	t.Run("empty body", func(t *testing.T) {
		ctx := context.Background()
		channel, err := conn.Channel()
		assert.NoError(t, err)

		done := make(chan bool, 1)

		logger.EXPECT().Info("received created book message")
		logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Do(func(_ any, _ ...any) { done <- true })
		err = channel.PublishWithContext(ctx, bookCreatedExchange, "#", false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte{},
		})
		assert.NoError(t, err)
		err = channel.Close()
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.FailNow(t, "service not returning error message within 5 second")
		}
	})

	t.Run("wrong id received", func(t *testing.T) {
		ctx := context.Background()
		channel, err := conn.Channel()
		assert.NoError(t, err)

		done := make(chan bool, 1)

		logger.EXPECT().Info("received created book message")
		logger.EXPECT().Infof(gomock.Any(), gomock.Any())
		logger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(_ any, _ ...any) { done <- true })

		book := &books_pb.BookModel{Genre: &genres_pb.GenreModel{Id: "bad id"}}
		body, err := json.Marshal(book)
		assert.NoError(t, err)
		err = channel.PublishWithContext(ctx, bookCreatedExchange, "#", false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		assert.NoError(t, err)
		err = channel.Close()
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.FailNow(t, "service not returning error message within 5 second")
		}
	})

	t.Run("error from storage received", func(t *testing.T) {
		ctx := context.Background()
		channel, err := conn.Channel()
		assert.NoError(t, err)

		done := make(chan bool, 1)

		logger.EXPECT().Info("received created book message")
		logger.EXPECT().Infof(gomock.Any(), gomock.Any())
		logger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Do(func(_ any, _ ...any) { done <- true })

		id := primitive.NewObjectID()
		book := &books_pb.BookModel{Genre: &genres_pb.GenreModel{Id: id.Hex()}}
		body, err := json.Marshal(book)
		assert.NoError(t, err)
		err = channel.PublishWithContext(ctx, bookCreatedExchange, "#", false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		storage.EXPECT().IncreateBookCount(gomock.Any(), id).Return(errors.New("error produced"))

		assert.NoError(t, err)
		err = channel.Close()
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.FailNow(t, "service not returning error message within 5 second")
		}
	})

	t.Run("successed receiver", func(t *testing.T) {
		ctx := context.Background()
		channel, err := conn.Channel()
		assert.NoError(t, err)

		done := make(chan bool, 1)

		logger.EXPECT().Info("received created book message")
		logger.EXPECT().Infof(gomock.Any(), gomock.Any())

		id := primitive.NewObjectID()
		book := &books_pb.BookModel{Genre: &genres_pb.GenreModel{Id: id.Hex()}}
		body, err := json.Marshal(book)
		assert.NoError(t, err)
		err = channel.PublishWithContext(ctx, bookCreatedExchange, "#", false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		storage.EXPECT().IncreateBookCount(gomock.Any(), id).DoAndReturn(func(_ any, rid primitive.ObjectID) error {
			if assert.Equal(t, id, rid) {
				done <- true
			}
			return nil
		})

		assert.NoError(t, err)
		err = channel.Close()
		assert.NoError(t, err)

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			assert.FailNow(t, "storage not received call or id was incorrect")
		}
	})
}
