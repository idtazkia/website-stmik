# Campus Website - Deployment Guide

## Table of Contents
- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
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
- [ ] Cloudflare account (for Pages and Workers)
- [ ] Google Cloud account (for OAuth OIDC)
- [ ] VPS account with local cloud provider

### Required Software
- Node.js 20+ (`node --version`)
- npm 9+ (`npm --version`)
- Git (`git --version`)
- PostgreSQL 14+ (on VPS)
- SSH client (for VPS access)

### VPS Requirements
- Ubuntu 22.04 LTS or similar
- 2GB RAM minimum
- 20GB storage minimum
- Public IP address
- Root or sudo access

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

## Local Development

### Option 1: Run Everything Together

```bash
# From project root
npm run dev

# This starts:
# - Frontend dev server (http://localhost:4321)
# - Backend dev server (http://localhost:3000)
```

### Option 2: Run Separately

**Terminal 1 - Frontend (Astro + Cloudflare Workers):**
```bash
cd frontend
npm install
cp .dev.vars.example .dev.vars  # Configure environment variables
npm run dev                      # http://localhost:4321
```

**Terminal 2 - Backend (Express.js):**
```bash
cd backend
npm install
cp .env.example .env             # Configure database, JWT secret
npm run dev                      # http://localhost:3000
```

**Terminal 3 - Database:**
```bash
# If running PostgreSQL locally
psql -U postgres
CREATE DATABASE campus;
\q

# Run migrations
cd backend
npm run migrate
```

### Local Testing Workflow

1. **Test static pages:** Visit `http://localhost:4321`
2. **Test authentication:** Try login/register flows
3. **Test API:** Use Postman or `curl` to test `http://localhost:3000/api/*`
4. **Test BFF:** Cloudflare Workers run locally via Wrangler

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

## Backend Deployment (VPS)

### 1. Prepare VPS

**Connect to VPS:**
```bash
ssh user@your-vps-ip
```

**Update system:**
```bash
sudo apt update && sudo apt upgrade -y
```

**Install Node.js 20:**
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node --version  # Verify v20.x
```

**Install PostgreSQL:**
```bash
sudo apt install postgresql postgresql-contrib -y
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

**Install pm2 (Process Manager):**
```bash
sudo npm install -g pm2
```

**Install nginx (Reverse Proxy):**
```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

### 2. Set Up Database

```bash
# Switch to postgres user
sudo -u postgres psql

-- Create database and user
CREATE DATABASE campus;
CREATE USER campus_user WITH ENCRYPTED PASSWORD 'your-strong-password';
GRANT ALL PRIVILEGES ON DATABASE campus TO campus_user;
\q
```

### 3. Clone Repository to VPS

```bash
# Create directory
sudo mkdir -p /var/www/campus-website
sudo chown $USER:$USER /var/www/campus-website

# Clone repository
cd /var/www
git clone https://github.com/yourorg/campus-website.git
cd campus-website
```

### 4. Configure Backend

```bash
cd backend

# Install dependencies
npm install --production

# Create .env file
cp .env.example .env
nano .env
```

**Edit .env:**
```bash
DATABASE_URL=postgresql://campus_user:your-strong-password@localhost:5432/campus
JWT_SECRET=generate-a-very-long-random-secret-here
PORT=3000
NODE_ENV=production
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

**Run migrations:**
```bash
npm run migrate
```

### 5. Start Backend with pm2

```bash
cd /var/www/campus-website/backend

# Start application
pm2 start src/index.js --name campus-backend

# Save pm2 configuration
pm2 save

# Set up pm2 to start on boot
pm2 startup
# Follow the command it outputs
```

**Verify backend is running:**
```bash
pm2 status
curl http://localhost:3000/health  # Should return 200 OK
```

### 6. Configure nginx Reverse Proxy

**Create nginx configuration:**
```bash
sudo nano /etc/nginx/sites-available/campus-backend
```

**Add configuration:**
```nginx
server {
    listen 80;
    server_name api.youruni.edu;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Enable site:**
```bash
sudo ln -s /etc/nginx/sites-available/campus-backend /etc/nginx/sites-enabled/
sudo nginx -t  # Test configuration
sudo systemctl reload nginx
```

### 7. Set Up SSL with Let's Encrypt

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx -y

# Get SSL certificate
sudo certbot --nginx -d api.youruni.edu

# Test auto-renewal
sudo certbot renew --dry-run
```

