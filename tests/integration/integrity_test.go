package integration

import (
	"avocato-db/src/core/integrity"
	"avocato-db/src/core/ledger"
	"avocato-db/src/storage/wal"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegrityCheck_CorruptionDetection(t *testing.T) {
	walPath := "test_ledger.wal"
	defer os.Remove(walPath)

	// 1. Create a WAL with valid records
	writer, err := wal.NewWriter(walPath)
	assert.NoError(t, err)

	gc := ledger.NewGroupCommit(writer, 1*time.Millisecond)

	// Add block 1
	payload1 := []byte(`{"msg":"first"}`)
	hash1, canon1, _ := ledger.CalculateBlockHash(1, 1000, [32]byte{}, payload1)
	b1 := &ledger.Block{Index: 1, Timestamp: 1000, PrevHash: [32]byte{}, Payload: canon1, Hash: hash1}
	<-gc.SubmitBlock(b1)

	// Add block 2
	payload2 := []byte(`{"msg":"second"}`)
	hash2, canon2, _ := ledger.CalculateBlockHash(2, 2000, hash1, payload2)
	b2 := &ledger.Block{Index: 2, Timestamp: 2000, PrevHash: hash1, Payload: canon2, Hash: hash2}
	<-gc.SubmitBlock(b2)

	gc.Stop()
	writer.Close()

	// 2. Run Boot Check (Should Pass)
	// We mock db as nil since our BootCheck can handle missing checkpoints if we mock LoadLatestCheckpoint
	// Actually, boot_check.go uses pgx.Conn, so we can't easily run it without a DB unless we mock.
	// But since this is a POC test, we'll just test the core logic of corruption.
	// Let's manually corrupt the file.
	f, err := os.OpenFile(walPath, os.O_RDWR, 0644)
	assert.NoError(t, err)

	// Corrupt a byte in the payload of block 1
	// The binary format is: [uint32:length][uint64:index][int64:ts][32b:prev_hash][uint32:payload_len][varb:payload][32b:hash][uint32:checksum]
	// Payload starts at 4+8+8+32+4 = 56
	_, err = f.WriteAt([]byte("x"), 57) 
	assert.NoError(t, err)
	f.Close()

	// 3. Re-read to trigger checksum error or hash mismatch
	f2, err := os.Open(walPath)
	assert.NoError(t, err)
	defer f2.Close()

	block1_corrupted, err := ledger.DecodeBlock(f2)
	// Because of CRC32 checksum at the end of the block, it should fail immediately at DecodeBlock
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checksum mismatch")
	assert.Nil(t, block1_corrupted)
}
