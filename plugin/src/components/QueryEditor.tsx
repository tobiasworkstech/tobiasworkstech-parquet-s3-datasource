import React, { ChangeEvent, useEffect, useState } from 'react';
import { InlineField, Input, Select, MultiSelect, TextArea, InlineFieldRow } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource/datasource';
import { ParquetS3DataSourceOptions, ParquetS3Query, defaultQuery, ParquetColumn } from '../types';

type Props = QueryEditorProps<DataSource, ParquetS3Query, ParquetS3DataSourceOptions>;

export function QueryEditor(props: Props) {
  const { datasource, query, onChange, onRunQuery } = props;
  const [buckets, setBuckets] = useState<SelectableValue[]>([]);
  const [columns, setColumns] = useState<ParquetColumn[]>([]);
  const [loadingBuckets, setLoadingBuckets] = useState(false);
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

  // Auto-load schema when bucket changes
  useEffect(() => {
    const autoLoadSchema = async () => {
      if (!query.bucket) {
        setColumns([]);
        return;
      }
      setLoadingColumns(true);
      try {
        // Get list of files and load schema from the first parquet file
        const result = await datasource.listFiles(query.bucket, query.pathPrefix);
        const parquetFiles = result.files.filter((f) => f.key.endsWith('.parquet'));
        if (parquetFiles.length > 0) {
          const schema = await datasource.getSchema(query.bucket, parquetFiles[0].key);
          setColumns(schema.columns);
        } else {
          setColumns([]);
        }
      } catch (error) {
        console.error('Failed to auto-load schema:', error);
        setColumns([]);
      } finally {
        setLoadingColumns(false);
      }
    };
    autoLoadSchema();
  }, [datasource, query.bucket, query.pathPrefix]);

  const onBucketChange = (value: SelectableValue<string> | null) => {
    onChange({ ...query, bucket: value?.value || '' });
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

  const onTimeColumnChange = (value: SelectableValue<string> | null) => {
    onChange({ ...query, timeColumn: value?.value || '' });
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

  const columnOptions = columns.map((c) => ({
    label: `${c.name} (${c.type})`,
    value: c.name,
  }));

  return (
    <>
      <InlineFieldRow>
        <InlineField label="Bucket" labelWidth={12} tooltip="S3 bucket containing Parquet files">
          <Select
            width={25}
            options={buckets}
            value={query.bucket ? { label: query.bucket, value: query.bucket } : null}
            onChange={onBucketChange}
            isLoading={loadingBuckets}
            placeholder="Select bucket..."
            isClearable
          />
        </InlineField>

        <InlineField label="Path Prefix" labelWidth={12} tooltip="Filter files by path prefix (optional, leave empty for root)">
          <Input
            width={20}
            value={query.pathPrefix || ''}
            onChange={onPathPrefixChange}
            onBlur={onRunQuery}
            placeholder="(root)"
          />
        </InlineField>

        <InlineField label="File Pattern" labelWidth={12} tooltip="File pattern to match">
          <Input
            width={15}
            value={query.filePattern || defaultQuery.filePattern}
            onChange={onFilePatternChange}
            onBlur={onRunQuery}
            placeholder="*.parquet"
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField
          label="Time Column"
          labelWidth={12}
          tooltip="Column to use as time field for time series (leave empty for non-time-series data)"
        >
          <Select
            width={25}
            options={columnOptions}
            value={query.timeColumn ? { label: query.timeColumn, value: query.timeColumn } : null}
            onChange={onTimeColumnChange}
            isLoading={loadingColumns}
            placeholder="(none)"
            isClearable
            allowCustomValue
          />
        </InlineField>

        <InlineField
          label="Columns"
          labelWidth={12}
          tooltip="Columns to include (leave empty for all)"
        >
          <MultiSelect
            width={40}
            options={columnOptions}
            value={(query.columns || []).map((c) => ({ label: c, value: c }))}
            onChange={onColumnsChange}
            isLoading={loadingColumns}
            placeholder="All columns"
            allowCustomValue
          />
        </InlineField>

        <InlineField label="Max Rows" labelWidth={10} tooltip="Maximum rows to return">
          <Input
            width={10}
            type="number"
            value={query.maxRows ?? defaultQuery.maxRows}
            onChange={onMaxRowsChange}
            onBlur={onRunQuery}
          />
        </InlineField>
      </InlineFieldRow>

      <InlineFieldRow>
        <InlineField
          label="Filter"
          labelWidth={12}
          tooltip="Filter condition (optional)"
          grow
        >
          <TextArea
            value={query.queryText || ''}
            onChange={onQueryTextChange}
            onBlur={onRunQuery}
            placeholder="(none)"
            rows={1}
          />
        </InlineField>
      </InlineFieldRow>
    </>
  );
}
