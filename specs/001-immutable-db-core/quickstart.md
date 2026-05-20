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

### Retrieve by Hash
```bash
curl http://localhost:8080/v1/block/<hash>
```

## Security Note
The engine is strictly append-only. Any attempt to modify the ledger file or PostgreSQL records directly will cause the Boot Integrity Check to fail on the next restart.
