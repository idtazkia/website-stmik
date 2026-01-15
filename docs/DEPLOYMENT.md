# Campus Website - Deployment Guide

## Table of Contents
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [VPS Setup](#vps-setup)
- [PostgreSQL Setup](#postgresql-setup)
- [Go Backend Setup](#go-backend-setup)
- [Local Development](#local-development)
- [Frontend Deployment (Cloudflare)](#frontend-deployment-cloudflare)
- [Backend Deployment (VPS)](#backend-deployment-vps)
- [Automated Deployment (GitHub Actions)](#automated-deployment-github-actions)
- [Environment Variables](#environment-variables)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Accounts
- [ ] GitHub account (for repository and CI/CD)
- [ ] Cloudflare account (for Pages and CDN)
- [ ] VPS provider account (DigitalOcean, Hetzner, Vultr, or similar)
- [ ] Google Cloud account (for OAuth OIDC)

### Required Software
- Node.js 20+ (`node --version`)
- npm 9+ (`npm --version`)
- Git (`git --version`)
- Go 1.21+ (`go version`) - for backend development
- PostgreSQL client (`psql --version`) - for database management

### Cost Summary
| Service | Tier | Cost |
|---------|------|------|
| Cloudflare Pages | Free (unlimited bandwidth) | $0 |
| Cloudflare CDN | Free (DDoS protection) | $0 |
| VPS (1GB RAM) | Go + PostgreSQL + Nginx | $5/mo |
| Google OAuth | Free (unlimited) | $0 |
| Let's Encrypt | Free (SSL certificate) | $0 |
| **Total** | | **$5/month**

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

## VPS Setup

A 1GB RAM VPS can comfortably host the Go backend, PostgreSQL, and Nginx for 100,000+ leads.

### 1. Choose VPS Provider

Recommended providers for $5/month VPS (1GB RAM, 1 vCPU, 25GB SSD):
- **DigitalOcean** - $6/month (Singapore region)
- **Hetzner** - €4.51/month (Germany, best value)
- **Vultr** - $5/month (Singapore region)
- **Linode** - $5/month (Singapore region)

### 2. Provision VPS

1. Create account with chosen provider
2. Create a new VPS/Droplet:
   - **OS:** Ubuntu 24.04 LTS
   - **Size:** 1GB RAM, 1 vCPU
   - **Region:** Asia (Singapore) or closest to users
   - **SSH Key:** Add your public key
3. Note the public IP address

### 3. Initial Server Setup

```bash
# SSH into server
ssh root@YOUR_VPS_IP

# Update system
apt update && apt upgrade -y

# Create non-root user
adduser campus
usermod -aG sudo campus

# Copy SSH key to new user
mkdir -p /home/campus/.ssh
cp ~/.ssh/authorized_keys /home/campus/.ssh/
chown -R campus:campus /home/campus/.ssh

# Disable root SSH login
sed -i 's/PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
systemctl restart sshd

# Configure firewall
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

### 4. Install Required Software

```bash
# Install Nginx
apt install -y nginx

# Install PostgreSQL 18
apt install -y postgresql-18 postgresql-contrib-16

# Install Go (latest stable)
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
go version

# Install Certbot for SSL
apt install -y certbot python3-certbot-nginx
```

---

## PostgreSQL Setup

### 1. Configure PostgreSQL

```bash
# Switch to postgres user
sudo -u postgres psql

# Create database and user
CREATE USER campus_app WITH PASSWORD 'YOUR_STRONG_PASSWORD';
CREATE DATABASE campus OWNER campus_app;
GRANT ALL PRIVILEGES ON DATABASE campus TO campus_app;
\q
```

### 2. Configure Authentication

Edit `/etc/postgresql/18/main/pg_hba.conf`:

```bash
# Allow local connections with password
# Add this line:
local   campus   campus_app   scram-sha-256
```

Edit `/etc/postgresql/18/main/postgresql.conf`:

```bash
# Listen only on localhost (Go app is on same server)
listen_addresses = 'localhost'
```

Restart PostgreSQL:

```bash
sudo systemctl restart postgresql
```

### 3. Test Connection

```bash
psql -U campus_app -d campus -h localhost
# Enter password when prompted
```

### 4. Create Tables

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

-- Sessions table
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

---

## Go Backend Setup

### 1. Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Auth, logging, etc.
│   ├── models/              # Database models
│   └── repository/          # Database queries
├── go.mod
├── go.sum
└── .env
```

### 2. Initialize Go Module

```bash
cd backend
go mod init github.com/yourorg/campus-website/backend

# Install dependencies
go get github.com/jackc/pgx/v5          # PostgreSQL driver
go get github.com/go-chi/chi/v5         # Router
go get github.com/golang-jwt/jwt/v5     # JWT
go get golang.org/x/crypto/bcrypt       # Password hashing
```

### 3. Environment Variables

Create `backend/.env`:

```bash
# Database
DATABASE_URL=postgres://campus_app:PASSWORD@localhost:5432/campus?sslmode=disable

# Authentication
JWT_SECRET=generate-a-long-random-secret-min-32-chars
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret

# Application
APP_URL=https://yourdomain.com
PORT=3000
```

### 4. Build and Run

```bash
cd backend

# Development
go run ./cmd/server

# Production build
go build -o campus-api ./cmd/server
./campus-api
```

### 5. Create Systemd Service

Create `/etc/systemd/system/campus-api.service`:

```ini
[Unit]
Description=Campus API Server
After=network.target postgresql.service

[Service]
Type=simple
User=campus
WorkingDirectory=/home/campus/backend
ExecStart=/home/campus/backend/campus-api
Restart=always
RestartSec=5
Environment=GIN_MODE=release
EnvironmentFile=/home/campus/backend/.env

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable campus-api
sudo systemctl start campus-api
sudo systemctl status campus-api
```

### 6. Configure Nginx Reverse Proxy

Create `/etc/nginx/sites-available/campus`:

```nginx
# Rate limiting zones
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=auth:10m rate=1r/s;

server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        return 301 https://$server_name$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL configuration (managed by Certbot)
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";

    # API proxy
    location /api/ {
        limit_req zone=api burst=20 nodelay;

        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Auth endpoints (stricter rate limit)
    location /api/auth/ {
        limit_req zone=auth burst=5 nodelay;

        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable site and get SSL:

```bash
sudo ln -s /etc/nginx/sites-available/campus /etc/nginx/sites-enabled/
sudo nginx -t
sudo certbot --nginx -d api.yourdomain.com
sudo systemctl reload nginx
```

### VPS Resource Usage (3,000 leads)

| Component | RAM Usage | CPU Usage |
|-----------|-----------|-----------|
| Go Backend | ~50MB | <1% |
| PostgreSQL | ~200MB | <3% |
| Nginx | ~10MB | <1% |
| OS | ~100MB | - |
| **Total** | ~360MB | <5% |

**Capacity:** Can scale to 100,000+ leads on same 1GB VPS

---

## Local Development

### Frontend Development (Static Site Only)

```bash
cd frontend
npm install
npm run dev                      # http://localhost:4321
```

### Full Stack Development (with Go Backend)

**Terminal 1 - Go Backend:**
```bash
cd backend

# Start local PostgreSQL (if not running)
# macOS: brew services start postgresql@18
# Linux: sudo systemctl start postgresql

# Create local database (first time only)
createdb campus
psql campus < schema.sql

# Run backend
cp .env.example .env            # Configure environment
go run ./cmd/server             # http://localhost:3000
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev                      # http://localhost:4321
```

**backend/.env configuration:**
```bash
# Database (local PostgreSQL)
DATABASE_URL=postgres://campus_app:PASSWORD@localhost:5432/campus?sslmode=disable

# Authentication
JWT_SECRET=local-dev-secret-at-least-32-characters
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here

# Application
APP_URL=http://localhost:4321
PORT=3000
```

### Local PostgreSQL Setup (macOS)

```bash
# Install PostgreSQL
brew install postgresql@18
brew services start postgresql@18

# Create user and database
createuser -P campus_app        # Enter password when prompted
createdb -O campus_app campus

# Test connection
psql -U campus_app -d campus -h localhost
```

### Local PostgreSQL Setup (Linux)

```bash
# Install PostgreSQL
sudo apt install postgresql-18

# Create user and database
sudo -u postgres createuser -P campus_app
sudo -u postgres createdb -O campus_app campus

# Test connection
psql -U campus_app -d campus -h localhost
```

### Local Testing Workflow

1. **Test static pages:** Visit `http://localhost:4321`
2. **Test authentication:** Login flows via `http://localhost:3000/api/auth/`
3. **Test API:** `curl http://localhost:3000/api/health`
4. **Test database:** `psql -U campus_app -d campus -h localhost`

---

## Frontend Deployment (Cloudflare)

The static Astro site is deployed to Cloudflare Pages via Git integration (auto-deploy on push).

### Method 1: Cloudflare Dashboard (Recommended)

1. Go to [Cloudflare Dashboard](https://dash.cloudflare.com) → Pages
2. Click **Create a project** → **Connect to Git**
3. Select your GitHub repository
4. Configure build settings:
   - **Framework preset:** Astro
   - **Build command:** `cd frontend && npm install && npm run build`
   - **Build output directory:** `frontend/dist`
   - **Root directory:** `/`
5. Click **Save and Deploy**

After initial setup, every push to `main` triggers automatic deployment.

### Method 2: Manual CLI Deploy

```bash
cd frontend

# Build Astro site
npm run build

# Deploy to Cloudflare Pages (requires wrangler login)
npx wrangler pages deploy dist --project-name=campus-website
```

### Configure Custom Domain

1. Add your domain to Cloudflare (update nameservers at registrar)
2. In Cloudflare Dashboard → Pages → your project
3. Click **Custom domains** → **Set up a custom domain**
4. Add `yourdomain.com` and `www.yourdomain.com`
5. SSL/TLS is configured automatically

### Configure Cloudflare CDN for VPS

To route API traffic through Cloudflare CDN (hiding VPS IP):

1. In Cloudflare DNS, add A record:
   - **Name:** `api`
   - **IPv4 address:** Your VPS IP
   - **Proxy status:** Proxied (orange cloud)
2. API requests to `api.yourdomain.com` now go through Cloudflare
3. VPS only sees Cloudflare IPs (DDoS protection)

---

## Backend Deployment (VPS)

### Manual Deployment

```bash
# SSH into VPS
ssh campus@YOUR_VPS_IP

# Pull latest code
cd ~/backend
git pull origin main

# Build Go binary
go build -o campus-api ./cmd/server

# Restart service
sudo systemctl restart campus-api

# Verify
sudo systemctl status campus-api
curl http://localhost:3000/api/health
```

### Database Migrations

```bash
# SSH into VPS
ssh campus@YOUR_VPS_IP

# Run migrations
cd ~/backend
psql -U campus_app -d campus -h localhost < migrations/001_create_tables.sql

# Verify
psql -U campus_app -d campus -h localhost -c "\dt"
```

### Rollback

```bash
# SSH into VPS
ssh campus@YOUR_VPS_IP

# Checkout previous version
cd ~/backend
git checkout <previous-commit>

# Rebuild and restart
go build -o campus-api ./cmd/server
sudo systemctl restart campus-api
```

---

## Automated Deployment (GitHub Actions)

### 1. Set Up GitHub Secrets

Go to GitHub Repository → Settings → Secrets and variables → Actions

Add the following secrets:

**For Cloudflare Pages Deployment:**
- `CLOUDFLARE_API_TOKEN` - Your Cloudflare API token
- `CLOUDFLARE_ACCOUNT_ID` - Your Cloudflare account ID

**For VPS Deployment:**
- `VPS_HOST` - VPS IP address or hostname
- `VPS_USER` - SSH username (e.g., `campus`)
- `VPS_SSH_KEY` - Private SSH key for deployment

### 2. Create Workflow Files

**Frontend Deployment (`.github/workflows/deploy-frontend.yml`):**

```yaml
name: Deploy Frontend

on:
  push:
    branches: [main]
    paths:
      - 'frontend/**'
      - '.github/workflows/deploy-frontend.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
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
```

**Backend Deployment (`.github/workflows/deploy-backend.yml`):**

```yaml
name: Deploy Backend

on:
  push:
    branches: [main]
    paths:
      - 'backend/**'
      - '.github/workflows/deploy-backend.yml'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Build Go binary
        working-directory: ./backend
        run: |
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o campus-api ./cmd/server

      - name: Deploy to VPS
        uses: appleboy/scp-action@v0.1.7
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          source: "backend/campus-api"
          target: "~/"

      - name: Restart service
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            mv ~/backend/campus-api ~/campus-api.new
            sudo systemctl stop campus-api
            mv ~/campus-api ~/campus-api.old
            mv ~/campus-api.new ~/campus-api
            sudo systemctl start campus-api
            curl -f http://localhost:3000/api/health || (mv ~/campus-api.old ~/campus-api && sudo systemctl start campus-api && exit 1)
```

### 3. Test Deployment

**Trigger frontend deployment:**
```bash
# Make a change to frontend
echo "<!-- test -->" >> frontend/src/pages/index.astro
git add frontend/
git commit -m "Test frontend deployment"
git push origin main
```

**Trigger backend deployment:**
```bash
# Make a change to backend
echo "// test" >> backend/cmd/server/main.go
git add backend/
git commit -m "Test backend deployment"
git push origin main
```

**Verify deployment:**
1. Check GitHub Actions tab for green checkmark
2. Visit your Cloudflare Pages URL for frontend
3. Test API endpoint: `curl https://api.yourdomain.com/api/health`

---

## Environment Variables

### VPS Backend Environment Variables

Create `/home/campus/backend/.env` on the VPS:

```bash
# Database (local PostgreSQL)
DATABASE_URL=postgres://campus_app:PASSWORD@localhost:5432/campus?sslmode=disable

# Authentication
JWT_SECRET=generate-a-long-random-secret-min-32-chars
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here

# Application
APP_URL=https://yourdomain.com
PORT=3000

# Optional: Email (for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@yourdomain.com
SMTP_PASS=your-app-password
```

### Local Development Environment

Create `backend/.env` for local development:

```bash
# Database (local PostgreSQL)
DATABASE_URL=postgres://campus_app:PASSWORD@localhost:5432/campus?sslmode=disable

# Authentication
JWT_SECRET=local-dev-secret-at-least-32-characters
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here

# Application
APP_URL=http://localhost:4321
PORT=3000
```

### How to Generate JWT Secret

```bash
# Option 1: OpenSSL
openssl rand -base64 48

# Option 2: Go
go run -e 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 48); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

### Secure Environment File on VPS

```bash
# Set correct permissions (owner read/write only)
chmod 600 /home/campus/backend/.env

# Verify
ls -la /home/campus/backend/.env
# Should show: -rw------- 1 campus campus
```

---

## Monitoring & Maintenance

### Cloudflare Monitoring

**Cloudflare Dashboard:**
- **Pages → Analytics**: Page views, unique visitors, bandwidth
- **DNS → Analytics**: Request count, cache hit ratio

**Set up alerts:**
1. Cloudflare Dashboard → Notifications
2. Create alerts for:
   - Origin server errors (5xx)
   - SSL certificate expiration

### VPS Monitoring

**Check system resources:**
```bash
# SSH into VPS
ssh campus@YOUR_VPS_IP

# Check disk usage
df -h

# Check memory usage
free -h

# Check CPU usage
top -bn1 | head -5

# Check running services
sudo systemctl status campus-api postgresql nginx
```

**Check service logs:**
```bash
# Go backend logs
sudo journalctl -u campus-api -f

# Nginx access logs
sudo tail -f /var/log/nginx/access.log

# Nginx error logs
sudo tail -f /var/log/nginx/error.log

# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-18-main.log
```

### PostgreSQL Monitoring

**Check database status:**
```bash
# Connect to database
psql -U campus_app -d campus -h localhost

# Check table sizes
SELECT table_name,
       pg_size_pretty(pg_total_relation_size(table_name::regclass)) as size
FROM information_schema.tables
WHERE table_schema = 'public';

# Check row counts
SELECT 'users' as table_name, count(*) as rows FROM users
UNION ALL
SELECT 'applications', count(*) FROM applications
UNION ALL
SELECT 'sessions', count(*) FROM sessions;

# Check active connections
SELECT count(*) FROM pg_stat_activity WHERE datname = 'campus';
```

### Backup Strategy

**PostgreSQL automated backup:**

Create `/home/campus/scripts/backup.sh`:
```bash
#!/bin/bash
BACKUP_DIR="/home/campus/backups"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U campus_app -h localhost campus > "$BACKUP_DIR/campus_$DATE.sql"
# Keep only last 7 days
find "$BACKUP_DIR" -name "campus_*.sql" -mtime +7 -delete
```

**Schedule with cron:**
```bash
# Add to crontab (runs daily at 2 AM)
crontab -e
0 2 * * * /home/campus/scripts/backup.sh
```

**Manual backup:**
```bash
# Export full database
pg_dump -U campus_app -h localhost campus > backup.sql

# Export specific tables
pg_dump -U campus_app -h localhost -t users -t applications campus > backup.sql
```

**Restore from backup:**
```bash
# Restore full database
psql -U campus_app -h localhost campus < backup.sql
```

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

### Go Backend Issues

**Backend not starting:**
```bash
# Check service status
sudo systemctl status campus-api

# View error logs
sudo journalctl -u campus-api -n 50

# Common issues:
# - Missing .env file
# - Database connection failed
# - Port already in use

# Test manually
cd /home/campus/backend
./campus-api
```

**Backend returning errors:**
```bash
# Check logs in real-time
sudo journalctl -u campus-api -f

# Test health endpoint
curl http://localhost:3000/api/health

# Common issues:
# - Environment variables not loaded
# - Database connection timeout
# - Invalid JWT secret
```

**Authentication not working:**
```bash
# Check .env file has correct values:
# - GOOGLE_CLIENT_ID
# - GOOGLE_CLIENT_SECRET
# - JWT_SECRET

# Verify Google OAuth redirect URIs match exactly in Google Console
# Test OAuth flow: https://api.yourdomain.com/api/auth/google/login
```

### PostgreSQL Issues

**Connection errors:**
```bash
# Test connection
psql -U campus_app -d campus -h localhost

# Check PostgreSQL is running
sudo systemctl status postgresql

# Check pg_hba.conf for authentication rules
sudo cat /etc/postgresql/18/main/pg_hba.conf

# View PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-18-main.log
```

**Query timeout or slow queries:**
```sql
-- Check for missing indexes
\di

-- Check query execution plan
EXPLAIN ANALYZE SELECT * FROM applications WHERE user_id = 1;

-- Check active queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '5 seconds';
```

### Nginx Issues

**502 Bad Gateway:**
```bash
# Check if Go backend is running
curl http://localhost:3000/api/health

# Check Nginx configuration
sudo nginx -t

# View Nginx error logs
sudo tail -f /var/log/nginx/error.log

# Restart Nginx
sudo systemctl restart nginx
```

**SSL certificate issues:**
```bash
# Check certificate status
sudo certbot certificates

# Renew certificate
sudo certbot renew

# Force renewal
sudo certbot renew --force-renewal
```

### GitHub Actions Issues

**Deployment failing:**
```bash
# Check GitHub Actions logs
# Common issues:
# - Missing secrets (VPS_SSH_KEY, VPS_HOST, etc.)
# - SSH key permissions incorrect
# - VPS firewall blocking connection

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

**API latency:**
```bash
# Check Go backend response time
time curl http://localhost:3000/api/health

# Check PostgreSQL query time
psql -U campus_app -d campus -h localhost -c "EXPLAIN ANALYZE SELECT * FROM users LIMIT 10;"

# If latency is high:
# 1. Add indexes for slow queries
# 2. Check VPS CPU/memory usage
# 3. Optimize Go code with profiling
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

**Via CLI:**
```bash
cd frontend
git checkout <previous-commit>
npm run build
npx wrangler pages deploy dist --project-name=campus-website
```

### Rollback Go Backend

**Via SSH:**
```bash
# SSH into VPS
ssh campus@YOUR_VPS_IP

# Stop current service
sudo systemctl stop campus-api

# Rollback to previous binary (if kept)
mv ~/campus-api ~/campus-api.failed
mv ~/campus-api.old ~/campus-api

# Or rebuild from previous commit
cd ~/backend
git checkout <previous-commit>
go build -o campus-api ./cmd/server

# Restart
sudo systemctl start campus-api
sudo systemctl status campus-api
```

### Rollback Database Changes

**Option 1: Manual SQL rollback**
```sql
-- Connect to database
psql -U campus_app -d campus -h localhost

-- Run rollback commands manually
DROP INDEX IF EXISTS idx_new_index;
ALTER TABLE users DROP COLUMN IF EXISTS new_column;
```

**Option 2: Restore from backup**
```bash
# Restore from daily backup
psql -U campus_app -d campus -h localhost < /home/campus/backups/campus_YYYYMMDD.sql
```

**Option 3: Point-in-time recovery (requires WAL archiving)**
```bash
# If you have WAL archiving configured, restore to specific time
# This requires advanced PostgreSQL configuration
```

---

## Security Checklist

- [ ] All environment variables use strong, unique secrets
- [ ] JWT_SECRET is at least 32 characters, randomly generated
- [ ] PostgreSQL password is strong (20+ chars, mixed case, numbers, symbols)
- [ ] VPS firewall (ufw) enabled with only required ports (22, 80, 443)
- [ ] SSH root login disabled
- [ ] SSH key authentication enabled (password auth disabled)
- [ ] GitHub secrets configured (never commit secrets to repo)
- [ ] Google OAuth redirect URIs match exactly (production URLs only)
- [ ] .env file has restrictive permissions (chmod 600)
- [ ] Rate limiting enabled via Nginx
- [ ] SSL/HTTPS enforced (Let's Encrypt + Nginx)
- [ ] Cloudflare proxy enabled (VPS IP hidden)
- [ ] PostgreSQL listens only on localhost
- [ ] Regular security updates applied to VPS

---

## Support & Resources

### Documentation Links
- [Astro Documentation](https://docs.astro.build)
- [Cloudflare Pages](https://developers.cloudflare.com/pages)
- [Go Documentation](https://go.dev/doc/)
- [Chi Router](https://go-chi.io/)
- [pgx PostgreSQL Driver](https://github.com/jackc/pgx)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/18/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Certbot Documentation](https://certbot.eff.org/docs/)
- [Let's Encrypt](https://letsencrypt.org/docs/)

### Useful Commands

```bash
# Frontend Development
cd frontend && npm run dev          # Dev server with hot reload
cd frontend && npm run build        # Build for production
cd frontend && npm run preview      # Preview production build

# Backend Development
cd backend && go run ./cmd/server   # Run Go backend locally
cd backend && go build -o campus-api ./cmd/server  # Build binary

# Cloudflare Deployment
npx wrangler pages deploy dist      # Deploy Pages

# VPS Operations
ssh campus@YOUR_VPS_IP              # Connect to VPS
sudo systemctl status campus-api    # Check backend status
sudo journalctl -u campus-api -f    # Stream backend logs

# PostgreSQL
psql -U campus_app -d campus -h localhost  # Connect to local database
pg_dump -U campus_app campus > backup.sql  # Backup database
```

### Quick Reference

| Task | Command |
|------|---------|
| Start frontend dev | `cd frontend && npm run dev` |
| Start backend dev | `cd backend && go run ./cmd/server` |
| Build frontend | `cd frontend && npm run build` |
| Build backend | `cd backend && go build -o campus-api ./cmd/server` |
| Deploy frontend | `npx wrangler pages deploy dist` |
| Deploy backend | `ssh + git pull + go build + systemctl restart` |
| View backend logs | `sudo journalctl -u campus-api -f` |
| Connect to database | `psql -U campus_app -d campus -h localhost` |

---

**Last Updated:** 2026-01-15
**Version:** 3.0 (VPS + Go + PostgreSQL Architecture)
