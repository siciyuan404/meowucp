# Observability Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add structured logging, audit logs, and basic metrics exposure.

**Architecture:** Introduce an audit log table, middleware for logging, and optional /metrics endpoint guarded for internal use.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add audit log model + migration

**Files:**
- Create: `migrations/014_audit_logs.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGSERIAL PRIMARY KEY,
  actor TEXT NOT NULL,
  action TEXT NOT NULL,
  target TEXT NOT NULL,
  payload JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Step 2: Add domain model**

```go
type AuditLog struct {
  ID        int64 `gorm:"primary_key"`
  Actor     string
  Action    string
  Target    string
  Payload   *string
  CreatedAt time.Time
}
```

**Step 3: Commit**

```bash
git add migrations/014_audit_logs.sql internal/domain/models.go
git commit -m "feat(obs): add audit log model"
```

### Task 2: Add logging middleware

**Files:**
- Create: `internal/middleware/request_logger.go`
- Modify: `cmd/api/main.go`
- Test: `internal/middleware/request_logger_test.go`

**Step 1: Write failing test**

```go
func TestRequestLoggerIncludesRequestID(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/middleware -run RequestLogger`
Expected: FAIL

**Step 3: Implement middleware**

```go
// log method, path, status, duration, request_id
```

**Step 4: Run tests**

Run: `go test ./internal/middleware -run RequestLogger`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/middleware/request_logger.go internal/middleware/request_logger_test.go cmd/api/main.go
git commit -m "feat(obs): add structured request logging"
```

### Task 3: Add audit logging in admin actions

**Files:**
- Modify: `internal/service/*` (admin actions)
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/audit_log_repository.go`
- Test: `internal/service/*_test.go`

**Step 1: Write failing test**

```go
func TestAdminOrderShipWritesAuditLog(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/service -run AuditLog`
Expected: FAIL

**Step 3: Implement audit log writes**

```go
func (s *AuditLogService) Record(actor, action, target string, payload string) error {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run AuditLog`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/repository.go internal/repository/audit_log_repository.go internal/service internal/service/*_test.go
git commit -m "feat(obs): add audit logging for admin actions"
```

### Task 4: Metrics endpoint

**Files:**
- Create: `internal/api/metrics_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/metrics_handler_test.go`

**Step 1: Write failing test**

```go
func TestMetricsEndpointRequiresInternalAuth(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/api -run MetricsEndpoint`
Expected: FAIL

**Step 3: Implement handler**

```go
GET /metrics
```

**Step 4: Run tests**

Run: `go test ./internal/api -run MetricsEndpoint`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/metrics_handler.go internal/api/metrics_handler_test.go cmd/api/main.go
git commit -m "feat(obs): add metrics endpoint"
```

### Task 5: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/middleware`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit remaining changes**

```bash
git add -A
git commit -m "chore: finalize observability"
```
