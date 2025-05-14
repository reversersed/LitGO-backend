package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (r *RabbitService) SendBookRatingChangedMessage(ctx context.Context, model *RatingChangeModel) error {
	if model == nil {
		r.logger.Errorf("received nil model: %v", model)
		return errors.New("received nil model")
	}
	channel, err := r.conn.Channel()
	if err != nil {
		r.logger.Errorf("error opening channel: %v", err)
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(bookRatingChangedExchange, false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error creating queue: %v", err)
		return err
	}

	err = channel.ExchangeDeclare(bookRatingChangedExchange, "fanout", false, false, false, false, nil)
	if err != nil {
		r.logger.Errorf("error creating exchange: %v", err)
		return err
	}

	err = channel.QueueBind(queue.Name, "#", bookRatingChangedExchange, false, nil)
	if err != nil {
		r.logger.Errorf("error binding queue to exchange: %v", err)
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 500*time.Second)
	defer cancel()

	body, err := json.Marshal(model)
	if err != nil {
		return err
	}
	err = channel.PublishWithContext(ctx, bookRatingChangedExchange, "#", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		r.logger.Errorf("error publishing message: %v", err)
		return err
	}
	r.logger.Infof("sended book rating changed message: %v", model)
	return nil
}
