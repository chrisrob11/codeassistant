# Makefile for basic Go project tasks

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOCMD=$(GOBASE)/cmd
GO_ENV_PATH=$(shell go env GOPATH)

# List of target OS/ARCH combinations
PLATFORMS = linux/amd64 darwin/amd64 darwin/arm64

# Compile and place binaries into the bin directory.
build:
	@echo "  >  Building binaries..."
	@GOBIN=$(GOBIN) CGO_ENABLED=0 go install  $(GOCMD)/...

# Cross-compile for multiple platforms
build-all:
	@echo "  >  Building binaries for all platforms..."
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; ARCH=$${platform#*/}; \
		OUTPUT_DIR=$(GOBIN)/$$OS/$$ARCH; \
		mkdir -p $$OUTPUT_DIR; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=0 go build -o $$OUTPUT_DIR/ca $(GOCMD)/... || exit 1; \
	done

# Run linters using golangci-lint.
lint:
	@echo "  >  Linting code..."
	${GO_ENV_PATH}/bin/golangci-lint --color always run

test:
	@echo ">> running tests"
	@go test -race ./...

coverage:
	@echo ">> running tests with coverage on"
	@go test -coverprofile out.coverage -v -race ./...
	@go tool cover -html out.coverage
# Execute all steps: generate code, lint, and build.
all: generate lint build ./...

tools:
	@echo "  >  Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

migrate-up:
	@echo "  >  Migrating sqlite..."
	sql-migrate up --config=./sqllite_migration_config.yml
	@echo "  >  Migrating sqlite test db..."
	sql-migrate up --config=./sqllite_migration_test_config.yml

# Clean up generated files and binaries.
clean:
	@echo "  >  Cleaning build cache"
	@go clean ./...
	@echo "  >  Removing binaries..."
	@rm -rf $(GOBIN)/*

.PHONY: build generate lint all clean docker-build migrate-up
