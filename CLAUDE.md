# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Campus website for STMIK Tazkia - a bilingual (Indonesian/English) marketing and admission system built with Astro. The project follows a hybrid static site + BFF (Backend-For-Frontend) architecture pattern designed for cost-effectiveness ($5-10/month) while supporting 300 registrants per admission cycle.

**Current Status:** Phase 2 Complete - Frontend foundation with bilingual support and brand colors implemented.

## Architecture

### High-Level Structure

```
Hybrid Static Site + BFF Pattern

Browser
  ├── Cloudflare Pages (Astro Static Site) - FREE
  ├── Cloudflare Workers (BFF Layer) - FREE [Planned]
  └── Express.js Backend (VPS) - $5-10/mo [Planned]
            └── PostgreSQL Database [Planned]
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
│   └── astro.config.mjs             # Astro configuration
├── backend/                           # Express.js API [Planned]
├── shared/                            # Shared TypeScript types [Planned]
├── docs/
│   ├── ARCHITECTURE.md               # Technical design details
│   └── DEPLOYMENT.md                 # Deployment guide
└── TODO.md                            # 8-phase implementation checklist
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
| **Styling** | Tailwind CSS + custom design system | Brand colors: primary (teal), secondary (orange) |
| **i18n** | Custom implementation | See frontend/src/utils/i18n.ts |
| **BFF** | Cloudflare Workers | [Planned] Auth handlers, API proxy |
| **Backend** | Express.js + PostgreSQL | [Planned] REST API, file storage |
| **Deployment** | Cloudflare Pages | [Planned] GitHub Actions CI/CD |

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
/* Primary: Teal/Green */
primary-50, primary-100, ..., primary-900, primary-950

/* Secondary: Orange */
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

**Tables:**
- `users` - User accounts (registrants + staff)
- `applications` - Application submissions
- `sessions` - Optional session storage

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

**Current Phase: Phase 2 Complete**

See `TODO.md` for the complete 8-phase implementation plan:
- Phase 1: Project Setup [PARTIAL - Frontend initialized]
- Phase 2: Backend Development [PLANNED]
- Phase 3: Frontend Development [IN PROGRESS - Static pages done]
- Phase 4: BFF Layer [PLANNED]
- Phase 5: Shared Code [PLANNED]
- Phase 6: Deployment & DevOps [PLANNED]
- Phase 7: Testing & Polish [PLANNED]
- Phase 8: Launch Preparation [PLANNED]

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

- 300 registrants per admission cycle
- <2s page load time
- 99% uptime
- $5-10/month hosting cost
- <1.5% Cloudflare Workers usage (well under 100k req/day limit)
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
- Astro-i18next is installed in root package.json but NOT used. The project uses a custom i18n implementation.
- GitHub Actions workflows for CI/CD will be added in Phase 6.
- File uploads and Cloudflare R2 integration planned for Phase 2-4.
