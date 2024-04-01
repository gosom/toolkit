MODULE_PKG := toolkit

GIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y%m%dT%H%M%S')

LDFLAGS := -ldflags="-s -w -X ${MODULE_PKG}/app.BuildDate=${BUILD_DATE} -X $(MODULE_PKG)/app.Commit=${GIT_HASH}"

default: help

help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


lint: ## runs the linter
	@golangci-lint -v run ./...
.PHONY: lint

build: ## builds the binary
	@go build -o /dev/null $(LDFLAGS) ./...

test: ## runs the unit tests
	@go test -v -race -timeout 5m ./... -failfast
.PHONY: test

vet: ## runs go vet
	@go vet ./...
.PHONY: vet

format: ## runs go fmt
	@gofmt -s -w .
.PHONY: format
