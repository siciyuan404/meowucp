# Localization and Multi-Currency Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Support locale-aware responses and multi-currency price display.

**Architecture:** Store currency rates and i18n strings, then update catalog responses to respect locale/currency query params.

**Tech Stack:** Go 1.21, Gin, GORM, PostgreSQL

---

### Task 1: Add currency rate + i18n models and migration

**Files:**
- Create: `migrations/015_i18n_currency.sql`
- Modify: `internal/domain/models.go`

**Step 1: Write migration**

```sql
CREATE TABLE IF NOT EXISTS currency_rates (
  id BIGSERIAL PRIMARY KEY,
  base TEXT NOT NULL,
  target TEXT NOT NULL,
  rate NUMERIC(18,6) NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS i18n_strings (
  id BIGSERIAL PRIMARY KEY,
  key TEXT NOT NULL,
  locale TEXT NOT NULL,
  value TEXT NOT NULL
);
```

**Step 2: Add domain models**

```go
type CurrencyRate struct {
  ID        int64 `gorm:"primary_key"`
  Base      string
  Target    string
  Rate      float64
  UpdatedAt time.Time
}

type I18nString struct {
  ID     int64 `gorm:"primary_key"`
  Key    string
  Locale string
  Value  string
}
```

**Step 3: Commit**

```bash
git add migrations/015_i18n_currency.sql internal/domain/models.go
git commit -m "feat(i18n): add currency and locale models"
```

### Task 2: Add localization service

**Files:**
- Create: `internal/service/localization_service.go`
- Test: `internal/service/localization_service_test.go`

**Step 1: Write failing tests**

```go
func TestConvertCurrency(t *testing.T) {}
func TestTranslateString(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/service -run ConvertCurrency|TranslateString`
Expected: FAIL

**Step 3: Implement service**

```go
func (s *LocalizationService) Convert(amount float64, base, target string) (float64, error) {}
func (s *LocalizationService) Translate(key, locale string) (string, error) {}
```

**Step 4: Run tests**

Run: `go test ./internal/service -run ConvertCurrency|TranslateString`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/service/localization_service.go internal/service/localization_service_test.go
git commit -m "feat(i18n): add localization service"
```

### Task 3: Apply locale/currency to catalog APIs

**Files:**
- Modify: `internal/api/product_handler.go`
- Modify: `internal/api/category_handler.go`
- Test: `internal/api/product_handler_test.go`

**Step 1: Write failing test**

```go
func TestProductListRespectsCurrency(t *testing.T) {}
```

**Step 2: Run tests**

Run: `go test ./internal/api -run ProductListRespectsCurrency`
Expected: FAIL

**Step 3: Implement query params**

```go
GET /api/v1/products?currency=USD&locale=en-US
```

**Step 4: Run tests**

Run: `go test ./internal/api -run ProductListRespectsCurrency`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/api/product_handler.go internal/api/category_handler.go internal/api/product_handler_test.go
git commit -m "feat(i18n): apply locale and currency to catalog"
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
git commit -m "chore: finalize i18n and currency"
```
