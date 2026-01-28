# Implementation Plan - STMIK Tazkia Campus Website

**Live Site:** https://stmik.tazkia.ac.id/

**Last Updated:** 2026-01-28 (SEO launch complete)

---

## Progress Overview

| Phase | Description | Status | Progress |
|-------|-------------|--------|----------|
| 1 | Infrastructure Setup | Partial | 60% |
| 2 | Frontend Marketing Site | In Progress | 85% |
| 3 | Backend API (Go) | In Progress | 10% |
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

## Phase 2: Frontend Marketing Site (85%)

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

### 2.4 Programs Pages (Done)
- [x] `src/pages/programs/index.astro` - Program listing
- [x] `src/pages/en/programs/index.astro`
- [x] `src/pages/programs/information-systems.astro` - SI detail
- [x] `src/pages/programs/computer-engineering.astro` - TI detail
- [x] `src/pages/en/programs/information-systems.astro`
- [x] `src/pages/en/programs/computer-engineering.astro`
- [x] Program cards with specializations
- [x] Curriculum overview (8 semesters)
- [x] Career prospects section
- [x] "Apply Now" CTA

### 2.5 Contact Page (Done)
- [x] `src/pages/contact.astro`
- [x] `src/pages/en/contact.astro`
- [x] Contact information (phone, email, address)
- [x] Embedded Google Maps
- [x] Social media links (Instagram, YouTube, LinkedIn)
- [x] Office hours
- [x] Contact form (static, backend submission deferred to Phase 3)

### 2.6 Admissions Page (Done)
- [x] `src/pages/admissions.astro`
- [x] `src/pages/en/admissions.astro`
- [x] Intake periods (3 periods with discount benefits)
- [x] Registration process (4 steps)
- [x] Required documents list (6 items)
- [x] Scholarship programs (4 types)
- [x] FAQ section (5 questions with accordion)
- [x] CTA sections

### 2.7 News/Blog System (Deferred)
> Deferred until content migration from old website is complete

- [ ] Content collection for news
- [ ] `src/pages/news/index.astro` - News listing
- [ ] `src/pages/news/[slug].astro` - News detail
- [ ] `src/pages/en/news/index.astro`
- [ ] `src/pages/en/news/[slug].astro`
- [ ] Pagination
- [ ] Categories/tags
- [ ] Related posts

#### Content Strategy - Article Topics for SEO

**Target Keywords:** kuliah IT Bogor, kampus informatika Bogor, jurusan sistem informasi, teknik informatika terbaik

**Program & Career Articles:**
- [ ] Prospek Kerja Lulusan Sistem Informasi 2026
- [ ] Perbedaan Sistem Informasi vs Teknik Informatika: Mana yang Cocok?
- [ ] 10 Skill yang Harus Dikuasai Mahasiswa IT
- [ ] Gaji Fresh Graduate IT di Indonesia
- [ ] Peluang Karir di Bidang Data Science dan AI

**Industry Trends:**
- [ ] Tren Teknologi 2026 yang Wajib Dipelajari Mahasiswa
- [ ] Peluang Karir di Industri Halal Tech
- [ ] Cybersecurity: Profesi Masa Depan dengan Permintaan Tinggi
- [ ] Machine Learning untuk Pemula: Panduan Lengkap
- [ ] Cloud Computing dan Peluang Sertifikasi

**Student Life & Campus:**
- [ ] Kehidupan Mahasiswa di STMIK Tazkia
- [ ] Fasilitas Kampus STMIK Tazkia Bogor
- [ ] Project-Based Learning: Belajar Sambil Praktek Langsung
- [ ] Testimoni Alumni STMIK Tazkia

**Admissions & Guides:**
- [ ] Panduan Lengkap Pendaftaran STMIK Tazkia
- [ ] Biaya Kuliah IT di Bogor: Perbandingan Kampus
- [ ] Tips Memilih Jurusan IT yang Tepat
- [ ] Beasiswa Kuliah IT untuk Mahasiswa Berprestasi

**Local SEO (Bogor):**
- [ ] Kampus IT Terbaik di Bogor
- [ ] Kuliah Sambil Kerja di Bogor: Tips dan Rekomendasi
- [ ] Komunitas IT di Bogor yang Wajib Diikuti

#### News Topics (Regular Updates)

**Academic Events:**
- [ ] Pembukaan Pendaftaran Mahasiswa Baru Gelombang [X]
- [ ] Wisuda STMIK Tazkia Angkatan [X]
- [ ] Seminar Nasional Teknologi Informasi
- [ ] Workshop Industry Partnership dengan [Mitra]
- [ ] Kunjungan Industri ke [Perusahaan]

**Student Achievements:**
- [ ] Mahasiswa STMIK Tazkia Juara Kompetisi [X]
- [ ] Tim STMIK Tazkia Raih Penghargaan di Hackathon
- [ ] Karya Mahasiswa Digunakan oleh [Klien Industri]
- [ ] Alumni STMIK Tazkia Sukses di [Perusahaan]

**Partnerships & MoU:**
- [ ] STMIK Tazkia Jalin Kerjasama dengan [Perusahaan Tech]
- [ ] MoU dengan [Universitas/Institusi]
- [ ] Program Magang di [Perusahaan Mitra]
- [ ] Sertifikasi [Vendor] untuk Mahasiswa STMIK Tazkia

