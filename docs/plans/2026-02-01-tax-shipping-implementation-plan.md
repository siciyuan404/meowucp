# Tax and Shipping Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Provide tax calculation, shipping rates, and address validation for checkout.

**Architecture:** Add tax and shipping rule tables, a service for rate calculation, and API endpoints used by checkout.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add tax and shipping rule models + migration

**Files:**
- Create: `migrations/012_tax_shipping.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS tax_rules (
  id BIGSERIAL PRIMARY KEY,
  region TEXT NOT NULL,
  category TEXT NOT NULL,
  rate NUMERIC(6,4) NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS shipping_rules (
  id BIGSERIAL PRIMARY KEY,
  region TEXT NOT NULL,
  method TEXT NOT NULL,
  base_amount NUMERIC(12,2) NOT NULL,
  per_item_amount NUMERIC(12,2) NOT NULL
);
```

**Step 2: Add domain models**

```go
type TaxRule struct {
  ID          int64 `gorm:"primary_key"`
  Region      string
  Category    string
  Rate        float64
  EffectiveAt time.Time
}

type ShippingRule struct {
  ID            int64 `gorm:"primary_key"`
  Region        string
  Method        string
  BaseAmount    float64
  PerItemAmount float64
}
```

**Step 3: Commit**

```bash
git add migrations/012_tax_shipping.sql internal/domain/models.go
git commit -m "feat(tax): add tax and shipping rule models"
```

### Task 2: Add rate calculation service

**Files:**
- Modify: `internal/service/checkout_session_service.go`
- Create: `internal/service/tax_shipping_service.go`
- Test: `internal/service/tax_shipping_service_test.go`

**Step 1: Write failing tests**

```go
func TestTaxRateAppliedByRegion(t *testing.T) {}
func TestShippingRateByItems(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/service -run Tax|Shipping`
Expected: FAIL

**Step 3: Implement rate calculation**

```go
func (s *TaxShippingService) Quote(region string, items []domain.OrderItem) (tax float64, shipping float64, err error) {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run Tax|Shipping`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/checkout_session_service.go internal/service/tax_shipping_service.go internal/service/tax_shipping_service_test.go
git commit -m "feat(tax): add tax and shipping quote service"
```

### Task 3: Add API endpoints

**Files:**
- Create: `internal/api/shipping_rate_handler.go`
- Create: `internal/api/address_validation_handler.go`
- Modify: `cmd/api/main.go`
- Test: `internal/api/shipping_rate_handler_test.go`

**Step 1: Write failing tests**

```go
func TestShippingRateEndpoint(t *testing.T) {}
func TestAddressValidationEndpoint(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run ShippingRate|AddressValidation`
Expected: FAIL

**Step 3: Implement routes**

```go
GET /api/v1/shipping/rates
POST /api/v1/address/validate
```

**Step 4: Run tests**

Run: `go test ./internal/api -run ShippingRate|AddressValidation`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/shipping_rate_handler.go internal/api/address_validation_handler.go internal/api/shipping_rate_handler_test.go cmd/api/main.go
git commit -m "feat(api): add shipping rate and address validation"
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
git commit -m "chore: finalize tax and shipping"
```
