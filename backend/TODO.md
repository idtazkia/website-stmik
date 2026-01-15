# Backend TODO - STMIK Tazkia Admission System

## Status: NOT STARTED

Go-based sales funnel management system. See `README.md` for architecture details.

---

## Phase 1: Project Setup

### 1.1 Initialize Project
- [ ] Create `go.mod`
- [ ] Create directory structure (`cmd/`, `internal/`, `templates/`, `static/`, `migrations/`)
- [ ] Create `.env.example`
- [ ] Create `cmd/server/main.go` entry point

### 1.2 Install Dependencies
```bash
go get github.com/jackc/pgx/v5
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/a-h/templ
go get github.com/golang-migrate/migrate/v4
```

### 1.3 Configuration
- [ ] Create `internal/config/config.go`
- [ ] Load environment variables
- [ ] Validate required config on startup

### 1.4 Static Assets Setup
- [ ] Initialize Tailwind CSS in `static/`
- [ ] Download HTMX (`static/js/htmx.min.js`)
- [ ] Download Alpine.js (`static/js/alpine.min.js`)
- [ ] Configure Tailwind with brand colors

---

## Phase 2: Database

### 2.1 Connection
- [ ] Create `internal/database/database.go`
- [ ] Configure pgx connection pool
- [ ] Add health check query

### 2.2 Migrations
- [ ] Create `migrations/001_create_users.up.sql`
- [ ] Create `migrations/001_create_users.down.sql`
- [ ] Create `migrations/002_create_lead_profiles.up.sql`
- [ ] Create `migrations/003_create_applications.up.sql`
- [ ] Create `migrations/004_create_documents.up.sql`
- [ ] Create `migrations/005_create_activity_log.up.sql`
- [ ] Create `cmd/migrate/main.go` CLI tool

### 2.3 Models
- [ ] Create `internal/models/user.go`
- [ ] Create `internal/models/lead_profile.go`
- [ ] Create `internal/models/application.go`
- [ ] Create `internal/models/document.go`

### 2.4 Repository Layer
- [ ] Create `internal/repository/user.go`
  - [ ] `Create(user)`
  - [ ] `FindByEmail(email)`
  - [ ] `FindByID(id)`
  - [ ] `UpdateProfile(id, data)`
- [ ] Create `internal/repository/lead.go`
  - [ ] `CreateLead(email, name, source)`
  - [ ] `GetProfile(userID)`
  - [ ] `UpdateProfile(userID, data)`
  - [ ] `ListLeads(filters, pagination)`
- [ ] Create `internal/repository/application.go`
  - [ ] `Create(userID, program)`
  - [ ] `GetByID(id)`
  - [ ] `GetByUserID(userID)`
  - [ ] `UpdateStatus(id, status, reviewerID)`
  - [ ] `ListApplications(filters, pagination)`

---

## Phase 3: Authentication

### 3.1 JWT Utilities
- [ ] Create `internal/services/auth.go`
  - [ ] `GenerateToken(userID, role)`
  - [ ] `ValidateToken(token)`
  - [ ] `GetUserFromToken(token)`

### 3.2 Password Utilities
- [ ] `HashPassword(password)`
- [ ] `VerifyPassword(password, hash)`

### 3.3 Middleware
- [ ] Create `internal/middleware/auth.go`
  - [ ] `RequireAuth` - verify JWT, attach user to context
  - [ ] `RequireRole(role)` - check user role
- [ ] Create `internal/middleware/logging.go`
- [ ] Create `internal/middleware/recovery.go`

### 3.4 Google OAuth
- [ ] Create `internal/services/oauth.go`
  - [ ] `GetGoogleAuthURL(redirectURL)`
  - [ ] `ExchangeGoogleCode(code)`
  - [ ] `GetGoogleUserInfo(accessToken)`
- [ ] Handle staff domain restriction

---

## Phase 4: Templates

### 4.1 Layouts
- [ ] Create `templates/layouts/base.templ`
  - [ ] HTML5 structure
  - [ ] Tailwind CSS include
  - [ ] HTMX include
  - [ ] Alpine.js include
- [ ] Create `templates/layouts/portal.templ`
- [ ] Create `templates/layouts/admin.templ`

### 4.2 Components
- [ ] Create `templates/components/button.templ`
- [ ] Create `templates/components/input.templ`
- [ ] Create `templates/components/card.templ`
- [ ] Create `templates/components/modal.templ`
- [ ] Create `templates/components/table.templ`
- [ ] Create `templates/components/pagination.templ`
- [ ] Create `templates/components/alert.templ`
- [ ] Create `templates/components/dropdown.templ` (Alpine.js)

### 4.3 Portal Pages
- [ ] Create `templates/pages/portal/login.templ`
- [ ] Create `templates/pages/portal/register.templ`
- [ ] Create `templates/pages/portal/profile.templ`
- [ ] Create `templates/pages/portal/application.templ`
- [ ] Create `templates/pages/portal/documents.templ`
- [ ] Create `templates/pages/portal/status.templ`

