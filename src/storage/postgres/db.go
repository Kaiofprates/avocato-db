package postgres

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5"
)

//go:embed schema.sql
var schemaFS embed.FS

func Connect(ctx context.Context, url string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Initialize schema from embedded file
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded schema: %w", err)
	}

	if _, err := conn.Exec(ctx, string(schema)); err != nil {
		return nil, fmt.Errorf("failed to execute schema: %w", err)
	}

	return conn, nil
}
