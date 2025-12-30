import { DataQuery, DataSourceJsonData } from '@grafana/data';

/**
 * Query interface for Parquet S3 datasource
 */
export interface ParquetS3Query extends DataQuery {
  /** S3 bucket name */
  bucket?: string;
  /** Path prefix in the bucket (e.g., 'data/2024/') */
  pathPrefix?: string;
  /** Specific file pattern to match (supports wildcards like '*.parquet') */
  filePattern?: string;
  /** SQL-like query to filter/transform data */
  queryText?: string;
  /** Column to use as time field */
  timeColumn?: string;
  /** Columns to select (empty means all) */
  columns?: string[];
  /** Maximum number of rows to return */
  maxRows?: number;
}

/**
 * Default values for new queries
 */
export const defaultQuery: Partial<ParquetS3Query> = {
  bucket: '',
  pathPrefix: '',
  filePattern: '*.parquet',
  queryText: '',
  timeColumn: '',
  columns: [],
  maxRows: 10000,
};

/**
 * Configuration stored in Grafana's encrypted storage
 */
export interface ParquetS3SecureJsonData {
  /** S3 Secret Access Key */
  secretAccessKey?: string;
}

/**
 * Main datasource configuration options
 */
export interface ParquetS3DataSourceOptions extends DataSourceJsonData {
  /** S3 endpoint URL (e.g., 'http://minio:9000' or 'https://s3.amazonaws.com') */
  endpoint?: string;
  /** S3 region (e.g., 'us-east-1') */
  region?: string;
  /** S3 Access Key ID */
  accessKeyId?: string;
  /** Default bucket to use if not specified in query */
  defaultBucket?: string;
  /** Use path-style URLs (required for MinIO and some S3-compatible services) */
  pathStyle?: boolean;
  /** Skip SSL certificate verification (use with caution) */
  insecureSkipVerify?: boolean;
}

/**
 * Schema information for a Parquet file
 */
export interface ParquetSchema {
  columns: ParquetColumn[];
}

/**
 * Column definition from Parquet schema
 */
export interface ParquetColumn {
  name: string;
  type: string;
  nullable: boolean;
}

/**
 * Response from listing files in S3
 */
export interface S3FileListResponse {
  files: S3FileInfo[];
  truncated: boolean;
  continuationToken?: string;
}

/**
 * Information about a single S3 file
 */
export interface S3FileInfo {
  key: string;
  size: number;
  lastModified: string;
  etag?: string;
}
