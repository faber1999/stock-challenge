# Infrastructure Setup

Terraform para desplegar:
- Frontend estatico en S3 + CloudFront
- Backend en ECS Fargate + ALB
- Conexion a base externa mediante Secrets Manager (`DATABASE_URL`)

## Estructura
- `modules/network`: VPC, subredes publicas, routing
- `modules/ecs_api`: ECS service, ALB, IAM, logs
- `modules/static_frontend`: S3 + CloudFront (+ proxy `/api/*` al ALB)
- `environments/dev`: composicion ejecutable para entorno dev

## Prerrequisitos
- Terraform >= 1.6
- AWS CLI configurada (`aws configure`)
- Imagen de backend publicada en ECR
- Secreto de `DATABASE_URL` en Secrets Manager

## Setup dev
1. Entra a `infrastructure/environments/dev`.
2. Copia variables:
```powershell
Copy-Item terraform.tfvars.example terraform.tfvars
```
3. Ajusta `terraform.tfvars` con valores reales:
- `backend_container_image`
- `database_url_secret_arn`
- `aws_region`
- `name_prefix`

4. Ejecuta:
```powershell
terraform init
terraform plan
terraform apply
```

## Actualizar infraestructura
Cada cambio en archivos `.tf` requiere:
```powershell
terraform plan
terraform apply
```

## Publicacion frontend (manual)
Despues del `apply`, construye y sube frontend:
```powershell
aws s3 sync ../../frontend/dist s3://<BUCKET_SALIDA_TERRAFORM> --delete
aws cloudfront create-invalidation --distribution-id <DIST_ID_SALIDA_TERRAFORM> --paths "/*"
```

## Importante
- No versionar `terraform.tfvars`, `terraform.tfstate` ni `.terraform/`.
- Usa `.env.example` y `terraform.tfvars.example` como plantillas.
