package usecases

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSession     = errors.New("invalid session")
)

type Auth interface {
	Register(login, password string) error
	Login(login, password string) (string, error)
	Authenticate(token string) error
	Logout(token string) error
}
