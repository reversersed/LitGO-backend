package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *RabbitService) InitiateBookRatingChangedReceiver() error {
	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Errorf("error while opening channel: %v", err)
		return err
	}
	r.channels = append(r.channels, ch)

	queue, err := ch.QueueDeclare(bookRatingChangedExchange, false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error while declaring queue: %v", err)
		return err
	}
	err = ch.ExchangeDeclare(bookRatingChangedExchange, "fanout", false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error while declaring exchange: %v", err)
		return err
	}

	err = ch.QueueBind(queue.Name, "#", bookRatingChangedExchange, false, nil)
	if err != nil {
		r.logger.Errorf("error while binding queue: %v", err)
		return err
	}

	messages, err := ch.Consume(queue.Name, "BookAPI", true, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error consuming from channel: %v", err)
		return err
	}
	go func(channel *amqp091.Channel) {
		for message := range messages {
			if channel.IsClosed() || r.conn.IsClosed() {
				return
			}
			r.logger.Info("received rating changed message")
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			var model RatingChangeModel
			if err := json.Unmarshal(message.Body, &model); err != nil {
				r.logger.Errorf("error while unmarshalling message body: %v", err)
			} else {
				if id, err := primitive.ObjectIDFromHex(model.BookId); err != nil {
					r.logger.Errorf("error while creating id from %s: %v", model.BookId, err)
				} else {
					err = r.storage.ChangeBookRating(ctx, id, model.Rating, model.TotalReviews)
					if err != nil {
						r.logger.Errorf("error while updating rating: %v", err)
					}
					r.cache.Delete([]byte("book_" + model.BookId))
				}
			}
		}
	}(ch)
	r.logger.Infof("Waiting for rating changing...")
	return nil
}