**nginx will be automatically configured for HTTPS**

### 8. Configure Firewall

```bash
# Allow SSH, HTTP, HTTPS
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw enable
sudo ufw status
```

---

## Automated Deployment (GitHub Actions)

### 1. Set Up GitHub Secrets

Go to GitHub Repository → Settings → Secrets and variables → Actions

Add the following secrets:

**For Frontend Deployment:**
- `CLOUDFLARE_API_TOKEN` - Your Cloudflare API token
- `CLOUDFLARE_ACCOUNT_ID` - Your Cloudflare account ID

**For Backend Deployment:**
- `VPS_HOST` - Your VPS IP address or hostname
- `VPS_USERNAME` - SSH username (e.g., `ubuntu`)
- `VPS_SSH_KEY` - Private SSH key for VPS access

**For Application:**
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `JWT_SECRET` - JWT signing secret
- `DATABASE_URL` - PostgreSQL connection string

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

### 3. Test Deployments

**Trigger frontend deployment:**
```bash
# Make a change to frontend
echo "# Test" >> frontend/README.md
git add frontend/README.md
git commit -m "Test frontend deployment"
git push origin main

# Check GitHub Actions tab for deployment status
```

**Trigger backend deployment:**
```bash
# Make a change to backend
echo "# Test" >> backend/README.md
git add backend/README.md
git commit -m "Test backend deployment"
git push origin main

# Check GitHub Actions tab for deployment status
```

---

## Environment Variables

### Frontend (.dev.vars for local, Cloudflare dashboard for production)

```bash
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here
BACKEND_URL=https://api.youruni.edu
APP_URL=https://youruni.edu
```

### Backend (.env on VPS)

```bash
# Database
DATABASE_URL=postgresql://campus_user:password@localhost:5432/campus

# Authentication
JWT_SECRET=generate-a-long-random-secret-min-32-chars
GOOGLE_CLIENT_ID=your-client-id-here
GOOGLE_CLIENT_SECRET=your-client-secret-here

# Server
PORT=3000
NODE_ENV=production

# Optional: Email (for notifications)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@youruni.edu
SMTP_PASS=your-app-password
```

### How to Generate JWT Secret

```bash
# Option 1: OpenSSL
openssl rand -base64 48

# Option 2: Node.js
node -e "console.log(require('crypto').randomBytes(48).toString('base64'))"
```

---

## Monitoring & Maintenance

### Frontend Monitoring (Cloudflare)

**Cloudflare Dashboard:**
- Pages → Analytics
- Workers → Metrics
- Monitor: Request count, errors, response time

**Set up alerts:**
- Cloudflare Dashboard → Notifications
- Create alerts for error rate spikes

### Backend Monitoring (VPS)

**pm2 monitoring:**
```bash
# View logs
pm2 logs campus-backend

# Monitor resources
pm2 monit

# Restart if needed
pm2 restart campus-backend
```

**nginx logs:**
```bash
# Access logs
sudo tail -f /var/log/nginx/access.log

# Error logs
sudo tail -f /var/log/nginx/error.log
```

**PostgreSQL monitoring:**
```bash
# Connect to database
psql -U campus_user -d campus

-- Check active connections
SELECT count(*) FROM pg_stat_activity;

-- Check database size
SELECT pg_size_pretty(pg_database_size('campus'));

-- Check slow queries
SELECT query, mean_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### Backup Strategy

**Database backups:**
```bash
# Create backup script
sudo nano /usr/local/bin/backup-db.sh
```

```bash
#!/bin/bash
BACKUP_DIR=/var/backups/postgres
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

pg_dump -U campus_user campus | gzip > $BACKUP_DIR/campus_$DATE.sql.gz

# Keep only last 7 days
find $BACKUP_DIR -name "campus_*.sql.gz" -mtime +7 -delete
```

```bash
# Make executable
sudo chmod +x /usr/local/bin/backup-db.sh

# Schedule daily backups via cron
crontab -e
# Add: 0 2 * * * /usr/local/bin/backup-db.sh
```

---

## Troubleshooting

### Frontend Issues

**Build fails on Cloudflare:**
```bash
# Check build logs in Cloudflare Dashboard
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
```

### Backend Issues

**Backend not starting:**
```bash
# Check pm2 logs
pm2 logs campus-backend --lines 100

