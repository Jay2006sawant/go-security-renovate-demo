# Git Repository Security Analyzer - Makefile
BINARY_NAME=analyzer
BUILD_DIR=bin
MAIN_PATH=cmd/analyzer/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
BUILD_FLAGS=-ldflags="-w -s"
RACE_FLAGS=-race

.PHONY: all build clean test deps run demo vuln-info help

# Default target
all: clean deps build test

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build with race detection
build-race:
	@echo "🔨 Building $(BINARY_NAME) with race detection..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(RACE_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-race $(MAIN_PATH)
	@echo "✅ Race detection build complete: $(BUILD_DIR)/$(BINARY_NAME)-race"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies ready"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ Tests complete"

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Run the application directly
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH) --help

# Run demo analysis
demo:
	@echo "🎯 Running demo analysis..."
	$(GOCMD) run $(MAIN_PATH) demo

# Show vulnerability information
vuln-info:
	@echo "🔒 Displaying vulnerability information..."
	$(GOCMD) run $(MAIN_PATH) vulnerability

# Analyze a specific repository
analyze-repo:
	@echo "🔍 Analyzing repository..."
	$(GOCMD) run $(MAIN_PATH) analyze --repo https://github.com/go-git/go-git

# Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@which golangci-lint > /dev/null || (echo "❌ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run
	@echo "✅ Linting complete"

# Security scan (requires gosec)
security-scan:
	@echo "🔒 Running security scan..."
	@which gosec > /dev/null || (echo "❌ gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...
	@echo "✅ Security scan complete"

# Check for vulnerabilities in dependencies
vuln-check:
	@echo "🔍 Checking for vulnerabilities in dependencies..."
	@which govulncheck > /dev/null || (echo "❌ govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest" && exit 1)
	govulncheck ./...
	@echo "✅ Vulnerability check complete"

# Install the application
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "✅ Installation complete"

# Create release build
release: clean
	@echo "🎁 Creating release build..."
	@mkdir -p $(BUILD_DIR)
	
	# Build for multiple platforms
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "✅ Release builds complete:"
	@ls -la $(BUILD_DIR)/

# Quick development workflow
dev: fmt build demo

# Full CI workflow
ci: clean deps fmt lint test security-scan build

# Show project status
status:
	@echo "📊 Project Status"
	@echo "=================="
	@echo "🔒 Vulnerable Dependency: go-git v5.4.2 (CVE-2023-49568)"
	@echo "🤖 Renovate: Configured for vulnerability-only updates"
	@echo "🎯 Purpose: Demonstrate automated security dependency management"
	@echo ""
	@echo "📦 Dependencies:"
	@$(GOMOD) list -m all | head -10
	@echo ""
	@echo "🔧 Available Commands:"
	@echo "  make build     - Build the application"
	@echo "  make demo      - Run demo analysis"
	@echo "  make vuln-info - Show vulnerability details"
	@echo "  make test      - Run tests"
	@echo "  make ci        - Full CI pipeline"

# Help target
help:
	@echo "📚 Git Repository Security Analyzer - Available Commands"
	@echo "========================================================"
	@echo ""
	@echo "🔨 Build Commands:"
	@echo "  build         Build the application"
	@echo "  build-race    Build with race detection"
	@echo "  clean         Clean build artifacts"
	@echo "  deps          Download dependencies"
	@echo "  install       Install to GOPATH/bin"
	@echo "  release       Create multi-platform release builds"
	@echo ""
	@echo "🧪 Test Commands:"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  lint          Run linter (requires golangci-lint)"
	@echo "  security-scan Run security scanner (requires gosec)"
	@echo "  vuln-check    Check for dependency vulnerabilities"
	@echo ""
	@echo "🚀 Run Commands:"
	@echo "  run           Run application with help"
	@echo "  demo          Run demo analysis"
	@echo "  vuln-info     Show vulnerability information"
	@echo "  analyze-repo  Analyze go-git repository"
	@echo ""
	@echo "🔧 Development Commands:"
	@echo "  fmt           Format code"
	@echo "  dev           Quick development workflow (fmt + build + demo)"
	@echo "  ci            Full CI pipeline"
	@echo "  status        Show project status"
	@echo ""
	@echo "💡 Example Usage:"
	@echo "  make build && ./bin/analyzer demo"
	@echo "  make vuln-info"
	@echo "  make ci"