**Campus Updates:**
- [ ] Fasilitas Baru di Kampus STMIK Tazkia
- [ ] Kurikulum Baru Berbasis Industri 4.0
- [ ] Akreditasi Program Studi [X]
- [ ] Dosen STMIK Tazkia Raih [Penghargaan/Sertifikasi]

#### Partner Site Articles (Backlink Strategy)

> Articles to be published on partner websites with backlinks to stmik.tazkia.ac.id

**For Tech Media (Techinasia, Dailysocial, etc):**
- [ ] Bagaimana Kampus Menyiapkan Talenta untuk Industri Halal Tech
- [ ] Project-Based Learning: Model Pendidikan IT Masa Depan
- [ ] Kolaborasi Kampus-Industri dalam Mencetak Developer Siap Kerja

**For Education Portals (Rencanamu, Quipper, Youthmanual):**
- [ ] Mengenal Jurusan Sistem Informasi: Prospek dan Peluang
- [ ] 5 Alasan Memilih Kuliah IT di Bogor
- [ ] Kampus dengan Kurikulum Berbasis Proyek Nyata

**For Business/Halal Industry Media:**
- [ ] Peran Teknologi Informasi dalam Ekosistem Industri Halal
- [ ] Digital Transformation untuk UMKM Halal
- [ ] Kebutuhan Talenta IT di Industri Keuangan Syariah

**For Local Bogor Media:**
- [ ] STMIK Tazkia: Kampus IT Unggulan di Kota Bogor
- [ ] Kontribusi STMIK Tazkia untuk Pengembangan SDM Bogor
- [ ] Kerjasama STMIK Tazkia dengan Pelaku Usaha Lokal Bogor

**For Tazkia Group Sites (tazkia.ac.id, STEI Tazkia):**
- [ ] Sinergi Pendidikan Ekonomi Islam dan Teknologi Informasi
- [ ] Program Double Degree STEI Tazkia - STMIK Tazkia
- [ ] Kolaborasi Riset Fintech Syariah

### 2.8 SEO Enhancements (Done)
- [x] JSON-LD structured data (Organization, EducationalOrganization)
- [x] Course schema on program pages
- [x] Breadcrumb navigation component
- [x] Meta descriptions per page
- [x] Canonical URL (`<link rel="canonical">`)
- [x] Open Graph image (`og:image`, `og:image:width`, `og:image:height`)
- [x] Open Graph URL (`og:url`)
- [x] Twitter card image (`twitter:image`)
- [x] Per-page dynamic hreflang (id, en, x-default)
- [x] Default OG image created (1200x630px with branding)
- [x] BaseLayout `image` prop for custom OG images per page
- [x] Search engine indexing enabled (noIndex default changed to false)
- [x] robots.txt updated to allow crawling with sitemap reference
- [ ] Internal linking optimization
- [ ] Submit sitemap to Google Search Console

### 2.9 Performance & Accessibility (Partial)
- [ ] Astro Image component optimization
- [x] Lazy loading for images
- [x] Explicit image dimensions (width/height attributes for CLS)
- [x] ARIA labels on interactive elements
- [x] Keyboard navigation (dropdown menus, mobile menu)
- [x] Color contrast verified (WCAG AA)
- [x] Cloudflare WAF skip rule for stmik.tazkia.ac.id (performance optimization)

---

## Phase 3: Backend API - Go (10%)

### 3.1 Project Setup (Done)
- [x] Docker Compose for local development
  - [x] `docker-compose.yml` with PostgreSQL 18
  - [x] Volume for data persistence
  - [x] Health check configuration
- [x] Initialize `go.mod`
- [x] Create directory structure
  - [x] `cmd/server/main.go`
  - [x] `cmd/migrate/main.go`
  - [x] `handler/`
  - [x] `model/`
  - [x] `migrations/`
  - [x] `templates/`
  - [x] `static/`
- [x] Create `.env.example`
- [x] `config/config.go` - Environment configuration

### 3.2 Dependencies
- [ ] pgx/v5 - PostgreSQL driver
- [ ] golang-jwt/jwt/v5 - JWT handling
- [ ] golang.org/x/crypto/bcrypt - Password hashing
- [ ] a-h/templ - HTML templates
- [ ] golang-migrate/migrate/v4 - Database migrations
- [ ] segmentio/kafka-go - Kafka consumer

### 3.3 Static Assets
- [ ] Tailwind CSS v4 setup (standalone CLI)
- [ ] HTMX download
- [ ] Alpine.js + alpine-csp download
- [ ] Brand colors from frontend design system

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
- [x] Home page live (ID + EN)
- [x] About page live (ID + EN)
- [x] Lecturers pages live (ID + EN)
- [x] Programs pages live (ID + EN)
- [x] Contact page live (ID + EN)
- [x] Admissions page live (ID + EN)
- [ ] News/Blog system (ID + EN)
- [x] Mobile responsive
- [x] SEO score 100 (Lighthouse)
- [x] Accessibility score 96 (Lighthouse)
- [ ] Performance score 90+ (blocked by Cloudflare challenge script ~50)
- [x] Search engine indexing enabled

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
