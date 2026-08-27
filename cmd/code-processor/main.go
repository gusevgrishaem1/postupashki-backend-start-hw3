package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"postupashki-backend-start-hw3/internal/code-processor/broker"
	"postupashki-backend-start-hw3/internal/code-processor/committer"
	"postupashki-backend-start-hw3/internal/code-processor/config"
	"postupashki-backend-start-hw3/internal/code-processor/runner"
	"postupashki-backend-start-hw3/internal/contracts"
)

func main() {
	queue := connect()
	defer queue.Close()

	deliveries, err := queue.Consume()
	if err != nil {
		log.Fatal(err)
	}
	codeRunner, err := runner.New()
	if err != nil {
		log.Fatal(err)
	}
	defer codeRunner.Close()

	resultCommitter := committer.NewCommitter(config.TaskServiceURL())
	for delivery := range deliveries {
		var submission contracts.Submission
		if err := json.Unmarshal(delivery.Body, &submission); err != nil {
			log.Printf("discard invalid submission: %v", err)
			_ = delivery.Reject(false)
			continue
		}

		runResult, err := codeRunner.Run(context.Background(), runner.RunRequest{
			Runtime: submission.Language,
			Code:    submission.Code,
			Input:   submission.Input,
		})
		result := contracts.Result{
			Stdout: runResult.Stdout, Stderr: runResult.Stderr, ExitCode: runResult.ExitCode,
		}
		if err != nil {
			result.Stderr = err.Error()
			result.ExitCode = -1
		}
		if err := resultCommitter.Commit(submission.ID, result); err != nil {
			log.Printf("commit task %s: %v", submission.ID, err)
			_ = delivery.Nack(false, true)
			continue
		}
		_ = delivery.Ack(false)
	}
}

func connect() *broker.RabbitMQ {
	for {
		queue, err := broker.NewRabbitMQ(config.RabbitMQURL())
		if err == nil {
			return queue
		}
		log.Printf("connect to RabbitMQ: %v", err)
		time.Sleep(2 * time.Second)
	}
}
