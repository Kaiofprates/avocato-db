# Quickstart: avocato-db

## Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local development)

## Running with Docker

1. **Start the database and backend**:
   ```bash
   docker-compose up -d
   ```

2. **Wait for Boot Integrity Check**:
   Check logs to ensure the chain is verified:
   ```bash
   docker-compose logs -f backend
   ```
   You should see: `[avocato-db] Boot Integrity Check PASSED. API ready on :8080`

## Basic Operations

### Append a Record
```bash
curl -X POST http://localhost:8080/v1/append \
     -H "Content-Type: application/json" \
     -d '{"payload": {"msg": "First Immutable Record"}}'
```

### Verify Integrity
```bash
curl http://localhost:8080/v1/integrity
```

### Retrieve by Index or Hash
```bash
# By Index
curl http://localhost:8080/v1/block/1

# By Hash
curl http://localhost:8080/v1/block/<hash>
```

### Generate Proof (Stub)
```bash
curl http://localhost:8080/v1/proof/1
```

## Security Note
The engine is strictly append-only. Any attempt to modify the ledger file or PostgreSQL records directly will cause the Boot Integrity Check to fail on the next restart.
