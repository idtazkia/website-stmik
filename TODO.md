# Campus Website - Implementation TODO

## Current Status (Updated: 2025-01-19)

**Live Site:** https://dev.stmik.tazkia.ac.id/

### ✅ Completed
- **Phase 1 (Partial):** Repository setup, Cloudflare Pages deployment
- **Phase 3 (Partial):** Marketing website foundation
  - Bilingual static site (Indonesian & English)
  - Homepage, About page, Lecturer profiles
  - Responsive design with Tailwind CSS
  - Custom i18n system
  - SEO optimization (meta tags, sitemap)

### 🚧 In Progress
- Phase 3: Additional marketing pages (Programs, Contact, Admissions, News)

### 📋 Pending
- Phase 1: VPS & backend infrastructure (on hold - static site first)
- Phase 2: Backend API development
- Phase 3: Authentication pages, Application portal, Admin interface
- Phase 4-8: BFF layer, deployment, testing, launch

### 📊 Progress
- Phase 1: 40% complete (infrastructure pending)
- Phase 2: 0% complete
- Phase 3: 30% complete (foundation done, dynamic features pending)
- Overall: ~20% complete

---

## Phase 1: Project Setup

### Infrastructure Setup
- [ ] Set up VPS with local cloud provider (DEFERRED - static site first)
  - [ ] Install Node.js 20+
  - [ ] Install PostgreSQL 14+
  - [ ] Configure firewall (allow ports 80, 443, 22)
  - [ ] Install pm2 for process management
  - [ ] Set up nginx as reverse proxy
- [x] Create Cloudflare account
  - [ ] Set up Cloudflare Workers (pending)
  - [x] Set up Cloudflare Pages
  - [x] Configure custom domain (dev.stmik.tazkia.ac.id)
- [ ] Set up Google Cloud project for OIDC (pending)
  - [ ] Create OAuth 2.0 credentials
  - [ ] Configure redirect URIs
  - [ ] Save client ID and client secret

### Repository Setup
- [x] Initialize monorepo structure
  - [x] Create `frontend/` directory
  - [ ] Create `backend/` directory (pending)
  - [ ] Create `shared/` directory (pending)
  - [x] Create `docs/` directory
  - [ ] Set up root `package.json` with workspaces (using frontend only currently)
  - [x] Create `.gitignore`
