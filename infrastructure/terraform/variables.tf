variable "aws_region" {
  description = "AWS region for all resources"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Deployment environment (staging or production)"
  type        = string
  validation {
    condition     = contains(["staging", "production"], var.environment)
    error_message = "Environment must be 'staging' or 'production'."
  }
}

variable "project_name" {
  description = "Project name used for resource naming"
  type        = string
  default     = "learnbot"
}

# ─── VPC ─────────────────────────────────────────────────────────────────────

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b", "us-east-1c"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
}

# ─── ECS / Compute ───────────────────────────────────────────────────────────

variable "api_gateway_image" {
  description = "Docker image for the API gateway service"
  type        = string
  default     = "ghcr.io/learnbot/api-gateway:latest"
}

variable "resume_parser_image" {
  description = "Docker image for the resume parser service"
  type        = string
  default     = "ghcr.io/learnbot/resume-parser:latest"
}

variable "job_aggregator_image" {
  description = "Docker image for the job aggregator service"
  type        = string
  default     = "ghcr.io/learnbot/job-aggregator:latest"
}

variable "learning_resources_image" {
  description = "Docker image for the learning resources service"
  type        = string
  default     = "ghcr.io/learnbot/learning-resources:latest"
}

variable "api_gateway_cpu" {
  description = "CPU units for the API gateway task (1024 = 1 vCPU)"
  type        = number
  default     = 512
}

variable "api_gateway_memory" {
  description = "Memory (MiB) for the API gateway task"
  type        = number
  default     = 1024
}

variable "api_gateway_desired_count" {
  description = "Desired number of API gateway tasks"
  type        = number
  default     = 2
}

variable "api_gateway_min_count" {
  description = "Minimum number of API gateway tasks for auto-scaling"
  type        = number
  default     = 1
}

variable "api_gateway_max_count" {
  description = "Maximum number of API gateway tasks for auto-scaling"
  type        = number
  default     = 10
}

variable "resume_parser_cpu" {
  description = "CPU units for the resume parser task"
  type        = number
  default     = 1024
}

variable "resume_parser_memory" {
  description = "Memory (MiB) for the resume parser task"
  type        = number
  default     = 2048
}

variable "resume_parser_desired_count" {
  description = "Desired number of resume parser tasks"
  type        = number
  default     = 2
}

# ─── Database ────────────────────────────────────────────────────────────────

variable "db_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = "db.t3.medium"
}

variable "db_allocated_storage" {
  description = "Allocated storage in GiB for RDS"
  type        = number
  default     = 20
}

variable "db_max_allocated_storage" {
  description = "Maximum allocated storage in GiB for RDS autoscaling"
  type        = number
  default     = 100
}

variable "db_name" {
  description = "Name of the PostgreSQL database"
  type        = string
  default     = "learnbot"
}

variable "db_username" {
  description = "Master username for the PostgreSQL database"
  type        = string
  default     = "learnbot_admin"
  sensitive   = true
}

variable "db_backup_retention_days" {
  description = "Number of days to retain automated database backups"
  type        = number
  default     = 7
}

variable "db_multi_az" {
  description = "Enable Multi-AZ deployment for RDS"
  type        = bool
  default     = false
}

# ─── S3 / Storage ────────────────────────────────────────────────────────────

variable "resume_bucket_name" {
  description = "S3 bucket name for resume storage (must be globally unique)"
  type        = string
  default     = ""
}

# ─── SSL / Domain ────────────────────────────────────────────────────────────

variable "domain_name" {
  description = "Primary domain name for the application"
  type        = string
  default     = "learnbot.example.com"
}

variable "acm_certificate_arn" {
  description = "ARN of the ACM certificate for HTTPS"
  type        = string
  default     = ""
}

# ─── Monitoring ──────────────────────────────────────────────────────────────

variable "alert_email" {
  description = "Email address for CloudWatch alarm notifications"
  type        = string
  default     = "ops@learnbot.example.com"
}

variable "log_retention_days" {
  description = "Number of days to retain CloudWatch logs"
  type        = number
  default     = 30
}
