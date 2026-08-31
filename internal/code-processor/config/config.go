package config

import "os"

func RabbitMQURL() string {
	return os.Getenv("RABBITMQ_URL")
}

func TaskServiceURL() string {
	return os.Getenv("TASK_SERVICE_URL")
}
