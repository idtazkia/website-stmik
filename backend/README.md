# Backend - STMIK Tazkia Admission System

Go-based sales funnel management system for campus admissions.

## Overview

The backend handles the complete lead-to-registration journey:
- Lead capture from Astro landing page
- Portal for leads to complete applications
- Admin dashboard for marketing staff
- Application review and approval workflow

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Language | Go 1.25+ | Performance, single binary deployment |
| Router | net/http (stdlib) | No external router dependency |
| Database | PostgreSQL 16 + pgx/v5 | Native driver, connection pooling |
| Templates | Templ | Type-safe, compiled HTML templates |
| Interactivity | HTMX | Server-driven UI updates |
| Client State | Alpine.js | Dropdowns, modals, form interactions |
| Styling | Tailwind CSS | Utility-first, consistent with landing page |
| Auth | golang-jwt/jwt/v5 | JWT tokens |
| Password | x/crypto/bcrypt | Secure password hashing |
| Migrations | golang-migrate | SQL-based migrations |
| Logging | slog (stdlib) | Structured logging |

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
                                         │  (Reverse Proxy)  │
                                         │  + Rate Limiting  │
                                         │  + SSL Termination│
                                         └───────────────────┘
                                                   │
                                                   ▼
                                         ┌───────────────────┐
                                         │    Go Backend     │
                                         │  ┌─────────────┐  │
                                         │  │   Handlers  │  │
                                         │  ├─────────────┤  │
                                         │  │  Services   │  │
                                         │  ├─────────────┤  │
                                         │  │ Repository  │  │
                                         │  └─────────────┘  │
                                         └───────────────────┘
                                                   │
                                                   ▼
                                         ┌───────────────────┐
                                         │   PostgreSQL 16   │
                                         └───────────────────┘
```

## User Flows

### Lead Journey
```
Astro Landing Page
       │
       ▼ (Submit name + email)
Go Backend: Create Lead
       │
       ▼ (Redirect)
Portal: Complete Profile
       │
       ▼
Portal: Fill Application Form
       │
       ▼
Portal: Upload Documents
       │
       ▼
Portal: Track Status
```

### Staff Journey
```
Admin Login (Google OIDC)
       │
       ▼
Dashboard: View Pipeline
       │
       ├── Filter/Search Leads
       ├── Review Applications
       ├── Approve/Reject
       └── Send Notifications
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── database/
│   │   └── database.go          # PostgreSQL connection
│   ├── handlers/
│   │   ├── api.go               # JSON API handlers (lead capture)
│   │   ├── portal.go            # Lead portal pages
│   │   ├── admin.go             # Admin dashboard pages
│   │   └── auth.go              # Authentication handlers
│   ├── middleware/
│   │   ├── auth.go              # JWT verification
│   │   ├── logging.go           # Request logging
│   │   └── recovery.go          # Panic recovery
│   ├── models/
│   │   ├── user.go              # User model
│   │   ├── lead.go              # Lead model
│   │   └── application.go       # Application model
│   ├── repository/
│   │   ├── user.go              # User queries
│   │   ├── lead.go              # Lead queries
│   │   └── application.go       # Application queries
│   └── services/
│       ├── auth.go              # Authentication logic
│       ├── lead.go              # Lead management
│       └── application.go       # Application processing
├── migrations/
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   └── ...
├── templates/
│   ├── components/
│   │   ├── button.templ
│   │   ├── card.templ
│   │   ├── form.templ
│   │   ├── modal.templ
│   │   └── table.templ
│   ├── layouts/
│   │   ├── base.templ           # HTML base layout
│   │   ├── portal.templ         # Portal layout
│   │   └── admin.templ          # Admin layout
│   ├── pages/
│   │   ├── portal/
│   │   │   ├── login.templ
│   │   │   ├── profile.templ
│   │   │   ├── application.templ
│   │   │   └── status.templ
│   │   └── admin/
│   │       ├── dashboard.templ
│   │       ├── leads.templ
│   │       ├── lead_detail.templ
│   │       └── settings.templ
│   └── emails/
│       ├── welcome.templ
│       └── status_update.templ
├── static/
│   ├── css/
│   │   └── app.css              # Tailwind output
│   └── js/
│       ├── htmx.min.js
│       └── alpine.min.js
├── go.mod
├── go.sum
├── .env.example
├── README.md
└── TODO.md
```

## Request Flow

### HTML Pages (Portal/Admin)
```
Request → Middleware Chain → Handler → Service → Repository → DB
                                           │
                                           ▼
                                    Templ Template
                                           │
                                           ▼
                                    HTML Response
```

### HTMX Partial Updates
```
HTMX Request (HX-Request header)
       │
       ▼
Handler detects HTMX
       │
       ▼
Returns HTML fragment (not full page)
       │
       ▼
HTMX swaps DOM element
```

### JSON API (Lead Capture)
```
POST /api/leads (from Astro)
       │
       ▼
Handler → Service → Repository → DB
       │
       ▼
JSON Response
```

## Routes

### Public API
```
POST /api/leads                  # Create lead (from landing page)
GET  /api/health                 # Health check
```

### Portal (Leads)
```
GET  /portal/login               # Login page
POST /portal/login               # Login submit
GET  /portal/register            # Register page
POST /portal/register            # Register submit
GET  /portal/auth/google         # Google OAuth initiate
GET  /portal/auth/google/callback

GET  /portal/profile             # Profile form
POST /portal/profile             # Update profile (HTMX)

GET  /portal/application         # Application form
POST /portal/application         # Submit application (HTMX)

GET  /portal/documents           # Document upload
POST /portal/documents           # Upload file (HTMX)

