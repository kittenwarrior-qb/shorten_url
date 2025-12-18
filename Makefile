.PHONY: run build test tidy swagger fmt lint

# Run development server
run:
	go run cmd/main.go

# Build binary
build:
	go build -o bin/app cmd/main.go

# Run tests
test:
	go test -v -cover ./...

# Tidy dependencies
tidy:
	go mod tidy

# Generate Swagger docs
swagger:
	swag init -g cmd/main.go -o docs

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run
