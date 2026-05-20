package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Initialize schema
	schema, err := os.ReadFile("src/storage/postgres/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	if _, err := conn.Exec(ctx, string(schema)); err != nil {
		return nil, fmt.Errorf("failed to execute schema: %w", err)
	}

	return conn, nil
}
