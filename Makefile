.PHONY: help deps docker-up docker-down run migrate test test-coverage test-html test-verbose test-bdd ac-coverage bdd-report clean

help:
	@echo "Coolmate eCommerce Backend - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make deps         - Download Go dependencies"
	@echo "  make run          - Start the server"
	@echo "  make docker-up    - Start PostgreSQL, Redis, MinIO"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make build        - Build binary"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run unit tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make test-coverage - Run tests and generate HTML report"
	@echo "  make test-html    - Run tests, generate HTML report & open it"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make docker-clean - Remove Docker volumes"

deps:
	@echo "Downloading dependencies..."
	go mod tidy
	go mod download

run:
	@echo "Starting server on port 8080..."
	go run cmd/api/main.go

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "MinIO: localhost:9000 (http://localhost:9001 for console)"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-clean:
	@echo "Removing Docker volumes..."
	docker-compose down -v

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage and generating HTML report..."
	go test ./... -coverprofile=coverage.out -timeout 120s
	go tool cover -html=coverage.out -o coverage_report.html
	@echo "✅ HTML Report: coverage_report.html"
	@echo "📈 Coverage: $$(go tool cover -func=coverage.out | tail -1)"

test-html:
	@echo "Running tests and generating HTML report..."
	go test ./... -coverprofile=coverage.out -timeout 120s
	go tool cover -html=coverage.out -o coverage_report.html
	@echo "✅ Done! Opening report..."
	start coverage_report.html 2>/dev/null || open coverage_report.html 2>/dev/null || xdg-open coverage_report.html

test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -cover ./... -timeout 120s

test-bdd:
	@echo "Running BDD acceptance tests (TestUS_*)..."
	go test -v -run "^TestUS_" -timeout 120s ./internal/...

ac-coverage:
	@bash scripts/ac-coverage.sh

bdd-report: test-bdd ac-coverage
	@echo ""
	@echo "BDD report complete. Line coverage stays green via existing AAA tests."

clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage_report.html bin/api
	go clean

build:
	@echo "Building binary..."
	go build -o bin/api cmd/api/main.go
