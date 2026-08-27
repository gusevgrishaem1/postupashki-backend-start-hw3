package repository

import "postupashki-backend-start-hw3/internal/task-service/domain"

type Task interface {
	Save(domain.Task)
	Get(string) (domain.Task, bool)
}
