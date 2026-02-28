# LearnBot Infrastructure

This directory contains all infrastructure-as-code, deployment scripts, monitoring configuration, and operational documentation for the LearnBot platform.

## Directory Structure

```
infrastructure/
├── terraform/              # AWS infrastructure (IaC)
│   ├── providers.tf        # Terraform providers and backend config
│   ├── variables.tf        # Input variables
│   ├── main.tf             # Core infrastructure resources
│   ├── outputs.tf          # Output values
│   ├── terraform.tfvars.staging     # Staging environment values
│   └── terraform.tfvars.production  # Production environment values
├── kubernetes/             # Kubernetes manifests (alternative to ECS)
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── ingress.yaml
│   └── api-gateway/
│       ├── deployment.yaml
│       ├── service.yaml
│       └── hpa.yaml
├── monitoring/             # Observability stack
│   ├── prometheus.yml      # Prometheus scrape config
│   ├── alerts.yml          # Alerting rules
│   └── grafana/
│       ├── provisioning/   # Auto-provisioned datasources & dashboards
│       └── dashboards/     # Dashboard JSON definitions
├── nginx/
│   └── nginx.local.conf    # Local development reverse proxy
├── scripts/
│   ├── deploy.sh           # Manual deployment script
│   ├── db-backup.sh        # Database backup script
│   └── db-restore.sh       # Database restore script
└── docs/
    └── RUNBOOK.md          # Operations runbook
```

## Architecture Overview

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

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.6.0
- [AWS CLI](https://aws.amazon.com/cli/) >= 2.0 configured with appropriate credentials
- [Docker](https://www.docker.com/) >= 24.0
- [Docker Compose](https://docs.docker.com/compose/) >= 2.0
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (for Kubernetes deployments)

## Quick Start: Local Development

```bash
# 1. Copy environment template
cp .env.example .env
# Edit .env with your local values

# 2. Start all services
docker compose up -d

# 3. Run database migrations
docker compose exec postgres psql -U learnbot_admin -d learnbot \
  -f /docker-entrypoint-initdb.d/001_create_enums.sql

# 4. Access services
# Frontend:    http://localhost:3000
# API Gateway: http://localhost:8090
# Grafana:     http://localhost:3001 (admin/admin)
# Prometheus:  http://localhost:9090
```

## Terraform: Provision Infrastructure

### First-time Setup

```bash
# 1. Create S3 bucket for Terraform state (one-time)
aws s3 mb s3://learnbot-terraform-state --region us-east-1
aws s3api put-bucket-versioning \
  --bucket learnbot-terraform-state \
  --versioning-configuration Status=Enabled

# 2. Create DynamoDB table for state locking (one-time)
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

### Deploy Staging

```bash
cd infrastructure/terraform
terraform workspace new staging || terraform workspace select staging
terraform plan -var-file=terraform.tfvars.staging -out=staging.plan
terraform apply staging.plan
```

### Deploy Production

```bash
cd infrastructure/terraform
terraform workspace new production || terraform workspace select production
terraform plan -var-file=terraform.tfvars.production -out=production.plan
# Review the plan carefully!
terraform apply production.plan
```

### Destroy (Staging Only)

```bash
cd infrastructure/terraform
terraform workspace select staging
terraform destroy -var-file=terraform.tfvars.staging
```

## CI/CD Pipeline

### Workflows

| Workflow | Trigger | Description |
|----------|---------|-------------|
| `test.yml` | Push/PR | Run all unit, integration, and E2E tests |
| `security.yml` | Push/PR/Weekly | Security scanning (govulncheck, Trivy, gosec, Gitleaks) |
| `deploy-staging.yml` | Push to `main` | Auto-deploy to staging after tests pass |
| `deploy-production.yml` | Manual | Deploy specific SHA to production with approval gate |

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

### Deployment Flow

```
Developer pushes to main
        │
        ▼
  Run all tests ──── FAIL ──→ Block deployment, notify
        │
       PASS
        │
        ▼
  Build Docker images
        │
        ▼
  Push to ECR (staging)
        │
        ▼
  Run DB migrations (staging)
        │
        ▼
  Deploy to ECS (staging)
        │
        ▼
  Smoke tests ──── FAIL ──→ Alert team
        │
       PASS
        │
        ▼
  Notify: "Ready for production"
        │
        ▼ (manual trigger)
  Production deployment
  (requires SHA + "deploy" confirmation)
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

## Monitoring

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
| CPU utilization | > 85% | - |
| Memory utilization | > 90% | - |
| DB connections | > 80% of max | - |
| Disk space | < 15% free | < 5% free |

## Security

### Encryption

- **Data at rest:** RDS encrypted with AWS KMS, S3 with SSE-KMS
- **Data in transit:** TLS 1.3 enforced on ALB, HTTPS redirect
- **Secrets:** Stored in AWS Secrets Manager, never in environment files or code

### Network Security

- Services run in private subnets (no public IPs)
- ALB is the only public-facing component
- Security groups follow least-privilege principle
- VPC Flow Logs enabled for network monitoring

### Access Control

- ECS tasks use IAM roles (no static credentials)
- S3 bucket access restricted to ECS task role
- RDS accessible only from ECS security group
- No SSH access to containers (use ECS Exec or SSM)

## Database Backup Schedule

| Environment | Frequency | Retention |
|-------------|-----------|-----------|
| Staging | Daily (automated RDS) | 3 days |
| Production | Daily (automated RDS) + Pre-deployment snapshot | 14 days |

Backups are stored in S3 with SSE-KMS encryption and lifecycle policies.

## Cost Optimization

- ECS Fargate Spot used for non-critical workloads
- RDS storage autoscaling enabled
- CloudFront caching reduces origin requests
- S3 lifecycle policies move old data to cheaper storage classes
- Auto-scaling prevents over-provisioning

## Support

- **Slack:** `#learnbot-ops`
- **Runbook:** [infrastructure/docs/RUNBOOK.md](docs/RUNBOOK.md)
- **On-call:** PagerDuty rotation
