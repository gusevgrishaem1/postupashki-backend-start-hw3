package postgres

import (
	"context"
	"database/sql"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Open(ctx context.Context, url string) (*sql.DB, error) {
	database, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, err
	}
	source := &migrate.EmbedFileSystemMigrationSource{FileSystem: migrations, Root: "migrations"}
	if _, err := migrate.ExecContext(ctx, database, "postgres", source, migrate.Up); err != nil {
		_ = database.Close()
		return nil, err
	}
	return database, nil
}
