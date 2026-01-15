# Infrastructure TODO

**Status:** Not Started (Deferred to Phase 2)

This directory will contain Ansible playbooks and scripts for VPS provisioning and deployment automation.

---

## Phase 1: Ansible Setup

### Prerequisites
- [ ] Choose VPS provider (Hetzner, Vultr, or local Indonesian provider)
- [ ] Provision Ubuntu 24.04 LTS VPS (2GB RAM minimum)
- [ ] Obtain VPS root access credentials
- [ ] Configure DNS records for backend subdomain (api.stmik.tazkia.ac.id)

### Ansible Installation (Local Machine)
- [ ] Install Ansible (2.14+)
- [ ] Install required Ansible collections
  - [ ] `ansible-galaxy collection install community.general`
  - [ ] `ansible-galaxy collection install community.postgresql`
- [ ] Test Ansible installation: `ansible --version`

### SSH Configuration
- [ ] Generate deployment SSH key
- [ ] Copy SSH key to VPS: `ssh-copy-id user@vps-ip`
- [ ] Test SSH connection: `ssh user@vps-ip`
- [ ] Add VPS to known_hosts

---

## Phase 2: Ansible Playbook Development

### 2.1 Initial VPS Setup Playbook (`playbooks/setup-vps.yml`)
- [ ] System updates and security patches
  - [ ] Update apt packages
  - [ ] Enable automatic security updates
  - [ ] Configure timezone (Asia/Jakarta)
- [ ] Create deployment user
  - [ ] Create `deploy` user with sudo access
  - [ ] Configure SSH key authentication
  - [ ] Disable root SSH login
  - [ ] Disable password authentication
- [ ] Install Node.js 20.x LTS
  - [ ] Add NodeSource repository
  - [ ] Install Node.js and npm
  - [ ] Verify installation
- [ ] Install PostgreSQL 18
  - [ ] Add PostgreSQL APT repository
  - [ ] Install PostgreSQL
  - [ ] Configure PostgreSQL for remote connections (optional)
  - [ ] Create database and user
- [ ] Install Nginx
  - [ ] Install nginx package
  - [ ] Configure as reverse proxy
  - [ ] Enable and start nginx
- [ ] Install PM2
  - [ ] Install PM2 globally: `npm install -g pm2`
  - [ ] Configure PM2 startup script
- [ ] Install Certbot for SSL
  - [ ] Install certbot and nginx plugin
  - [ ] Configure certbot auto-renewal
- [ ] Configure UFW firewall
  - [ ] Allow SSH (port 22)
  - [ ] Allow HTTP (port 80)
  - [ ] Allow HTTPS (port 443)
  - [ ] Enable UFW
- [ ] Install Fail2ban
  - [ ] Install fail2ban package
  - [ ] Configure SSH protection
  - [ ] Enable fail2ban service

### 2.2 Database Setup Playbook (`playbooks/setup-database.yml`)
- [ ] Create PostgreSQL database
  - [ ] Database name: `campus`
  - [ ] Database user: `campus_user`
  - [ ] Generate strong password
- [ ] Configure PostgreSQL
  - [ ] Set max_connections
  - [ ] Configure shared_buffers
  - [ ] Enable logging
- [ ] Run database migrations
  - [ ] Clone repository to VPS
  - [ ] Run migration script
  - [ ] Seed initial data (if needed)
- [ ] Verify database
  - [ ] Test connection
  - [ ] Check tables created

### 2.3 Backend Deployment Playbook (`playbooks/deploy-backend.yml`)
- [ ] Clone/pull repository
  - [ ] Git clone if first deployment
  - [ ] Git pull if updating
  - [ ] Checkout specific branch/tag
- [ ] Install dependencies
  - [ ] Run `npm ci --production`
  - [ ] Verify package-lock.json
- [ ] Configure environment
  - [ ] Create `.env` file from template
  - [ ] Set DATABASE_URL
  - [ ] Set JWT_SECRET
  - [ ] Set NODE_ENV=production
- [ ] Run database migrations
  - [ ] Execute migration command
  - [ ] Verify migration success
- [ ] Start/restart application with PM2
  - [ ] PM2 start or restart
  - [ ] Save PM2 configuration
  - [ ] Verify application running
- [ ] Configure Nginx
  - [ ] Create Nginx site configuration
  - [ ] Enable site
  - [ ] Test Nginx config
  - [ ] Reload Nginx
- [ ] SSL certificate
  - [ ] Run certbot for domain
  - [ ] Verify SSL certificate
  - [ ] Test HTTPS access

### 2.4 Maintenance Playbook (`playbooks/maintenance.yml`)
- [ ] Database backup
  - [ ] Run pg_dump
  - [ ] Compress backup
  - [ ] Store in /var/backups/
  - [ ] Clean old backups (keep 7 days)
- [ ] System updates
  - [ ] Update apt packages
  - [ ] Update npm packages (backend)
  - [ ] Restart services if needed
- [ ] Health checks
  - [ ] Check PM2 status
  - [ ] Check Nginx status
  - [ ] Check PostgreSQL status
  - [ ] Check disk space
  - [ ] Check memory usage
- [ ] Log rotation
  - [ ] Configure logrotate for PM2 logs
  - [ ] Configure logrotate for Nginx logs
- [ ] SSL certificate renewal
  - [ ] Run certbot renew
  - [ ] Reload Nginx if renewed

---

## Phase 3: Ansible Roles (Optional - for reusability)

