#!/bin/bash

# Test runner script for dockenv

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

print_status "Running dockenv test suite..."
print_status "Project root: $PROJECT_ROOT"

# Change to project root
cd "$PROJECT_ROOT"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3)
print_status "Using Go version: $GO_VERSION"

# Run go mod tidy to ensure dependencies are correct
print_status "Tidying Go modules..."
go mod tidy

# Run unit tests
print_status "Running unit tests..."
if go test -v ./tests/unit/...; then
    print_status "Unit tests passed ✓"
else
    print_error "Unit tests failed ✗"
    exit 1
fi

# Build the binary for integration tests
print_status "Building dockenv binary..."
if go build -o build/dockenv ./main.go; then
    print_status "Build successful ✓"
else
    print_error "Build failed ✗"
    exit 1
fi

# Run integration tests
print_status "Running integration tests..."
cd tests/integration
if go test -v .; then
    print_status "Integration tests passed ✓"
else
    print_error "Integration tests failed ✗"
    cd "$PROJECT_ROOT"
    exit 1
fi

cd "$PROJECT_ROOT"

# Run tests with coverage
print_status "Running tests with coverage..."
if go test -coverprofile=coverage.out ./...; then
    print_status "Coverage analysis completed ✓"
    go tool cover -html=coverage.out -o coverage.html
    print_status "Coverage report generated: coverage.html"
else
    print_warning "Coverage analysis failed, but continuing..."
fi

# Run go vet for static analysis
print_status "Running go vet..."
if go vet ./...; then
    print_status "Static analysis passed ✓"
else
    print_warning "Static analysis found issues"
fi

# Check for go fmt issues
print_status "Checking code formatting..."
UNFORMATTED=$(gofmt -l .)
if [ -z "$UNFORMATTED" ]; then
    print_status "Code formatting is correct ✓"
else
    print_warning "The following files need formatting:"
    echo "$UNFORMATTED"
    print_warning "Run 'go fmt ./...' to fix formatting issues"
fi

# Check for ineffassign (if available)
if command -v ineffassign &> /dev/null; then
    print_status "Running ineffassign..."
    if ineffassign ./...; then
        print_status "No ineffective assignments found ✓"
    else
        print_warning "Ineffective assignments found"
    fi
else
    print_warning "ineffassign not installed, skipping check"
    print_status "Install with: go install github.com/gordonklaus/ineffassign@latest"
fi

# Check for misspell (if available)
if command -v misspell &> /dev/null; then
    print_status "Running misspell check..."
    if misspell -error .; then
        print_status "No spelling errors found ✓"
    else
        print_warning "Spelling errors found"
    fi
else
    print_warning "misspell not installed, skipping check"
    print_status "Install with: go install github.com/client9/misspell/cmd/misspell@latest"
fi

print_status "All tests completed successfully! 🎉"
print_status ""
print_status "Summary:"
print_status "✓ Unit tests passed"
print_status "✓ Integration tests passed"
print_status "✓ Binary builds successfully"
print_status "✓ Static analysis completed"
print_status ""
print_status "Ready for publication! 🚀"