# Check if port is in use
sudo lsof -i :3000

# Check environment variables
pm2 env campus-backend
```

**Database connection issues:**
```bash
# Test database connection
psql -U campus_user -d campus -h localhost

# Check PostgreSQL is running
sudo systemctl status postgresql

# Check database credentials in .env
cat backend/.env | grep DATABASE_URL
```

**502 Bad Gateway from nginx:**
```bash
# Check if backend is running
pm2 status

# Check nginx error logs
sudo tail -f /var/log/nginx/error.log

# Test backend directly
curl http://localhost:3000/health
```

### GitHub Actions Issues

**Deployment failing:**
```bash
# Check GitHub Actions logs
# Common issues:
# - Missing secrets
# - SSH key format incorrect
# - VPS unreachable

# Test SSH connection manually
ssh -i path/to/key user@vps-ip
```

**Path triggers not working:**
```yaml
# Ensure paths are correct
paths:
  - 'frontend/**'  # Note: no leading slash
  - 'shared/**'
```

### Performance Issues

**Slow API responses:**
```sql
-- Check for missing indexes
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE tablename = 'applications'
ORDER BY n_distinct DESC;

-- Add indexes if needed
CREATE INDEX idx_applications_user_id ON applications(user_id);
CREATE INDEX idx_applications_status ON applications(status);
```

**High memory usage:**
```bash
# Check pm2 memory usage
pm2 monit

# Restart if needed
pm2 restart campus-backend
```

---

## Rollback Procedures

### Rollback Frontend

```bash
# Via Cloudflare Dashboard
# Pages → Deployments → Select previous deployment → Rollback

# Via Wrangler
cd frontend
git checkout <previous-commit>
npm run build
npx wrangler pages deploy dist
```

### Rollback Backend

```bash
# SSH to VPS
ssh user@vps

# Navigate to project
cd /var/www/campus-website

# Checkout previous version
git log --oneline  # Find previous commit
git checkout <commit-hash>

# Reinstall dependencies
cd backend
npm install --production

# Restart
pm2 restart campus-backend
```

### Rollback Database Migration

```bash
# If migrations support down migration
cd backend
npm run migrate:rollback

# Otherwise, restore from backup
gunzip < /var/backups/postgres/campus_YYYYMMDD.sql.gz | psql -U campus_user campus
```

---

## Security Checklist

- [ ] All environment variables use strong, unique secrets
- [ ] JWT_SECRET is at least 32 characters, randomly generated
- [ ] Database passwords are strong (20+ chars, mixed case, numbers, symbols)
- [ ] SSH key authentication enabled, password auth disabled
- [ ] Firewall configured (ufw enable, only necessary ports open)
- [ ] SSL certificates installed and auto-renewing
- [ ] GitHub secrets configured (never commit secrets to repo)
- [ ] Google OAuth redirect URIs match exactly (production URLs only)
- [ ] CORS configured to allow only your frontend domain
- [ ] Rate limiting enabled on authentication endpoints
- [ ] Database backups running daily
- [ ] pm2 startup configured (backend restarts after reboot)

---

## Support & Resources

### Documentation Links
- [Astro Documentation](https://docs.astro.build)
- [Cloudflare Pages](https://developers.cloudflare.com/pages)
- [Cloudflare Workers](https://developers.cloudflare.com/workers)
- [Express.js Guide](https://expressjs.com)
- [PostgreSQL Docs](https://www.postgresql.org/docs)
- [pm2 Documentation](https://pm2.keymetrics.io)
- [nginx Documentation](https://nginx.org/en/docs)

### Useful Commands

```bash
# Frontend
cd frontend && npm run build        # Build Astro site
cd frontend && npm run dev          # Dev server
npx wrangler pages deploy dist      # Deploy to Cloudflare

# Backend
cd backend && npm run dev           # Dev server
cd backend && npm run migrate       # Run migrations
pm2 logs campus-backend             # View logs
pm2 restart campus-backend          # Restart

# System
sudo systemctl status nginx         # Check nginx
sudo systemctl status postgresql    # Check database
sudo ufw status                     # Check firewall
df -h                              # Check disk space
free -h                            # Check memory
```

---

**Last Updated:** 2025-11-19
**Version:** 1.0
