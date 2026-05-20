package integrity

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Checkpoint struct {
	LastIndex uint64
	MMRRoot   string
	LastHash  string
}

func LoadLatestCheckpoint(ctx context.Context, db *pgx.Conn) (*Checkpoint, error) {
	var cp Checkpoint
	err := db.QueryRow(ctx, "SELECT last_index, mmr_root, last_hash FROM checkpoints ORDER BY last_index DESC LIMIT 1").Scan(
		&cp.LastIndex, &cp.MMRRoot, &cp.LastHash,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No checkpoint yet
		}
		return nil, fmt.Errorf("failed to load checkpoint: %w", err)
	}
	return &cp, nil
}

func SaveCheckpoint(ctx context.Context, db *pgx.Conn, cp *Checkpoint) error {
	_, err := db.Exec(ctx, "INSERT INTO checkpoints (last_index, mmr_root, last_hash) VALUES ($1, $2, $3) ON CONFLICT (last_index) DO NOTHING",
		cp.LastIndex, cp.MMRRoot, cp.LastHash,
	)
	if err != nil {
		return fmt.Errorf("failed to save checkpoint: %w", err)
	}
	return nil
}
