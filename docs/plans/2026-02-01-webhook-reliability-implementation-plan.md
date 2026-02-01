# Webhook Reliability Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add DLQ, replay tooling, and alert deduplication for webhook delivery.

**Architecture:** Extend queue persistence with DLQ table, add admin APIs for listing and replaying jobs, and implement alert aggregation with rate limits.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add DLQ and replay log models + migration

**Files:**
- Create: `migrations/010_webhook_dlq.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS webhook_dlq (
  id BIGSERIAL PRIMARY KEY,
  job_id BIGINT NOT NULL,
  reason TEXT NOT NULL,
  payload TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_replay_logs (
  id BIGSERIAL PRIMARY KEY,
  job_id BIGINT NOT NULL,
  replay_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  result TEXT NOT NULL
);
```

**Step 2: Add domain models**

```go
type WebhookDLQ struct {
  ID        int64 `gorm:"primary_key"`
  JobID     int64
  Reason    string
  Payload   string
  CreatedAt time.Time
}

type WebhookReplayLog struct {
  ID       int64 `gorm:"primary_key"`
  JobID    int64
  ReplayAt time.Time
  Result   string
}
```

**Step 3: Commit**

```bash
git add migrations/010_webhook_dlq.sql internal/domain/models.go
git commit -m "feat(webhook): add DLQ and replay models"
```

### Task 2: Add repository + service methods

**Files:**
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/webhook_dlq_repository.go`
- Create: `internal/repository/webhook_replay_log_repository.go`
- Modify: `internal/service/webhook_queue_service.go`
- Test: `internal/service/webhook_queue_service_test.go`

**Step 1: Write failing tests**

```go
func TestWebhookQueueMovesToDLQAfterMaxAttempts(t *testing.T) {}
func TestWebhookReplayLogsResult(t *testing.T) {}
```

**Step 2: Run tests to verify failure**

Run: `go test ./internal/service -run Webhook.*DLQ`
Expected: FAIL

**Step 3: Implement DLQ + replay logging**

```go
func (s *WebhookQueueService) MoveToDLQ(job *domain.UCPWebhookJob, reason string) error {}
func (s *WebhookQueueService) ReplayJob(jobID int64) error {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run Webhook.*DLQ`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/repository.go internal/repository/webhook_dlq_repository.go internal/repository/webhook_replay_log_repository.go internal/service/webhook_queue_service.go internal/service/webhook_queue_service_test.go
git commit -m "feat(webhook): add DLQ and replay support"
```

### Task 3: Add admin APIs for DLQ and replay

**Files:**
- Create: `internal/api/admin_webhook_dlq_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/admin_webhook_dlq_handler_test.go`

**Step 1: Write failing tests**

```go
func TestAdminListWebhookDLQ(t *testing.T) {}
func TestAdminReplayWebhookDLQ(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run AdminWebhookDLQ`
Expected: FAIL

**Step 3: Implement handler + routes**

```go
GET /api/v1/admin/webhooks/dlq
POST /api/v1/admin/webhooks/dlq/:id/replay
```

**Step 4: Run tests**

Run: `go test ./internal/api -run AdminWebhookDLQ`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/admin_webhook_dlq_handler.go internal/api/admin_webhook_dlq_handler_test.go cmd/api/main.go
git commit -m "feat(api): add webhook DLQ admin endpoints"
```

### Task 4: Alert deduplication

**Files:**
- Modify: `internal/service/webhook_alert_service.go`
- Test: `internal/service/webhook_alert_service_test.go`

**Step 1: Write failing test**

```go
func TestWebhookAlertDedupWithinWindow(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/service -run WebhookAlertDedup`
Expected: FAIL

**Step 3: Implement dedup window**

```go
func (s *WebhookAlertService) CreateDedup(alert *domain.UCPWebhookAlert, window time.Duration) error {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run WebhookAlertDedup`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/webhook_alert_service.go internal/service/webhook_alert_service_test.go
git commit -m "feat(webhook): add alert deduplication"
```

### Task 5: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/service`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit remaining changes**

```bash
git add -A
git commit -m "chore: finalize webhook reliability"
```
