# Claude Code Context

## Project Overview

This is a **Grafana datasource plugin** that reads Parquet files from S3-compatible storage (AWS S3, MinIO, etc.). It's a backend plugin written in Go with a React/TypeScript frontend.

## Project Structure

```
tobiasworkstech-parquet-s3-datasource/
├── plugin/                     # Main plugin code
│   ├── src/                    # Frontend (React/TypeScript)
│   │   ├── components/         # React components (QueryEditor, ConfigEditor)
│   │   ├── datasource/         # Datasource implementation
│   │   ├── img/                # Plugin logo
│   │   └── plugin.json         # Plugin manifest
│   ├── pkg/                    # Backend (Go)
│   │   ├── plugin/             # Main datasource logic
│   │   ├── parquet/            # Parquet file reader
│   │   └── s3/                 # S3 client setup
│   ├── dist/                   # Build output (generated)
│   ├── .config/                # Webpack configuration
│   ├── package.json            # Frontend dependencies
│   └── go.mod                  # Go dependencies
├── docker/                     # Docker development environment
│   ├── docker-compose.yml      # Grafana + MinIO setup
│   └── provisioning/           # Auto-provisioned datasources
├── screenshots.py              # Playwright screenshot generator
└── PUBLISHING.md               # Plugin publishing guide
```

## Development Commands

### From plugin/ directory:
```bash
# Frontend
npm install
npm run build
npm run dev         # Watch mode

# Backend
go mod download
go build -o dist/gpx_parquet_s3_datasource ./pkg

# Sign plugin
export GRAFANA_ACCESS_POLICY_TOKEN=your-token
npx @grafana/sign-plugin@latest
```

### Docker development:
```bash
cd docker
docker-compose up -d        # Start Grafana + MinIO
docker-compose down         # Stop
docker-compose logs -f      # View logs
```

Access points:
- Grafana: http://localhost:3000 (admin/admin123)
- MinIO Console: http://localhost:9001 (minioadmin/minioadmin123)

## Key Files

| File | Purpose |
|------|---------|
| `plugin/src/plugin.json` | Plugin manifest (ID, version, dependencies) |
| `plugin/pkg/plugin/datasource.go` | Main backend datasource logic |
| `plugin/pkg/parquet/reader.go` | Parquet file reading |
| `plugin/src/components/QueryEditor.tsx` | Query UI component |
| `plugin/src/components/ConfigEditor.tsx` | Datasource config UI |
| `plugin/src/datasource/datasource.ts` | Frontend datasource class |

## Plugin ID

`tobiasworkstech-parquet-s3-datasource`

## Build Targets

Backend builds for:
- linux/amd64, linux/arm64
- darwin/amd64, darwin/arm64
- windows/amd64

## Common Tasks

### Add a new query option:
1. Add field to `QueryModel` struct in `plugin/pkg/plugin/datasource.go`
2. Add field to `ParquetS3Query` interface in `plugin/src/types.ts`
3. Add UI control in `plugin/src/components/QueryEditor.tsx`
4. Handle the option in query execution logic

### Add a new config option:
1. Add field to `DatasourceSettings` in `plugin/pkg/plugin/datasource.go`
2. Add field to `ParquetS3DataSourceOptions` in `plugin/src/types.ts`
3. Add UI control in `plugin/src/components/ConfigEditor.tsx`

### Debug backend:
```bash
# View Grafana logs
docker logs grafana -f

# Or check plugin logs
docker exec grafana cat /var/log/grafana/grafana.log
```

## Grafana Plugin Rules

When modifying Go code, avoid:
- Direct filesystem access with `os` package in plugin code (use only in build tools)
- Hardcoded credentials
- `console.log` in production frontend code

The `magefile.go` is intentionally in the repo root (not `plugin/`) to avoid Grafana's code validator flagging build tool code.

## Testing the Plugin

1. Start docker environment: `cd docker && docker-compose up -d`
2. Open Grafana: http://localhost:3000
3. Go to Explore, select "Parquet S3 (MinIO)" datasource
4. Select bucket "parquet-data", time column "timestamp"
5. Run query to see sample data

## Publishing

See `PUBLISHING.md` for complete publishing instructions. Quick steps:
1. Build: `cd plugin && npm run build && go build...`
2. Sign: `npx @grafana/sign-plugin@latest`
3. Package: `cd dist && zip -r ../plugin-1.0.0.zip .`
4. Create GitHub release with zip file
5. Submit to Grafana plugin catalog
