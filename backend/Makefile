.PHONY: build dev test test-e2e test-quick test-full clean css templ

# Version info from git
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -X github.com/idtazkia/stmik-admission-api/version.GitCommit=$(GIT_COMMIT) \
           -X github.com/idtazkia/stmik-admission-api/version.GitBranch=$(GIT_BRANCH) \
           -X github.com/idtazkia/stmik-admission-api/version.BuildTime=$(BUILD_TIME)

# Build the server binary
build: templ css
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/server ./cmd/server

# Build for Linux (deployment)
build-linux: templ css
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/server-linux ./cmd/server

# Run development server
dev: templ
	go run -ldflags "$(LDFLAGS)" ./cmd/server

# Run all tests
test:
	go test -v ./...

# Run E2E tests (requires running server)
test-e2e:
	npm run test:e2e

# Quick test: rebuild and run tests (no DB reset)
test-quick:
	./scripts/test-quick.sh

# Full test suite: clean build, reset DB, run all tests
test-full:
	./scripts/test-full.sh

# Generate templ files
templ:
	~/go/bin/templ generate

# Build CSS
css:
	npm run css:build

# Watch CSS changes
css-watch:
	npm run css:watch

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf playwright-report/
	rm -rf test-results/

# Install dependencies
deps:
	go mod download
	npm install
	go install github.com/a-h/templ/cmd/templ@latest

# Run database migrations
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

# Show version info
version:
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Branch: $(GIT_BRANCH)"
	@echo "Build:  $(BUILD_TIME)"
