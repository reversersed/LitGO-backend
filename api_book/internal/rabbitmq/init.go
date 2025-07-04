package rabbitmq

import (
	"context"
	"io"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=init.go -destination=mocks/rabbit.go
type logger interface {
	Info(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}
type storage interface {
	ChangeBookRating(ctx context.Context, bookId primitive.ObjectID, rating float64, totalReviews int) error
}
type cache interface {
	Delete([]byte) bool
}
type RabbitService struct {
	cache    cache
	conn     *amqp.Connection
	logger   logger
	storage  storage
	channels []io.Closer
}

func New(connection *amqp.Connection, logger logger, storage storage, cache cache) *RabbitService {
	return &RabbitService{
		conn:    connection,
		logger:  logger,
		storage: storage,
		cache:   cache,
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
