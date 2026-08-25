package inmemory

import (
	"sync"

	"postupashki-backend-start-hw3/internal/domain"
)

type User struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewUser() *User {
	return &User{users: make(map[string]domain.User)}
}

func (r *User) Save(user domain.User) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Login]; exists {
		return false
	}
	r.users[user.Login] = user
	return true
}

func (r *User) GetByLogin(login string) (domain.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[login]
	return user, ok
}
