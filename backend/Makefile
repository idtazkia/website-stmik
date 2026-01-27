.PHONY: build dev test test-e2e clean clean-test css templ

# Version info from git
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -X github.com/idtazkia/stmik-admission-api/internal/version.GitCommit=$(GIT_COMMIT) \
           -X github.com/idtazkia/stmik-admission-api/internal/version.GitBranch=$(GIT_BRANCH) \
           -X github.com/idtazkia/stmik-admission-api/internal/version.BuildTime=$(BUILD_TIME)

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

# Full test suite using testcontainers (like mvn clean test)
# Starts fresh PostgreSQL container, runs migrations, seeds data, runs all tests
clean-test: clean templ css
	go run ./cmd/testrunner

# Generate templ files
templ:
	~/go/bin/templ generate

# Build CSS
css:
	npm run css:build

# Watch CSS changes
css-watch:
	npm run css:watch

# Clean build artifacts and generated files
clean:
	rm -rf bin/
	rm -rf playwright-report/
	rm -rf test-results/
	rm -f web/static/css/output.css
	find web/templates -name "*_templ.go" -delete

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
