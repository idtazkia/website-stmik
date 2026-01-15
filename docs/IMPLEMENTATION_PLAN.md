# Implementation Plan - STMIK Tazkia Campus Website

**Live Site:** https://stmik.tazkia.ac.id/

---

## Progress Overview

| Phase | Description | Status | Progress |
|-------|-------------|--------|----------|
| 1 | Infrastructure Setup | Partial | 60% |
| 2 | Frontend Marketing Site | In Progress | 30% |
| 3 | Backend API (Go) | Not Started | 0% |
| 4 | Deployment & DevOps | Partial | 20% |
| 5 | Testing & Polish | Not Started | 0% |
| 6 | Launch Preparation | Not Started | 0% |

**Current Focus:** Phase 2 - Frontend Marketing Site

---

## Phase 1: Infrastructure Setup (60%)

### 1.1 Repository & Documentation
- [x] Repository structure created
- [x] README.md
- [x] CLAUDE.md
- [x] docs/ARCHITECTURE.md
- [x] docs/DEPLOYMENT.md
- [x] docs/IMPLEMENTATION_PLAN.md

### 1.2 Frontend Infrastructure
- [x] Astro 5.x project setup
- [x] Tailwind CSS 4.x configuration
- [x] Custom i18n implementation
- [x] ESLint + Prettier configuration
- [x] TypeScript configuration

### 1.3 Deployment Infrastructure
- [x] Cloudflare Pages configured
- [x] Auto-deploy on git push
- [x] Custom domain (stmik.tazkia.ac.id)
- [x] SSL/TLS enabled

### 1.4 Backend Infrastructure (Deferred)
- [ ] VPS provisioning
- [ ] PostgreSQL 18 installation
- [ ] Nginx reverse proxy
- [ ] systemd service configuration
- [ ] SSL with Certbot

### 1.5 External Services (Deferred)
- [ ] Google OAuth credentials
- [ ] WhatsApp Business API access
- [ ] Kafka cluster access

---

## Phase 2: Frontend Marketing Site (30%)

### 2.1 Layouts & Components (Done)
- [x] BaseLayout.astro
- [x] MarketingLayout.astro
- [x] Header component
- [x] Footer component
- [x] Navigation component
- [x] Card component
- [x] Button component
- [x] Container component
- [x] Section component

### 2.2 Core Pages (Done)
- [x] Homepage (ID/EN)
- [x] About page (ID/EN)

### 2.3 Lecturer Profiles (Done)
- [x] Content collection setup
- [x] Lecturer list page (ID/EN)
- [x] Lecturer detail pages (ID/EN)

### 2.4 Programs Pages (Not Started)
- [ ] Content collection for programs
- [ ] `src/pages/programs/index.astro` - Program listing
- [ ] `src/pages/programs/[slug].astro` - Program detail
- [ ] `src/pages/en/programs/index.astro`
- [ ] `src/pages/en/programs/[slug].astro`
- [ ] Program cards with images
- [ ] Curriculum overview
- [ ] Career prospects section
- [ ] Admission requirements
- [ ] "Apply Now" CTA

### 2.5 Contact Page (Not Started)
- [ ] `src/pages/contact.astro`
- [ ] `src/pages/en/contact.astro`
- [ ] Contact information (phone, email, address)
- [ ] Embedded Google Maps
- [ ] Social media links
- [ ] Office hours
- [ ] Contact form (static, submission deferred)

### 2.6 Admissions Page (Not Started)
- [ ] `src/pages/admissions.astro`
- [ ] `src/pages/en/admissions.astro`
- [ ] Admission requirements
- [ ] Application process steps
- [ ] Important dates/calendar
- [ ] Required documents list
- [ ] FAQ section
- [ ] "Apply Now" CTA

### 2.7 News/Blog System (Not Started)
- [ ] Content collection for news
- [ ] `src/pages/news/index.astro` - News listing
- [ ] `src/pages/news/[slug].astro` - News detail
- [ ] `src/pages/en/news/index.astro`
- [ ] `src/pages/en/news/[slug].astro`
- [ ] Pagination
- [ ] Categories/tags
- [ ] Related posts

### 2.8 SEO Enhancements (Not Started)
- [ ] JSON-LD structured data (Organization, EducationalOrganization, Course)
- [ ] Breadcrumb navigation
- [ ] Internal linking optimization
- [ ] Meta descriptions per page

### 2.9 Performance & Accessibility (Not Started)
- [ ] Astro Image component optimization
- [ ] Lazy loading for images
- [ ] ARIA labels
- [ ] Keyboard navigation
- [ ] Color contrast (WCAG AA)

---

## Phase 3: Backend API - Go (0%)

### 3.1 Project Setup
- [ ] Initialize `go.mod`
- [ ] Create directory structure
  - [ ] `cmd/server/main.go`
  - [ ] `cmd/migrate/main.go`
  - [ ] `handler/`
  - [ ] `model/`
  - [ ] `migrations/`
  - [ ] `templates/`
  - [ ] `static/`
