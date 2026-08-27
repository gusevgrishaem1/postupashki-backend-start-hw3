package service

import (
	"github.com/google/uuid"

	"postupashki-backend-start-hw3/internal/contracts"
	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type Task struct {
	repository repository.Task
	publisher  Publisher
}

type Publisher interface {
	Publish(contracts.Submission) error
}

func NewTask(repository repository.Task, publisher Publisher) *Task {
	service := &Task{repository: repository, publisher: publisher}
	return service
}

func (s *Task) Create(code, language string) (string, error) {
	id := uuid.NewString()
	s.repository.Save(domain.Task{ID: id, Status: domain.InProgress})
	if err := s.publisher.Publish(contracts.Submission{ID: id, Code: code, Language: language}); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Task) Commit(id string, result domain.Result) error {
	task, ok := s.repository.Get(id)
	if !ok {
		return usecases.ErrNotFound
	}
	task.Status = domain.Ready
	task.Result = result
	s.repository.Save(task)
	return nil
}

func (s *Task) Get(id string) (domain.Task, error) {
	task, ok := s.repository.Get(id)
	if !ok {
		return domain.Task{}, usecases.ErrNotFound
	}
	return task, nil
}
