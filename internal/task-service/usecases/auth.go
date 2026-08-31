package usecases

import (
	"context"
	"errors"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid session")
)

type Auth interface {
	Register(context.Context, string, string) error
	Login(context.Context, string, string) (string, error)
	Authenticate(context.Context, string) error
	Logout(context.Context, string) error
}
