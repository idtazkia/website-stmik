# Campus Website - Technical Architecture

## Table of Contents
- [Architecture Pattern](#architecture-pattern)
- [Technology Stack](#technology-stack)
- [Authentication System](#authentication-system)
- [Authentication Flows](#authentication-flows)
- [Database Schema](#database-schema)
- [Session Management](#session-management)
- [Traffic Analysis](#traffic-analysis)
- [DDoS Protection & Cost Safety](#ddos-protection--cost-safety)
- [Security Considerations](#security-considerations)
- [Scalability](#scalability)

---

## Architecture Pattern

**Static Site + Go Backend Pattern**

```mermaid
graph TB
    subgraph "User Layer"
        Browser[Browser]
    end

    subgraph "Edge Layer - Cloudflare (Free)"
        CDN[Cloudflare CDN<br/>DDoS Protection]
        StaticSite[Static Marketing Site<br/>Astro + Markdown<br/>Cloudflare Pages]
    end

    subgraph "VPS Layer ($5/month)"
        Nginx[Nginx<br/>Reverse Proxy + SSL]
        GoAPI[Go Backend<br/>REST API + Auth]
        Database[(PostgreSQL<br/>Database)]
        FileStorage[Local Disk<br/>File Storage]
    end

    Browser -->|HTTPS| CDN
    CDN -->|Static| StaticSite
    CDN -->|API Calls| Nginx
    Nginx -->|Proxy| GoAPI
    GoAPI -->|SQL| Database
    GoAPI -->|Files| FileStorage

    StaticSite -.->|Links to| CDN
```

###  Key Architectural Decisions

**Why Go for Backend?**
1. **Performance**: Compiled binary, minimal resource usage (~30MB RAM)
2. **Simplicity**: Single binary deployment, no runtime dependencies
3. **Concurrency**: Goroutines handle thousands of concurrent requests
4. **Type Safety**: Compile-time error checking
5. **Low Latency**: No cold starts, always running

**Why Static Site Generation?**
1. **Performance**: Pre-rendered HTML served from global CDN
2. **SEO**: Excellent crawlability for marketing pages
3. **Cost**: Unlimited bandwidth on Cloudflare Pages
4. **Developer Experience**: Markdown-based content management
5. **Git-based Workflow**: Version control for content

**Why VPS over Serverless?**
1. **Simplicity**: Single server vs multiple cloud services
2. **Local Database**: <1ms latency vs ~30ms remote
3. **No Cold Starts**: Always running, instant response
4. **Full Control**: Direct access, easy debugging
5. **Predictable Cost**: Fixed $5/month, no usage surprises

---

## Technology Stack

| Component | Technology | Hosting | Cost |
|-----------|-----------|---------|------|
| **Static Site** | Astro + Markdown | Cloudflare Pages | Free |
| **CDN/DDoS** | Cloudflare | Cloudflare | Free |
| **Backend API** | Go (Golang) | VPS | $5/mo |
| **Database** | PostgreSQL 18 | VPS | Included |
| **File Storage** | Local Disk | VPS | Included |
| **Reverse Proxy** | Nginx + Let's Encrypt | VPS | Included |
| **Build/Deploy** | GitHub Actions | GitHub | Free |

**Total Monthly Cost: $5** (fixed, predictable)

### Technology Rationale

**Frontend: Astro**
- Static site generation for marketing pages
- Server-side rendering capability for dynamic pages
- Excellent performance (minimal JavaScript)
- MDX support for rich content
- Component islands architecture

**Backend: Go**
- Compiled binary (~15MB), no runtime dependencies
- Minimal memory usage (~30-50MB for this workload)
- Excellent concurrency with goroutines
- Strong standard library (net/http, database/sql)
- Fast compilation and deployment
- Single binary deployment via scp or rsync

**Database: PostgreSQL**
- Industry standard, 30+ years of reliability
- ACID compliance for critical data
- JSON/JSONB support for flexible schemas
- Excellent performance with proper indexing
- Local to application (<1ms latency)
- Easy backup with pg_dump

**VPS Benefits**
- Single server to manage (simpler than 3+ cloud services)
- Full control over configuration
- Easy debugging and log access
- No vendor lock-in
- Predictable monthly cost

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
- HttpOnly cookies (XSS protection)
- JWT tokens managed by Go backend
- Tokens never exposed to client JavaScript
- Rate limiting via Nginx
- DDoS protection via Cloudflare CDN
- bcrypt password hashing (10 salt rounds)
- JWT expiration (7 days)
- SQL injection prevention (parameterized queries with pgx)

---

## Authentication Flows

### Google OIDC Authentication Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Nginx
    participant Go as Go Backend
    participant Google as Google OAuth
    participant DB as PostgreSQL

    Note over Browser,DB: Login Initiation
    Browser->>Nginx: GET /auth/google/login
    Nginx->>Go: Proxy request
    Go->>Browser: HTTP 302 Redirect to Google
    Browser->>Google: User logs in with Google
    Google->>Browser: Redirect with authorization code

    Note over Browser,DB: Token Exchange
    Browser->>Go: GET /auth/google/callback?code=AUTHCODE
    Go->>Google: POST /token (code, client_id, client_secret)
    Google->>Go: {id_token, access_token}

    Note over Go: Verify Google ID Token<br/>Extract: {email, name, sub}

    Go->>DB: SELECT/INSERT/UPDATE user
    DB->>Go: User data {id: UUID, email, role}

    Note over Go: Generate JWT<br/>{userId: UUID, role, email}
    Note over Go: Set HttpOnly Cookie
    Go->>Browser: HTTP 302 /dashboard<br/>Set-Cookie: token=JWT

    Note over Browser: Cookie stored<br/>User logged in!
```

### Traditional Email/Password Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Go as Go Backend
    participant DB as PostgreSQL

    Browser->>Go: POST /auth/login<br/>{email, password}

    Go->>DB: SELECT * FROM users<br/>WHERE email = ? AND provider = 'local'
    DB->>Go: User data + password_hash

    Note over Go: bcrypt.CompareHashAndPassword()

    Note over Go: Generate JWT<br/>{userId: UUID, email, role}
    Note over Go: Set HttpOnly Cookie
    Go->>Browser: HTTP 302 /dashboard<br/>Set-Cookie: token=JWT

    Note over Browser: Cookie stored<br/>User logged in!
```

### Authenticated Request Flow

```mermaid
sequenceDiagram
    participant Browser
    participant Go as Go Backend
    participant DB as PostgreSQL

    Note over Browser: Cookie automatically sent
    Browser->>Go: GET /portal/applications<br/>Cookie: token=JWT

    Note over Go: Extract JWT from cookie<br/>Verify with jwt.Parse()

    Go->>DB: SELECT * FROM applications<br/>WHERE user_id = UUID
    DB->>Go: Application records
    Go->>Browser: Render Templ template
```

---

## Session Identification: Two Token Pattern

### Token 1: Google's ID Token (OIDC)
**Used only during login, then discarded**

| Property | Value |
|----------|-------|
| **Issued by** | Google |
| **Signed with** | Google's private key |
| **Verified by** | Go backend (using Google's public cert) |
| **Contains** | `sub` (Google user ID), `email`, `name`, `iss`, `aud` |
| **Purpose** | Prove user identity during OAuth login |
| **Lifetime** | 1 hour (but only used once) |
| **Stored** | NOWHERE (discarded after verification) |

### Token 2: Your Backend's JWT
**Used for all authenticated requests**

| Property | Value |
|----------|-------|
| **Issued by** | Go backend |
| **Signed with** | `JWT_SECRET` |
| **Verified by** | Go backend |
| **Contains** | `userId` (UUID), `email`, `role`, `provider` |
| **Purpose** | Session management for all requests |
| **Lifetime** | 7 days |
| **Stored** | HttpOnly cookie (set by Go backend) |

### Flow Summary

```mermaid
graph TB
    subgraph "Login Phase"
        A[Google ID Token] -->|Verify & Extract| B[User Info]
        B -->|Create/Update User| C[Go Backend Generates JWT]
    end

    subgraph "Session Phase"
        C -->|Set HttpOnly Cookie| D[Browser]
        D -->|Cookie sent with requests| E[Go Backend verifies JWT]
    end

    style A fill:#fbbc04
    style C fill:#34a853
    style E fill:#34a853
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
        uuid id PK
        varchar email UK
        varchar name
        varchar password_hash "NULL for Google users"
        varchar provider "google or local"
        varchar provider_id "Google user ID"
        varchar role "registrant or staff"
        boolean is_active
        timestamp last_assigned_at "For round-robin"
        timestamp created_at
        timestamp updated_at
    }

    APPLICATIONS {
        uuid id PK
        uuid user_id FK
        uuid prospect_id FK
        uuid intake_id FK
        uuid program_id FK
        uuid track_id FK
        varchar status "pending_review/approved/enrolled/cancelled"
        varchar va_number
        timestamp paid_at
        timestamp created_at
        timestamp updated_at
    }

```

### Table Descriptions

**USERS Table**
- Stores both local and Google OAuth users
- `provider` field determines authentication method
- `password_hash` is NULL for Google users
- `provider_id` stores Google's unique user ID
- `role` determines access level (registrant/staff)
- `last_assigned_at` for round-robin staff assignment

**APPLICATIONS Table**
- Linked to prospects, programs, tracks, intakes
- `status` enum: pending_review, approved, enrolled, cancelled
- `va_number` stores virtual account for payment

See `backend/README.md` for complete database schema with all tables.

### PostgreSQL Configuration

**Connection string format (local on VPS):**
```
postgresql://campus_app:password@localhost:5432/campus?sslmode=disable
```

**Recommended PostgreSQL settings (postgresql.conf):**
```ini
# Memory (for 1GB VPS)
shared_buffers = 256MB
effective_cache_size = 512MB
work_mem = 4MB

# Connections
max_connections = 100

# Logging
log_statement = 'none'  # Set to 'all' for debugging
log_min_duration_statement = 1000  # Log queries > 1 second
```

**Go database driver:**
```go
import "github.com/jackc/pgx/v5"

// Connection with pool
pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
```

---

## Session Management

### JWT in HttpOnly Cookie (Recommended)

```mermaid
sequenceDiagram
    Browser->>GoAPI: Cookie: token (HttpOnly)
    Note over GoAPI: jwt.Verify(token)<br/>Decode: {userId: 123, role}
    GoAPI->>DB: Query with userId
    DB->>GoAPI: User data
    GoAPI->>Browser: Response + Set-Cookie (if refresh needed)
```

**Implementation in Go:**
```go
// Set cookie after login
http.SetCookie(w, &http.Cookie{
    Name:     "token",
    Value:    jwtToken,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode,
    MaxAge:   7 * 24 * 60 * 60, // 7 days
    Path:     "/",
})

// Read cookie in middleware
cookie, err := r.Cookie("token")
claims, err := jwt.Verify(cookie.Value, jwtSecret)
```

**Pros:**
- ✅ Simple implementation
- ✅ Stateless (no session storage needed)
- ✅ Fast (no database lookup for auth)
- ✅ XSS-resistant (HttpOnly cookie)
- ✅ Good for <50,000 users

**Cons:**
- ⚠️ Can't revoke individual sessions (use short expiry + refresh)
- ⚠️ Token in cookie adds to request size (~500 bytes)

---

## VPS Capacity Analysis

### Real-World Usage Estimation

All API requests go through the Go backend on VPS.

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
- Create account: 1 request
- Login: 1 request
Total: 2 requests

Days 2-5: Application Preparation
- Daily login: 1 request
- View form: 1 request
- Auto-save drafts: 5 requests
Total: 7 requests/day × 4 days = 28 requests

Day 6: Application Submission
- Login: 1 request
- Load form: 1 request
- Upload documents: 3 requests
- Submit: 1 request
Total: 6 requests

Days 7-30: Status Checking
- Login: 1 request
- Check status: 1 request
- View updates: 1 request
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

| Scenario | Daily Requests | VPS Capacity | Status |
|----------|---------------|--------------|--------|
| **Average Day** | 2,155 | 1M+ req/day | ✅ 0.2% used |
| **Peak Day** (100 simultaneous) | 2,655 | 1M+ req/day | ✅ 0.3% used |
| **10x Traffic Spike** | 21,550 | 1M+ req/day | ✅ 2% used |
| **Quiet Period** | 200 | 1M+ req/day | ✅ Negligible |

### VPS Resource Usage (1GB RAM, 1 vCPU)

| Lead Count | RAM Usage | CPU Usage | Storage | Status |
|------------|-----------|-----------|---------|--------|
| **3,000 leads** | ~300MB | <5% | ~50MB | ✅ Comfortable |
| **10,000 leads** | ~350MB | <10% | ~150MB | ✅ Comfortable |
| **30,000 leads** | ~400MB | <20% | ~500MB | ✅ Good |
| **100,000 leads** | ~500MB | <30% | ~1.5GB | ⚠️ Consider 2GB VPS |

### Database Capacity (PostgreSQL on VPS)

| Lead Count | Data Size | Index Size | Total | 25GB Disk |
|------------|-----------|------------|-------|-----------|
| **3,000 leads** | ~5MB | ~2MB | ~10MB | ✅ 0.04% |
| **10,000 leads** | ~15MB | ~5MB | ~25MB | ✅ 0.1% |
| **100,000 leads** | ~150MB | ~50MB | ~250MB | ✅ 1% |
| **1,000,000 leads** | ~1.5GB | ~500MB | ~2.5GB | ✅ 10% |

### Important: Static vs Dynamic Traffic

```mermaid
graph LR
    subgraph "Website Traffic Distribution"
        A[1000 Daily Visitors]
        A -->|95%| B[950 visitors<br/>Browse Marketing Pages<br/>Static Content]
        A -->|5%| C[50 visitors<br/>Login & Apply<br/>VPS API]
    end

    B -->|Static| D[Cloudflare Pages CDN<br/>UNLIMITED bandwidth<br/>$0 cost]
    C -->|API| E[VPS Go Backend<br/>1M+ req/day capacity<br/>$5/month]
```

**Key Insight:**
- **Marketing pages** (programs, about, contact): Unlimited traffic via Cloudflare CDN
- **Application portal** (login, submit, status): Goes through VPS (~5% of traffic)
- **VPS handles 10,000x more traffic than needed** for 3,000 leads

### Verdict for 3,000 Leads + 5 Admins

**VPS (1GB RAM, 1 vCPU, $5/month):**
✅ Using <5% of CPU capacity
✅ Using ~30% of RAM (comfortable headroom)
✅ Using <1% of storage
✅ Can handle 100,000+ leads without upgrade
✅ Fixed $5/month cost regardless of traffic

**Database (PostgreSQL on VPS):**
✅ Local database = <1ms latency
✅ No external dependencies
✅ 25GB disk can store 1M+ leads
✅ Full PostgreSQL features available

**Conclusion: A $5/month VPS provides massive headroom for 3,000 leads and can scale to 100,000+ without upgrade.**

---

## DDoS Protection & Cost Safety

### Cloudflare Protection (VPS behind Cloudflare Proxy)

| Feature | Protection | Cost |
|---------|------------|------|
| **Static Site** | Cloudflare CDN | $0 |
| **API Endpoint** | Cloudflare Proxy (orange cloud) | $0 |
| **DDoS Mitigation** | Automatic L3/L4/L7 filtering | $0 |
| **Bot Detection** | Automatic blocking | $0 |
| **Rate Limiting** | 10,000 free rules/month | $0 |

### VPS Protection Strategy

```mermaid
graph LR
    A[Attacker] -->|DDoS| B[Cloudflare Edge]
    B -->|Filtered Traffic| C[VPS]
    B -->|Blocked| D[Dropped]
```

**Configuration:**
1. Point domain DNS to Cloudflare (proxy enabled - orange cloud)
2. VPS firewall allows only Cloudflare IPs
3. Nginx rate limiting as secondary defense
4. fail2ban for SSH protection

### Cost Safety Summary
- ✅ **Maximum DDoS cost: $5/month** (fixed VPS cost)
- ✅ VPS IP hidden behind Cloudflare proxy
- ✅ No surprise bills (unlike serverless with usage-based pricing)
- ✅ Cloudflare absorbs attack traffic
- ✅ VPS only sees legitimate requests

---

## Security Considerations

### Implemented Security Measures

- ✅ **HttpOnly Cookies** - XSS-resistant token storage
- ✅ **HTTPS Enforced** - Let's Encrypt via Nginx
- ✅ **Rate Limiting** - Nginx + Cloudflare
- ✅ **Password Hashing** - bcrypt with salt rounds
- ✅ **JWT Expiration** - 7-day token lifetime
- ✅ **Email Domain Validation** - Staff role enforcement
- ✅ **CORS Configuration** - Restricted origins
- ✅ **Input Validation** - Go struct validation
- ✅ **DDoS Protection** - Cloudflare proxy
- ✅ **Bot Detection** - Cloudflare automatic blocking
- ✅ **SQL Injection Prevention** - Parameterized queries (pgx)
- ✅ **Firewall** - UFW allowing only SSH, HTTP, HTTPS

### Security Middleware (Go)

```go
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // HSTS
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        // Content Security Policy
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        // Prevent clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        // Prevent MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        // XSS Protection
        w.Header().Set("X-XSS-Protection", "1; mode=block")

        next.ServeHTTP(w, r)
    })
}
```

### Rate Limiting Strategy

**Nginx (primary):**
```nginx
# /etc/nginx/conf.d/rate-limit.conf
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=auth:10m rate=1r/s;

server {
    location /api/ {
        limit_req zone=api burst=20 nodelay;
        proxy_pass http://localhost:8080;
    }

    location /api/auth/ {
        limit_req zone=auth burst=5 nodelay;
        proxy_pass http://localhost:8080;
    }
}
```

**Go (secondary, for specific endpoints):**
```go
import "golang.org/x/time/rate"

var loginLimiter = rate.NewLimiter(rate.Every(time.Second), 5)

func loginHandler(w http.ResponseWriter, r *http.Request) {
    if !loginLimiter.Allow() {
        http.Error(w, "Too many login attempts", http.StatusTooManyRequests)
        return
    }
    // handle login
}
```

---

## Scalability

### Current Capacity (3,000 leads, 300 registrations)

| Component | Capacity | Current Usage | Bottleneck |
|-----------|----------|---------------|------------|
| **Cloudflare Pages** | Unlimited | N/A | None |
| **VPS CPU** | 1 vCPU | <5% | None |
| **VPS RAM** | 1 GB | ~30% | None |
| **PostgreSQL** | 100 connections | ~5 | None |
| **Disk** | 25 GB | <1% | None |

### Scaling Strategy

```mermaid
graph TB
    subgraph "Phase 1: <100K leads - $5/mo"
        VPS1[1GB VPS<br/>Go + PostgreSQL]
    end

    subgraph "Phase 2: 100K-500K leads - $10-20/mo"
        VPS2[2-4GB VPS<br/>Go + PostgreSQL]
    end

    subgraph "Phase 3: 500K+ leads - $50+/mo"
        LB[Load Balancer]
        VPS3A[VPS: Go API]
        VPS3B[VPS: Go API]
        DB3[Managed PostgreSQL<br/>or Dedicated VPS]
    end

    Phase1 --> Phase2
    Phase2 --> Phase3
```

**Scaling Options:**
1. **Phase 1 (current):** 1GB VPS - handles up to 100K leads ($5/mo)
2. **Phase 2:** Upgrade to 2-4GB VPS ($10-20/mo)
3. **Phase 3:** Multiple VPS behind load balancer + managed PostgreSQL
4. **Alternative:** Migrate to Kubernetes if needed

### Performance Optimization Strategies

**Caching:**
- Static assets: Aggressive caching via Cloudflare CDN
- API responses: In-memory cache in Go (sync.Map or groupcache)
- Database queries: Connection pooling via pgx

**Database Optimization:**
- Indexes on frequently queried fields (user_id, email, status)
- Connection pooling (pgxpool)
- Query optimization with EXPLAIN ANALYZE
- Regular VACUUM (automatic in PostgreSQL 13+)

**Frontend Optimization:**
- Static site generation
- Image optimization
- Code splitting
- Service Worker for offline support

---

## Key Benefits

### Cost-Effective
- **$5/month** fixed cost (VPS)
- No usage-based surprises
- Predictable, simple billing
- Can scale to 100K+ leads without cost increase

### DDoS-Proof
- Cloudflare proxy absorbs attacks
- VPS IP hidden from public
- **Maximum attack cost: $5/month** (fixed)

### Simple Architecture
- Single server to manage
- No vendor lock-in
- Easy debugging and maintenance
- Full control over stack

### Performant
- Go backend: minimal latency, high throughput
- Local PostgreSQL: <1ms database queries
- No cold starts: always running
- Handles 1M+ requests/day

### Secure
- Industry-standard authentication (OIDC)
- HttpOnly cookies prevent XSS
- Rate limiting at Nginx + app level
- Firewall + fail2ban

### Scalable
- Static site scales infinitely via CDN
- VPS can upgrade vertically ($5 → $10 → $20)
- Can add load balancer when needed
- PostgreSQL handles millions of rows

### Developer-Friendly
- Modern tech stack (Astro, Go)
- Git-based workflow
- Simple deployment (scp binary + systemctl restart)
- Easy local development

### SEO-Optimized
- Static site generation for marketing pages
- Fast page loads
- Excellent Core Web Vitals

---

## Contributors

**Architecture Design:** 2025
**Version:** 2.0 (VPS + Go + PostgreSQL)
**Last Updated:** 2026-01-15
