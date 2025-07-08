# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=dehydrated-api-metadata-plugin-netscaler

.PHONY: all build clean test deps tidy test-integration test-all

all: clean build

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test: ## Run unit tests
	$(GOTEST) -v ./...

test-integration: ## Run integration tests
	$(GOTEST) -v -tags=integration ./...

test-all: test test-integration ## Run all tests

deps: ## Install dependencies
	$(GOGET) -v ./...

tidy: ## Tidy up dependencies
	$(GOMOD) tidy

# Development helpers
.PHONY: lint fmt

lint: ## Run linter
	@golangci-lint run

lint-fix: ## Run linter (and fix issues if possible)
	@golangci-lint run --fix

fmt: ## Format the code
	$(GOCMD) fmt ./...

release: ## Create a release with goreleaser
	@goreleaser release --snapshot --clean

#
# Help
#

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Integration Test Environment Variables:"
	@echo "  NETSCALER_ENDPOINT   - Netscaler endpoint URL"
	@echo "  NETSCALER_USERNAME   - Netscaler username"
	@echo "  NETSCALER_PASSWORD   - Netscaler password"
	@echo "  NETSCALER_PREFIX     - Certificate prefix (default: test-)"
	@echo "  NETSCALER_SSL_VERIFY - Enable SSL verification (default: false)"
	@echo "----------------------------------------"
	@echo "For more information, see the README.md file."
