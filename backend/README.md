# Backend - STMIK Tazkia Admission System

Go-based sales funnel management system for campus admissions.

## Overview

Complete lead-to-enrollment journey management:
- Lead capture from landing page
- Portal for prospects to complete applications
- Admin CRM dashboard for marketing staff
- Document review with checklists
- Payment verification
- WhatsApp notifications

## Quick Start

```bash
# 1. Install dependencies
make deps

# 2. Start PostgreSQL with Docker Compose
docker compose up -d

# 3. Configure environment
cp .env.example .env
# Edit .env (set JWT_SECRET and ENCRYPTION_KEY)

# 4. Run migrations
go run ./cmd/migrate up

# 5. Seed test data (optional)
go run ./cmd/seedtest

# 6. Start server
make dev
```

Server runs at http://localhost:8080

**Test Endpoints:**
- Health: http://localhost:8080/health
- Admin Auto-login: http://localhost:8080/test/login/admin
- Candidate Auto-login: http://localhost:8080/test/login/candidate

## Tech Stack

| Layer | Technology |
|-------|------------|
| Language | Go 1.25+ |
| Router | net/http (stdlib) |
| Database | PostgreSQL 18 + pgx/v5 |
| Templates | Templ |
| Interactivity | HTMX + Alpine.js CSP |
| Styling | Tailwind CSS |
| Auth | golang-jwt/jwt/v5 |
| Migrations | golang-migrate |
| Testing | Testcontainers + Playwright |

## Documentation

