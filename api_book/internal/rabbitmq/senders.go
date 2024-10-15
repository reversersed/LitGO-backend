package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rabbitmq/amqp091-go"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
)

func (r *RabbitService) SendBookCreatedMessage(ctx context.Context, book *books_pb.BookModel) error {
	if book == nil {
		r.logger.Errorf("received nil book: %v", book)
		return errors.New("received nil book")
	}
	channel, err := r.conn.Channel()
	if err != nil {
		r.logger.Errorf("error opening channel: %v", err)
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(bookCreatedQueue, false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error creating queue: %v", err)
		return err
	}

	err = channel.ExchangeDeclare(bookCreatedExchange, "fanout", false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error creating exchange: %v", err)
		return err
	}

	err = channel.QueueBind(queue.Name, "#", bookCreatedExchange, false, nil)
	if err != nil {
		r.logger.Errorf("error binding queue to exchange: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
	defer cancel()

	body, err := json.Marshal(book)
	if err != nil {
		return err
	}
	err = channel.PublishWithContext(ctx, bookCreatedExchange, "#", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		r.logger.Errorf("error publishing message: %v", err)
		return err
	}
	r.logger.Infof("sended book created message")
	return nil
}
