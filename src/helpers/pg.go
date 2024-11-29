package helpers

import (
	"context"
	"database/sql"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func InitTestContainersPostgres(ctx context.Context) (*sql.DB, *postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:16.6-bullseye",
		postgres.WithInitScripts(filepath.Join("..", "scripts", "init-db.sql")),
		postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
			wait.ForExec([]string{"sleep", "2"}),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	return db, pgContainer, nil
}
