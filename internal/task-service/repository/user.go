package repository

import (
	"context"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

type User interface {
	Save(context.Context, domain.User) error
	GetByLogin(context.Context, string) (domain.User, error)
}
