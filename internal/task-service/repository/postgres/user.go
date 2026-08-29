package postgres

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"

	"postupashki-backend-start-hw3/internal/task-service/domain"
)

type User struct {
	database *sql.DB
}

func NewUser(database *sql.DB) *User {
	return &User{database: database}
}

func (r *User) Save(user domain.User) bool {
	_, err := r.database.Exec(`INSERT INTO users (id, login, password) VALUES ($1, $2, $3)`, user.ID, user.Login, user.Password)
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return false
	}
	if err != nil {
		log.Printf("save user: %v", err)
		return false
	}
	return true
}

func (r *User) GetByLogin(login string) (domain.User, bool) {
	var user domain.User
	err := r.database.QueryRow(`SELECT id, login, password FROM users WHERE login = $1`, login).Scan(&user.ID, &user.Login, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, false
	}
	if err != nil {
		log.Printf("get user: %v", err)
		return domain.User{}, false
	}
	return user, true
}
