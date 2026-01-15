# Backend TODO - STMIK Tazkia Admission System

## Status: NOT STARTED

Go-based sales funnel management system. See `README.md` for architecture and feature details.

---

## Phase 1: Project Setup

### 1.1 Initialize Project
- [ ] Create `go.mod`
- [ ] Create directory structure (`handler/`, `model/`, `migrations/`, `templates/`, `static/`)
- [ ] Create `.env.example`
- [ ] Create `cmd/server/main.go`
- [ ] Create `cmd/migrate/main.go`

### 1.2 Dependencies
```bash
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/a-h/templ
go get github.com/golang-migrate/migrate/v4
go get github.com/segmentio/kafka-go
```

### 1.3 Configuration
- [ ] `config.go` - Load and validate environment variables

### 1.4 Static Assets
- [ ] Initialize Tailwind CSS in `static/`
- [ ] Download HTMX
- [ ] Download Alpine.js
- [ ] Configure Tailwind with brand colors

---

## Phase 2: Database

### 2.1 Connection
- [ ] `model/db.go` - pgx connection pool

### 2.2 Migrations
- [ ] `001_create_users` - Users table
- [ ] `002_create_intakes` - Intake periods
- [ ] `003_create_programs` - Programs (SI, TI)
- [ ] `004_create_tracks` - Funding tracks
- [ ] `005_create_cancel_reasons` - Cancel reason categories
- [ ] `006_create_referrers` - Referral tracking
- [ ] `007_create_campaigns` - Ad campaign tracking
- [ ] `008_create_prospects` - Prospects with UTM fields
- [ ] `009_create_applications` - Applications
- [ ] `010_create_documents` - Uploaded documents
- [ ] `011_create_document_checklists` - Checklist templates
- [ ] `012_create_document_reviews` - Review results
- [ ] `013_create_activity_log` - Audit trail
- [ ] `014_create_communication_log` - WhatsApp/email log
- [ ] `015_seed_data` - Programs, tracks, reasons, checklists

### 2.3 Models (structs + queries in same file)
- [ ] `model/user.go` - User struct, Create, FindByEmail, FindByID
- [ ] `model/prospect.go` - Prospect struct + CRUD + list queries
- [ ] `model/application.go` - Application struct + CRUD
- [ ] `model/document.go` - Document struct + CRUD + review
- [ ] `model/lookup.go` - Intake, Program, Track, CancelReason, Referrer, Campaign

---

## Phase 3: Authentication

### 3.1 Auth Utils
- [ ] `auth.go`
  - [ ] GenerateToken(userID, role)
  - [ ] ValidateToken(token)
  - [ ] HashPassword / VerifyPassword
  - [ ] GetAuthURL(redirectURL)
  - [ ] ExchangeCode(code)
  - [ ] GetUserInfo(accessToken)
  - [ ] Staff domain check

### 3.2 Middleware
- [ ] `handler/router.go`
  - [ ] RequireAuth middleware
  - [ ] RequireRole(role) middleware
  - [ ] Logging middleware
  - [ ] Recovery middleware

---

## Phase 4: Templates

### 4.1 Layouts
- [ ] `templates/layout.templ` - Base HTML, portal layout, admin layout

### 4.2 Components
- [ ] `templates/components.templ` - button, input, select, card, modal, table, pagination, alert, dropdown, file_upload, checklist, timeline, stats_card

### 4.3 Portal Pages
- [ ] `templates/portal.templ`
  - [ ] Login page
  - [ ] Register page
  - [ ] Dashboard (status overview)
  - [ ] Application form (program/track selection)
  - [ ] Documents (upload KTP/Ijazah)
  - [ ] Cancel confirmation

### 4.4 Admin Pages
- [ ] `templates/admin.templ`
  - [ ] Login page
  - [ ] Dashboard (funnel stats)
  - [ ] Prospects list (with filters)
  - [ ] Prospect detail (timeline, actions)
  - [ ] Applications list
  - [ ] Application detail (docs, review)
  - [ ] Document review modal (checklist)
  - [ ] Reports (funnel, conversion, source/campaign/referrer)

### 4.5 Admin Settings Pages
- [ ] `templates/settings.templ`
  - [ ] Intakes management
  - [ ] Tracks management
  - [ ] Cancel reasons management
  - [ ] Document checklists management
  - [ ] Staff management
  - [ ] Referrers management
  - [ ] Campaigns management

