CREATE TABLE IF NOT EXISTS ucp_webhook_audit (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(64),
    reason VARCHAR(64) NOT NULL,
    signature_header TEXT,
    key_id VARCHAR(128),
    payload_hash VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
