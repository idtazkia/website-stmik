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
| Client State | Alpine.js | Dropdowns, modals, form interactions |
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
│   ├── server/main.go           # Entry point
│   └── migrate/main.go          # Migration CLI
├── handler/                      # HTTP handlers
│   ├── router.go                # Routes + middleware
│   ├── auth.go                  # Login, OAuth, logout
│   ├── api.go                   # JSON API (lead capture)
│   ├── portal.go                # Registrant portal pages
│   └── admin.go                 # Staff dashboard + settings
├── model/                        # Data structs + DB queries
│   ├── db.go                    # Connection pool
│   ├── user.go
│   ├── prospect.go
│   ├── application.go
│   ├── document.go
│   └── lookup.go                # Programs, tracks, intakes, etc.
├── migrations/                   # SQL files
├── templates/                    # Templ files (split as needed)
├── static/                       # CSS, JS
├── config.go
├── auth.go                       # JWT + bcrypt
├── kafka.go                      # Payment consumer
├── whatsapp.go                   # WA client
├── go.mod
└── .env.example
```

**Design decisions:**

- No `internal/` - not a library, single executable, no external importers
- No `repository/` layer - model handles its own queries
- No `services/` layer - handlers call models, extract only when reused
- Top-level files for small concerns (config, auth, integrations)
- Split into packages only when files grow large

**References:**
- [Official Go module layout](https://go.dev/doc/modules/layout)
- [Alex Edwards: 11 tips for structuring Go projects](https://www.alexedwards.net/blog/11-tips-for-structuring-your-go-projects)
- [Rost Glukhov: Go Project Structure](https://www.glukhov.org/post/2025/12/go-project-structure/)

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

```sql
-- UUID primary keys (PostgreSQL 13+ has gen_random_uuid() built-in)

-- Users (registrants and staff)
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255),
    provider        VARCHAR(50) NOT NULL DEFAULT 'local',
    provider_id     VARCHAR(255),
    role            VARCHAR(50) NOT NULL DEFAULT 'registrant',
    is_active       BOOLEAN DEFAULT TRUE,
    last_assigned_at TIMESTAMP,  -- For round-robin (staff only)
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Intakes (registration periods)
CREATE TABLE intakes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    year            INTEGER NOT NULL,
    period          VARCHAR(20) NOT NULL,  -- 'ganjil', 'genap'
    open_date       DATE NOT NULL,
    close_date      DATE NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Programs
CREATE TABLE programs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(10) NOT NULL,  -- 'SI', 'TI'
    name            VARCHAR(255) NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE
);

-- Tracks (funding types)
CREATE TABLE tracks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) NOT NULL,
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL,  -- 'private', 'government', 'internal'
    is_active       BOOLEAN DEFAULT TRUE
);

-- Cancel reasons
CREATE TABLE cancel_reasons (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) NOT NULL,
    label           VARCHAR(255) NOT NULL,
    applies_to      VARCHAR(50) NOT NULL,  -- 'prospect', 'application', 'both'
    is_active       BOOLEAN DEFAULT TRUE
);

-- Referrers (people who refer prospects)
CREATE TABLE referrers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255),
    whatsapp        VARCHAR(20),
    type            VARCHAR(50) NOT NULL,  -- 'student', 'alumni', 'partner', 'staff'
    user_id         UUID REFERENCES users(id),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Campaigns (for tracking ad spend and ROI)
