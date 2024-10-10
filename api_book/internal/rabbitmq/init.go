package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

//go:generate mockgen -source=init.go -destination=mocks/rabbit.go
type logger interface {
	Info(...any)
	Infof(string, ...any)
	Errorf(string, ...any)
}
type storage interface {
}
type RabbitService struct {
	conn    *amqp.Connection
	logger  logger
	storage storage
}

func New(connection *amqp.Connection, logger logger, storage storage) *RabbitService {
	return &RabbitService{
		conn:    connection,
		logger:  logger,
		storage: storage,
	}
}

func (s *RabbitService) Close() error {
	return s.conn.Close()
}
