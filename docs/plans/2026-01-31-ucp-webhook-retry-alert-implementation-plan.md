# UCP Webhook Retry and Alerting Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Improve webhook retry behavior and alerting thresholds for outbound delivery stability.

**Architecture:** Extend webhook queue worker to support exponential backoff and configurable retry ceilings. Update alerting to include webhook event context.

**Tech Stack:** Go, Redis Stream, existing webhook queue/worker

---

### Task 1: Add failing test for retry backoff

**Files:**
- Modify: `internal/ucp/worker/webhook_worker_test.go`

**Step 1: Write the failing test**

```go
func TestWebhookRetryBackoff(t *testing.T) {
  // simulate failed delivery, expect next_retry_at to increase exponentially
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ucp/worker -run TestWebhookRetryBackoff`
Expected: FAIL

**Step 3: Commit**

```bash
```

### Task 2: Implement backoff logic

**Files:**
- Modify: `internal/ucp/worker/webhook_worker.go`

**Step 1: Implement minimal backoff**

```go
func nextRetry(attempts int, base time.Duration) time.Time {
  // exponential backoff capped to max
}
```

**Step 2: Run test**

Run: `go test ./internal/ucp/worker -run TestWebhookRetryBackoff`
Expected: PASS

**Step 3: Commit**

```bash
```

### Task 3: Alert context enhancement

**Files:**
- Modify: `internal/service/webhook_alert_service.go`
- Modify: `internal/api/ucp_webhook_alert_handler.go`
- Test: `internal/service/webhook_alert_service_test.go`

**Step 1: Write failing test**

```go
func TestWebhookAlertIncludesEventContext(t *testing.T) {
  // expect event_id/order_id in alert details
}
```

**Step 2: Run test**

Run: `go test ./internal/service -run TestWebhookAlertIncludesEventContext`
Expected: FAIL

**Step 3: Implement context**

Add event_id/order_id to alert details payload.

**Step 4: Run test**

Run: `go test ./internal/service -run TestWebhookAlertIncludesEventContext`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/ucp/worker`
Expected: PASS

Run: `go test ./internal/service`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
```
