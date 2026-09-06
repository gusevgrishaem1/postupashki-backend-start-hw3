package config

import (
	"fmt"
	"os"
	"strconv"
)

func Address() (string, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65_535 {
		return "", fmt.Errorf("port must be an integer between 0 and 65535")
	}

	return fmt.Sprintf(":%d", port), nil
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
