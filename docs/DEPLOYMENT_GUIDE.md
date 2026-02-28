# LearnBot — Deployment Guide

## Table of Contents

1. [Infrastructure Overview](#1-infrastructure-overview)
2. [Prerequisites](#2-prerequisites)
3. [Local Development Deployment](#3-local-development-deployment)
4. [Staging Deployment](#4-staging-deployment)
5. [Production Deployment](#5-production-deployment)
6. [Configuration Management](#6-configuration-management)
7. [Database Management](#7-database-management)
8. [CI/CD Pipeline](#8-cicd-pipeline)
9. [Security](#9-security)
10. [Monitoring & Alerting](#10-monitoring--alerting)
11. [Rollback Procedures](#11-rollback-procedures)
12. [Troubleshooting Guide](#12-troubleshooting-guide)

---

## 1. Infrastructure Overview

LearnBot runs on AWS using a containerized microservices architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                          AWS Cloud                               │
│                                                                  │
│  ┌──────────┐    ┌──────────────────────────────────────────┐   │
│  │CloudFront│    │              VPC (10.0.0.0/16)           │   │
│  │   CDN    │    │                                          │   │
│  └────┬─────┘    │  ┌─────────────────────────────────┐    │   │
│       │          │  │     Public Subnets               │    │   │
│       │          │  │  ┌──────────────────────────┐   │    │   │
│       │          │  │  │  Application Load Balancer│   │    │   │
│       │          │  │  │  (HTTPS/TLS termination) │   │    │   │
│       │          │  │  └──────────┬───────────────┘   │    │   │
│       │          │  └─────────────┼───────────────────┘    │   │
│       │          │                │                         │   │
│       │          │  ┌─────────────┼───────────────────┐    │   │
│       │          │  │     Private Subnets              │    │   │
│       │          │  │             │                    │    │   │
│       │          │  │  ┌──────────▼──────────────┐    │    │   │
│       │          │  │  │     ECS Fargate Cluster  │    │    │   │
│       │          │  │  │  ┌─────────────────────┐│    │    │   │
│       │          │  │  │  │   api-gateway :8090  ││    │    │   │
│       │          │  │  │  │   resume-parser:8080 ││    │    │   │
│       │          │  │  │  │   job-aggregator:8081││    │    │   │
│       │          │  │  │  │   learning-res.:8082 ││    │    │   │
│       │          │  │  │  └─────────────────────┘│    │    │   │
│       │          │  │  └──────────────────────────┘    │    │   │
│       │          │  │                                   │    │   │
│       │          │  │  ┌────────────────────────────┐  │    │   │
│       │          │  │  │  RDS PostgreSQL (Multi-AZ) │  │    │   │
│       │          │  │  └────────────────────────────┘  │    │   │
│       │          │  └───────────────────────────────────┘    │   │
│       │          └──────────────────────────────────────────┘   │
│       │                                                          │
│  ┌────▼─────┐    ┌──────────────────────────────────────────┐   │
│  │  S3      │    │  S3 Buckets                              │   │
│  │ (static) │    │  • learnbot-{env}-resumes (encrypted)    │   │
│  └──────────┘    │  • learnbot-{env}-alb-logs               │   │
│                  └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### AWS Services Used

| Service | Purpose |
|---------|---------|
| ECS Fargate | Container orchestration (serverless) |
| RDS PostgreSQL | Primary database (Multi-AZ in production) |
| ElastiCache Redis | Rate limiting and caching |
| S3 | Resume file storage, ALB access logs |
| CloudFront | Frontend CDN |
| ALB | Load balancing, TLS termination |
| ECR | Docker image registry |
| Secrets Manager | Secrets storage |
| CloudWatch | Logs and metrics |
| KMS | Encryption key management |

---

## 2. Prerequisites

### Tools Required

```bash
# Terraform >= 1.6.0
terraform --version

# AWS CLI >= 2.0
aws --version

# Docker >= 24.0
docker --version

# kubectl (for Kubernetes deployments)
kubectl version --client
```

### AWS Credentials

Configure AWS CLI with appropriate credentials:

```bash
aws configure
# AWS Access Key ID: <your-key>
# AWS Secret Access Key: <your-secret>
# Default region: us-east-1
# Default output format: json
```

Or use AWS SSO:
```bash
aws sso login --profile learnbot-deploy
```

### Required IAM Permissions

The deployment IAM role needs:
- `ecs:*` — ECS task and service management
- `ecr:*` — Docker image push/pull
- `rds:*` — Database management
- `s3:*` — Bucket and object management
- `secretsmanager:*` — Secrets management
- `iam:PassRole` — ECS task role assignment
- `elasticloadbalancing:*` — ALB management
- `cloudfront:*` — CDN management

---

## 3. Local Development Deployment

### Quick Start

```bash
# 1. Clone and configure
git clone https://github.com/learnbot/learnbot.git
cd learnbot
cp .env.example .env
# Edit .env as needed

# 2. Start all services
docker compose up -d

# 3. Verify health
curl http://localhost:8090/health
open http://localhost:3000
```

### Service Ports

| Service | Port | URL |
|---------|------|-----|
| Frontend | 3000 | http://localhost:3000 |
| API Gateway | 8090 | http://localhost:8090 |
| Resume Parser | 8080 | http://localhost:8080 |
| Job Aggregator | 8081 | http://localhost:8081 |
| Learning Resources | 8082 | http://localhost:8082 |
| PostgreSQL | 5432 | localhost:5432 |
| Redis | 6379 | localhost:6379 |
| Prometheus | 9090 | http://localhost:9090 |
| Grafana | 3001 | http://localhost:3001 |
| Nginx (proxy) | 80 | http://localhost |

### Stopping Services

```bash
# Stop all services (keep data)
docker compose down

# Stop and remove all data (clean slate)
docker compose down -v
```

---

## 4. Staging Deployment

### First-Time Terraform Setup

```bash
# 1. Create S3 bucket for Terraform state
aws s3 mb s3://learnbot-terraform-state --region us-east-1
aws s3api put-bucket-versioning \
  --bucket learnbot-terraform-state \
  --versioning-configuration Status=Enabled

# 2. Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name learnbot-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1

# 3. Initialize Terraform
cd infrastructure/terraform
terraform init
```

### Deploy Staging Infrastructure

```bash
cd infrastructure/terraform

# Select or create staging workspace
terraform workspace new staging || terraform workspace select staging

# Plan (review changes)
terraform plan -var-file=terraform.tfvars.staging -out=staging.plan

# Apply
terraform apply staging.plan
```

### Deploy Application to Staging

```bash
# Using the deployment script
./infrastructure/scripts/deploy.sh staging <git-sha>

# Or manually:
# 1. Build and push Docker images
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

docker build -t learnbot/api-gateway:staging .
docker tag learnbot/api-gateway:staging \
  <account-id>.dkr.ecr.us-east-1.amazonaws.com/learnbot-staging-api-gateway:latest
docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/learnbot-staging-api-gateway:latest

# 2. Run database migrations
aws ecs run-task \
  --cluster learnbot-staging \
  --task-definition learnbot-staging-migrate \
  --launch-type FARGATE \
  --network-configuration "..."

# 3. Update ECS service
aws ecs update-service \
  --cluster learnbot-staging \
  --service learnbot-staging-api-gateway \
  --force-new-deployment
```

---

## 5. Production Deployment

> ⚠️ **Production deployments require manual approval.** Never deploy directly to production without review.

### Pre-Deployment Checklist

- [ ] All tests pass on `main` branch
- [ ] Staging deployment is healthy
- [ ] Database backup completed
- [ ] Change request approved
- [ ] Rollback plan documented
- [ ] Team notified in `#learnbot-ops`

### Deploy to Production

```bash
cd infrastructure/terraform

# Select production workspace
terraform workspace select production

# Plan (review carefully!)
terraform plan -var-file=terraform.tfvars.production -out=production.plan

# Review the plan output thoroughly
cat production.plan

# Apply (requires confirmation)
terraform apply production.plan
```

### Application Deployment

```bash
# Using the deployment script (recommended)
./infrastructure/scripts/deploy.sh production <git-sha>
```

The deployment script:
1. Validates the git SHA exists
2. Backs up the production database
3. Builds and pushes Docker images to ECR
4. Runs database migrations
5. Deploys to ECS with rolling update
6. Runs health checks
7. Tags the release in git
8. Notifies Slack

### Post-Deployment Verification

```bash
# Check service health
curl https://api.learnbot.example.com/health

# Check ECS service status
aws ecs describe-services \
  --cluster learnbot-production \
  --services learnbot-production-api-gateway

# Check CloudWatch logs
aws logs tail /ecs/learnbot-production-api-gateway --follow
```

---

## 6. Configuration Management

### Secrets Management

All secrets are stored in **AWS Secrets Manager**. Never store secrets in:
- Environment files (`.env`)
- Docker Compose files
- Terraform state
- Git repository

#### Creating Secrets

```bash
# JWT Secret
aws secretsmanager create-secret \
  --name "learnbot/production/jwt-secret" \
  --secret-string "$(openssl rand -base64 48)"

# Database password
aws secretsmanager create-secret \
  --name "learnbot/production/db-password" \
  --secret-string "$(openssl rand -base64 32)"
```

#### Rotating Secrets

```bash
# Enable automatic rotation (90-day cycle)
aws secretsmanager rotate-secret \
  --secret-id "learnbot/production/jwt-secret" \
  --rotation-rules AutomaticallyAfterDays=90
```

### Environment-Specific Configuration

| Config | Staging | Production |
|--------|---------|------------|
| RDS instance | `db.t3.medium` | `db.r6g.large` |
| ECS CPU | 256 | 512 |
| ECS Memory | 512 MB | 1024 MB |
| Min tasks | 1 | 2 |
| Max tasks | 3 | 10 |
| RDS Multi-AZ | No | Yes |
| RDS backup retention | 3 days | 14 days |

### Kubernetes Configuration (Alternative)

For Kubernetes deployments, manifests are in `infrastructure/kubernetes/`:

```bash
# Apply namespace
kubectl apply -f infrastructure/kubernetes/namespace.yaml

# Apply ConfigMap
kubectl apply -f infrastructure/kubernetes/configmap.yaml

# Deploy API Gateway
kubectl apply -f infrastructure/kubernetes/api-gateway/

# Apply Ingress
kubectl apply -f infrastructure/kubernetes/ingress.yaml
```

---

## 7. Database Management

### Running Migrations

```bash
# Staging
./infrastructure/scripts/deploy.sh staging <sha> --migrate-only

# Production (always backup first!)
./infrastructure/scripts/db-backup.sh production
./infrastructure/scripts/deploy.sh production <sha> --migrate-only
```

### Database Backup

```bash
# Manual backup
./infrastructure/scripts/db-backup.sh production

# Backup is stored in S3:
# s3://learnbot-production-backups/postgres/YYYY-MM-DD-HH-MM-SS.sql.gz
```

### Database Restore

```bash
# List available backups
aws s3 ls s3://learnbot-production-backups/postgres/

# Restore from backup
./infrastructure/scripts/db-restore.sh production 2024-01-15-02-00-00
```

### Backup Schedule

| Environment | Frequency | Retention |
|-------------|-----------|-----------|
| Staging | Daily (automated RDS) | 3 days |
| Production | Daily (automated RDS) + Pre-deployment snapshot | 14 days |

---

## 8. CI/CD Pipeline

### Workflow Overview

```
Developer pushes to main
        │
        ▼
  Run all tests ──── FAIL ──→ Block deployment, notify Slack
        │
       PASS
        │
        ▼
  Security scan (govulncheck, Trivy, gosec)
        │
        ▼
  Build Docker images
        │
        ▼
  Push to ECR (staging tags)
        │
        ▼
  Run DB migrations (staging)
        │
        ▼
  Deploy to ECS (staging)
        │
        ▼
  Smoke tests ──── FAIL ──→ Alert team, auto-rollback
        │
       PASS
        │
        ▼
  Notify: "Staging healthy, ready for production"
        │
        ▼ (manual trigger with SHA + "deploy" confirmation)
  Production deployment
        │
        ▼
  Backup production DB
        │
        ▼
  Run DB migrations (production)
        │
        ▼
  Deploy to ECS (production)
        │
        ▼
  Health checks ──── FAIL ──→ Auto-rollback + Alert
        │
       PASS
        │
        ▼
  Tag release in git
```

### GitHub Actions Workflows

| Workflow | Trigger | Description |
|----------|---------|-------------|
| `test.yml` | Push/PR to `main` | Run all unit, integration, E2E tests |
| `security.yml` | Push/PR/Weekly | Security scanning |
| `deploy-staging.yml` | Push to `main` (after tests) | Auto-deploy to staging |
| `deploy-production.yml` | Manual dispatch | Deploy to production with approval |

### Required GitHub Secrets

| Secret | Description |
|--------|-------------|
| `AWS_ACCOUNT_ID` | AWS account ID |
| `AWS_ACCESS_KEY_ID` | AWS access key for deployments |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key for deployments |
| `STAGING_PRIVATE_SUBNET_IDS` | Comma-separated private subnet IDs (staging) |
| `STAGING_ECS_SG_ID` | ECS security group ID (staging) |
| `PROD_PRIVATE_SUBNET_IDS` | Comma-separated private subnet IDs (production) |
| `PROD_ECS_SG_ID` | ECS security group ID (production) |
| `SLACK_WEBHOOK_URL` | Slack webhook for deployment notifications |

---

## 9. Security

### Encryption

| Layer | Mechanism |
|-------|-----------|
| Data in transit | TLS 1.3 (enforced by ALB) |
| Data at rest (RDS) | AES-256 via AWS KMS |
| Data at rest (S3) | SSE-KMS |
| Passwords | bcrypt (cost factor 12) |
| JWT tokens | HS256 with 32+ char secret |

### Network Security

- All services run in **private subnets** (no public IPs)
- ALB is the only internet-facing component
- Security groups follow **least-privilege** principle:
  - ALB → ECS: port 8090 only
  - ECS → RDS: port 5432 only
  - ECS → Redis: port 6379 only
- VPC Flow Logs enabled for network monitoring
- No SSH access to containers (use ECS Exec or SSM)

### Access Control

- ECS tasks use **IAM roles** (no static credentials in containers)
- S3 bucket access restricted to ECS task role
- RDS accessible only from ECS security group
- Secrets Manager access restricted to specific task roles

### Security Scanning

The CI pipeline runs:
- `govulncheck` — Go vulnerability scanning
- `gosec` — Go security linter
- `Trivy` — Container image vulnerability scanning
- `Gitleaks` — Secret detection in git history

---

## 10. Monitoring & Alerting

### Accessing Dashboards

**Local:**
- Grafana: http://localhost:3001 (admin/admin)
- Prometheus: http://localhost:9090

**Production:**
- Grafana: https://monitoring.learnbot.example.com
- CloudWatch: AWS Console → CloudWatch → Dashboards

### Key Metrics

| Metric | Warning | Critical |
|--------|---------|----------|
| API error rate | > 5% | > 20% |
| P99 latency | > 2s | > 10s |
| CPU utilization | > 85% | — |
| Memory utilization | > 90% | — |
| DB connections | > 80% of max | — |
| Disk space | < 15% free | < 5% free |

### Alert Channels

- **Slack:** `#learnbot-ops` for all alerts
- **PagerDuty:** Critical alerts trigger on-call rotation
- **Email:** Weekly summary reports

### Grafana Dashboards

The `learnbot-overview` dashboard includes:
- Request rate and error rate per service
- P50/P95/P99 latency percentiles
- Database connection pool utilization
- Resume parse success/failure rate
- Job scraping statistics
- Go runtime metrics (GC, goroutines, memory)

---

## 11. Rollback Procedures

### Automatic Rollback

ECS automatically rolls back if health checks fail during deployment. The deployment script also triggers rollback on smoke test failure.

### Manual Rollback

```bash
# 1. Identify the previous stable task definition
aws ecs describe-task-definition \
  --task-definition learnbot-production-api-gateway \
  --query 'taskDefinition.revision'

# 2. Update service to use previous revision
aws ecs update-service \
  --cluster learnbot-production \
  --service learnbot-production-api-gateway \
  --task-definition learnbot-production-api-gateway:<previous-revision>

# 3. Wait for rollback to complete
aws ecs wait services-stable \
  --cluster learnbot-production \
  --services learnbot-production-api-gateway
```

### Database Rollback

If a migration needs to be rolled back:

```bash
# 1. Restore from pre-deployment backup
./infrastructure/scripts/db-restore.sh production <backup-timestamp>

# 2. Verify data integrity
psql "postgres://..." -c "SELECT COUNT(*) FROM users;"
```

> ⚠️ Database rollbacks are destructive. Always verify the backup before restoring.

---

## 12. Troubleshooting Guide

### Service Won't Start

```bash
# Check ECS task logs
aws logs tail /ecs/learnbot-production-api-gateway --follow

# Check task stopped reason
aws ecs describe-tasks \
  --cluster learnbot-production \
  --tasks <task-arn> \
  --query 'tasks[0].stoppedReason'
```

**Common causes:**
- Missing environment variable → Check Secrets Manager
- Database connection failed → Check security group rules
- Out of memory → Increase ECS task memory in Terraform

### Database Connection Issues

```bash
# Test connectivity from ECS
aws ecs execute-command \
  --cluster learnbot-production \
  --task <task-arn> \
  --container api-gateway \
  --interactive \
  --command "nc -zv <rds-endpoint> 5432"
```

**Common causes:**
- Security group not allowing ECS → RDS traffic
- Wrong DB credentials in Secrets Manager
- RDS instance not running

### High Error Rate

```bash
# Check recent errors in CloudWatch
aws logs filter-log-events \
  --log-group-name /ecs/learnbot-production-api-gateway \
  --filter-pattern "ERROR" \
  --start-time $(date -d '1 hour ago' +%s000)
```

### Resume Upload Failures

```bash
# Check S3 bucket permissions
aws s3api get-bucket-policy --bucket learnbot-production-resumes

# Check ECS task role has S3 access
aws iam simulate-principal-policy \
  --policy-source-arn <task-role-arn> \
  --action-names s3:PutObject \
  --resource-arns "arn:aws:s3:::learnbot-production-resumes/*"
```

### Slow Queries

```bash
# Enable slow query logging in RDS
aws rds modify-db-parameter-group \
  --db-parameter-group-name learnbot-production \
  --parameters "ParameterName=log_min_duration_statement,ParameterValue=1000,ApplyMethod=immediate"

# View slow queries
aws logs filter-log-events \
  --log-group-name /aws/rds/instance/learnbot-production/postgresql \
  --filter-pattern "duration"
```

### Terraform State Issues

```bash
# Unlock stuck state
terraform force-unlock <lock-id>

# Import existing resource
terraform import aws_ecs_service.api_gateway <cluster>/<service>

# Refresh state
terraform refresh -var-file=terraform.tfvars.production
```

### Operations Runbook

For detailed operational procedures, see [infrastructure/docs/RUNBOOK.md](../infrastructure/docs/RUNBOOK.md).
