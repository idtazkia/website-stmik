# Campus Website - Deployment Guide

## Table of Contents
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [CockroachDB Setup](#cockroachdb-setup)
- [Cloudflare R2 Setup](#cloudflare-r2-setup)
- [Local Development](#local-development)
- [Frontend Deployment (Cloudflare)](#frontend-deployment-cloudflare)
- [Automated Deployment (GitHub Actions)](#automated-deployment-github-actions)
- [Environment Variables](#environment-variables)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Accounts
- [ ] GitHub account (for repository and CI/CD)
- [ ] Cloudflare account (for Pages, Workers, and R2)
- [ ] CockroachDB account (for Serverless database)
- [ ] Google Cloud account (for OAuth OIDC)

### Required Software
- Node.js 20+ (`node --version`)
- npm 9+ (`npm --version`)
- Git (`git --version`)
- Wrangler CLI (`npm install -g wrangler`)

### Cost Summary
| Service | Free Tier | Cost |
|---------|-----------|------|
| Cloudflare Pages | Unlimited bandwidth | $0 |
| Cloudflare Workers | 100k req/day | $0 |
| Cloudflare R2 | 10GB storage | $0 |
| CockroachDB Serverless | 50M RUs/month, 10GB | $0 |
| Google OAuth | Unlimited | $0 |
| **Total** | | **$0/month**

---

## Initial Setup

### 1. Clone Repository

```bash
git clone https://github.com/yourorg/campus-website.git
cd campus-website
```

### 2. Install Dependencies (npm Workspaces)

```bash
# Install all dependencies at once
npm install

# This installs dependencies for:
# - frontend/
# - backend/
# - shared/
```

### 3. Google Cloud Setup (OIDC)

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new project (free, no credit card required)
3. Enable "Google+ API" or "Google Identity"
4. Navigate to **Credentials** → **Create Credentials** → **OAuth 2.0 Client ID**
5. Configure:
   - **Application type:** Web application
   - **Authorized JavaScript origins:**
     - `https://youruni.edu` (production)
     - `http://localhost:4321` (development)
   - **Authorized redirect URIs:**
     - `https://youruni.edu/api/auth/google/callback`
     - `http://localhost:4321/api/auth/google/callback`
6. Copy **Client ID** and **Client Secret**
7. Save for environment variables

**Cost:** Free forever

---

## CockroachDB Setup

CockroachDB Serverless provides a PostgreSQL-compatible database with a generous free tier (50M Request Units/month, 10GB storage).

### 1. Create CockroachDB Account

1. Go to [CockroachDB Cloud](https://cockroachlabs.cloud/)
2. Sign up with GitHub, Google, or email
3. Verify your email address

### 2. Create a Serverless Cluster

1. Click **Create Cluster**
2. Select **Serverless** (free tier)
3. Choose configuration:
   - **Cloud provider:** AWS or GCP (either works)
   - **Region:** Choose closest to your users (e.g., `asia-southeast1` for Indonesia)
   - **Spend limit:** Set to $0 (free tier only)
4. Name your cluster (e.g., `campus-website`)
5. Click **Create cluster**

### 3. Create Database and User

1. Once cluster is created, click **SQL Users** in left sidebar
2. Click **Add User**
3. Create a user:
   - **Username:** `campus_app`
   - **Password:** Generate a strong password (save it securely)
4. Click **Databases** in left sidebar
5. Click **Create Database**
6. Name it `campus`

### 4. Get Connection String

1. Click **Connect** button (top right)
2. Select **Connection string**
3. Choose your language: **Node.js**
4. Copy the connection string, it looks like:
   ```
   postgresql://campus_app:PASSWORD@cluster-name-1234.abc.cockroachlabs.cloud:26257/campus?sslmode=verify-full
   ```
5. Save this for environment variables

### 5. Download CA Certificate (Optional)

For local development, you may need the CA certificate:

1. In the Connect dialog, click **Download CA Cert**
2. Save to `~/.postgresql/root.crt`
3. Or use `sslmode=require` instead of `verify-full` for simplicity

### 6. Test Connection

```bash
# Using psql (PostgreSQL client)
psql "postgresql://campus_app:PASSWORD@cluster-name.cockroachlabs.cloud:26257/campus?sslmode=verify-full"

# Or using CockroachDB CLI
cockroach sql --url "postgresql://campus_app:PASSWORD@cluster-name.cockroachlabs.cloud:26257/campus?sslmode=verify-full"
```

### 7. Create Tables

Run the schema creation SQL:

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    provider VARCHAR(50) NOT NULL DEFAULT 'local',
    provider_id VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'registrant',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Applications table
CREATE TABLE applications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    program VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    submitted_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table (optional)
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_provider ON users(provider, provider_id);
CREATE INDEX idx_applications_user_id ON applications(user_id);
CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
```

### CockroachDB Free Tier Limits

| Resource | Limit | Notes |
|----------|-------|-------|
| Storage | 10 GB | Sufficient for 300K+ leads |
| Request Units | 50M/month | ~1% usage for 3,000 leads |
| Clusters | 5 | One per project |
| Burst capacity | Available | Handles traffic spikes |

### CockroachDB vs PostgreSQL Compatibility

For this application, CockroachDB works identically to PostgreSQL:

| Feature | Status |
|---------|--------|
| SERIAL/BIGSERIAL | ✅ Works |
| JOINs, subqueries | ✅ Works |
| Indexes (B-tree, GIN) | ✅ Works |
| JSON/JSONB | ✅ Works |
| Transactions | ✅ Works (SERIALIZABLE default) |
| Foreign keys | ✅ Works |

**Minor differences:**
- Default port is 26257 (not 5432)
- Uses `sslmode=verify-full` by default
- Some PostgreSQL-specific functions may need alternatives

---

## Cloudflare R2 Setup

Cloudflare R2 provides S3-compatible object storage for file uploads (documents, images).

### 1. Enable R2 in Cloudflare

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Select your account
3. Click **R2** in the left sidebar
4. Click **Create bucket**

### 2. Create a Bucket

1. **Bucket name:** `campus-uploads`
2. **Location:** Automatic (or choose specific region)
3. Click **Create bucket**

### 3. Create API Token for R2

1. Go to **R2** → **Manage R2 API Tokens**
2. Click **Create API Token**
3. Configure:
   - **Token name:** `campus-website-r2`
   - **Permissions:** Object Read & Write
   - **Bucket:** Select `campus-uploads`
4. Click **Create API Token**
5. Copy and save:
   - **Access Key ID**
   - **Secret Access Key**
   - **Endpoint URL** (looks like `https://ACCOUNT_ID.r2.cloudflarestorage.com`)

### 4. Configure CORS (for browser uploads)

1. Go to **R2** → **campus-uploads** bucket
2. Click **Settings**
3. Under **CORS Policy**, add:

```json
[
  {
    "AllowedOrigins": ["https://yourdomain.com", "http://localhost:4321"],
    "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
    "AllowedHeaders": ["*"],
    "MaxAgeSeconds": 3600
  }
]
```

### R2 Free Tier Limits

| Resource | Limit |
|----------|-------|
| Storage | 10 GB |
| Class A operations (writes) | 1M/month |
| Class B operations (reads) | 10M/month |
| Egress | Free (no bandwidth charges) |

---

## Local Development

### Frontend Development (Static Site)

```bash
cd frontend
npm install
npm run dev                      # http://localhost:4321
```

### Full Stack Development (with Workers)

**Terminal 1 - Frontend with Workers emulation:**
```bash
cd frontend
npm install
cp .dev.vars.example .dev.vars  # Configure environment variables
npm run dev                      # http://localhost:4321
```

**.dev.vars configuration:**
```bash
# CockroachDB connection (use your actual connection string)
DATABASE_URL=postgresql://campus_app:PASSWORD@cluster.cockroachlabs.cloud:26257/campus?sslmode=verify-full

# Authentication
JWT_SECRET=generate-a-long-random-secret-min-32-chars
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here

# Cloudflare R2
R2_ACCESS_KEY_ID=your-r2-access-key
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_ENDPOINT=https://ACCOUNT_ID.r2.cloudflarestorage.com
R2_BUCKET=campus-uploads
```

### Database Access During Development

Since CockroachDB is cloud-hosted, no local database setup is needed:

```bash
# Connect to CockroachDB for debugging
psql "postgresql://campus_app:PASSWORD@cluster.cockroachlabs.cloud:26257/campus?sslmode=verify-full"

# Or use CockroachDB CLI
cockroach sql --url "YOUR_CONNECTION_STRING"
```

### Local Testing Workflow

1. **Test static pages:** Visit `http://localhost:4321`
2. **Test authentication:** Try login/register flows
3. **Test API:** Cloudflare Workers run locally via Wrangler
4. **Test database:** Queries go directly to CockroachDB Serverless

---

## Frontend Deployment (Cloudflare)

### Prerequisites

1. **Install Wrangler CLI:**
```bash
npm install -g wrangler
```

2. **Login to Cloudflare:**
```bash
wrangler login
```

3. **Get API Token:**
   - Go to Cloudflare Dashboard
   - My Profile → API Tokens
   - Create Token → "Edit Cloudflare Workers" template
   - Save token securely

### Deploy Astro Site to Cloudflare Pages

**Method 1: Via Wrangler (Recommended)**

```bash
cd frontend

# Build Astro site
npm run build

# Deploy to Cloudflare Pages
npx wrangler pages deploy dist --project-name=campus-website

# Deploy Cloudflare Workers (BFF)
npx wrangler deploy
```

**Method 2: Via Cloudflare Dashboard**

1. Go to Cloudflare Dashboard → Pages
2. Connect to Git → Select GitHub repository
3. Configure build settings:
   - **Build command:** `cd frontend && npm install && npm run build`
   - **Build output directory:** `frontend/dist`
   - **Root directory:** `/`
4. Add environment variables (see Environment Variables section)
5. Deploy

### Configure Custom Domain

1. Add domain to Cloudflare
2. Update DNS records (nameservers)
3. In Cloudflare Pages → Custom Domains
4. Add `youruni.edu` and `www.youruni.edu`
5. SSL/TLS will be configured automatically

### Deploy Cloudflare Workers (BFF)

**wrangler.toml configuration:**

```toml
# frontend/wrangler.toml
name = "campus-website-bff"
main = "functions/**/*.js"
compatibility_date = "2024-01-01"

[env.production]
vars = { BACKEND_URL = "https://api.youruni.edu" }

[env.development]
vars = { BACKEND_URL = "http://localhost:3000" }
```

**Deploy:**
```bash
cd frontend
npx wrangler deploy
```

---

## Automated Deployment (GitHub Actions)

### 1. Set Up GitHub Secrets

Go to GitHub Repository → Settings → Secrets and variables → Actions

Add the following secrets:

**For Cloudflare Deployment:**
- `CLOUDFLARE_API_TOKEN` - Your Cloudflare API token
- `CLOUDFLARE_ACCOUNT_ID` - Your Cloudflare account ID

**For Database (CockroachDB):**
- `DATABASE_URL` - CockroachDB connection string

**For Authentication:**
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `JWT_SECRET` - JWT signing secret

**For File Storage (R2):**
- `R2_ACCESS_KEY_ID` - Cloudflare R2 access key
- `R2_SECRET_ACCESS_KEY` - Cloudflare R2 secret key

### 2. Create Workflow Files

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

### 3. Test Deployment

**Trigger deployment:**
```bash
# Make a change to frontend
echo "# Test" >> frontend/README.md
git add frontend/README.md
git commit -m "Test deployment"
git push origin main

# Check GitHub Actions tab for deployment status
```

**Verify deployment:**
1. Check GitHub Actions tab for green checkmark
2. Visit your Cloudflare Pages URL
3. Test authentication and database connectivity

---

## Environment Variables

### Cloudflare Workers Environment Variables

Configure in Cloudflare Dashboard → Workers → Settings → Variables, or in `wrangler.toml`:

```bash
# Database (CockroachDB Serverless)
DATABASE_URL=postgresql://campus_app:PASSWORD@cluster.cockroachlabs.cloud:26257/campus?sslmode=verify-full

# Authentication
JWT_SECRET=generate-a-long-random-secret-min-32-chars
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here

# Application
APP_URL=https://youruni.edu
NODE_ENV=production

# Cloudflare R2 (File Storage)
R2_ACCESS_KEY_ID=your-r2-access-key
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_ENDPOINT=https://ACCOUNT_ID.r2.cloudflarestorage.com
R2_BUCKET=campus-uploads

# Optional: Email (for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@youruni.edu
SMTP_PASS=your-app-password
```

### Local Development (.dev.vars)

Create `.dev.vars` file in frontend directory:

```bash
DATABASE_URL=postgresql://campus_app:PASSWORD@cluster.cockroachlabs.cloud:26257/campus?sslmode=verify-full
JWT_SECRET=local-dev-secret-at-least-32-characters
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here
R2_ACCESS_KEY_ID=your-r2-access-key
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_ENDPOINT=https://ACCOUNT_ID.r2.cloudflarestorage.com
R2_BUCKET=campus-uploads
```

### How to Generate JWT Secret

```bash
# Option 1: OpenSSL
openssl rand -base64 48

# Option 2: Node.js
node -e "console.log(require('crypto').randomBytes(48).toString('base64'))"
```

### Setting Variables in Cloudflare Dashboard

1. Go to Cloudflare Dashboard → Workers & Pages
2. Select your Worker
3. Go to **Settings** → **Variables**
4. Add each variable as **Encrypted** (for secrets) or **Plain text**
5. Click **Save and Deploy**

---

## Monitoring & Maintenance

### Cloudflare Monitoring

**Cloudflare Dashboard:**
- **Pages → Analytics**: Page views, unique visitors, bandwidth
- **Workers → Metrics**: Request count, errors, CPU time, response time
- **R2 → Metrics**: Storage usage, operations count

**Set up alerts:**
1. Cloudflare Dashboard → Notifications
2. Create alerts for:
   - Error rate spikes (Workers)
   - High CPU usage (Workers)
   - Storage threshold (R2)

### CockroachDB Monitoring

**CockroachDB Cloud Console:**
1. Go to [CockroachDB Cloud](https://cockroachlabs.cloud)
2. Select your cluster
3. View metrics:
   - **SQL Activity**: Queries/second, latency
   - **Request Units**: RU consumption vs. limit
   - **Storage**: Data size, index size
   - **Sessions**: Active connections

**Check usage via SQL:**
```sql
-- Connect to your cluster
cockroach sql --url "YOUR_CONNECTION_STRING"

-- Check table sizes
SELECT table_name,
       pg_size_pretty(pg_total_relation_size(table_name::regclass)) as size
FROM information_schema.tables
WHERE table_schema = 'public';

-- Check row counts
SELECT 'users' as table_name, count(*) as rows FROM users
UNION ALL
SELECT 'applications', count(*) FROM applications
UNION ALL
SELECT 'sessions', count(*) FROM sessions;
```

**Set up alerts in CockroachDB:**
1. CockroachDB Console → Alerts
2. Configure alerts for:
   - RU consumption > 80% of limit
   - Storage > 8GB (80% of free tier)
   - Connection failures

### Backup Strategy

**CockroachDB Serverless backups:**
- CockroachDB automatically handles backups
- Point-in-time recovery available
- No manual backup configuration needed

**Manual export (if needed):**
```bash
# Export data using cockroach dump
cockroach dump campus --url "YOUR_CONNECTION_STRING" > backup.sql

# Or export specific tables
cockroach dump campus users applications --url "YOUR_CONNECTION_STRING" > backup.sql
```

**R2 file backups:**
- Cloudflare R2 provides 99.999999999% durability
- No additional backup configuration needed
- Consider cross-region replication for critical files

---

## Troubleshooting

### Frontend/Pages Issues

**Build fails on Cloudflare:**
```bash
# Check build logs in Cloudflare Dashboard → Pages → Deployments
# Common issues:
# - Missing environment variables
# - Node version mismatch
# - Dependency installation failure

# Test build locally
cd frontend
npm run build
```

**Workers not deploying:**
```bash
# Check wrangler.toml configuration
# Verify API token has correct permissions
# Try manual deployment:
cd frontend
npx wrangler deploy --verbose

# Check Workers logs in Cloudflare Dashboard
```

### CockroachDB Issues

**Connection errors:**
```bash
# Test connection
psql "YOUR_CONNECTION_STRING"

# Common issues:
# - Wrong password
# - SSL mode not specified
# - IP not allowlisted (check CockroachDB Console → Networking)

# Check if cluster is running
# CockroachDB Console → Overview → Status should be "Running"
```

**Query timeout or slow queries:**
```sql
-- Check for missing indexes
SHOW INDEXES FROM users;
SHOW INDEXES FROM applications;

-- Add indexes if needed
CREATE INDEX idx_applications_user_id ON applications(user_id);
CREATE INDEX idx_applications_status ON applications(status);

-- Check query execution plan
EXPLAIN ANALYZE SELECT * FROM applications WHERE user_id = 1;
```

**RU consumption too high:**
```sql
-- Check which queries consume most RUs
-- Go to CockroachDB Console → SQL Activity → Statements

-- Optimize queries:
-- 1. Add indexes for frequently filtered columns
-- 2. Use LIMIT for large result sets
-- 3. Avoid SELECT * - select only needed columns
```

### Cloudflare Workers Issues

**Workers returning errors:**
```bash
# Check Workers logs:
# Cloudflare Dashboard → Workers → Your Worker → Logs

# Common issues:
# - Missing environment variables
# - Database connection timeout
# - Invalid JWT secret

# Test locally with wrangler
cd frontend
npx wrangler dev
```

**Authentication not working:**
```bash
# Check:
# 1. GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET are set
# 2. Redirect URIs match exactly in Google Console
# 3. JWT_SECRET is configured

# Test Google OAuth flow manually
# Visit: https://yoursite.com/api/auth/google/login
```

### R2 Storage Issues

**Upload failures:**
```bash
# Check:
# 1. R2_ACCESS_KEY_ID and R2_SECRET_ACCESS_KEY are set
# 2. CORS policy allows your domain
# 3. Bucket name is correct

# Test R2 connection using AWS CLI (S3-compatible)
aws s3 ls s3://campus-uploads --endpoint-url https://ACCOUNT_ID.r2.cloudflarestorage.com
```

### GitHub Actions Issues

**Deployment failing:**
```bash
# Check GitHub Actions logs
# Common issues:
# - Missing secrets (CLOUDFLARE_API_TOKEN, etc.)
# - API token permissions incorrect

# Verify secrets are set:
# Repository → Settings → Secrets → Actions
```

**Path triggers not working:**
```yaml
# Ensure paths are correct in workflow file
paths:
  - 'frontend/**'  # Note: no leading slash
  - 'shared/**'
```

### Performance Issues

**Slow page loads:**
```bash
# Check Cloudflare Analytics for:
# - Cache hit ratio (should be high for static assets)
# - Response times by region

# Optimize:
# 1. Enable caching headers in Astro
# 2. Use Cloudflare CDN (automatic with Pages)
# 3. Optimize images
```

**Database latency:**
```bash
# Check CockroachDB Console → SQL Activity → Latency

# If latency is high:
# 1. Check if region is close to users
# 2. Add indexes for slow queries
# 3. Consider connection pooling
```

---

## Rollback Procedures

### Rollback Frontend (Cloudflare Pages)

**Via Cloudflare Dashboard (Recommended):**
1. Go to Cloudflare Dashboard → Pages
2. Select your project
3. Click **Deployments**
4. Find the previous working deployment
5. Click **...** → **Rollback to this deployment**

**Via Wrangler CLI:**
```bash
cd frontend
git checkout <previous-commit>
npm run build
npx wrangler pages deploy dist --project-name=campus-website
```

### Rollback Workers

**Via Cloudflare Dashboard:**
1. Go to Cloudflare Dashboard → Workers
2. Select your Worker
3. Click **Deployments**
4. Select previous version → **Rollback**

**Via Wrangler CLI:**
```bash
cd frontend
git checkout <previous-commit>
npx wrangler deploy
```

### Rollback Database Changes

**CockroachDB does not support automatic rollback migrations.**

**Option 1: Manual SQL rollback**
```sql
-- If you have rollback SQL prepared
-- Connect to your cluster
cockroach sql --url "YOUR_CONNECTION_STRING"

-- Run rollback commands manually
DROP INDEX IF EXISTS idx_new_index;
ALTER TABLE users DROP COLUMN IF EXISTS new_column;
```

**Option 2: Restore from export**
```bash
# If you exported data before migration
cockroach sql --url "YOUR_CONNECTION_STRING" < backup.sql
```

**Option 3: Point-in-time recovery (Enterprise feature)**
- Contact CockroachDB support for point-in-time recovery
- Available for Serverless clusters with support plan

---

## Security Checklist

- [ ] All environment variables use strong, unique secrets
- [ ] JWT_SECRET is at least 32 characters, randomly generated
- [ ] CockroachDB password is strong (20+ chars, mixed case, numbers, symbols)
- [ ] All Cloudflare secrets are set as "Encrypted" type
- [ ] GitHub secrets configured (never commit secrets to repo)
- [ ] Google OAuth redirect URIs match exactly (production URLs only)
- [ ] CORS configured to allow only your frontend domain
- [ ] R2 bucket CORS policy restricts allowed origins
- [ ] Rate limiting enabled via Cloudflare (automatic)
- [ ] SSL/HTTPS enforced (automatic with Cloudflare)
- [ ] CockroachDB network allowlist configured (if using IP restrictions)
- [ ] No secrets in wrangler.toml (use environment variables instead)

---

## Support & Resources

### Documentation Links
- [Astro Documentation](https://docs.astro.build)
- [Cloudflare Pages](https://developers.cloudflare.com/pages)
- [Cloudflare Workers](https://developers.cloudflare.com/workers)
- [Cloudflare R2](https://developers.cloudflare.com/r2)
- [CockroachDB Serverless](https://www.cockroachlabs.com/docs/cockroachcloud/serverless)
- [CockroachDB SQL Reference](https://www.cockroachlabs.com/docs/stable/sql-statements)
- [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler)

### Useful Commands

```bash
# Frontend Development
cd frontend && npm run dev          # Dev server with hot reload
cd frontend && npm run build        # Build for production
cd frontend && npm run preview      # Preview production build

# Cloudflare Deployment
npx wrangler login                  # Login to Cloudflare
npx wrangler pages deploy dist      # Deploy Pages
npx wrangler deploy                 # Deploy Workers
npx wrangler tail                   # Stream Workers logs

# CockroachDB
cockroach sql --url "CONNECTION_STRING"    # Connect to database
psql "CONNECTION_STRING"                    # Alternative: use psql

# Local Development
cp .dev.vars.example .dev.vars      # Create local env file
npx wrangler dev                    # Run Workers locally
```

### Quick Reference

| Task | Command |
|------|---------|
| Start dev server | `cd frontend && npm run dev` |
| Build for production | `cd frontend && npm run build` |
| Deploy to Cloudflare | `npx wrangler pages deploy dist` |
| Deploy Workers | `npx wrangler deploy` |
| View Workers logs | `npx wrangler tail` |
| Connect to database | `cockroach sql --url "URL"` |

---

**Last Updated:** 2026-01-15
**Version:** 2.0 (Serverless Architecture)
