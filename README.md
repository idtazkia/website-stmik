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

## Project Structure

```
campus-website/
├── marketing/                    # Static site repo
│   ├── .github/
│   │   └── workflows/
│   │       └── deploy.yml       # GitHub Actions
│   ├── src/
│   │   ├── content/             # Markdown content
│   │   │   ├── programs/
│   │   │   ├── about/
│   │   │   └── admissions/
│   │   ├── pages/               # Astro pages
│   │   ├── components/          # Reusable components
│   │   └── layouts/             # Page layouts
│   ├── functions/               # Cloudflare Workers (BFF)
│   │   ├── auth/
│   │   │   ├── google/
│   │   │   │   ├── login.js
│   │   │   │   └── callback.js
│   │   │   ├── login.js         # Email/password
│   │   │   └── logout.js
│   │   └── applications/
│   │       ├── submit.js
│   │       └── status.js
│   ├── wrangler.toml            # Cloudflare config
│   ├── astro.config.mjs
│   └── package.json
│
└── backend/                     # Backend API repo
    ├── src/
    │   ├── routes/
    │   │   ├── auth.js
    │   │   ├── applications.js
    │   │   └── users.js
    │   ├── middleware/
    │   │   ├── authenticateToken.js
    │   │   └── rateLimiter.js
    │   ├── models/
    │   │   ├── User.js
    │   │   └── Application.js
    │   └── index.js             # Express app
    ├── migrations/              # Database migrations
    ├── .env.example
    └── package.json
```

---

## Development Workflow

### Local Development

```bash
# Frontend (Astro + BFF)
cd marketing
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

**Frontend to Cloudflare:**
```bash
# Deploy static site + workers
npx wrangler deploy
```

**Backend to VPS:**
```bash
# SSH to VPS
ssh user@your-vps.com

# Pull latest code
git pull origin main

# Install dependencies
npm install

# Run migrations
npm run migrate

# Restart service
pm2 restart backend
```

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
