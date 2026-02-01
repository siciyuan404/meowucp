CREATE INDEX IF NOT EXISTS orders_user_id_created_at_idx
  ON orders (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS payments_order_id_idx
  ON payments (order_id);

CREATE INDEX IF NOT EXISTS products_sku_idx
  ON products (sku);