- [ ] Create `.env.example`
- [ ] `config.go` - Environment configuration

### 3.2 Dependencies
- [ ] pgx/v5 - PostgreSQL driver
- [ ] golang-jwt/jwt/v5 - JWT handling
- [ ] golang.org/x/crypto/bcrypt - Password hashing
- [ ] a-h/templ - HTML templates
- [ ] golang-migrate/migrate/v4 - Database migrations
- [ ] segmentio/kafka-go - Kafka consumer

### 3.3 Static Assets
- [ ] Tailwind CSS setup
- [ ] HTMX download
- [ ] Alpine.js + alpine-csp download
- [ ] Brand colors configuration

### 3.4 Database Migrations
- [ ] `001_create_users`
- [ ] `002_create_intakes`
- [ ] `003_create_programs`
- [ ] `004_create_tracks`
- [ ] `005_create_cancel_reasons`
- [ ] `006_create_referrers`
- [ ] `007_create_campaigns`
- [ ] `008_create_prospects`
- [ ] `009_create_applications`
- [ ] `010_create_documents`
- [ ] `011_create_document_checklists`
- [ ] `012_create_document_reviews`
- [ ] `013_create_activity_log`
- [ ] `014_create_communication_log`
- [ ] `015_seed_data`

### 3.5 Models
- [ ] `model/db.go` - Connection pool
- [ ] `model/user.go`
- [ ] `model/prospect.go`
- [ ] `model/application.go`
- [ ] `model/document.go`
- [ ] `model/lookup.go` (Intake, Program, Track, CancelReason, Referrer, Campaign)

### 3.6 Authentication
- [ ] `auth.go`
  - [ ] JWT generation/validation
  - [ ] Password hashing
  - [ ] Google OAuth helpers
  - [ ] Staff domain check
- [ ] `handler/router.go`
  - [ ] RequireAuth middleware
  - [ ] RequireRole middleware
  - [ ] Logging middleware
  - [ ] Recovery middleware

### 3.7 Templ Templates
- [ ] `templates/layout.templ` - Base, portal, admin layouts
- [ ] `templates/components.templ` - UI components
- [ ] `templates/portal.templ` - Registrant pages
- [ ] `templates/admin.templ` - Staff pages
- [ ] `templates/settings.templ` - Admin settings

### 3.8 Handlers
- [ ] `handler/router.go` - Route registration
- [ ] `handler/api.go` - Public API endpoints
- [ ] `handler/auth.go` - Authentication
- [ ] `handler/portal.go` - Registrant portal
- [ ] `handler/admin.go` - Admin dashboard

### 3.9 Integrations
- [ ] `whatsapp.go` - WhatsApp Business API
- [ ] `kafka.go` - Payment event consumer

### 3.10 File Upload
- [ ] Type validation (PDF, JPG, PNG)
- [ ] Size validation (max 5MB)
- [ ] Unique filename generation
- [ ] Auth-protected file serving

---

## Phase 4: Deployment & DevOps (20%)

### 4.1 Frontend Deployment (Done)
- [x] Cloudflare Pages configuration
- [x] Auto-deploy on git push
- [x] Custom domain with SSL

### 4.2 Backend Deployment (Not Started)
- [ ] VPS provisioning (Ubuntu 24.04 LTS)
- [ ] Go installation
- [ ] PostgreSQL 18 installation
- [ ] Nginx reverse proxy configuration
- [ ] SSL with Certbot
- [ ] systemd service file
- [ ] UFW firewall configuration
- [ ] Fail2ban setup

