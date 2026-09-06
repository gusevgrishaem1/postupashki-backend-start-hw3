package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type Auth struct {
	users      repository.User
	sessions   repository.Session
	sessionTTL time.Duration
}

const defaultSessionTTL = 24 * time.Hour

func NewAuth(users repository.User, sessions repository.Session, sessionTTL time.Duration) *Auth {
	if sessionTTL <= 0 {
		sessionTTL = defaultSessionTTL
	}
	return &Auth{users: users, sessions: sessions, sessionTTL: sessionTTL}
}

func (s *Auth) Register(ctx context.Context, login, password string) error {
	hashedPassword, err := passwordHash(password)
	if err != nil {
		log.Printf("register: hash password: %v", err)
		return usecases.ErrServiceUnavailable
	}
	user := domain.User{ID: uuid.NewString(), Login: login, Password: hashedPassword}
	if err := s.users.Save(ctx, user); errors.Is(err, repository.ErrAlreadyExists) {
		existingUser, getErr := s.users.GetByLogin(ctx, login)
		if getErr != nil {
			log.Printf("register: get existing user: %v", getErr)
			return usecases.ErrServiceUnavailable
		}
		if bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)) == nil {
			return nil
		}
		return usecases.ErrUserExists
	} else if err != nil {
		log.Printf("register: save user: %v", err)
		return usecases.ErrServiceUnavailable
	}
	return nil
}

func (s *Auth) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if errors.Is(err, repository.ErrNotFound) {
		return "", usecases.ErrInvalidCredentials
	}
	if err != nil {
		log.Printf("login: get user: %v", err)
		return "", usecases.ErrServiceUnavailable
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", usecases.ErrInvalidCredentials
	}

	token := uuid.NewString()
	if err := s.sessions.Save(ctx, domain.Session{UserID: user.ID, SessionID: token, ExpiresAt: time.Now().Add(s.sessionTTL)}); err != nil {
		log.Printf("login: save session: %v", err)
		return "", usecases.ErrServiceUnavailable
	}
	return token, nil
}

func (s *Auth) Authenticate(ctx context.Context, token string) error {
	_, err := s.sessions.Get(ctx, token)
	if errors.Is(err, repository.ErrNotFound) {
		return usecases.ErrInvalidSession
	}
	if err != nil {
		log.Printf("authenticate: get session: %v", err)
		return usecases.ErrServiceUnavailable
	}
	return nil
}

func (s *Auth) Logout(ctx context.Context, token string) error {
	if err := s.sessions.Delete(ctx, token); errors.Is(err, repository.ErrNotFound) {
		return usecases.ErrInvalidSession
	} else if err != nil {
		log.Printf("logout: delete session: %v", err)
		return usecases.ErrServiceUnavailable
	}
	return nil
}

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
