package ledger

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"github.com/google/uuid"
)

type Block struct {
	Index     uint64
	UUID      string
	Timestamp int64
	PrevHash  [32]byte
	Payload   []byte
	Hash      [32]byte
}

// Binary encoding format:
// [uint32:length][uint64:index][16b:uuid][int64:ts][32b:prev_hash][uint32:payload_len][varb:payload][32b:hash][uint32:checksum]

func (b *Block) Encode(w io.Writer) error {
	payloadLen := uint32(len(b.Payload))
	totalLen := 4 + 8 + 16 + 8 + 32 + 4 + payloadLen + 32 + 4

	buf := make([]byte, totalLen)
	binary.BigEndian.PutUint32(buf[0:4], totalLen)
	binary.BigEndian.PutUint64(buf[4:12], b.Index)
	
	uid, err := uuid.Parse(b.UUID)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	copy(buf[12:28], uid[:])
	
	binary.BigEndian.PutUint64(buf[28:36], uint64(b.Timestamp))
	copy(buf[36:68], b.PrevHash[:])
	binary.BigEndian.PutUint32(buf[68:72], payloadLen)
	copy(buf[72:72+payloadLen], b.Payload)
	copy(buf[72+payloadLen:104+payloadLen], b.Hash[:])

	// Calculate checksum
	checksum := crc32.ChecksumIEEE(buf[0 : totalLen-4])
	binary.BigEndian.PutUint32(buf[totalLen-4:totalLen], checksum)

	_, err = w.Write(buf)
	return err
}

func DecodeBlock(r io.Reader) (*Block, error) {
	var totalLen uint32
	if err := binary.Read(r, binary.BigEndian, &totalLen); err != nil {
		return nil, err
	}

	// Minimum possible length: 4 (totalLen) + 8 (index) + 16 (uuid) + 8 (ts) + 32 (prevHash) + 4 (payloadLen) + 0 (payload) + 32 (hash) + 4 (checksum) = 108
	if totalLen < 108 {
		return nil, fmt.Errorf("block too small: %d", totalLen)
	}

	// Prevent huge allocations (e.g., 100MB limit)
	if totalLen > 100*1024*1024 {
		return nil, fmt.Errorf("block too large: %d", totalLen)
	}

	buf := make([]byte, totalLen)
	binary.BigEndian.PutUint32(buf[0:4], totalLen)
	if _, err := io.ReadFull(r, buf[4:]); err != nil {
		return nil, err
	}

	// Verify checksum
	expectedChecksum := binary.BigEndian.Uint32(buf[totalLen-4 : totalLen])
	actualChecksum := crc32.ChecksumIEEE(buf[0 : totalLen-4])
	if expectedChecksum != actualChecksum {
		return nil, fmt.Errorf("checksum mismatch")
	}

	b := &Block{}
	b.Index = binary.BigEndian.Uint64(buf[4:12])
	
	uid, err := uuid.FromBytes(buf[12:28])
	if err != nil {
		return nil, fmt.Errorf("invalid uuid in block: %w", err)
	}
	b.UUID = uid.String()
	
	b.Timestamp = int64(binary.BigEndian.Uint64(buf[28:36]))
	copy(b.PrevHash[:], buf[36:68])
	
	payloadLen := binary.BigEndian.Uint32(buf[68:72])
	
	// Validate payloadLen against totalLen
	// totalLen = header(72) + payload(payloadLen) + hash(32) + checksum(4)
	if 72+payloadLen+32+4 != totalLen {
		return nil, fmt.Errorf("payload length mismatch: payload %d, total %d", payloadLen, totalLen)
	}

	b.Payload = make([]byte, payloadLen)
	copy(b.Payload, buf[72:72+payloadLen])
	copy(b.Hash[:], buf[72+payloadLen:104+payloadLen])

	return b, nil
}
