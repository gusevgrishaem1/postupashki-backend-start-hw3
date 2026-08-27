package main

import (
	"log"
	"net/http"
	"time"

	"postupashki-backend-start-hw3/internal/task-service/broker"
	"postupashki-backend-start-hw3/internal/task-service/config"
	taskhttp "postupashki-backend-start-hw3/internal/task-service/http"
	"postupashki-backend-start-hw3/internal/task-service/repository/inmemory"
	"postupashki-backend-start-hw3/internal/task-service/usecases/service"
)

func main() {
	queue := connect()
	defer queue.Close()

	taskService := service.NewTask(inmemory.NewTask(), queue)
	authService := service.NewAuth(inmemory.NewUser(), inmemory.NewSession())
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
