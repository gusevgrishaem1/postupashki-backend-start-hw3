package service

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"

	"github.com/google/uuid"

	"postupashki-backend-start-hw3/internal/domain"
	"postupashki-backend-start-hw3/internal/repository"
	"postupashki-backend-start-hw3/internal/usecases"
)

type Auth struct {
	users    repository.User
	sessions repository.Session
}

func NewAuth(users repository.User, sessions repository.Session) *Auth {
	return &Auth{users: users, sessions: sessions}
}

func (s *Auth) Register(login, password string) error {
	user := domain.User{ID: uuid.NewString(), Login: login, Password: passwordHash(password)}
	if !s.users.Save(user) {
		return usecases.ErrUserExists
	}
	return nil
}

func (s *Auth) Login(login, password string) (string, error) {
	user, ok := s.users.GetByLogin(login)
	if !ok || subtle.ConstantTimeCompare([]byte(user.Password), []byte(passwordHash(password))) != 1 {
		return "", usecases.ErrInvalidCredentials
	}

	token := uuid.NewString()
	s.sessions.Save(domain.Session{UserID: user.ID, SessionID: token})
	return token, nil
}

func (s *Auth) Authenticate(token string) error {
	if _, ok := s.sessions.Get(token); !ok {
		return usecases.ErrInvalidSession
	}
	return nil
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
