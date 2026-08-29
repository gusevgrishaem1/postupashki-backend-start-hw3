package main

import (
	"database/sql"
	"log"
	"net/http"
	"postupashki-backend-start-hw3/internal/task-service/api/http"
	"time"

	"postupashki-backend-start-hw3/internal/task-service/config"
	"postupashki-backend-start-hw3/internal/task-service/repository/postgres"
	"postupashki-backend-start-hw3/internal/task-service/repository/rabbitmq"
	redisRepository "postupashki-backend-start-hw3/internal/task-service/repository/redis"
	"postupashki-backend-start-hw3/internal/task-service/usecases/service"
)

func main() {
	queue := connect()
	defer queue.Close()
	database, sessions := connectStorage()
	defer database.Close()
	defer sessions.Close()

	taskService := service.NewTask(postgres.NewTask(database), queue)
	authService := service.NewAuth(postgres.NewUser(database), sessions)
	server := taskhttp.NewServer(taskService, authService)

	log.Fatal(http.ListenAndServe(config.Address(), server.Handler()))
}

func connectStorage() (*sql.DB, *redisRepository.Session) {
	for {
		database, databaseErr := postgres.Open(config.PostgreSQLURL())
		sessions, redisErr := redisRepository.NewSession(config.RedisURL())
		if databaseErr == nil && redisErr == nil {
			return database, sessions
		}
		if database != nil {
			_ = database.Close()
		}
		if sessions != nil {
			_ = sessions.Close()
		}
		log.Printf("connect to storage: postgres=%v redis=%v", databaseErr, redisErr)
		time.Sleep(2 * time.Second)
	}
}

func connect() *rabbitmq.RabbitMQ {
	for {
		queue, err := rabbitmq.NewRabbitMQ(config.RabbitMQURL())
		if err == nil {
			return queue
		}
		log.Printf("connect to RabbitMQ: %v", err)
		time.Sleep(2 * time.Second)
	}
}
