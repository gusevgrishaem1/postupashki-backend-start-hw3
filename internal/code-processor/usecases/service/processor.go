package service

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"postupashki-backend-start-hw3/internal/code-processor/repository"
	"postupashki-backend-start-hw3/internal/code-processor/usecases"
	"postupashki-backend-start-hw3/internal/contracts"
)

type Processor struct {
	consumer  repository.Consumer
	runner    usecases.Runner
	committer repository.Committer
	duration  *prometheus.HistogramVec
	usage     *prometheus.CounterVec
}

func NewProcessor(consumer repository.Consumer, runner usecases.Runner, committer repository.Committer, registerer prometheus.Registerer) *Processor {
	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "code_processor_execution_duration_seconds",
		Help: "Time spent executing submitted code.",
	}, []string{"translator"})
	usage := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "code_processor_translator_usage_total",
		Help: "Number of submissions handled by each translator.",
	}, []string{"translator"})
	registerer.MustRegister(duration, usage)
	return &Processor{consumer: consumer, runner: runner, committer: committer, duration: duration, usage: usage}
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

		startedAt := time.Now()
		s.usage.WithLabelValues(submission.Language).Inc()
		runResult, err := s.runner.Run(ctx, usecases.RunRequest{
			Runtime: submission.Language,
			Code:    submission.Code,
			Input:   submission.Input,
		})
		s.duration.WithLabelValues(submission.Language).Observe(time.Since(startedAt).Seconds())
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
