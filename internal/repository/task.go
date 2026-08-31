package repository

import "postupashki-backend-start-hw3/internal/domain"

type Task interface {
	Save(domain.Task)
	Get(string) (domain.Task, bool)
}
