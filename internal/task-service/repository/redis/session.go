package redis

import (
	"context"
	"log"

	redisclient "github.com/redis/go-redis/v9"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

const sessionPrefix = "session:"

type Session struct {
	client *redisclient.Client
}

func NewSession(url string) (*Session, error) {
	options, err := redisclient.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redisclient.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &Session{client: client}, nil
}

func (r *Session) Close() error {
	return r.client.Close()
}

func (r *Session) Save(session domain.Session) error {
	return r.client.Set(context.Background(), sessionPrefix+session.SessionID, session.UserID, 0).Err()
}

func (r *Session) Get(id string) (domain.Session, bool) {
	userID, err := r.client.Get(context.Background(), sessionPrefix+id).Result()
	if err != nil {
		log.Printf("get session: %v", err)
		return domain.Session{}, false
	}
	return domain.Session{UserID: userID, SessionID: id}, true
}
