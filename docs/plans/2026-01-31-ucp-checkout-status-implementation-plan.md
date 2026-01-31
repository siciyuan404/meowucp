# UCP Checkout Status Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement deterministic checkout status and error messaging rules for UCP Checkout based on `requires_buyer_input` and `recoverable` messages.

**Architecture:** Extend `internal/ucp/api/checkout_handler.go` with a status/message resolver used by Create and Update. Tests live in `internal/ucp/api/checkout_handler_test.go` using existing fake repositories.

**Tech Stack:** Go, Gin, GORM, UCP model types

---

### Task 1: Add failing tests for requires_escalation rules

**Files:**
- Modify: `internal/ucp/api/checkout_handler_test.go`

**Step 1: Write the failing test**

```go
func TestCheckoutCreateRequiresEscalationWhenNoHandlers(t *testing.T) {
  // configure services with no payment handlers
  // POST /ucp/v1/checkout-sessions
  // expect status=requires_escalation, continue_url present
  // expect messages include code payment_handlers_missing with requires_buyer_input
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ucp/api -run TestCheckoutCreateRequiresEscalationWhenNoHandlers`
Expected: FAIL with status not requires_escalation or missing message

**Step 3: Write the failing test for recoverable**

```go
func TestCheckoutCreateIncompleteOnRecoverableErrors(t *testing.T) {
  // missing currency or line_items
  // expect status=incomplete and recoverable message
}
```

**Step 4: Run test to verify it fails**

Run: `go test ./internal/ucp/api -run TestCheckoutCreateIncompleteOnRecoverableErrors`
Expected: FAIL

**Step 5: Commit**

```bash
git add internal/ucp/api/checkout_handler_test.go
git commit -m "test(ucp): add checkout status rule coverage"
```

### Task 2: Implement status/message resolver

**Files:**
- Modify: `internal/ucp/api/checkout_handler.go`

**Step 1: Implement minimal resolver**

```go
func resolveMessagesAndStatus(hasHandlers bool, recoverable []model.Message, buyerInput []model.Message) (string, []model.Message) {
  // if !hasHandlers -> add payment_handlers_missing requires_buyer_input
  // combine messages, compute status via rules
}
```

**Step 2: Use resolver in Create/Update**

Replace existing `resolveCheckoutStatus` calls with resolver output.

**Step 3: Run tests**

Run: `go test ./internal/ucp/api -run CheckoutCreate`
Expected: PASS

**Step 4: Commit**

```bash
git add internal/ucp/api/checkout_handler.go
git commit -m "feat(ucp): enforce checkout escalation rules"
```

### Task 3: Update flow and messages for Update endpoint

**Files:**
- Modify: `internal/ucp/api/checkout_handler_test.go`
- Modify: `internal/ucp/api/checkout_handler.go`

**Step 1: Write failing test for Update**

```go
func TestCheckoutUpdateRequiresEscalationWhenSignInNeeded(t *testing.T) {
  // update path returns requires_sign_in with requires_buyer_input
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ucp/api -run TestCheckoutUpdateRequiresEscalationWhenSignInNeeded`
Expected: FAIL

**Step 3: Implement minimal update behavior**

Add rules to create buyer-input messages when required (e.g., missing buyer data or custom flag from request).

**Step 4: Run tests**

Run: `go test ./internal/ucp/api -run CheckoutUpdate`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/ucp/api/checkout_handler.go internal/ucp/api/checkout_handler_test.go
git commit -m "feat(ucp): apply escalation rules on checkout update"
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/ucp/api`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
git add -A
git commit -m "chore: finalize checkout status rules"
```
