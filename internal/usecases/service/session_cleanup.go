package service

import (
	"context"
	"time"

	"postupashki-backend-start-hw3/internal/repository"
)

const sessionCleanupInterval = time.Minute

type sessionCleanup struct {
	sessions repository.Session
	interval time.Duration
}

func newSessionCleanup(sessions repository.Session, interval time.Duration) *sessionCleanup {
	if interval <= 0 {
		interval = sessionCleanupInterval
	}
	return &sessionCleanup{sessions: sessions, interval: interval}
}

func (w *sessionCleanup) run(ctx context.Context) {
	w.sessions.RemoveExpired(time.Now())

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case now := <-ticker.C:
			w.sessions.RemoveExpired(now)
		case <-ctx.Done():
			return
		}
	}
}
