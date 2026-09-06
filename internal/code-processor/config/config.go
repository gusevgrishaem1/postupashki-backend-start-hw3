package config

import "os"

func RabbitMQURL() string {
	return os.Getenv("RABBITMQ_URL")
}

func TaskServiceURL() string {
	return os.Getenv("TASK_SERVICE_URL")
}

func MetricsAddress() string {
	if port := os.Getenv("METRICS_PORT"); port != "" {
		return ":" + port
	}
	return ":9000"
}
