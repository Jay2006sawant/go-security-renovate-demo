# Makefile for Git Repository Security Analyzer

BINARY_NAME=analyzer
BUILD_DIR=bin
MAIN_PATH=cmd/analyzer/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Build flags
BUILD_FLAGS=-ldflags="-w -s"
RACE_FLAGS=-race

.PHONY: all build clean test run deps fmt lint security-scan vuln-check install release help

# Default target: clean, dependencies, build, test
all: clean deps build test

# Build the Go application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build with race detector enabled
build-race:
	@echo "🔨 Building $(BINARY_NAME) with race detection..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(RACE_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-race $(MAIN_PATH)
	@echo "✅ Race detection build complete: $(BUILD_DIR)/$(BINARY_NAME)-race"

# Clean build artifacts and bin folder
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Download and tidy dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies ready"

# Run tests verbosely
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ Tests complete"

# Run the app with --help flag
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH) --help

# Format all Go code
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

# Lint the code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@which golangci-lint > /dev/null || (echo "❌ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run
	@echo "✅ Linting complete"

# Run security scan (requires gosec)
security-scan:
	@echo "🔒 Running security scan..."
	@which gosec > /dev/null || (echo "❌ gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...
	@echo "✅ Security scan complete"

# Check for vulnerabilities in dependencies (requires govulncheck)
vuln-check:
	@echo "🔍 Checking for vulnerabilities..."
	@which govulncheck > /dev/null || (echo "❌ govulncheck not installed. Install with: go install golang.org/x/vuln/cmd/govulncheck@latest" && exit 1)
	govulncheck ./...
	@echo "✅ Vulnerability check complete"

# Install binary to GOPATH/bin
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "✅ Installation complete"

# Create release builds for multiple platforms
release: clean
	@echo "🎁 Creating release builds..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "✅ Release builds created:"
	@ls -la $(BUILD_DIR)/

# Show help message with available commands
help:
	@echo "📚 Available Makefile commands:"
	@echo "  all          : clean, download deps, build, test"
	@echo "  build        : build the application"
	@echo "  build-race   : build with race detector"
	@echo "  clean        : remove build artifacts"
	@echo "  deps         : download Go dependencies"
	@echo "  test         : run all tests"
	@echo "  run          : run application with --help"
	@echo "  fmt          : format source code"
	@echo "  lint         : lint source code (requires golangci-lint)"
	@echo "  security-scan: run gosec security scan"
	@echo "  vuln-check   : check for vulnerable dependencies"
	@echo "  install      : install binary to GOPATH/bin"
	@echo "  release      : build multi-platform release binaries"
	@echo "  help         : show this help message"
