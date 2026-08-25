package repository

import "postupashki-backend-start-hw3/internal/domain"

type Session interface {
	Save(domain.Session)
	Get(string) (domain.Session, bool)
}
