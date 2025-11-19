# Campus Website - Architecture Design

## Project Overview

Campus website for marketing and admission processing, serving:
- **Registrants** (prospective students): Submit applications via Google SSO or email/password
- **Marketing Staff**: Manage applications and content via Google SSO only

---

## Architecture Pattern

**Hybrid Static Site + BFF (Backend-For-Frontend) Pattern**

```mermaid
graph TB
    subgraph "User Layer"
        Browser[Browser]
    end

    subgraph "Edge Layer - Cloudflare (Free)"
        StaticSite[Static Marketing Site<br/>Astro + Markdown<br/>Cloudflare Pages]
        BFF[BFF Layer<br/>Cloudflare Workers<br/>API Routes]
    end

    subgraph "VPS Layer - Local Cloud Provider"
        Backend[Express.js Backend<br/>REST API]
        Database[(PostgreSQL<br/>Database)]
        FileStorage[File Storage<br/>Documents/Uploads]
    end

    Browser -->|HTTPS| StaticSite
    Browser -->|API Calls<br/>Cookie: token| BFF
    BFF -->|Authorization: Bearer JWT| Backend
    Backend --> Database
    Backend --> FileStorage

    StaticSite -.->|Links to| BFF
```

---

## Technology Stack

| Component | Technology | Hosting | Cost |
|-----------|-----------|---------|------|
| **Static Site** | Astro + Markdown | Cloudflare Pages | Free |
| **BFF Layer** | Cloudflare Workers | Cloudflare | Free (100k req/day) |
| **Backend API** | Express.js (Node.js) | Local VPS Provider | $5-10/mo |
| **Database** | PostgreSQL | Local VPS Provider | Included in VPS |
| **File Storage** | Cloudflare R2 / VPS | Cloudflare / VPS | Free tier / Included |
| **Build/Deploy** | GitHub Actions | GitHub | Free |

**Total Monthly Cost: $5-10** (VPS only)

---

## Authentication System

### Hybrid OIDC + Traditional Authentication

```mermaid
graph LR
    subgraph "Registrants"
        R1[Google OIDC<br/>Recommended]
        R2[Email/Password<br/>Traditional]
    end

    subgraph "Marketing Staff"
        S1[Google OIDC Only<br/>@youruni.edu required]
    end

    R1 --> Auth[Authentication System]
    R2 --> Auth
    S1 --> Auth

    Auth --> DB[(Users Database<br/>provider: google/local<br/>role: registrant/staff)]
```

### Authentication Methods

#### Registrants (Prospective Students):
- **Option 1: Google OIDC** (recommended)
  - One-click login
  - No password management
  - Email automatically verified
- **Option 2: Email/Password** (traditional)
  - Privacy-conscious option
  - Self-managed credentials

#### Marketing Staff:
- **Google OIDC only**
- Enforced via `@youruni.edu` email domain check
- Automatically assigned `staff` role

### Security Features
- ✅ HttpOnly cookies (XSS protection)
- ✅ JWT tokens managed by BFF layer
- ✅ Tokens never exposed to client JavaScript
- ✅ Same-origin requests (no CORS issues)
- ✅ Rate limiting built-in (Cloudflare)
- ✅ DDoS protection included

---

## Google OIDC Authentication Flow

```mermaid
sequenceDiagram
    participant Browser
    participant BFF as BFF Layer<br/>(Cloudflare Workers)
    participant Google as Google OAuth
    participant Backend as Express.js<br/>Backend
    participant DB as PostgreSQL

    Note over Browser,DB: Login Initiation
    Browser->>BFF: GET /api/auth/google/login
    BFF->>Browser: HTTP 302 Redirect to Google
    Browser->>Google: User logs in with Google
    Google->>Browser: Redirect with authorization code

    Note over Browser,DB: Token Exchange
    Browser->>BFF: GET /callback?code=AUTHCODE
    BFF->>Google: POST /token<br/>{code, client_id, client_secret}
    Google->>BFF: {id_token, access_token}

    Note over BFF: Verify Google ID Token<br/>Extract: {email, name, sub}

    BFF->>Backend: POST /auth/google<br/>{email, name, provider_id}
    Backend->>DB: SELECT/INSERT/UPDATE user
    DB->>Backend: User data {id: 123, email, role}

    Note over Backend: Generate YOUR JWT<br/>{userId: 123, role, email}
    Backend->>BFF: {token: YOUR_JWT, user: {...}}

    Note over BFF: Set HttpOnly Cookie<br/>token=YOUR_JWT
    BFF->>Browser: HTTP 302 /dashboard<br/>Set-Cookie: token=YOUR_JWT

    Note over Browser: Cookie stored<br/>User logged in!
```

