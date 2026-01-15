# Backend TODO - STMIK Tazkia Admission System

## Status: NOT STARTED

Go-based sales funnel management system. See `README.md` for architecture and feature details.

---

## Phase 1: Project Setup

### 1.1 Initialize Project
- [ ] Create `go.mod`
- [ ] Create directory structure
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
- [ ] `internal/config/config.go` - Load and validate environment variables

### 1.4 Static Assets
- [ ] Initialize Tailwind CSS in `static/`
- [ ] Download HTMX
- [ ] Download Alpine.js
- [ ] Configure Tailwind with brand colors

---

## Phase 2: Database

### 2.1 Connection
- [ ] `internal/database/database.go` - pgx connection pool

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

### 2.3 Models
- [ ] `internal/models/user.go`
- [ ] `internal/models/intake.go`
- [ ] `internal/models/program.go`
- [ ] `internal/models/track.go`
- [ ] `internal/models/referrer.go`
- [ ] `internal/models/campaign.go`
- [ ] `internal/models/prospect.go`
- [ ] `internal/models/application.go`
- [ ] `internal/models/document.go`

### 2.4 Repository
- [ ] `internal/repository/user.go`
  - [ ] Create, FindByEmail, FindByID
- [ ] `internal/repository/intake.go`
  - [ ] Create, List, GetActive, Update
- [ ] `internal/repository/referrer.go`
  - [ ] Create, FindByCode, List, Update, GetStats
- [ ] `internal/repository/campaign.go`
  - [ ] Create, FindByUTM, List, Update, GetStats
- [ ] `internal/repository/prospect.go`
  - [ ] Create (with UTM + referrer), FindByID, FindByEmail
  - [ ] List (filter by source, campaign, referrer), UpdateStatus, Cancel
- [ ] `internal/repository/application.go`
  - [ ] Create, FindByID, FindByProspectID, List, UpdateStatus, Cancel
- [ ] `internal/repository/document.go`
  - [ ] Create, FindByApplicationID, UpdateStatus
- [ ] `internal/repository/activity.go`
  - [ ] Log, GetByEntity

---

## Phase 3: Authentication

### 3.1 Auth Service
- [ ] `internal/services/auth.go`
  - [ ] GenerateToken(userID, role)
  - [ ] ValidateToken(token)
  - [ ] HashPassword / VerifyPassword

### 3.2 Google OAuth
- [ ] `internal/services/oauth.go`
  - [ ] GetAuthURL(redirectURL)
  - [ ] ExchangeCode(code)
  - [ ] GetUserInfo(accessToken)
  - [ ] Staff domain check

### 3.3 Middleware
- [ ] `internal/middleware/auth.go`
  - [ ] RequireAuth
  - [ ] RequireRole(role)
- [ ] `internal/middleware/logging.go`
- [ ] `internal/middleware/recovery.go`

---

## Phase 4: Templates

### 4.1 Layouts
- [ ] `templates/layouts/base.templ` - HTML structure, CSS/JS includes
- [ ] `templates/layouts/portal.templ` - Registrant layout
- [ ] `templates/layouts/admin.templ` - Staff layout

### 4.2 Components
- [ ] `templates/components/button.templ`
- [ ] `templates/components/input.templ`
- [ ] `templates/components/select.templ`
- [ ] `templates/components/card.templ`
- [ ] `templates/components/modal.templ` (Alpine.js)
- [ ] `templates/components/table.templ`
- [ ] `templates/components/pagination.templ`
- [ ] `templates/components/alert.templ`
- [ ] `templates/components/dropdown.templ` (Alpine.js)
- [ ] `templates/components/file_upload.templ`
- [ ] `templates/components/checklist.templ`
- [ ] `templates/components/timeline.templ`
- [ ] `templates/components/stats_card.templ`

### 4.3 Portal Pages
- [ ] `templates/pages/portal/login.templ`
- [ ] `templates/pages/portal/register.templ`
- [ ] `templates/pages/portal/dashboard.templ` - Status overview
- [ ] `templates/pages/portal/application.templ` - Program/track selection
- [ ] `templates/pages/portal/documents.templ` - Upload KTP/Ijazah
- [ ] `templates/pages/portal/cancel.templ` - Cancel confirmation

### 4.4 Admin Pages
- [ ] `templates/pages/admin/login.templ`
- [ ] `templates/pages/admin/dashboard.templ` - Funnel stats
- [ ] `templates/pages/admin/prospects.templ` - List with filters
- [ ] `templates/pages/admin/prospect_detail.templ` - Timeline, actions
- [ ] `templates/pages/admin/applications.templ` - List with filters
- [ ] `templates/pages/admin/application_detail.templ` - Docs, review
- [ ] `templates/pages/admin/document_review.templ` - Checklist modal
- [ ] `templates/pages/admin/reports.templ` - Funnel, conversion, source/campaign/referrer

