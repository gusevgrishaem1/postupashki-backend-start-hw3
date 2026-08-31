package domain

type Task struct {
	ID     string
	Status string
	Result string
}

const (
	InProgress = "in_progress"
	Ready      = "ready"
)
