# UCP Order Webhook Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Emit outbound UCP order events (created/paid/shipped/cancelled) via async queue by default, with an admin sync trigger.

**Architecture:** Extend OrderService to enqueue order events after status changes. Reuse existing webhook queue/worker and delivery URL. Add an admin endpoint for synchronous delivery.

**Tech Stack:** Go, Gin, GORM, Redis Stream, existing webhook queue/worker

---

### Task 1: Add outbound event payload builder

**Files:**
- Modify: `internal/service/webhook_queue_service.go`
- Modify: `internal/ucp/model` (if new payload struct needed)
- Test: `internal/service/webhook_queue_service_test.go`

**Step 1: Write the failing test**

```go
func TestBuildOrderWebhookPayload(t *testing.T) {
  // build payload for event_type=paid
  // expect order_no, status, event_type, timestamp
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service -run TestBuildOrderWebhookPayload`
Expected: FAIL

**Step 3: Write minimal implementation**

```go
func buildOrderWebhookPayload(order *domain.Order, eventType string) ([]byte, error) {
  // map to JSON with order_no, status, event_type, created_at
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/service -run TestBuildOrderWebhookPayload`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/webhook_queue_service.go internal/service/webhook_queue_service_test.go
git commit -m "feat(ucp): add outbound order payload builder"
```

### Task 2: Enqueue order events on status change

**Files:**
- Modify: `internal/service/order_service.go`
- Modify: `internal/service/webhook_queue_service.go`
- Test: `internal/service/order_service_test.go`

**Step 1: Write the failing test**

```go
func TestOrderServiceEnqueuesOrderWebhookOnStatusChange(t *testing.T) {
  // update order status to paid
  // assert enqueue called with event_type=paid
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service -run TestOrderServiceEnqueuesOrderWebhookOnStatusChange`
Expected: FAIL

**Step 3: Implement minimal enqueue hook**

```go
// after successful status change
webhookQueue.EnqueueOrderEvent(order, "paid")
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/service -run TestOrderServiceEnqueuesOrderWebhookOnStatusChange`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 3: Add admin manual trigger (sync/async)

**Files:**
- Create: `internal/api/admin_order_webhook_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/admin_order_webhook_handler_test.go`

**Step 1: Write the failing test**

```go
func TestAdminOrderWebhookSync(t *testing.T) {
  // POST /api/v1/admin/orders/:id/webhook with event_type=paid mode=sync
  // expect DeliverOrderEvent called
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/api -run TestAdminOrderWebhookSync`
Expected: FAIL

**Step 3: Implement handler**

```go
type AdminOrderWebhookHandler struct {
  orderSvc OrderService
  webhookSvc WebhookQueueService
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/api -run TestAdminOrderWebhookSync`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/service`
Expected: PASS

Run: `go test ./internal/api`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
```
