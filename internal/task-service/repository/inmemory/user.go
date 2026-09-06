package inmemory

import (
	"context"
	"sync"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
)

type User struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewUser() *User {
	return &User{users: make(map[string]domain.User)}
}

func (r *User) Save(_ context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Login]; exists {
		return repository.ErrAlreadyExists
	}
	r.users[user.Login] = user
	return nil
}

func (r *User) GetByLogin(_ context.Context, login string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[login]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return user, nil
}
