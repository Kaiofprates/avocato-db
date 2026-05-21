# API Contracts: avocato-db REST API

## 1. Append Data
**Endpoint**: `POST /v1/append`

**Request Body**:
```json
{
  "payload": {
    "key": "value",
    "metadata": { "user": "alice" }
  }
}
```

**Response (201 Created)**:
```json
{
  "status": "success",
  "data": {
    "index": 123,
    "hash": "a5f3...e21",
    "prev_hash": "b2e1...f0a",
    "timestamp": 1716200000000000000
  }
}
```

---

## 2. Retrieve Block
**Endpoint**: `GET /v1/block/{index_or_hash}`

**Response (200 OK)**:
```json
{
  "index": 123,
  "hash": "a5f3...e21",
  "prev_hash": "b2e1...f0a",
  "timestamp": 1716200000000000000,
  "payload": {
    "key": "value",
    "metadata": { "user": "alice" }
  }
}
```

---

## 3. Integrity Status
**Endpoint**: `GET /v1/integrity`

**Response (200 OK)**:
```json
{
  "total_blocks": 1234,
  "mmr_root": "d41d...8cd",
  "last_checkpoint": {
    "index": 1000,
    "timestamp": 1716100000000000000
  },
  "status": "VALID"
}
```

---

## 4. Verification Proof
**Endpoint**: `GET /v1/proof/{index}`

**Description**: Returns an MMR proof that the block at `{index}` is part of the chain root.

**Response (200 OK)**:
```json
{
  "index": 123,
  "root": "d41d...8cd",
  "proof": [
    "e3b0...c1a",
    "81fe...221"
  ]
}
```
