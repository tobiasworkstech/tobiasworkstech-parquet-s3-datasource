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

## Grafana Validation Status

**Submitted:** 2026-01-25

### Current Errors (Need to fix submission URL)

| Item | Status | Notes |
|------|--------|-------|
| no-go-manifest | ❌ ERROR | Grafana can't find go.mod - need correct ref |
| js-map-no-match | ❌ ERROR | Source maps don't match - Grafana checking wrong ref |
| unsigned-plugin | ⚠️ Expected | Will be signed after Grafana approves |

### Root Cause

**Grafana is checking the wrong git ref!** The source code at `main` branch has commits after v1.0.15, causing the mismatch.

### Solution: Provide Source Code URL with Tag Reference

When submitting to Grafana, use **one of these URL formats**:

**Option 1 - Tree URL:**
```
https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/tree/v1.0.15
```

**Option 2 - Archive URL (source zip):**
```
https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/archive/refs/tags/v1.0.15.zip
```

**Option 3 - If there's a separate ref field:**
- Repository: `https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource`
- Ref/Tag: `v1.0.15`

### Previously Fixed Issues

| Item | Status | Notes |
|------|--------|-------|
| binary-executable-permissions | ✅ FIXED | All binaries have 0755 |
| missing-readme | ✅ FIXED | README.md included |
| screenshots | ✅ FIXED | All PNGs included |

---

## Grafana Plugin Submission

**Plugin ID:** `tobiasworkstech-parquets3-datasource`

**Repository:** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource

**Source Code (with tag):** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/tree/v1.0.15

To submit/update on Grafana Plugin Catalog:
1. Go to https://grafana.com/auth/sign-in/
2. Navigate to My Account > My Plugins
3. Click 'Submit Plugin' or 'Update Submission'
4. **IMPORTANT:** Use the source code URL with the tag reference:
   ```
   https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/tree/v1.0.15
   ```
5. Follow the submission wizard

---

## Previous Issues Fixed

| Version | Issue | Fix |
|---------|-------|-----|
| v1.0.15 | Screenshots not in git | Removed `screenshot*.png` from .gitignore, added PNGs |
| v1.0.13 | Binary permissions lost | Added chmod 755 in sign-and-release job (after artifact download) |
| v1.0.12 | README.md path wrong | Fixed webpack path from `../README.md` to `README.md` |
| v1.0.11 | golangci-lint errors | Fixed errcheck and staticcheck issues |
