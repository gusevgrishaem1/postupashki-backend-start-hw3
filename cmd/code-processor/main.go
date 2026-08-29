package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	registry := prometheus.NewRegistry()
	go serveMetrics(registry)
	processor := service.NewProcessor(queue, codeRunner, httpRepository.NewCommitter(config.TaskServiceURL()), registry)
	if err := processor.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func serveMetrics(gatherer prometheus.Gatherer) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(config.MetricsAddress(), mux); err != nil {
		log.Printf("serve metrics: %v", err)
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
