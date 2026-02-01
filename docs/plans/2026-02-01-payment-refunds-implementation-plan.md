# Payment Refunds and Chargebacks Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add refund and chargeback support with order state reconciliation and webhook events.

**Architecture:** Introduce refund records and payment event history, then wire service logic to update order/payment states and emit webhooks. Keep handlers thin and use service methods.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add refund and payment event models + migration

**Files:**
- Create: `migrations/007_payment_refunds.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS payment_refunds (
  id BIGSERIAL PRIMARY KEY,
  payment_id BIGINT NOT NULL,
  amount NUMERIC(12,2) NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  reason TEXT,
  external_ref TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS payment_refunds_payment_id_idx
  ON payment_refunds (payment_id);

CREATE TABLE IF NOT EXISTS payment_events (
  id BIGSERIAL PRIMARY KEY,
  payment_id BIGINT NOT NULL,
  event_type TEXT NOT NULL,
  payload JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS payment_events_payment_id_idx
  ON payment_events (payment_id);
```

**Step 2: Add domain models**

```go
type PaymentRefund struct {
  ID          int64 `gorm:"primary_key"`
  PaymentID   int64 `gorm:"not null"`
  Amount      float64
  Status      string
  Reason      string
  ExternalRef *string
  CreatedAt   time.Time
  UpdatedAt   time.Time
}

type PaymentEvent struct {
  ID        int64 `gorm:"primary_key"`
  PaymentID int64 `gorm:"not null"`
  EventType string
  Payload   *string
  CreatedAt time.Time
}
```

**Step 3: Commit**

```bash
git add migrations/007_payment_refunds.sql internal/domain/models.go
git commit -m "feat(payment): add refund and payment event models"
```

### Task 2: Add repositories and service APIs

**Files:**
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/payment_refund_repository.go`
- Create: `internal/repository/payment_event_repository.go`
- Modify: `internal/service/payment_service.go`
- Test: `internal/service/payment_service_test.go`

**Step 1: Add failing test**

```go
func TestPaymentServiceCreateRefundUpdatesPaymentAndOrder(t *testing.T) {
  // set up payment + order
  // call CreateRefund
  // expect payment status refunded/partial and order status updated
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service -run TestPaymentServiceCreateRefundUpdatesPaymentAndOrder`
Expected: FAIL (method missing)

**Step 3: Implement minimal service**

```go
func (s *PaymentService) CreateRefund(paymentID int64, amount float64, reason string) (*domain.PaymentRefund, error) {
  // load payment + order, create refund record, update payment/order status, enqueue webhook
}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run TestPaymentServiceCreateRefundUpdatesPaymentAndOrder`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/repository.go internal/repository/payment_refund_repository.go internal/repository/payment_event_repository.go internal/service/payment_service.go internal/service/payment_service_test.go
git commit -m "feat(payment): add refund service and repositories"
```

### Task 3: Add refund API + webhook events

**Files:**
- Create: `internal/api/payment_refund_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/payment_refund_handler_test.go`

**Step 1: Write failing handler test**

```go
func TestPaymentRefundHandlerCreatesRefund(t *testing.T) {
  // POST /api/v1/payments/:id/refund
  // expect 200 and refund payload
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/api -run TestPaymentRefundHandlerCreatesRefund`
Expected: FAIL

**Step 3: Implement handler + route**

```go
type PaymentRefundRequest struct {
  Amount float64 `json:"amount"`
  Reason string  `json:"reason"`
}
```

**Step 4: Run tests**

Run: `go test ./internal/api -run TestPaymentRefundHandlerCreatesRefund`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/payment_refund_handler.go internal/api/payment_refund_handler_test.go cmd/api/main.go
git commit -m "feat(api): add payment refund endpoint"
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
git commit -m "chore: finalize refund flow"
```
