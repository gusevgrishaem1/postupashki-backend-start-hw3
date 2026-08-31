package main

import (
	"log"
	"net/http"
	"time"

	taskhttp "postupashki-backend-start-hw3/internal/task-service/api/http"
	"postupashki-backend-start-hw3/internal/task-service/config"
	"postupashki-backend-start-hw3/internal/task-service/repository/inmemory"
	"postupashki-backend-start-hw3/internal/task-service/repository/rabbitmq"
	"postupashki-backend-start-hw3/internal/task-service/usecases/service"
)

func main() {
	addr, err := config.Address()
	if err != nil {
		log.Fatal(err)
	}

	queue := connect()
	defer queue.Close()

	taskService := service.NewTask(inmemory.NewTask(), queue)
	const sessionTTL = 24 * time.Hour
	authService := service.NewAuth(inmemory.NewUser(), inmemory.NewSession(), sessionTTL)
	defer authService.Close()
	server := taskhttp.NewServer(taskService, authService)

	log.Fatal(http.ListenAndServe(addr, server.Handler()))
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
