package repository

import "postupashki-backend-start-hw3/internal/contracts"

type Publisher interface {
	Publish(contracts.Submission) error
}
