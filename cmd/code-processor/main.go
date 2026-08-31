package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	queue, err := connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer queue.Close()

	codeRunner, err := runner.New()
	if err != nil {
		log.Printf("create runner: %v", err)
		return
	}
	defer codeRunner.Close()

	registry := prometheus.NewRegistry()
	metricsServer := newMetricsServer(registry)
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("serve metrics: %v", err)
		}
	}()
	processor := service.NewProcessor(queue, codeRunner, httpRepository.NewCommitter(config.TaskServiceURL()), registry)
	if err := processor.Run(ctx); err != nil && ctx.Err() == nil {
		log.Printf("run processor: %v", err)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown metrics server: %v", err)
	}
}

func newMetricsServer(gatherer prometheus.Gatherer) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	return &http.Server{
		Addr:              config.MetricsAddress(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func connect(ctx context.Context) (*rabbitmq.RabbitMQ, error) {
	for {
		queue, err := rabbitmq.NewRabbitMQ(config.RabbitMQURL())
		if err == nil {
			return queue, nil
		}
		log.Printf("connect to RabbitMQ: %v", err)
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
