CREATE TABLE IF NOT EXISTS tax_rules (
  id BIGSERIAL PRIMARY KEY,
  region TEXT NOT NULL,
  category TEXT NOT NULL,
  rate NUMERIC(6,4) NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS shipping_rules (
  id BIGSERIAL PRIMARY KEY,
  region TEXT NOT NULL,
  method TEXT NOT NULL,
  base_amount NUMERIC(12,2) NOT NULL,
  per_item_amount NUMERIC(12,2) NOT NULL
);
