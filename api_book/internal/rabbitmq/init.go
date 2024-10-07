package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type logger interface {
	Info(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}
type RabbitService struct {
	conn   *amqp.Connection
	logger logger
}

func New(connection *amqp.Connection, logger logger) *RabbitService {
	return &RabbitService{
		conn:   connection,
		logger: logger,
	}
}

func (s *RabbitService) Close() error {
	return s.conn.Close()
}
