package usecases

import (
	"context"
	"errors"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

var ErrNotFound = errors.New("task not found")

type Task interface {
	Create(context.Context, string, string, string) (string, error)
	Get(context.Context, string) (domain.Task, error)
	Commit(context.Context, string, domain.Result) error
}
