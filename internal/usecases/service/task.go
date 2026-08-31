package service

import (
	"time"

	"github.com/google/uuid"

	"postupashki-backend-start-hw3/internal/domain"
	"postupashki-backend-start-hw3/internal/repository"
	"postupashki-backend-start-hw3/internal/usecases"
)

type Task struct {
	repository repository.Task
}

func NewTask(repository repository.Task) *Task {
	return &Task{repository: repository}
}

func (s *Task) Create() string {
	id := uuid.NewString()
	s.repository.Save(domain.Task{ID: id, Status: domain.InProgress})

	go func() {
		time.Sleep(time.Second * 15)
		s.repository.Save(domain.Task{ID: id, Status: domain.Ready, Result: "done"})
	}()

	return id
}

func (s *Task) Get(id string) (domain.Task, error) {
	task, ok := s.repository.Get(id)
	if !ok {
		return domain.Task{}, usecases.ErrNotFound
	}
	return task, nil
}