---

## Traditional Email/Password Flow

```mermaid
sequenceDiagram
    participant Browser
    participant BFF as BFF Layer<br/>(Cloudflare Workers)
    participant Backend as Express.js<br/>Backend
    participant DB as PostgreSQL

    Browser->>BFF: POST /api/auth/login<br/>{email, password}
    BFF->>Backend: POST /auth/login<br/>{email, password}

    Backend->>DB: SELECT * FROM users<br/>WHERE email = ? AND provider = 'local'
    DB->>Backend: User data + password_hash

    Note over Backend: bcrypt.compare(password, hash)

    Note over Backend: Generate JWT<br/>{userId: 123, email, role}
    Backend->>BFF: {token: JWT, user: {...}}

    Note over BFF: Set HttpOnly Cookie
    BFF->>Browser: {user: {...}}<br/>Set-Cookie: token=JWT

    Note over Browser: Cookie stored<br/>User logged in!
```

---

## Authenticated Request Flow

```mermaid
sequenceDiagram
    participant Browser
    participant BFF as BFF Layer<br/>(Cloudflare Workers)
    participant Backend as Express.js<br/>Backend
    participant DB as PostgreSQL

    Note over Browser: Cookie automatically sent
    Browser->>BFF: GET /api/applications<br/>Cookie: token=YOUR_JWT

    Note over BFF: Extract JWT from cookie<br/>const token = cookies.token

    BFF->>Backend: GET /applications<br/>Authorization: Bearer YOUR_JWT

    Note over Backend: jwt.verify(token, SECRET)<br/>Decode: {userId: 123, role}

    Backend->>DB: SELECT * FROM applications<br/>WHERE user_id = 123
    DB->>Backend: Application records
    Backend->>BFF: [{application data}]
    BFF->>Browser: [{application data}]
```

---

## Session Identification: Two Token Pattern

### Token 1: Google's ID Token (OIDC)
**Used only during login, then discarded**

