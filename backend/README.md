# Backend - STMIK Tazkia Admission System

Go-based sales funnel management system for campus admissions.

## Overview

The backend handles the complete lead-to-registration journey:
- Lead capture from Astro landing page
- Portal for prospects to complete applications
- Admin dashboard for marketing staff
- Document review with checklist
- Payment verification via Kafka integration
- WhatsApp notifications via REST API

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Language | Go 1.25+ | Performance, single binary deployment |
| Router | net/http (stdlib) | No external router dependency |
| Database | PostgreSQL 18 + pgx/v5 | Native driver, connection pooling |
| Templates | Templ | Type-safe, compiled HTML templates |
| Interactivity | HTMX | Server-driven UI updates |
| Client State | Alpine.js + alpine-csp | Dropdowns, modals, CSP-compliant |
| Styling | Tailwind CSS | Utility-first, consistent with landing page |
| Auth | golang-jwt/jwt/v5 | JWT tokens |
| Password | x/crypto/bcrypt | Secure password hashing |
| Migrations | golang-migrate | SQL-based migrations |
| Logging | slog (stdlib) | Structured logging |
| Messaging | Kafka (segmentio) | Payment event integration |

## Sales Funnel

### Status Flow

```
PROSPECT                           APPLICATION
────────                           ───────────
   │                                    │
   ▼                                    ▼
┌─────┐    ┌───────────┐    ┌─────────────────┐    ┌──────────┐    ┌──────────┐
│ New │───▶│ Contacted │───▶│ Applicant       │───▶│ Approved │───▶│ Enrolled │
└─────┘    └───────────┘    │ (uploads docs)  │    │ (payment)│    │ (paid)   │
   │            │           └─────────────────┘    └──────────┘    └──────────┘
   │            │                   │                   │
   └────────────┴───────────────────┴───────────────────┘
                                    ▼
                              Cancelled
                         (by registrant)
```

### Prospect Status

| Status | Description |
|--------|-------------|
| `new` | Just submitted form from landing page |
| `contacted` | Staff has made contact |
| `applicant` | Started application (uploaded documents) |
| `cancelled` | Cancelled by registrant |

### Application Status

| Status | Description |
|--------|-------------|
| `pending_review` | Documents uploaded, waiting for review |
| `revision_required` | Documents need revision |
| `approved` | All documents approved, awaiting payment |
| `enrolled` | Payment verified |
| `cancelled` | Cancelled by registrant |

### Document Status

| Status | Description |
|--------|-------------|
| `pending` | Uploaded, not yet reviewed |
| `approved` | Passed all checklist items |
| `revision_required` | Failed checklist, needs re-upload |

## Features

### 1. Lead Management

- Create prospects from landing page API
- List with filters (status, source, intake, assigned staff)
- Round-robin auto-assignment to staff
- Manual assignment/reassignment
- Activity timeline per prospect
- Mark as cancelled with reason

### 2. Application Management

- Select program (SI / TI)
- Select track (Regular, KIP-K, LPDP, Internal Scholarship)
- Select intake period (Ganjil / Genap)
- Upload documents (KTP, Ijazah)
- Document review with checklist
- Approve application (when all docs approved)
- Cancel with reason

### 3. Document Review

Each document type has a checklist:

**KTP Checklist:**
- Nama jelas terbaca
- NIK lengkap (16 digit)
- Foto tidak buram
- Masih berlaku

**Ijazah Checklist:**
- Nama sesuai KTP
- Tahun lulus terlihat
- Stempel sekolah ada
- Tanda tangan kepala sekolah

Staff marks each item pass/fail. If any fail, document status becomes `revision_required` with remarks.

### 4. Intake Management

- Create intake periods (e.g., "2025 Ganjil", "2025 Genap")
- Set registration open/close dates
- 2 intakes per year
- No quota limits

### 5. Staff Management

- Round-robin lead assignment
- Set staff active/inactive for assignment
- View assigned leads count
- Manual lead reassignment

### 6. Communication

**WhatsApp (REST API):**
- Welcome message (new prospect)
- Follow-up reminder (3 days no activity)
- Document reminder (7 days incomplete)
- Revision request (with remarks)
- Application approved (with VA number)
- Payment confirmed

**Email (SMTP):**
- Same templates as WhatsApp
- Formal acceptance letter (PDF)

### 7. Payment Integration (Kafka)

