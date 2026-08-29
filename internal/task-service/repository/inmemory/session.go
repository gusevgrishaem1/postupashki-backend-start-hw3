package inmemory

import (
	"sync"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
)

type Session struct {
	mu       sync.RWMutex
	sessions map[string]domain.Session
}

func NewSession() *Session {
	return &Session{sessions: make(map[string]domain.Session)}
}

func (r *Session) Save(session domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *Session) Get(id string) (domain.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[id]
	if !ok {
		return domain.Session{}, repository.ErrNotFound
	}
	return session, nil
}
