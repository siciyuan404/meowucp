CREATE TABLE IF NOT EXISTS data_retention_policies (
  id BIGSERIAL PRIMARY KEY,
  entity TEXT NOT NULL,
  ttl_days INT NOT NULL,
  strategy TEXT NOT NULL
);
