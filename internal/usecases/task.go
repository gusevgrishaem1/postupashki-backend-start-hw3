package usecases

import (
	"errors"

	"postupashki-backend-start-hw3/internal/domain"
)

var ErrNotFound = errors.New("task not found")

type Task interface {
	Create() string
	Get(string) (domain.Task, error)
}
