# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Campus website for STMIK Tazkia - a bilingual (Indonesian/English) marketing and admission system built with Astro. The project follows a static site + Go backend architecture pattern designed for simplicity and cost-effectiveness ($5/month VPS) while supporting 3,000 leads per admission cycle (300 registrations at 10% conversion).

**Current Status:** Marketing Site Phase (Phase 3 - 30% Complete) - Deployed to Cloudflare Pages at https://stmik.tazkia.ac.id/

**Completed:**
- Bilingual static site (Indonesian & English)
- Homepage, About page, Lecturer profiles
- Responsive design with Tailwind CSS 4.x
- Custom i18n implementation
- SEO optimization (meta tags, sitemap)
- Cloudflare Pages deployment

**In Progress:**
- Additional marketing pages (Programs, Contact, Admissions, News)

## Architecture

### High-Level Structure

```
Static Site + Go Backend Pattern

Browser
  ├── Cloudflare Pages (Astro Static Site) - FREE
  ├── Cloudflare CDN (DDoS Protection) - FREE
  └── VPS ($5/month)
        ├── Nginx (Reverse Proxy + SSL)
        ├── Go Backend (REST API)
        └── PostgreSQL (Database)
```

**Key Architectural Decisions:**
- **Astro Static Site**: SEO-optimized marketing pages with minimal JavaScript
- **Go Backend**: High-performance REST API with minimal resource usage
- **Local PostgreSQL**: <1ms latency, no external database dependencies
- **Cloudflare CDN**: DDoS protection, VPS IP hidden behind proxy
- **Bilingual**: Indonesian (default) and English with custom i18n implementation

### Project Structure

```
website-stmik/
├── frontend/                          # Astro application
│   ├── src/
│   │   ├── components/               # Reusable UI components
│   │   ├── layouts/                  # Page layouts
│   │   ├── pages/                    # Route pages
│   │   │   ├── index.astro          # Homepage (Indonesian)
│   │   │   └── en/index.astro       # Homepage (English)
│   │   ├── utils/                    # Utilities
│   │   │   └── i18n.ts              # Custom i18n implementation
│   │   ├── content/                  # Content collections
│   │   └── styles/                   # Global styles
│   ├── public/
│   │   └── locales/                  # Translation JSON files
│   │       ├── id/common.json       # Indonesian translations
│   │       └── en/common.json       # English translations
│   ├── astro.config.mjs             # Astro configuration
│   └── TODO.md                       # Frontend implementation tasks
├── backend/                           # Go API
│   ├── cmd/                          # Application entry points
│   │   ├── server/                   # Main server
│   │   ├── migrate/                  # Database migrations
│   │   ├── seedtest/                 # Test data seeder
│   │   └── testrunner/               # Test orchestrator
│   ├── internal/                     # Private application code
│   │   ├── auth/                     # Session management
│   │   ├── config/                   # Configuration
│   │   ├── handler/                  # HTTP handlers
│   │   ├── model/                    # Data models + queries
│   │   └── storage/                  # File storage
│   ├── web/                          # Web assets
│   │   ├── static/                   # CSS, JS
│   │   └── templates/                # Templ files
│   ├── test/                         # Test files
│   │   ├── e2e/                      # Playwright E2E tests
│   │   └── testutil/                 # Test utilities
│   └── migrations/                   # SQL migration files
├── docs/
│   ├── ARCHITECTURE.md               # Technical design details
│   └── DEPLOYMENT.md                 # Deployment guide
├── TODO.md                            # High-level project overview
└── CLAUDE.md                          # This file
```

## Development Commands

### Frontend Development

```bash
# Development
cd frontend
npm run dev                   # Start dev server at http://localhost:4321

# Build & Preview
npm run build                 # Build for production
npm run preview               # Preview production build

# Code Quality
npm run typecheck             # Run TypeScript checks
npm run lint                  # Run ESLint
npm run lint:fix              # Auto-fix ESLint issues
npm run format                # Format with Prettier
npm run format:check          # Check formatting
```

### Backend Commands

