package repository

import "postupashki-backend-start-hw3/internal/domain"

type User interface {
	Save(domain.User) bool
	GetByLogin(string) (domain.User, bool)
}
