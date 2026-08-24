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
	repository := inmemory.NewTask()
	usecase := service.NewTask(repository)
	handler := httpapi.NewTask(usecase)

	log.Fatal(http.ListenAndServe(config.Address(), handler.Handler()))
}
