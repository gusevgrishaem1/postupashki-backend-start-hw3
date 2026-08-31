package repository

import (
	"context"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

type Session interface {
	Save(context.Context, domain.Session) error
	Get(context.Context, string) (domain.Session, error)
	Delete(context.Context, string) error
}
