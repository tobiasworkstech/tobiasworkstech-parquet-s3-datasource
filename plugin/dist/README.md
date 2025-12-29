# Grafana Parquet S3 Datasource Plugin

A Grafana datasource plugin for reading Apache Parquet files stored in S3-compatible storage (AWS S3, MinIO, etc.).

## Features

- Read Parquet files directly from S3-compatible storage
- Support for AWS S3, MinIO, and other S3-compatible services
- Time series and table data visualization
- Column selection and filtering
- SQL-like WHERE clause filtering
- Template variable support for dynamic queries
- Automatic schema detection

## Installation

### From Grafana Plugin Catalog

1. Go to **Configuration > Plugins** in Grafana
2. Search for "Parquet S3"
3. Click **Install**

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/tobiasworkstech/grafana-parquet-s3-datasource/releases)
2. Extract to your Grafana plugins directory:
   ```bash
   unzip tobiasworkstech-parquet-s3-datasource-*.zip -d /var/lib/grafana/plugins/
   ```
3. Restart Grafana

### Using Docker

```bash
docker run -d \
  -p 3000:3000 \
  -e GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=tobiasworkstech-parquet-s3-datasource \
  -v /path/to/plugin:/var/lib/grafana/plugins/tobiasworkstech-parquet-s3-datasource \
  grafana/grafana:latest
```

## Configuration

### Adding the Datasource

1. Go to **Configuration > Data Sources**
2. Click **Add data source**
3. Search for "Parquet S3"
4. Configure the connection settings

### Connection Settings

| Setting | Description | Example |
|---------|-------------|---------|
| **Endpoint URL** | S3-compatible endpoint | `http://minio:9000` or `https://s3.amazonaws.com` |
| **Region** | AWS region | `us-east-1` |
| **Access Key ID** | S3 access key | `AKIAIOSFODNN7EXAMPLE` |
| **Secret Access Key** | S3 secret key (stored securely) | `wJalrXUtnFEMI...` |
| **Default Bucket** | Default bucket for queries | `my-data-bucket` |
| **Path Style** | Use path-style URLs (required for MinIO) | `true` |
| **Skip TLS Verify** | Skip SSL certificate verification | `false` |

### Query Configuration

| Field | Description |
|-------|-------------|
| **Bucket** | S3 bucket name |
| **Path Prefix** | Filter files by path prefix (e.g., `data/2024/`) |
| **File Pattern** | File pattern to match (supports wildcards like `*.parquet`) |
| **Time Column** | Column to use as time field for time series |
| **Columns** | Specific columns to select (empty = all) |
| **SQL Filter** | WHERE clause for filtering (e.g., `value > 100 AND status = 'active'`) |
| **Max Rows** | Maximum rows to return |

### Template Variables

The plugin supports template variable queries:

- `buckets()` - List all available buckets
- `files(bucket, prefix)` - List files in a bucket
- `columns(bucket, key)` - List columns from a parquet file schema

Example:
```
files(my-bucket, data/)
```

## Development

### Prerequisites

- Node.js 18+
- Go 1.21+
- Docker (for testing)

### Building

```bash
# Install dependencies
npm install
go mod download

# Build frontend
npm run build

# Build backend (all platforms)
mage buildBackend

# Or build for current platform only
mage buildDev
```

### Running locally

```bash
# Start development server with Docker
cd docker
docker-compose up -d

# Frontend watch mode
npm run dev
```

### Testing

```bash
# Frontend tests
npm run test:ci

# Backend tests
go test -v ./...

# Lint
npm run lint
go run github.com/golangci/golangci-lint/cmd/golangci-lint run
```

## Docker Deployment

A complete Docker Compose setup is provided with:
- Grafana with the plugin pre-installed
- MinIO for S3-compatible storage
- Sample Parquet data
- Pre-configured datasource and dashboard

### Quick Start

```bash
cd docker
cp .env.example .env
# Edit .env if needed
docker-compose up -d
```

Access:
- **Grafana**: http://localhost:3000 (admin/admin123)
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin123)

## Publishing

### Signing the Plugin

1. Create a Grafana Cloud account
2. Generate an access policy token
3. Set the environment variable:
   ```bash
   export GRAFANA_ACCESS_POLICY_TOKEN=your-token
   ```
4. Sign the plugin:
   ```bash
   npx @grafana/sign-plugin@latest
   ```

### Submitting to Grafana Plugin Catalog

1. Ensure your plugin meets the [publishing requirements](https://grafana.com/developers/plugin-tools/publish-a-plugin/publishing-best-practices)
2. Create a release on GitHub with the signed plugin
3. Submit via the [Grafana Plugin Submission](https://grafana.com/auth/sign-in/) portal

## CI/CD

The repository includes GitHub Actions workflows for:

- **CI/CD** (`ci-cd.yml`): Builds, tests, and creates releases on tags
- **Release** (`release.yml`): Manual workflow to create new versions

### Creating a Release

1. Push a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
2. Or use the Release workflow from GitHub Actions

## License

Apache License 2.0

## Contributing

1. Fork the repository
2. Create a feature branch
3. Submit a pull request

## Support

- [GitHub Issues](https://github.com/tobiasworkstech/grafana-parquet-s3-datasource/issues)
- [Grafana Community](https://community.grafana.com/)
