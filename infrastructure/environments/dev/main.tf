locals {
  backend_secrets = merge(
    {
      DATABASE_URL = var.database_url_secret_arn
    },
    var.stocks_api_token_secret_arn != "" ? {
      STOCKS_API_TOKEN = var.stocks_api_token_secret_arn
    } : {}
  )
}

module "network" {
  source = "../../modules/network"

  name_prefix         = var.name_prefix
  vpc_cidr            = var.vpc_cidr
  azs                 = var.azs
  public_subnet_cidrs = var.public_subnet_cidrs
  tags                = var.tags
}

module "backend_api" {
  source = "../../modules/ecs_api"

  name_prefix       = var.name_prefix
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
  region            = var.aws_region

  container_image = var.backend_container_image
  container_port  = 8080
  desired_count   = var.backend_desired_count

  environment = {
    PORT                 = "8080"
    STOCKS_API_URL       = var.stocks_api_url
    SYNC_TIMEOUT_SECONDS = tostring(var.sync_timeout_seconds)
    SYNC_MAX_PAGES       = tostring(var.sync_max_pages)
    AUTO_SYNC_ON_STARTUP = tostring(var.auto_sync_on_startup)
    CORS_ALLOWED_ORIGINS = var.cors_allowed_origins
  }

  secrets = local.backend_secrets
  tags    = var.tags
}

module "frontend_static" {
  source = "../../modules/static_frontend"

  name_prefix            = var.name_prefix
  bucket_force_destroy   = var.frontend_bucket_force_destroy
  aliases                = var.frontend_aliases
  acm_certificate_arn    = var.frontend_acm_certificate_arn
  api_origin_domain_name = module.backend_api.alb_dns_name
  tags                   = var.tags
}
