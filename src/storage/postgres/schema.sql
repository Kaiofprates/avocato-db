-- avocato-db Schema
-- Version: 1.1.0 (Added UUID)

-- 1. Create the Immutable Ledger Table
CREATE TABLE IF NOT EXISTS blocks (
    index BIGINT PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL,
    hash TEXT UNIQUE NOT NULL,
    prev_hash TEXT NOT NULL,
    created_at BIGINT NOT NULL,
    payload JSONB NOT NULL
);

-- 2. Create Checkpoints Table for Fast Boot
CREATE TABLE IF NOT EXISTS checkpoints (
    last_index BIGINT PRIMARY KEY,
    mmr_root TEXT NOT NULL,
    last_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Prevent Modifications (Immutability Gate)
CREATE OR REPLACE FUNCTION block_immutable_changes() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'avocato-db: Table is immutable. UPDATE and DELETE are strictly forbidden by Constitution v1.1.0.';
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_lock_ledger') THEN
        CREATE TRIGGER trg_lock_ledger
        BEFORE UPDATE OR DELETE ON blocks
        FOR EACH ROW EXECUTE FUNCTION block_immutable_changes();
    END IF;
END $$;

-- 4. High-Performance Indexing
-- BRIN is efficient for naturally ordered data like the ledger index.
CREATE INDEX IF NOT EXISTS idx_blocks_brin ON blocks USING BRIN (index);
CREATE INDEX IF NOT EXISTS idx_blocks_hash ON blocks (hash);
CREATE INDEX IF NOT EXISTS idx_blocks_uuid ON blocks (uuid);
