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
	addr, err := config.Address()
	if err != nil {
		log.Fatal(err)
	}

	taskRepository := inmemory.NewTask()
	taskService := service.NewTask(taskRepository)
	server := httpapi.NewServer(taskService)

	log.Fatal(http.ListenAndServe(addr, server.Handler()))
}
