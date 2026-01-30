-- UCP checkout sessions
CREATE TABLE IF NOT EXISTS checkout_sessions (
    id VARCHAR(64) PRIMARY KEY,
    status VARCHAR(32) NOT NULL,
    currency VARCHAR(8) NOT NULL,
    line_items JSONB NOT NULL,
    totals JSONB NOT NULL,
    buyer JSONB,
    messages JSONB,
    links JSONB,
    continue_url TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- UCP payment handlers
CREATE TABLE IF NOT EXISTS payment_handlers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    version VARCHAR(32) NOT NULL,
    spec TEXT NOT NULL,
    config_schema TEXT NOT NULL,
    config JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- UCP webhook events for idempotency
CREATE TABLE IF NOT EXISTS ucp_webhook_events (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(64) UNIQUE NOT NULL,
    event_type VARCHAR(64) NOT NULL,
    order_id VARCHAR(64) NOT NULL,
    status VARCHAR(32) NOT NULL,
    payload_hash VARCHAR(64) NOT NULL,
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);

-- Extend payments for handler payload
ALTER TABLE payments ADD COLUMN IF NOT EXISTS payment_payload JSONB;
