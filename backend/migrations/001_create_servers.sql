CREATE TABLE IF NOT EXISTS servers (
    id BIGSERIAL PRIMARY KEY,
    hostname VARCHAR(255) NOT NULL UNIQUE,
    ip INET NOT NULL,
    os VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'unknown',
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_servers_status ON servers(status);
CREATE INDEX IF NOT EXISTS idx_servers_last_seen ON servers(last_seen);
