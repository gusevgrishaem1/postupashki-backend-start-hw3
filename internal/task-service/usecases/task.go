package usecases

import (
	"errors"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

var ErrNotFound = errors.New("task not found")

type Task interface {
	Create(code, language, input string) (string, error)
	Get(string) (domain.Task, error)
	Commit(string, domain.Result) error
}
