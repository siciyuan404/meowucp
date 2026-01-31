# UCP Profile and Capabilities Enhancement Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Expand the UCP business profile response with accurate services, capabilities, and version metadata.

**Architecture:** Update profile builder in `internal/ucp/api/profile_handler.go` to include Order capability and any supported extensions. Keep URLs derived from config.

**Tech Stack:** Go, Gin, UCP model types

---

### Task 1: Add failing test for profile capabilities

**Files:**
- Modify: `internal/ucp/api/profile_handler_test.go`

**Step 1: Write the failing test**

```go
func TestProfileIncludesOrderCapability(t *testing.T) {
  // GET /.well-known/ucp
  // expect dev.ucp.shopping.order capability present
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ucp/api -run TestProfileIncludesOrderCapability`
Expected: FAIL

**Step 3: Commit**

```bash
```

### Task 2: Update profile builder

**Files:**
- Modify: `internal/ucp/api/profile_handler.go`

**Step 1: Implement capability expansion**

```go
// add dev.ucp.shopping.order capability with spec + schema URLs
```

**Step 2: Run test**

Run: `go test ./internal/ucp/api -run TestProfileIncludesOrderCapability`
Expected: PASS

**Step 3: Commit**

```bash
```

### Task 3: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/ucp/api`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
```