```
Incoming: payment.completed
├── va_number
├── amount
├── paid_at
└── → Update application: approved → enrolled

Outgoing: application.approved (optional)
├── application_id
├── prospect_name
└── → Trigger VA generation in payment app
```

### 8. Referral Tracking

- Generate unique referral codes for referrers
- Referrer types: student, alumni, partner, staff
- Track referral source on prospect registration
- Referral conversion stats (referred → enrolled)
- Referrer leaderboard
- Referral rewards tracking (optional)

**Referral Flow:**
```
Referrer gets code → Shares link → Prospect registers with code
                                          │
                                          ▼
                              Prospect linked to referrer
                                          │
                                          ▼
                              Track conversion to enrolled
```

### 9. Ad Campaign Tracking

Track marketing campaign effectiveness via UTM parameters:

| Parameter | Purpose | Example |
|-----------|---------|---------|
| `utm_source` | Traffic source | google, facebook, instagram |
| `utm_medium` | Marketing medium | cpc, social, email, banner |
| `utm_campaign` | Campaign name | intake_2025_ganjil |
| `utm_term` | Paid keywords | kuliah_it_jakarta |
| `utm_content` | Ad variation | banner_v1, video_ad |

**Additional Tracking:**
- Landing page URL
- First touch attribution (first campaign that brought them)
- Device type (mobile/desktop)
- Registration timestamp

**Reports:**
- Conversion by source (Google vs Facebook vs Instagram)
- Conversion by campaign
- Conversion by medium (paid vs organic)
- Cost per acquisition (if cost data provided)
- ROI per campaign

### 10. Reports & Dashboard

- Funnel overview (counts per stage)
- This intake vs previous comparison
- Leads by source/campaign
- Leads by program/track
- Leads by referrer
- Staff conversion leaderboard
- Campaign performance comparison
- Referrer leaderboard
- Export to CSV

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Browser                               │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┴─────────────────────┐
        ▼                                           ▼
┌───────────────────┐                    ┌───────────────────┐
│  Cloudflare Pages │                    │  Cloudflare CDN   │
│  (Astro Landing)  │                    │  (Proxy to VPS)   │
└───────────────────┘                    └───────────────────┘
                                                   │
                                                   ▼
                                         ┌───────────────────┐
                                         │      Nginx        │
                                         │  + Rate Limiting  │
                                         │  + SSL            │
                                         └───────────────────┘
                                                   │
                              ┌────────────────────┼────────────────────┐
                              ▼                    ▼                    ▼
                    ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
                    │  Go Backend  │     │    Kafka     │     │  WhatsApp    │
                    │              │◀───▶│   (Payment)  │     │  REST API    │
                    └──────────────┘     └──────────────┘     └──────────────┘
                              │
                              ▼
                    ┌──────────────┐
                    │ PostgreSQL   │
                    └──────────────┘
```

## Project Structure

```
backend/
├── cmd/
│   ├── server/          # Main application entry point
│   ├── migrate/         # Database migration CLI
│   ├── seedtest/        # Test data seeder
│   ├── mockup/          # Screenshot generator
│   └── testrunner/      # Test orchestrator (like mvn clean test)
├── internal/            # Private application code
│   ├── auth/            # Session management, Google OAuth
│   ├── config/          # Configuration loading
│   ├── handler/         # HTTP handlers (admin, portal, public)
│   ├── integration/     # External services (Resend, WhatsApp)
│   ├── model/           # Data models + database queries
│   ├── mockdata/        # Mock data generators
│   ├── pkg/crypto/      # Encryption utilities
│   ├── storage/         # File storage abstraction
│   └── version/         # Build version info
├── web/                 # Web assets
│   ├── static/          # CSS, JS (Tailwind, HTMX, Alpine)
│   └── templates/       # Templ files (admin, portal, email)
├── test/                # Test files
│   ├── e2e/             # Playwright E2E tests
│   └── testutil/        # Test utilities (testcontainers)
├── migrations/          # SQL migration files
├── Makefile             # Build commands
└── go.mod
```

**Design decisions:**

- Standard Go layout with `internal/` for private packages
- No `repository/` layer - model handles its own queries
- No `services/` layer - handlers call models, extract only when reused
- `web/` for all frontend assets (templates, static files)
- `test/` for all test-related code (e2e, utilities)

**References:**
- [Official Go module layout](https://go.dev/doc/modules/layout)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Routes

### Public API

```
POST /api/prospects              # Create prospect (from landing page)
GET  /api/health                 # Health check
GET  /api/referrers/{code}       # Validate referral code (for landing page)
```

**POST /api/prospects** payload:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "whatsapp": "08123456789",
  "intake_id": 1,
  "referral_code": "REF123",
  "utm_source": "google",
  "utm_medium": "cpc",
  "utm_campaign": "intake_2025_ganjil",
  "utm_term": "kuliah it jakarta",
  "utm_content": "banner_v1",
  "landing_page": "https://stmik.tazkia.ac.id/promo",
  "device_type": "mobile"
}
```

