# Payment Callback and Order Sync Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Sync payment completion callbacks into order status updates and trigger paid outbound events.

**Architecture:** Add a payment callback handler that verifies payloads, updates Payment + Order status to paid, and emits outbound order events via the webhook queue.

**Tech Stack:** Go, Gin, GORM, existing PaymentService/OrderService/WebhookQueueService

---

### Task 1: Add failing test for payment callback updates order

**Files:**
- Create: `internal/api/payment_callback_handler_test.go`

**Step 1: Write the failing test**

```go
func TestPaymentCallbackMarksOrderPaid(t *testing.T) {
  // POST /api/v1/payment/callback
  // expect OrderService.UpdateOrderStatus called with paid
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/api -run TestPaymentCallbackMarksOrderPaid`
Expected: FAIL

**Step 3: Commit**

```bash
```

### Task 2: Implement callback handler

**Files:**
- Create: `internal/api/payment_callback_handler.go`
- Modify: `cmd/api/main.go`

**Step 1: Implement minimal handler**

```go
func (h *PaymentCallbackHandler) Handle(c *gin.Context) {
  // parse payload, validate order_id/transaction_id
  // update payment + order to paid
}
```

**Step 2: Run test to verify it passes**

Run: `go test ./internal/api -run TestPaymentCallbackMarksOrderPaid`
Expected: PASS

**Step 3: Commit**

```bash
```

### Task 3: Trigger outbound paid event

**Files:**
- Modify: `internal/service/order_service.go`
- Modify: `internal/service/webhook_queue_service.go`
- Test: `internal/service/order_service_test.go`

**Step 1: Write failing test**

```go
func TestPaymentCallbackTriggersPaidWebhook(t *testing.T) {
  // after payment callback, expect enqueue paid event
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/service -run TestPaymentCallbackTriggersPaidWebhook`
Expected: FAIL

**Step 3: Implement**

Hook into paid status update to enqueue `paid` event.

**Step 4: Run test**

Run: `go test ./internal/service -run TestPaymentCallbackTriggersPaidWebhook`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/api`
Expected: PASS

Run: `go test ./internal/service`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
```
