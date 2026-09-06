FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /code-processor ./cmd/code-processor

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /code-processor .
EXPOSE 9000
CMD ["./code-processor"]
