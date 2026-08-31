APP_NAME := task-service
PORT ?= 8000

.PHONY: run build test docker-build docker-run docker-stop

run:
	go run ./cmd/api

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api

test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME) .

docker-run: docker-build
	docker run --rm --name $(APP_NAME) -p $(PORT):8000 $(APP_NAME)

docker-stop:
	docker stop $(APP_NAME)