### Portal (Registrants)

```
GET  /portal/login               # Login page
POST /portal/login               # Login submit
GET  /portal/register            # Register page
POST /portal/register            # Register submit
GET  /portal/auth/google         # Google OAuth
GET  /portal/auth/google/callback
POST /portal/logout

GET  /portal                     # Dashboard (status overview)
GET  /portal/application         # Application form
POST /portal/application         # Create/update application (HTMX)
GET  /portal/documents           # Document upload page
POST /portal/documents           # Upload document (HTMX)
DELETE /portal/documents/{id}    # Remove document (HTMX)

GET  /portal/cancel              # Cancel confirmation page
POST /portal/cancel              # Submit cancellation
```

### Admin (Staff)

```
GET  /admin/login                # Login (Google only)
GET  /admin/auth/google
GET  /admin/auth/google/callback
POST /admin/logout

GET  /admin                      # Dashboard
GET  /admin/prospects            # Prospect list
GET  /admin/prospects/{id}       # Prospect detail
POST /admin/prospects/{id}/assign        # Assign to staff (HTMX)
POST /admin/prospects/{id}/status        # Update status (HTMX)
POST /admin/prospects/{id}/whatsapp      # Send WhatsApp (HTMX)
POST /admin/prospects/{id}/cancel        # Mark cancelled (HTMX)

GET  /admin/applications         # Application list
GET  /admin/applications/{id}    # Application detail
GET  /admin/applications/{id}/documents/{docId}/review  # Review modal
POST /admin/applications/{id}/documents/{docId}/review  # Submit review (HTMX)
POST /admin/applications/{id}/approve    # Approve application (HTMX)
POST /admin/applications/{id}/cancel     # Mark cancelled (HTMX)

GET  /admin/referrers            # Referrer list
GET  /admin/referrers/{id}       # Referrer detail + stats
POST /admin/referrers            # Create referrer (HTMX)
PUT  /admin/referrers/{id}       # Update referrer (HTMX)
POST /admin/referrers/{id}/toggle        # Toggle active (HTMX)

GET  /admin/campaigns            # Campaign list
GET  /admin/campaigns/{id}       # Campaign detail + stats
POST /admin/campaigns            # Create campaign (HTMX)
PUT  /admin/campaigns/{id}       # Update campaign (HTMX)

GET  /admin/settings             # Settings overview
GET  /admin/settings/intakes     # Intake management
POST /admin/settings/intakes     # Create intake (HTMX)
PUT  /admin/settings/intakes/{id}        # Update intake (HTMX)
GET  /admin/settings/tracks      # Track management
GET  /admin/settings/reasons     # Cancel reasons
GET  /admin/settings/checklists  # Document checklists
GET  /admin/settings/staff       # Staff management
POST /admin/settings/staff/{id}/toggle   # Toggle active (HTMX)

GET  /admin/reports              # Reports page
GET  /admin/reports/funnel       # Funnel data (HTMX)
GET  /admin/reports/sources      # Conversion by source (HTMX)
GET  /admin/reports/campaigns    # Campaign performance (HTMX)
GET  /admin/reports/referrers    # Referrer leaderboard (HTMX)
GET  /admin/reports/export       # CSV export
```

## Database Schema

See `migrations/` folder for SQL migration files. Run migrations with:

```bash
go run ./cmd/migrate up
```

## Cancel Reasons (Seed Data)

| Code | Label | Applies To |
|------|-------|------------|
| `no_response` | Tidak merespon | prospect |
| `chose_other` | Memilih kampus lain | both |
| `financial` | Kendala biaya | both |
| `changed_mind` | Berubah pikiran | both |
| `age_requirement` | Tidak memenuhi syarat usia | prospect |
| `invalid_document` | Dokumen tidak valid | application |
| `other` | Lainnya | both |

## Document Checklists (Seed Data)

**KTP:**
| Check Item |
|------------|
| Nama jelas terbaca |
| NIK lengkap (16 digit) |
| Foto tidak buram |
| Masih berlaku |

