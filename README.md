# Campus Website - STMIK Tazkia

> Modern, secure, and cost-effective campus admission system built with Astro, Express.js, and Cloudflare.

**🌐 Live Site:** https://dev.stmik.tazkia.ac.id/

## 🎯 Current Status (Updated: 2025-11-19)

**Phase:** Marketing Site Foundation (Phase 3 - 30% Complete)

### ✅ Completed
- Bilingual static site (Indonesian & English)
- Homepage with hero, features, programs overview
- About page (vision, mission, programs)
- Lecturer profiles with content collections
- Responsive design with Tailwind CSS 4.x
- Custom i18n system (no external dependencies)
- SEO optimization (meta tags, sitemap, Open Graph)
- Cloudflare Pages deployment with auto-deploy

### 🚧 Next Steps
- Programs detail pages
- Contact page
- Admissions information page
- News/blog system
- Backend API development
- Authentication system
- Application portal

---

## 📋 Project Overview

Campus website for marketing and admission processing, serving:
- **Registrants** (prospective students): Submit applications via Google SSO or email/password
- **Marketing Staff**: Manage applications and content via Google SSO only

**Target Scale:** 3,000 leads per admission cycle (300 registrations at 10% conversion), 5 staff members
**Monthly Cost:** $5 (VPS with Go + PostgreSQL)

---

## 🏗️ Architecture

**Static Site + Go Backend Pattern**

```mermaid
graph TB
    Browser[🌐 Browser]

    subgraph "Edge Layer - Cloudflare (FREE)"
        CDN[🛡️ Cloudflare CDN<br/>DDoS Protection]
        Pages[📄 Cloudflare Pages<br/>Static Marketing Site]
    end

    subgraph "VPS Layer ($5/month)"
        Nginx[🔒 Nginx<br/>Reverse Proxy + SSL]
        GoAPI[🚀 Go Backend<br/>REST API]
        Database[(🗄️ PostgreSQL<br/>Database)]
    end

    Browser -->|HTTPS| CDN
    CDN -->|Static| Pages
    CDN -->|API| Nginx
    Nginx -->|Proxy| GoAPI
    GoAPI -->|SQL| Database

    style Browser fill:#e1f5ff
    style CDN fill:#c8e6c9
    style Pages fill:#c8e6c9
    style Nginx fill:#fff9c4
    style GoAPI fill:#fff9c4
    style Database fill:#fff9c4
```

**Why This Stack?**
- ✅ **Cost-effective**: Fixed $5/month (no usage surprises)
- ✅ **DDoS-proof**: Cloudflare proxy, VPS IP hidden
- ✅ **Simple**: Single server, full control
- ✅ **Fast**: Go + local PostgreSQL (<1ms latency)
- ✅ **Scalable**: Handles 100,000+ leads on $5 VPS

📖 **[Read Full Architecture Documentation →](docs/ARCHITECTURE.md)**

---

## 🚀 Tech Stack

| Component | Technology | Hosting | Cost |
|-----------|-----------|---------|------|
| **Static Site** | Astro + Markdown | Cloudflare Pages | Free |
| **CDN/DDoS** | Cloudflare | Cloudflare | Free |
| **Backend API** | Go (Golang) | VPS | $5/mo |
| **Database** | PostgreSQL 18 | VPS | Included |
| **Reverse Proxy** | Nginx + Let's Encrypt | VPS | Included |
| **Build/Deploy** | GitHub Actions | GitHub | Free |

---

## 📁 Repository Structure

```
website-stmik/
├── frontend/                         # Astro static site
│   ├── src/
│   │   ├── content/                 # Markdown content (lecturers, programs, etc.)
│   │   ├── pages/                   # Astro pages (bilingual routing)
│   │   ├── components/              # Reusable UI components
│   │   ├── layouts/                 # Page layouts (Base, Marketing)
│   │   ├── styles/                  # Global styles (Tailwind CSS 4.x)
│   │   └── utils/                   # Utilities (i18n, etc.)
│   ├── public/
│   │   ├── images/                  # Static images
│   │   └── locales/                 # Translation JSON files (id/en)
│   ├── astro.config.mjs             # Astro configuration
│   ├── TODO.md                      # Frontend implementation tasks
│   └── package.json
├── backend/                          # Express.js API (Phase 2)
│   └── TODO.md                      # Backend implementation plan
├── shared/                           # Shared TypeScript types (Phase 5)
│   └── TODO.md                      # Shared code plan
├── infrastructure/                   # Infrastructure as Code
│   ├── ansible/                     # Ansible playbooks for VPS
│   │   ├── playbooks/              # Deployment automation
│   │   ├── inventory/              # Server inventory
│   │   ├── roles/                  # Reusable Ansible roles
│   │   └── README.md               # Ansible usage guide
│   ├── scripts/                     # Deployment scripts
│   ├── TODO.md                      # Infrastructure tasks
│   └── README.md                    # Infrastructure overview
├── docs/                             # Documentation
│   ├── ARCHITECTURE.md              # Technical design details
│   └── DEPLOYMENT.md                # Deployment guide
├── tests/                            # E2E tests (Playwright)
│   └── deployment-check.spec.ts    # Browser tests for deployed site
├── TODO.md                           # High-level project overview
├── CLAUDE.md                         # Claude Code guidance
├── playwright.config.ts             # Playwright configuration
└── README.md                         # This file
```

