package rabbitmq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"

	"postupashki-backend-start-hw3/internal/code-processor/repository"
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

func (b *RabbitMQ) Consume() (<-chan repository.Delivery, error) {
	if err := b.channel.Qos(1, 0, false); err != nil {
		return nil, err
	}
	deliveries, err := b.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	result := make(chan repository.Delivery)
	go func() {
		defer close(result)
		for delivery := range deliveries {
			result <- rabbitDelivery{delivery: delivery}
		}
	}()
	return result, nil
}

func (b *RabbitMQ) Close() error {
	_ = b.channel.Close()
	return b.connection.Close()
}

type rabbitDelivery struct {
	delivery amqp.Delivery
}

func (d rabbitDelivery) Submission() (contracts.Submission, error) {
	var submission contracts.Submission
	err := json.Unmarshal(d.delivery.Body, &submission)
	return submission, err
}

func (d rabbitDelivery) Ack() error {
	return d.delivery.Ack(false)
}

func (d rabbitDelivery) Reject(requeue bool) error {
	return d.delivery.Reject(requeue)
}

func (d rabbitDelivery) Nack(requeue bool) error {
	return d.delivery.Nack(false, requeue)
}
