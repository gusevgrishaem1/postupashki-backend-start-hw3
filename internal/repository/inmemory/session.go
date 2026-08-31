package inmemory

import (
	"sync"
	"time"

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
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[id]
	if ok && !session.ExpiresAt.IsZero() && !time.Now().Before(session.ExpiresAt) {
		delete(r.sessions, id)
		return domain.Session{}, false
	}
	return session, ok
}

func (r *Session) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sessions[id]; !ok {
		return false
	}
	delete(r.sessions, id)
	return true
}

func (r *Session) RemoveExpired(now time.Time) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	removed := 0
	for id, session := range r.sessions {
		if !session.ExpiresAt.IsZero() && !now.Before(session.ExpiresAt) {
			delete(r.sessions, id)
			removed++
		}
	}

	return removed
}
