package repository

import "postupashki-backend-start-hw3/internal/task-service/domain"

type Session interface {
	Save(domain.Session) error
	Get(string) (domain.Session, error)
}
