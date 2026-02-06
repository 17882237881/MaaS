.PHONY: help build run test clean

help:
	@echo "Available targets:"
	@echo "  build      - Build all services"
	@echo "  run-api    - Run API Gateway"
	@echo "  run-model  - Run Model Registry"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"

build:
	@echo "Building API Gateway..."
	@cd api-gateway && go build -o ../bin/api-gateway ./cmd/main.go
	@echo "Building Model Registry..."
	@cd model-registry && go build -o ../bin/model-registry ./cmd/main.go
	@echo "Build complete!"

run-api:
	@cd api-gateway && go run ./cmd/main.go

run-model:
	@cd model-registry && go run ./cmd/main.go

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@rm -rf bin/
	@echo "Clean complete!"
