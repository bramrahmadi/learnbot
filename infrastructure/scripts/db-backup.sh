#!/usr/bin/env bash
# ─── LearnBot Database Backup Script ─────────────────────────────────────────
# Usage:
#   ./db-backup.sh [--env staging|production] [--bucket <s3-bucket>]
#
# Environment variables (can also be set via .env):
#   DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
#   S3_BACKUP_BUCKET, AWS_REGION
#
# Dependencies: pg_dump, aws-cli, gzip

set -euo pipefail

# ─── Defaults ────────────────────────────────────────────────────────────────
ENVIRONMENT="${ENVIRONMENT:-production}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-learnbot}"
DB_USER="${DB_USER:-learnbot_admin}"
S3_BACKUP_BUCKET="${S3_BACKUP_BUCKET:-learnbot-${ENVIRONMENT}-backups}"
AWS_REGION="${AWS_REGION:-us-east-1}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"
BACKUP_DIR="${BACKUP_DIR:-/tmp/learnbot-backups}"

# ─── Parse Arguments ─────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case $1 in
    --env)       ENVIRONMENT="$2"; shift 2 ;;
    --bucket)    S3_BACKUP_BUCKET="$2"; shift 2 ;;
    --retention) RETENTION_DAYS="$2"; shift 2 ;;
    *) echo "Unknown argument: $1"; exit 1 ;;
  esac
done

# ─── Logging ─────────────────────────────────────────────────────────────────
log() { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [INFO] $*"; }
error() { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [ERROR] $*" >&2; }

# ─── Validate Dependencies ───────────────────────────────────────────────────
for cmd in pg_dump gzip aws; do
  if ! command -v "$cmd" &>/dev/null; then
    error "Required command not found: $cmd"
    exit 1
  fi
done

# ─── Create Backup ───────────────────────────────────────────────────────────
TIMESTAMP=$(date -u +%Y%m%d_%H%M%S)
BACKUP_FILENAME="${DB_NAME}_${ENVIRONMENT}_${TIMESTAMP}.sql.gz"
BACKUP_PATH="${BACKUP_DIR}/${BACKUP_FILENAME}"

mkdir -p "$BACKUP_DIR"

log "Starting backup of database '${DB_NAME}' on ${DB_HOST}:${DB_PORT}"
log "Backup file: ${BACKUP_PATH}"

PGPASSWORD="${DB_PASSWORD}" pg_dump \
  --host="$DB_HOST" \
  --port="$DB_PORT" \
  --username="$DB_USER" \
  --dbname="$DB_NAME" \
  --format=plain \
  --no-owner \
  --no-acl \
  --verbose \
  2>/tmp/pg_dump_stderr.log \
  | gzip -9 > "$BACKUP_PATH"

BACKUP_SIZE=$(du -sh "$BACKUP_PATH" | cut -f1)
log "Backup completed. Size: ${BACKUP_SIZE}"

# ─── Upload to S3 ────────────────────────────────────────────────────────────
S3_KEY="backups/${ENVIRONMENT}/${TIMESTAMP:0:8}/${BACKUP_FILENAME}"
S3_URI="s3://${S3_BACKUP_BUCKET}/${S3_KEY}"

log "Uploading backup to ${S3_URI}"

aws s3 cp "$BACKUP_PATH" "$S3_URI" \
  --region "$AWS_REGION" \
  --sse aws:kms \
  --storage-class STANDARD_IA \
  --metadata "environment=${ENVIRONMENT},db=${DB_NAME},timestamp=${TIMESTAMP}"

log "Upload completed: ${S3_URI}"

# ─── Verify Upload ───────────────────────────────────────────────────────────
REMOTE_SIZE=$(aws s3 ls "$S3_URI" --region "$AWS_REGION" | awk '{print $3}')
LOCAL_SIZE=$(stat -c%s "$BACKUP_PATH")

if [ "$REMOTE_SIZE" != "$LOCAL_SIZE" ]; then
  error "Size mismatch! Local: ${LOCAL_SIZE}, Remote: ${REMOTE_SIZE}"
  exit 1
fi
log "Backup verified. Remote size: ${REMOTE_SIZE} bytes"

# ─── Clean Up Local File ─────────────────────────────────────────────────────
rm -f "$BACKUP_PATH"
log "Local backup file removed"

# ─── Prune Old Backups ───────────────────────────────────────────────────────
log "Pruning backups older than ${RETENTION_DAYS} days from S3..."

CUTOFF_DATE=$(date -u -d "${RETENTION_DAYS} days ago" +%Y-%m-%d 2>/dev/null || \
              date -u -v-${RETENTION_DAYS}d +%Y-%m-%d)

aws s3 ls "s3://${S3_BACKUP_BUCKET}/backups/${ENVIRONMENT}/" \
  --region "$AWS_REGION" \
  --recursive \
  | awk '{print $4}' \
  | while read -r key; do
      FILE_DATE=$(echo "$key" | grep -oP '\d{8}' | head -1 | sed 's/\(....\)\(..\)\(..\)/\1-\2-\3/')
      if [[ "$FILE_DATE" < "$CUTOFF_DATE" ]]; then
        log "Deleting old backup: s3://${S3_BACKUP_BUCKET}/${key}"
        aws s3 rm "s3://${S3_BACKUP_BUCKET}/${key}" --region "$AWS_REGION"
      fi
    done

log "Backup process completed successfully"

# ─── Send Metrics to CloudWatch ──────────────────────────────────────────────
if command -v aws &>/dev/null; then
  aws cloudwatch put-metric-data \
    --namespace "LearnBot/Database" \
    --metric-name "BackupSuccess" \
    --value 1 \
    --unit Count \
    --dimensions Environment="${ENVIRONMENT}" \
    --region "$AWS_REGION" 2>/dev/null || true
fi
