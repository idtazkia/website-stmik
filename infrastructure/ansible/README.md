# Ansible Automation - STMIK Tazkia

Ansible playbooks for automated VPS provisioning and application deployment.

**Target OS:** Ubuntu 24.04 LTS (or 22.04 LTS)
**Status:** Not Started (Phase 2)

---

## Directory Structure

```
ansible/
├── playbooks/              # Ansible playbooks
│   ├── setup-vps.yml      # Initial VPS setup (Node.js, PostgreSQL, Nginx, etc.)
│   ├── deploy-backend.yml # Deploy Express.js backend
│   ├── setup-database.yml # PostgreSQL configuration and migrations
│   └── maintenance.yml    # Backups, updates, health checks
├── inventory/             # Server inventory
│   ├── production.ini     # Production VPS
│   └── staging.ini        # Staging VPS (optional)
├── roles/                 # Reusable Ansible roles
│   ├── nodejs/           # Node.js installation
│   ├── postgresql/       # PostgreSQL setup
│   ├── nginx/            # Nginx reverse proxy
│   ├── pm2/              # PM2 process manager
│   └── ssl-certbot/      # SSL certificate management
├── group_vars/           # Variable files
│   ├── all.yml          # Global variables
│   └── production.yml   # Production-specific variables
├── ansible.cfg           # Ansible configuration
└── README.md             # This file
```

---

## Prerequisites

### 1. Install Ansible (Local Machine)

**macOS:**
```bash
brew install ansible
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ansible
```

**Verify installation:**
```bash
ansible --version  # Should be 2.14 or higher
```

### 2. Install Ansible Collections

```bash
ansible-galaxy collection install community.general
ansible-galaxy collection install community.postgresql
```

### 3. VPS Requirements

- Ubuntu 24.04 LTS (or 22.04 LTS)
- 2GB RAM minimum
- 20GB storage minimum
- Public IP address
- Root or sudo access

### 4. SSH Access

Generate SSH key (if not already done):
```bash
ssh-keygen -t ed25519 -C "deployment@stmik-tazkia"
```

Copy SSH key to VPS:
```bash
ssh-copy-id -i ~/.ssh/id_ed25519.pub root@YOUR_VPS_IP
```

Test SSH connection:
```bash
ssh root@YOUR_VPS_IP
```

---

## Configuration

### 1. Update Inventory File

Edit `inventory/production.ini`:

```ini
[webservers]
backend-vps ansible_host=YOUR_VPS_IP ansible_user=root ansible_port=22

[databases]
backend-vps ansible_host=YOUR_VPS_IP ansible_user=root ansible_port=22
```

Replace `YOUR_VPS_IP` with your actual VPS IP address.

### 2. Configure Variables

Edit `group_vars/all.yml`:

```yaml
# Application
app_name: "campus-backend"
app_dir: "/var/www/campus-backend"
app_port: 3000

# Versions
nodejs_version: "20.x"
postgresql_version: "16"

# Domain
domain: "api.stmik.tazkia.ac.id"
```

Edit `group_vars/production.yml` and encrypt sensitive data:

```yaml
# Database
database_name: campus
database_user: campus_user
database_password: !vault |
  $ANSIBLE_VAULT;1.1;AES256
  [encrypted password here]

# Application
jwt_secret: !vault |
  $ANSIBLE_VAULT;1.1;AES256
  [encrypted JWT secret here]
```

To encrypt secrets:
```bash
ansible-vault encrypt_string 'your-database-password' --name 'database_password'
ansible-vault encrypt_string 'your-jwt-secret' --name 'jwt_secret'
```

---

## Usage

### Initial VPS Setup

Provision a fresh VPS with all required software:

```bash
ansible-playbook -i inventory/production.ini playbooks/setup-vps.yml
```

This will:
- Update system packages
- Create deployment user
- Install Node.js 20.x
- Install PostgreSQL 16
- Install Nginx
- Install PM2
- Install Certbot (SSL)
- Configure firewall (UFW)
- Install Fail2ban

**Duration:** ~10-15 minutes

### Deploy Backend Application

Deploy the Express.js backend:

```bash
ansible-playbook -i inventory/production.ini playbooks/deploy-backend.yml
```

This will:
- Clone/update repository
- Install dependencies
- Configure environment variables
- Run database migrations
- Start application with PM2
- Configure Nginx reverse proxy
- Obtain SSL certificate

**Duration:** ~5-10 minutes

### Database Setup

Initialize database and run migrations:

```bash
ansible-playbook -i inventory/production.ini playbooks/setup-database.yml
```

### Maintenance

Run backups, updates, and health checks:

```bash
ansible-playbook -i inventory/production.ini playbooks/maintenance.yml
```

