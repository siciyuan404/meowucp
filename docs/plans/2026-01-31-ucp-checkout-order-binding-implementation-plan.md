# UCP Checkout to Order Binding Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Bind UCP Checkout Complete to real order creation with inventory adjustment and idempotency.

**Architecture:** On Complete Checkout, map checkout session line_items to cart/order flow, enforce inventory checks, and create order/paid status in one transaction. Reuse existing OrderService idempotency and inventory logic.

**Tech Stack:** Go, Gin, GORM, UCP checkout models, OrderService

---

### Task 1: Add failing test for checkout completion creates order

**Files:**
- Modify: `internal/ucp/api/checkout_handler_test.go`

**Step 1: Write the failing test**

```go
func TestCheckoutCompleteCreatesOrder(t *testing.T) {
  // create checkout, then complete
  // expect order confirmation present in response
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ucp/api -run TestCheckoutCompleteCreatesOrder`
Expected: FAIL

**Step 3: Commit**

```bash
```

### Task 2: Map checkout line_items to order creation

**Files:**
- Modify: `internal/ucp/api/checkout_handler.go`
- Modify: `internal/service/order_service.go`
- Modify: `internal/ucp/model/checkout.go`

**Step 1: Implement mapping helper**

```go
func buildOrderFromCheckout(session *domain.CheckoutSession, payment model.PaymentInstrument) (*domain.Order, error) {
  // map buyer + line_items to order items and totals
}
```

**Step 2: Use OrderService.CreateOrder**

Call CreateOrder with an idempotency key derived from checkout id, and mark order paid if payment instrument present.

**Step 3: Run tests**

Run: `go test ./internal/ucp/api -run TestCheckoutCompleteCreatesOrder`
Expected: PASS

**Step 4: Commit**

```bash
```

### Task 3: Inventory and idempotency integration

**Files:**
- Modify: `internal/service/order_service.go`
- Modify: `internal/ucp/api/checkout_handler.go`
- Test: `internal/ucp/api/checkout_handler_test.go`

**Step 1: Write failing test for duplicate complete**

```go
func TestCheckoutCompleteIdempotent(t *testing.T) {
  // complete checkout twice, ensure same order returned
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./internal/ucp/api -run TestCheckoutCompleteIdempotent`
Expected: FAIL

**Step 3: Implement idempotency + inventory usage**

Use idempotency key based on checkout id and route through existing order inventory logic.

**Step 4: Run tests**

Run: `go test ./internal/ucp/api -run TestCheckoutCompleteIdempotent`
Expected: PASS

**Step 5: Commit**

```bash
```

### Task 4: End-to-end verification

**Step 1: Run tests**

Run: `go test ./internal/ucp/api`
Expected: PASS

Run: `go test ./...`
Expected: PASS

**Step 2: Commit any remaining changes**

```bash
```
