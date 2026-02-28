#!/usr/bin/env bash
# ─── LearnBot Database Restore Script ────────────────────────────────────────
# Usage:
#   ./db-restore.sh --backup-key <s3-key> [--env staging|production]
#   ./db-restore.sh --list                  # list available backups
#
# Environment variables:
#   DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
#   S3_BACKUP_BUCKET, AWS_REGION

set -euo pipefail

ENVIRONMENT="${ENVIRONMENT:-production}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-learnbot}"
DB_USER="${DB_USER:-learnbot_admin}"
S3_BACKUP_BUCKET="${S3_BACKUP_BUCKET:-learnbot-${ENVIRONMENT}-backups}"
AWS_REGION="${AWS_REGION:-us-east-1}"
BACKUP_KEY=""
LIST_ONLY=false
RESTORE_DIR="/tmp/learnbot-restore"

log()   { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [INFO] $*"; }
error() { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [ERROR] $*" >&2; }
warn()  { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [WARN] $*"; }

while [[ $# -gt 0 ]]; do
  case $1 in
    --backup-key) BACKUP_KEY="$2"; shift 2 ;;
    --env)        ENVIRONMENT="$2"; shift 2 ;;
    --list)       LIST_ONLY=true; shift ;;
    *) echo "Unknown argument: $1"; exit 1 ;;
  esac
done

# ─── List Available Backups ──────────────────────────────────────────────────
if [ "$LIST_ONLY" = true ]; then
  log "Available backups for environment: ${ENVIRONMENT}"
  aws s3 ls "s3://${S3_BACKUP_BUCKET}/backups/${ENVIRONMENT}/" \
    --region "$AWS_REGION" \
    --recursive \
    | sort -r \
    | head -20
  exit 0
fi

if [ -z "$BACKUP_KEY" ]; then
  error "No backup key specified. Use --backup-key <s3-key> or --list to see available backups."
  exit 1
fi

# ─── Safety Confirmation ─────────────────────────────────────────────────────
warn "⚠️  WARNING: This will OVERWRITE the database '${DB_NAME}' on ${DB_HOST}"
warn "⚠️  Backup to restore: s3://${S3_BACKUP_BUCKET}/${BACKUP_KEY}"
echo ""
read -r -p "Type 'RESTORE' to confirm: " CONFIRM
if [ "$CONFIRM" != "RESTORE" ]; then
  log "Restore cancelled."
  exit 0
fi

# ─── Download Backup ─────────────────────────────────────────────────────────
mkdir -p "$RESTORE_DIR"
LOCAL_FILE="${RESTORE_DIR}/$(basename "$BACKUP_KEY")"

log "Downloading backup from s3://${S3_BACKUP_BUCKET}/${BACKUP_KEY}"
aws s3 cp "s3://${S3_BACKUP_BUCKET}/${BACKUP_KEY}" "$LOCAL_FILE" \
  --region "$AWS_REGION"

log "Download completed: ${LOCAL_FILE}"

# ─── Create Pre-restore Backup ───────────────────────────────────────────────
log "Creating pre-restore backup of current database..."
PRE_RESTORE_FILE="${RESTORE_DIR}/pre_restore_$(date -u +%Y%m%d_%H%M%S).sql.gz"

PGPASSWORD="${DB_PASSWORD}" pg_dump \
  --host="$DB_HOST" \
  --port="$DB_PORT" \
  --username="$DB_USER" \
  --dbname="$DB_NAME" \
  --format=plain \
  | gzip -9 > "$PRE_RESTORE_FILE"

log "Pre-restore backup saved: ${PRE_RESTORE_FILE}"

# ─── Terminate Active Connections ────────────────────────────────────────────
log "Terminating active connections to database..."
PGPASSWORD="${DB_PASSWORD}" psql \
  --host="$DB_HOST" \
  --port="$DB_PORT" \
  --username="$DB_USER" \
  --dbname="postgres" \
  --command="SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${DB_NAME}' AND pid <> pg_backend_pid();" \
  > /dev/null

# ─── Restore Database ────────────────────────────────────────────────────────
log "Restoring database from backup..."

gunzip -c "$LOCAL_FILE" | PGPASSWORD="${DB_PASSWORD}" psql \
  --host="$DB_HOST" \
  --port="$DB_PORT" \
  --username="$DB_USER" \
  --dbname="$DB_NAME" \
  --single-transaction \
  --set ON_ERROR_STOP=on

log "Database restore completed successfully"

# ─── Verify Restore ──────────────────────────────────────────────────────────
log "Verifying restore..."
TABLE_COUNT=$(PGPASSWORD="${DB_PASSWORD}" psql \
  --host="$DB_HOST" \
  --port="$DB_PORT" \
  --username="$DB_USER" \
  --dbname="$DB_NAME" \
  --tuples-only \
  --command="SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")

log "Tables in restored database: ${TABLE_COUNT}"

# ─── Clean Up ────────────────────────────────────────────────────────────────
rm -f "$LOCAL_FILE"
log "Temporary files cleaned up"
log "Pre-restore backup retained at: ${PRE_RESTORE_FILE}"
log "Restore process completed successfully"