### 4.4 Admin Pages
- [ ] Create `templates/pages/admin/login.templ`
- [ ] Create `templates/pages/admin/dashboard.templ`
- [ ] Create `templates/pages/admin/leads.templ`
- [ ] Create `templates/pages/admin/lead_detail.templ`
- [ ] Create `templates/pages/admin/applications.templ`
- [ ] Create `templates/pages/admin/application_detail.templ`

---

## Phase 5: Handlers

### 5.1 Router Setup
- [ ] Create `internal/handlers/router.go`
  - [ ] Configure stdlib mux
  - [ ] Apply middleware chain
  - [ ] Register all routes
  - [ ] Serve static files

### 5.2 API Handlers
- [ ] Create `internal/handlers/api.go`
  - [ ] `POST /api/leads` - create lead from landing page
  - [ ] `GET /api/health` - health check

### 5.3 Auth Handlers
- [ ] Create `internal/handlers/auth.go`
  - [ ] `GET /portal/login` - render login page
  - [ ] `POST /portal/login` - process login
  - [ ] `GET /portal/register` - render register page
  - [ ] `POST /portal/register` - process registration
  - [ ] `GET /portal/auth/google` - initiate OAuth
  - [ ] `GET /portal/auth/google/callback` - OAuth callback
  - [ ] `POST /portal/logout` - clear session
  - [ ] Same for `/admin/` routes (Google only)

### 5.4 Portal Handlers
- [ ] Create `internal/handlers/portal.go`
  - [ ] `GET /portal/profile` - profile form
  - [ ] `POST /portal/profile` - update profile (HTMX)
  - [ ] `GET /portal/application` - application form
  - [ ] `POST /portal/application` - submit (HTMX)
  - [ ] `GET /portal/documents` - upload page
  - [ ] `POST /portal/documents` - file upload (HTMX)
  - [ ] `GET /portal/status` - status page

### 5.5 Admin Handlers
- [ ] Create `internal/handlers/admin.go`
  - [ ] `GET /admin/dashboard` - stats dashboard
  - [ ] `GET /admin/leads` - lead list with HTMX filters
  - [ ] `GET /admin/leads/:id` - lead detail
  - [ ] `POST /admin/leads/:id/status` - update status (HTMX)
  - [ ] `GET /admin/applications` - application list
  - [ ] `GET /admin/applications/:id` - application detail
  - [ ] `POST /admin/applications/:id/approve` (HTMX)
  - [ ] `POST /admin/applications/:id/reject` (HTMX)

---

## Phase 6: File Upload

### 6.1 Storage
- [ ] Create `internal/services/storage.go`
  - [ ] `SaveFile(file, docType)` - save to disk
  - [ ] `GetFilePath(filename)`
  - [ ] `DeleteFile(filename)`
  - [ ] Validate file type (PDF, JPG, PNG)
  - [ ] Validate file size (max 5MB)

### 6.2 Document Handlers
- [ ] Handle multipart form uploads
- [ ] Generate unique filenames
- [ ] Store metadata in database
- [ ] Serve files with auth check

---

## Phase 7: Services

### 7.1 Lead Service
- [ ] Create `internal/services/lead.go`
  - [ ] `CreateLead(email, name, source)`
  - [ ] `GetLeadPipeline(filters)`
  - [ ] `UpdateLeadStatus(id, status)`

### 7.2 Application Service
- [ ] Create `internal/services/application.go`
  - [ ] `CreateDraft(userID, program)`
  - [ ] `SubmitApplication(id)`
  - [ ] `ReviewApplication(id, status, reviewerID)`
  - [ ] `GetApplicationWithDocuments(id)`

### 7.3 Notification Service (optional)
- [ ] Create `internal/services/notification.go`
  - [ ] `SendWelcomeEmail(user)`
  - [ ] `SendStatusUpdate(user, status)`

---

## Phase 8: Testing

### 8.1 Unit Tests
- [ ] Test JWT utilities
- [ ] Test password hashing
- [ ] Test repository queries
- [ ] Test services

### 8.2 Integration Tests
- [ ] Test API endpoints
- [ ] Test auth flow
- [ ] Test file uploads

### 8.3 E2E Tests
- [ ] Test lead journey (landing → application)
- [ ] Test admin workflow

---

## Phase 9: Deployment

### 9.1 Build
- [ ] Create build script
- [ ] Build production binary
- [ ] Build production CSS

### 9.2 Systemd Service
- [ ] Create service file
- [ ] Configure auto-restart
- [ ] Set up log rotation

### 9.3 Nginx Config
- [ ] Reverse proxy to Go app
- [ ] Static file serving
- [ ] SSL with Certbot
- [ ] Rate limiting

---

## Dependencies

```go
// go.mod
require (
    github.com/jackc/pgx/v5 v5.x
    github.com/golang-jwt/jwt/v5 v5.x
    golang.org/x/crypto v0.x
    github.com/a-h/templ v0.x
    github.com/golang-migrate/migrate/v4 v4.x
)
```

---

## Success Criteria

- [ ] Lead can register, complete profile, submit application
- [ ] Staff can login via Google, view/filter leads, approve/reject
- [ ] File uploads working (documents)
- [ ] HTMX interactions smooth (no full page reloads)
- [ ] Mobile responsive
- [ ] All forms have validation
- [ ] Deployed and running on VPS
