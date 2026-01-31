# Order Idempotency and Concurrency Design

Date: 2026-01-31

## Background and Goals

We need to complete the order pipeline by adding concurrency safety and idempotent order creation. The goal is to prevent overselling under concurrent requests and to avoid duplicate orders on retries. The solution must be transactional and auditable.

## Scope and Non-Goals

Scope:
- Add atomic stock updates for concurrent safety.
- Add an idempotency record table keyed by user and idempotency key.
- Enforce idempotent create order behavior in the order service.

Non-goals:
- Full retry recovery for incomplete idempotency records.
- Background reconciliation jobs.
- Payment flow changes.

## Approach Summary

Use database atomic stock updates to prevent oversell and add a dedicated idempotency table with a unique constraint on `(user_id, idempotency_key)`. The order service will create or read the idempotency record inside the order transaction and return the existing order when the key was already used.

## Data Model

Add table `order_idempotencies` with:
- `id` (PK)
- `user_id`
- `idempotency_key`
- `order_id` (nullable)
- `status` (optional: pending/completed)
- `created_at`, `updated_at`

Unique index: `(user_id, idempotency_key)`.

## Flow and Transaction Boundary

1) Start transaction in `OrderService.CreateOrder`.
2) Create or fetch idempotency record by `(user_id, key)`.
   - If record exists and `order_id` is set, return that order.
   - If record exists but `order_id` is empty, return a conflict/pending error.
3) Read cart and product snapshots.
4) Validate stock, then create order and order items.
5) Adjust stock via `InventoryService.AdjustStock` using atomic update path.
6) Increment sales.
7) Clear cart.
8) Update idempotency record with `order_id` and status `completed`.
9) Commit.

Any failure rolls back the transaction, keeping data consistent.

## Error Handling

- Insufficient stock: return `insufficient_stock`, no order created.
- Idempotency record exists with order: return existing order.
- Idempotency record exists but pending: return `idempotency_conflict` or `order_pending`.
- Database or inventory failures: return a generic error and roll back.

## Testing Strategy

Service tests:
- Same idempotency key returns the same order without double stock adjustments.
- Pending idempotency record returns conflict.
- Atomic stock update failure rolls back order and items.

Repository tests:
- Atomic update rejects when stock is insufficient.
- Atomic update succeeds when stock is sufficient.

Transaction tests:
- Rollback keeps inventory logs, cart, and order tables consistent.

## Future Extensions

- Automatic recovery for pending idempotency records.
- Periodic cleanup of old idempotency keys.
