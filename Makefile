# Variables
GO := go
AGENT_BINARY := app_agent.linux64
API_BINARY := app_api.linux64
CLI_BINARY := app_cli.linux64
ORGANIZATION_BINARY := app_org.linux64
PLATFORM_BINARY := app_platform.linux64
VERSION := v1.0.0
BUILD_TIME := $(shell date -u +'%Y-%m-%d %H:%M:%S')
LDFLAGS := -ldflags="-s -w\
	-X 'main.version=$(VERSION)' \
	-X 'main.buildTime=$(BUILD_TIME)'"

# Tools
GOLANGCI_LINT := $(shell which golangci-lint 2>/dev/null || echo "")

.PHONY: all build clean test lint fmt vet run help

all: build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(AGENT_BINARY) ./cmd/agent/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(API_BINARY) ./cmd/api/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(CLI_BINARY) ./cmd/cli/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(ORGANIZATION_BINARY) ./cmd/organization/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(PLATFORM_BINARY) ./cmd/platform/main.go

clean:
	rm -rf $(AGENT_BINARY)
	rm -rf $(API_BINARY)
	rm -rf $(CLI_BINARY)
	rm -rf $(ORGANIZATION_BINARY)
	rm -rf $(PLATFORM_BINARY)

lint:
	@if [ -z "$(GOLANGCI_LINT)" ]; then \
		echo "golangci-lint not found. Install from https://golangci-lint.run/"; \
		exit 1; \
	fi
	golangci-lint run
fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

run: build
	./$(BINARY)

help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove binary"
	@echo "  test          - Run tests"
	@echo "  lint          - Run golangci-lint"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  run           - Build and run"
