package main

import (
	"log"
	"net/http"

	httpapi "postupashki-backend-start-hw3/internal/api/http"
	"postupashki-backend-start-hw3/internal/config"
	"postupashki-backend-start-hw3/internal/repository/inmemory"
	"postupashki-backend-start-hw3/internal/usecases/service"
)

// @title Task service
// @version 1.0
// @description Asynchronous task processing service.
// @host localhost:8000
// @BasePath /
func main() {
	taskRepository := inmemory.NewTask()
	taskService := service.NewTask(taskRepository)
	server := httpapi.NewServer(taskService)

	log.Fatal(http.ListenAndServe(config.Address(), server.Handler()))
}
