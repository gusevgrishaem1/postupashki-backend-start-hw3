package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	taskhttp "postupashki-backend-start-hw3/internal/task-service/api/http"
	"postupashki-backend-start-hw3/internal/task-service/config"
	"postupashki-backend-start-hw3/internal/task-service/repository/postgres"
	"postupashki-backend-start-hw3/internal/task-service/repository/rabbitmq"
	redisRepository "postupashki-backend-start-hw3/internal/task-service/repository/redis"
	"postupashki-backend-start-hw3/internal/task-service/usecases/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr, err := config.Address()
	if err != nil {
		log.Fatal(err)
	}
	queue, err := connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer queue.Close()
	database, sessions, err := connectStorage(ctx)
	if err != nil {
		log.Printf("connect to storage: %v", err)
		return
	}
	defer database.Close()
	defer sessions.Close()

	taskService := service.NewTask(postgres.NewTask(database), queue)
	authService := service.NewAuth(postgres.NewUser(database), sessions, 24*time.Hour)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           taskhttp.NewServer(taskService, authService).Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverError := make(chan error, 1)
	go func() { serverError <- httpServer.ListenAndServe() }()
	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("serve HTTP: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown HTTP server: %v", err)
		}
	}
}

func connectStorage(ctx context.Context) (*sql.DB, *redisRepository.Session, error) {
	for {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		database, databaseErr := postgres.Open(attemptCtx, config.PostgreSQLURL())
		sessions, redisErr := redisRepository.NewSession(attemptCtx, config.RedisURL())
		cancel()
		if databaseErr == nil && redisErr == nil {
			return database, sessions, nil
		}
		if database != nil {
			_ = database.Close()
		}
		if sessions != nil {
			_ = sessions.Close()
		}
		log.Printf("connect to storage: postgres=%v redis=%v", databaseErr, redisErr)
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
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
