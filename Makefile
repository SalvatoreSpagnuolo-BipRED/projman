# Makefile for Projman

.PHONY: help build test clean lint snapshot release

BINARY_NAME=projman
LDFLAGS=-ldflags="-s -w -X github.com/SalvatoreSpagnuolo-BipRED/projman/cmd.Version=dev"

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@go build $(LDFLAGS) -trimpath -o $(BINARY_NAME) .

test: ## Run tests
	@go test -v ./...

clean: ## Remove artifacts
	@rm -f $(BINARY_NAME)
	@rm -rf dist/

lint: ## Format and vet
	@go fmt ./...
	@go vet ./...

snapshot: ## Test release locally
	@goreleaser release --snapshot --clean

release: ## Trigger release (make release VERSION=v1.2.3)
ifndef VERSION
	@echo "❌ Usage: make release VERSION=v1.2.3"
	@exit 1
endif
	@echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$' || (echo "❌ Invalid version format. Use vX.Y.Z" && exit 1)
	@command -v gh >/dev/null 2>&1 || (echo "❌ GitHub CLI not found. Install: https://cli.github.com/" && exit 1)
	@gh workflow run release.yml -f version=$(VERSION)
	@echo "✓ Release $(VERSION) triggered"
	@echo "  Monitor: gh run list --workflow=release.yml"




