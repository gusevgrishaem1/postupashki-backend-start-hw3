package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/repository"
)

type User struct {
	database *sql.DB
}

func NewUser(database *sql.DB) *User {
	return &User{database: database}
}

func (r *User) Save(ctx context.Context, user domain.User) error {
	_, err := r.database.ExecContext(ctx, `INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`, user.ID, user.Login, user.Password)
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return repository.ErrAlreadyExists
	}
	return err
}

func (r *User) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	var user domain.User
	err := r.database.QueryRowContext(ctx, `SELECT id, login, password FROM users WHERE login = $1`, login).Scan(&user.ID, &user.Login, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