**Ijazah:**
| Check Item |
|------------|
| Nama sesuai KTP |
| Tahun lulus terlihat |
| Stempel sekolah ada |
| Tanda tangan kepala sekolah |

## WhatsApp Templates

| Template | Trigger | Variables |
|----------|---------|-----------|
| `welcome` | New prospect | `name` |
| `followup_3d` | 3 days no activity | `name` |
| `document_reminder` | 7 days incomplete | `name`, `missing_docs` |
| `revision_request` | Doc needs revision | `name`, `doc_type`, `remarks` |
| `approved` | Application approved | `name`, `program`, `va_number` |
| `enrolled` | Payment confirmed | `name`, `program` |

## Configuration

```bash
# Server
PORT=3000
APP_URL=https://yourdomain.com

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/campus?sslmode=disable

# Authentication
JWT_SECRET=your-32-char-minimum-secret
JWT_EXPIRY=168h
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx
STAFF_EMAIL_DOMAIN=tazkia.ac.id

# File uploads
UPLOAD_DIR=/var/www/uploads
MAX_FILE_SIZE=5242880

# WhatsApp API
WHATSAPP_API_URL=https://api.whatsapp.example.com
WHATSAPP_API_KEY=xxx

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CONSUMER_GROUP=campus-backend
KAFKA_PAYMENT_TOPIC=payment.completed
```

## Alpine.js CSP Guidelines

This project uses **Alpine.js CSP build** (`@alpinejs/csp`) for Content Security Policy compliance. This is required for PCI-DSS 4.0 compliance (effective April 2025).

### Key Differences from Standard Alpine.js

1. **No inline expressions** - You cannot write JavaScript directly in HTML attributes
2. **Use `Alpine.data()`** - All component logic must be registered via `Alpine.data()`
3. **Reference by name only** - HTML attributes reference methods/properties by name

### ❌ DON'T (Standard Alpine.js)

```html
<div x-data="{ open: false }">
    <button @click="open = !open">Toggle</button>
    <div x-show="open">Content</div>
</div>
```

### ✅ DO (Alpine.js CSP)

```html
<div x-data="dropdown">
    <button @click="toggle">Toggle</button>
    <div x-show="open">Content</div>
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

### Reference

- [Alpine.js CSP Documentation](https://alpinejs.dev/advanced/csp)

## Local Installation for Manual Testing

### Prerequisites

- **Go 1.25+** - [Download Go](https://go.dev/dl/)
- **Docker & Docker Compose** - [Install Docker Desktop](https://www.docker.com/products/docker-desktop/) (recommended for database)
- **Node.js 20+** - [Download Node.js](https://nodejs.org/)
- **Make** - Usually pre-installed on macOS/Linux

**Alternative:** Install PostgreSQL 18 locally instead of using Docker (see Step 2, Option B).

### Step 1: Install Dependencies

```bash
# Install Go modules, npm packages, and templ CLI
make deps
```

### Step 2: Setup PostgreSQL Database

#### Option A: Using Docker Compose (Recommended)

```bash
# Start PostgreSQL container
docker compose up -d

# Check container is running
docker compose ps

# View logs (optional)
docker compose logs -f postgres

# Test connection (requires psql client)
psql -U stmik -d stmik_admission -h localhost

# Stop container (when done)
docker compose down

# Stop and remove data volume (reset database)
docker compose down -v
```

**Benefits:**
- No PostgreSQL installation required
- Isolated environment
- Easy to reset database
- Consistent across all developers

#### Option B: Using Local PostgreSQL

```bash
# Start PostgreSQL service
# macOS (Homebrew):
brew services start postgresql@18

# Linux (systemd):
sudo systemctl start postgresql

# Create database and user
psql postgres -c "CREATE USER stmik WITH PASSWORD 'stmik_dev_password';"
psql postgres -c "CREATE DATABASE stmik_admission OWNER stmik;"
psql postgres -c "GRANT ALL PRIVILEGES ON DATABASE stmik_admission TO stmik;"

# Test connection
psql -U stmik -d stmik_admission -h localhost
```

### Step 3: Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your settings
# Required:
# - DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME
# - JWT_SECRET (generate a random 32+ character string)
# - ENCRYPTION_KEY (generate with: openssl rand -hex 32)
#
# Optional for local testing:
# - Google OAuth (can skip for basic testing)
# - Resend, WhatsApp, Kafka (optional integrations)
```

**Generate required secrets:**

