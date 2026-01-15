# Campus Website - Technical Architecture

## Table of Contents
- [Architecture Pattern](#architecture-pattern)
- [Technology Stack](#technology-stack)
- [Authentication System](#authentication-system)
- [Authentication Flows](#authentication-flows)
- [Database Schema](#database-schema)
- [Session Management](#session-management)
- [BFF Traffic Analysis](#bff-traffic-analysis)
- [DDoS Protection & Cost Safety](#ddos-protection--cost-safety)
- [Security Considerations](#security-considerations)
- [Scalability](#scalability)

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
        FileStorage[Cloudflare R2<br/>File Storage]
    end

    subgraph "Database Layer - Cockroach Labs (Free)"
        Database[(CockroachDB Serverless<br/>50M RUs/month)]
    end

    Browser -->|HTTPS| StaticSite
    Browser -->|API Calls<br/>Cookie: token| BFF
    BFF -->|SQL queries| Database
    BFF -->|Upload/Download| FileStorage

    StaticSite -.->|Links to| BFF
```

###  Key Architectural Decisions

**Why BFF Pattern?**
1. **Security**: HttpOnly cookies prevent XSS attacks
2. **Token Management**: BFF handles token extraction and forwarding
3. **API Aggregation**: Can combine multiple backend calls in future
4. **Rate Limiting**: Edge-level protection via Cloudflare
5. **CORS Simplification**: Same-origin requests from browser

**Why Static Site Generation?**
1. **Performance**: Pre-rendered HTML served from global CDN
2. **SEO**: Excellent crawlability for marketing pages
3. **Cost**: Unlimited bandwidth on Cloudflare Pages
4. **Developer Experience**: Markdown-based content management
5. **Git-based Workflow**: Version control for content

**Why Serverless for BFF?**
1. **Zero Cost**: Cloudflare Workers free tier (100k req/day)
2. **Global Distribution**: Runs at edge, low latency
3. **Auto-scaling**: No infrastructure management
4. **DDoS Protection**: Built-in hard limits

---

## Technology Stack

| Component | Technology | Hosting | Cost |
|-----------|-----------|---------|------|
| **Static Site** | Astro + Markdown | Cloudflare Pages | Free |
| **BFF Layer** | Cloudflare Workers | Cloudflare | Free (100k req/day) |
| **Backend API** | Express.js (Node.js) | Cloudflare Workers / VPS | Free / $5-10/mo |
| **Database** | CockroachDB Serverless | Cockroach Labs | Free (50M RUs/mo) |
| **File Storage** | Cloudflare R2 | Cloudflare | Free tier (10GB) |
| **Build/Deploy** | GitHub Actions | GitHub | Free |

**Total Monthly Cost: $0-5** (potentially fully free with serverless backend)

### Technology Rationale

**Frontend: Astro**
- Static site generation for marketing pages
- Server-side rendering capability for dynamic pages
- Excellent performance (minimal JavaScript)
- MDX support for rich content
- Component islands architecture

**BFF: Cloudflare Workers**
- Edge computing for low latency
- JavaScript/TypeScript runtime
- Built-in KV storage for caching
- Generous free tier
- Global distribution

**Backend: Express.js**
- Mature, battle-tested framework
- Large ecosystem of middleware
- Easy to understand and maintain
- Your 20+ years of experience applicable
- Good PostgreSQL integration

**Database: CockroachDB Serverless**
- PostgreSQL wire-compatible (standard CRUD, JOINs, indexes work identically)
- ACID compliance for critical data
- JSON/JSONB support for flexible schemas
- Generous free tier: 10GB storage, 50M Request Units/month
- No cold start or auto-pause (unlike Neon/Supabase free tiers)
- Supports 3,000 leads at ~1% of free tier capacity
- Automatic scaling and replication
- No database maintenance required

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
  - Leverages Google's 2FA
- **Option 2: Email/Password** (traditional)
  - Privacy-conscious option
  - Self-managed credentials
  - Requires email verification
  - Password reset flow needed

#### Marketing Staff:
- **Google OIDC only**
- Enforced via `@youruni.edu` email domain check
- Automatically assigned `staff` role
- Leverages institutional Google Workspace

### Security Features
- ✅ HttpOnly cookies (XSS protection)
- ✅ JWT tokens managed by BFF layer
- ✅ Tokens never exposed to client JavaScript
- ✅ Same-origin requests (no CORS issues)
- ✅ Rate limiting built-in (Cloudflare)
- ✅ DDoS protection included
- ✅ bcrypt password hashing (10 salt rounds)
- ✅ JWT expiration (7 days)
- ✅ SQL injection prevention (parameterized queries)

---

## Authentication Flows

### Google OIDC Authentication Flow

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

### Traditional Email/Password Flow

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

### Authenticated Request Flow

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

### Table Descriptions

**USERS Table**
- Stores both local and Google OAuth users
- `provider` field determines authentication method
- `password_hash` is NULL for Google users
- `provider_id` stores Google's unique user ID
- `role` determines access level (registrant/staff)

**APPLICATIONS Table**
- One-to-many relationship with USERS
- `status` enum: pending, approved, rejected
- Additional fields can include: program, documents, etc.

**SESSIONS Table** (Optional - for Option B session management)
- Tracks active sessions
- Enables "logout from all devices" feature
- Can store additional metadata

### CockroachDB Compatibility Notes

CockroachDB is PostgreSQL wire-compatible. For this application, all standard operations work identically:

| Feature | Status | Notes |
|---------|--------|-------|
| CRUD operations | ✅ Works | SELECT, INSERT, UPDATE, DELETE |
| JOINs | ✅ Works | All join types supported |
| Indexes | ✅ Works | B-tree, GIN, inverted indexes |
| Transactions | ✅ Works | ACID compliant, SERIALIZABLE by default |
| SERIAL/BIGSERIAL | ✅ Works | Auto-incrementing IDs |
| JSON/JSONB | ✅ Works | Full JSON support |
| Timestamps | ✅ Works | TIMESTAMP, TIMESTAMPTZ |
| VARCHAR, TEXT | ✅ Works | String types |
| BOOLEAN | ✅ Works | True/false values |

**Minor differences (not relevant for this app):**
- Some PostgreSQL-specific functions may need alternatives
- `pg_*` system tables are limited
- Sequences behave slightly differently (use SERIAL instead)

**Connection string format:**
```
postgresql://username:password@host:26257/database?sslmode=verify-full
```

---

## Session Management

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

## BFF Traffic Analysis

### Real-World Usage Estimation

All authenticated API requests go through the BFF (Cloudflare Workers) layer.

### Traffic Breakdown: 3,000 Leads + 5 Admins

**Funnel Assumptions:**
- Target registrations: 300
- Conversion rate: 10%
- Total leads: 3,000

#### Lead Journey (30-day admission period)

**Non-Converting Leads (2,700 users - 90%):**
```
- Account creation: 2 requests
- Browse/form views: 5-8 requests
- Abandon without completing
Total per non-converting lead: ~10 requests
```

#### Registrant Journey (300 users - 10% who convert)

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

**3,000 Leads (300 converting + 2,700 non-converting):**
- Non-converting leads: 2,700 × 10 = **27,000 requests**
- Converting registrants: 300 × 108 = **32,400 requests**
- Total over 30 days: **59,400 requests**
- Daily average: **1,980 requests/day**

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
| **Average Day** | 2,155 | **2.2%** | ✅ Very Safe |
| **Peak Day** (100 simultaneous submissions) | 2,655 | **2.7%** | ✅ Very Safe |
| **10x Traffic Spike** | 21,550 | **21.6%** | ✅ Safe |
| **Quiet Period** (post-admission) | 200 | **0.2%** | ✅ Very Safe |

### Scalability Analysis (BFF Layer - Cloudflare Workers)

| Lead Count | Daily BFF Requests | % of Free Tier | Monthly Cost | Status |
|------------|-------------------|----------------|--------------|--------|
| **3,000 leads** | 2,155 | 2.2% | $0 | ✅ Very Safe |
| **10,000 leads** | 7,200 | 7.2% | $0 | ✅ Very Safe |
| **30,000 leads** | 21,500 | 21.5% | $0 | ✅ Safe |
| **100,000 leads** | 72,000 | 72% | $0 | ✅ Safe |
| **140,000 leads** | 99,000 | 99% | $0 | ⚠️ Near Limit |
| **150,000+ leads** | 110,000+ | 110%+ | $5/mo* | ✅ Upgrade Available |

*Cloudflare Workers Paid: $5/month for 10 million requests

### Database Capacity (CockroachDB Serverless)

| Lead Count | Monthly RUs | % of Free Tier | Monthly Cost | Status |
|------------|-------------|----------------|--------------|--------|
| **3,000 leads** | ~500K | 1% | $0 | ✅ Very Safe |
| **10,000 leads** | ~1.7M | 3.4% | $0 | ✅ Very Safe |
| **30,000 leads** | ~5M | 10% | $0 | ✅ Safe |
| **100,000 leads** | ~17M | 34% | $0 | ✅ Safe |
| **150,000 leads** | ~25M | 50% | $0 | ✅ Safe |
| **300,000 leads** | ~50M | 100% | $0 | ⚠️ Near Limit |

*CockroachDB Serverless: Free tier includes 50M RUs/month, 10GB storage

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

### Verdict for 3,000 Leads + 5 Admins

**BFF Layer (Cloudflare Workers):**
✅ Using only 2.2% of the free tier
✅ 97.8% buffer for traffic spikes
✅ Even 10x growth keeps you under 22% usage
✅ Zero cost risk - hard limits prevent overages
✅ Could scale to 100,000 leads and still be at 72%

**Database (CockroachDB Serverless):**
✅ Using only 1% of the free tier (50M RUs)
✅ No cold start latency (unlike Neon/Supabase free tiers)
✅ No auto-pause policy
✅ Could scale to 150,000 leads and still be at 50%

**Conclusion: Both BFF and database have massive headroom. The system can scale to 100,000+ leads on fully free infrastructure.**

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

### Comparison with Other Platforms

| Platform | DDoS Risk | Hard Limits | Max Cost on Attack |
|----------|-----------|-------------|-------------------|
| **Cloudflare Pages + Workers** | ✅ VERY LOW | ✅ Yes | **$0** |
| **Vercel (Free)** | ⚠️ HIGH | ❌ No | **$500+** |
| **Netlify (Free)** | ✅ LOW | ✅ Yes | $0 (stops) |

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

### Security Headers (Express.js)

```javascript
const helmet = require('helmet');

app.use(helmet({
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'"],
      scriptSrc: ["'self'"],
      imgSrc: ["'self'", "data:", "https:"],
    }
  },
  hsts: {
    maxAge: 31536000,
    includeSubDomains: true,
    preload: true
  }
}));
```

### Rate Limiting Strategy

**Cloudflare Workers (BFF):**
```javascript
// Per-IP rate limiting
const rateLimiter = {
  limit: 10,           // 10 requests
  window: 60,          // per 60 seconds
};
```

**Express.js Backend:**
```javascript
const rateLimit = require('express-rate-limit');

const authLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 5,                    // 5 attempts
  message: 'Too many login attempts'
});

app.use('/auth/login', authLimiter);
```

---

## Scalability

### Current Capacity (3,000 leads, 300 registrations)

| Component | Capacity | Current Usage | Bottleneck |
|-----------|----------|---------------|------------|
| **Cloudflare Pages** | Unlimited | N/A | None |
| **Cloudflare Workers** | 100k req/day | 2.2% | None |
| **CockroachDB Serverless** | 50M RUs/month | 1% | None |
| **CockroachDB Storage** | 10 GB | <0.1% | None |

### Scaling Strategy

```mermaid
graph TB
    subgraph "Phase 1: <100K leads - FREE"
        CF1[Cloudflare Free<br/>Pages + Workers]
        DB1[CockroachDB Serverless Free<br/>50M RUs/month]
    end

    subgraph "Phase 2: 100K-300K leads"
        CF2[Cloudflare Workers Paid<br/>$5/10M requests]
        DB2[CockroachDB Serverless<br/>Pay-as-you-go]
    end

    subgraph "Phase 3: 300K+ leads"
        CF3[Cloudflare Workers Paid]
        DB3[CockroachDB Dedicated<br/>or Self-hosted PostgreSQL]
    end

    Phase1 --> Phase2
    Phase2 --> Phase3
```

**Scaling Options:**
1. **Phase 1 (current):** Fully free - CockroachDB Serverless + Cloudflare free tiers
2. **Phase 2:** Pay-as-you-go on both platforms (~$10-20/month for 100K-300K leads)
3. **Phase 3:** Dedicated clusters or migrate to self-hosted PostgreSQL on VPS
4. **Alternative:** Can always fall back to VPS + PostgreSQL ($5-10/month) if CockroachDB pricing changes

### Performance Optimization Strategies

**Caching:**
- Static assets: Aggressive caching via CDN
- API responses: Cache in Cloudflare Workers KV
- Database queries: Redis for frequently accessed data

**Database Optimization:**
- Indexes on frequently queried fields (user_id, email, status)
- Connection pooling (handled by CockroachDB Serverless)
- Query optimization
- No manual maintenance required (CockroachDB handles vacuuming, compaction)

**Frontend Optimization:**
- Static site generation
- Image optimization
- Code splitting
- Service Worker for offline support

---

## Key Benefits

### Cost-Effective
- **$0/month** for up to 100,000 leads (fully free tier)
- No hidden fees
- Predictable costs
- Fallback to VPS ($5-10/month) always available

### DDoS-Proof
- Hard limits prevent cost overruns
- **Maximum attack cost: $0**
- Automatic traffic filtering

### Secure
- Industry-standard authentication (OIDC)
- HttpOnly cookies prevent XSS
- Rate limiting prevents abuse

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

## Contributors

**Architecture Design:** 2025
**Version:** 1.0
**Last Updated:** 2025-11-19
