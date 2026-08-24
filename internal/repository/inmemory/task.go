package inmemory

import (
	"sync"

	"postupashki-backend-start-hw3/internal/domain"
)

type Task struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

func NewTask() *Task {
	return &Task{tasks: make(map[string]domain.Task)}
}

func (r *Task) Save(task domain.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
}

func (r *Task) Get(id string) (domain.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	return task, ok
}
