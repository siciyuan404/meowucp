# OAuth Client Management Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Support multiple OAuth clients, token revocation, and scoped access control.

**Architecture:** Persist clients and tokens in DB, expose admin CRUD, and enforce scope checks in OAuth handlers.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add OAuth client/token models + migration

**Files:**
- Create: `migrations/011_oauth_clients.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS oauth_clients (
  id BIGSERIAL PRIMARY KEY,
  client_id TEXT NOT NULL UNIQUE,
  secret_hash TEXT NOT NULL,
  scopes TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS oauth_tokens (
  id BIGSERIAL PRIMARY KEY,
  token TEXT NOT NULL UNIQUE,
  client_id TEXT NOT NULL,
  user_id BIGINT,
  scopes TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ
);
```

**Step 2: Add domain models**

```go
type OAuthClient struct {
  ID         int64 `gorm:"primary_key"`
  ClientID   string
  SecretHash string
  Scopes     string
  Status     string
  CreatedAt  time.Time
}

type OAuthToken struct {
  ID        int64 `gorm:"primary_key"`
  Token     string
  ClientID  string
  UserID    *int64
  Scopes    string
  ExpiresAt time.Time
  RevokedAt *time.Time
}
```

**Step 3: Commit**

```bash
git add migrations/011_oauth_clients.sql internal/domain/models.go
git commit -m "feat(oauth): add client and token models"
```

### Task 2: Add repositories + admin CRUD

**Files:**
- Modify: `internal/repository/repository.go`
- Create: `internal/repository/oauth_client_repository.go`
- Create: `internal/repository/oauth_token_repository.go`
- Create: `internal/api/admin_oauth_client_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/admin_oauth_client_handler_test.go`

**Step 1: Write failing tests**

```go
func TestAdminCreatesOAuthClient(t *testing.T) {}
func TestAdminListsOAuthClients(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run AdminOAuthClient`
Expected: FAIL

**Step 3: Implement handler + routes**

```go
POST /api/v1/admin/oauth/clients
GET /api/v1/admin/oauth/clients
```

**Step 4: Run tests**

Run: `go test ./internal/api -run AdminOAuthClient`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/repository/repository.go internal/repository/oauth_client_repository.go internal/repository/oauth_token_repository.go internal/api/admin_oauth_client_handler.go internal/api/admin_oauth_client_handler_test.go cmd/api/main.go
git commit -m "feat(api): add OAuth client admin endpoints"
```

### Task 3: Token revocation + scope enforcement

**Files:**
- Modify: `internal/api/oauth_token_handler.go`
- Create: `internal/api/oauth_revoke_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/oauth_token_handler_test.go`

**Step 1: Write failing tests**

```go
func TestOAuthTokenScopesRestricted(t *testing.T) {}
func TestOAuthRevokeToken(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run OAuth.*Token`
Expected: FAIL

**Step 3: Implement revoke endpoint**

```go
POST /oauth2/revoke
```

**Step 4: Run tests**

Run: `go test ./internal/api -run OAuth.*Token`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/oauth_token_handler.go internal/api/oauth_revoke_handler.go internal/api/oauth_token_handler_test.go cmd/api/main.go
git commit -m "feat(oauth): add token revocation and scope checks"
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/api -run OAuth`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit remaining changes**

```bash
git add -A
git commit -m "chore: finalize OAuth client management"
```