**Current Focus:** Frontend marketing site (Phase 3)
**Deferred:** Backend, shared types, infrastructure automation (Phase 2+)

---

## 🔐 Authentication

**Hybrid OIDC + Traditional Authentication**

### For Registrants (Prospective Students)
- **Google OIDC** (recommended): One-click, no password management
- **Email/Password** (traditional): Privacy-conscious option

### For Marketing Staff
- **Google OIDC only**: Enforced via `@youruni.edu` domain check

**Security Features:**
- HttpOnly cookies (XSS protection)
- JWT tokens (7-day expiration)
- Rate limiting (brute force prevention)
- DDoS protection (Cloudflare edge)
- bcrypt password hashing
- SQL injection prevention

📖 **[Read Authentication Details →](docs/ARCHITECTURE.md#authentication-system)**

---

## 💰 Cost Breakdown

### Monthly Costs

| Service | Usage | Cost |
|---------|-------|------|
| **Cloudflare Pages** | Unlimited bandwidth | **$0** |
| **Cloudflare CDN** | DDoS protection, caching | **$0** |
| **VPS (1GB RAM)** | Go + PostgreSQL + Nginx | **$5** |
| **GitHub Actions** | Build/deploy automation | **$0** |
| **Google OAuth** | OIDC authentication | **$0** |
| **Let's Encrypt** | SSL certificate | **$0** |

**Total: $5/month** 🎉

### Capacity Analysis (3,000 Leads)

- **VPS CPU usage:** <5% (Go is efficient)
- **VPS RAM usage:** ~300MB of 1GB (~30%)
- **Database size:** ~10MB (tiny)
- **Can scale to:** 100,000+ leads without upgrade
- **DDoS attack cost:** $5 (fixed, Cloudflare absorbs traffic)

📖 **[Read Traffic Analysis →](docs/ARCHITECTURE.md#bff-traffic-analysis)**

---

## 🚦 Quick Start

### Prerequisites
- Node.js 20+
- npm or yarn

### Local Development (Frontend Only - Currently Available)

```bash
# Clone repository
git clone https://github.com/idtazkia/website-stmik.git
cd website-stmik

# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Visit: http://localhost:4321
# - Homepage (ID): http://localhost:4321/
# - Homepage (EN): http://localhost:4321/en/
# - About: http://localhost:4321/about
# - Lecturers: http://localhost:4321/lecturers
```

### Build for Production

```bash
cd frontend
npm run build        # Build static site
npm run preview      # Preview production build
```

### For Backend Development (Coming Soon)
Backend API, authentication, and application portal are planned for Phase 2. See [TODO.md](TODO.md) for the implementation roadmap.

📖 **[Read Full Deployment Guide →](docs/DEPLOYMENT.md)**

---

## 📝 Implementation Checklist

See **[TODO.md](TODO.md)** for the complete implementation roadmap.

### Quick Status

**Phase 1: Project Setup** - [x] 40% Complete
- [x] Repository initialization
- [x] Cloudflare Pages setup
- [x] Custom domain configured
- [ ] VPS setup (deferred)
- [ ] Google OAuth setup (pending)

**Phase 2: Backend Development** - [ ] 0% Complete
- [ ] Database schema
- [ ] Authentication system
- [ ] API endpoints
- [ ] File upload system

**Phase 3: Frontend Development** - [x] 30% Complete
- [x] Astro site setup with Tailwind CSS
- [x] Bilingual system (ID/EN)
- [x] Homepage, About, Lecturers pages
- [ ] Programs, Contact, Admissions pages
- [ ] Application portal
- [ ] Admin interface

**Phase 4: BFF Layer** - [ ] 0% Complete
- [ ] Authentication handlers
- [ ] API proxy functions

**Phase 5-8:** Shared code, deployment, testing, launch (0% complete)

📖 **[View Full Checklist →](TODO.md)**

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** | Technical design, authentication flows, database schema, security, scalability |
| **[DEPLOYMENT.md](docs/DEPLOYMENT.md)** | Step-by-step deployment guide, environment setup, monitoring, troubleshooting |
| **[TODO.md](TODO.md)** | Complete implementation checklist with 8 phases and timeline |

---

## 🎯 Key Features

### For Registrants
- Browse programs and admission requirements (static pages)
- Create account (Google or email/password)
- Fill out application form
- Upload supporting documents
- Track application status

### For Marketing Staff
- Login via Google Workspace
- Review submitted applications
- Approve or reject applications
- Manage user accounts
- View analytics dashboard

### Technical Features
- Static site generation (SEO-optimized)
- Real-time form validation
- Auto-save drafts
- File upload with validation
- Responsive design (mobile-friendly)
- Accessibility compliant

---

## 🔒 Security

- ✅ **HttpOnly Cookies** - XSS-resistant token storage
- ✅ **HTTPS Enforced** - All connections encrypted
- ✅ **Rate Limiting** - Prevent brute force attacks
- ✅ **OIDC Authentication** - Google's security infrastructure
- ✅ **DDoS Protection** - Cloudflare edge filtering
- ✅ **Input Validation** - Prevent injection attacks
- ✅ **SQL Injection Prevention** - Parameterized queries
- ✅ **CORS Configuration** - Restricted origins

📖 **[Read Security Details →](docs/ARCHITECTURE.md#security-considerations)**

---

## 📈 Scalability

### Current Capacity (3,000 leads)
- VPS using <5% CPU, ~30% RAM
- PostgreSQL using <1% of disk
- Can scale to 100,000+ leads on same VPS

### Scaling Path
1. **0-100K leads:** Current 1GB VPS ($5/mo)
2. **100K-500K leads:** Upgrade to 2-4GB VPS ($10-20/mo)
3. **500K+ leads:** Load balancer + multiple VPS ($50+/mo)

📖 **[Read Scalability Strategy →](docs/ARCHITECTURE.md#scalability)**

---

## 🛠️ Development

### Commands

```bash
# Frontend development (current)
cd frontend
npm install                    # Install dependencies
npm run dev                    # Start dev server (http://localhost:4321)
npm run build                  # Build for production
npm run preview                # Preview production build

# Code quality
npm run typecheck              # Run TypeScript checks
npm run lint                   # Run ESLint
npm run format                 # Format with Prettier

# Deployment
git push                       # Cloudflare Pages auto-deploys on push
```

**Backend & Monorepo Commands (coming in Phase 2):**
Backend API, database migrations, and npm workspace commands will be available after Phase 2 implementation.

---

## 🐛 Troubleshooting

### Common Issues

**Build fails:**
```bash
# Clear node_modules and reinstall
rm -rf node_modules frontend/node_modules backend/node_modules
npm install
```

**Database connection error:**
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U campus_app -d campus -h localhost

# Check logs
sudo tail -f /var/log/postgresql/postgresql-18-main.log
```

**Go backend not responding:**
```bash
# Check service status
sudo systemctl status campus-api

# View logs
sudo journalctl -u campus-api -f

# Restart service
sudo systemctl restart campus-api
```

📖 **[Read Full Troubleshooting Guide →](docs/DEPLOYMENT.md#troubleshooting)**

---

## 📅 Timeline

**Estimated time to MVP:** 7-10 weeks

- **Phase 1-2:** Backend foundation (1-2 weeks)
- **Phase 3-4:** Frontend + BFF (2-3 weeks)
- **Phase 5:** Shared code (3-5 days)
- **Phase 6:** Deployment (1 week)
- **Phase 7:** Testing & polish (1-2 weeks)
- **Phase 8:** Launch prep (1 week)

📖 **[View Detailed Timeline →](TODO.md#notes)**

---

## 🤝 Contributing

This is a private project for STMIK Campus. Only authorized developers should contribute.

### Workflow
1. Create feature branch from `main`
2. Implement changes
3. Test locally
4. Create pull request
5. After review, merge to `main`
6. GitHub Actions deploys automatically

---

## 📞 Support

### For Development Issues
- Check **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** for technical questions
- Check **[DEPLOYMENT.md](docs/DEPLOYMENT.md)** for deployment issues
- Check **[TODO.md](TODO.md)** for implementation guidance

### For Production Issues
- Check Cloudflare Analytics for traffic/errors
- Check Go logs: `sudo journalctl -u campus-api -f`
- Check Nginx logs: `sudo tail -f /var/log/nginx/error.log`
- Check PostgreSQL: `sudo tail -f /var/log/postgresql/*.log`

---

## 📄 License

This project is proprietary and confidential. All rights reserved to STMIK.

---

## 🎓 Project Info

**Architecture Design:** 2025
**Version:** 1.0
**Target Deployment:** Q1 2025
**Maintenance:** Active

---

## ✨ Highlights

### Why This Solution?

1. **Cost-Effective:** Fixed $5/month (no usage surprises)
2. **Simple:** Single VPS to manage (not 3+ cloud services)
3. **Fast:** Go backend + local PostgreSQL (<1ms latency)
4. **Scalable:** Can grow from 3,000 to 100,000+ leads on same VPS
5. **Secure:** OIDC, HttpOnly cookies, Cloudflare DDoS protection
6. **Developer-Friendly:** Modern stack (Astro, Go), simple deployment
7. **SEO-Optimized:** Static site generation for excellent search rankings
8. **Full Control:** Direct server access, easy debugging

### Success Metrics

- [ ] 3,000 leads with 300 registrations (10% conversion)
- [ ] Zero security incidents
- [ ] 99% uptime
- [ ] <2s page load time
- [ ] <30% VPS resource usage
- [ ] $5/month hosting cost maintained

---

**Ready to start?** Check out the **[TODO.md](TODO.md)** for the implementation plan!
