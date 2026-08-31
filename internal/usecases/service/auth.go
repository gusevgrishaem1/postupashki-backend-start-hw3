package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"postupashki-backend-start-hw3/internal/domain"
	"postupashki-backend-start-hw3/internal/repository"
	"postupashki-backend-start-hw3/internal/usecases"
)

type Auth struct {
	users      repository.User
	sessions   repository.Session
	sessionTTL time.Duration
	cancel     context.CancelFunc
}

const defaultSessionTTL = 24 * time.Hour

func NewAuth(users repository.User, sessions repository.Session, sessionTTL time.Duration) *Auth {
	if sessionTTL <= 0 {
		sessionTTL = defaultSessionTTL
	}
	ctx, cancel := context.WithCancel(context.Background())
	auth := &Auth{users: users, sessions: sessions, sessionTTL: sessionTTL, cancel: cancel}
	go newSessionCleanup(sessions, sessionCleanupInterval).run(ctx)

	return auth
}

func (s *Auth) Close() {
	s.cancel()
}

func (s *Auth) Register(login, password string) error {
	hash, err := passwordHash(password)
	if err != nil {
		return err
	}
	user := domain.User{ID: uuid.NewString(), Login: login, Password: hash}
	if !s.users.Save(user) {
		return usecases.ErrUserExists
	}
	return nil
}

func (s *Auth) Login(login, password string) (string, error) {
	user, ok := s.users.GetByLogin(login)
	if !ok || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", usecases.ErrInvalidCredentials
	}

	token := uuid.NewString()
	s.sessions.Save(domain.Session{
		UserID:    user.ID,
		SessionID: token,
		ExpiresAt: time.Now().Add(s.sessionTTL),
	})
	return token, nil
}

func (s *Auth) Authenticate(token string) error {
	if _, ok := s.sessions.Get(token); !ok {
		return usecases.ErrInvalidSession
	}
	return nil
}

func (s *Auth) Logout(token string) error {
	if !s.sessions.Delete(token) {
		return usecases.ErrInvalidSession
	}
	return nil
}

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}
