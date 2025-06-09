# Makefile for rsss

.PHONY: build test clean run-tui run-cli install lint fmt vet

# Binary name
BINARY_NAME=rsss
BUILD_DIR=build
CMD_DIR=cmd/rsss

# Build the binary
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	go clean

# Run TUI mode
run-tui:
	go run ./$(CMD_DIR) --menu

# Run CLI mode with BBC News
run-cli:
	go run ./$(CMD_DIR) https://feeds.bbci.co.uk/news/rss.xml

# Install the binary to GOPATH/bin
install:
	go install ./$(CMD_DIR)

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  run-tui       - Run in TUI mode"
	@echo "  run-cli       - Run in CLI mode with BBC News"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  check         - Run fmt, vet, and test"
	@echo "  help          - Show this help"