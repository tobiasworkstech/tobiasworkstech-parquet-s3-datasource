import React, { ChangeEvent, useCallback, useEffect, useState } from 'react';
import { InlineField, Input, Select, MultiSelect, TextArea, InlineFieldRow } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource/datasource';
import { ParquetS3DataSourceOptions, ParquetS3Query, defaultQuery, ParquetColumn } from '../types';

type Props = QueryEditorProps<DataSource, ParquetS3Query, ParquetS3DataSourceOptions>;

export function QueryEditor(props: Props) {
  const { datasource, query, onChange, onRunQuery } = props;
  const [buckets, setBuckets] = useState<SelectableValue[]>([]);
  const [files, setFiles] = useState<SelectableValue[]>([]);
  const [columns, setColumns] = useState<ParquetColumn[]>([]);
  const [loadingBuckets, setLoadingBuckets] = useState(false);
  const [loadingFiles, setLoadingFiles] = useState(false);
  const [loadingColumns, setLoadingColumns] = useState(false);

  // Load buckets on mount
  useEffect(() => {
    const loadBuckets = async () => {
      setLoadingBuckets(true);
      try {
        const bucketList = await datasource.listBuckets();
        setBuckets(bucketList.map((b) => ({ label: b, value: b })));
      } catch (error) {
        console.error('Failed to load buckets:', error);
      } finally {
        setLoadingBuckets(false);
      }
    };
    loadBuckets();
  }, [datasource]);

  // Load files when bucket or prefix changes
  useEffect(() => {
    const loadFiles = async () => {
      if (!query.bucket) {
        setFiles([]);
        return;
      }
      setLoadingFiles(true);
      try {
        const result = await datasource.listFiles(query.bucket, query.pathPrefix);
        setFiles(
          result.files
            .filter((f) => f.key.endsWith('.parquet'))
            .map((f) => ({ label: f.key, value: f.key }))
        );
      } catch (error) {
        console.error('Failed to load files:', error);
      } finally {
        setLoadingFiles(false);
      }
    };
    loadFiles();
  }, [datasource, query.bucket, query.pathPrefix]);

  // Load schema when a file is selected
  const loadSchema = useCallback(
    async (bucket: string, key: string) => {
      setLoadingColumns(true);
      try {
        const schema = await datasource.getSchema(bucket, key);
        setColumns(schema.columns);
      } catch (error) {
        console.error('Failed to load schema:', error);
        setColumns([]);
      } finally {
        setLoadingColumns(false);
      }
    },
    [datasource]
  );

  const onBucketChange = (value: SelectableValue<string>) => {
    onChange({ ...query, bucket: value.value || '' });
  };

  const onPathPrefixChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, pathPrefix: event.target.value });
  };

  const onFilePatternChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, filePattern: event.target.value });
  };

  const onQueryTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    onChange({ ...query, queryText: event.target.value });
  };

  const onTimeColumnChange = (value: SelectableValue<string>) => {
    onChange({ ...query, timeColumn: value.value || '' });
    onRunQuery();
  };

  const onColumnsChange = (values: SelectableValue<string>[]) => {
    onChange({ ...query, columns: values.map((v) => v.value || '') });
    onRunQuery();
  };

  const onMaxRowsChange = (event: ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(event.target.value, 10);
    onChange({ ...query, maxRows: isNaN(value) ? defaultQuery.maxRows : value });
  };

  const onFileSelect = (value: SelectableValue<string>) => {
    if (value.value && query.bucket) {
      loadSchema(query.bucket, value.value);
    }
  };

  const columnOptions = columns.map((c) => ({
    label: `${c.name} (${c.type})`,
    value: c.name,
  }));

  return (
    <>
      <InlineFieldRow>
        <InlineField label="Bucket" labelWidth={14} tooltip="S3 bucket containing Parquet files">
          <Select
            width={30}
            options={buckets}
            value={query.bucket ? { label: query.bucket, value: query.bucket } : null}
            onChange={onBucketChange}
            isLoading={loadingBuckets}
            placeholder="Select bucket..."
            isClearable
          />
        </InlineField>

        <InlineField label="Path Prefix" labelWidth={14} tooltip="Filter files by path prefix">
          <Input
            width={30}
            value={query.pathPrefix || ''}
            onChange={onPathPrefixChange}
            onBlur={onRunQuery}
            placeholder="data/2024/"
          />
        </InlineField>

        <InlineField label="File Pattern" labelWidth={14} tooltip="File pattern to match (supports wildcards)">
          <Input
            width={20}
            value={query.filePattern || defaultQuery.filePattern}
            onChange={onFilePatternChange}
            onBlur={onRunQuery}
            placeholder="*.parquet"
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField
          label="Load Schema From"
          labelWidth={14}
          tooltip="Select a file to load column schema"
        >
          <Select
            width={40}
            options={files}
            onChange={onFileSelect}
            isLoading={loadingFiles}
            placeholder="Select file to load schema..."
            isClearable
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField
          label="Time Column"
          labelWidth={14}
          tooltip="Column to use as time field for time series data"
        >
          <Select
            width={30}
            options={columnOptions}
            value={query.timeColumn ? { label: query.timeColumn, value: query.timeColumn } : null}
            onChange={onTimeColumnChange}
            isLoading={loadingColumns}
            placeholder="Select time column..."
            isClearable
            allowCustomValue
          />
        </InlineField>

        <InlineField
          label="Columns"
          labelWidth={14}
          tooltip="Columns to include (leave empty for all)"
        >
          <MultiSelect
            width={50}
            options={columnOptions}
            value={(query.columns || []).map((c) => ({ label: c, value: c }))}
            onChange={onColumnsChange}
            isLoading={loadingColumns}
            placeholder="Select columns..."
            allowCustomValue
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField
          label="Max Rows"
          labelWidth={14}
          tooltip="Maximum number of rows to return"
        >
          <Input
            width={15}
            type="number"
            value={query.maxRows ?? defaultQuery.maxRows}
            onChange={onMaxRowsChange}
            onBlur={onRunQuery}
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField
          label="SQL Filter"
          labelWidth={14}
          tooltip="SQL-like WHERE clause to filter data (e.g., 'value > 100 AND status = active')"
          grow
        >
          <TextArea
            value={query.queryText || ''}
            onChange={onQueryTextChange}
            onBlur={onRunQuery}
            placeholder="value > 100 AND status = 'active'"
            rows={2}
          />
        </InlineField>
      </InlineFieldRow>
    </>
  );
}
