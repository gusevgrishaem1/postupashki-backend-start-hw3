package repository

import "postupashki-backend-start-hw3/internal/contracts"

type Committer interface {
	Commit(taskID string, result contracts.Result) error
}
