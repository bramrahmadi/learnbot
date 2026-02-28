#!/usr/bin/env bash
# ─── LearnBot Deployment Script ──────────────────────────────────────────────
# Usage:
#   ./deploy.sh --env staging --tag <sha>
#   ./deploy.sh --env production --tag <sha>
#   ./deploy.sh --env staging --tag latest --service api-gateway
#
# Prerequisites: aws-cli, docker (for local builds)

set -euo pipefail

ENVIRONMENT=""
IMAGE_TAG=""
SERVICE="all"
DRY_RUN=false
AWS_REGION="${AWS_REGION:-us-east-1}"

log()   { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [INFO] $*"; }
error() { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [ERROR] $*" >&2; }
warn()  { echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [WARN] $*"; }

while [[ $# -gt 0 ]]; do
  case $1 in
    --env)     ENVIRONMENT="$2"; shift 2 ;;
    --tag)     IMAGE_TAG="$2"; shift 2 ;;
    --service) SERVICE="$2"; shift 2 ;;
    --dry-run) DRY_RUN=true; shift ;;
    *) error "Unknown argument: $1"; exit 1 ;;
  esac
done

# ─── Validation ──────────────────────────────────────────────────────────────
if [[ -z "$ENVIRONMENT" ]]; then
  error "Environment is required. Use --env staging|production"
  exit 1
fi

if [[ -z "$IMAGE_TAG" ]]; then
  error "Image tag is required. Use --tag <sha|latest>"
  exit 1
fi

if [[ "$ENVIRONMENT" != "staging" && "$ENVIRONMENT" != "production" ]]; then
  error "Invalid environment: $ENVIRONMENT. Must be 'staging' or 'production'"
  exit 1
fi

CLUSTER_NAME="learnbot-${ENVIRONMENT}-cluster"
ECR_REGISTRY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

SERVICES=(api-gateway resume-parser job-aggregator learning-resources)
if [[ "$SERVICE" != "all" ]]; then
  SERVICES=("$SERVICE")
fi

log "Deploying to ${ENVIRONMENT} environment"
log "Image tag: ${IMAGE_TAG}"
log "Services: ${SERVICES[*]}"
log "Cluster: ${CLUSTER_NAME}"

if [[ "$DRY_RUN" == "true" ]]; then
  warn "DRY RUN mode - no changes will be made"
fi

# ─── Deploy Each Service ─────────────────────────────────────────────────────
for svc in "${SERVICES[@]}"; do
  log "Deploying ${svc}..."

  IMAGE_URI="${ECR_REGISTRY}/learnbot-${ENVIRONMENT}/${svc}:${IMAGE_TAG}"
  ECS_SERVICE="learnbot-${ENVIRONMENT}-${svc}"

  # Verify image exists
  if ! aws ecr describe-images \
    --repository-name "learnbot-${ENVIRONMENT}/${svc}" \
    --image-ids imageTag="${IMAGE_TAG}" \
    --region "$AWS_REGION" > /dev/null 2>&1; then
    error "Image not found: ${IMAGE_URI}"
    exit 1
  fi

  if [[ "$DRY_RUN" == "true" ]]; then
    log "[DRY RUN] Would update service ${ECS_SERVICE} with image ${IMAGE_URI}"
    continue
  fi

  # Get current task definition
  TASK_DEF_ARN=$(aws ecs describe-services \
    --cluster "$CLUSTER_NAME" \
    --services "$ECS_SERVICE" \
    --region "$AWS_REGION" \
    --query 'services[0].taskDefinition' \
    --output text)

  # Update task definition with new image
  NEW_TASK_DEF=$(aws ecs describe-task-definition \
    --task-definition "$TASK_DEF_ARN" \
    --region "$AWS_REGION" \
    --query 'taskDefinition' \
    | jq --arg IMAGE "$IMAGE_URI" --arg SVC "$svc" \
      '.containerDefinitions |= map(if .name == $SVC then .image = $IMAGE else . end)
       | del(.taskDefinitionArn, .revision, .status, .requiresAttributes, .placementConstraints, .compatibilities, .registeredAt, .registeredBy)')

  NEW_TASK_DEF_ARN=$(echo "$NEW_TASK_DEF" | aws ecs register-task-definition \
    --region "$AWS_REGION" \
    --cli-input-json /dev/stdin \
    --query 'taskDefinition.taskDefinitionArn' \
    --output text)

  log "Registered new task definition: ${NEW_TASK_DEF_ARN}"

  # Update service
  aws ecs update-service \
    --cluster "$CLUSTER_NAME" \
    --service "$ECS_SERVICE" \
    --task-definition "$NEW_TASK_DEF_ARN" \
    --force-new-deployment \
    --region "$AWS_REGION" > /dev/null

  log "Service update initiated for ${svc}"

  # Wait for stability
  log "Waiting for ${svc} to stabilize..."
  aws ecs wait services-stable \
    --cluster "$CLUSTER_NAME" \
    --services "$ECS_SERVICE" \
    --region "$AWS_REGION"

  log "✅ ${svc} deployed successfully"
done

log "Deployment completed successfully"
