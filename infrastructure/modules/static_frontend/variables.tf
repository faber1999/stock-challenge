variable "name_prefix" {
  type = string
}

variable "bucket_force_destroy" {
  type    = bool
  default = false
}

variable "price_class" {
  type    = string
  default = "PriceClass_100"
}

variable "aliases" {
  type    = list(string)
  default = []
}

variable "acm_certificate_arn" {
  description = "ACM cert in us-east-1 for CloudFront custom domains"
  type        = string
  default     = ""
}

variable "default_root_object" {
  type    = string
  default = "index.html"
}

variable "api_origin_domain_name" {
  description = "Optional API origin domain (for example an ALB DNS name) used for /api/* proxying"
  type        = string
  default     = ""
}

variable "tags" {
  type    = map(string)
  default = {}
}