This will:
- Backup PostgreSQL database
- Update system packages
- Rotate logs
- Check service health
- Renew SSL certificates if needed

---

## Playbook Details

### setup-vps.yml

Initial VPS provisioning playbook.

**Tasks:**
- System updates and security patches
- User management (create deploy user, disable root login)
- Software installation (Node.js, PostgreSQL, Nginx, PM2, Certbot)
- Firewall configuration (UFW)
- Fail2ban setup for SSH protection

**Usage:**
```bash
ansible-playbook -i inventory/production.ini playbooks/setup-vps.yml
```

**First run only:** Use root user, subsequent runs use deploy user.

---

### deploy-backend.yml

Backend application deployment playbook.

**Tasks:**
- Clone/update Git repository
- Install npm dependencies
- Create/update .env file
- Run database migrations
- Start/restart PM2 application
- Configure Nginx reverse proxy
- Obtain/renew SSL certificate

**Usage:**
```bash
ansible-playbook -i inventory/production.ini playbooks/deploy-backend.yml
```

**Safe to run multiple times:** Idempotent, won't break running application.

---

### setup-database.yml

Database initialization playbook.

**Tasks:**
- Create PostgreSQL database
- Create database user
- Configure PostgreSQL settings
- Run migrations
- Verify database health

**Usage:**
```bash
ansible-playbook -i inventory/production.ini playbooks/setup-database.yml
```

---

### maintenance.yml

Automated maintenance tasks.

**Tasks:**
- Database backup (pg_dump)
- System package updates
- Log rotation
- Health checks (PM2, Nginx, PostgreSQL, disk space)
- SSL certificate renewal

**Usage:**
```bash
ansible-playbook -i inventory/production.ini playbooks/maintenance.yml
```

**Recommended:** Run daily via cron job.

---

## Testing

### Dry Run (Check Mode)

Test playbook without making changes:

```bash
ansible-playbook -i inventory/production.ini playbooks/setup-vps.yml --check
```

### Verbose Output

Run with detailed output for debugging:

```bash
ansible-playbook -i inventory/production.ini playbooks/deploy-backend.yml -vvv
```

### Specific Tags

Run only specific tasks:

```bash
ansible-playbook -i inventory/production.ini playbooks/maintenance.yml --tags backup
```

---

## Troubleshooting

### SSH Connection Failed

**Error:** `Failed to connect to the host via ssh`

**Solution:**
1. Verify VPS IP in inventory file
2. Test SSH manually: `ssh user@vps-ip`
3. Check SSH key: `ssh-add -l`
4. Verify `~/.ssh/config` if using custom SSH config

### Permission Denied

**Error:** `Permission denied (publickey)`

**Solution:**
1. Copy SSH key to VPS: `ssh-copy-id user@vps-ip`
2. Verify user has sudo access
3. Check `ansible.cfg` for correct `remote_user`

### Playbook Failed

**Solution:**
1. Run with verbose output: `-vvv`
2. Check task that failed
3. SSH to VPS and check logs: `sudo journalctl -xe`
4. Fix issue and re-run playbook (idempotent)

---

## Best Practices

1. **Always use ansible-vault for secrets:**
   ```bash
   ansible-vault encrypt_string 'secret-value' --name 'variable_name'
   ```

2. **Test playbooks in staging first:**
   ```bash
   ansible-playbook -i inventory/staging.ini playbooks/deploy-backend.yml
   ```

3. **Run maintenance regularly:**
   - Setup cron job for daily maintenance playbook
   - Monitor disk space and logs

4. **Keep playbooks idempotent:**
   - Safe to run multiple times
   - Use `changed_when` and `failed_when` conditions

5. **Version control all changes:**
   - Commit playbook changes to Git
   - Document changes in commit messages

---

## Security

- ✅ SSH key-only authentication (passwords disabled)
- ✅ UFW firewall (ports 22, 80, 443 only)
- ✅ Fail2ban protecting SSH from brute-force
- ✅ Secrets encrypted with ansible-vault
- ✅ SSL/TLS with Let's Encrypt
- ✅ Regular security updates
- ✅ Non-root deployment user

---

## Future Improvements

- [ ] Multi-server load balancing
- [ ] Database replication (primary + replica)
- [ ] Monitoring (Prometheus + Grafana)
- [ ] Log aggregation (ELK stack)
- [ ] Blue-green deployment
- [ ] Automated rollback on failure

---

## Documentation

- **Infrastructure Overview:** `../README.md`
- **TODO List:** `../TODO.md`
- **Deployment Guide:** `../../docs/DEPLOYMENT.md`
- **Architecture:** `../../docs/ARCHITECTURE.md`
