package config

import "os"

func Address() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ":8000"
}

func RabbitMQURL() string {
	return os.Getenv("RABBITMQ_URL")
}

func PostgreSQLURL() string {
	return os.Getenv("POSTGRES_URL")
}

func RedisURL() string {
	return os.Getenv("REDIS_URL")
}
