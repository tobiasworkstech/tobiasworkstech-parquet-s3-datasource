# Publishing Guide for Grafana Parquet S3 Datasource Plugin

This guide covers the complete process for signing and publishing the plugin to the Grafana Plugin Catalog.

## Prerequisites

1. **Grafana Cloud Account**: Create one at https://grafana.com/auth/sign-in/
2. **GitHub Repository**: Your plugin must be in a public GitHub repository
3. **Valid plugin.json**: Ensure all required fields are populated

## Quick Start (Recommended)

The easiest way to release is using GitHub Actions automation:

### One-Time Setup

1. **Create Access Policy Token**:
   - Go to https://grafana.com/auth/sign-in/
   - Navigate to **My Account** → **Access Policies**
   - Click **Create access policy**
   - Name: `plugin-signing`, Scopes: `plugins:write`
   - Create a token and copy it

2. **Add Token to GitHub**:
   - Go to your repository → **Settings** → **Secrets and variables** → **Actions**
   - Create secret: `GRAFANA_ACCESS_POLICY_TOKEN`
   - Paste your token

### Release a New Version

```bash
cd plugin

# Bump version (updates package.json and creates git tag)
npm version patch   # or 'minor' or 'major'

# Push changes and tag to trigger workflow
git push origin main --tags
```

The GitHub Actions workflow will automatically:
1. ✅ Build frontend and backend for all platforms
2. ✅ Sign the plugin with your Grafana token
3. ✅ Package correctly (with proper folder structure)
4. ✅ Create a GitHub Release with the zip attached

Monitor progress at: `https://github.com/YOUR_USERNAME/YOUR_REPO/actions`

---

## Manual Process (Alternative)

If you prefer to build and release manually:

### Step 1: Build the Plugin

```bash
cd plugin

# Install dependencies
npm ci
go mod download

# Build frontend
npm run build

# Build backend for all platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o dist/gpx_parquet_s3_datasource_linux_amd64 ./pkg
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o dist/gpx_parquet_s3_datasource_linux_arm64 ./pkg
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o dist/gpx_parquet_s3_datasource_darwin_amd64 ./pkg
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" -o dist/gpx_parquet_s3_datasource_darwin_arm64 ./pkg
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o dist/gpx_parquet_s3_datasource_windows_amd64.exe ./pkg
```

### Step 2: Sign the Plugin

```bash
cd plugin

# Set your token
export GRAFANA_ACCESS_POLICY_TOKEN=your-token-here

# Sign the plugin
npx @grafana/sign-plugin@latest

# Verify the signature
cat dist/MANIFEST.txt
```

### Step 3: Package the Plugin

**Important**: The zip must contain a folder named exactly as your plugin ID.

```bash
cd plugin

# Rename dist folder to plugin ID
mv dist tobiasworkstech-parquet-s3-datasource

# Create the zip (zip the FOLDER, not contents)
zip -r tobiasworkstech-parquet-s3-datasource-1.0.0.zip tobiasworkstech-parquet-s3-datasource

# Rename back for continued development
mv tobiasworkstech-parquet-s3-datasource dist
```

**Verify the structure:**
```bash
unzip -l tobiasworkstech-parquet-s3-datasource-1.0.0.zip
```

Expected output - all files prefixed with the plugin ID folder:
```
tobiasworkstech-parquet-s3-datasource-1.0.0.zip
└── tobiasworkstech-parquet-s3-datasource/
    ├── plugin.json
    ├── module.js
    ├── MANIFEST.txt
    ├── img/logo.svg
    ├── gpx_parquet_s3_datasource_linux_amd64
    ├── gpx_parquet_s3_datasource_linux_arm64
    ├── gpx_parquet_s3_datasource_darwin_amd64
    ├── gpx_parquet_s3_datasource_darwin_arm64
    └── gpx_parquet_s3_datasource_windows_amd64.exe
```

❌ **Wrong** (don't do this):
```bash
cd dist
zip -r ../plugin.zip .   # This zips contents without parent folder!
```

### Step 4: Create GitHub Release

1. Go to your GitHub repository → **Releases** → **Create a new release**
2. Create tag (e.g., `v1.0.0`)
3. Upload the zip file
4. Publish the release

---

## Submit to Grafana Plugin Catalog

1. Go to https://grafana.com/auth/sign-in/
2. Navigate to **My Account** → **My Plugins**
3. Click **Submit Plugin**
4. Enter your GitHub repository URL
5. Follow the submission wizard:
   - Select the release version
   - Review plugin metadata
   - Accept the terms
6. Submit for review

## Plugin Review Process

Grafana reviews plugins for:

- **Security**: No malicious code, proper handling of credentials
- **Quality**: Works as described, good error handling
- **Documentation**: Clear README, configuration instructions
- **Best Practices**: Follows Grafana plugin guidelines

Review typically takes 1-2 weeks. You'll receive feedback via email if changes are needed.

## Updating Published Plugins

Using GitHub Actions (recommended):
```bash
cd plugin
npm version patch  # bumps version and creates tag
git push origin main --tags
```

The catalog will automatically detect the new version from your GitHub release.

## Troubleshooting

### Signing Fails

```bash
# Check token is set
echo $GRAFANA_ACCESS_POLICY_TOKEN

# Verify plugin structure
ls -la dist/

# Check plugin.json is valid
cat dist/plugin.json | jq .
```

### Wrong Zip Structure

```bash
# Verify zip contents
unzip -l your-plugin.zip

# All files should be prefixed with: tobiasworkstech-parquet-s3-datasource/
```

### Plugin Not Loading

1. Check Grafana logs: `docker logs grafana`
2. Verify plugin is in correct directory
3. Check for unsigned plugin errors

## Resources

- [Grafana Plugin Development Guide](https://grafana.com/developers/plugin-tools/)
- [Build Automation](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation)
- [Package a Plugin](https://grafana.com/developers/plugin-tools/publish-a-plugin/package-a-plugin)
- [Plugin Signing](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)
