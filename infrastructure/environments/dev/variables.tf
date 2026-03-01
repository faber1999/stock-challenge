variable "aws_region" {
  description = "AWS region for environment"
  type        = string
  default     = "us-east-1"
}

variable "name_prefix" {
  description = "Global prefix for resource names"
  type        = string
  default     = "stock-challenge-dev"
}

variable "vpc_cidr" {
  type    = string
  default = "10.40.0.0/16"
}

variable "azs" {
  type    = list(string)
  default = ["us-east-1a", "us-east-1b"]
}

variable "public_subnet_cidrs" {
  type    = list(string)
  default = ["10.40.1.0/24", "10.40.2.0/24"]
}

variable "backend_container_image" {
  description = "Container image URI for backend API"
  type        = string
}

variable "backend_desired_count" {
  type    = number
  default = 1
}

variable "database_url_secret_arn" {
  description = "Secrets Manager ARN containing CockroachDB DATABASE_URL"
  type        = string
}

variable "stocks_api_token_secret_arn" {
  description = "Optional secret ARN for STOCKS_API_TOKEN"
  type        = string
  default     = ""
}

variable "stocks_api_url" {
  type    = string
  default = "https://api.karenai.click/swechallenge/list"
}

variable "sync_timeout_seconds" {
  type    = number
  default = 20
}

variable "sync_max_pages" {
  type    = number
  default = 50
}

variable "auto_sync_on_startup" {
  type    = bool
  default = false
}

variable "cors_allowed_origins" {
  description = "CORS origins for backend, comma-separated"
  type        = string
  default     = "*"
}

variable "frontend_bucket_force_destroy" {
  type    = bool
  default = false
}

variable "frontend_aliases" {
  description = "Optional custom domains for CloudFront"
  type        = list(string)
  default     = []
}

variable "frontend_acm_certificate_arn" {
  description = "Optional ACM certificate ARN in us-east-1"
  type        = string
  default     = ""
}

variable "tags" {
  type = map(string)
  default = {
    Project     = "stock-challenge"
    Environment = "dev"
    ManagedBy   = "terraform"
  }
}
