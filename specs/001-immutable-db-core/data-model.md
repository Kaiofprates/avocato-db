# Data Model: avocato-db Core

## Entities

### 1. Block (The Ledger Entry)
The atomic unit of data in the system. Each block is cryptographically linked to its predecessor.

| Field | Type | Description |
|-------|------|-------------|
| `index` | `uint64` | Monotonically increasing sequence number. |
| `timestamp` | `int64` | Unix timestamp in nanoseconds (deterministic). |
| `prev_hash` | `string` | SHA-256 hex string of the previous block's hash. |
| `payload` | `JSON` | User-provided data. Must be JCS-canonicalized before hashing. |
| `hash` | `string` | SHA-256 hex string of the canonicalized block. |

**Hash Calculation**:
`hash = SHA256(JCS(index + timestamp + prev_hash + JCS(payload)))`

### 2. MMR (Merkle Mountain Range)
A companion structure for fast integrity proofs.
- **Peaks**: The list of Merkle roots for each "mountain" in the range.
- **Root**: The "bagged" root of all peaks.

### 3. Checkpoint
Persisted state used to accelerate the Boot Integrity Check.

| Field | Type | Description |
|-------|------|-------------|
| `last_index` | `uint64` | Index of the last verified block in this checkpoint. |
| `mmr_root` | `string` | The MMR root at this index. |
| `last_hash` | `string` | The hash of the block at `last_index`. |
| `signature` | `string` | (Optional) Internal signature for checkpoint authenticity. |

## Persistence Mapping

### PostgreSQL (Secondary Storage & Indexing)
```sql
CREATE TABLE blocks (
    index BIGINT PRIMARY KEY,
    hash TEXT UNIQUE NOT NULL,
    prev_hash TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    payload JSONB NOT NULL
);

CREATE INDEX idx_blocks_hash ON blocks (hash);
CREATE INDEX idx_blocks_brin ON blocks USING BRIN (index);
```

### WAL (Primary Append Log)
- **Format**: Binary flat file.
- **Record**: `[uint32:length][uint64:index][int64:ts][32b:prev_hash][varb:payload][32b:hash][uint32:checksum]`
- **Sync Strategy**: O_APPEND with Group Commit.
