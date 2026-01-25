# Release Information

## Latest Release: v1.0.15

**Release Date:** 2026-01-25

**Download:** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/releases/tag/v1.0.15

### Checksums

**File:** `tobiasworkstech-parquets3-datasource-1.0.15.zip`

| Algorithm | Checksum |
|-----------|----------|
| MD5 | `179c4f119a378f74f370b90d92f16912` |
| SHA1 | `9e02ccfadc0da40e14f56d943efe27c42af916e5` |
| SHA256 | `e05bb4a5d559ccaf62157756488c53bb9a6f963a9a879e3710210a768f920dce` |

### Changes in v1.0.15

- Add screenshot PNG files to git repository (were in .gitignore)
- Set executable permissions (0755) on all binary files
- Copy README.md to dist folder
- All Grafana validation errors fixed

### Package Contents

- Binary executables (linux/darwin/windows, amd64/arm64)
- README.md
- plugin.json
- Screenshots in img/
- go.mod

---

## Grafana Plugin Submission Form Values

**Use these exact values when submitting/updating on Grafana:**

| Field | Value |
|-------|-------|
| **Plugin URL (zip file)** | `https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/releases/download/v1.0.15/tobiasworkstech-parquets3-datasource-1.0.15.zip` |
| **MD5 or SHA1** | `179c4f119a378f74f370b90d92f16912` |
| **Source code URL** | `https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/tree/v1.0.15` |
| **Testing guidance** | `Use docker-compose in /docker folder to spin up Grafana + MinIO with sample data. See README for details.` |
| **Provisioning provided** | Yes |

### Important Notes

- **Source code URL must include the tag** (`/tree/v1.0.15`) so Grafana checks the correct ref
- Without the tag, Grafana checks `main` branch which has commits after v1.0.15, causing mismatch errors

---

## Grafana Validation Status

**Submitted:** 2026-01-25

### Previous Errors (Fixed by using correct URLs)

| Error | Cause | Fix |
|-------|-------|-----|
| no-go-manifest | Grafana checking wrong ref | Use `/tree/v1.0.15` in Source code URL |
| js-map-no-match | Source at main ≠ source at tag | Use `/tree/v1.0.15` in Source code URL |

### Fixed Issues

| Item | Status |
|------|--------|
| binary-executable-permissions | ✅ FIXED |
| missing-readme | ✅ FIXED |
| screenshots | ✅ FIXED |
| unsigned-plugin | ⚠️ Expected (will be signed after approval) |

---

## Plugin Information

**Plugin ID:** `tobiasworkstech-parquets3-datasource`

**Repository:** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource

**Description:** Grafana datasource plugin for reading Parquet files from S3-compatible storage

---

## Previous Issues Fixed

| Version | Issue | Fix |
|---------|-------|-----|
| v1.0.15 | Screenshots not in git | Removed `screenshot*.png` from .gitignore, added PNGs |
| v1.0.13 | Binary permissions lost | Added chmod 755 in sign-and-release job (after artifact download) |
| v1.0.12 | README.md path wrong | Fixed webpack path from `../README.md` to `README.md` |
| v1.0.11 | golangci-lint errors | Fixed errcheck and staticcheck issues |
