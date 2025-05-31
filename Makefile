# 🔥 Azure OpenAI Plugin for Firebase Genkit - Makefile
# Pure vibe coded SDK ☕️

.PHONY: help build test test-verbose test-race test-coverage clean lint format check install-tools deps update security bench example

# Default target
all: check test build

# 📋 Show help information
help:
	@echo "🔥 Azure OpenAI Plugin for Firebase Genkit - Development Commands"
	@echo ""
	@echo "📋 Available commands:"
	@echo "  build          🏗️  Build the project"
	@echo "  test           🧪  Run tests"
	@echo "  test-verbose   🔍  Run tests with verbose output"
	@echo "  test-race      🏃  Run tests with race detection"
	@echo "  test-coverage  📊  Run tests with coverage report"
	@echo "  lint           🧹  Run linter"
	@echo "  format         📝  Format code"
	@echo "  check          ✅  Run all quality checks"
	@echo "  clean          🗑️   Clean build artifacts"
	@echo "  install-tools  🔧  Install development tools"
	@echo "  deps           📦  Download dependencies"
	@echo "  update         ⬆️   Update dependencies"
	@echo "  security       🔒  Run security scans"
	@echo "  bench          ⚡  Run benchmarks"
	@echo "  example        🚀  Run example"
	@echo ""

# 🏗️ Build the project
build:
	@echo "🏗️ Building project..."
	go build -v ./...
	@echo "✅ Build complete!"

# 🧪 Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...
	@echo "✅ Tests passed!"

# 🔍 Run tests with verbose output
test-verbose:
	@echo "🔍 Running tests (verbose)..."
	go test -v ./...

# 🏃 Run tests with race detection
test-race:
	@echo "🏃 Running tests with race detection..."
	go test -race ./...
	@echo "✅ Race tests passed!"

# 📊 Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# 🧹 Run linter
lint:
	@echo "🧹 Running linter..."
	@which golangci-lint > /dev/null || (echo "❌ golangci-lint not found. Run 'make install-tools' first" && exit 1)
	golangci-lint run
	@echo "✅ Linting complete!"

# 📝 Format code
format:
	@echo "📝 Formatting code..."
	gofmt -s -w .
	@which goimports > /dev/null && goimports -w . || true
	@echo "✅ Code formatted!"

# ✅ Run all quality checks
check: format lint test-race
	@echo "✅ All quality checks passed!"

# 🗑️ Clean build artifacts
clean:
	@echo "🗑️ Cleaning build artifacts..."
	go clean -cache -testcache -modcache
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete!"

# 🔧 Install development tools
install-tools:
	@echo "🔧 Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "✅ Development tools installed!"

# 📦 Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod verify
	@echo "✅ Dependencies ready!"

# ⬆️ Update dependencies
update:
	@echo "⬆️ Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies updated!"

# 🔒 Run security scans
security:
	@echo "🔒 Running security scans..."
	@which gosec > /dev/null || (echo "❌ gosec not found. Run 'make install-tools' first" && exit 1)
	@which govulncheck > /dev/null || (echo "❌ govulncheck not found. Run 'make install-tools' first" && exit 1)
	gosec ./...
	govulncheck ./...
	@echo "✅ Security scans complete!"

# ⚡ Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	go test -bench=. -benchmem ./...
	@echo "✅ Benchmarks complete!"

# 🚀 Run example
example:
	@echo "🚀 Running example..."
	go test -v -run Example ./...
	@echo "✅ Example complete!"

# 📝 Check go mod tidy
mod-tidy:
	@echo "📝 Checking go mod tidy..."
	go mod tidy
	@git diff --exit-code go.mod go.sum || (echo "❌ go.mod or go.sum is not tidy. Run 'go mod tidy' and commit the changes." && exit 1)
	@echo "✅ go.mod and go.sum are tidy!"

# 🔍 Show project info
info:
	@echo "🔥 Azure OpenAI Plugin for Firebase Genkit"
	@echo "📁 Project: $(shell basename $(PWD))"
	@echo "🔧 Go version: $(shell go version)"
	@echo "📦 Module: $(shell go list -m)"
	@echo "🌟 Pure vibe coded SDK - Built with Cursor, Claude Sonnet 4, and yuanyang ☕️" 