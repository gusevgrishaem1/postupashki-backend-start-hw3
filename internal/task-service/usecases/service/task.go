package service

import (
	"errors"
	"log"

	"github.com/google/uuid"

	"postupashki-backend-start-hw3/internal/contracts"
	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type Task struct {
	repository repository.Task
	publisher  repository.Publisher
}

func NewTask(repository repository.Task, publisher repository.Publisher) *Task {
	service := &Task{repository: repository, publisher: publisher}
	return service
}

func (s *Task) Create(code, language, input string) (string, error) {
	id := uuid.NewString()
	if err := s.repository.Save(domain.Task{ID: id, Status: domain.InProgress}); err != nil {
		log.Printf("create task: save task: %v", err)
		return "", usecases.ErrServiceUnavailable
	}
	if err := s.publisher.Publish(contracts.Submission{ID: id, Code: code, Language: language, Input: input}); err != nil {
		log.Printf("create task: publish submission: %v", err)
		return "", usecases.ErrServiceUnavailable
	}
	return id, nil
}

func (s *Task) Commit(id string, result domain.Result) error {
	task, err := s.repository.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		return usecases.ErrNotFound
	}
	if err != nil {
		log.Printf("commit task: get task: %v", err)
		return usecases.ErrServiceUnavailable
	}
	task.Status = domain.Ready
	task.Result = result
	if err := s.repository.Save(task); err != nil {
		log.Printf("commit task: save task: %v", err)
		return usecases.ErrServiceUnavailable
	}
	return nil
}

func (s *Task) Get(id string) (domain.Task, error) {
	task, err := s.repository.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		return domain.Task{}, usecases.ErrNotFound
	}
	if err != nil {
		log.Printf("get task: %v", err)
		return domain.Task{}, usecases.ErrServiceUnavailable
	}
	return task, nil
}
