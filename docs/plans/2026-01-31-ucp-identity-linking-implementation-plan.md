# UCP Identity Linking Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add OAuth 2.0 endpoints and metadata to support UCP Identity Linking capability for platform authorization.

**Architecture:** Implement OAuth metadata endpoints and token issuance with minimal scopes for checkout. Use existing user model and JWT where possible, but keep OAuth tokens distinct.

**Tech Stack:** Go, Gin, OAuth 2.0 concepts, JWT

---

### Task 1: Add OAuth metadata endpoints

**Files:**
- Create: `internal/api/oauth_metadata_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/oauth_metadata_handler_test.go`

**Step 1: Write the failing test**

```go
func TestOAuthMetadataWellKnown(t *testing.T) {
  // GET /.well-known/oauth-authorization-server
  // expect issuer, token_endpoint, authorization_endpoint, scopes_supported
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/api -run TestOAuthMetadataWellKnown`
Expected: FAIL

**Step 3: Implement metadata handler**

```go
func (h *OAuthMetadataHandler) WellKnown(c *gin.Context) {
  // return RFC 8414 response
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/api -run TestOAuthMetadataWellKnown`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 2: Token endpoint (authorization code exchange)

**Files:**
- Create: `internal/api/oauth_token_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/oauth_token_handler_test.go`

**Step 1: Write the failing test**

```go
func TestOAuthTokenExchange(t *testing.T) {
  // POST /oauth2/token with client credentials
  // expect access_token, token_type, expires_in
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/api -run TestOAuthTokenExchange`
Expected: FAIL

**Step 3: Implement token exchange**

```go
func (h *OAuthTokenHandler) Token(c *gin.Context) {
  // validate client_id/client_secret and code
  // issue access_token for ucp:scopes:checkout_session
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/api -run TestOAuthTokenExchange`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 3: Authorization endpoint (minimal)

**Files:**
- Create: `internal/api/oauth_authorize_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/oauth_authorize_handler_test.go`

**Step 1: Write the failing test**

```go
func TestOAuthAuthorizeRedirect(t *testing.T) {
  // GET /oauth2/authorize?response_type=code&client_id=...&redirect_uri=...
  // expect redirect with code
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/api -run TestOAuthAuthorizeRedirect`
Expected: FAIL

**Step 3: Implement authorize**

```go
func (h *OAuthAuthorizeHandler) Authorize(c *gin.Context) {
  // validate client and generate code
  // redirect to redirect_uri
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./internal/api -run TestOAuthAuthorizeRedirect`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/api`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
```
