# LearnBot Operations Runbook

## Table of Contents

1. [Service Overview](#service-overview)
2. [Deployment Procedures](#deployment-procedures)
3. [Rollback Procedures](#rollback-procedures)
4. [Database Operations](#database-operations)
5. [Incident Response](#incident-response)
6. [Scaling Operations](#scaling-operations)
7. [Monitoring & Alerting](#monitoring--alerting)
8. [Security Operations](#security-operations)
9. [Common Issues & Resolutions](#common-issues--resolutions)

---

## Service Overview

| Service | Port | Description | Health Endpoint |
|---------|------|-------------|-----------------|
| api-gateway | 8090 | Main API gateway, JWT auth, rate limiting | `/health` |
| resume-parser | 8080 | PDF/DOCX resume parsing and analysis | `/health` |
| job-aggregator | 8081 | Job scraping and aggregation | `/health` |
| learning-resources | 8082 | Learning resource management | `/health` |
| frontend | 3000 | Next.js frontend application | `/` |
| PostgreSQL | 5432 | Primary database | `pg_isready` |

### Architecture

```
Internet → ALB → api-gateway → resume-parser
                             → job-aggregator
                             → learning-resources
                             → PostgreSQL (RDS)
                             → S3 (resume storage)
```

---

## Deployment Procedures

### Staging Deployment (Automatic)

Staging deploys automatically on every push to `main` via GitHub Actions.

**Monitor at:** GitHub Actions → Deploy to Staging workflow

### Production Deployment (Manual)

1. **Verify staging is healthy:**
   ```bash
   curl https://staging-api.learnbot.example.com/health
   ```

2. **Get the image tag to deploy** (from successful staging deployment):
   ```bash
   # Get the latest staging deployment SHA
   git log --oneline -5
   ```

3. **Trigger production deployment** via GitHub Actions:
   - Go to Actions → "Deploy to Production"
   - Click "Run workflow"
   - Enter the full 40-character SHA as `image_tag`
   - Type `deploy` in the confirmation field

4. **Monitor deployment:**
   ```bash
   # Watch ECS service events
   aws ecs describe-services \
     --cluster learnbot-production-cluster \
     --services learnbot-production-api-gateway \
     --query 'services[0].events[:5]'
   ```

5. **Verify production health:**
   ```bash
   curl https://api.learnbot.example.com/health
   ```

### Manual Deployment via Script

```bash
# Deploy all services to staging
export AWS_ACCOUNT_ID=123456789012
./infrastructure/scripts/deploy.sh --env staging --tag <sha>

# Deploy single service
./infrastructure/scripts/deploy.sh --env staging --tag <sha> --service api-gateway

# Dry run
./infrastructure/scripts/deploy.sh --env production --tag <sha> --dry-run
```

---

## Rollback Procedures

### Automatic Rollback

ECS deployment circuit breaker is enabled. If a deployment fails health checks, ECS automatically rolls back to the previous task definition.

### Manual Rollback via GitHub Actions

1. Go to Actions → "Deploy to Production"
2. Use the previous known-good SHA as `image_tag`

### Manual Rollback via CLI

```bash
# Get previous task definition revision
aws ecs describe-services \
  --cluster learnbot-production-cluster \
  --services learnbot-production-api-gateway \
  --query 'services[0].deployments'

# Roll back to previous revision
FAMILY="learnbot-production-api-gateway"
PREV_REV=42  # Replace with actual previous revision

aws ecs update-service \
  --cluster learnbot-production-cluster \
  --service learnbot-production-api-gateway \
  --task-definition "${FAMILY}:${PREV_REV}" \
  --force-new-deployment

# Wait for stability
aws ecs wait services-stable \
  --cluster learnbot-production-cluster \
  --services learnbot-production-api-gateway
```

### Database Rollback

If a migration caused issues:

```bash
# List available backups
./infrastructure/scripts/db-restore.sh --list --env production

# Restore from pre-deployment snapshot
./infrastructure/scripts/db-restore.sh \
  --backup-key backups/production/20240101/learnbot_production_20240101_120000.sql.gz \
  --env production
```

---

## Database Operations

### Connect to Production Database

```bash
# Via AWS SSM Session Manager (no direct SSH needed)
aws ssm start-session \
  --target <bastion-instance-id> \
  --document-name AWS-StartPortForwardingSessionToRemoteHost \
  --parameters '{"host":["<rds-endpoint>"],"portNumber":["5432"],"localPortNumber":["5432"]}'

# Then connect locally
PGPASSWORD=$(aws secretsmanager get-secret-value \
  --secret-id learnbot-production/db/password \
  --query SecretString --output text) \
psql -h localhost -U learnbot_admin -d learnbot
```

### Run Database Migrations

```bash
# Migrations run automatically during deployment
# To run manually:
aws ecs run-task \
  --cluster learnbot-production-cluster \
  --task-definition learnbot-production-migrate \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=DISABLED}"
```

### Manual Database Backup

```bash
export DB_HOST=$(aws rds describe-db-instances \
  --db-instance-identifier learnbot-production-postgres \
  --query 'DBInstances[0].Endpoint.Address' --output text)
export DB_PASSWORD=$(aws secretsmanager get-secret-value \
  --secret-id learnbot-production/db/password \
  --query SecretString --output text)

./infrastructure/scripts/db-backup.sh --env production
```

### Check Database Health

```bash
# Connection count
psql -c "SELECT count(*) FROM pg_stat_activity WHERE datname = 'learnbot';"

# Long-running queries
psql -c "SELECT pid, now() - pg_stat_activity.query_start AS duration, query
         FROM pg_stat_activity
         WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';"

# Table sizes
psql -c "SELECT schemaname, tablename,
         pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
         FROM pg_tables WHERE schemaname = 'public'
         ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

---

## Incident Response

### Severity Levels

| Level | Description | Response Time | Examples |
|-------|-------------|---------------|---------|
| P1 - Critical | Complete service outage | 15 minutes | All services down, data loss |
| P2 - High | Major feature unavailable | 1 hour | Resume parsing failing, auth broken |
| P3 - Medium | Degraded performance | 4 hours | High latency, partial failures |
| P4 - Low | Minor issue | Next business day | UI glitch, non-critical errors |

### P1 Incident Response

1. **Acknowledge** the alert in PagerDuty/Slack
2. **Assess** the scope:
   ```bash
   # Check all service health
   for svc in api-gateway resume-parser job-aggregator learning-resources; do
     echo -n "$svc: "
     curl -s -o /dev/null -w "%{http_code}" \
       https://api.learnbot.example.com/health || echo "UNREACHABLE"
   done
   ```
3. **Check recent deployments:**
   ```bash
   git log --oneline -10
   ```
4. **Check ECS service events:**
   ```bash
   aws ecs describe-services \
     --cluster learnbot-production-cluster \
     --services learnbot-production-api-gateway \
     --query 'services[0].events[:10]'
   ```
5. **Check CloudWatch logs:**
   ```bash
   aws logs tail /ecs/learnbot-production/api-gateway --follow
   ```
6. **Rollback if deployment-related** (see Rollback Procedures)
7. **Communicate** status to stakeholders
8. **Post-mortem** within 48 hours

### Checking Logs

```bash
# Tail logs for a service
aws logs tail /ecs/learnbot-production/api-gateway --follow

# Search for errors in last hour
aws logs filter-log-events \
  --log-group-name /ecs/learnbot-production/api-gateway \
  --start-time $(date -d '1 hour ago' +%s000) \
  --filter-pattern '"level":"error"'

# Get logs for specific request ID
aws logs filter-log-events \
  --log-group-name /ecs/learnbot-production/api-gateway \
  --filter-pattern '"request_id":"<REQUEST_ID>"'
```

---

## Scaling Operations

### Manual Scaling

```bash
# Scale up API gateway
aws ecs update-service \
  --cluster learnbot-production-cluster \
  --service learnbot-production-api-gateway \
  --desired-count 5

# Check current scaling
aws ecs describe-services \
  --cluster learnbot-production-cluster \
  --services learnbot-production-api-gateway \
  --query 'services[0].{desired:desiredCount,running:runningCount,pending:pendingCount}'
```

### Auto-scaling Configuration

Auto-scaling is configured via Terraform. Current thresholds:
- **Scale out:** CPU > 70% or Memory > 80% for 60 seconds
- **Scale in:** CPU < 70% and Memory < 80% for 300 seconds
- **Min replicas:** 2 (production), 1 (staging)
- **Max replicas:** 10 (production), 3 (staging)

To modify thresholds, update `infrastructure/terraform/main.tf` and apply.

---

## Monitoring & Alerting

### Dashboards

| Dashboard | URL | Description |
|-----------|-----|-------------|
| Service Overview | http://grafana:3001/d/learnbot-overview | Request rates, errors, latency |
| Database | http://grafana:3001/d/postgres | Connection pool, query performance |
| Infrastructure | http://grafana:3001/d/node | CPU, memory, disk |

### Alert Channels

- **Slack:** `#learnbot-alerts` channel
- **Email:** ops@learnbot.example.com
- **PagerDuty:** For P1/P2 incidents (on-call rotation)

### Silencing Alerts

```bash
# During planned maintenance, silence alerts via Alertmanager API
curl -X POST http://alertmanager:9093/api/v2/silences \
  -H "Content-Type: application/json" \
  -d '{
    "matchers": [{"name": "environment", "value": "production", "isRegex": false}],
    "startsAt": "2024-01-01T00:00:00Z",
    "endsAt": "2024-01-01T02:00:00Z",
    "comment": "Planned maintenance window",
    "createdBy": "ops-team"
  }'
```

---

## Security Operations

### Rotating Secrets

```bash
# Rotate JWT secret
NEW_SECRET=$(openssl rand -base64 48)
aws secretsmanager put-secret-value \
  --secret-id learnbot-production/app/jwt-secret \
  --secret-string "{\"jwt_secret\":\"${NEW_SECRET}\"}"

# Force service restart to pick up new secret
aws ecs update-service \
  --cluster learnbot-production-cluster \
  --service learnbot-production-api-gateway \
  --force-new-deployment

# Rotate DB password
NEW_DB_PASS=$(openssl rand -base64 32)
aws secretsmanager put-secret-value \
  --secret-id learnbot-production/db/password \
  --secret-string "$NEW_DB_PASS"

# Update RDS password
aws rds modify-db-instance \
  --db-instance-identifier learnbot-production-postgres \
  --master-user-password "$NEW_DB_PASS" \
  --apply-immediately
```

### Reviewing Access Logs

```bash
# Download ALB access logs from S3
aws s3 sync \
  s3://learnbot-production-alb-logs/alb/ \
  /tmp/alb-logs/ \
  --region us-east-1

# Analyze for suspicious patterns
grep -h "5[0-9][0-9]" /tmp/alb-logs/*.log | \
  awk '{print $3}' | sort | uniq -c | sort -rn | head -20
```

---

## Common Issues & Resolutions

### Issue: Service fails to start after deployment

**Symptoms:** ECS tasks keep stopping, deployment circuit breaker triggers

**Resolution:**
1. Check task logs: `aws logs tail /ecs/learnbot-production/<service> --follow`
2. Common causes:
   - Missing environment variable → Check Secrets Manager and task definition
   - Database connection failure → Verify security groups and DB credentials
   - Port conflict → Check task definition port mappings
3. Roll back if needed (see Rollback Procedures)

### Issue: High database connection count

**Symptoms:** `PostgreSQLHighConnections` alert fires

**Resolution:**
1. Check current connections:
   ```sql
   SELECT count(*), state FROM pg_stat_activity GROUP BY state;
   ```
2. Kill idle connections:
   ```sql
   SELECT pg_terminate_backend(pid)
   FROM pg_stat_activity
   WHERE state = 'idle'
   AND query_start < now() - interval '10 minutes';
   ```
3. Consider enabling PgBouncer connection pooling

### Issue: Resume parsing failures

**Symptoms:** `ResumeParseFailureRate` alert fires

**Resolution:**
1. Check resume-parser logs for error patterns
2. Verify S3 bucket permissions
3. Check disk space in `/tmp/uploads`
4. Verify PDF/DOCX library dependencies in container

### Issue: Job scraping not running

**Symptoms:** `JobScrapingFailed` alert fires

**Resolution:**
1. Check job-aggregator logs
2. Verify external job board URLs are accessible
3. Check if IP is blocked by job boards (rotate if needed)
4. Manually trigger scrape via admin endpoint:
   ```bash
   curl -X POST https://api.learnbot.example.com/admin/scrape/trigger \
     -H "Authorization: Bearer <admin-token>"
   ```

### Issue: High API latency

**Symptoms:** `HighP99Latency` alert fires

**Resolution:**
1. Check if specific endpoints are slow:
   ```bash
   aws logs filter-log-events \
     --log-group-name /ecs/learnbot-production/api-gateway \
     --filter-pattern '"duration_ms" > 2000'
   ```
2. Check database query performance
3. Check if auto-scaling is keeping up with load
4. Consider adding caching for frequently accessed data
