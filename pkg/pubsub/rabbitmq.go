package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alancesar/imgur-fetcher/pkg/media"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	RabbitMQPublisher struct {
		channel  *amqp.Channel
		exchange string
		key      string
	}

	Consumer func(ctx context.Context, m media.Media) error
)

func NewRabbitMQPublisher(connection *amqp.Connection, exchange, key string) (*RabbitMQPublisher, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQPublisher{
		channel:  channel,
		exchange: exchange,
		key:      key,
	}, nil
}

func (p RabbitMQPublisher) Publish(ctx context.Context, m media.Media) error {
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return p.channel.PublishWithContext(
		ctx,
		p.exchange,
		p.key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p RabbitMQPublisher) Close() error {
	return p.channel.Close()
}
