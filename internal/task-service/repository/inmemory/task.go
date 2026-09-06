package inmemory

import (
	"context"
	"sync"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
)

type Task struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

func NewTask() *Task {
	return &Task{tasks: make(map[string]domain.Task)}
}

func (r *Task) Save(_ context.Context, task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *Task) Get(_ context.Context, id string) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return domain.Task{}, repository.ErrNotFound
	}
	return task, nil
}

func (r *Task) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tasks, id)
	return nil
}
