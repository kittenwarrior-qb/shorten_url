.PHONY: run build test tidy swagger fmt lint

# Run development server
run:
	go run cmd/main.go

# Build binary
build:
	go build -o bin/app cmd/main.go

# Run tests (basic)
test:
	go test -v -cover ./tests/...

# Run tests with srcipt report
test-report:
	powershell -ExecutionPolicy Bypass -File scripts/test-report.ps1

# Tidy dependencies
tidy:
	go mod tidy

# Generate Swagger docs
swagger:
	swag init -g cmd/main.go -o docs

# Format code
fmt:
	go fmt ./...


