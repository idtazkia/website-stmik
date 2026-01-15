# Infrastructure TODO

**Status:** Not Started (Deferred to Phase 2)

This directory will contain Ansible playbooks and scripts for VPS provisioning and deployment automation.

---

## Phase 1: Ansible Setup

### Prerequisites
- [ ] Choose VPS provider (Hetzner, Vultr, or local Indonesian provider)
- [ ] Provision Ubuntu 24.04 LTS VPS (1GB RAM minimum)
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
- [ ] Install Go 1.25+
  - [ ] Download Go binary
  - [ ] Extract to /usr/local
  - [ ] Configure PATH
  - [ ] Verify installation: `go version`
- [ ] Install PostgreSQL 18
  - [ ] Add PostgreSQL APT repository
  - [ ] Install PostgreSQL
  - [ ] Create database and user
- [ ] Install Nginx
  - [ ] Install nginx package
  - [ ] Configure as reverse proxy
  - [ ] Enable and start nginx
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
  - [ ] Database user: `campus_app`
  - [ ] Generate strong password
- [ ] Configure PostgreSQL
  - [ ] Set max_connections
  - [ ] Configure shared_buffers
  - [ ] Enable logging
- [ ] Run database migrations
  - [ ] Transfer migration files to VPS
  - [ ] Run: `go run ./cmd/migrate up`
  - [ ] Seed initial data
- [ ] Verify database
  - [ ] Test connection
  - [ ] Check tables created

### 2.3 Backend Deployment Playbook (`playbooks/deploy-backend.yml`)
- [ ] Build Go binary locally
  - [ ] `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o campus-api ./cmd/server`
- [ ] Transfer binary to VPS
  - [ ] SCP or rsync binary
  - [ ] Set executable permissions
- [ ] Configure environment
  - [ ] Create `.env` file from template
  - [ ] Set DATABASE_URL
  - [ ] Set JWT_SECRET
  - [ ] Set Google OAuth credentials
- [ ] Configure systemd service
  - [ ] Create /etc/systemd/system/campus-api.service
  - [ ] Enable service: `systemctl enable campus-api`
  - [ ] Start service: `systemctl start campus-api`
- [ ] Configure Nginx
  - [ ] Create virtual host config
  - [ ] Set up SSL with Certbot
  - [ ] Configure rate limiting
- [ ] Health check
  - [ ] Test API endpoint: `curl localhost:3000/api/health`
  - [ ] Test via domain

### 2.4 Maintenance Playbook (`playbooks/maintenance.yml`)
- [ ] Backup tasks
  - [ ] PostgreSQL dump: `pg_dump`
  - [ ] Compress and timestamp backup
  - [ ] Cleanup old backups
- [ ] Update tasks
  - [ ] Check for Go binary updates
  - [ ] Run database migrations
- [ ] Health check tasks
  - [ ] Check systemd service status
  - [ ] Check PostgreSQL status
  - [ ] Check Nginx status
  - [ ] Check disk space
- [ ] Log management
  - [ ] Configure logrotate for application logs

---

## Phase 3: Ansible Roles

### roles/common/
- [ ] tasks/main.yml - System updates, timezone, basic packages
- [ ] vars/main.yml - Common variables

### roles/go/
- [ ] tasks/main.yml - Install Go
- [ ] vars/main.yml - Go version

### roles/postgresql/
- [ ] tasks/main.yml - Install PostgreSQL
- [ ] vars/main.yml - PostgreSQL version
- [ ] templates/pg_hba.conf.j2 - PostgreSQL config

### roles/nginx/
- [ ] tasks/main.yml - Install and configure Nginx
- [ ] templates/backend.conf.j2 - Nginx virtual host

### roles/certbot/
- [ ] tasks/main.yml - Install Certbot and obtain certificates

### roles/campus-api/
- [ ] tasks/main.yml - Deploy Go binary
- [ ] templates/campus-api.service.j2 - systemd service
- [ ] templates/.env.j2 - Environment config

---

## Phase 4: Deployment Scripts

### scripts/deploy.sh
- [ ] Pull latest code
- [ ] Build Go binary
- [ ] Transfer to VPS
- [ ] Restart service
- [ ] Health check

### scripts/backup.sh
- [ ] PostgreSQL backup
- [ ] Compress with timestamp
- [ ] Optional: upload to remote storage

### scripts/rollback.sh
- [ ] Stop service
- [ ] Restore previous binary
- [ ] Start service
- [ ] Health check

---

## Phase 5: Inventory and Variables

### inventory/production.ini
```ini
[webservers]
api.stmik.tazkia.ac.id ansible_user=deploy

[database]
api.stmik.tazkia.ac.id ansible_user=deploy
```

### group_vars/all.yml
```yaml
app_name: "campus-api"
app_dir: "/var/www/campus-api"
app_port: 3000
app_user: "deploy"

go_version: "1.25.0"
postgresql_version: "18"

domain: "api.stmik.tazkia.ac.id"
ssl_email: "admin@stmik.tazkia.ac.id"
```

---

## Security Checklist

- [ ] SSH key-only authentication
- [ ] UFW firewall enabled
- [ ] Fail2ban configured
- [ ] PostgreSQL strong password (ansible-vault encrypted)
- [ ] SSL/TLS enabled
- [ ] Automatic security updates

---

## Success Criteria

- [ ] VPS can be provisioned with single command
- [ ] Backend can be deployed with single command
- [ ] Database backups run automatically
- [ ] SSL certificate auto-renews
- [ ] Deployment takes <5 minutes
- [ ] Rollback possible within 2 minutes
