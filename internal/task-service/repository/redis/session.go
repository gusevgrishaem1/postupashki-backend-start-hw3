package redis

import (
	"context"
	"errors"
	"time"

	redisclient "github.com/redis/go-redis/v9"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
)

const sessionPrefix = "session:"

type Session struct {
	client *redisclient.Client
}

func NewSession(ctx context.Context, url string) (*Session, error) {
	options, err := redisclient.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redisclient.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &Session{client: client}, nil
}

func (r *Session) Close() error {
	return r.client.Close()
}

func (r *Session) Save(ctx context.Context, session domain.Session) error {
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return repository.ErrNotFound
	}
	return r.client.Set(ctx, sessionPrefix+session.SessionID, session.UserID, ttl).Err()
}

func (r *Session) Get(ctx context.Context, id string) (domain.Session, error) {
	userID, err := r.client.Get(ctx, sessionPrefix+id).Result()
	if errors.Is(err, redisclient.Nil) {
		return domain.Session{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Session{}, err
	}
	return domain.Session{UserID: userID, SessionID: id}, nil
}

func (r *Session) Delete(ctx context.Context, id string) error {
	deleted, err := r.client.Del(ctx, sessionPrefix+id).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return repository.ErrNotFound
	}
	return nil
}
