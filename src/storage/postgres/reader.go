package postgres

import (
	"avocato-db/src/core/ledger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetBlockByIndex(ctx context.Context, db *pgx.Conn, index uint64) (*ledger.Block, error) {
	var hashStr, prevHashStr string
	var b ledger.Block
	err := db.QueryRow(ctx, "SELECT index, hash, prev_hash, created_at, payload FROM blocks WHERE index = $1", index).Scan(
		&b.Index, &hashStr, &prevHashStr, &b.Timestamp, &b.Payload,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get block by index: %w", err)
	}

	fmt.Sscanf(hashStr, "%x", &b.Hash)
	fmt.Sscanf(prevHashStr, "%x", &b.PrevHash)

	return &b, nil
}

func GetBlockByHash(ctx context.Context, db *pgx.Conn, hash string) (*ledger.Block, error) {
	var hashStr, prevHashStr string
	var b ledger.Block
	err := db.QueryRow(ctx, "SELECT index, hash, prev_hash, created_at, payload FROM blocks WHERE hash = $1", hash).Scan(
		&b.Index, &hashStr, &prevHashStr, &b.Timestamp, &b.Payload,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get block by hash: %w", err)
	}

	fmt.Sscanf(hashStr, "%x", &b.Hash)
	fmt.Sscanf(prevHashStr, "%x", &b.PrevHash)

	return &b, nil
}
