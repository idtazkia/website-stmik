# Frontend - Campus Website

> Static site + BFF layer built with Astro and Cloudflare Workers

## Overview

This frontend handles:
- **Static marketing pages** (programs, about, admissions)
- **Application portal** (login, register, apply, dashboard)
- **BFF layer** (Cloudflare Workers for API proxy and authentication)

**Deployment:**
- Static site → Cloudflare Pages (free, unlimited bandwidth)
- BFF functions → Cloudflare Workers (free, 100k req/day)

---

## Tech Stack

- **Framework:** Astro 4.x (static site generation)
- **Styling:** Tailwind CSS 3.x
- **Components:** Astro components + islands
- **BFF Runtime:** Cloudflare Workers (serverless functions)
- **Content:** Markdown (MDX support)
- **Type Safety:** TypeScript (strict mode)

---

## Directory Structure

```
frontend/
├── src/
│   ├── content/                  # Markdown content
│   │   ├── programs/             # Program descriptions
│   │   │   ├── computer-science.md
│   │   │   ├── business.md
│   │   │   └── engineering.md
│   │   ├── about/                # About campus
│   │   │   ├── history.md
│   │   │   ├── vision-mission.md
│   │   │   └── facilities.md
│   │   └── admissions/           # Admission info
│   │       ├── requirements.md
│   │       ├── process.md
│   │       └── calendar.md
│   │
│   ├── pages/                    # Astro pages (routes)
│   │   ├── index.astro          # Homepage
│   │   ├── programs/
│   │   │   ├── index.astro      # Programs listing
│   │   │   └── [slug].astro     # Program detail page
│   │   ├── about.astro
│   │   ├── admissions.astro
│   │   ├── contact.astro
│   │   ├── login.astro          # Login page
│   │   ├── register.astro       # Register page
│   │   ├── dashboard.astro      # User dashboard
│   │   ├── apply.astro          # Application form
│   │   └── admin/
│   │       └── applications.astro  # Admin panel
│   │
│   ├── layouts/                  # Page layouts
│   │   ├── BaseLayout.astro     # Base HTML structure
│   │   ├── MarketingLayout.astro # For static pages
│   │   └── DashboardLayout.astro # For authenticated pages
│   │
│   ├── components/               # Reusable components
│   │   ├── Header.astro         # Site header
│   │   ├── Footer.astro         # Site footer
│   │   ├── Navigation.astro     # Navigation menu
│   │   ├── Hero.astro           # Hero section
│   │   ├── ProgramCard.astro    # Program display card
│   │   ├── ApplicationForm.astro # Application form
│   │   └── StatusBadge.astro    # Application status badge
│   │
│   ├── scripts/                  # Client-side JavaScript
│   │   ├── auth.ts              # Authentication utilities
│   │   ├── api.ts               # API client
│   │   └── form-validation.ts   # Form validation logic
│   │
│   ├── styles/                   # Global styles
│   │   ├── global.css           # Global CSS
│   │   └── tailwind.css         # Tailwind directives
│   │
│   └── env.d.ts                  # TypeScript environment types
│
├── functions/                    # Cloudflare Workers (BFF)
│   ├── auth/
│   │   ├── google/
│   │   │   ├── login.ts         # Initiate Google OIDC
│   │   │   └── callback.ts      # Handle Google callback
│   │   ├── login.ts             # Email/password login
│   │   ├── register.ts          # User registration
│   │   └── logout.ts            # Logout
│   │
│   ├── applications/
│   │   ├── submit.ts            # Submit application
│   │   ├── list.ts              # Get user's applications
│   │   └── status.ts            # Get application status
│   │
│   ├── users/
│   │   ├── me.ts                # Get current user
│   │   └── update.ts            # Update user profile
│   │
│   ├── files/
│   │   └── upload.ts            # File upload handler
│   │
│   └── _middleware/              # BFF middleware
│       ├── auth.ts              # Authentication check
│       ├── cors.ts              # CORS headers
│       └── error-handler.ts     # Error handling
│
├── public/                       # Static assets
│   ├── images/
│   │   ├── logo-stmik-tazkia.svg  # Campus logo
│   │   └── hero-bg.jpg          # Hero background
│   ├── favicon.ico
│   └── robots.txt
│
├── .dev.vars.example             # Environment variables template
├── .dev.vars                     # Local environment (git-ignored)
├── wrangler.toml                 # Cloudflare Workers config
├── astro.config.mjs              # Astro configuration
├── tailwind.config.mjs           # Tailwind configuration
├── tsconfig.json                 # TypeScript configuration
├── package.json
└── README.md                     # This file
```

