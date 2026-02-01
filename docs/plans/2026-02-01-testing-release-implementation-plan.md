# Testing and Release Readiness Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add end-to-end regression coverage and release tooling.

**Architecture:** Implement high-value E2E tests for checkout/payment/order flow and provide release metadata tracking.

**Tech Stack:** Go 1.21, testing, Gin, GORM

---

### Task 1: Add release history model + migration

**Files:**
- Create: `migrations/018_release_history.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS release_history (
  id BIGSERIAL PRIMARY KEY,
  version TEXT NOT NULL,
  deployed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  operator TEXT NOT NULL
);
```

**Step 2: Add domain model**

```go
type ReleaseHistory struct {
  ID         int64 `gorm:"primary_key"`
  Version    string
  DeployedAt time.Time
  Operator   string
}
```

**Step 3: Commit**

```bash
git add migrations/018_release_history.sql internal/domain/models.go
git commit -m "feat(release): add release history model"
```

### Task 2: Add E2E regression tests

**Files:**
- Create: `internal/api/e2e_checkout_flow_test.go`

**Step 1: Write failing test**

```go
func TestE2ECheckoutPaymentOrder(t *testing.T) {
  // create checkout -> complete -> payment callback -> order paid
}
```

**Step 2: Run test**

Run: `go test ./internal/api -run TestE2ECheckoutPaymentOrder`
Expected: FAIL

**Step 3: Implement test using existing fakes**

```go
// use in-memory repos and handlers to simulate flow
```

**Step 4: Run test**

Run: `go test ./internal/api -run TestE2ECheckoutPaymentOrder`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/e2e_checkout_flow_test.go
git commit -m "test(e2e): add checkout-payment-order flow"
```

### Task 3: Add release notes endpoint

**Files:**
- Create: `internal/api/admin_release_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/admin_release_handler_test.go`

**Step 1: Write failing tests**

```go
func TestAdminReleaseHistoryList(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run ReleaseHistory`
Expected: FAIL

**Step 3: Implement handler + route**

```go
GET /api/v1/admin/releases
```

**Step 4: Run tests**

Run: `go test ./internal/api -run ReleaseHistory`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/admin_release_handler.go internal/api/admin_release_handler_test.go cmd/api/main.go
git commit -m "feat(api): add release history endpoint"
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/api`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit remaining changes**

```bash
git add -A
git commit -m "chore: finalize testing and release readiness"
```
