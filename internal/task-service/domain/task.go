package domain

import "postupashki-backend-start-hw3/internal/contracts"

type Task struct {
	ID     string
	Status string
	Result contracts.Result
}

type Result = contracts.Result

const (
	InProgress = "in_progress"
	Ready      = "ready"
)
