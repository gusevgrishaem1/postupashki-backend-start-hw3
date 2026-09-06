package repository

import "errors"

var (
	ErrNotFound      = errors.New("repository entity not found")
	ErrAlreadyExists = errors.New("repository entity already exists")
)
