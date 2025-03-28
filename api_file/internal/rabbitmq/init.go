package rabbitmq

import (
	"io"

	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate mockgen -source=init.go -destination=mocks/rabbit.go
type logger interface {
	Info(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}
type RabbitService struct {
	conn     *amqp.Connection
	logger   logger
	channels []io.Closer
}

func New(connection *amqp.Connection, logger logger) *RabbitService {
	return &RabbitService{
		conn:   connection,
		logger: logger,
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
