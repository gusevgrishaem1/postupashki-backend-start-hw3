package config

import "os"

func RabbitMQURL() string {
	return os.Getenv("RABBITMQ_URL")
}
