#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Quick Test - Rebuild & Test (no DB reset)${NC}"
echo -e "${BLUE}========================================${NC}"

# Detect environment
if [ -n "$CI" ]; then
    echo -e "${YELLOW}Running in CI environment${NC}"
    TEMPL_CMD="templ"
else
    echo -e "${YELLOW}Running in local environment${NC}"
    TEMPL_CMD="$HOME/go/bin/templ"
fi

# Step 1: Kill any running server
echo -e "\n${YELLOW}[1/4] Stopping any running server...${NC}"
pkill -f "go run ./cmd/server" 2>/dev/null || true
pkill -f "bin/server" 2>/dev/null || true
sleep 1

# Step 2: Generate templ files
echo -e "\n${YELLOW}[2/4] Generating templ files...${NC}"
$TEMPL_CMD generate

# Step 3: Build CSS
echo -e "\n${YELLOW}[3/4] Building CSS...${NC}"
npm run css:build

# Step 4: Run E2E tests
echo -e "\n${YELLOW}[4/4] Running E2E tests...${NC}"
echo -e "${BLUE}----------------------------------------${NC}"
npm run test:e2e

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}Quick test completed!${NC}"
echo -e "${GREEN}========================================${NC}"