```bash
# Generate JWT secret (32+ characters)
openssl rand -base64 32

# Generate encryption key (32-byte hex)
openssl rand -hex 32
```

### Step 4: Run Database Migrations

```bash
# Apply all migrations
go run ./cmd/migrate up

# Verify tables created
psql -U stmik -d stmik_admission -c "\dt"
```

### Step 5: Seed Test Data (Optional)

```bash
# Seed test data for manual testing
go run ./cmd/seedtest

# This creates:
# - Test users for each role (admin, supervisor, consultant, finance, academic)
# - Test candidates
# - Sample applications and documents
```

### Step 6: Generate Templates and CSS

```bash
# Generate Templ templates
make templ

# Build Tailwind CSS
make css
```

### Step 7: Run Development Server

```bash
# Start server
make dev

# Server will start at http://localhost:8080
```

### Step 8: Access the Application

**Public Endpoints:**
- Health Check: http://localhost:8080/health
- Portal Login: http://localhost:8080/portal/login
- Portal Register: http://localhost:8080/portal/register

**Test Endpoints (Development Only):**
- Test Admin: http://localhost:8080/test/admin
- Test Portal: http://localhost:8080/test/portal
- Auto-login as Admin: http://localhost:8080/test/login/admin
- Auto-login as Candidate: http://localhost:8080/test/login/candidate

**Admin Panel:**
- Admin Login: http://localhost:8080/admin/login
- Admin Dashboard: http://localhost:8080/admin

### Quick Start (All Commands)

```bash
# From backend directory
make deps
docker compose up -d
cp .env.example .env
# Edit .env and set required values (JWT_SECRET, ENCRYPTION_KEY)
go run ./cmd/migrate up
go run ./cmd/seedtest  # Optional
make dev
```

### Troubleshooting

**Docker Compose issues:**
```bash
# Check if container is running
docker compose ps

# View container logs
docker compose logs postgres

# Restart container
docker compose restart

# Reset database completely
docker compose down -v
docker compose up -d
```

**Database connection error:**
```bash
# Check PostgreSQL is running
docker compose ps

# Or if using local PostgreSQL
psql -U stmik -d stmik_admission -h localhost

# Verify .env has correct database credentials
cat .env | grep DATABASE
```

**Port 8080 already in use:**
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
SERVER_PORT=8081
```

**Templates not loading:**
```bash
# Regenerate templates
make templ

# Check generated files
ls -la web/templates/*_templ.go
```

**CSS not loading:**
```bash
# Rebuild CSS
make css

# Check output
ls -la web/static/css/output.css
```

### Development Workflow

```bash
# Start database (terminal 1)
docker compose up -d

# Watch for template changes (terminal 2)
watch make templ

# Watch for CSS changes (terminal 3)
make css-watch

# Run server (terminal 4)
make dev

# When done, stop database
docker compose down
```

### Common Docker Compose Commands

```bash
# Start database in background
docker compose up -d

# View database logs
docker compose logs -f postgres

# Stop database
docker compose down

# Stop and remove data (reset)
docker compose down -v

# Check database status
docker compose ps

# Connect to database CLI
docker compose exec postgres psql -U stmik -d stmik_admission
```

## Development

```bash
# Install all dependencies (Go, npm, templ)
make deps

# Run development server
make dev

# Generate templates
make templ

# Build CSS
make css
```

## Testing

```bash
# Full test suite (like mvn clean test)
# - Cleans generated files
# - Starts PostgreSQL testcontainer
# - Runs migrations
# - Seeds test data
# - Runs unit tests and E2E tests
make clean-test

# Run unit tests only
make test

# Run E2E tests only (requires running server)
make test-e2e

# Clean build artifacts
make clean
```

The `make clean-test` command uses [testcontainers](https://testcontainers.com/) to automatically manage a PostgreSQL container, similar to how Maven's `clean test` works in Java projects.

## Deployment

```bash
# Build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o campus-api ./cmd/server

# Deploy
scp campus-api user@vps:~/
ssh user@vps sudo systemctl restart campus-api
```

## Dependencies

```
github.com/jackc/pgx/v5              # PostgreSQL driver
github.com/golang-jwt/jwt/v5         # JWT tokens
golang.org/x/crypto                  # bcrypt password hashing
github.com/a-h/templ                 # Type-safe templates
github.com/golang-migrate/migrate    # Database migrations
github.com/segmentio/kafka-go        # Kafka integration
github.com/testcontainers/testcontainers-go  # Test containers
```
