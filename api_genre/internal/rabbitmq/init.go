package rabbitmq

import (
	"context"
	"io"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO write tests for receivers
// go:generate mockgen -source=init.go -destination=mocks/rabbit.go
type logger interface {
	Info(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}
type storage interface {
	IncreateBookCount(context.Context, primitive.ObjectID) error
}
type RabbitService struct {
	conn     *amqp.Connection
	logger   logger
	storage  storage
	channels []io.Closer
}

func New(connection *amqp.Connection, logger logger, storage storage) *RabbitService {
	return &RabbitService{
		conn:    connection,
		logger:  logger,
		storage: storage,
	}
}

func (s *RabbitService) Close() error {
	for _, c := range s.channels {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return s.conn.Close()
}