---

## Content Organization

### Static Marketing Content (Markdown)

Content files are written in Markdown and stored in `src/content/`.

**Example: `src/content/programs/computer-science.md`**
```markdown
---
title: "Computer Science"
slug: "computer-science"
description: "Become a software engineer"
duration: "4 years"
degree: "Bachelor of Science"
---

# Computer Science Program

Learn programming, algorithms, and software engineering...
```

### Dynamic Pages (Astro)

Pages that require authentication or dynamic data use Astro components.

**Example: `src/pages/dashboard.astro`**
```astro
---
import DashboardLayout from '../layouts/DashboardLayout.astro';
// Server-side code runs on Cloudflare Workers
const user = await getCurrentUser(Astro.cookies);
---

<DashboardLayout title="Dashboard">
  <h1>Welcome, {user.name}</h1>
  <!-- Dashboard content -->
</DashboardLayout>
```

---

## Authentication Flow

### Client-Side (Browser)

```typescript
// src/scripts/auth.ts
class AuthClient {
  async login(email: string, password: string) {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include', // Send cookies
      body: JSON.stringify({ email, password })
    });

    if (response.ok) {
      window.location.href = '/dashboard';
    }
  }

  async loginWithGoogle() {
    window.location.href = '/api/auth/google/login';
  }
}
```

### Server-Side (BFF - Cloudflare Workers)

```typescript
// functions/auth/login.ts
export async function onRequestPost(context) {
  const { request, env } = context;
  const { email, password } = await request.json();

  // Call Express.js backend
  const response = await fetch(`${env.BACKEND_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });

  const { token, user } = await response.json();

  // Set HttpOnly cookie
  return new Response(JSON.stringify({ user }), {
    headers: {
      'Set-Cookie': `token=${token}; HttpOnly; Secure; SameSite=Strict; Path=/; Max-Age=604800`,
      'Content-Type': 'application/json'
    }
  });
}
```

---

## Development Workflow

### Install Dependencies

```bash
cd frontend
npm install
```

### Local Development

```bash
# Start Astro dev server
npm run dev

# Runs on http://localhost:4321
# Hot module reload enabled
# Cloudflare Workers run locally via Wrangler
```

### Build

```bash
# Build for production
npm run build

# Output: dist/ directory (static HTML/CSS/JS)
```

### Preview Production Build

```bash
# Test production build locally
npm run preview
```

### Deploy

```bash
# Deploy to Cloudflare Pages + Workers
npm run deploy

# Or deploy separately:
npx wrangler pages deploy dist  # Deploy static site
npx wrangler deploy            # Deploy Workers (BFF)
```

---

## Environment Variables

### Local Development (`.dev.vars`)

```bash
# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Backend API
BACKEND_URL=http://localhost:3000

# Application
APP_URL=http://localhost:4321
```

### Production (Cloudflare Dashboard)

Set these in Cloudflare Pages → Settings → Environment Variables:
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `BACKEND_URL` (https://api.youruni.edu)
- `APP_URL` (https://youruni.edu)

---

## Styling

### Tailwind CSS

This project uses Tailwind CSS for styling.

**Common classes:**
```html
<!-- Container -->
<div class="container mx-auto px-4">

<!-- Button -->
<button class="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded">

