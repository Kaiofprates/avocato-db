package ledger

import (
	"avocato-db/src/core"
	"avocato-db/src/storage/wal"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	mu          sync.RWMutex
	lastIndex   uint64
	lastHash    [32]byte
	groupCommit *GroupCommit
	db          *pgx.Conn
	cache       *Cache
}

func NewService(walPath string, db *pgx.Conn) (*Service, error) {
	writer, err := wal.NewWriter(walPath)
	if err != nil {
		return nil, err
	}

	gc := NewGroupCommit(writer, 5*time.Millisecond)

	s := &Service{
		groupCommit: gc,
		db:          db,
		cache:       NewCache(1000), // Cache last 1000 blocks
	}

	if err := s.recoverState(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) recoverState() error {
	var index uint64
	var hashStr string
	err := s.db.QueryRow(context.Background(), "SELECT index, hash FROM blocks ORDER BY index DESC LIMIT 1").Scan(&index, &hashStr)
	if err != nil {
		if err == pgx.ErrNoRows {
			s.lastIndex = 0
			s.lastHash = [32]byte{}
			return nil
		}
		return err
	}

	s.lastIndex = index
	fmt.Sscanf(hashStr, "%x", &s.lastHash)
	return nil
}

func (s *Service) Append(ctx context.Context, payload interface{}) (*Block, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := s.lastIndex + 1
	timestamp := time.Now().UnixNano()
	prevHash := s.lastHash

	hash, canonPayload, err := CalculateBlockHash(index, timestamp, prevHash, payload)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Index:     index,
		Timestamp: timestamp,
		PrevHash:  prevHash,
		Payload:   canonPayload,
		Hash:      hash,
	}

	// 1. Submit to Group Commit (WAL)
	done := s.groupCommit.SubmitBlock(block)
	
	if err := <-done; err != nil {
		return nil, fmt.Errorf("WAL write failed: %w", err)
	}

	// 2. Persist to Postgres
	_, err = s.db.Exec(ctx, "INSERT INTO blocks (index, hash, prev_hash, created_at, payload) VALUES ($1, $2, $3, $4, $5)",
		block.Index, fmt.Sprintf("%x", block.Hash), fmt.Sprintf("%x", block.PrevHash), block.Timestamp, block.Payload)
	if err != nil {
		return nil, fmt.Errorf("postgres write failed: %w", err)
	}

	s.lastIndex = index
	s.lastHash = hash

	// Add to cache
	s.cache.Put(fmt.Sprintf("%d", block.Index), block)
	s.cache.Put(fmt.Sprintf("%x", block.Hash), block)

	core.LogInfo("Appended block %d with hash %x", index, hash)
	return block, nil
}

func (s *Service) GetBlock(ctx context.Context, id string) (*Block, error) {
	// 1. Try Cache
	if b, ok := s.cache.Get(id); ok {
		return b, nil
	}

	// 2. Try Postgres
	b, err := s.readFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	if b != nil {
		s.cache.Put(fmt.Sprintf("%d", b.Index), b)
		s.cache.Put(fmt.Sprintf("%x", b.Hash), b)
	}

	return b, nil
}

func (s *Service) readFromDB(ctx context.Context, id string) (*Block, error) {
	var hashStr, prevHashStr string
	var b Block
	var query string
	
	// Determine if id is index or hash
	var index uint64
	if _, err := fmt.Sscanf(id, "%d", &index); err == nil {
		query = "SELECT index, hash, prev_hash, created_at, payload FROM blocks WHERE index = $1"
	} else {
		query = "SELECT index, hash, prev_hash, created_at, payload FROM blocks WHERE hash = $1"
	}
	
	err := s.db.QueryRow(ctx, query, id).Scan(
		&b.Index, &hashStr, &prevHashStr, &b.Timestamp, &b.Payload,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	fmt.Sscanf(hashStr, "%x", &b.Hash)
	fmt.Sscanf(prevHashStr, "%x", &b.PrevHash)

	return &b, nil
}

func (s *Service) GetLastState() (uint64, [32]byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastIndex, s.lastHash
}
