package main

import (
	"context"
	"log"
	"time"

	"postupashki-backend-start-hw3/internal/code-processor/config"
	httpRepository "postupashki-backend-start-hw3/internal/code-processor/repository/http"
	"postupashki-backend-start-hw3/internal/code-processor/repository/rabbitmq"
	"postupashki-backend-start-hw3/internal/code-processor/usecases/runner"
	"postupashki-backend-start-hw3/internal/code-processor/usecases/service"
)

func main() {
	queue := connect()
	defer queue.Close()

	codeRunner, err := runner.New()
	if err != nil {
		log.Fatal(err)
	}
	defer codeRunner.Close()

	processor := service.NewProcessor(queue, codeRunner, httpRepository.NewCommitter(config.TaskServiceURL()))
	if err := processor.Run(context.Background()); err != nil {
		log.Fatal(err)
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
