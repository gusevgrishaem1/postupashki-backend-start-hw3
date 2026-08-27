APP_NAME := task-service
PORT ?= 8000

.PHONY: run code-processor build test integration-test up down

run:
	go run ./cmd/task-service

code-processor:
	go run ./cmd/code-processor

build:
	mkdir -p bin
	go build -o bin/task-service ./cmd/task-service
	go build -o bin/code-processor ./cmd/code-processor

test:
	go test ./...

integration-test: up
	pytest -q tests/tests.py

up:
	docker compose up --build -d --wait --wait-timeout 60

down:
	docker compose down --remove-orphans
