ALTER TABLE ucp_webhook_jobs ADD COLUMN IF NOT EXISTS last_error TEXT;
ALTER TABLE ucp_webhook_jobs ADD COLUMN IF NOT EXISTS last_attempt_at TIMESTAMP;

CREATE TABLE IF NOT EXISTS ucp_webhook_alerts (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(64),
    reason VARCHAR(64) NOT NULL,
    details TEXT,
    attempts INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
