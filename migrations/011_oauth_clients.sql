CREATE TABLE IF NOT EXISTS oauth_clients (
  id BIGSERIAL PRIMARY KEY,
  client_id TEXT NOT NULL UNIQUE,
  secret_hash TEXT NOT NULL,
  scopes TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth_tokens (
  id BIGSERIAL PRIMARY KEY,
  token TEXT NOT NULL UNIQUE,
  client_id TEXT NOT NULL,
  user_id BIGINT,
  scopes TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ
);
