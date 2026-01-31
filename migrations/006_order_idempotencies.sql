CREATE TABLE IF NOT EXISTS order_idempotencies (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    idempotency_key TEXT NOT NULL,
    order_id BIGINT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS order_idempotencies_user_key_uidx
    ON order_idempotencies (user_id, idempotency_key);
