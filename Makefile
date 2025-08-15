# SQL CLI Tool Makefile

.PHONY: build install clean test deps check-deps run

# Variables
BINARY_NAME=db-cli
BUILD_DIR=./build
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
all: build

# Build the binary
build: deps
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Install to system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/"
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Check system dependencies
check-deps:
	@echo "Checking system dependencies..."
	@command -v psql >/dev/null 2>&1 || echo "Warning: psql not found (required for PostgreSQL)"
	@command -v mysql >/dev/null 2>&1 || echo "Warning: mysql not found (required for MySQL)"
	@command -v sqlite3 >/dev/null 2>&1 || echo "Warning: sqlite3 not found (required for SQLite)"
	@command -v redis-cli >/dev/null 2>&1 || echo "Warning: redis-cli not found (required for Redis)"

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	go clean

# Run the application (for development)
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Development helpers
dev-run:
	go run . $(ARGS)

format:
	go fmt ./...

lint:
	golangci-lint run

# # Create example config
# example-config:
# 	@mkdir -p ~/.sql
# 	@cat > ~/.sql/config.toml << 'EOF'
# [hosts]
# cd = "postgres://postgres:@127.0.0.1:5432/database"
# pg = "postgres://postgres:@127.0.0.1:5432/postgres"  
# sakila = "sqlite3:///$$HOME/Projects/personal/sql/datasets/sakila.db"
# red = "redis://127.0.0.1:6379"
# EOF
# 	@echo "Example config created at ~/.sql/config.toml"

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  install      - Install to /usr/local/bin"
	@echo "  deps         - Install Go dependencies"
	@echo "  check-deps   - Check system dependencies"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run (use ARGS=... for arguments)"
	@echo "  dev-run      - Run without building (use ARGS=... for arguments)"
	@echo "  format       - Format code"
	@echo "  lint         - Run linter"
	@echo "  example-config - Create example config file"
	@echo "  help         - Show this help"
