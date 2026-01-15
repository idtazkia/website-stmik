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

**Step 1: Create Vault Password File**

```bash
# Create vault password file (NEVER commit this!)
echo "your-strong-vault-password" > ~/.ansible_vault_password

# Secure the file
chmod 600 ~/.ansible_vault_password
```

**Step 2: Configure Ansible to Use Vault Password**

Edit `ansible.cfg`:
```ini
[defaults]
vault_password_file = ~/.ansible_vault_password
```

**Step 3: Create Non-Sensitive Variables**

Edit `group_vars/all.yml`:

```yaml
# Application
app_name: "campus-backend"
app_dir: "/var/www/campus-backend"
app_port: 3000

# Versions
nodejs_version: "20.x"
postgresql_version: "18"

# Domain
domain: "api.stmik.tazkia.ac.id"

# Environment
env: production
```

**Step 4: Create and Encrypt Sensitive Variables**

Create `group_vars/production.yml` with **encrypted** secrets:

```bash
# Generate strong random passwords
DB_PASSWORD=$(openssl rand -base64 32)
JWT_SECRET=$(openssl rand -base64 64)

# Encrypt and add to group_vars/production.yml
ansible-vault encrypt_string "$DB_PASSWORD" --name 'database_password' >> group_vars/production.yml
ansible-vault encrypt_string "$JWT_SECRET" --name 'jwt_secret' >> group_vars/production.yml
ansible-vault encrypt_string 'your-google-client-id' --name 'google_client_id' >> group_vars/production.yml
ansible-vault encrypt_string 'your-google-client-secret' --name 'google_client_secret' >> group_vars/production.yml
```

**Example `group_vars/production.yml` (encrypted):**

```yaml
# Database Configuration
database_name: campus
database_user: campus_user
database_password: !vault |
          $ANSIBLE_VAULT;1.1;AES256
          36623262613465646437393739653234656538306461653266623938316566653932643234663331
          3265326562623539653939636233306333643035613639640a653265343435653831623339616236
          65346435656335363265626538646164633236353937353364326232343239323963396363663237
          3034613033313866660a313337613762396165386636393765656438393166653335656565313831
          6664

# Application Secrets
jwt_secret: !vault |
          $ANSIBLE_VAULT;1.1;AES256
          62303936313636323531356137396362646638656564613738346237613962353866613363613262
          3962333538643233363030353365313436316630303066310a623934323632636339393139653033
          33643834383439633831323161646537353465666231656234343332336366393738343037633931
          6364303633656130370a643466373339353030353831383264383361343734386335643335343130
          3934

# OAuth Configuration
google_client_id: !vault |
          $ANSIBLE_VAULT;1.1;AES256
          [encrypted client ID]

google_client_secret: !vault |
          $ANSIBLE_VAULT;1.1;AES256
          [encrypted client secret]
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
- Install PostgreSQL 18
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

### Ansible Vault Best Practices

**CRITICAL: All sensitive data MUST be encrypted with Ansible Vault**

1. **Create Vault Password File:**
   ```bash
   # Generate strong vault password
   openssl rand -base64 32 > ~/.ansible_vault_password
   chmod 600 ~/.ansible_vault_password
   ```

2. **Configure ansible.cfg:**
   ```ini
   [defaults]
   vault_password_file = ~/.ansible_vault_password
   ```

3. **Never Commit:**
   - ❌ `~/.ansible_vault_password` - NEVER commit vault password
   - ❌ Unencrypted secrets in `group_vars/`
   - ✅ Encrypted vault strings in `group_vars/production.yml`

4. **Encrypt All Secrets:**
   ```bash
   # Database password
   ansible-vault encrypt_string "$(openssl rand -base64 32)" --name 'database_password'

   # JWT secret
   ansible-vault encrypt_string "$(openssl rand -base64 64)" --name 'jwt_secret'

   # API keys
   ansible-vault encrypt_string 'your-api-key' --name 'api_key'
   ```

5. **View Encrypted Data:**
   ```bash
   # View encrypted file
   ansible-vault view group_vars/production.yml

   # Edit encrypted file
   ansible-vault edit group_vars/production.yml
   ```

6. **Rotate Vault Password:**
   ```bash
   # Rekey with new password
   ansible-vault rekey group_vars/production.yml
   ```

### Security Checklist

- ✅ SSH key-only authentication (passwords disabled)
- ✅ UFW firewall (ports 22, 80, 443 only)
- ✅ Fail2ban protecting SSH from brute-force
- ✅ **ALL secrets encrypted with ansible-vault**
- ✅ **Vault password file secured (chmod 600)**
- ✅ **Vault password NOT in version control**
- ✅ SSL/TLS with Let's Encrypt
- ✅ Regular security updates
- ✅ Non-root deployment user
- ✅ Strong random passwords (32+ characters)
- ✅ Database password encrypted
- ✅ JWT secret encrypted
- ✅ OAuth secrets encrypted

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
