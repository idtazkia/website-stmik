#!/usr/bin/env bash
# One-time server setup for spmb.stmik.tazkia.ac.id
# Run locally: ./deploy/setup-server.sh
set -euo pipefail

SSH_HOST="endymuhardin@spmb.stmik.tazkia.ac.id"
APP_NAME="stmik-admission"
APP_PORT="10002"
APP_DIR="/opt/${APP_NAME}"
DB_NAME="stmik_admission"
DB_USER="stmik_admission"

echo "=== Setting up ${APP_NAME} on spmb.stmik.tazkia.ac.id ==="

# Generate a random database password
DB_PASSWORD=$(openssl rand -base64 24 | tr -d '/+=' | head -c 24)
echo "Generated DB password: ${DB_PASSWORD}"
echo "Save this — you'll need it for the .env file."

ssh "${SSH_HOST}" bash -s "${APP_NAME}" "${APP_PORT}" "${APP_DIR}" "${DB_NAME}" "${DB_USER}" "${DB_PASSWORD}" << 'REMOTE_SCRIPT'
set -euo pipefail

APP_NAME="$1"
APP_PORT="$2"
APP_DIR="$3"
DB_NAME="$4"
DB_USER="$5"
DB_PASSWORD="$6"

echo "--- Creating application directory ---"
sudo mkdir -p "${APP_DIR}"/{bin,web/static,web/templates/email,migrations,uploads}
sudo chown -R endymuhardin:endymuhardin "${APP_DIR}"

echo "--- Creating PostgreSQL database and user ---"
# Check if user exists
if sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='${DB_USER}'" | grep -q 1; then
    echo "PostgreSQL user '${DB_USER}' already exists, updating password"
    sudo -u postgres psql -c "ALTER USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';"
else
    sudo -u postgres psql -c "CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';"
fi

# Check if database exists
if sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" | grep -q 1; then
    echo "Database '${DB_NAME}' already exists"
else
    sudo -u postgres createdb -O "${DB_USER}" "${DB_NAME}"
    echo "Database '${DB_NAME}' created"
fi

echo "--- Creating systemd service ---"
sudo tee /etc/systemd/system/${APP_NAME}.service > /dev/null << EOF
[Unit]
Description=STMIK Admission API
After=network.target postgresql@18-main.service
Requires=postgresql@18-main.service

[Service]
Type=simple
User=endymuhardin
Group=endymuhardin
WorkingDirectory=${APP_DIR}
ExecStart=${APP_DIR}/bin/server
Restart=always
RestartSec=5
EnvironmentFile=${APP_DIR}/.env

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths=${APP_DIR}/uploads
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable "${APP_NAME}"
echo "Systemd service created and enabled"

echo "--- Creating nginx site config ---"
sudo tee /etc/nginx/sites-available/spmb.stmik.tazkia.ac.id > /dev/null << 'EOF'
# Rate limiting
limit_req_zone $binary_remote_addr zone=spmb_api:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=spmb_auth:10m rate=1r/s;

server {
    listen 80;
    listen [::]:80;
    server_name spmb.stmik.tazkia.ac.id;

    # Max upload size
    client_max_body_size 10M;

    location / {
        limit_req zone=spmb_api burst=20 nodelay;

        proxy_pass http://127.0.0.1:10002;
        proxy_http_version 1.1;
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Port $server_port;
    }

    location /auth/ {
        limit_req zone=spmb_auth burst=5 nodelay;

        proxy_pass http://127.0.0.1:10002;
        proxy_http_version 1.1;
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Port $server_port;
    }

    # Deny access to test endpoints in production
    location /test/ {
        return 404;
    }
}
EOF

# Enable site
sudo ln -sf /etc/nginx/sites-available/spmb.stmik.tazkia.ac.id /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
echo "Nginx configured and reloaded"

echo "--- Obtaining SSL certificate ---"
sudo certbot --nginx -d spmb.stmik.tazkia.ac.id --non-interactive --agree-tos --email admin@tazkia.ac.id --redirect
echo "SSL certificate obtained"

echo ""
echo "=== Server setup complete ==="
echo ""
echo "Create the .env file at ${APP_DIR}/.env with:"
echo ""
echo "SERVER_PORT=${APP_PORT}"
echo "SERVER_HOST=0.0.0.0"
echo "SECURE_COOKIE=true"
echo "DATABASE_HOST=localhost"
echo "DATABASE_PORT=5432"
echo "DATABASE_USER=${DB_USER}"
echo "DATABASE_PASSWORD=${DB_PASSWORD}"
echo "DATABASE_NAME=${DB_NAME}"
echo "DATABASE_SSL_MODE=disable"
echo "JWT_SECRET=<generate with: openssl rand -base64 48>"
echo "JWT_EXPIRATION_HOURS=168"
echo "GOOGLE_CLIENT_ID=<your-google-client-id>"
echo "GOOGLE_CLIENT_SECRET=<your-google-client-secret>"
echo "GOOGLE_REDIRECT_URL=https://spmb.stmik.tazkia.ac.id/admin/auth/google/callback"
echo "STAFF_EMAIL_DOMAIN=tazkia.ac.id"
echo "UPLOAD_DIR=${APP_DIR}/uploads"
echo "MAX_UPLOAD_SIZE_MB=5"
echo "ENCRYPTION_KEY=<generate with: openssl rand -hex 32>"
REMOTE_SCRIPT

echo ""
echo "=== Next steps ==="
echo "1. SSH into the server and create ${APP_DIR}/.env with the values printed above"
echo "2. Run: ./deploy/deploy.sh  (to build and deploy the app)"
