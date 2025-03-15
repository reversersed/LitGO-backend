package rabbitmq

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/golang/mock/gomock"
	"github.com/rabbitmq/amqp091-go"
	"github.com/reversersed/LitGO-backend-pkg/rabbitmq"
	mock_rabbitmq "github.com/reversersed/LitGO-backend/tree/main/api_book/internal/rabbitmq/mocks"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			WaitingFor:   wait.ForLog("Server startup complete"),
			Env:          map[string]string{"RABBITMQ_DEFAULT_USER": "user", "RABBITMQ_DEFAULT_PASS": "password"},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.AutoRemove = true
				hc.PortBindings = nat.PortMap{"5672/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "54004"}}}
			},
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
func TestBookCreatedSender(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mock_rabbitmq.NewMocklogger(ctrl)
	service := New(conn, logger, nil)
	defer service.Close()

	ctx := context.Background()
	logger.EXPECT().Errorf(gomock.Any(), gomock.Any())
	err := service.SendBookCreatedMessage(ctx, nil)
	assert.Error(t, err)

	logger.EXPECT().Infof("sended book created message")
	book := &books_pb.BookModel{Name: "книга", Genre: &books_pb.GenreModel{Id: primitive.NewObjectID().Hex(), Name: "жанр"}}
	err = service.SendBookCreatedMessage(ctx, book)
	assert.NoError(t, err)

	ch, err := conn.Channel()
	assert.NoError(t, err)

	queue, err := ch.QueueDeclare(bookCreatedQueue, false, false, false, false, nil)
	assert.NoError(t, err)

	err = ch.QueueBind(queue.Name, "#", bookCreatedExchange, false, nil)
	assert.NoError(t, err)

	messages, err := ch.Consume(queue.Name, "TestingAPI", true, false, false, false, nil)
	assert.NoError(t, err)

	data := <-messages
	var received books_pb.BookModel

	err = json.Unmarshal(data.Body, &received)
	assert.NoError(t, err)

	assert.Equal(t, book, &received)
}
