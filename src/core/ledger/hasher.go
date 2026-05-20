package ledger

import (
	"avocato-db/src/core/crypto"
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// CalculateBlockHash computes the SHA-256 hash of a block using JCS canonicalization.
// Hash = SHA256(JCS(index + uuid + timestamp + prev_hash + JCS(payload)))
func CalculateBlockHash(index uint64, uid string, timestamp int64, prevHash [32]byte, payload interface{}) ([32]byte, []byte, error) {
	// 1. Canonicalize payload
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	payloadCanon, err := crypto.Canonicalize(payloadRaw)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("failed to canonicalize payload: %w", err)
	}

	// 2. Prepare block structure for hashing
	blockMap := map[string]interface{}{
		"index":     index,
		"uuid":      uid,
		"timestamp": timestamp,
		"prev_hash": fmt.Sprintf("%x", prevHash),
		"payload":   json.RawMessage(payloadCanon),
	}

	blockRaw, err := json.Marshal(blockMap)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("failed to marshal block map: %w", err)
	}
	blockCanon, err := crypto.Canonicalize(blockRaw)
	if err != nil {
		return [32]byte{}, nil, fmt.Errorf("failed to canonicalize block: %w", err)
	}

	// 3. Hash the canonicalized block
	hash := sha256.Sum256(blockCanon)
	return hash, payloadCanon, nil
}
