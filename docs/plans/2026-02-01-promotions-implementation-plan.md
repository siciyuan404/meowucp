# Promotions and Coupons Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add coupon validation and promotion rules for checkout pricing.

**Architecture:** Introduce coupon and promotion tables, apply rule evaluation in checkout totals, and expose validation APIs.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add coupon and promotion models + migration

**Files:**
- Create: `migrations/013_promotions.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS coupons (
  id BIGSERIAL PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL,
  value NUMERIC(12,2) NOT NULL,
  min_spend NUMERIC(12,2) NOT NULL DEFAULT 0,
  usage_limit INT NOT NULL DEFAULT 0,
  used_count INT NOT NULL DEFAULT 0,
  starts_at TIMESTAMPTZ,
  ends_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS promotions (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  rules JSONB NOT NULL,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  status TEXT NOT NULL DEFAULT 'active'
);
```

**Step 2: Add domain models**

```go
type Coupon struct {
  ID         int64 `gorm:"primary_key"`
  Code       string
  Type       string
  Value      float64
  MinSpend   float64
  UsageLimit int
  UsedCount  int
  StartsAt   *time.Time
  EndsAt     *time.Time
}

type Promotion struct {
  ID       int64 `gorm:"primary_key"`
  Name     string
  Rules    string
  StartsAt time.Time
  EndsAt   time.Time
  Status   string
}
```

**Step 3: Commit**

```bash
git add migrations/013_promotions.sql internal/domain/models.go
git commit -m "feat(promo): add coupon and promotion models"
```

### Task 2: Add rule evaluation service

**Files:**
- Create: `internal/service/promotion_service.go`
- Test: `internal/service/promotion_service_test.go`
- Modify: `internal/service/checkout_session_service.go`

**Step 1: Write failing tests**

```go
func TestCouponValidation(t *testing.T) {}
func TestPromotionAppliesToTotals(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/service -run Coupon|Promotion`
Expected: FAIL

**Step 3: Implement service**

```go
func (s *PromotionService) ValidateCoupon(code string, subtotal float64) (*domain.Coupon, error) {}
func (s *PromotionService) ApplyPromotions(subtotal float64, promotions []domain.Promotion) (float64, error) {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run Coupon|Promotion`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/promotion_service.go internal/service/promotion_service_test.go internal/service/checkout_session_service.go
git commit -m "feat(promo): add promotion evaluation service"
```

### Task 3: Add validation APIs

**Files:**
- Create: `internal/api/coupon_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/coupon_handler_test.go`

**Step 1: Write failing test**

```go
func TestCouponValidateEndpoint(t *testing.T) {}
```

**Step 2: Run test**

Run: `go test ./internal/api -run CouponValidate`
Expected: FAIL

**Step 3: Implement handler + route**

```go
POST /api/v1/coupons/validate
```

**Step 4: Run tests**

Run: `go test ./internal/api -run CouponValidate`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/coupon_handler.go internal/api/coupon_handler_test.go cmd/api/main.go
git commit -m "feat(api): add coupon validation endpoint"
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
git commit -m "chore: finalize promotions"
```
