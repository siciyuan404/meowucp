.PHONY: all run build test clean migrate-up migrate-down help

all: run

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/admin cmd/admin/main.go

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

migrate-up:
	psql -h localhost -U meowucp -d meowucp -f migrations/001_init.sql

migrate-down:
	psql -h localhost -U meowucp -d meowucp -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

help:
	@echo "Available commands:"
	@echo "  make run           - Run the API server"
	@echo "  make build         - Build the applications"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make help          - Show this help message"
