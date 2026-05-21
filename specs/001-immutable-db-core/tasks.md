# Tasks: avocato-db Core Engine

**Input**: Design documents from `/specs/001-immutable-db-core/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/rest-api.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Initialize Go project with `go mod init avocato-db`
- [x] T002 [P] Create project structure `src/core`, `src/storage`, `src/api`, `src/config`
- [x] T003 [P] Add primary dependencies: `json-canon`, `mmr`, `pgx/v5`
- [x] T004 Create `docker-compose.yml` with PostgreSQL 16 and backend service
- [x] T005 [P] Setup `src/config/config.go` for environment variable management

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure for Hybrid Storage and Deterministic Serialization

- [x] T006 Implement deterministic JCS serialization wrapper in `src/core/crypto/jcs.go`
- [x] T007 Setup PostgreSQL schema with immutability triggers in `src/storage/postgres/schema.sql`
- [x] T008 Implement low-level WAL file writer with `O_APPEND` in `src/storage/wal/writer.go`
- [x] T009 [P] Implement Group Commit coordinator in `src/core/ledger/group_commit.go`
- [x] T010 Create base `Block` struct and binary encoding in `src/core/ledger/block.go`
- [x] T011 [P] Setup central error handling and logging in `src/core/logger.go`

**Checkpoint**: Foundation ready - hybrid storage and crypto primitives are available.

---

## Phase 3: User Story 1 - Secure Data Append (Priority: P1) 🎯 MVP

**Goal**: Append records to the ledger with O(1) performance and cryptographic linking.

**Independent Test**: Send POST requests to `/v1/append` and verify binary WAL and Postgres table growth.

### Implementation for User Story 1

- [x] T012 [US1] Implement hash calculation logic in `src/core/ledger/hasher.go`
- [x] T013 [US1] Implement `AppendBlock` service logic in `src/core/ledger/service.go`
- [x] T014 [US1] Create REST API handler for `/v1/append` in `src/api/handlers/append.go`
- [x] T015 [US1] Integrate WAL and Postgres persistence in the append flow
- [x] T016 [US1] Implement response serialization with new block hash and index
- [ ] T017 [US1] Add integration test for concurrent appends in `tests/integration/append_test.go`

**Checkpoint**: User Story 1 (MVP) is fully functional. Records can be appended securely.

---

## Phase 4: User Story 2 - Cryptographic Integrity Verification (Priority: P2)

**Goal**: Verify the entire chain at boot and provide status/proofs.

**Independent Test**: Restart container and verify "Boot Integrity Check PASSED" in logs.

### Implementation for User Story 2

- [x] T018 [US2] Integrate Merkle Mountain Range (MMR) in `src/core/integrity/mmr_manager.go`
- [x] T019 [US2] Implement Checkpoint persistence in `src/storage/postgres/checkpoints.sql`
- [x] T020 [US2] Implement Boot Integrity Check loop in `src/core/integrity/boot_check.go`
- [x] T021 [US2] Add Checkpoint logic to skip O(n) scan in `src/core/integrity/checkpoints.go`
- [x] T022 [US2] Create REST API handler for `/v1/integrity` in `src/api/handlers/integrity.go`
- [x] T023 [US2] Create REST API handler for `/v1/proof/{index}` in `src/api/handlers/proof.go`
- [x] T024 [US2] Add integration test for chain corruption detection in `tests/integration/integrity_test.go`

**Checkpoint**: User Story 2 is complete. System ensures chain validity at every startup.

---

## Phase 5: User Story 3 - High-Performance Read by Hash (Priority: P3)

**Goal**: Retrieve records by index or hash with <10ms latency.

**Independent Test**: Query GET `/v1/block/{hash}` and measure response time.

### Implementation for User Story 3

- [x] T025 [US3] Implement Postgres lookup with BRIN index in `src/storage/postgres/reader.go`
- [x] T026 [US3] Create REST API handler for `/v1/block/{id}` in `src/api/handlers/get_block.go`
- [x] T027 [US3] Implement caching for the last N blocks in `src/core/ledger/cache.go`
- [x] T028 [US3] Add performance benchmark test for reads in `tests/bench/read_bench_test.go`

---

## Phase 6: Polish & Cross-Cutting Concerns

- [ ] T029 [P] Generate Swagger/OpenAPI documentation from contracts
- [x] T030 [P] Implement graceful shutdown handling in `main.go`
- [x] T031 Final code cleanup and linter audit (Principle I)
- [x] T032 [P] Update `quickstart.md` with final API examples

---

## Dependencies & Execution Order

### Phase Dependencies
- **Phase 1 & 2** are strictly required before any User Story.
- **Phase 3 (Append)** is required for Phase 4 & 5 to have data to work with.

### Parallel Opportunities
- T002-T005 can run in parallel.
- T009 (Group Commit) and T011 (Logger) can run in parallel while WAL is being built.
- Once Foundation is ready, US1 implementation can be split (API handlers vs Core logic).

---

## Implementation Strategy
1. **MVP (US1)**: Focus on reaching a state where `docker-compose up` allows appending a record that persists in both WAL and Postgres.
2. **Hardening (US2)**: Implement the "Boot Integrity Check" which is a core Constitution mandate.
3. **Optimization (US3)**: Polish the read paths and ensure latency targets.
