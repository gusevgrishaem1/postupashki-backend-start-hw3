package repository

import "postupashki-backend-start-hw3/internal/task-service/domain"

type User interface {
	Save(domain.User) error
	GetByLogin(string) (domain.User, error)
}
