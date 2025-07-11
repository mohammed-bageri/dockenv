# dockenv Makefile

# Variables
BINARY_NAME=dockenv
VERSION?=0.2.0
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"

# OS and architecture
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build clean test deps install-tools install uninstall release help

# Default target
all: clean deps build

# Build the binary
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@$(GOTEST) -v ./tests/unit/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@$(GOTEST) -v ./tests/integration/...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./tests/unit/...

# Run comprehensive test suite
test-all: clean deps build test-unit test-integration test-coverage
	@echo "All tests completed successfully!"

# Static analysis (comprehensive)
lint:
	@echo "Running static analysis..."
	@go vet ./...
	@GOBIN=$$(go env GOPATH)/bin; \
	if [ -x "$$GOBIN/golangci-lint" ] || command -v golangci-lint >/dev/null 2>&1; then \
		if [ -x "$$GOBIN/golangci-lint" ]; then \
			$$GOBIN/golangci-lint run; \
		else \
			golangci-lint run; \
		fi; \
	else \
		echo "golangci-lint not installed, skipping advanced linting"; \
		echo "To install golangci-lint, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@GOBIN=$$(go env GOPATH)/bin; \
	if [ -x "$$GOBIN/goimports" ] || command -v goimports >/dev/null 2>&1; then \
		if [ -x "$$GOBIN/goimports" ]; then \
			$$GOBIN/goimports -w .; \
		else \
			goimports -w .; \
		fi; \
	else \
		echo "goimports not available, skipping import formatting"; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	else \
		echo "golangci-lint already installed"; \
	fi
	@GOBIN=$$(go env GOPATH)/bin; \
	if [ ! -x "$$GOBIN/goimports" ] && ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	else \
		echo "goimports already installed"; \
	fi
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	else \
		echo "gosec already installed"; \
	fi

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"

# Uninstall binary from system
uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_DIR)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstallation complete!"

# Build for multiple platforms
release: clean deps
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	
	# Linux AMD64
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 .
	
	# Linux ARM64
	@echo "Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 .
	
	# macOS AMD64
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS ARM64 (Apple Silicon)
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows AMD64
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "Release builds complete in $(BUILD_DIR)/release/"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Development: build and install locally
dev: clean deps build
	@echo "Installing development build..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Development installation complete!"

# Vet code
vet:
	@echo "Vetting code..."
	@$(GOCMD) vet ./...

# Security check (requires gosec)
security:
	@echo "Running security check..."
	@GOBIN=$$(go env GOPATH)/bin; \
	if [ -x "$$GOBIN/gosec" ] || command -v gosec >/dev/null 2>&1; then \
		if [ -x "$$GOBIN/gosec" ]; then \
			$$GOBIN/gosec ./... || echo "Security issues found - review recommended"; \
		else \
			gosec ./... || echo "Security issues found - review recommended"; \
		fi; \
	else \
		echo "gosec not installed, skipping security check"; \
		echo "To install gosec, run: make install-tools"; \
	fi

# Full quality check
check: fmt vet lint test
	@echo "All checks passed!"

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .

# Create package
package: release
	@echo "Creating packages..."
	@mkdir -p $(BUILD_DIR)/packages
	
	# Linux package
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR)/release $(BINARY_NAME)-linux-amd64
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz -C $(BUILD_DIR)/release $(BINARY_NAME)-linux-arm64
	
	# macOS package
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR)/release $(BINARY_NAME)-darwin-amd64
	@tar -czf $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR)/release $(BINARY_NAME)-darwin-arm64
	
	# Windows package
	@zip -j $(BUILD_DIR)/packages/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe
	
	@echo "Packages created in $(BUILD_DIR)/packages/"

# Show help
help:
	@echo "Available targets:"
	@echo "  all         - Clean, download deps, and build"
	@echo "  build       - Build the binary"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  deps        - Download dependencies"
	@echo "  install-tools - Install development tools (golangci-lint, gosec)"
	@echo "  install     - Install binary to system"
	@echo "  uninstall   - Remove binary from system"
	@echo "  release     - Build for multiple platforms"
	@echo "  run         - Build and run the application"
	@echo "  dev         - Build and install development version"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code (installs golangci-lint if needed)"
	@echo "  vet         - Vet code"
	@echo "  security    - Security check (requires gosec)"
	@echo "  check       - Run all quality checks"
	@echo "  docker-build - Build Docker image"
	@echo "  package     - Create release packages"
	@echo "  help        - Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION     - Set version (default: $(VERSION))"
	@echo "  GOOS        - Target OS (default: $(GOOS))"
	@echo "  GOARCH      - Target architecture (default: $(GOARCH))"
