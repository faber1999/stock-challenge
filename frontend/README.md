# Frontend Setup

Aplicacion Vue 3 + TypeScript + Pinia + Tailwind.

## Requisitos
- Node.js 20+
- pnpm 9/10

## Configuracion
1. Copia variables de entorno:
```powershell
Copy-Item .env.example .env
```
2. Ajusta `VITE_API_BASE_URL`:
- Local: `http://localhost:8080`
- Con CloudFront + proxy API: `/api`

## Ejecutar local
```powershell
pnpm install
pnpm dev
```

## Build
```powershell
pnpm build
pnpm preview
```

## Publicar en S3 + CloudFront
```powershell
aws s3 sync dist s3://<NOMBRE_BUCKET> --delete
aws cloudfront create-invalidation --distribution-id <DISTRIBUTION_ID> --paths "/*"
```

## Funcionalidades principales
- Sync manual de datos (`POST /stocks/sync`)
- Listado con busqueda, filtros, orden y paginacion
- Top recomendaciones (`GET /stocks/recommendations`)
- Detalle por ticker (`GET /stocks/{ticker}`)
- Tema claro/oscuro
- Internacionalizacion (ES/EN)
