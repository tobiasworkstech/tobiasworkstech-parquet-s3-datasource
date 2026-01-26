# Release Information

## Latest Release: v1.0.16

**Release Date:** 2026-01-26

**Download:** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/releases/tag/v1.0.16

### Checksums

**File:** `tobiasworkstech-parquets3-datasource-1.0.16.zip`

| Algorithm | Checksum |
|-----------|----------|
| MD5 | `e55450c073bd3e72db413ceedef1395b` |
| SHA1 | `e0948403e6c3158ea38a5bc36c8e154ad346dbb6` |
| SHA256 | `bdb77b75b39bba654901f6de1495b7cf1bfe5f69b5e7f60db792d38d2cf1f240` |

### Changes in v1.0.16

- Version bump for resubmission

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
| **Plugin URL (zip file)** | `https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/releases/download/v1.0.16/tobiasworkstech-parquets3-datasource-1.0.16.zip` |
| **MD5 or SHA1** | `e55450c073bd3e72db413ceedef1395b` |
| **Source code URL** | `https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/tree/v1.0.16` |
| **Testing guidance** | `Use docker-compose in /docker folder to spin up Grafana + MinIO with sample data. See README for details.` |
| **Provisioning provided** | Yes |

### Important Notes

- **Source code URL must include the tag** (`/tree/v1.0.16`) so Grafana checks the correct ref
- Without the tag, Grafana checks `main` branch which has commits after v1.0.16, causing mismatch errors
- Using `/releases/tag/` instead of `/tree/` causes validation failures

---

## Grafana Validation Status

**Submitted:** 2026-01-26 (v1.0.16)

### Previous Errors (Fixed by using correct URLs)

| Error | Cause | Fix |
|-------|-------|-----|
| no-go-manifest | Grafana checking wrong ref | Use `/tree/v1.0.16` in Source code URL |
| js-map-no-match | Source at main ≠ source at tag | Use `/tree/v1.0.16` in Source code URL |

**Note:** v1.0.16 first submission used `/releases/tag/` URL which failed. Resubmit with `/tree/` URL.

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
