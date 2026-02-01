# Order Lifecycle Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add cancel/ship/receive transitions with inventory rollback and shipment tracking.

**Architecture:** Extend order domain with shipment + status logs, expose service methods for transitions, and add admin endpoints for lifecycle actions.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add shipment + status log models and migration

**Files:**
- Create: `migrations/008_order_lifecycle.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS shipments (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  carrier TEXT NOT NULL,
  tracking_no TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'created',
  shipped_at TIMESTAMPTZ,
  delivered_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_status_logs (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  from_status TEXT NOT NULL,
  to_status TEXT NOT NULL,
  reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Step 2: Add domain models**

```go
type Shipment struct {
  ID          int64 `gorm:"primary_key"`
  OrderID     int64 `gorm:"not null"`
  Carrier     string
  TrackingNo  string
  Status      string
  ShippedAt   *time.Time
  DeliveredAt *time.Time
  CreatedAt   time.Time
  UpdatedAt   time.Time
}

type OrderStatusLog struct {
  ID         int64 `gorm:"primary_key"`
  OrderID    int64 `gorm:"not null"`
  FromStatus string
  ToStatus   string
  Reason     string
  CreatedAt  time.Time
}
```

**Step 3: Commit**

```bash
git add migrations/008_order_lifecycle.sql internal/domain/models.go
git commit -m "feat(order): add shipment and status log models"
```

### Task 2: Add repositories + lifecycle service

**Files:**
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/shipment_repository.go`
- Create: `internal/repository/order_status_log_repository.go`
- Modify: `internal/service/order_service.go`
- Test: `internal/service/order_service_test.go`

**Step 1: Write failing tests**

```go
func TestOrderServiceCancelRollsBackInventory(t *testing.T) {}
func TestOrderServiceShipCreatesShipmentAndLog(t *testing.T) {}
func TestOrderServiceReceiveMarksDelivered(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/service -run TestOrderServiceCancelRollsBackInventory`
Expected: FAIL

**Step 3: Implement minimal lifecycle methods**

```go
func (s *OrderService) CancelOrder(id int64, reason string) error {}
func (s *OrderService) ShipOrder(id int64, carrier, tracking string) (*domain.Shipment, error) {}
func (s *OrderService) ReceiveOrder(id int64) error {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run OrderService.*Lifecycle`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/repository.go internal/repository/shipment_repository.go internal/repository/order_status_log_repository.go internal/service/order_service.go internal/service/order_service_test.go
git commit -m "feat(order): add lifecycle service methods"
```

### Task 3: Add admin endpoints

**Files:**
- Modify: `internal/api/admin_order_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/admin_order_handler_test.go`

**Step 1: Write failing tests**

```go
func TestAdminCancelOrder(t *testing.T) {}
func TestAdminShipOrder(t *testing.T) {}
func TestAdminReceiveOrder(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run Admin.*Order`
Expected: FAIL

**Step 3: Implement handlers + routes**

```go
POST /api/v1/admin/orders/:id/cancel
POST /api/v1/admin/orders/:id/ship
POST /api/v1/admin/orders/:id/receive
```

**Step 4: Run tests**

Run: `go test ./internal/api -run Admin.*Order`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/admin_order_handler.go internal/api/admin_order_handler_test.go cmd/api/main.go
git commit -m "feat(api): add order lifecycle admin endpoints"
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/service`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit remaining changes**

```bash
git add -A
git commit -m "chore: finalize order lifecycle"
```
