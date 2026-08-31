FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /task-service ./cmd/task-service

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /task-service .
EXPOSE 8000
CMD ["./task-service"]
