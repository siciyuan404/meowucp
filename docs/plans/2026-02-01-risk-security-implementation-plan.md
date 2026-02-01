# Risk and Security Hardening Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Enforce idempotency across write APIs, add replay protection, and unify signature validation.

**Architecture:** Introduce an idempotency record store, middleware for request fingerprinting, and shared signature verification helpers used by callbacks and webhook receivers.

**Tech Stack:** Go 1.21, Gin, GORM, Redis (optional), PostgreSQL

---

### Task 1: Add idempotency record model + migration

**Files:**
- Create: `migrations/009_idempotency_keys.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS idempotency_keys (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  key TEXT NOT NULL,
  request_hash TEXT NOT NULL,
  response_snapshot JSONB,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idempotency_keys_user_key_uidx
  ON idempotency_keys (user_id, key);
```

**Step 2: Add domain model**

```go
type IdempotencyKey struct {
  ID               int64 `gorm:"primary_key"`
  UserID           int64 `gorm:"not null"`
  Key              string
  RequestHash      string
  ResponseSnapshot *string
  Status           string
  CreatedAt        time.Time
  UpdatedAt        time.Time
}
```

**Step 3: Commit**

```bash
git add migrations/009_idempotency_keys.sql internal/domain/models.go
git commit -m "feat(security): add idempotency key model"
```

### Task 2: Add middleware and repository

**Files:**
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/idempotency_repository.go`
- Create: `internal/middleware/idempotency.go`
- Test: `internal/middleware/idempotency_test.go`

**Step 1: Write failing test**

```go
func TestIdempotencyMiddlewareReplaysResponse(t *testing.T) {
  // same key + same hash returns cached response
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/middleware -run TestIdempotencyMiddlewareReplaysResponse`
Expected: FAIL

**Step 3: Implement middleware**

```go
// Extract Idempotency-Key header, hash body, store/return response snapshot
```

**Step 4: Run tests**

Run: `go test ./internal/middleware -run TestIdempotencyMiddlewareReplaysResponse`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/repository.go internal/repository/idempotency_repository.go internal/middleware/idempotency.go internal/middleware/idempotency_test.go
git commit -m "feat(security): add idempotency middleware"
```

### Task 3: Signature and replay protection

**Files:**
- Modify: `internal/ucp/security` (or `internal/ucp/api` helpers)
- Modify: webhook/callback handlers
- Test: `internal/ucp/security/*_test.go`

**Step 1: Write failing tests**

```go
func TestWebhookRejectsReplayNonce(t *testing.T) {}
func TestCallbackRejectsInvalidSignature(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/ucp/security -run Replay|Signature`
Expected: FAIL

**Step 3: Implement nonce store + verifier**

```go
type NonceStore interface { Seen(nonce string) bool; Mark(nonce string) error }
```

**Step 4: Run tests**

Run: `go test ./internal/ucp/security -run Replay|Signature`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/ucp/security internal/ucp/api internal/api
git commit -m "feat(security): enforce signature and replay protection"
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/middleware`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit remaining changes**

```bash
git add -A
git commit -m "chore: finalize security hardening"
```
