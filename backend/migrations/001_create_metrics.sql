CREATE TABLE IF NOT EXISTS metrics (
    id BIGSERIAL PRIMARY KEY,
    cpu NUMERIC(5, 2) NOT NULL,
    memory NUMERIC(5, 2) NOT NULL,
    disk NUMERIC(5, 2) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp);
