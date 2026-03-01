# Backend Setup

API en Go (pgx) para consulta, sincronizacion y recomendacion de stocks.

## Requisitos
- Go 1.22+
- Acceso a una base de datos CockroachDB/PostgreSQL compatible

## Configuracion
1. Copia variables de entorno:
```powershell
Copy-Item .env.example .env
```
2. Completa `DATABASE_URL` en `.env`.

Variables soportadas:
- `PORT` (default `8080`)
- `DATABASE_URL` (requerida)
- `STOCKS_API_URL` (default `https://api.karenai.click/swechallenge/list`)
- `STOCKS_API_TOKEN` (default `1`)
- `SYNC_TIMEOUT_SECONDS` (default `20`)
- `SYNC_MAX_PAGES` (default `50`)
- `AUTO_SYNC_ON_STARTUP` (default `false`)
- `CORS_ALLOWED_ORIGINS` (default `*`)

## Ejecutar local
```powershell
go mod tidy
go run .\cmd\api
```

## Test
```powershell
go test ./...
```

## Docker (opcional)
```powershell
docker build -t stock-backend:latest .
docker run --rm -p 8080:8080 --env-file .env stock-backend:latest
```

## Endpoints principales
- `GET /health`
- `GET /db-time`
- `POST /stocks/sync`
- `GET /stocks`
- `GET /stocks/{ticker}`
- `GET /stocks/recommendations`
- `GET /swagger`

## Nota de migraciones
- Las migraciones se ejecutan automaticamente al iniciar (`internal/db/migrate.go`).