CREATE TABLE campaigns (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    utm_campaign    VARCHAR(100) UNIQUE NOT NULL,
    start_date      DATE,
    end_date        DATE,
    budget          DECIMAL(15,2),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Prospects (leads)
CREATE TABLE prospects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    whatsapp        VARCHAR(20) NOT NULL,

    -- Assignment
    intake_id       UUID REFERENCES intakes(id),
    assigned_to     UUID REFERENCES users(id),
    status          VARCHAR(50) DEFAULT 'new',

    -- Referral tracking
    referrer_id     UUID REFERENCES referrers(id),

    -- UTM tracking
    utm_source      VARCHAR(100),
    utm_medium      VARCHAR(100),
    utm_campaign    VARCHAR(100),
    utm_term        VARCHAR(255),
    utm_content     VARCHAR(255),

    -- Additional tracking
    landing_page    VARCHAR(500),
    device_type     VARCHAR(50),

    -- Cancellation
    cancel_reason_id UUID REFERENCES cancel_reasons(id),
    cancel_remarks  TEXT,
    cancelled_at    TIMESTAMP,

    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Applications
CREATE TABLE applications (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prospect_id     UUID REFERENCES prospects(id),
    user_id         UUID REFERENCES users(id),
    intake_id       UUID REFERENCES intakes(id),
    program_id      UUID REFERENCES programs(id),
    track_id        UUID REFERENCES tracks(id),
    status          VARCHAR(50) DEFAULT 'pending_review',
    va_number       VARCHAR(50),
    paid_at         TIMESTAMP,
    cancel_reason_id UUID REFERENCES cancel_reasons(id),
    cancel_remarks  TEXT,
    cancelled_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- Documents
CREATE TABLE documents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id  UUID REFERENCES applications(id),
    doc_type        VARCHAR(50) NOT NULL,  -- 'ktp', 'ijazah'
    filename        VARCHAR(255) NOT NULL,
    filepath        VARCHAR(500) NOT NULL,
    status          VARCHAR(50) DEFAULT 'pending',
    remarks         TEXT,
    reviewed_by     UUID REFERENCES users(id),
    reviewed_at     TIMESTAMP,
    uploaded_at     TIMESTAMP DEFAULT NOW()
);

-- Document checklist templates
CREATE TABLE document_checklists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doc_type        VARCHAR(50) NOT NULL,
    check_item      VARCHAR(255) NOT NULL,
    sort_order      INTEGER DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE
);

-- Document review results
CREATE TABLE document_reviews (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id     UUID REFERENCES documents(id),
    checklist_id    UUID REFERENCES document_checklists(id),
    passed          BOOLEAN NOT NULL,
    reviewed_by     UUID REFERENCES users(id),
    reviewed_at     TIMESTAMP DEFAULT NOW()
);

-- Activity log
CREATE TABLE activity_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID REFERENCES users(id),
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       UUID NOT NULL,
    action          VARCHAR(100) NOT NULL,
    metadata        JSONB,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Communication log
CREATE TABLE communication_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prospect_id     UUID REFERENCES prospects(id),
    channel         VARCHAR(20) NOT NULL,  -- 'whatsapp', 'email'
    template        VARCHAR(100),
    message         TEXT,
    status          VARCHAR(50),  -- 'sent', 'delivered', 'failed'
    sent_by         UUID REFERENCES users(id),
    sent_at         TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_prospects_status ON prospects(status);
CREATE INDEX idx_prospects_intake ON prospects(intake_id);
CREATE INDEX idx_prospects_assigned ON prospects(assigned_to);
CREATE INDEX idx_prospects_referrer ON prospects(referrer_id);
CREATE INDEX idx_prospects_utm_source ON prospects(utm_source);
CREATE INDEX idx_prospects_utm_campaign ON prospects(utm_campaign);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_intake ON applications(intake_id);
CREATE INDEX idx_documents_application ON documents(application_id);
CREATE INDEX idx_activity_entity ON activity_log(entity_type, entity_id);
CREATE INDEX idx_referrers_code ON referrers(code);
CREATE INDEX idx_campaigns_utm ON campaigns(utm_campaign);
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

## Development

```bash
# Install dependencies
go mod download
go install github.com/a-h/templ/cmd/templ@latest

# Setup database
createdb campus
go run ./cmd/migrate up

# Generate templates
templ generate

# Run server
go run ./cmd/server

# With hot reload
air
```

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
github.com/jackc/pgx/v5           # PostgreSQL
github.com/golang-jwt/jwt/v5       # JWT
golang.org/x/crypto                # bcrypt
github.com/a-h/templ               # Templates
github.com/golang-migrate/migrate  # Migrations
github.com/segmentio/kafka-go      # Kafka
```