- [x] Set up GitHub repository (https://github.com/endymuhardin/website-stmik)
- [ ] Set up GitHub Actions (REMOVED - using Cloudflare Pages auto-deploy)
  - [x] ~~Create `.github/workflows/deploy-frontend.yml`~~ (removed in favor of Cloudflare)
  - [ ] Create `.github/workflows/deploy-backend.yml` (pending)

---

## Phase 2: Backend Development

### Database Setup
- [ ] Design and create database schema
  - [ ] Create `users` table
  - [ ] Create `applications` table
  - [ ] Create `sessions` table (if using Option B)
  - [ ] Create migration files
  - [ ] Write seed data for testing
- [ ] Set up database connection
  - [ ] Install `pg` (node-postgres)
  - [ ] Create `backend/src/config/database.js`
  - [ ] Implement connection pooling

### Authentication System
- [ ] Implement JWT utilities
  - [ ] Token generation
  - [ ] Token verification
  - [ ] Token refresh logic
- [ ] Create authentication middleware
  - [ ] `authenticateToken` middleware
  - [ ] Role-based access control
- [ ] Build authentication endpoints
  - [ ] `POST /auth/login` (email/password)
  - [ ] `POST /auth/register`
  - [ ] `POST /auth/google` (OIDC handler)
  - [ ] `POST /auth/logout`
  - [ ] `POST /auth/refresh-token` (optional)

### API Development
- [ ] User management endpoints
  - [ ] `GET /users/me` (get current user)
  - [ ] `PATCH /users/me` (update profile)
  - [ ] `GET /users` (admin only - list users)
- [ ] Application endpoints
  - [ ] `POST /applications` (submit application)
  - [ ] `GET /applications` (list user's applications)
  - [ ] `GET /applications/:id` (get single application)
  - [ ] `PATCH /applications/:id` (update application - draft mode)
  - [ ] `PATCH /applications/:id/status` (admin - approve/reject)
  - [ ] `DELETE /applications/:id` (delete draft)

### File Upload
- [ ] Set up file storage
  - [ ] Configure Cloudflare R2 or VPS storage
  - [ ] Implement file upload handler
  - [ ] File validation (size, type)
  - [ ] Generate secure file URLs
- [ ] File management endpoints
  - [ ] `POST /files/upload`
  - [ ] `GET /files/:id`
  - [ ] `DELETE /files/:id`

### Testing & Security
- [ ] Write backend tests
  - [ ] Authentication tests
  - [ ] Application CRUD tests
  - [ ] Authorization tests
- [ ] Implement security measures
  - [ ] Rate limiting
  - [ ] Input validation (express-validator)
  - [ ] SQL injection prevention
  - [ ] CORS configuration
  - [ ] Helmet.js security headers

---

## Phase 3: Frontend Development

### Astro Site Setup
- [x] Initialize Astro project
  - [x] Install Astro 5.x
  - [x] Configure `astro.config.mjs` (site, i18n, integrations)
  - [x] Set up Tailwind CSS 4.x
  - [x] Create custom component library (Card, Button, Container, Section)
- [x] Create layouts
  - [x] `BaseLayout.astro` (main layout with SEO, i18n)
  - [x] `MarketingLayout.astro` (public pages with header/footer)
  - [ ] `DashboardLayout.astro` (authenticated layout) - pending
- [x] Create components
  - [x] `Header.astro`
  - [x] `Footer.astro`
  - [x] `Navigation.astro` (desktop & mobile with language switcher)
  - [x] `LanguageSwitcher.astro`
  - [x] `Card.astro`, `Button.astro`, `Container.astro`, `Section.astro`
- [x] Implement custom i18n system (no astro-i18next dependency)
  - [x] `utils/i18n.ts` with locale detection and translation
  - [x] Translation files: `public/locales/id/common.json`, `public/locales/en/common.json`
  - [x] Bilingual routing: `/` (ID), `/en/` (EN)

### Marketing Pages (Static)
- [x] Create content collection system
  - [x] `content/lecturers/*.md` (faculty profiles with frontmatter)
  - [ ] `content/programs/*.md` (program descriptions) - pending
  - [ ] `content/about/*.md` (about campus) - pending
  - [ ] `content/admissions/*.md` (admission requirements) - pending
- [x] Build static pages (bilingual)
  - [x] Homepage (`pages/index.astro`, `pages/en/index.astro`)
  - [x] About page (`pages/about.astro`, `pages/en/about.astro`)
  - [x] Lecturers listing (`pages/lecturers/index.astro`, `pages/en/lecturers/index.astro`)
  - [x] Lecturer profiles (`pages/lecturers/[slug].astro`, `pages/en/lecturers/[slug].astro`)
  - [ ] Programs listing (`pages/programs/index.astro`) - pending
  - [ ] Program detail pages (`pages/programs/[slug].astro`) - pending
  - [ ] Admissions page (`pages/admissions.astro`) - pending
  - [ ] News listing/detail pages - pending
  - [ ] Contact page (`pages/contact.astro`) - pending

### Authentication Pages
- [ ] Build auth UI
  - [ ] Login page (`pages/login.astro`)
  - [ ] Register page (`pages/register.astro`)
  - [ ] Client-side auth utilities (`src/scripts/auth.js`)
  - [ ] Protected route wrapper

### Application Portal
- [ ] Dashboard
  - [ ] `pages/dashboard.astro` (user dashboard)
  - [ ] Display user info
  - [ ] List user's applications
  - [ ] Show application status
- [ ] Application form
  - [ ] `pages/apply.astro` (application form)
  - [ ] Form validation (client-side)
  - [ ] File upload UI
  - [ ] Draft auto-save
  - [ ] Submit handler
- [ ] Admin interface
  - [ ] `pages/admin/applications.astro` (list all applications)
  - [ ] Application review UI
  - [ ] Approve/reject actions
  - [ ] User management interface

---

## Phase 4: BFF Layer (Cloudflare Workers)

### Authentication Handlers
- [ ] Google OIDC flow
  - [ ] `functions/auth/google/login.js` (initiate OIDC)
  - [ ] `functions/auth/google/callback.js` (handle callback)
  - [ ] Token verification
  - [ ] Cookie management
- [ ] Traditional auth
  - [ ] `functions/auth/login.js` (email/password)
  - [ ] `functions/auth/register.js`
  - [ ] `functions/auth/logout.js`

### API Proxy Functions
- [ ] Application handlers
  - [ ] `functions/applications/submit.js`
  - [ ] `functions/applications/list.js`
  - [ ] `functions/applications/status.js`
- [ ] User handlers
  - [ ] `functions/users/me.js`
  - [ ] `functions/users/update.js`
- [ ] File upload handlers
  - [ ] `functions/files/upload.js`

### BFF Utilities
- [ ] Cookie parser
- [ ] Token extraction
- [ ] Error handling
- [ ] Rate limiting (Cloudflare Workers)

---

## Phase 5: Shared Code

### TypeScript Types
- [ ] Create shared types
  - [ ] `shared/types/User.ts`
  - [ ] `shared/types/Application.ts`
  - [ ] `shared/types/Auth.ts`
  - [ ] `shared/types/index.ts` (exports)

### Constants
- [ ] Define shared constants
  - [ ] `shared/constants/applicationStatus.ts`
  - [ ] `shared/constants/userRoles.ts`
  - [ ] `shared/constants/programs.ts`

### Validators
- [ ] Create validation schemas
  - [ ] `shared/validators/applicationSchema.ts`
  - [ ] `shared/validators/userSchema.ts`
  - [ ] Can use Zod or Joi

---

## Phase 6: Deployment & DevOps

### Frontend Deployment
- [ ] Configure Cloudflare Workers
  - [ ] Create `wrangler.toml`
  - [ ] Set up environment variables
  - [ ] Test Workers locally
- [ ] Deploy to Cloudflare Pages
  - [ ] Build Astro site
  - [ ] Deploy static assets
  - [ ] Configure custom domain
  - [ ] Set up SSL/TLS

### Backend Deployment
- [ ] Prepare VPS
  - [ ] Clone repository to VPS
  - [ ] Install dependencies
  - [ ] Set up `.env` file
  - [ ] Run database migrations
  - [ ] Configure pm2
  - [ ] Set up nginx reverse proxy
- [ ] Configure SSL
  - [ ] Install certbot
  - [ ] Generate SSL certificates
  - [ ] Configure nginx for HTTPS

### CI/CD
- [ ] Test GitHub Actions workflows
  - [ ] Test frontend deployment
  - [ ] Test backend deployment
  - [ ] Verify path-based triggers
- [ ] Set up monitoring
  - [ ] Cloudflare Analytics
  - [ ] VPS monitoring (pm2, nginx logs)
  - [ ] Database monitoring

---

## Phase 7: Testing & Polish

### Testing
- [ ] End-to-end testing
  - [ ] User registration flow
  - [ ] Google OIDC login flow
  - [ ] Email/password login flow
  - [ ] Application submission
  - [ ] Admin review workflow
  - [ ] File upload/download
- [ ] Performance testing
  - [ ] Load testing (simulate 300 users)
  - [ ] Check Cloudflare Workers metrics
  - [ ] Database query optimization
- [ ] Security audit
  - [ ] Test XSS protection
  - [ ] Test CSRF protection
  - [ ] Verify JWT expiration
  - [ ] Test rate limiting
  - [ ] Check file upload security

### Documentation
- [ ] API documentation
  - [ ] Document all endpoints
  - [ ] Request/response examples
  - [ ] Authentication requirements
- [ ] User guides
  - [ ] Registrant guide (how to apply)
  - [ ] Admin guide (how to review applications)
- [ ] Developer documentation
  - [ ] Local development setup
  - [ ] Deployment process
  - [ ] Troubleshooting guide

### Polish
- [ ] UI/UX improvements
  - [ ] Mobile responsiveness
  - [ ] Loading states
  - [ ] Error messages
  - [ ] Success notifications
- [ ] Accessibility
  - [ ] ARIA labels
  - [ ] Keyboard navigation
  - [ ] Screen reader testing
- [ ] SEO optimization
  - [ ] Meta tags
  - [ ] Sitemap
  - [ ] robots.txt
  - [ ] Open Graph tags

---

## Phase 8: Launch Preparation

### Pre-launch Checklist
- [ ] Environment variables verified
  - [ ] Production Cloudflare API tokens
  - [ ] Google OAuth credentials (production)
  - [ ] JWT secrets (strong, unique)
  - [ ] Database credentials
- [ ] Backups configured
  - [ ] Database backup script
  - [ ] File storage backup
  - [ ] Automated daily backups
- [ ] Monitoring set up
  - [ ] Error tracking
  - [ ] Performance monitoring
  - [ ] Uptime monitoring
- [ ] Security review
  - [ ] SSL certificates valid
  - [ ] Security headers configured
  - [ ] Rate limiting tested
  - [ ] DDoS protection enabled

### Soft Launch
- [ ] Test with small group (5-10 users)
- [ ] Collect feedback
- [ ] Fix critical bugs
- [ ] Verify email notifications work
- [ ] Test under realistic load

### Full Launch
- [ ] Announce to students
- [ ] Monitor traffic and errors
- [ ] Be ready for support requests
- [ ] Track application submissions

---

## Future Enhancements (Post-Launch)

### Phase 9: Advanced Features
- [ ] Email verification for local accounts
- [ ] Password reset flow
- [ ] Multi-factor authentication (2FA)
- [ ] Advanced RBAC (reviewers, super-admin)
- [ ] Application workflow stages
- [ ] Email notifications for status updates
- [ ] Admin analytics dashboard
- [ ] Document preview in browser
- [ ] Application deadline management
- [ ] Bulk application processing
- [ ] Export to CSV/Excel
- [ ] Real-time notifications (WebSocket)
- [ ] Audit logs
- [ ] Integration with student information system

### Phase 10: Optimization
- [ ] Implement caching (BFF layer)
- [ ] Request batching
- [ ] Database indexing
- [ ] Image optimization
- [ ] Code splitting
- [ ] Service Worker for offline support

---

## Notes

### Priority Levels
- **P0 (Critical)**: Must have for launch
- **P1 (High)**: Important but can be added shortly after launch
- **P2 (Medium)**: Nice to have, can be planned for future
- **P3 (Low)**: Future enhancement, not urgent

### Estimated Timeline
- Phase 1-2: 1-2 weeks (Backend foundation)
- Phase 3-4: 2-3 weeks (Frontend + BFF)
- Phase 5: 3-5 days (Shared code)
- Phase 6: 1 week (Deployment)
- Phase 7: 1-2 weeks (Testing & polish)
- Phase 8: 1 week (Launch prep)

**Total: 7-10 weeks** for MVP launch

### Team Requirements
- 1 Full-stack developer (can handle all phases)
- OR
- 1 Frontend developer + 1 Backend developer (parallel work)

### Success Metrics
- [ ] 300 registrants successfully submit applications
- [ ] Zero security incidents
- [ ] 99% uptime
- [ ] <2s page load time
- [ ] <1.5% BFF traffic usage (well under limit)
- [ ] $5-10/month hosting cost maintained
