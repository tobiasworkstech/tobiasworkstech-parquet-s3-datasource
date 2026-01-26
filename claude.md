# Claude Code Context

## Grafana Plugin Development Guidelines

**CRITICAL**: This is a Grafana plugin. Follow these rules strictly:

1. **Do NOT modify the `.config` folder** - It's managed by Grafana plugin tools (`@grafana/create-plugin`). Any changes will be overwritten and cause validation failures.

2. **Use Magefile.go for Go builds** - The Grafana plugin SDK provides standard build methods. Never write custom Go build scripts.

3. **Use standard GitHub Actions** - Use `grafana/plugin-actions/build-plugin` for CI/CD. Do not create custom build workflows.

4. **Extend webpack only when necessary** - If you must customize webpack, do it by extending (not modifying) the base config.

5. **For Grafana plugin documentation**, fetch it directly from grafana.com

## Project Overview

This is a **Grafana datasource plugin** that reads Parquet files from S3-compatible storage (AWS S3, MinIO, etc.). It's a backend plugin written in Go with a React/TypeScript frontend.

## Project Structure

```
tobiasworkstech-parquets3-datasource/
├── src/                        # Frontend (React/TypeScript)
│   ├── components/             # React components (QueryEditor, ConfigEditor)
│   ├── datasource/             # Datasource implementation
│   ├── img/                    # Plugin logo and screenshots
│   └── plugin.json             # Plugin manifest
├── pkg/                        # Backend (Go)
│   ├── plugin/                 # Main datasource logic
│   ├── parquet/                # Parquet file reader
│   └── s3/                     # S3 client setup
├── dist/                       # Build output (generated)
├── .config/                    # Webpack configuration (DO NOT MODIFY)
├── docker/                     # Docker development environment
│   ├── docker-compose.yml      # Grafana + MinIO setup
│   └── provisioning/           # Auto-provisioned datasources
├── .github/workflows/          # GitHub Actions
│   ├── ci.yml                  # CI on push to main
│   └── release.yml             # Release on tag push
├── Magefile.go                 # Go build system (Grafana SDK)
├── package.json                # Frontend dependencies
├── go.mod                      # Go dependencies
└── PUBLISHING.md               # Plugin publishing guide
```

## Development Commands

### Frontend
```bash
npm install
npm run build           # Production build
npm run dev             # Watch mode
npm run lint            # Run ESLint
npm run test:ci         # Run tests
```

### Backend (using Mage)
```bash
# Install mage (one-time)
go install github.com/magefile/mage@latest

# Build all platforms
mage -v

# Or build specific target
mage build:linux
mage build:darwin
mage build:windows
```

### Manual Go build (for development)
```bash
go mod download
go build -o dist/gpx_parquet_s3_datasource_darwin_arm64 ./pkg
```

### Docker development
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
| `src/plugin.json` | Plugin manifest (ID, version, dependencies) |
| `pkg/plugin/datasource.go` | Main backend datasource logic |
| `pkg/parquet/reader.go` | Parquet file reading |
| `src/components/QueryEditor.tsx` | Query UI component |
| `src/components/ConfigEditor.tsx` | Datasource config UI |
| `src/datasource/datasource.ts` | Frontend datasource class |
| `Magefile.go` | Go build configuration (Grafana SDK) |

## Plugin ID

`tobiasworkstech-parquets3-datasource`

## Build Targets

Backend builds for (via Mage):
- linux/amd64, linux/arm64
- darwin/amd64, darwin/arm64
- windows/amd64

## Common Tasks

### Add a new query option:
1. Add field to `QueryModel` struct in `pkg/plugin/datasource.go`
2. Add field to `ParquetS3Query` interface in `src/types.ts`
3. Add UI control in `src/components/QueryEditor.tsx`
4. Handle the option in query execution logic

### Add a new config option:
1. Add field to `DatasourceSettings` in `pkg/plugin/datasource.go`
2. Add field to `ParquetS3DataSourceOptions` in `src/types.ts`
3. Add UI control in `src/components/ConfigEditor.tsx`

### Debug backend:
```bash
# View Grafana logs
docker logs grafana -f

# Or check plugin logs
docker exec grafana cat /var/log/grafana/grafana.log
```

## Grafana Plugin Rules

When modifying Go code, avoid:
- Direct filesystem access with `os` package in plugin code
- Accessing environment variables with `os.Getenv` in plugin code
- Hardcoded credentials
- `console.log` in production frontend code

## Testing the Plugin

1. Start docker environment: `cd docker && docker-compose up -d`
2. Open Grafana: http://localhost:3000
3. Go to Explore, select "Parquet S3 (MinIO)" datasource
4. Select bucket "parquet-data", time column "timestamp"
5. Run query to see sample data

## Publishing (Automated)

GitHub Actions handles everything using official `grafana/plugin-actions/build-plugin`:

**On push to `main`**: Runs CI (build, lint, test)
**On push of version tag (`v*`)**: Runs full release (build, sign, package, release)

### To Release a New Version

```bash
npm version patch   # bumps version, commits, and creates tag
git push origin main --tags   # push commit AND tag to trigger release
```

**Required**: Add `GRAFANA_ACCESS_POLICY_TOKEN` secret to your GitHub repository.

The release workflow will:
1. Build frontend (webpack) + backend (mage, all platforms)
2. Sign the plugin with Grafana token
3. Package correctly (zip with plugin ID folder)
4. Create GitHub Release with zip attached

See `PUBLISHING.md` for manual instructions.