| Property | Value |
|----------|-------|
| **Issued by** | Google |
| **Signed with** | Google's private key |
| **Verified by** | BFF (using Google's public cert) |
| **Contains** | `sub` (Google user ID), `email`, `name`, `iss`, `aud` |
| **Purpose** | Prove user identity to BFF |
| **Lifetime** | 1 hour (but only used once) |
| **Stored** | NOWHERE (discarded after verification) |

### Token 2: Your Backend's JWT
**Used for all authenticated requests**

| Property | Value |
|----------|-------|
| **Issued by** | YOUR Express.js backend |
| **Signed with** | YOUR `JWT_SECRET` |
| **Verified by** | YOUR backend |
| **Contains** | `userId` (your DB ID), `email`, `role`, `provider` |
| **Purpose** | Session management for all API requests |
| **Lifetime** | 7 days |
| **Stored** | HttpOnly cookie (managed by BFF) |

### Flow Summary

```mermaid
graph TB
    subgraph "Login Phase"
        A[Google ID Token] -->|Verify & Extract| B[User Info]
        B -->|Create/Update User| C[Backend Generates<br/>YOUR JWT]
    end

    subgraph "Session Phase"
        C -->|Store in Cookie| D[Browser]
        D -->|All Requests| E[BFF extracts JWT]
        E -->|Authorization Header| F[Backend verifies<br/>YOUR JWT]
    end

    style A fill:#fbbc04
    style C fill:#34a853
    style F fill:#34a853
```

**Key Insight:**
- **OIDC is for AUTHENTICATION** (proving who you are during login)
- **YOUR JWT is for AUTHORIZATION** (accessing resources after login)

---

## Database Schema

```mermaid
erDiagram
    USERS ||--o{ APPLICATIONS : submits
    USERS {
        serial id PK
        varchar email UK
        varchar name
        varchar password_hash "NULL for Google users"
        varchar provider "google or local"
        varchar provider_id "Google user ID"
        varchar role "registrant or staff"
        boolean email_verified
        timestamp created_at
        timestamp updated_at
    }

    APPLICATIONS {
        serial id PK
        integer user_id FK
        varchar program
        varchar status "pending/approved/rejected"
        timestamp submitted_at
        timestamp updated_at
    }

    SESSIONS {
        serial id PK
        integer user_id FK
        varchar token_hash
        timestamp expires_at
        timestamp created_at
    }
```

---

## DDoS Protection & Cost Safety

### Cloudflare Protection Features

| Feature | Free Tier Limit | Behavior on Excess | Cost Risk |
|---------|----------------|-------------------|-----------|
| **Bandwidth** | Unlimited | N/A | **$0** |
| **Requests** | 100,000/day | Blocked (HTTP 429) | **$0** |
| **Workers Execution** | 100k req/day | Blocked | **$0** |
| **DDoS Protection** | Included | Automatic filtering | **$0** |
| **Bot Detection** | Included | Automatic blocking | **$0** |

### Cost Safety Summary
- ✅ **Maximum DDoS cost: $0** (hard limits prevent overages)
- ✅ VPS protected behind Cloudflare proxy
- ✅ No surprise bills
- ✅ Automatic rate limiting
- ✅ World-class DDoS protection (Cloudflare's core business)

### Comparison with Vercel

| Platform | DDoS Risk | Hard Limits | Max Cost on Attack |
|----------|-----------|-------------|-------------------|
| **Cloudflare Pages + Workers** | ✅ VERY LOW | ✅ Yes | **$0** |
| **Vercel (Free)** | ⚠️ HIGH | ❌ No | **$500+** |
| **Netlify (Free)** | ✅ LOW | ✅ Yes | $0 (stops) |

---

## BFF Traffic Analysis

### Real-World Usage Estimation

All authenticated API requests go through the BFF (Cloudflare Workers) layer. Let's analyze realistic traffic patterns.

### Traffic Breakdown: 300 Registrants + 5 Admins

#### Registrant Journey (30-day admission period)

```
Day 1: Account Creation
- Create account (BFF): 1 request
- Login (BFF): 1 request
Total: 2 requests

Days 2-5: Application Preparation
- Daily login (BFF): 1 request
- View form (BFF): 1 request
- Auto-save drafts (BFF): 5 requests
Total: 7 requests/day × 4 days = 28 requests

Day 6: Application Submission
- Login (BFF): 1 request
- Load form (BFF): 1 request
- Upload documents (BFF): 3 requests
- Submit (BFF): 1 request
Total: 6 requests

Days 7-30: Status Checking
- Login (BFF): 1 request
- Check status (BFF): 1 request
- View updates (BFF): 1 request
Total: 3 requests/day × 24 days = 72 requests

Total per registrant: 2 + 28 + 6 + 72 = 108 requests/30 days
Daily average: 108 ÷ 30 = 3.6 requests/day
```

**300 Registrants:**
- Total over 30 days: 300 × 108 = **32,400 requests**
- Daily average: **1,080 requests/day**

#### Admin Activity (Daily)

```
Per Admin per Day:
- Morning: Login, dashboard, review 10 apps, update statuses
  = 1 + 1 + 1 + 10 + 10 = 23 requests
- Afternoon: Login, dashboard, review 5 apps, notifications
  = 1 + 1 + 5 + 5 = 12 requests

Total per admin: 35 requests/day
```

**5 Admins:**
- Daily total: 5 × 35 = **175 requests/day**
- 30-day total: **5,250 requests**

### Total Traffic Summary

| Scenario | Daily BFF Requests | % of Free Tier | Status |
|----------|-------------------|----------------|--------|
| **Average Day** | 1,255 | **1.3%** | ✅ Very Safe |
| **Peak Day** (50 simultaneous submissions) | 1,425 | **1.4%** | ✅ Very Safe |
| **10x Traffic Spike** | 14,250 | **14.3%** | ✅ Safe |
| **Quiet Period** (post-admission) | 150 | **0.2%** | ✅ Very Safe |

### Scalability Analysis

| User Count | Daily BFF Requests | % of Free Tier | Monthly Cost | Status |
|------------|-------------------|----------------|--------------|--------|
| **300 registrants** | 1,255 | 1.3% | $0 | ✅ Very Safe |
| **1,000 registrants** | 3,775 | 3.8% | $0 | ✅ Very Safe |
| **3,000 registrants** | 11,500 | 11.5% | $0 | ✅ Safe |
| **10,000 registrants** | 38,000 | 38% | $0 | ✅ Safe |
| **27,000 registrants** | 99,000 | 99% | $0 | ⚠️ Near Limit |
| **30,000+ registrants** | 110,000+ | 110%+ | $5/mo* | ✅ Upgrade Available |

*Cloudflare Workers Paid: $5/month for 10 million requests

### Important: Static vs Dynamic Traffic

```mermaid
graph LR
    subgraph "Website Traffic Distribution"
        A[1000 Daily Visitors]
        A -->|95%| B[950 visitors<br/>Browse Marketing Pages<br/>Static Content]
        A -->|5%| C[50 visitors<br/>Login & Apply<br/>BFF Layer]
    end

    B -->|NO BFF| D[Cloudflare Pages CDN<br/>UNLIMITED bandwidth<br/>$0 cost]
    C -->|Through BFF| E[Cloudflare Workers<br/>100k req/day limit<br/>Currently 1.3% used]
```

**Key Insight:**
- **Marketing pages** (programs, about, contact): Unlimited traffic, $0 cost, NO BFF
- **Application portal** (login, submit, status): Only active applicants, goes through BFF
- **Estimated 95%+ of traffic** doesn't touch BFF at all

### Bandwidth Analysis

**BFF Request Size:**
```
Average authenticated request:
- Cookie: ~200 bytes
- Headers: ~500 bytes
- Request body (POST): ~5 KB
- Response body: ~10 KB
Total: ~15 KB per request

Daily bandwidth (1,255 requests):
1,255 × 15 KB = 18.8 MB/day

Monthly bandwidth:
18.8 MB × 30 = 564 MB/month
```

**Cloudflare Workers:** CPU time matters, not bandwidth
**Cloudflare Pages (static):** Unlimited bandwidth ✅

### Optimization Strategies (If Needed)

#### 1. Caching Strategy
```javascript
// Cache application status for 5 minutes
const cacheKey = `status:${userId}`;
const cached = await cache.get(cacheKey);

if (cached) return cached; // Saves BFF request

const fresh = await fetchFromBackend();
await cache.put(cacheKey, fresh, { expirationTtl: 300 });
return fresh;
```

**Impact:**
- Status checks: 72 → 12 requests/user (83% reduction)
- Total traffic: 1,255 → 550 requests/day
- **New usage: 0.55% of free tier**

#### 2. Request Batching
```javascript
// Before: Multiple BFF calls
GET /api/user/profile       // 1 request
GET /api/applications       // 1 request
GET /api/notifications      // 1 request

// After: Single BFF call
GET /api/dashboard          // 1 request (returns all data)
```

**Impact:** ~30% reduction in BFF requests

#### 3. Monitoring & Alerts

Set up Cloudflare Workers analytics to track:
- Daily request count
- Request patterns by endpoint
- Alert if approaching 50% of daily limit (50,000 requests)

### Cost Breakdown

| Component | Traffic Type | Daily Volume | Limit | Cost |
|-----------|-------------|--------------|-------|------|
| **Static Pages** | Marketing content | 100,000+ page views | Unlimited | **$0** |
| **BFF Workers** | Authenticated API | 1,255 requests | 100k/day | **$0** |
| **Backend API** | Database queries | 1,255 requests | VPS capacity | **$5-10/mo** |
| **Total** | | | | **$5-10/mo** |

### When to Upgrade

You'll need to upgrade to Cloudflare Workers Paid ($5/month) when:

- **Daily requests exceed 100,000** (sustained)
- **User count reaches ~27,000 registrants**
- **You want additional features** (longer CPU time, KV storage, Durable Objects)

At that scale, $5/month is negligible compared to the value delivered.

### Verdict for 300 Registrants + 5 Admins

✅ **You're using only 1.3% of the free tier**
✅ **98.7% buffer for traffic spikes**
✅ **Even 10x growth keeps you under 15% usage**
✅ **Zero cost risk - hard limits prevent overages**
✅ **Could scale to 3,000 users and still be at 11.5%**

**Conclusion: No need to worry about BFF traffic limits for your use case.**

---

## Content Management

### Marketing Content (Static)
```mermaid
graph LR
    A[Markdown Files<br/>in Git] -->|Push| B[GitHub]
    B -->|Trigger| C[GitHub Actions]
    C -->|Build Astro| D[Static HTML/CSS/JS]
    D -->|Deploy| E[Cloudflare Pages]
    E -->|Serve| F[Users]
```

**Content Structure:**
```
content/
├── programs/
│   ├── computer-science.md
│   ├── business.md
│   └── engineering.md
├── about/
│   ├── history.md
│   └── campus.md
└── admissions/
    ├── requirements.md
    └── process.md
```

**Workflow:**
1. Edit Markdown files
2. Commit to Git
3. GitHub Actions builds Astro site
4. Auto-deploy to Cloudflare Pages
5. Global CDN distribution

### Application Data (Dynamic)
- Stored in PostgreSQL
- Managed via Express.js API
- Accessed through BFF layer
- Real-time updates

---

## Session Management Options

### Option A: Backend JWT (Current Design - Recommended)

```mermaid
sequenceDiagram
    Browser->>BFF: Cookie: token
    BFF->>Backend: Authorization: Bearer JWT
    Note over Backend: jwt.verify(token)<br/>Decode: {userId: 123}<br/>No DB lookup needed
    Backend->>DB: Query with userId
```

**Pros:**
- ✅ Simple implementation
- ✅ Stateless backend (easy scaling)
- ✅ Fast (no session lookup)
- ✅ Good for <5,000 users

**Cons:**
- ⚠️ Can't revoke sessions immediately
- ⚠️ Token refresh requires new flow

### Option B: BFF Session + Backend JWT (Advanced)

```mermaid
sequenceDiagram
    Browser->>BFF: Cookie: sessionId
    BFF->>Redis: GET session:sessionId
    Redis->>BFF: {userId, role}
    BFF->>Backend: Generate fresh JWT
    BFF->>Backend: Authorization: Bearer JWT
```

**Pros:**
- ✅ Can revoke sessions immediately
- ✅ Track active sessions
- ✅ Better security

**Cons:**
- ⚠️ More complex
- ⚠️ Requires Redis/session storage
- ⚠️ Adds latency

**Recommendation:** Start with Option A, migrate to Option B if needed.

---

## Repository Strategy

### Monorepo Approach (Recommended)

This project uses a **single repository** (monorepo) containing both frontend and backend code.

**Why monorepo?**
- ✅ Single source of truth for all code
- ✅ Easier to keep frontend/backend in sync
- ✅ Atomic commits for cross-stack features
- ✅ Can share TypeScript types, constants, validation schemas
- ✅ Better code review (see full context of changes)
- ✅ Simpler for small teams

**Independent deployments:**
- GitHub Actions can deploy frontend and backend separately
- Changes to `frontend/` only trigger frontend deployment
- Changes to `backend/` only trigger backend deployment

---

## Project Structure

```
campus-website/                       # Single monorepo
├── .github/
│   └── workflows/
│       ├── deploy-frontend.yml      # Deploy Astro + Workers to Cloudflare
│       └── deploy-backend.yml       # Deploy Express.js to VPS
│
├── frontend/                        # Astro + Cloudflare Workers
│   ├── src/
│   │   ├── content/                 # Markdown content
│   │   │   ├── programs/
│   │   │   │   ├── computer-science.md
│   │   │   │   ├── business.md
│   │   │   │   └── engineering.md
│   │   │   ├── about/
│   │   │   │   ├── history.md
│   │   │   │   └── campus.md
│   │   │   └── admissions/
│   │   │       ├── requirements.md
│   │   │       └── process.md
│   │   ├── pages/                   # Astro pages
│   │   │   ├── index.astro         # Homepage
│   │   │   ├── login.astro
│   │   │   ├── register.astro
│   │   │   ├── dashboard.astro
│   │   │   └── apply.astro
│   │   ├── components/              # Reusable components
│   │   │   ├── Header.astro
│   │   │   ├── Footer.astro
│   │   │   └── ApplicationForm.astro
│   │   └── layouts/                 # Page layouts
│   │       ├── BaseLayout.astro
│   │       └── DashboardLayout.astro
│   ├── functions/                   # Cloudflare Workers (BFF)
│   │   ├── auth/
│   │   │   ├── google/
│   │   │   │   ├── login.js        # Initiate Google OIDC
│   │   │   │   └── callback.js     # Handle Google callback
│   │   │   ├── login.js            # Email/password login
│   │   │   ├── register.js         # User registration
│   │   │   └── logout.js           # Logout
│   │   └── applications/
│   │       ├── submit.js           # Submit application
│   │       ├── status.js           # Get application status
│   │       └── list.js             # List user's applications
│   ├── public/                      # Static assets
│   │   ├── favicon.ico
│   │   └── images/
│   ├── wrangler.toml                # Cloudflare Workers config
│   ├── astro.config.mjs             # Astro config
│   └── package.json
│
├── backend/                         # Express.js API
│   ├── src/
│   │   ├── routes/
│   │   │   ├── auth.js             # Authentication endpoints
│   │   │   ├── applications.js     # Application CRUD
│   │   │   └── users.js            # User management
│   │   ├── middleware/
│   │   │   ├── authenticateToken.js # JWT verification
│   │   │   ├── rateLimiter.js      # Rate limiting
│   │   │   └── errorHandler.js     # Error handling
│   │   ├── models/
│   │   │   ├── User.js             # User model
│   │   │   └── Application.js      # Application model
│   │   ├── config/
│   │   │   └── database.js         # Database connection
│   │   └── index.js                # Express app entry point
│   ├── migrations/                  # Database migrations
│   │   ├── 001_create_users.sql
│   │   └── 002_create_applications.sql
│   ├── tests/                       # Backend tests
│   │   ├── auth.test.js
│   │   └── applications.test.js
│   ├── .env.example                 # Environment variables template
│   └── package.json
│
├── shared/                          # Shared code (optional but recommended)
│   ├── types/                       # TypeScript types
│   │   ├── User.ts
│   │   ├── Application.ts
│   │   └── index.ts
│   ├── constants/
│   │   ├── applicationStatus.ts    # 'pending' | 'approved' | 'rejected'
│   │   └── userRoles.ts            # 'registrant' | 'staff'
│   └── validators/
│       └── applicationSchema.ts    # Shared validation logic
│
├── docs/                            # Documentation
│   ├── API.md                       # API documentation
│   ├── DEPLOYMENT.md                # Deployment guide
│   └── DEVELOPMENT.md               # Development setup
│
├── package.json                     # Root package.json (npm workspaces)
├── .gitignore
└── README.md
```

### Using npm Workspaces

**Root package.json:**
```json
{
  "name": "campus-website",
  "private": true,
  "workspaces": [
    "frontend",
    "backend",
    "shared"
  ],
  "scripts": {
    "dev": "npm run dev --workspaces",
    "build": "npm run build --workspaces",
    "deploy:frontend": "npm run deploy -w frontend",
    "deploy:backend": "npm run deploy -w backend",
    "test": "npm run test --workspaces"
  }
}
```

**Benefits:**
- Single `npm install` at root installs all packages
- Shared dependencies are deduplicated
- Easy to run scripts across workspaces
- `shared/` package can be imported by both frontend and backend

---

## Development Workflow

### Initial Setup

```bash
# Clone repository
git clone https://github.com/yourorg/campus-website.git
cd campus-website

# Install all dependencies (uses npm workspaces)
npm install

# This installs dependencies for:
# - frontend/
# - backend/
# - shared/
```

### Local Development

```bash
# Option 1: Run everything
npm run dev

# Option 2: Run individually

# Frontend (Astro + BFF)
cd frontend
npm install
npm run dev                      # http://localhost:4321

# Backend (Express.js)
cd backend
npm install
cp .env.example .env             # Configure database, JWT secret
npm run dev                      # http://localhost:3000
```

### Environment Variables

**Frontend (.dev.vars for Cloudflare Workers):**
```bash
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
BACKEND_URL=http://localhost:3000
APP_URL=http://localhost:4321
```

**Backend (.env):**
```bash
DATABASE_URL=postgresql://user:pass@localhost:5432/campus
JWT_SECRET=your-super-secret-key-change-this
PORT=3000
NODE_ENV=development
```

### Deployment

#### Automated Deployment via GitHub Actions

The monorepo uses path-based triggers to deploy only what changed.

**Frontend Deployment (`.github/workflows/deploy-frontend.yml`):**
```yaml
name: Deploy Frontend

on:
  push:
    branches: [main]
    paths:
      - 'frontend/**'
      - 'shared/**'
      - '.github/workflows/deploy-frontend.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci

      - name: Build Astro site
        working-directory: ./frontend
        run: npm run build

      - name: Deploy to Cloudflare Pages
        working-directory: ./frontend
        run: npx wrangler pages deploy dist --project-name=campus-website
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}

      - name: Deploy Cloudflare Workers (BFF)
        working-directory: ./frontend
        run: npx wrangler deploy
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
```

**Backend Deployment (`.github/workflows/deploy-backend.yml`):**
```yaml
name: Deploy Backend

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
      - 'shared/**'
      - '.github/workflows/deploy-backend.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Deploy to VPS
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USERNAME }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            cd /var/www/campus-website
            git pull origin main
            cd backend
            npm install --production
            npm run migrate
            pm2 restart campus-backend
```

**Key Features:**
- ✅ Frontend deploys only when `frontend/` or `shared/` changes
- ✅ Backend deploys only when `backend/` or `shared/` changes
- ✅ Both can deploy independently and simultaneously
- ✅ Shared types trigger both deployments (ensures consistency)

#### Manual Deployment

**Frontend to Cloudflare:**
```bash
cd frontend
npx wrangler pages deploy dist --project-name=campus-website
npx wrangler deploy  # Deploy Workers (BFF)
```

**Backend to VPS:**
```bash
# SSH to VPS
ssh user@your-vps.com

# Navigate to project
cd /var/www/campus-website

# Pull latest code
git pull origin main

# Install dependencies
cd backend
npm install --production

# Run migrations
npm run migrate

# Restart service
pm2 restart campus-backend
```

### Sharing Code Between Frontend and Backend

The `shared/` directory contains code used by both frontend and backend.

**Example: Shared TypeScript Types**

**shared/types/Application.ts:**
```typescript
export interface Application {
  id: number;
  userId: number;
  program: string;
  status: 'pending' | 'approved' | 'rejected';
  submittedAt: Date;
  documents?: {
    transcript: string;
    idCard: string;
    photo: string;
  };
}

export type ApplicationStatus = Application['status'];
```

**Backend usage (backend/src/routes/applications.js):**
```javascript
// Import shared types
import type { Application } from '../../../shared/types/Application';

router.get('/applications', async (req, res) => {
  const applications: Application[] = await db.query(
    'SELECT * FROM applications WHERE user_id = $1',
    [req.user.userId]
  );

  res.json(applications);
});
```

**Frontend usage (frontend/functions/applications/list.js):**
```javascript
// Import shared types
import type { Application } from '../../../shared/types/Application';

export async function onRequestGet(context) {
  const response = await fetch(`${env.BACKEND_URL}/applications`, {
    headers: { 'Authorization': `Bearer ${token}` }
  });

  const applications: Application[] = await response.json();
  return new Response(JSON.stringify(applications));
}
```

**Benefits:**
- ✅ Type safety across full stack
- ✅ Single source of truth for data structures
- ✅ Change types once, reflected everywhere
- ✅ Catch breaking changes at compile time

---

## Google Cloud Setup (OIDC)

### Steps to Configure Google OAuth 2.0

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project (free, no credit card required)
3. Enable "Google+ API" or "Google Identity"
4. Navigate to **Credentials** → **Create Credentials** → **OAuth 2.0 Client ID**
5. Configure:
   - **Application type:** Web application
   - **Authorized JavaScript origins:**
     - `https://youruni.edu`
     - `http://localhost:4321` (development)
   - **Authorized redirect URIs:**
     - `https://youruni.edu/api/auth/google/callback`
     - `http://localhost:4321/api/auth/google/callback`
6. Copy **Client ID** and **Client Secret**
7. Add to environment variables

**Cost:** Free forever

---

## Security Considerations

### Implemented Security Measures

- ✅ **HttpOnly Cookies** - XSS-resistant token storage
- ✅ **HTTPS Enforced** - All connections encrypted
- ✅ **Rate Limiting** - Prevent brute force attacks
- ✅ **Password Hashing** - bcrypt with salt rounds
- ✅ **JWT Expiration** - 7-day token lifetime
- ✅ **Email Domain Validation** - Staff role enforcement
- ✅ **CORS Configuration** - Restricted origins
- ✅ **Input Validation** - express-validator on all inputs
- ✅ **DDoS Protection** - Cloudflare edge filtering
- ✅ **Bot Detection** - Automatic bot blocking
- ✅ **SQL Injection Prevention** - Parameterized queries

---

## Scalability

### Current Capacity (5,000 users)

| Component | Capacity | Bottleneck |
|-----------|----------|------------|
| **Cloudflare Pages** | Unlimited | None |
| **Cloudflare Workers** | 100k req/day | ✅ Sufficient |
| **VPS (2GB RAM)** | ~5k users | CPU/Memory |
| **PostgreSQL** | 50k+ rows | None |

### Scaling Strategy

```mermaid
graph TB
    subgraph "Phase 1: <5K users"
        CF1[Cloudflare Free]
        VPS1[VPS 2GB RAM<br/>$5-10/mo]
    end

    subgraph "Phase 2: 5K-20K users"
        CF2[Cloudflare Free]
        VPS2[VPS 4GB RAM<br/>$15-20/mo]
    end

    subgraph "Phase 3: 20K+ users"
        CF3[Cloudflare Workers Paid<br/>$5/10M requests]
        LB[Load Balancer]
        VPS3A[VPS 4GB]
        VPS3B[VPS 4GB]
        DBManaged[Managed PostgreSQL]
    end

    Phase1 --> Phase2
    Phase2 --> Phase3
```

**Scaling Options:**
1. **Vertical Scaling:** Upgrade VPS (2GB → 4GB → 8GB)
2. **Horizontal Scaling:** Add load balancer + multiple VPS instances
3. **Database Scaling:** Move to managed PostgreSQL with replication
4. **Workers Scaling:** Upgrade to Cloudflare Workers Paid ($5/10M requests)

---

## Key Benefits

### Cost-Effective
- **$5-10/month total** for <5,000 users
- No hidden fees
- Predictable costs

### DDoS-Proof
- Hard limits prevent cost overruns
- **Maximum attack cost: $0**
- Automatic traffic filtering

### Secure
- Industry-standard authentication (OIDC)
- HttpOnly cookies prevent XSS
- Rate limiting prevents abuse

### Flexible
- Supports both SSO and traditional auth
- Multiple user roles (staff/registrant)
- Easy to extend

### Scalable
- Static site scales infinitely via CDN
- Backend can scale vertically or horizontally
- Database can migrate to managed service

### Developer-Friendly
- Modern tech stack (Astro, Express.js)
- Git-based workflow
- Easy local development

### SEO-Optimized
- Static site generation for marketing pages
- Fast page loads
- Excellent Core Web Vitals

### Low Maintenance
- Cloudflare handles edge layer
- Minimal operational overhead
- Automated deployments

---

## Future Enhancements

### Potential Features

- [ ] Email verification for local accounts
- [ ] Password reset flow (forgot password)
- [ ] Multi-factor authentication (2FA)
- [ ] Advanced RBAC (reviewers, admins, super-admin)
- [ ] Application workflow (stages, reviews, approvals)
- [ ] Email notifications (application status updates)
- [ ] Admin dashboard with analytics
- [ ] Document preview and validation
- [ ] Application deadline management
- [ ] Bulk application processing
- [ ] Export applications to CSV/Excel
- [ ] Mobile app (React Native using same backend)
- [ ] Real-time notifications (WebSocket)
- [ ] Audit logs for compliance
- [ ] Integration with student information system

### Migration Paths

**If you need better session control:**
- Migrate from Option A (stateless JWT) to Option B (Redis sessions)
- Add session tracking dashboard
- Implement "logout from all devices"

**If you outgrow free tier:**
- Upgrade Cloudflare Workers: $5/month for 10M requests
- Bigger VPS: $20-40/month for 8GB RAM
- Managed PostgreSQL: $15-50/month depending on provider

---

## Support & Resources

### Documentation
- [Astro Documentation](https://docs.astro.build)
- [Cloudflare Pages](https://developers.cloudflare.com/pages)
- [Cloudflare Workers](https://developers.cloudflare.com/workers)
- [Express.js Guide](https://expressjs.com)
- [PostgreSQL Docs](https://www.postgresql.org/docs)
- [Google OAuth 2.0](https://developers.google.com/identity/protocols/oauth2)

### Monitoring & Debugging
- Cloudflare Analytics (built-in, free)
- Cloudflare Logs (Workers logs)
- VPS logging (pm2 logs, nginx logs)
- Database monitoring (pg_stat_activity)

---

## License

This project architecture is designed for STMIK Campus Website.

---

## Contributors

**Architecture Design:** 2025
**Version:** 1.0
**Last Updated:** 2025-11-19
