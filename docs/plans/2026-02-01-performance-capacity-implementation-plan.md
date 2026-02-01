# Performance and Capacity Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Improve latency and capacity with caching, indexing, and async processing.

**Architecture:** Add targeted DB indexes, cache hot read endpoints, and move non-critical work to async worker where possible.

**Tech Stack:** Go 1.21, Gin, GORM, Redis, PostgreSQL

---

### Task 1: Add critical indexes

**Files:**
- Create: `migrations/016_performance_indexes.sql`

**Step 1: Write migration**

```sql
CREATE INDEX IF NOT EXISTS orders_user_id_created_at_idx
  ON orders (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS payments_order_id_idx
  ON payments (order_id);

CREATE INDEX IF NOT EXISTS products_sku_idx
  ON products (sku);
```

**Step 2: Commit**

```bash
git add migrations/016_performance_indexes.sql
git commit -m "perf(db): add critical indexes"
```

### Task 2: Add product list caching

**Files:**
- Modify: `internal/service/product_service.go`
- Test: `internal/service/product_service_test.go`

**Step 1: Write failing test**

```go
func TestProductListCachesResults(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/service -run ProductListCachesResults`
Expected: FAIL

**Step 3: Implement cache layer**

```go
// key: products:list:<filters>
```

**Step 4: Run tests**

Run: `go test ./internal/service -run ProductListCachesResults`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/product_service.go internal/service/product_service_test.go
git commit -m "perf(cache): add product list caching"
```

### Task 3: Async background updates

**Files:**
- Modify: `internal/ucp/worker` or `cmd/worker`
- Modify: `internal/service/order_service.go`
- Test: `internal/ucp/worker/*_test.go`

**Step 1: Write failing test**

```go
func TestAsyncOrderFollowupEnqueued(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/ucp/worker -run AsyncOrderFollowupEnqueued`
Expected: FAIL

**Step 3: Implement async enqueue**

```go
// enqueue non-critical work (email, analytics)
```

**Step 4: Run tests**

Run: `go test ./internal/ucp/worker -run AsyncOrderFollowupEnqueued`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/ucp/worker internal/service/order_service.go
git commit -m "perf(async): enqueue non-critical order followups"
```

### Task 4: Capacity test harness

**Files:**
- Create: `cmd/loadtest/main.go`
- Create: `cmd/loadtest/README.md`

**Step 1: Add load test harness**

```go
// minimal concurrency runner for product list and checkout creation
```

**Step 2: Document usage**

```bash
go run cmd/loadtest/main.go --endpoint http://localhost:8080
```

**Step 3: Commit**

```bash
git add cmd/loadtest/main.go cmd/loadtest/README.md
git commit -m "perf(test): add load test harness"
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
git commit -m "chore: finalize performance and capacity"
```