<!-- Card -->
<div class="bg-white shadow-lg rounded-lg p-6">
```

### Custom CSS

Global styles in `src/styles/global.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
.btn-primary {
  @apply bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded transition;
}
```

---

## Type Safety

### Shared Types

Import types from `../shared/types/`:

```typescript
// src/scripts/api.ts
import type { Application } from '../../shared/types/Application';
import type { User } from '../../shared/types/User';

async function getApplications(): Promise<Application[]> {
  const response = await fetch('/api/applications');
  return response.json();
}
```

### Content Collections

Astro content collections provide type-safe access to Markdown files:

```typescript
// src/content/config.ts
import { defineCollection, z } from 'astro:content';

const programs = defineCollection({
  schema: z.object({
    title: z.string(),
    slug: z.string(),
    description: z.string(),
    duration: z.string(),
    degree: z.string(),
  }),
});

export const collections = { programs };
```

---

## Key Pages

### Marketing Pages (Static)

- **Homepage** (`/`): Hero, featured programs, call-to-action
- **Programs** (`/programs`): List of all programs
- **Program Detail** (`/programs/[slug]`): Individual program page
- **About** (`/about`): Campus information
- **Admissions** (`/admissions`): Requirements and process
- **Contact** (`/contact`): Contact form

### Application Portal (Dynamic)

- **Login** (`/login`): Google SSO + email/password
- **Register** (`/register`): Create account
- **Dashboard** (`/dashboard`): User home, application status
- **Apply** (`/apply`): Application form with file upload
- **Admin** (`/admin/applications`): Review applications (staff only)

---

## API Routes (BFF)

All API routes are handled by Cloudflare Workers in `functions/`:

### Authentication
- `POST /api/auth/login` - Email/password login
- `POST /api/auth/register` - Create account
- `GET /api/auth/google/login` - Initiate Google OIDC
- `GET /api/auth/google/callback` - Handle Google callback
- `POST /api/auth/logout` - Logout

### Applications
- `POST /api/applications` - Submit application
- `GET /api/applications` - List user's applications
- `GET /api/applications/:id` - Get single application

### Users
- `GET /api/users/me` - Get current user
- `PATCH /api/users/me` - Update profile

### Files
- `POST /api/files/upload` - Upload document

---

## Testing

### Unit Tests

```bash
npm run test
```

### E2E Tests

```bash
npm run test:e2e
```

### Type Checking

```bash
npm run typecheck
```

---

## Performance

### Build Optimization

- Static site generation (pre-rendered HTML)
- Image optimization (Astro Image)
- CSS purging (Tailwind)
- JavaScript code splitting
- Cloudflare CDN (global distribution)

### Lighthouse Scores Target

- Performance: 90+
- Accessibility: 90+
- Best Practices: 90+
- SEO: 90+

---

## Deployment

### Via GitHub Actions

Automatically deploys when changes are pushed to `main` branch.

See `.github/workflows/deploy-frontend.yml`

### Manual Deployment

```bash
# Build
npm run build

# Deploy to Cloudflare Pages
npx wrangler pages deploy dist --project-name=campus-website

# Deploy Cloudflare Workers (BFF)
npx wrangler deploy
```

---

## Troubleshooting

### Build Errors

```bash
# Clear cache and rebuild
rm -rf node_modules .astro dist
npm install
npm run build
```

### Workers Not Deploying

```bash
# Re-login to Cloudflare
npx wrangler login

# Deploy with verbose output
npx wrangler deploy --verbose
```

### Type Errors

```bash
# Check types
npm run typecheck

# Update types
npm run astro sync
```

---

## Resources

- [Astro Documentation](https://docs.astro.build)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Cloudflare Pages](https://developers.cloudflare.com/pages)
- [Cloudflare Workers](https://developers.cloudflare.com/workers)

---

## License

Proprietary - STMIK Campus Website

**Version:** 1.0
**Last Updated:** 2025-11-19
