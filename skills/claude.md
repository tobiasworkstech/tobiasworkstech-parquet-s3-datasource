# Release Information

## Latest Release: v1.0.12

**Release Date:** 2026-01-25

**Download:** https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource/releases/tag/v1.0.12

### Checksums

**File:** `tobiasworkstech-parquets3-datasource-1.0.12.zip`

| Algorithm | Checksum |
|-----------|----------|
| MD5 | `e8224b892556772078983b4cf8ca8653` |
| SHA1 | `bbbf08b0b12bd3c3c6ce80d209e3867a499f19ed` |
| SHA256 | `11169acd3e5b72ad607db3d22adaa7d2957662e4fd7b25ff0a7a270272574ca6` |

### Changes in v1.0.12

- Set executable permissions (0755) on all binary files
- Copy README.md to dist folder
- Fix webpack README.md path
- Update plugin version to 1.0.12

Fixes Grafana validation errors:
- binary-executable-permissions errors
- missing-readme error
- path-not-exists errors for screenshots

---

## Previous Release: v1.0.11

**File:** `tobiasworkstech-parquets3-datasource-1.0.11.zip`

| Algorithm | Checksum |
|-----------|----------|
| MD5 | `490721da62bccb5b6e171d41483f4e21` |
| SHA1 | `315f18d41feda0dbe7e6d94d8b3c94c9c72ec405` |
| SHA256 | `95d11c510b1ac17bc3760c182c799c7e10a09a473085a75c64fabcd0127bba01` |

### Changes in v1.0.11

- Fixed golangci-lint errors:
  - Check error return values for `Body.Close()` and `rows.Close()`
  - Lifted `maxRows` condition into loop condition (staticcheck QF1006)

---

## Grafana Plugin Submission

**Plugin ID:** `tobiasworkstech-parquets3-datasource`

To submit to Grafana Plugin Catalog:
1. Go to https://grafana.com/auth/sign-in/
2. Navigate to My Account > My Plugins
3. Click 'Submit Plugin' or 'Update Submission'
4. Enter repository URL: https://github.com/tobiasworkstech/tobiasworkstech-parquets3-datasource
5. Follow the submission wizard
