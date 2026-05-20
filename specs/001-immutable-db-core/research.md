# Research Report: avocato-db Core Engine

This report details the technical foundations for building `avocato-db`, a high-performance, immutable blockchain-inspired database engine.

---

## 1. Deterministic Serialization in Go (RFC 8785)

**Context**: Cryptographic consistency requires that any given object always serializes to the exact same byte sequence. Standard JSON encoders (like Go's `encoding/json`) do not guarantee key order, which breaks hashing.

### Concrete Library Recommendations
- **`github.com/lattice-substrate/json-canon` (Recommended)**: 
  - **Why**: Specifically designed for cryptographic use. It includes a strict parser that rejects ambiguous JSON (e.g., duplicate keys, leading zeros) and implements the ECMA-262 compatible number formatting required by RFC 8785.
  - **Best for**: Infrastructure-grade security.
- **`github.com/ucarion/jcs`**:
  - **Why**: More idiomatic and easy to use if you are working with `interface{}` or custom structs.
  - **Best for**: Rapid development where strict input validation is handled elsewhere.

### Code Pattern
```go
import "github.com/lattice-substrate/json-canon/jcs"

func CanonicalHash(data interface{}) ([]byte, error) {
    // 1. Serialize to standard JSON first or pass raw bytes
    raw, _ := json.Marshal(data)
    
    // 2. Canonicalize according to RFC 8785
    canonical, err := jcs.Canonicalize(raw)
    if err != nil {
        return nil, err
    }
    
    // 3. Hash the deterministic output
    hash := sha256.Sum256(canonical)
    return hash[:], nil
}
```

---

## 2. PostgreSQL as an Immutable Ledger

**Context**: Using PostgreSQL as a storage backend provides ACID compliance while triggers and specific indexing strategies ensure immutability and performance.

### Architectural Recommendations
- **Trigger-Based Immutability**: Block `UPDATE` and `DELETE` at the database level to ensure data integrity even against application-level bugs.
- **BRIN (Block Range Indexing)**:
  - **Why**: Ledgers are naturally ordered (by ID or timestamp). B-Tree indexes for millions of rows become massive. BRIN indexes are ~100x smaller and significantly faster for append-only workloads.
  - **Config**: Set `autosummarize = on` to ensure the index stays current as data is appended.

### Implementation Pattern
```sql
-- 1. Create the Immutable Ledger
CREATE TABLE ledger (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    payload JSONB NOT NULL,
    hash BYTEA NOT NULL
);

-- 2. Prevent Modifications
CREATE OR REPLACE FUNCTION block_immutable_changes() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Table is immutable. Update/Delete not allowed.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_lock_ledger
BEFORE UPDATE OR DELETE ON ledger
FOR EACH ROW EXECUTE FUNCTION block_immutable_changes();

-- 3. High-Performance Indexing
CREATE INDEX idx_ledger_brin ON ledger USING BRIN (id, created_at) WITH (autosummarize = on);
```

---

## 3. Boot Integrity Check Performance

**Context**: Scanning 1M+ records at boot to verify chain integrity is O(n). To achieve fast boot, we use Merkle Mountain Ranges (MMR) and Checkpointing.

### Concrete Library
- **`github.com/discretemind/mmr`**: A high-performance MMR implementation in Go.

### Architectural Strategy
- **Merkle Mountain Range (MMR)**: Unlike standard Merkle Trees, MMRs are append-only and allow for efficient "bagging the peaks" to calculate a global root.
- **Checkpointing**:
  - Periodically (e.g., every 10,000 blocks) save a "State Snapshot" containing the MMR root and the last processed index.
  - At boot, the engine only verifies the chain from the **last trusted checkpoint** to the tip. This turns boot verification from O(n) to O(k), where k is the distance from the last checkpoint.
- **Pruning**: Old MMR nodes can be pruned if only the current root is needed for verification.

---

## 4. Go High-Performance File I/O

**Context**: Minimizing disk latency while ensuring durability (WAL).

### Architectural Recommendations
- **Group Commit**: Instead of calling `fsync` for every write, batch multiple concurrent requests.
- **mmap vs. Standard I/O**:
  - **Use Standard I/O (O_APPEND)**: Best for the Ledger/WAL. It's safer, simpler, and better aligned with the Go scheduler (which understands blocking syscalls but not page faults).
  - **Use mmap**: Best for reading large indexes or state trees where random access is frequent. It bypasses Go's GC pressure for large datasets.

### Group Commit Code Pattern
```go
func (w *WAL) runFlusher() {
    const linger = 5 * time.Millisecond
    for {
        batch := collectRequests(w.queue, linger)
        if len(batch) == 0 { continue }

        for _, req := range batch {
            w.writer.Write(req.Data)
        }
        w.writer.Flush()
        w.file.Sync() // Single fsync for the entire batch

        for _, req := range batch {
            req.Done <- nil
        }
    }
}
```

### Recommendation Summary
- **Storage**: PostgreSQL with BRIN + Custom WAL for high-throughput appends.
- **Integrity**: MMR with Checkpointing every 50k records.
- **Serialization**: `json-canon` for all hashed structures.
- **Concurrency**: Group Commit with 5-10ms linger time for O(1) append performance under load.
