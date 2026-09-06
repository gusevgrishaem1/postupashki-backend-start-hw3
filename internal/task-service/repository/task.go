package repository

import (
	"context"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

type Task interface {
	Save(context.Context, domain.Task) error
	Get(context.Context, string) (domain.Task, error)
	Delete(context.Context, string) error
}
