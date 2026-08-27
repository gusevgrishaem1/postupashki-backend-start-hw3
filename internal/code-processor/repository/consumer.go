package repository

import "postupashki-backend-start-hw3/internal/contracts"

type Delivery interface {
	Submission() (contracts.Submission, error)
	Ack() error
	Reject(requeue bool) error
	Nack(requeue bool) error
}

type Consumer interface {
	Consume() (<-chan Delivery, error)
}
