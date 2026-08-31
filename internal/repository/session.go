package repository

import (
	"time"

	"postupashki-backend-start-hw3/internal/domain"
)

type Session interface {
	Save(domain.Session)
	Get(string) (domain.Session, bool)
	Delete(string) bool
	RemoveExpired(time.Time) int
}
