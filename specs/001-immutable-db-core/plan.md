# Implementation Plan: avocato-db Core Engine

**Branch**: `001-immutable-db-core` | **Date**: 2026-05-20 | **Spec**: [specs/001-immutable-db-core/spec.md](spec.md)
**Input**: Feature specification from `/specs/001-immutable-db-core/spec.md`

## Summary

The `avocato-db` core engine will be implemented as a high-performance Go-based backend service running in Docker. It uses a **Group Commit** pattern for O(1) append performance and **Merkle Mountain Ranges (MMR)** for efficient cryptographic integrity verification. Data persistence follows a hybrid approach: a high-speed custom **Write-Ahead Log (WAL)** for the immutable ledger and **PostgreSQL** with **BRIN indexing** for efficient querying and secondary storage, as requested by the user.

## Technical Context

**Language/Version**: Go 1.22+  
**Primary Dependencies**: `github.com/lattice-substrate/json-canon` (RFC 8785), `github.com/discretemind/mmr` (Integrity), `github.com/jackc/pgx/v5` (Postgres driver)  
**Storage**: PostgreSQL 16 (BRIN indexing, Trigger-based immutability) + Go-managed binary WAL  
**Testing**: `go test -race`, Integration tests (Crash-consistency injection), Cryptographic mutation testing  
**Target Platform**: Docker (Alpine-based, non-root user)
**Project Type**: backend service / database engine  
**Performance Goals**: 500+ append/sec (Group Commit batching), <10ms P99 read latency  
**Constraints**: Append-only (no UPDATE/DELETE), Deterministic JCS serialization, Boot Integrity Check  
**Scale/Scope**: Support for 1M+ records per ledger instance with <5s boot verification time via MMR checkpoints.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Code Quality**: Architecture follows SOLID; clear separation between Storage and Integrity layers.
- [x] **Testing Standards**: Plan includes crash-consistency and cryptographic validation.
- [x] **UX Consistency**: API provides explicit integrity feedback (hashes).
- [x] **Performance**: Group commit and BRIN indexing directly address O(1) and latency requirements.
- [x] **Scalability**: Docker-native with volume-mapped storage supports vertical scaling and isolation.

## Project Structure

### Documentation (this feature)

```text
specs/001-immutable-db-core/
├── plan.md              # This file
├── research.md          # Research on JCS, MMR, and Postgres optimization
├── data-model.md        # Block and MMR structure (Phase 1)
├── quickstart.md        # API usage and Docker setup (Phase 1)
├── contracts/           # gRPC/REST API definitions (Phase 1)
└── tasks.md             # Implementation tasks (Phase 2)
```

### Source Code (repository root)

```text
src/
├── api/
│   ├── handlers/        # API route handlers
│   └── proto/           # gRPC definitions (if used)
├── core/
│   ├── ledger/          # WAL, Group Commit, Block Appender
│   ├── integrity/       # MMR peaks, Checkpointing, Hash Chain logic
│   └── crypto/          # JCS Canonicalization, Hashing
├── storage/
│   ├── postgres/        # DB schema, triggers, BRIN indexes
│   └── wal/             # Low-level file I/O and sync
├── config/              # Environment-based config
└── main.go              # Entry point and Boot Integrity Check trigger

tests/
├── integration/         # Integrity breaking tests, Docker-based tests
└── internal/            # Unit tests for core logic
```

**Structure Decision**: Hybrid structure separating the Go-native engine (core/) from the persistence layer (storage/) and API layer (api/). This allows testing the integrity logic independently of Postgres.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Hybrid Storage (WAL + Postgres) | User requested Postgres; Constitution mandates Log Appender. | Pure Postgres might struggle with O(1) append at very high scale without complex tuning; Pure WAL lacks rich querying. |
| MMR Checkpoints | 1M records requirement for <5s boot. | Full O(n) scan would exceed the 5s boot target as the ledger grows. |
