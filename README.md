# Stock Challenge Monorepo

Monorepo con aplicacion fullstack para consulta y recomendacion de stocks:
- `backend`: API en Go
- `frontend`: app web en Vue 3 + TypeScript
- `infrastructure`: IaC con Terraform para AWS

## Estructura
- `backend/`: servicio API, conexion a base de datos y logica de negocio
- `frontend/`: interfaz de usuario y consumo de API
- `infrastructure/`: modulos Terraform (network, ecs_api, static_frontend) y entorno `dev`

## Setup local (rapido)
1. Backend
```powershell
cd backend
Copy-Item .env.example .env
# completa DATABASE_URL
go mod tidy
go run .\cmd\api
```

2. Frontend
```powershell
cd frontend
Copy-Item .env.example .env
# para local: VITE_API_BASE_URL=http://localhost:8080
pnpm install
pnpm dev
```

## Setup infraestructura (AWS)
```powershell
cd infrastructure\environments\dev
Copy-Item terraform.tfvars.example terraform.tfvars
# completa valores reales (imagen ECR, secret ARN, etc.)
terraform init
terraform plan
terraform apply
```

## Deploy frontend en S3/CloudFront
```powershell
cd frontend
pnpm build
aws s3 sync dist s3://<BUCKET_NAME> --delete
aws cloudfront create-invalidation --distribution-id <DISTRIBUTION_ID> --paths "/*"
```

## Notas
- No subir archivos sensibles (`.env`, `terraform.tfvars`, `terraform.tfstate`).
- Usa los README internos de cada carpeta para detalle:
  - `backend/README.md`
  - `frontend/README.md`
  - `infrastructure/README.md`