### 4.3 CI/CD (Not Started)
- [ ] GitHub Actions: Frontend deploy (on frontend/** changes)
- [ ] GitHub Actions: Backend deploy (on backend/** changes)
- [ ] Database migration automation
- [ ] Health check verification

### 4.4 Backup & Maintenance (Not Started)
- [ ] PostgreSQL backup script (pg_dump)
- [ ] Backup rotation (keep 7 days)
- [ ] Log rotation configuration

---

## Phase 5: Testing & Polish (0%)

### 5.1 Frontend Testing
- [ ] Manual testing on mobile devices
- [ ] Cross-browser testing (Chrome, Firefox, Safari, Edge)
- [ ] Navigation link verification
- [ ] Language switcher testing
- [ ] Responsive breakpoint testing

### 5.2 Performance Testing
- [ ] Lighthouse audit (target 90+)
- [ ] Core Web Vitals optimization
- [ ] Page load time <2s verification

### 5.3 Backend Testing
- [ ] Unit tests (auth, assignment logic, validation)
- [ ] Integration tests (API endpoints, auth flow)
- [ ] E2E tests (prospect → enrolled journey)

### 5.4 Security Audit
- [ ] Input validation review
- [ ] SQL injection prevention verification
- [ ] XSS protection verification
- [ ] CSRF protection verification
- [ ] Authentication flow review

---

## Phase 6: Launch Preparation (0%)

### 6.1 Environment Verification
- [ ] Production environment variables set
- [ ] Database connection verified
- [ ] External service connections verified (Google OAuth, WhatsApp, Kafka)

### 6.2 Monitoring Setup
- [ ] systemd service monitoring
- [ ] PostgreSQL status monitoring
- [ ] Nginx status monitoring
- [ ] Disk space alerts

### 6.3 Security Review
- [ ] SSH key-only authentication
- [ ] Firewall rules verified
- [ ] SSL certificate auto-renewal verified
- [ ] Fail2ban active

### 6.4 Soft Launch
- [ ] Internal testing with staff
- [ ] Test registration flow
- [ ] Test admin workflow
- [ ] Bug fixes

### 6.5 Production Launch
- [ ] DNS switch to production domain
- [ ] Monitor for issues
- [ ] Documentation finalization

---

## Milestones

| Milestone | Description | Target | Status |
|-----------|-------------|--------|--------|
| M1 | Marketing site live | - | Done |
| M2 | All marketing pages complete | - | In Progress |
| M3 | Backend API functional | - | Not Started |
| M4 | Registration portal live | - | Not Started |
| M5 | Admin dashboard live | - | Not Started |
| M6 | Full system launch | - | Not Started |

---

## Dependencies

```
Phase 1 (Infrastructure) ─┬─► Phase 2 (Frontend) ─► Phase 5 (Testing)
                          │
                          └─► Phase 3 (Backend) ─┬─► Phase 4 (Deploy Backend)
                                                 │
                                                 └─► Phase 5 (Testing) ─► Phase 6 (Launch)
```

**Critical Path:** Phase 1 → Phase 3 → Phase 4 → Phase 5 → Phase 6

**Parallel Work:** Phase 2 (Frontend) can proceed independently of Phase 3 (Backend)

---

## Success Criteria

### Marketing Site (Phase 2)
- [ ] All static pages live (Home, About, Programs, Contact, Admissions, News)
- [ ] Bilingual content complete (ID + EN)
- [ ] Mobile responsive
- [ ] Page load <2s
- [ ] Lighthouse score 90+

### Full System (Phase 3-6)
- [ ] 3,000 leads, 300 registrations per cycle (10% conversion)
- [ ] Staff can review and approve applications
- [ ] Document upload with checklist review
- [ ] Referral and campaign tracking functional
- [ ] WhatsApp notifications working
- [ ] Kafka payment integration working
- [ ] 99% uptime
- [ ] $5/month VPS cost maintained
- [ ] Zero security incidents

---

## Architecture Summary

```
┌─────────────────────────────────────────────────────────────┐
│                      Cloudflare (Free)                       │
│  ┌─────────────────┐  ┌─────────────────────────────────┐   │
│  │ Cloudflare CDN  │  │ Cloudflare Pages                │   │
│  │ DDoS Protection │  │ Astro Static Site (Marketing)   │   │
│  └────────┬────────┘  └─────────────────────────────────┘   │
└───────────┼─────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────────────────────────────────────────────┐
│                      VPS ($5/month)                          │
│  ┌─────────────────┐  ┌─────────────────────────────────┐   │
│  │ Nginx           │  │ Go Backend                      │   │
│  │ Reverse Proxy   │─►│ - Templ + HTMX + Alpine.js      │   │
│  │ Rate Limiting   │  │ - JWT + Google OAuth            │   │
│  │ SSL/TLS         │  │ - File uploads                  │   │
│  └─────────────────┘  └────────────────┬────────────────┘   │
│                                        │                     │
│                       ┌────────────────▼────────────────┐   │
│                       │ PostgreSQL 18                   │   │
│                       │ UUID primary keys               │   │
│                       └─────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Quick Reference

### Frontend Development
```bash
cd frontend
npm run dev          # http://localhost:4321
npm run build        # Production build
npm run preview      # Preview build
```

### Backend Development
```bash
cd backend
go run ./cmd/server              # Run server
go run ./cmd/migrate up          # Run migrations
templ generate                   # Generate templates
CGO_ENABLED=0 go build -o campus-api ./cmd/server
```

### Deployment
```bash
# Frontend: Auto-deploys on git push to main
git push

# Backend: Manual deploy via Ansible
cd infrastructure/ansible
ansible-playbook -i inventory/production.ini playbooks/deploy-backend.yml
```

---

## Related Documentation

| Document | Description |
|----------|-------------|
| `README.md` | Project overview |
| `CLAUDE.md` | Development guidance |
| `docs/ARCHITECTURE.md` | Technical design |
| `docs/DEPLOYMENT.md` | Deployment guide |
| `frontend/TODO.md` | Frontend task details |
| `backend/TODO.md` | Backend task details |
| `backend/README.md` | Backend architecture |
| `infrastructure/TODO.md` | Infrastructure tasks |
