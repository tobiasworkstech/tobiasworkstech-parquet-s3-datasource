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

---

## Grafana Validation Status

**Submitted:** 2026-01-25
**Status:** ⚠️ Warnings only (can be reviewed)

### Validation Results

| Item | Status | Notes |
|------|--------|-------|
| binary-executable-permissions | ✅ FIXED | All binaries have 0755 |
| missing-readme | ✅ FIXED | README.md included |
| screenshots | ✅ FIXED | All PNGs included |
| unsigned-plugin | ⚠️ Expected | Will be signed after Grafana approves |
| source-code-not-provided | ⚠️ Check URL | Ensure correct repo URL was entered |
| sponsorshiplink | 💡 Optional | Just a suggestion |

### Notes

- **unsigned-plugin** warning is expected for new plugins - Grafana will sign it after approval
- For source-code warning, ensure the exact repository URL was entered:
  ```
  https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource
  ```
  No trailing slashes, no `.git` extension, no branch/tag references.

---

## Grafana Plugin Submission

**Plugin ID:** `tobiasworkstech-parquets3-datasource`

**Repository:** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource

To submit/update on Grafana Plugin Catalog:
1. Go to https://grafana.com/auth/sign-in/
2. Navigate to My Account > My Plugins
3. Click 'Submit Plugin' or 'Update Submission'
4. Enter repository URL: https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource
5. Follow the submission wizard

---

## Previous Issues Fixed

| Version | Issue | Fix |
|---------|-------|-----|
| v1.0.15 | Screenshots not in git | Removed `screenshot*.png` from .gitignore, added PNGs |
| v1.0.13 | Binary permissions lost | Added chmod 755 in sign-and-release job (after artifact download) |
| v1.0.12 | README.md path wrong | Fixed webpack path from `../README.md` to `README.md` |
| v1.0.11 | golangci-lint errors | Fixed errcheck and staticcheck issues |
