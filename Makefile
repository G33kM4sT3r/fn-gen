BINARY_NAME := fn-gen
CMD_PATH := ./cmd/fn-gen
BUILD_DIR := dist

# Version info (can be overridden)
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build build-all run install test lint fmt vet check deps clean help

## build: Build binary for current platform
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

## build-all: Build binaries for all platforms
build-all: clean
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}$$([ "$${platform%/*}" = "windows" ] && echo ".exe") $(CMD_PATH); \
	done

## run: Run the application
run:
	@go run $(CMD_PATH)

## install: Install binary to GOPATH/bin
install: build
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## test: Run tests
test:
	go test -race ./...

## lint: Run linter
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	gofmt -s -w .

## vet: Run go vet
vet:
	go vet ./...

## check: Run all checks (fmt, vet, lint, test)
check: vet lint test

## deps: Download and tidy dependencies
deps:
	go mod download
	go mod tidy

clean:
	@rm -rf $(BUILD_DIR)

help:
	@echo "Available targets:"
	@echo "  build     - Build for current platform"
	@echo "  build-all - Build for all platforms"
	@echo "  run       - Run the application"
	@echo "  install   - Install binary to GOPATH/bin"
	@echo "  test      - Run tests"
	@echo "  lint      - Run linter"
	@echo "  fmt       - Format code"
	@echo "  vet       - Run go vet"
	@echo "  check     - Run all checks (vet, lint, test)"
	@echo "  deps      - Download and tidy dependencies"
	@echo "  clean     - Remove build artifacts"
	@echo "  help      - Show this help"
