# Campus Website - Project TODO

**Live Site:** https://dev.stmik.tazkia.ac.id/

---

## Overall Progress

| Phase | Component | Status | Progress |
|-------|-----------|--------|----------|
| **Phase 1** | Infrastructure Setup | Partial | 60% |
| **Phase 2** | Backend API (Go) | Deferred | 0% |
| **Phase 3** | Frontend Marketing Site | In Progress | 30% |
| **Phase 4** | Deployment & DevOps | Partial | 20% |
| **Phase 5** | Testing & Polish | Not Started | 0% |
| **Phase 6** | Launch Preparation | Not Started | 0% |

---

## Current Focus: Phase 3 - Marketing Site

### Completed
- Astro 5.x + Tailwind CSS 4.x setup
- Bilingual system (Indonesian/English) with custom i18n
- Layouts: BaseLayout, MarketingLayout
- Components: Header, Footer, Navigation, Card, Button, Container, Section
- Pages: Homepage, About, Lecturer Profiles
- Content Collections: Lecturers
- SEO: Meta tags, sitemap, Open Graph tags
- Deployment: Cloudflare Pages with auto-deploy

### Next Tasks
- [ ] Programs listing and detail pages
- [ ] Contact page
- [ ] Admissions information page
- [ ] News/blog system

See `frontend/TODO.md` for detailed frontend tasks.

---

## Phase-by-Phase Overview

### Phase 1: Infrastructure Setup (60% Complete)

#### Completed
- Repository structure created
- Cloudflare Pages configured and deployed
- Custom domain configured (dev.stmik.tazkia.ac.id)
- Documentation created (README, CLAUDE, ARCHITECTURE, DEPLOYMENT)

#### Deferred
- VPS provisioning (when backend implementation starts)
- PostgreSQL 18 setup
- Google OAuth setup

---

### Phase 2: Backend Development (0% Complete)
**Status:** Deferred - backend not needed for marketing site

Will implement when authentication and application portal are needed.

**Stack:**
- Go (Golang) with stdlib router
- PostgreSQL 18 (UUID primary keys)
- Templ templates + HTMX + Alpine.js
- JWT + Google OAuth authentication

**Scope:**
- Sales funnel management (leads → applications → enrolled)
- Document upload and review with checklist
- Referral tracking and UTM campaign tracking
- WhatsApp notifications (REST API)
- Kafka integration for payment events
- Admin dashboard with reports

See `backend/TODO.md` for detailed backend tasks.

---

### Phase 3: Frontend Development (30% Complete)

**Completed:**
- Static site infrastructure
- Homepage, About page
- Lecturer profiles with content collections
- Bilingual routing (ID/EN)
- Responsive design

**Remaining:**
- Programs pages
- Contact page
- Admissions page
- News/blog system

See `frontend/TODO.md` for detailed frontend tasks.

---

### Phase 4: Deployment & DevOps (20% Complete)

#### Completed
- Cloudflare Pages deployment configured
- Auto-deploy on git push working
- Custom domain with SSL/TLS

#### Pending
- VPS setup and configuration
- Go binary deployment with systemd
- Nginx reverse proxy
- Database migrations
- Backup automation

---

### Phase 5: Testing & Polish (0% Complete)

**Scope:**
- End-to-end testing
- Performance testing (Lighthouse audits)
- Security audit
- Mobile responsiveness verification

---

### Phase 6: Launch Preparation (0% Complete)

**Scope:**
- Environment variables verification
- Backup configuration
- Monitoring setup
- Security review
- Soft launch with test users

---

## Success Metrics

### Marketing Site (Phase 3)
- [ ] All static pages live (Home, About, Programs, Contact, Admissions)
- [ ] Bilingual content complete (ID + EN)
- [ ] Mobile responsive
- [ ] Page load <2s
- [ ] Lighthouse score 90+

### Full System (Phases 2-6)
- [ ] 3,000 leads, 300 registrations per cycle
- [ ] Staff can review and approve applications
- [ ] Document upload with checklist review
- [ ] 99% uptime
- [ ] $5/month VPS cost maintained
- [ ] Zero security incidents

---

## Task Organization

- **Frontend tasks:** `frontend/TODO.md`
- **Backend tasks:** `backend/TODO.md`
- **Infrastructure tasks:** `infrastructure/TODO.md`

---

## Quick Start

### Frontend (Current Phase)
```bash
cd frontend
npm install
npm run dev          # http://localhost:4321
```

### Backend (Future Phase)
```bash
cd backend
go run ./cmd/server
```

---

## Architecture Summary

- **Frontend:** Astro static site on Cloudflare Pages (free)
- **Backend:** Go on VPS ($5/month)
- **Database:** PostgreSQL 18 on VPS (UUID primary keys)
- **Reverse Proxy:** Nginx with rate limiting
- **CDN/DDoS:** Cloudflare (free)

Target scale: 3,000 leads → 300 registrations (10% conversion)