```bash
cd backend

# Install dependencies
make deps                     # Go, npm, and templ

# Development
make dev                      # Run development server
make templ                    # Generate templates
make css                      # Build CSS

# Database
go run ./cmd/migrate up       # Run migrations
go run ./cmd/migrate down     # Rollback migration

# Testing (uses testcontainers - like mvn clean test)
make clean-test               # Full test suite with testcontainers
make test                     # Run unit tests only
make test-e2e                 # Run E2E tests (requires running server)

# Build
make build                    # Build for current platform
make build-linux              # Build for Linux deployment
```

## Tech Stack

| Component | Technology | Notes |
|-----------|-----------|-------|
| **Frontend** | Astro 5.x + Tailwind CSS 4.x | Static site generation, component islands |
| **Styling** | Tailwind CSS + custom design system | Brand colors: primary (blue #194189), secondary (orange #EE7B1D) |
| **i18n** | Custom implementation | See frontend/src/utils/i18n.ts |
| **Backend** | Go 1.25+ | REST API + HTMX, runs on VPS |
| **Templates** | Templ | Type-safe, compiled HTML templates |
| **Interactivity** | HTMX + Alpine.js CSP | Server-driven UI, CSP-compliant |
| **Database** | PostgreSQL 18 + pgx/v5 | Local on VPS, <1ms latency |
| **Testing** | Testcontainers + Playwright | Automated E2E with real PostgreSQL |
| **Reverse Proxy** | Nginx + Let's Encrypt | SSL termination, rate limiting |
| **Deployment** | Cloudflare Pages + VPS | Frontend auto-deploy, backend via GitHub Actions |

## Internationalization (i18n)

The project uses a **custom i18n implementation** (not astro-i18next dependency).

### How it Works

1. **Route Structure:**
   - Indonesian (default): `/` `/programs` `/about`
   - English: `/en/` `/en/programs` `/en/about`

2. **Translation Files:**
   - Located in `frontend/public/locales/{locale}/common.json`
   - Structure: Nested JSON with dot notation keys

3. **Usage in Components:**
   ```typescript
   import { getLocaleFromUrl, t, localizePath } from '../utils/i18n';

   const locale = getLocaleFromUrl(Astro.url);
   const translate = (key: string) => t(locale, key);
   const localPath = (path: string) => localizePath(locale, path);
   ```

4. **Adding New Translations:**
   - Add keys to both `id/common.json` and `en/common.json`
   - Use dot notation: `"nav.home"`, `"button.submit"`

## Design System

### Brand Colors

```css
/* Primary: Blue (#194189 from logo) */
primary-50, primary-100, ..., primary-900, primary-950

/* Secondary: Orange (#EE7B1D from logo) */
secondary-50, secondary-100, ..., secondary-900, secondary-950

/* Backgrounds */
background="white"    /* White background */
background="gray"     /* Light gray (gray-50) */
background="gradient" /* Primary gradient (primary-50 to secondary-50) */
```

### Component Patterns

**Reusable Components:**
- `<Container>` - Content width container with responsive padding
- `<Section>` - Page section with configurable background/spacing
- `<Card>` - Content card with hover effects
- `<Button>` - Styled buttons with variants (primary, secondary, outline)

**Layouts:**
- `BaseLayout.astro` - HTML structure, SEO, language alternates
- `MarketingLayout.astro` - Header + Footer wrapper for public pages

## Authentication System (Planned)

**Hybrid OIDC + Traditional Auth:**
- **Registrants**: Google OIDC or email/password
- **Staff**: Google OIDC only (domain check: `@youruni.edu`)

**Security Features:**
- HttpOnly cookies (XSS protection)
- JWT tokens (7-day expiration)
- Rate limiting via Nginx
- bcrypt password hashing

See `docs/ARCHITECTURE.md` for authentication flows.

## Database Schema (Planned)

**Database:** PostgreSQL 18 (local on VPS)

**Tables:**
- `users` - User accounts (registrants + staff)
- `applications` - Application submissions
- `sessions` - Optional session storage

**Go Database Driver:**
- Using `github.com/jackc/pgx/v5` (recommended PostgreSQL driver for Go)
- Connection pooling via `pgxpool`
- Connection string: `postgresql://campus_app:password@localhost:5432/campus`

See `docs/ARCHITECTURE.md#database-schema` for detailed schema.

## Important Guidelines

### When Adding New Pages

1. Create Indonesian version in `frontend/src/pages/`
2. Create English version in `frontend/src/pages/en/`
3. Use `MarketingLayout` for public pages
4. Import i18n utilities: `getLocaleFromUrl`, `t`, `localizePath`
5. Add translations to both locale files

### When Modifying Translations

1. Update both `id/common.json` and `en/common.json`
2. Maintain consistent key structure
3. Use semantic keys (e.g., `nav.home` not `header_link_1`)

### Security Considerations

- Never commit `.env` files
- Never commit credentials or API keys
- Use environment variables for sensitive data
- Implement input validation for all user inputs
- Follow OWASP top 10 security practices

### Code Style

- Use TypeScript strict mode
- Follow ESLint configuration
- Format with Prettier before commits
- Prefer functional components in Astro
- Use semantic HTML elements

### Alpine.js CSP Guidelines (Backend)

The backend uses **Alpine.js CSP build** for Content Security Policy compliance. This is **mandatory** - do not use standard Alpine.js patterns.

**NEVER write inline expressions in HTML:**
```html
<!-- ❌ WRONG - Will not work with CSP build -->
<div x-data="{ open: false }">
    <button @click="open = !open">Toggle</button>
</div>
```

**ALWAYS use Alpine.data() and reference by name:**
```html
<!-- ✅ CORRECT - CSP compliant -->
<div x-data="dropdown">
    <button @click="toggle">Toggle</button>
</div>
```

```javascript
Alpine.data('dropdown', () => ({
    open: false,
    toggle() { this.open = !this.open }
}))
```

See `backend/README.md` for full Alpine CSP guidelines.

## Implementation Roadmap

**Current Phase: Phase 3 - Marketing Site (30% Complete)**

See `TODO.md` for the implementation plan:
- Phase 1: Project Setup [PARTIAL - Frontend deployed, backend infrastructure deferred]
- Phase 2: Backend Development (Go) [DEFERRED - not started]
- Phase 3: Frontend Development [IN PROGRESS - Homepage, About, Lecturers done]
- Phase 4-6: Deployment, Testing, Launch [NOT STARTED]

## Key Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Project overview, quick start, architecture summary |
| `docs/ARCHITECTURE.md` | Technical design, auth flows, database schema, security |
| `docs/DEPLOYMENT.md` | Step-by-step deployment guide, monitoring, troubleshooting |
| `TODO.md` | Complete implementation checklist with 8 phases |
| `CLAUDE.md` | This file - guidance for Claude Code |

## Target Metrics

- 3,000 leads per admission cycle (300 registrations at 10% conversion)
- <2s page load time
- 99% uptime
- $5/month hosting cost (VPS)
- <30% VPS resource usage (RAM, CPU)
- <1% disk usage
- Zero security incidents

## Common Tasks

### Adding a New Static Page

```bash
# 1. Create Indonesian page
touch frontend/src/pages/new-page.astro

# 2. Create English page
touch frontend/src/pages/en/new-page.astro

# 3. Add translations to both locale files
# 4. Add navigation link if needed
# 5. Test both language versions
```

### Updating Brand Colors

Colors are defined in Tailwind configuration. Use the predefined `primary-*` and `secondary-*` scales throughout the application.

### Working with Content Collections

Content collections are configured in `frontend/src/content/config.ts`. Use them for structured content like programs, news articles, or blog posts.

## Notes

- Backend follows standard Go project layout with `internal/`, `web/`, `test/` directories.
- Testing uses testcontainers for PostgreSQL - run `make clean-test` (like `mvn clean test`).
- The project uses a custom i18n implementation (no external dependencies).
- Frontend deployment uses Cloudflare Pages auto-deploy (triggered on git push to GitHub).
- Backend deployment uses GitHub Actions to deploy Go binary to VPS.
- Database uses PostgreSQL 18 running locally on VPS (<1ms latency).
- File uploads stored on VPS local disk.
- Simple architecture: single VPS ($5/month) for backend + database.
- Can scale to 100,000+ leads on same VPS before needing upgrade.
