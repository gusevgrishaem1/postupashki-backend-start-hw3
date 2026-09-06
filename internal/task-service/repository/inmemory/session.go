package inmemory

import (
	"context"
	"sync"
	"time"

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

func (r *Session) Save(_ context.Context, session domain.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.SessionID] = session
	return nil
}

func (r *Session) Get(_ context.Context, id string) (domain.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[id]
	if !ok || (!session.ExpiresAt.IsZero() && !time.Now().Before(session.ExpiresAt)) {
		delete(r.sessions, id)
		return domain.Session{}, repository.ErrNotFound
	}
	return session, nil
}

func (r *Session) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.sessions[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.sessions, id)
	return nil
}
