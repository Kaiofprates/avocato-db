package ledger

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
)

type Block struct {
	Index     uint64
	Timestamp int64
	PrevHash  [32]byte
	Payload   []byte
	Hash      [32]byte
}

// Binary encoding format:
// [uint32:length][uint64:index][int64:ts][32b:prev_hash][uint32:payload_len][varb:payload][32b:hash][uint32:checksum]

func (b *Block) Encode(w io.Writer) error {
	payloadLen := uint32(len(b.Payload))
	totalLen := 4 + 8 + 8 + 32 + 4 + payloadLen + 32 + 4

	buf := make([]byte, totalLen)
	binary.BigEndian.PutUint32(buf[0:4], totalLen)
	binary.BigEndian.PutUint64(buf[4:12], b.Index)
	binary.BigEndian.PutUint64(buf[12:20], uint64(b.Timestamp))
	copy(buf[20:52], b.PrevHash[:])
	binary.BigEndian.PutUint32(buf[52:56], payloadLen)
	copy(buf[56:56+payloadLen], b.Payload)
	copy(buf[56+payloadLen:88+payloadLen], b.Hash[:])

	// Calculate checksum of everything except the last 4 bytes
	checksum := crc32.ChecksumIEEE(buf[0 : totalLen-4])
	binary.BigEndian.PutUint32(buf[totalLen-4:totalLen], checksum)

	_, err := w.Write(buf)
	return err
}

func DecodeBlock(r io.Reader) (*Block, error) {
	var totalLen uint32
	if err := binary.Read(r, binary.BigEndian, &totalLen); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("checksum mismatch: expected %x, got %x", expectedChecksum, actualChecksum)
	}

	b := &Block{}
	b.Index = binary.BigEndian.Uint64(buf[4:12])
	b.Timestamp = int64(binary.BigEndian.Uint64(buf[12:20]))
	copy(b.PrevHash[:], buf[20:52])
	payloadLen := binary.BigEndian.Uint32(buf[52:56])
	b.Payload = make([]byte, payloadLen)
	copy(b.Payload, buf[56:56+payloadLen])
	copy(b.Hash[:], buf[56+payloadLen:88+payloadLen])

	return b, nil
}
