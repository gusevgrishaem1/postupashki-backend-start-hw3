package rabbitmq

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"postupashki-backend-start-hw3/internal/contracts"
)

const queue = "code_submissions"

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		_ = connection.Close()
		return nil, err
	}
	if _, err = channel.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		_ = channel.Close()
		_ = connection.Close()
		return nil, err
	}
	return &RabbitMQ{connection: connection, channel: channel}, nil
}

func (b *RabbitMQ) Publish(ctx context.Context, submission contracts.Submission) error {
	body, err := json.Marshal(submission)
	if err != nil {
		return err
	}
	return b.channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType: "application/json", DeliveryMode: amqp.Persistent, Body: body,
	})
}

func (b *RabbitMQ) Close() error {
	_ = b.channel.Close()
	return b.connection.Close()
}
