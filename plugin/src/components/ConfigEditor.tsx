import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput, Switch, FieldSet } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { ParquetS3DataSourceOptions, ParquetS3SecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<ParquetS3DataSourceOptions, ParquetS3SecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields, secureJsonData } = options;

  const onEndpointChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        endpoint: event.target.value,
      },
    });
  };

  const onRegionChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        region: event.target.value,
      },
    });
  };

  const onAccessKeyIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        accessKeyId: event.target.value,
      },
    });
  };

  const onSecretAccessKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...secureJsonData,
        secretAccessKey: event.target.value,
      },
    });
  };

  const onResetSecretAccessKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...secureJsonFields,
        secretAccessKey: false,
      },
      secureJsonData: {
        ...secureJsonData,
        secretAccessKey: '',
      },
    });
  };

  const onDefaultBucketChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        defaultBucket: event.target.value,
      },
    });
  };

  const onPathStyleChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        pathStyle: event.currentTarget.checked,
      },
    });
  };

  const onInsecureSkipVerifyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        insecureSkipVerify: event.currentTarget.checked,
      },
    });
  };

  return (
    <>
      <FieldSet label="S3 Connection">
        <InlineField
          label="Endpoint URL"
          labelWidth={20}
          tooltip="S3-compatible endpoint URL (e.g., http://minio:9000 or https://s3.amazonaws.com)"
        >
          <Input
            width={40}
            value={jsonData.endpoint || ''}
            onChange={onEndpointChange}
            placeholder="http://minio:9000"
          />
        </InlineField>

        <InlineField
          label="Region"
          labelWidth={20}
          tooltip="AWS region (e.g., us-east-1). For MinIO, use 'us-east-1' as default."
        >
          <Input
            width={40}
            value={jsonData.region || ''}
            onChange={onRegionChange}
            placeholder="us-east-1"
          />
        </InlineField>

        <InlineField
          label="Default Bucket"
          labelWidth={20}
          tooltip="Default bucket to use if not specified in queries"
        >
          <Input
            width={40}
            value={jsonData.defaultBucket || ''}
            onChange={onDefaultBucketChange}
            placeholder="my-bucket"
          />
        </InlineField>
      </FieldSet>

      <FieldSet label="Authentication">
        <InlineField
          label="Access Key ID"
          labelWidth={20}
          tooltip="S3 Access Key ID"
        >
          <Input
            width={40}
            value={jsonData.accessKeyId || ''}
            onChange={onAccessKeyIdChange}
            placeholder="AKIAIOSFODNN7EXAMPLE"
          />
        </InlineField>

        <InlineField
          label="Secret Access Key"
          labelWidth={20}
          tooltip="S3 Secret Access Key (stored securely)"
        >
          <SecretInput
            width={40}
            isConfigured={secureJsonFields?.secretAccessKey ?? false}
            value={secureJsonData?.secretAccessKey || ''}
            onChange={onSecretAccessKeyChange}
            onReset={onResetSecretAccessKey}
            placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
          />
        </InlineField>
      </FieldSet>

      <FieldSet label="Advanced Options">
        <InlineField
          label="Path Style"
          labelWidth={20}
          tooltip="Use path-style URLs (required for MinIO and some S3-compatible services)"
        >
          <Switch
            value={jsonData.pathStyle ?? true}
            onChange={onPathStyleChange}
          />
        </InlineField>

        <InlineField
          label="Skip TLS Verify"
          labelWidth={20}
          tooltip="Skip SSL/TLS certificate verification (use with caution)"
        >
          <Switch
            value={jsonData.insecureSkipVerify ?? false}
            onChange={onInsecureSkipVerifyChange}
          />
        </InlineField>
      </FieldSet>
    </>
  );
}
