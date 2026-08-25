package inmemory

import (
	"sync"

	"postupashki-backend-start-hw3/internal/domain"
)

type Session struct {
	mu       sync.RWMutex
	sessions map[string]domain.Session
}

func NewSession() *Session {
	return &Session{sessions: make(map[string]domain.Session)}
}

func (r *Session) Save(session domain.Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.SessionID] = session
}

func (r *Session) Get(id string) (domain.Session, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[id]
	return session, ok
}
