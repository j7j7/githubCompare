.PHONY: build build-all build-linux build-warwin build-windows test clean install-tools help

# Application name
APP_NAME := githubCompare

# Version
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build flags
LDFLAGS := -X main.version=$(VERSION)

# Directories
DIST_DIR := dist
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build for current platform"
	@echo "  build-all      - Build for all platforms (Linux, macOS, Windows)"
	@echo "  build-linux    - Build for Linux (amd64)"
	@echo "  build-darwin   - Build for macOS (amd64 and arm64)"
	@echo "  build-windows  - Build for Windows (amd64)"
	@echo "  test           - Run tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  clean          - Remove build artifacts"
	@echo "  install-tools  - Install required Go tools"
	@echo "  deps           - Download dependencies"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"

# Build for current platform
build: deps
	@echo "Building $(APP_NAME) for current platform..."
	@mkdir -p $(DIST_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME) .

# Build for all platforms
build-all: build-linux build-darwin build-windows

# Build for Linux
build-linux: deps
	@echo "Building $(APP_NAME) for Linux (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .

# Build for macOS (both Intel and Apple Silicon)
build-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64: deps
	@echo "Building $(APP_NAME) for macOS (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .

build-darwin-arm64: deps
	@echo "Building $(APP_NAME) for macOS (arm64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .

# Build for Windows
build-windows: deps
	@echo "Building $(APP_NAME) for Windows (amd64)..."
	@mkdir -p $(DIST_DIR)
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@go clean

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lint code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Verify everything
verify: fmt vet test
	@echo "All checks passed!"
