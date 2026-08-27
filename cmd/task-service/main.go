package main

import (
	"log"
	"net/http"
	"postupashki-backend-start-hw3/internal/task-service/broker"
	"time"

	"postupashki-backend-start-hw3/internal/task-service/config"
	"postupashki-backend-start-hw3/internal/task-service/http"
	"postupashki-backend-start-hw3/internal/task-service/repository/inmemory"
	"postupashki-backend-start-hw3/internal/task-service/usecases/service"
)

// @title Task service
// @version 1.0
// @description Asynchronous task processing service.
// @host localhost:8000
// @BasePath /
func main() {
	queue := connect()
	defer queue.Close()

	taskRepository := inmemory.NewTask()
	userRepository := inmemory.NewUser()
	sessionRepository := inmemory.NewSession()

	taskService := service.NewTask(taskRepository, queue)
	authService := service.NewAuth(userRepository, sessionRepository)

	server := taskhttp.NewServer(taskService, authService)

	log.Fatal(http.ListenAndServe(config.Address(), server.Handler()))
}

func connect() *broker.RabbitMQ {
	for {
		queue, err := broker.NewRabbitMQ(config.RabbitMQURL())
		if err == nil {
			return queue
		}
		log.Printf("connect to RabbitMQ: %v", err)
		time.Sleep(2 * time.Second)
	}
}
