package repository

import (
	"context"

	"postupashki-backend-start-hw3/internal/contracts"
)

type Publisher interface {
	Publish(context.Context, contracts.Submission) error
}
