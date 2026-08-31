package postgres

import (
	"context"
	"database/sql"
	"errors"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
)

type Task struct {
	database *sql.DB
}

func NewTask(database *sql.DB) *Task {
	return &Task{database: database}
}

func (r *Task) Save(ctx context.Context, task domain.Task) error {
	_, err := r.database.ExecContext(ctx, `
		INSERT INTO tasks (id, status, stdout, stderr, exit_code)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			stdout = EXCLUDED.stdout,
			stderr = EXCLUDED.stderr,
			exit_code = EXCLUDED.exit_code`,
		task.ID, task.Status, task.Result.Stdout, task.Result.Stderr, task.Result.ExitCode)
	return err
}

func (r *Task) Get(ctx context.Context, id string) (domain.Task, error) {
	var task domain.Task
	err := r.database.QueryRowContext(ctx, `SELECT id, status, stdout, stderr, exit_code FROM tasks WHERE id = $1`, id).Scan(
		&task.ID, &task.Status, &task.Result.Stdout, &task.Result.Stderr, &task.Result.ExitCode,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (r *Task) Delete(ctx context.Context, id string) error {
	_, err := r.database.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}
