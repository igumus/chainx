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

build-node: clean tidy ## Build basic build-node
	@GO111MODULE=on CGO_ENABLED=0 go build -ldflags="-w -s" -o ${PROJECT_BINARY_OUTPUT}/bin/node cmd/node/main.go

build-vnode: clean tidy ## Build basic build-node
	@GO111MODULE=on CGO_ENABLED=0 go build -ldflags="-w -s" -o ${PROJECT_BINARY_OUTPUT}/bin/vnode cmd/vnode/main.go

build: build-node build-vnode ## Builds project
	@echo "Building Status: DONE"

test: build ## Run unit tests
	@go clean -testcache
	@go test ./... 

pre-commit: test ## Checks everything is allright
	@echo "Commit Status: OK"

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
