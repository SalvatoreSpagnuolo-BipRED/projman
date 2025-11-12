# Makefile for Projman
# Cross-platform build automation

.PHONY: all help build test clean lint install snapshot release

# ============================================
# Configuration
# ============================================
BINARY_NAME := projman
VERSION := dev
LDFLAGS := -ldflags="-s -w -X github.com/SalvatoreSpagnuolo-BipRED/projman/cmd.Version=$(VERSION)"

# ============================================
# OS Detection
# ============================================
ifeq ($(OS),Windows_NT)
	DETECTED_OS := Windows
	BINARY := $(BINARY_NAME).exe
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
else
	DETECTED_OS := Unix
	BINARY := $(BINARY_NAME)
endif

# ============================================
# Default Target
# ============================================
all: build

# ============================================
# Help
# ============================================
help:
ifeq ($(DETECTED_OS),Windows)
	@Write-Host 'Available commands:' -ForegroundColor Cyan
	@Write-Host '  make build       - Build the binary'
	@Write-Host '  make test        - Run tests'
	@Write-Host '  make lint        - Format and vet code'
	@Write-Host '  make clean       - Remove build artifacts'
	@Write-Host '  make install     - Install binary to GOPATH/bin'
	@Write-Host '  make snapshot    - Test release locally (requires goreleaser)'
	@Write-Host '  make release     - Trigger GitHub release (requires VERSION=vX.Y.Z)'
else
	@echo "Available commands:"
	@echo "  make build       - Build the binary"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Format and vet code"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make install     - Install binary to GOPATH/bin"
	@echo "  make snapshot    - Test release locally (requires goreleaser)"
	@echo "  make release     - Trigger GitHub release (requires VERSION=vX.Y.Z)"
endif

# ============================================
# Build
# ============================================
build:
	@go build $(LDFLAGS) -trimpath -o $(BINARY) .

# ============================================
# Test
# ============================================
test:
	@go test -v ./...

# ============================================
# Lint
# ============================================
lint:
	@go fmt ./...
	@go vet ./...

# ============================================
# Install
# ============================================
install:
	@go install $(LDFLAGS) -trimpath .

# ============================================
# Clean
# ============================================
clean:
ifeq ($(DETECTED_OS),Windows)
	@Remove-Item -Force -ErrorAction SilentlyContinue $(BINARY)
	@Remove-Item -Recurse -Force -ErrorAction SilentlyContinue dist
else
	@rm -f $(BINARY)
	@rm -rf dist/
endif

# ============================================
# Snapshot (local test release)
# ============================================
snapshot:
	@goreleaser release --snapshot --clean

# ============================================
# Release (trigger GitHub Actions)
# ============================================
release:
ifndef VERSION
ifeq ($(DETECTED_OS),Windows)
	@Write-Host 'Error: VERSION required' -ForegroundColor Red
	@Write-Host 'Usage: make release VERSION=v1.2.3' -ForegroundColor Yellow
	@exit 1
else
	@echo "Error: VERSION required"
	@echo "Usage: make release VERSION=v1.2.3"
	@exit 1
endif
endif
ifeq ($(DETECTED_OS),Windows)
	@if ('$(VERSION)' -notmatch '^v[0-9]+\.[0-9]+\.[0-9]+$$') { Write-Host 'Invalid version format. Use vX.Y.Z' -ForegroundColor Red; exit 1 }
	@if (!(Get-Command gh -ErrorAction SilentlyContinue)) { Write-Host 'GitHub CLI not found: https://cli.github.com/' -ForegroundColor Red; exit 1 }
	@gh workflow run release.yml -f version=$(VERSION)
	@Write-Host 'Release $(VERSION) triggered' -ForegroundColor Green
else
	@echo "$(VERSION)" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$' || (echo "Invalid version format. Use vX.Y.Z" && exit 1)
	@command -v gh >/dev/null 2>&1 || (echo "GitHub CLI not found: https://cli.github.com/" && exit 1)
	@gh workflow run release.yml -f version=$(VERSION)
	@echo "Release $(VERSION) triggered"
endif




