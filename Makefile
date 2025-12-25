.PHONY: build run test clean install lint fmt

# Build variables
BINARY_NAME=typtui
BUILD_DIR=bin
CMD_DIR=cmd/typtui

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run:
	@go run ./$(CMD_DIR)

# Run with a specific file
run-file:
	@go run ./$(CMD_DIR) $(FILE)

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests for specific package
test-parser:
	@go test -v ./internal/parser/...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./$(CMD_DIR)
	@echo "Install complete"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@echo "Format complete"

# Run go mod tidy
tidy:
	@go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  run-file       - Run with specific file (use FILE=path/to/file.typ)"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-parser    - Run parser tests only"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install binary to GOPATH/bin"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code with gofmt"
	@echo "  tidy           - Run go mod tidy"
	@echo "  help           - Show this help message"
