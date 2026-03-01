output "backend_api_base_url" {
  value = module.backend_api.api_base_url
}

output "backend_alb_dns_name" {
  value = module.backend_api.alb_dns_name
}

output "frontend_bucket_name" {
  value = module.frontend_static.bucket_name
}

output "frontend_cloudfront_distribution_id" {
  value = module.frontend_static.cloudfront_distribution_id
}

output "frontend_cloudfront_domain_name" {
  value = module.frontend_static.cloudfront_domain_name
}