Create reusable Ansible roles:

### roles/nodejs/
- [ ] tasks/main.yml - Install Node.js
- [ ] vars/main.yml - Node.js version
- [ ] handlers/main.yml - Restart services

### roles/postgresql/
- [ ] tasks/main.yml - Install PostgreSQL
- [ ] tasks/database.yml - Create database
- [ ] templates/pg_hba.conf.j2 - PostgreSQL config
- [ ] vars/main.yml - Database settings

### roles/nginx/
- [ ] tasks/main.yml - Install Nginx
- [ ] templates/backend.conf.j2 - Site config
- [ ] handlers/main.yml - Reload Nginx

### roles/pm2/
- [ ] tasks/main.yml - Install PM2
- [ ] tasks/deploy.yml - Deploy application
- [ ] templates/ecosystem.config.js.j2 - PM2 config

### roles/ssl-certbot/
- [ ] tasks/main.yml - Install Certbot
- [ ] tasks/certificate.yml - Get SSL cert
- [ ] handlers/main.yml - Reload Nginx

---

## Phase 4: Deployment Scripts

### scripts/deploy.sh
```bash
#!/bin/bash
# Wrapper script for deployment
# Usage: ./scripts/deploy.sh production
```

- [ ] Parse environment argument (production/staging)
- [ ] Run Ansible playbook
- [ ] Show deployment status
- [ ] Run post-deployment checks

### scripts/backup.sh
```bash
#!/bin/bash
# Manual backup trigger
# Usage: ./scripts/backup.sh
```

- [ ] Run backup playbook
- [ ] Download backup to local machine
- [ ] Verify backup integrity

### scripts/rollback.sh
```bash
#!/bin/bash
# Rollback to previous version
# Usage: ./scripts/rollback.sh
```

- [ ] Stop PM2 application
- [ ] Git checkout previous commit
- [ ] Restore database backup
- [ ] Restart application

---

## Phase 5: Inventory & Variables

### inventory/production.ini
```ini
[webservers]
backend-vps ansible_host=YOUR_VPS_IP ansible_user=deploy

[databases]
backend-vps ansible_host=YOUR_VPS_IP ansible_user=deploy
```

- [ ] Add production VPS IP
- [ ] Configure SSH user
- [ ] Set Python interpreter path

### group_vars/all.yml
```yaml
nodejs_version: "20.x"
postgresql_version: "18"
app_name: "campus-backend"
app_port: 3000
domain: "api.stmik.tazkia.ac.id"
```

- [ ] Define global variables
- [ ] Set versions for all software
- [ ] Configure domain names

### group_vars/production.yml
```yaml
env: production
database_name: campus
database_user: campus_user
backup_retention_days: 7
```

- [ ] Production-specific variables
- [ ] Database credentials (encrypted with ansible-vault)
- [ ] Backup settings

---

## Phase 6: Ansible Configuration

### ansible.cfg
```ini
[defaults]
inventory = inventory/production.ini
remote_user = deploy
host_key_checking = False
retry_files_enabled = False

[privilege_escalation]
become = True
become_method = sudo
become_user = root
```

- [ ] Configure default inventory
- [ ] Set SSH options
- [ ] Configure privilege escalation

---

## Phase 7: Testing & Documentation

### Testing
- [ ] Test playbooks against staging VPS
- [ ] Verify idempotence (run playbook twice)
- [ ] Test rollback procedure
- [ ] Document common errors and fixes

### Documentation
- [ ] Write Ansible usage guide (ansible/README.md)
- [ ] Document deployment process
- [ ] Create troubleshooting guide
- [ ] Add runbook for common tasks

---

## Timeline Estimate

- **Phase 1:** Ansible setup - 0.5 day
- **Phase 2:** Playbook development - 2-3 days
- **Phase 3:** Ansible roles (optional) - 1-2 days
- **Phase 4:** Deployment scripts - 0.5 day
- **Phase 5:** Inventory & variables - 0.5 day
- **Phase 6:** Configuration - 0.5 day
- **Phase 7:** Testing & docs - 1 day

**Total: 5-8 days** (can be done in parallel with backend development)

---

## Success Criteria

- [ ] Ansible can provision fresh VPS from scratch
- [ ] Backend deploys successfully via Ansible
- [ ] Database backups run automatically
- [ ] SSL certificates auto-renew
- [ ] Rollback procedure tested and working
- [ ] Documentation complete and accurate

---

## Dependencies

- **Backend Phase 2:** Backend application must be developed first
- **Domain Configuration:** DNS records must point to VPS
- **VPS Provider:** VPS must be provisioned and accessible

---

## Security Checklist

- [ ] SSH key-only authentication (no passwords)
- [ ] UFW firewall configured
- [ ] Fail2ban protecting SSH
- [ ] PostgreSQL strong password (ansible-vault encrypted)
- [ ] JWT_SECRET strong and random (ansible-vault encrypted)
- [ ] SSL/TLS certificates installed
- [ ] Automatic security updates enabled
- [ ] Root login disabled
- [ ] Non-root deployment user created

---

## Future Enhancements (Post-MVP)

- [ ] Multi-server setup (load balancer + multiple backends)
- [ ] Database replication (primary + replica)
- [ ] Monitoring with Prometheus + Grafana
- [ ] Log aggregation with ELK stack
- [ ] Automated performance testing
- [ ] Blue-green deployment strategy
- [ ] Docker containerization (optional)
- [ ] Kubernetes deployment (if scale requires)
