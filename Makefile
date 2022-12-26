PROJECT_BINARY=chainx
PROJECT_BINARY_OUTPUT=output

.PHONY: all

all: help

## Build:
tidy: ## Tidy project
	@go mod tidy

clean: ## Cleans temporary folder
	@rm -rf ${PROJECT_BINARY_OUTPUT}
	@rm -rf ${PROJECT_RELEASER_OUTPUT}

build: clean tidy ## Builds project
	@GO111MODULE=on CGO_ENABLED=0 go build -ldflags="-w -s" -o ${PROJECT_BINARY_OUTPUT}/bin/${PROJECT_BINARY} main.go

run: build ## Run example
	@${PROJECT_BINARY_OUTPUT}/bin/${PROJECT_BINARY}

test: clean tidy ## Run unit tests
	@go clean -testcache
	@go test -v ./... 

coverage: clean tidy ## Run code coverage
	@go clean -testcache
	@go test ./... -coverprofile=./docs/coverage.out

bench: test ## Run benchmarks
	@go clean -testcache
	@go test ./... -bench=.

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    %-20s%s\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  %s\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