GET  /portal/status              # Application status
```

### Admin (Staff)
```
GET  /admin/login                # Staff login (Google only)
GET  /admin/auth/google
GET  /admin/auth/google/callback

GET  /admin/dashboard            # Dashboard with stats
GET  /admin/leads                # Lead list with filters
GET  /admin/leads/:id            # Lead detail
POST /admin/leads/:id/status     # Update status (HTMX)
POST /admin/leads/:id/notes      # Add note (HTMX)

GET  /admin/applications         # Application list
GET  /admin/applications/:id     # Application detail
POST /admin/applications/:id/approve
POST /admin/applications/:id/reject

GET  /admin/settings             # Admin settings
```

## Middleware Stack

```go
// Applied to all routes
mux.Handle("/",
    recovery(
        logging(
            securityHeaders(
                handler,
            ),
        ),
    ),
)

// Applied to portal routes
portalHandler = requireAuth(portalMux)

// Applied to admin routes
adminHandler = requireAuth(requireRole("staff", adminMux))
```

## Database Schema

### Core Tables
```sql
-- Users (both leads and staff)
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    provider      VARCHAR(50) NOT NULL DEFAULT 'local',
    provider_id   VARCHAR(255),
    role          VARCHAR(50) NOT NULL DEFAULT 'lead',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Lead profiles (extended info)
CREATE TABLE lead_profiles (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER UNIQUE REFERENCES users(id),
    phone         VARCHAR(20),
    address       TEXT,
    birth_date    DATE,
    high_school   VARCHAR(255),
    graduation_year INTEGER,
    source        VARCHAR(100),  -- How they found us
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Applications
CREATE TABLE applications (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER REFERENCES users(id),
    program       VARCHAR(100) NOT NULL,
    status        VARCHAR(50) DEFAULT 'draft',
    submitted_at  TIMESTAMP,
    reviewed_by   INTEGER REFERENCES users(id),
    reviewed_at   TIMESTAMP,
    notes         TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Documents
CREATE TABLE documents (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER REFERENCES users(id),
    application_id INTEGER REFERENCES applications(id),
    doc_type      VARCHAR(50) NOT NULL,
    filename      VARCHAR(255) NOT NULL,
    filepath      VARCHAR(500) NOT NULL,
    file_size     INTEGER,
    uploaded_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Activity log
CREATE TABLE activity_log (
    id            SERIAL PRIMARY KEY,
    user_id       INTEGER REFERENCES users(id),
    action        VARCHAR(100) NOT NULL,
    entity_type   VARCHAR(50),
    entity_id     INTEGER,
    metadata      JSONB,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Status Values
```
Lead Status:    new → contacted → qualified → converted → lost
Application:    draft → submitted → under_review → approved → rejected
```

## Configuration

Environment variables (`.env`):

```bash
# Server
PORT=3000
APP_URL=https://yourdomain.com

# Database
DATABASE_URL=postgres://campus_app:password@localhost:5432/campus?sslmode=disable

# Authentication
JWT_SECRET=your-32-char-minimum-secret
JWT_EXPIRY=168h  # 7 days

# Google OAuth
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx
GOOGLE_REDIRECT_URL=https://yourdomain.com/portal/auth/google/callback
ADMIN_GOOGLE_REDIRECT_URL=https://yourdomain.com/admin/auth/google/callback

# Staff domain restriction
STAFF_EMAIL_DOMAIN=tazkia.ac.id

# File uploads
UPLOAD_DIR=/var/www/uploads
MAX_FILE_SIZE=5242880  # 5MB
```

## Development

### Prerequisites
- Go 1.25+
- PostgreSQL 16+
- Node.js (for Tailwind CSS build)

### Setup
```bash
# Clone and enter directory
cd backend

# Install Go dependencies
go mod download

# Install Templ
go install github.com/a-h/templ/cmd/templ@latest

# Install Tailwind (via npm in static/)
cd static && npm install && cd ..

# Copy environment file
cp .env.example .env
# Edit .env with your values

# Run database migrations
go run ./cmd/migrate up

# Generate Templ files
templ generate

# Build Tailwind CSS
cd static && npm run build && cd ..

# Run development server
go run ./cmd/server
```

### Development Commands
```bash
# Run server with hot reload (using air)
air

# Generate Templ templates
templ generate

# Watch Templ files
templ generate --watch

# Build Tailwind CSS
cd static && npm run build

# Watch Tailwind
cd static && npm run watch

# Run tests
go test ./...

# Run migrations
go run ./cmd/migrate up
go run ./cmd/migrate down
```

## Deployment

### Build
```bash
# Build binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o campus-api ./cmd/server

# Build Tailwind for production
cd static && npm run build:prod
```

### Deploy to VPS
```bash
# Copy binary and static files
scp campus-api user@vps:~/
scp -r static user@vps:~/
scp -r templates user@vps:~/  # If using file-based templates

# SSH and restart service
ssh user@vps
sudo systemctl restart campus-api
```

See `docs/DEPLOYMENT.md` for full deployment instructions.

## Dependencies

```
# Core
github.com/jackc/pgx/v5           # PostgreSQL driver
github.com/golang-jwt/jwt/v5       # JWT handling
golang.org/x/crypto                # bcrypt

# Templates
github.com/a-h/templ               # Type-safe templates

# Migrations
github.com/golang-migrate/migrate  # Database migrations

# Development
github.com/air-verse/air           # Hot reload (dev only)
```

## Performance

Expected resource usage on 1GB VPS:
- Go binary: ~50MB RAM
- PostgreSQL: ~200MB RAM
- Nginx: ~10MB RAM
- Total: ~360MB (36% of 1GB)

Capacity: 100,000+ leads without VPS upgrade
