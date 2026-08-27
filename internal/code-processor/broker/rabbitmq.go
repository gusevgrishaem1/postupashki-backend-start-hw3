package broker

import amqp "github.com/rabbitmq/amqp091-go"

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

func (b *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	if err := b.channel.Qos(1, 0, false); err != nil {
		return nil, err
	}
	return b.channel.Consume(queue, "", false, false, false, false, nil)
}

func (b *RabbitMQ) Close() error {
	_ = b.channel.Close()
	return b.connection.Close()
}
