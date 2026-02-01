# Data Protection and Compliance Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add auditability, data retention policies, and PII masking for compliance.

**Architecture:** Introduce retention policies and audit logs, then apply masking at response boundaries and admin views.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add retention policy model + migration

**Files:**
- Create: `migrations/017_data_retention.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS data_retention_policies (
  id BIGSERIAL PRIMARY KEY,
  entity TEXT NOT NULL,
  ttl_days INT NOT NULL,
  strategy TEXT NOT NULL
);
```

**Step 2: Add domain model**

```go
type DataRetentionPolicy struct {
  ID       int64 `gorm:"primary_key"`
  Entity   string
  TTLDays  int
  Strategy string
}
```

**Step 3: Commit**

```bash
git add migrations/017_data_retention.sql internal/domain/models.go
git commit -m "feat(compliance): add data retention policy model"
```

### Task 2: Add masking helpers

**Files:**
- Create: `internal/service/masking_service.go`
- Test: `internal/service/masking_service_test.go`
- Modify: `internal/api/user_handler.go`

**Step 1: Write failing test**

```go
func TestMaskingServiceMasksEmail(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/service -run MaskingService`
Expected: FAIL

**Step 3: Implement masking**

```go
func MaskEmail(email string) string {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run MaskingService`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/masking_service.go internal/service/masking_service_test.go internal/api/user_handler.go
git commit -m "feat(compliance): add PII masking"
```

### Task 3: Add admin audit log listing

**Files:**
- Create: `internal/api/admin_audit_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/admin_audit_handler_test.go`

**Step 1: Write failing test**

```go
func TestAdminAuditList(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/api -run AdminAuditList`
Expected: FAIL

**Step 3: Implement handler + route**

```go
GET /api/v1/admin/audit-logs
```

**Step 4: Run tests**

Run: `go test ./internal/api -run AdminAuditList`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/admin_audit_handler.go internal/api/admin_audit_handler_test.go cmd/api/main.go
git commit -m "feat(api): add audit log listing"
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
git commit -m "chore: finalize data compliance"
```
