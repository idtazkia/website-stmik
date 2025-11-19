# Infrastructure - STMIK Tazkia

Infrastructure-as-Code (IaC) for STMIK Tazkia campus website deployment and management.

## Overview

This directory contains all infrastructure and deployment automation:
- **Ansible playbooks** for VPS provisioning and configuration
- **Deployment scripts** for automated releases
- **Maintenance scripts** for backups and updates

**Target Platform:** Ubuntu 24.04 LTS (or 22.04 LTS)
**Target Environment:** VPS (2GB RAM minimum)
**Monthly Cost:** $5-10

---

## Directory Structure

```
infrastructure/
├── ansible/                    # Ansible automation
│   ├── playbooks/             # Ansible playbooks
│   ├── inventory/             # Server inventory files
│   ├── roles/                 # Reusable Ansible roles
│   ├── group_vars/            # Variable files
│   ├── ansible.cfg            # Ansible configuration
│   └── README.md              # Ansible usage guide
├── scripts/                   # Deployment & maintenance scripts
│   ├── deploy.sh              # Deployment wrapper
│   ├── backup.sh              # Database backup
│   └── rollback.sh            # Rollback deployment
└── README.md                  # This file
```

---

## Quick Start

### Prerequisites

1. **Ansible installed locally:**
   ```bash
   # macOS
   brew install ansible

   # Ubuntu/Debian
   sudo apt update && sudo apt install ansible

   # Verify
   ansible --version  # Should be 2.14+
   ```

2. **SSH access to VPS:**
   ```bash
   # Generate SSH key if needed
   ssh-keygen -t ed25519 -C "deployment@stmik-tazkia"

   # Copy public key to VPS
   ssh-copy-id -i ~/.ssh/id_ed25519.pub user@your-vps-ip
   ```

3. **VPS requirements:**
   - Ubuntu 24.04 LTS (or 22.04 LTS)
   - 2GB RAM minimum
   - 20GB storage minimum
   - Root or sudo access

---

## Usage

### 1. Initial VPS Setup

Provision a fresh VPS with all required software:

```bash
cd infrastructure/ansible
ansible-playbook -i inventory/production.ini playbooks/setup-vps.yml
```

This installs:
- Node.js 20.x
- PostgreSQL 16
- Nginx
- PM2
- Certbot (SSL)

### 2. Deploy Backend

Deploy the Express.js backend application:

```bash
cd infrastructure/ansible
ansible-playbook -i inventory/production.ini playbooks/deploy-backend.yml
```

### 3. Database Setup

Initialize PostgreSQL database:

```bash
cd infrastructure/ansible
ansible-playbook -i inventory/production.ini playbooks/setup-database.yml
```

### 4. Maintenance

Run backups, updates, and health checks:

```bash
cd infrastructure/ansible
ansible-playbook -i inventory/production.ini playbooks/maintenance.yml
```

---

## Components Managed

### VPS Software Stack

| Component | Version | Purpose |
|-----------|---------|---------|
| Ubuntu LTS | 24.04 | Operating system |
| Node.js | 20.x LTS | JavaScript runtime |
| PostgreSQL | 16 | Database |
| Nginx | Latest | Reverse proxy |
| PM2 | Latest | Process manager |
| Certbot | Latest | SSL certificates |

### Application Components

- **Backend API:** Express.js application
- **Database:** PostgreSQL with automatic backups
- **SSL/TLS:** Let's Encrypt certificates
- **Process Management:** PM2 with auto-restart

---

## Security Features

- ✅ UFW firewall (ports 22, 80, 443 only)
- ✅ SSH key-only authentication (password login disabled)
- ✅ Automatic security updates
- ✅ SSL/TLS encryption (Let's Encrypt)
- ✅ PostgreSQL with strong passwords
- ✅ Fail2ban for brute-force protection

---

## Current Status

**Phase:** Not Started (Deferred to Phase 2)

See `TODO.md` for implementation checklist.

---

## Documentation

- **Ansible Setup:** See `ansible/README.md`
- **Deployment Guide:** See `../docs/DEPLOYMENT.md`
- **Architecture:** See `../docs/ARCHITECTURE.md`

---

## Cost Breakdown

**VPS Requirements (2GB RAM):**
- Digital Ocean: $12/month
- Vultr: $10/month
- Hetzner: $5/month (EU)
- Local Indonesian providers: $5-10/month

**Total Infrastructure Cost:** $5-10/month

---

## Support

For infrastructure issues:
1. Check Ansible logs: `ansible-playbook -vvv ...`
2. Check VPS logs: `ssh user@vps "sudo journalctl -u backend"`
3. Review docs: `docs/DEPLOYMENT.md`
