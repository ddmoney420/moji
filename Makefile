VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%d")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
LDFLAGS := -ldflags "-s -w \
	-X github.com/ddmoney420/moji/internal/ux.Version=$(VERSION) \
	-X github.com/ddmoney420/moji/internal/ux.BuildDate=$(BUILD_DATE) \
	-X github.com/ddmoney420/moji/internal/ux.GitCommit=$(GIT_COMMIT)"

BINARY := moji
GOBIN ?= /usr/local/bin

.PHONY: all build install test lint vet fmt ci clean help wasm web web-serve web-clean

all: ci ## Run full CI pipeline

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY) .

install: build ## Build and install to GOBIN
	cp $(BINARY) $(GOBIN)/$(BINARY)

test: ## Run tests
	go test ./... -v -count=1

test-short: ## Run tests (short mode)
	go test ./... -short -count=1

test-cover: ## Run tests with coverage
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -func=coverage.out
	@rm -f coverage.out

GOLANGCI_LINT := $(shell go env GOPATH)/bin/golangci-lint

lint: ## Run golangci-lint
	@test -f $(GOLANGCI_LINT) || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	$(GOLANGCI_LINT) run ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Check formatting
	@test -z "$$(gofmt -l .)" || { echo "Files need formatting:"; gofmt -l .; exit 1; }

fmt-fix: ## Fix formatting
	gofmt -w .

tidy: ## Tidy go.mod
	go mod tidy

ci: fmt vet lint test build ## Run full local CI (fmt, vet, lint, test, build)
	@echo ""
	@echo "=== CI PASSED ==="
	@echo "Binary: ./$(BINARY) ($(VERSION))"

ci-quick: fmt vet test-short build ## Quick CI (skip lint, short tests)
	@echo ""
	@echo "=== QUICK CI PASSED ==="

vhs: build ## Generate terminal recordings (requires vhs)
	@which vhs > /dev/null 2>&1 || { echo "Install vhs: https://github.com/charmbracelet/vhs"; exit 1; }
	vhs vhs/demo.tape
	vhs vhs/banner.tape
	vhs vhs/filters.tape
	vhs vhs/tui.tape

release: ## Create release (requires goreleaser)
	goreleaser release --clean

release-snapshot: ## Test release locally
	goreleaser release --snapshot --clean

wasm: ## Build WASM binary
	GOOS=js GOARCH=wasm go build -o web/wasm/moji.wasm ./cmd/wasm/
	@echo "WASM binary: web/wasm/moji.wasm ($$(du -h web/wasm/moji.wasm | cut -f1))"

GOROOT_DIR := $(shell go env GOROOT)
WASM_EXEC := $(shell test -f "$(GOROOT_DIR)/lib/wasm/wasm_exec.js" && echo "$(GOROOT_DIR)/lib/wasm/wasm_exec.js" || echo "$(GOROOT_DIR)/misc/wasm/wasm_exec.js")

web: wasm ## Build WASM and copy wasm_exec.js
	cp "$(WASM_EXEC)" web/js/wasm_exec.js
	@echo "Web assets ready in web/"

web-serve: web ## Serve web locally on :8080
	@echo "Serving at http://localhost:8080"
	cd web && python3 -m http.server 8080

web-clean: ## Remove WASM build artifacts
	rm -f web/wasm/moji.wasm web/js/wasm_exec.js

clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out
	rm -f vhs/*.gif

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
