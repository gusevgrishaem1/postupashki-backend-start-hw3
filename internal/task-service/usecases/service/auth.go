package service

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log"

	"github.com/google/uuid"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type Auth struct {
	users    repository.User
	sessions repository.Session
}

func NewAuth(users repository.User, sessions repository.Session) *Auth {
	return &Auth{users: users, sessions: sessions}
}

func (s *Auth) Register(login, password string) error {
	hashedPassword := passwordHash(password)
	user := domain.User{ID: uuid.NewString(), Login: login, Password: hashedPassword}
	if err := s.users.Save(user); errors.Is(err, repository.ErrAlreadyExists) {
		existingUser, getErr := s.users.GetByLogin(login)
		if getErr != nil {
			log.Printf("register: get existing user: %v", getErr)
			return usecases.ErrServiceUnavailable
		}
		if subtle.ConstantTimeCompare([]byte(existingUser.Password), []byte(hashedPassword)) == 1 {
			return nil
		}
		return usecases.ErrUserExists
	} else if err != nil {
		log.Printf("register: save user: %v", err)
		return usecases.ErrServiceUnavailable
	}
	return nil
}

func (s *Auth) Login(login, password string) (string, error) {
	user, err := s.users.GetByLogin(login)
	if errors.Is(err, repository.ErrNotFound) {
		return "", usecases.ErrInvalidCredentials
	}
	if err != nil {
		log.Printf("login: get user: %v", err)
		return "", usecases.ErrServiceUnavailable
	}
	if subtle.ConstantTimeCompare([]byte(user.Password), []byte(passwordHash(password))) != 1 {
		return "", usecases.ErrInvalidCredentials
	}

	token := uuid.NewString()
	if err := s.sessions.Save(domain.Session{UserID: user.ID, SessionID: token}); err != nil {
		log.Printf("login: save session: %v", err)
		return "", usecases.ErrServiceUnavailable
	}
	return token, nil
}

func (s *Auth) Authenticate(token string) error {
	_, err := s.sessions.Get(token)
	if errors.Is(err, repository.ErrNotFound) {
		return usecases.ErrInvalidSession
	}
	if err != nil {
		log.Printf("authenticate: get session: %v", err)
		return usecases.ErrServiceUnavailable
	}
	return nil
}

func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