---

## Phase 5: Handlers

### 5.1 Router Setup
- [ ] `handler/router.go`
  - [ ] Configure stdlib mux
  - [ ] Apply middleware
  - [ ] Register routes
  - [ ] Static file serving

### 5.2 API Handlers
- [ ] `handler/api.go`
  - [ ] `POST /api/prospects` - Create from landing page (with UTM + referral)
  - [ ] `GET /api/health`
  - [ ] `GET /api/referrers/{code}` - Validate referral code

### 5.3 Auth Handlers
- [ ] `handler/auth.go`
  - [ ] Portal: login, register, Google OAuth, logout
  - [ ] Admin: Google OAuth only, logout

### 5.4 Portal Handlers
- [ ] `handler/portal.go`
  - [ ] Dashboard
  - [ ] Application form (create/update)
  - [ ] Document upload/delete
  - [ ] Cancel application

### 5.5 Admin Handlers
- [ ] `handler/admin.go`
  - [ ] Dashboard (stats)
  - [ ] Prospects list/detail
  - [ ] Assign prospect (round-robin / manual)
  - [ ] Update prospect status
  - [ ] Cancel prospect
  - [ ] Applications list/detail
  - [ ] Document review (checklist)
  - [ ] Approve application
  - [ ] Cancel application
  - [ ] Send WhatsApp
  - [ ] Settings: Intakes, Tracks, Cancel reasons, Checklists, Staff, Referrers, Campaigns CRUD

### 5.6 Reports Handlers
- [ ] `handler/admin.go` (reports section)
  - [ ] Funnel data (by intake, date range)
  - [ ] Source breakdown (utm_source stats)
  - [ ] Campaign performance (leads, conversions, cost per lead)
  - [ ] Referrer performance (leads, conversions by referrer)
  - [ ] CSV export (prospects, applications, reports)

---

## Phase 6: Integrations

### 6.1 WhatsApp
- [ ] `whatsapp.go`
  - [ ] SendTemplate(phone, template, variables)
  - [ ] SendWelcome, SendFollowUp, SendDocumentReminder
  - [ ] SendRevisionRequest, SendApproved, SendEnrolled

### 6.2 Kafka Consumer
- [ ] `kafka.go`
  - [ ] Consumer for `payment.completed`
  - [ ] Match VA number to application
  - [ ] Update status: approved → enrolled
  - [ ] Log payment details

---

## Phase 7: File Upload

- [ ] `handler/portal.go` (document upload)
  - [ ] Validate type (PDF, JPG, PNG)
  - [ ] Validate size (max 5MB)
  - [ ] Generate unique filename
  - [ ] Save to upload directory
  - [ ] Serve files with auth check

---

## Phase 8: Testing

### 8.1 Unit Tests
- [ ] Auth (JWT, password hashing)
- [ ] Round-robin assignment
- [ ] Document validation

### 8.2 Integration Tests
- [ ] API endpoints
- [ ] Auth flow (local + Google)
- [ ] Document upload
- [ ] Kafka consumer

### 8.3 E2E Tests
- [ ] Prospect → Application → Enrolled journey
- [ ] Admin review workflow

---

## Phase 9: Deployment

### 9.1 Build
- [ ] Build script (binary + Tailwind)
- [ ] Dockerfile (optional)

### 9.2 VPS Setup
- [ ] Systemd service file
- [ ] Nginx reverse proxy config
- [ ] SSL with Certbot
- [ ] Kafka connection

### 9.3 CI/CD
- [ ] GitHub Actions workflow
- [ ] Auto-deploy on push

---

## Success Criteria

- [ ] Prospect can register from landing page
- [ ] Prospect can login, select program/track, upload docs
- [ ] Staff can login via Google
- [ ] Staff can view/filter prospects and applications
- [ ] Staff can review documents with checklist
- [ ] Staff can approve applications
- [ ] Round-robin assignment works
- [ ] WhatsApp notifications sent
- [ ] Kafka payment events update status to enrolled
- [ ] Dashboard shows funnel stats
- [ ] UTM parameters captured from landing page
- [ ] Referral codes tracked and linked to prospects
- [ ] Campaign performance reports available
- [ ] Referrer performance reports available
- [ ] All HTMX interactions work (no full page reloads)
- [ ] Mobile responsive
