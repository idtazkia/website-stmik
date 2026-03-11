#!/usr/bin/env bash
# Deploy backend to spmb.stmik.tazkia.ac.id
# Run locally: ./deploy/deploy.sh
set -euo pipefail

SSH_HOST="endymuhardin@spmb.stmik.tazkia.ac.id"
APP_NAME="stmik-admission"
APP_DIR="/opt/${APP_NAME}"
DEPLOY_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$(cd "${DEPLOY_DIR}/../backend" && pwd)"
ENV_FILE="${DEPLOY_DIR}/.env.production"

# Verify .env.production exists
if [ ! -f "${ENV_FILE}" ]; then
    echo "ERROR: ${ENV_FILE} not found"
    exit 1
fi

echo "=== Deploying ${APP_NAME} ==="

# Build
echo "--- Building Go binary (linux/amd64) ---"
cd "${BACKEND_DIR}"

# Generate templ files
echo "Generating templ templates..."
~/go/bin/templ generate

# Build CSS
echo "Building CSS..."
npm run css:build

# Cross-compile
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS="-X github.com/idtazkia/stmik-admission-api/internal/version.GitCommit=${GIT_COMMIT}"
LDFLAGS="${LDFLAGS} -X github.com/idtazkia/stmik-admission-api/internal/version.GitBranch=${GIT_BRANCH}"
LDFLAGS="${LDFLAGS} -X github.com/idtazkia/stmik-admission-api/internal/version.BuildTime=${BUILD_TIME}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o bin/server-linux ./cmd/server
echo "Binary built: bin/server-linux"

# Deploy
echo "--- Uploading .env ---"
scp "${ENV_FILE}" "${SSH_HOST}:${APP_DIR}/.env"
ssh "${SSH_HOST}" "chmod 600 ${APP_DIR}/.env"

echo "--- Uploading files ---"
rsync -avz --delete \
    bin/server-linux \
    "${SSH_HOST}:${APP_DIR}/bin/server.new"

rsync -avz --delete \
    web/static/ \
    "${SSH_HOST}:${APP_DIR}/web/static/"

rsync -avz --delete \
    web/templates/email/ \
    "${SSH_HOST}:${APP_DIR}/web/templates/email/"

rsync -avz --delete \
    migrations/ \
    "${SSH_HOST}:${APP_DIR}/migrations/"

echo "--- Running migrations ---"
ssh "${SSH_HOST}" bash -s "${APP_DIR}" << 'REMOTE_SCRIPT'
set -euo pipefail
APP_DIR="$1"

# Source env for database credentials
set -a
source "${APP_DIR}/.env"
set +a

DB_URL="postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=${DATABASE_SSL_MODE}"

# Run up migrations in order
for f in "${APP_DIR}"/migrations/*.up.sql; do
    echo "Applying: $(basename "$f")"
    psql "${DB_URL}" -f "$f" 2>&1 || true
done
echo "Migrations complete"
REMOTE_SCRIPT

echo "--- Restarting service ---"
ssh "${SSH_HOST}" bash -s "${APP_DIR}" "${APP_NAME}" << 'REMOTE_SCRIPT'
set -euo pipefail
APP_DIR="$1"
APP_NAME="$2"

# Swap binary
if [ -f "${APP_DIR}/bin/server" ]; then
    mv "${APP_DIR}/bin/server" "${APP_DIR}/bin/server.old"
fi
mv "${APP_DIR}/bin/server.new" "${APP_DIR}/bin/server"
chmod +x "${APP_DIR}/bin/server"

# Restart
sudo systemctl restart "${APP_NAME}"
sleep 2

# Health check
if curl -sf http://localhost:10002/health > /dev/null; then
    echo "Health check passed"
    rm -f "${APP_DIR}/bin/server.old"
else
    echo "Health check FAILED — rolling back"
    if [ -f "${APP_DIR}/bin/server.old" ]; then
        mv "${APP_DIR}/bin/server.old" "${APP_DIR}/bin/server"
        sudo systemctl restart "${APP_NAME}"
    fi
    exit 1
fi
REMOTE_SCRIPT

echo ""
echo "=== Deploy complete ==="
echo "https://spmb.stmik.tazkia.ac.id/health"
