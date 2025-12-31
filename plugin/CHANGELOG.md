# Changelog

All notable changes to the Grafana Parquet S3 Datasource plugin will be documented in this file.

## [1.0.0] - 2025-01-01

### Added

- Initial release of the Parquet S3 Datasource plugin
- Read Parquet files directly from S3-compatible storage (AWS S3, MinIO, etc.)
- Support for time series and table data visualization
- Column selection for optimized queries
- SQL-like WHERE clause filtering
- Template variable support with `buckets()`, `files()`, and `columns()` queries
- Automatic schema detection from Parquet files
- Configurable options:
  - Endpoint URL for S3-compatible services
  - Region configuration
  - Path-style URL support (required for MinIO)
  - TLS verification skip option
  - Default bucket setting
- Multi-platform backend support (Linux, macOS, Windows on amd64/arm64)
