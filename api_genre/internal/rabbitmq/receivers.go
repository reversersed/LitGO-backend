package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	books_pb "github.com/reversersed/LitGO-proto/gen/go/books"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO write tests for receiver
func (r *RabbitService) InitiateBookCreatedReceived() error {
	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Errorf("error while opening channel: %v", err)
		return err
	}
	r.channels = append(r.channels, ch)

	queue, err := ch.QueueDeclare(bookCreatedQueue, false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error while declaring queue: %v", err)
		return err
	}
	err = ch.ExchangeDeclare(bookCreatedExchange, "fanout", false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error while declaring exchange: %v", err)
		return err
	}

	err = ch.QueueBind(queue.Name, "#", bookCreatedExchange, false, nil)
	if err != nil {
		r.logger.Errorf("error while binding queue: %v", err)
		return err
	}

	messages, err := ch.Consume(queue.Name, "GenreAPI", true, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error consuming from channel: %v", err)
		return err
	}
	go func(channel *amqp091.Channel) {
		for message := range messages {
			if channel.IsClosed() || r.conn.IsClosed() {
				return
			}
			r.logger.Info("received created book message")
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			var book books_pb.BookModel
			if err := json.Unmarshal(message.Body, &book); err != nil {
				r.logger.Errorf("error while unmarshalling message body: %v", err)
			} else {
				r.logger.Infof("unmarshalled book model: %v", book.String())
				if id, err := primitive.ObjectIDFromHex(book.GetGenre().GetId()); err != nil {
					r.logger.Errorf("error while creating id from %s: %v", book.GetGenre().GetId(), err)
				} else {
					err = r.storage.IncreateBookCount(ctx, id)
					if err != nil {
						r.logger.Errorf("error while increasing book count: %v", err)
					}
				}
			}
		}
	}(ch)
	r.logger.Infof("Waiting for created books...")
	return nil
}
