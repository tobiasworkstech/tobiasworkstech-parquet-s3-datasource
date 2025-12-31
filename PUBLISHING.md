# Publishing Guide for Grafana Parquet S3 Datasource Plugin

This guide covers the complete process for signing and publishing the plugin to the Grafana Plugin Catalog.

## Prerequisites

1. **Grafana Cloud Account**: Create one at https://grafana.com/auth/sign-in/
2. **GitHub Repository**: Your plugin must be in a public GitHub repository
3. **Valid plugin.json**: Ensure all required fields are populated

## Step 1: Prepare Your Plugin

### Verify plugin.json

Ensure your `plugin.json` contains:

```json
{
  "type": "datasource",
  "name": "Parquet S3",
  "id": "tobiasworkstech-parquet-s3-datasource",
  "info": {
    "description": "A clear description of your plugin",
    "author": {
      "name": "Your Name",
      "url": "https://your-website.com"
    },
    "keywords": ["parquet", "s3", "datasource"],
    "logos": {
      "small": "img/logo.svg",
      "large": "img/logo.svg"
    },
    "version": "1.0.0"
  },
  "dependencies": {
    "grafanaDependency": ">=10.0.0",
    "plugins": []
  }
}
```

### Build the Plugin

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

## Step 2: Create Access Policy Token

1. Go to https://grafana.com/auth/sign-in/
2. Navigate to **My Account** → **Access Policies**
3. Click **Create access policy**
4. Configure the policy:
   - **Name**: `plugin-signing`
   - **Scopes**: Select `plugins:write`
5. Click **Create**
6. Create a token for this policy and save it securely

## Step 3: Sign the Plugin

### Local Signing

```bash
# Set the token
export GRAFANA_ACCESS_POLICY_TOKEN=your-token-here

# Sign the plugin
npx @grafana/sign-plugin@latest

# Verify the signature
cat dist/MANIFEST.txt
```

### Automated Signing (GitHub Actions)

1. Add your token as a repository secret:
   - Go to your repository → **Settings** → **Secrets and variables** → **Actions**
   - Create a new secret named `GRAFANA_ACCESS_POLICY_TOKEN`
   - Paste your token

2. Push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. The CI/CD workflow will automatically build, sign, and create a release

## Step 4: Package the Plugin

**Important**: The zip must contain a folder named exactly as your plugin ID, not loose files.

### Correct Packaging

```bash
cd plugin

# Rename dist folder to plugin ID
mv dist tobiasworkstech-parquet-s3-datasource

# Create the zip (zip the FOLDER, not contents)
zip -r tobiasworkstech-parquet-s3-datasource-1.0.0.zip tobiasworkstech-parquet-s3-datasource

# Rename back for continued development
mv tobiasworkstech-parquet-s3-datasource dist
```

### Verify Zip Structure

The zip should have this structure:
```
tobiasworkstech-parquet-s3-datasource-1.0.0.zip
└── tobiasworkstech-parquet-s3-datasource/    <-- folder with plugin ID
    ├── plugin.json
    ├── module.js
    ├── module.js.map
    ├── MANIFEST.txt
    ├── img/
    │   └── logo.svg
    ├── gpx_parquet_s3_datasource_linux_amd64
    ├── gpx_parquet_s3_datasource_linux_arm64
    ├── gpx_parquet_s3_datasource_darwin_amd64
    ├── gpx_parquet_s3_datasource_darwin_arm64
    └── gpx_parquet_s3_datasource_windows_amd64.exe
```

**Wrong** (don't do this):
```bash
cd dist
zip -r ../plugin.zip .   # ❌ This zips contents without parent folder
```

### Verify Before Uploading

```bash
# List zip contents to verify structure
unzip -l tobiasworkstech-parquet-s3-datasource-1.0.0.zip
```

You should see all files prefixed with `tobiasworkstech-parquet-s3-datasource/`.

## Step 5: Create GitHub Release

1. Go to your GitHub repository → **Releases** → **Create a new release**
2. Choose your tag (e.g., `v1.0.0`)
3. Upload the correctly packaged zip file
4. Publish the release

### Automated Release

The provided GitHub Actions workflow handles this automatically when you push a version tag.

## Step 5: Submit to Grafana Plugin Catalog

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

1. **Security**: No malicious code, proper handling of credentials
2. **Quality**: Works as described, good error handling
3. **Documentation**: Clear README, configuration instructions
4. **Best Practices**: Follows Grafana plugin guidelines

Review typically takes 1-2 weeks. You'll receive feedback via email if changes are needed.

## Best Practices for Approval

### Documentation
- Clear README with installation instructions
- Configuration examples
- Screenshots of the plugin in action
- Troubleshooting guide

### Code Quality
- Proper error handling
- No console.log statements in production
- TypeScript strict mode
- ESLint/golangci-lint clean

### Security
- Secrets stored in secureJsonData
- No hardcoded credentials
- Input validation
- Proper TLS handling

### User Experience
- Intuitive configuration UI
- Helpful error messages
- Loading states
- Responsive design

## Updating Published Plugins

1. Update the version in `plugin.json` and `package.json`
2. Build and sign the new version
3. Create a new GitHub release
4. The catalog will automatically detect the new version

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

### Plugin Not Loading

1. Check Grafana logs: `docker logs grafana`
2. Verify plugin is in correct directory
3. Check for unsigned plugin errors
4. Ensure `GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS` is set for unsigned plugins

### Catalog Submission Issues

- Ensure GitHub repository is public
- Verify all required plugin.json fields
- Check that releases contain the signed plugin zip
- Confirm plugin builds and works locally

## Resources

- [Grafana Plugin Development Guide](https://grafana.com/developers/plugin-tools/)
- [Publishing Best Practices](https://grafana.com/developers/plugin-tools/publish-a-plugin/publishing-best-practices)
- [Build Automation](https://grafana.com/developers/plugin-tools/publish-a-plugin/build-automation)
- [Plugin Signing](https://grafana.com/developers/plugin-tools/publish-a-plugin/sign-a-plugin)
- [Grafana Community Forum](https://community.grafana.com/)
