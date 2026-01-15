# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Campus website for STMIK Tazkia - a bilingual (Indonesian/English) marketing and admission system built with Astro. The project follows a hybrid static site + BFF (Backend-For-Frontend) architecture pattern designed for cost-effectiveness ($0/month on free tiers) while supporting 3,000 leads per admission cycle (300 registrations at 10% conversion).

**Current Status:** Marketing Site Phase (Phase 3 - 30% Complete) - Deployed to Cloudflare Pages at https://dev.stmik.tazkia.ac.id/

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
Hybrid Static Site + BFF Pattern (Fully Serverless)

Browser
  ├── Cloudflare Pages (Astro Static Site) - FREE
  ├── Cloudflare Workers (BFF Layer) - FREE [Planned]
  ├── Cloudflare R2 (File Storage) - FREE [Planned]
  └── CockroachDB Serverless (Database) - FREE [Planned]
```

**Key Architectural Decisions:**
- **Astro Static Site**: SEO-optimized marketing pages with minimal JavaScript
- **BFF Pattern**: HttpOnly cookies for security, edge-level rate limiting
- **Serverless Edge**: Cloudflare Workers handle auth/API proxying (100k req/day free)
- **Monorepo Structure**: Shared TypeScript types between frontend/backend
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
├── backend/                           # Express.js API [Deferred to Phase 2]
│   └── TODO.md                       # Backend implementation plan
├── shared/                            # Shared TypeScript types [Deferred to Phase 5]
│   └── TODO.md                       # Shared code plan
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

### Future Monorepo Commands (when backend is added)

```bash
# Root level
npm install                   # Install all workspace dependencies
npm run dev                   # Run all services
npm run dev -w frontend       # Frontend only
npm run dev -w backend        # Backend only
npm run build                 # Build all packages

# Backend specific
cd backend
npm run migrate               # Run database migrations
npm run migrate:rollback      # Rollback last migration
```

## Tech Stack

| Component | Technology | Notes |
|-----------|-----------|-------|
| **Frontend** | Astro 5.x + Tailwind CSS 4.x | Static site generation, component islands |
| **Styling** | Tailwind CSS + custom design system | Brand colors: primary (blue #194189), secondary (orange #EE7B1D) |
| **i18n** | Custom implementation | See frontend/src/utils/i18n.ts |
| **BFF** | Cloudflare Workers | [Planned] Auth handlers, API proxy (100k req/day free) |
| **Database** | CockroachDB Serverless | [Planned] PostgreSQL-compatible, 50M RUs/month free |
| **File Storage** | Cloudflare R2 | [Planned] 10GB free tier |
| **Deployment** | Cloudflare Pages | Auto-deploy on git push (Cloudflare integration) |

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
- Rate limiting via Cloudflare Workers
- bcrypt password hashing

See `docs/ARCHITECTURE.md` for authentication flows.

## Database Schema (Planned)

**Database:** CockroachDB Serverless (PostgreSQL wire-compatible)

**Tables:**
- `users` - User accounts (registrants + staff)
- `applications` - Application submissions
- `sessions` - Optional session storage

**CockroachDB Notes:**
- Standard CRUD, JOINs, indexes work identically to PostgreSQL
- Use SERIAL for auto-incrementing IDs
- Connection string: `postgresql://user:pass@host:26257/db?sslmode=verify-full`

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

## Implementation Roadmap

**Current Phase: Phase 3 - Marketing Site (30% Complete)**

See `TODO.md` for the complete 8-phase implementation plan:
- Phase 1: Project Setup [PARTIAL - Frontend deployed, backend infrastructure deferred]
- Phase 2: Backend Development [DEFERRED - not started]
- Phase 3: Frontend Development [IN PROGRESS - Homepage, About, Lecturers done; Programs, Contact, Admissions, News pending]
- Phase 4: BFF Layer [DEFERRED - not started]
- Phase 5: Shared Code [DEFERRED - not started]
- Phase 6: Deployment & DevOps [PARTIAL - Cloudflare Pages auto-deploy working]
- Phase 7: Testing & Polish [NOT STARTED]
- Phase 8: Launch Preparation [NOT STARTED]

Estimated time to MVP: 7-10 weeks

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
- $0/month hosting cost (fully on free tiers)
- <5% Cloudflare Workers usage (2.2% of 100k req/day limit)
- <5% CockroachDB usage (1% of 50M RUs/month limit)
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

- The project currently uses only frontend dependencies. Backend dependencies will be added in Phase 2.
- The project uses a custom i18n implementation (no external dependencies).
- Deployment uses Cloudflare Pages auto-deploy (triggered on git push to GitHub).
- Database uses CockroachDB Serverless (PostgreSQL-compatible, 50M RUs/month free tier).
- File uploads use Cloudflare R2 (10GB free tier).
- Fully serverless architecture - no VPS required, $0/month hosting cost.
- Can scale to 100,000+ leads on free tier before needing paid plans.
