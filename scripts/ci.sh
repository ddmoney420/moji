#!/usr/bin/env bash
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

PASS=0
FAIL=0
WARN=0

step() {
    echo -e "\n${CYAN}${BOLD}==> $1${NC}"
}

pass() {
    echo -e "  ${GREEN}✓${NC} $1"
    PASS=$((PASS + 1))
}

fail() {
    echo -e "  ${RED}✗${NC} $1"
    FAIL=$((FAIL + 1))
}

warn() {
    echo -e "  ${YELLOW}!${NC} $1"
    WARN=$((WARN + 1))
}

echo -e "${BOLD}moji local CI${NC}"
echo "────────────────────────────────────"

# Go version
step "Environment"
echo "  Go: $(go version | awk '{print $3}')"
echo "  OS: $(uname -s)/$(uname -m)"
echo "  Dir: $(pwd)"

# Format check
step "Formatting (gofmt)"
UNFORMATTED=$(gofmt -l . 2>&1 || true)
if [ -z "$UNFORMATTED" ]; then
    pass "All files formatted"
else
    fail "Files need formatting:"
    echo "$UNFORMATTED" | while read -r f; do echo "    $f"; done
fi

# Go vet
step "Vet (go vet)"
if go vet ./... 2>&1; then
    pass "No issues found"
else
    fail "go vet found issues"
fi

# Lint
step "Lint (golangci-lint)"
GOLANGCI_LINT="$(go env GOPATH)/bin/golangci-lint"
if [ -f "$GOLANGCI_LINT" ]; then
    if $GOLANGCI_LINT run ./... 2>&1; then
        pass "No lint issues"
    else
        fail "Lint issues found"
    fi
else
    warn "golangci-lint not installed (run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)"
fi

# Tests
step "Tests (go test)"
if go test ./... -count=1 2>&1; then
    pass "All tests pass"
else
    fail "Tests failed"
fi

# Test coverage
step "Coverage"
COVERAGE=$(go test ./... -coverprofile=/tmp/moji-coverage.out -count=1 2>&1 | grep -o '[0-9]*\.[0-9]*%' | head -1 || echo "0%")
rm -f /tmp/moji-coverage.out
echo "  Coverage: $COVERAGE"

# Build
step "Build"
if go build -o moji . 2>&1; then
    SIZE=$(ls -lh moji | awk '{print $5}')
    pass "Binary built (${SIZE})"
else
    fail "Build failed"
fi

# Quick smoke test
step "Smoke test"
if ./moji --version &>/dev/null; then
    pass "moji --version"
else
    fail "moji --version"
fi

if ./moji shrug 2>/dev/null | grep -q "ツ"; then
    pass "moji shrug"
else
    fail "moji shrug"
fi

if ./moji banner "CI" --font standard &>/dev/null; then
    pass "moji banner"
else
    fail "moji banner"
fi

if ./moji qr "test" &>/dev/null; then
    pass "moji qr"
else
    fail "moji qr"
fi

if ./moji doctor &>/dev/null; then
    pass "moji doctor"
else
    fail "moji doctor"
fi

# Summary
echo ""
echo "────────────────────────────────────"
echo -e "${BOLD}Results:${NC}"
echo -e "  ${GREEN}${PASS} passed${NC}  ${RED}${FAIL} failed${NC}  ${YELLOW}${WARN} warnings${NC}"
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo -e "${RED}${BOLD}CI FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}${BOLD}CI PASSED${NC}"
fi