### 4.5 Admin Settings Pages
- [ ] `templates/pages/admin/settings/intakes.templ`
- [ ] `templates/pages/admin/settings/tracks.templ`
- [ ] `templates/pages/admin/settings/cancel_reasons.templ`
- [ ] `templates/pages/admin/settings/checklists.templ`
- [ ] `templates/pages/admin/settings/staff.templ`
- [ ] `templates/pages/admin/settings/referrers.templ` - Referral partner management
- [ ] `templates/pages/admin/settings/campaigns.templ` - Ad campaign management

---

## Phase 5: Handlers

### 5.1 Router Setup
- [ ] `internal/handlers/router.go`
  - [ ] Configure stdlib mux
  - [ ] Apply middleware
  - [ ] Register routes
  - [ ] Static file serving

### 5.2 API Handlers
- [ ] `internal/handlers/api.go`
  - [ ] `POST /api/prospects` - Create from landing page
  - [ ] `GET /api/health`

### 5.3 Auth Handlers
- [ ] `internal/handlers/auth.go`
  - [ ] Portal: login, register, Google OAuth, logout
  - [ ] Admin: Google OAuth only, logout

### 5.4 Portal Handlers
- [ ] `internal/handlers/portal.go`
  - [ ] Dashboard
  - [ ] Application form (create/update)
  - [ ] Document upload/delete
  - [ ] Cancel application

### 5.5 Admin Handlers
- [ ] `internal/handlers/admin.go`
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

### 5.6 Admin Settings Handlers
- [ ] `internal/handlers/settings.go`
  - [ ] Intakes CRUD
  - [ ] Tracks CRUD
  - [ ] Cancel reasons CRUD
  - [ ] Document checklists CRUD
  - [ ] Staff toggle active
  - [ ] Referrers CRUD (code, name, type, contact)
  - [ ] Campaigns CRUD (name, utm_campaign, dates, budget)

### 5.7 Reports Handlers
- [ ] `internal/handlers/reports.go`
  - [ ] Funnel data (by intake, date range)
  - [ ] Source breakdown (utm_source stats)
  - [ ] Campaign performance (leads, conversions, cost per lead)
  - [ ] Referrer performance (leads, conversions by referrer)
  - [ ] CSV export (prospects, applications, reports)

---

## Phase 6: Services

### 6.1 Prospect Service
- [ ] `internal/services/prospect.go`
  - [ ] CreateProspect (from landing page)
  - [ ] AssignToStaff (round-robin)
  - [ ] UpdateStatus
  - [ ] Cancel

### 6.2 Application Service
- [ ] `internal/services/application.go`
  - [ ] Create
  - [ ] UpdateProgramTrack
  - [ ] SubmitForReview
  - [ ] Approve (check all docs approved)
  - [ ] Cancel

### 6.3 Document Service
- [ ] `internal/services/document.go`
  - [ ] Upload (validate type/size)
  - [ ] Delete
  - [ ] ReviewDocument (checklist)
  - [ ] GetMissingDocuments

### 6.4 Assignment Service
- [ ] `internal/assignment/roundrobin.go`
  - [ ] GetNextStaff
  - [ ] AssignProspect
  - [ ] ReassignProspect

---

## Phase 7: Integrations

### 7.1 WhatsApp Service
- [ ] `internal/services/whatsapp.go`
  - [ ] SendTemplate(phone, template, variables)
  - [ ] SendWelcome
  - [ ] SendFollowUp
  - [ ] SendDocumentReminder
  - [ ] SendRevisionRequest
  - [ ] SendApproved
  - [ ] SendEnrolled

### 7.2 Kafka Consumer
- [ ] `internal/services/kafka.go`
  - [ ] Consumer for `payment.completed`
  - [ ] Match VA number to application
  - [ ] Update status: approved → enrolled
  - [ ] Log payment details

### 7.3 Email Service (Optional)
- [ ] `internal/services/email.go`
  - [ ] Same templates as WhatsApp
  - [ ] PDF acceptance letter generation

---

## Phase 8: File Upload

- [ ] `internal/services/storage.go`
  - [ ] SaveFile (validate type: PDF, JPG, PNG)
  - [ ] SaveFile (validate size: max 5MB)
  - [ ] GenerateUniqueName
  - [ ] DeleteFile
  - [ ] ServeFile (with auth check)

---

## Phase 9: Testing

### 9.1 Unit Tests
- [ ] Auth service (JWT, password)
- [ ] Assignment service (round-robin)
- [ ] Document validation

### 9.2 Integration Tests
- [ ] API endpoints
- [ ] Auth flow (local + Google)
- [ ] Document upload
- [ ] Kafka consumer

### 9.3 E2E Tests
- [ ] Prospect → Application → Enrolled journey
- [ ] Admin review workflow

---

## Phase 10: Deployment

### 10.1 Build
- [ ] Build script (binary + Tailwind)
- [ ] Dockerfile (optional)

### 10.2 VPS Setup
- [ ] Systemd service file
- [ ] Nginx reverse proxy config
- [ ] SSL with Certbot
- [ ] Kafka connection

### 10.3 CI/CD
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
