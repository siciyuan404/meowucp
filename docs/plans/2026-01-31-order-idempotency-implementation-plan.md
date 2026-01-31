# Order Idempotency and Concurrency Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add idempotent order creation using client-provided keys while preserving atomic stock updates.

**Architecture:** Store idempotency records in a dedicated table keyed by `(user_id, idempotency_key)`; gate order creation in `OrderService` within the existing transaction, returning the prior order when the key is reused. Inventory remains concurrency-safe via `UpdateStockWithDelta`.

**Tech Stack:** Go, GORM, PostgreSQL, Gin

---

### Task 1: Add idempotency data model and migration

**Files:**
- Create: `migrations/006_order_idempotencies.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write the migration**

```sql
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
```

**Step 2: Add domain model**

```go
type OrderIdempotency struct {
  ID             int64 `gorm:"primary_key"`
  UserID         int64 `gorm:"index;not null"`
  IdempotencyKey string `gorm:"not null"`
  OrderID        *int64
  Status         string `gorm:"not null;default:'pending'"`
  CreatedAt      time.Time
  UpdatedAt      time.Time
}
```

**Step 3: Commit**

```bash
git add migrations/006_order_idempotencies.sql internal/domain/models.go
git commit -m "feat(order): add idempotency model and migration"
```

### Task 2: Add repository interface and implementation

**Files:**
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/order_idempotency_repository.go`

**Step 1: Add interface**

```go
type OrderIdempotencyRepository interface {
  Create(record *domain.OrderIdempotency) error
  FindByUserAndKey(userID int64, key string) (*domain.OrderIdempotency, error)
  Update(record *domain.OrderIdempotency) error
}
```

**Step 2: Register in repositories**

```go
OrderIdempotency OrderIdempotencyRepository
```

**Step 3: Implement repository**

```go
func (r *orderIdempotencyRepository) Create(record *domain.OrderIdempotency) error {
  return r.db.Create(record).Error
}

func (r *orderIdempotencyRepository) FindByUserAndKey(userID int64, key string) (*domain.OrderIdempotency, error) {
  var record domain.OrderIdempotency
  if err := r.db.Where("user_id = ? AND idempotency_key = ?", userID, key).First(&record).Error; err != nil {
    return nil, err
  }
  return &record, nil
}

func (r *orderIdempotencyRepository) Update(record *domain.OrderIdempotency) error {
  return r.db.Save(record).Error
}
```

**Step 4: Commit**

```bash
git add internal/repository/repository.go internal/repository/order_idempotency_repository.go
git commit -m "feat(order): add idempotency repository"
```

### Task 3: Add idempotency flow in OrderService (TDD)

**Files:**
- Modify: `internal/service/order_service.go`
- Modify: `internal/service/order_service_test.go`
- Modify: `internal/service/user_service.go` (Services wiring)

**Step 1: Write failing tests**

```go
func TestOrderServiceCreateOrderReturnsExistingOrderForIdempotencyKey(t *testing.T) {
  // create idempotency record with order_id set, verify CreateOrder returns that order
}

func TestOrderServiceCreateOrderReturnsConflictForPendingIdempotencyKey(t *testing.T) {
  // idempotency record exists with nil order_id, expect conflict error
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/service -run Idempotency`
Expected: FAIL with missing idempotency handling

**Step 3: Implement minimal idempotency flow**

```go
type OrderService struct {
  // add idempotencyRepo
}

func (s *OrderService) CreateOrder(userID int64, idempotencyKey string, shippingAddress, billingAddress, notes string, paymentMethod string) (*domain.Order, error) {
  // in transaction: create/find idempotency record
  // if order_id present -> return existing
  // if pending -> return conflict error
  // else proceed with order creation and set order_id
}
```

**Step 4: Update tests to use the new signature**

Update all existing `CreateOrder` call sites in tests to pass an idempotency key (use empty string where not needed).

**Step 5: Run tests to verify they pass**

Run: `go test ./internal/service -run Idempotency`
Expected: PASS

**Step 6: Commit**

```bash
git add internal/service/order_service.go internal/service/order_service_test.go internal/service/user_service.go
git commit -m "feat(order): add idempotent order creation"
```

### Task 4: API integration (if applicable)

**Files:**
- Modify: `internal/api` order creation handler (if present)

**Step 1: Add idempotency key extraction**

```go
key := c.GetHeader("Idempotency-Key")
```

**Step 2: Pass key into order service**

```go
order, err := orderService.CreateOrder(userID, key, ...)
```

**Step 3: Add tests for duplicate requests**

Run: `go test ./internal/api -run Idempotency`
Expected: PASS

**Step 4: Commit**

```bash
git add internal/api
git commit -m "feat(api): support idempotent order creation"
```

### Task 5: End-to-end verification

**Step 1: Run relevant tests**

Run: `go test ./internal/service`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
git add -A
git commit -m "chore: finalize idempotent order changes"
```