| Document | Purpose |
|----------|---------|
| **[FEATURES.md](FEATURES.md)** | Complete feature list with implementation status (70% complete) |
| **[API.md](API.md)** | REST API and web routes documentation |
| **[Local Installation](#local-installation)** | Step-by-step setup guide |

## Project Structure

```
backend/
├── cmd/                  # Application entry points
│   ├── server/          # Main server
│   ├── migrate/         # Database migrations
│   ├── seedtest/        # Test data seeder
│   └── testrunner/      # Test orchestrator
├── internal/            # Private application code
│   ├── auth/           # Session management, OAuth
│   ├── config/         # Configuration
│   ├── handler/        # HTTP handlers
│   ├── integration/    # External services
│   ├── model/          # Data models + queries
│   └── storage/        # File storage
├── web/                 # Web assets
│   ├── static/         # CSS, JS
│   └── templates/      # Templ files
├── test/                # Test files
│   ├── e2e/            # Playwright tests
│   └── testutil/       # Test utilities
└── migrations/          # SQL migrations
```

## Local Installation

### Prerequisites

- **Go 1.25+** - [Download](https://go.dev/dl/)
- **Docker Desktop** - [Install](https://www.docker.com/products/docker-desktop/)
- **Node.js 20+** - [Download](https://nodejs.org/)
- **Make** - Pre-installed on macOS/Linux

### Step-by-Step Setup

#### 1. Install Dependencies

```bash
make deps
```

Installs Go modules, npm packages, and templ CLI.

#### 2. Setup PostgreSQL Database

**Option A: Docker Compose (Recommended)**

```bash
docker compose up -d
docker compose ps
```

**Option B: Local PostgreSQL**

```bash
brew services start postgresql@18  # macOS
sudo systemctl start postgresql     # Linux

psql postgres -c "CREATE USER stmik WITH PASSWORD 'stmik_dev_password';"
psql postgres -c "CREATE DATABASE stmik_admission OWNER stmik;"
```

#### 3. Configure Environment

```bash
cp .env.example .env
```

**Required in .env:**
- `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_NAME`
- `JWT_SECRET` - Generate: `openssl rand -base64 32`
- `ENCRYPTION_KEY` - Generate: `openssl rand -hex 32`

**Optional:**
- `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` - For OAuth
- `RESEND_API_KEY` - For email notifications
- `WHATSAPP_API_URL`, `WHATSAPP_API_TOKEN` - For WhatsApp

#### 4. Run Database Migrations

```bash
go run ./cmd/migrate up
```

#### 5. Seed Test Data (Optional)

```bash
go run ./cmd/seedtest
```

Creates test users and candidates.

#### 6. Generate Templates and CSS

```bash
make templ
make css
```

#### 7. Run Development Server

```bash
make dev
```

Server starts at http://localhost:8080

### Access the Application

**Public Endpoints:**
- Health Check: http://localhost:8080/health
- Portal Login: http://localhost:8080/portal/login
- Portal Register: http://localhost:8080/portal/register

**Test Endpoints (Development Only):**
- Test Admin: http://localhost:8080/test/admin
- Test Portal: http://localhost:8080/test/portal
- Auto-login Admin: http://localhost:8080/test/login/admin
- Auto-login Candidate: http://localhost:8080/test/login/candidate

**Admin Panel:**
- Admin Login: http://localhost:8080/admin/login
- Admin Dashboard: http://localhost:8080/admin

### Troubleshooting

**Docker Compose:**
```bash
docker compose ps              # Check status
docker compose logs postgres   # View logs
docker compose restart         # Restart
docker compose down -v         # Reset database
```

**Database Connection:**
```bash
docker compose ps
psql -U stmik -d stmik_admission -h localhost
cat .env | grep DATABASE
```

**Port 8080 in use:**
```bash
lsof -i :8080
kill -9 <PID>
# Or change port in .env: SERVER_PORT=8081
```

**Templates/CSS:**
```bash
make templ
ls -la web/templates/*_templ.go

make css
ls -la web/static/css/output.css
```

### Development Workflow

```bash
# Terminal 1: Database
docker compose up -d

# Terminal 2: Watch templates
watch make templ

# Terminal 3: Watch CSS
make css-watch

# Terminal 4: Run server
make dev

# When done
docker compose down
```

### Docker Compose Commands

```bash
docker compose up -d              # Start in background
docker compose logs -f postgres   # View logs
docker compose down               # Stop
docker compose down -v            # Stop and remove data
docker compose ps                 # Status
docker compose exec postgres psql -U stmik -d stmik_admission  # CLI
```

## Development Commands

```bash
make deps         # Install dependencies
make dev          # Run development server
make templ        # Generate templates
make css          # Build CSS
make build        # Build for current platform
make build-linux  # Build for Linux deployment
```

## Testing

```bash
make clean-test   # Full test suite (like mvn clean test)
make test         # Unit tests only
make test-e2e     # E2E tests (requires running server)
make clean        # Clean build artifacts
```

Uses testcontainers for PostgreSQL isolation.

## Deployment

```bash
# Build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o campus-api ./cmd/server

# Deploy
scp campus-api user@vps:~/
ssh user@vps sudo systemctl restart campus-api
```

## Architecture

```
Browser
  ├── Cloudflare Pages (Astro Landing) - FREE
  ├── Cloudflare CDN (DDoS Protection) - FREE
  └── VPS ($5/month)
        ├── Nginx (Reverse Proxy + SSL)
        ├── Go Backend (REST API)
        └── PostgreSQL (Database)
```

**Key Features:**
- $5/month fixed cost
- Handles 100,000+ leads
- <1ms database latency
- DDoS protected via Cloudflare

## Alpine.js CSP Guidelines

This project uses Alpine.js CSP build for Content Security Policy compliance.

**❌ DON'T:**
```html
<div x-data="{ open: false }">
    <button @click="open = !open">Toggle</button>
</div>
```

**✅ DO:**
```html
<div x-data="dropdown">
    <button @click="toggle">Toggle</button>
</div>

<script>
document.addEventListener('alpine:init', () => {
    Alpine.data('dropdown', () => ({
        open: false,
        toggle() { this.open = !this.open }
    }))
})
</script>
```

Reference: [Alpine.js CSP Documentation](https://alpinejs.dev/advanced/csp)

## Dependencies

```
github.com/jackc/pgx/v5              # PostgreSQL driver
github.com/golang-jwt/jwt/v5         # JWT tokens
golang.org/x/crypto                  # bcrypt
github.com/a-h/templ                 # Templates
github.com/golang-migrate/migrate    # Migrations
github.com/testcontainers/testcontainers-go  # Testing
```

## Status

**Implementation Progress:** 70% (29/42 features)

**Completed Phases:**
- ✅ Phase 0: UI Mockup
- ✅ Phase 1: Admin Foundation
- ✅ Phase 2: Configuration
- ✅ Phase 4: CRM Operations
- ✅ Phase 6: Commissions

**In Progress:**
- 🔄 Phase 3: Registration & Portal (87%)
- 🔄 Phase 5: Commitment & Enrollment (80%)
- 🔄 Phase 7: Reporting (14%)

**Pending:**
- ⏳ Phase 8: Notifications

See [FEATURES.md](FEATURES.md) for detailed progress tracking.
