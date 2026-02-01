CREATE TABLE IF NOT EXISTS idempotency_keys (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  key TEXT NOT NULL,
  request_hash TEXT NOT NULL,
  response_snapshot JSONB,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idempotency_keys_user_key_uidx
  ON idempotency_keys (user_id, key);
