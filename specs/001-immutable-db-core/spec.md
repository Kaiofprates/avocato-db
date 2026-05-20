# Feature Specification: avocato-db Core Engine

**Feature Branch**: `001-immutable-db-core`  
**Created**: 2026-05-20  
**Status**: Draft  
**Input**: User description: "O avocato-db é um motor de banco de dados imutável, leve e focado estritamente em operações append-only protegidas por hashes criptográficos, rodando inteiramente de forma isolada dentro de um container Docker. O projeto deve seguir à risca as diretrizes da 'avocato-db Constitution v1.1.0', priorizando determinismo de dados, concorrência e performance."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Secure Data Append (Priority: P1)

As a data producer, I want to append records to the database so that I can store information in an immutable and chronologically ordered manner without risk of modification.

**Why this priority**: This is the core functionality of the engine. Without secure append, the database has no purpose.

**Independent Test**: Can be tested by sending a series of data records to the engine and verifying they are stored sequentially with correct cryptographic links.

**Acceptance Scenarios**:

1. **Given** a running avocato-db container, **When** I send a valid data record, **Then** the system returns a success status along with the cryptographic hash of the new record.
2. **Given** multiple concurrent append requests, **When** they are processed, **Then** all records are stored sequentially without data loss or race conditions.

---

### User Story 2 - Cryptographic Integrity Verification (Priority: P2)

As an auditor, I want to verify the integrity of the entire database chain so that I can be certain that no data has been tampered with or deleted since it was written.

**Why this priority**: Immutability is useless without a way to verify it. This ensures the "trustless" nature of the system.

**Independent Test**: Can be tested by manually altering a bit in the stored data file and running a verification tool, which must report an integrity failure.

**Acceptance Scenarios**:

1. **Given** a database with existing records, **When** I trigger a full integrity check, **Then** the system validates each block's hash against its predecessor and reports "Valid".
2. **Given** a database where a record was modified outside the engine, **When** the container starts up, **Then** the Boot Integrity Check fails and prevents the API from opening.

---

### User Story 3 - High-Performance Read by Hash (Priority: P3)

As a data consumer, I want to retrieve specific records by their cryptographic hash so that I can access stored information with minimal latency.

**Why this priority**: While append-only is the primary focus, the data must be retrievable to be useful for applications.

**Independent Test**: Can be tested by querying a known hash and measuring the response time.

**Acceptance Scenarios**:

1. **Given** a record exists with hash `H`, **When** I query for hash `H`, **Then** the system returns the data in under 10ms.

---

### Edge Cases

- **Boundary Condition**: What happens when the storage volume reaches its capacity limit?
- **Error Scenario**: How does the system handle a crash during the exact moment of a write operation (Atomicity/Crash-consistency)?
- **Concurrency**: How does the system handle 10,000 simultaneous append requests?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an append-only API (gRPC or lightweight HTTP) that accepts raw data payloads.
- **FR-002**: System MUST implement a "Chain Coupling" mechanism where each record contains the hash of the previous record.
- **FR-003**: System MUST use SHA-256 or Blake3 for all cryptographic hash calculations.
- **FR-004**: System MUST guarantee $O(1)$ write performance by appending directly to the end of the physical log file.
- **FR-005**: System MUST perform a full cryptographic integrity scan upon every Docker container startup (Boot Integrity Check).
- **FR-006**: System MUST ensure deterministic serialization (e.g., fixed key ordering) before hashing data payloads.
- **FR-007**: System MUST NOT expose any methods for `UPDATE`, `DELETE`, or `TRUNCATE` operations.
- **FR-UX**: System MUST provide clear, explicit error messages for any integrity violations or resource exhaustion.
- **FR-PERF**: System MUST maintain latency below 10ms for index-based or hash-based read operations.

### Key Entities *(include if feature involves data)*

- **Block/Record**: The fundamental unit of data. Contains the payload, a timestamp, the hash of the previous block, and its own cryptographic hash.
- **Ledger/Log**: The sequential collection of blocks persisted on disk.
- **Hash Index**: A memory-mapped or structured index mapping hashes to physical file offsets for $O(1)$ lookups.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: System sustains a throughput of at least 500 append operations per second on standard SSD hardware.
- **SC-002**: 100% of data corruption attempts (simulated via mutation) are detected during the Boot Integrity Check.
- **SC-003**: Cold startup time (including Integrity Check) is under 5 seconds for a ledger containing 1,000,000 records.
- **SC-004**: Average retrieval time for a record by hash is consistently under 10ms.

## Assumptions

- **A-001**: The system will run in a Docker environment with persistent volume mapping for the ledger.
- **A-002**: Go (Golang) will be used as the implementation language for the core engine as per the Constitution.
- **A-003**: Data payloads are relatively small (e.g., under 1MB per record) to maintain performance targets.
- **A-004**: The host system provides sufficient entropy for cryptographic operations.
