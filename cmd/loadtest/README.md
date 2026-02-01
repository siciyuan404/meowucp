# Load Test Harness

Run a basic concurrency test against product listing.

## Usage

```bash
go run cmd/loadtest/main.go --endpoint http://localhost:8080
```

## Options

- `--endpoint`: Base API endpoint (default `http://localhost:8080`)
- `--concurrency`: Concurrent workers (default `10`)
- `--iterations`: Requests per worker (default `50`)
