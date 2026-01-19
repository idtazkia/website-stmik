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
echo -e "${BLUE}Full Test Suite - Clean Build & Test${NC}"
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
echo -e "\n${YELLOW}[1/7] Stopping any running server...${NC}"
pkill -f "go run ./cmd/server" 2>/dev/null || true
pkill -f "bin/server" 2>/dev/null || true
sleep 1

# Step 2: Reset database (local only - CI uses services block)
if [ -z "$CI" ]; then
    echo -e "\n${YELLOW}[2/7] Resetting database...${NC}"
    CONTAINER_NAME="stmik-postgres"
    VOLUME_NAME="stmiktazkiaacid_postgres_data"

    if docker ps -q -f name="$CONTAINER_NAME" | grep -q .; then
        echo "  Stopping postgres container..."
        docker stop "$CONTAINER_NAME" >/dev/null
    fi

    if docker ps -aq -f name="$CONTAINER_NAME" | grep -q .; then
        echo "  Removing postgres container..."
        docker rm "$CONTAINER_NAME" >/dev/null
    fi

    if docker volume ls -q | grep -q "^${VOLUME_NAME}$"; then
        echo "  Removing postgres volume..."
        docker volume rm "$VOLUME_NAME" >/dev/null
    fi

    echo "  Creating fresh postgres container..."
    docker run -d \
        --name "$CONTAINER_NAME" \
        -e POSTGRES_USER=stmik \
        -e POSTGRES_PASSWORD='WHDnr908ovsvy1aIrBIhNSmFmiNVgbnGpewRqV+tBMQ=' \
        -e POSTGRES_DB=stmik_admission \
        -p 5432:5432 \
        -v "$VOLUME_NAME":/var/lib/postgresql/data \
        postgres:18 >/dev/null

    # Step 3: Wait for postgres
    echo -e "\n${YELLOW}[3/7] Waiting for PostgreSQL to be ready...${NC}"
    MAX_RETRIES=30
    RETRY_COUNT=0
    until docker exec "$CONTAINER_NAME" pg_isready -U stmik -d stmik_admission >/dev/null 2>&1; do
        RETRY_COUNT=$((RETRY_COUNT + 1))
        if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
            echo -e "${RED}ERROR: PostgreSQL did not become ready in time${NC}"
            exit 1
        fi
        echo -n "."
        sleep 1
    done
    echo -e " ${GREEN}Ready${NC}"
else
    echo -e "\n${YELLOW}[2/7] Skipping DB reset (CI uses services block)${NC}"
    echo -e "\n${YELLOW}[3/7] Skipping wait (CI handles health checks)${NC}"
fi

# Step 4: Generate templ files
echo -e "\n${YELLOW}[4/7] Generating templ files...${NC}"
$TEMPL_CMD generate

# Step 5: Build CSS
echo -e "\n${YELLOW}[5/7] Building CSS...${NC}"
npm run css:build

# Step 6: Run migrations
echo -e "\n${YELLOW}[6/7] Running migrations...${NC}"
go run ./cmd/migrate up

# Step 7: Run E2E tests
echo -e "\n${YELLOW}[7/7] Running E2E tests...${NC}"
echo -e "${BLUE}----------------------------------------${NC}"
npm run test:e2e

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}Test suite completed!${NC}"
echo -e "${GREEN}========================================${NC}"
