package postgres

import (
	"database/sql"
	"errors"
	"log"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

type Task struct {
	database *sql.DB
}

func NewTask(database *sql.DB) *Task {
	return &Task{database: database}
}

func (r *Task) Save(task domain.Task) {
	_, err := r.database.Exec(`
		INSERT INTO tasks (id, status, stdout, stderr, exit_code)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			stdout = EXCLUDED.stdout,
			stderr = EXCLUDED.stderr,
			exit_code = EXCLUDED.exit_code`,
		task.ID, task.Status, task.Result.Stdout, task.Result.Stderr, task.Result.ExitCode)
	if err != nil {
		log.Printf("save task: %v", err)
	}
}

func (r *Task) Get(id string) (domain.Task, bool) {
	var task domain.Task
	err := r.database.QueryRow(`SELECT id, status, stdout, stderr, exit_code FROM tasks WHERE id = $1`, id).Scan(
		&task.ID, &task.Status, &task.Result.Stdout, &task.Result.Stderr, &task.Result.ExitCode,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, false
	}
	if err != nil {
		log.Printf("get task: %v", err)
		return domain.Task{}, false
	}
	return task, true
}
