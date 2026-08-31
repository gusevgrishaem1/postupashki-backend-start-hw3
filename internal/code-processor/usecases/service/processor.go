package service

import (
	"context"
	"log"

	"postupashki-backend-start-hw3/internal/code-processor/repository"
	"postupashki-backend-start-hw3/internal/code-processor/usecases"
	"postupashki-backend-start-hw3/internal/contracts"
)

type Processor struct {
	consumer  repository.Consumer
	runner    usecases.Runner
	committer repository.Committer
}

func NewProcessor(consumer repository.Consumer, runner usecases.Runner, committer repository.Committer) *Processor {
	return &Processor{consumer: consumer, runner: runner, committer: committer}
}

func (s *Processor) Run(ctx context.Context) error {
	deliveries, err := s.consumer.Consume()
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		submission, err := delivery.Submission()
		if err != nil {
			log.Printf("discard invalid submission: %v", err)
			_ = delivery.Reject(false)
			continue
		}

		runResult, err := s.runner.Run(ctx, usecases.RunRequest{
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
		if err := s.committer.Commit(submission.ID, result); err != nil {
			log.Printf("commit task %s: %v", submission.ID, err)
			_ = delivery.Nack(true)
			continue
		}
		if err := delivery.Ack(); err != nil {
			log.Printf("ack task %s: %v", submission.ID, err)
		}
	}
	return nil
}
