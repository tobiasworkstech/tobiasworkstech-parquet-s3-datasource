import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  MutableDataFrame,
  FieldType,
  MetricFindValue,
} from '@grafana/data';
import { getBackendSrv, getTemplateSrv, FetchResponse } from '@grafana/runtime';
import { lastValueFrom, Observable } from 'rxjs';

import {
  ParquetS3Query,
  ParquetS3DataSourceOptions,
  defaultQuery,
  ParquetSchema,
  S3FileListResponse,
} from '../types';

export class DataSource extends DataSourceApi<ParquetS3Query, ParquetS3DataSourceOptions> {
  defaultBucket?: string;

  constructor(instanceSettings: DataSourceInstanceSettings<ParquetS3DataSourceOptions>) {
    super(instanceSettings);
    this.defaultBucket = instanceSettings.jsonData.defaultBucket;
  }

  /**
   * Get the base URL for backend resource API calls
   */
  private getResourceUrl(path: string): string {
    return `/api/datasources/uid/${this.uid}/resources${path}`;
  }

  /**
   * Apply template variables to query
   */
  private applyTemplateVariables(query: ParquetS3Query, scopedVars: any): ParquetS3Query {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      bucket: templateSrv.replace(query.bucket || '', scopedVars),
      pathPrefix: templateSrv.replace(query.pathPrefix || '', scopedVars),
      filePattern: templateSrv.replace(query.filePattern || '', scopedVars),
      queryText: templateSrv.replace(query.queryText || '', scopedVars),
      timeColumn: templateSrv.replace(query.timeColumn || '', scopedVars),
    };
  }

  /**
   * Main query method - called when dashboard panels request data
   */
  async query(options: DataQueryRequest<ParquetS3Query>): Promise<DataQueryResponse> {
    const { range, scopedVars } = options;
    
    const promises = options.targets.map(async (target) => {
      if (target.hide) {
        return new MutableDataFrame();
      }

      const query = this.applyTemplateVariables(
        { ...defaultQuery, ...target },
        scopedVars
      );

      try {
        const response = await lastValueFrom(
          getBackendSrv().fetch<any>({
            url: this.getResourceUrl('/query'),
            method: 'POST',
            data: {
              ...query,
              refId: target.refId,
              timeRange: {
                from: range?.from.valueOf(),
                to: range?.to.valueOf(),
              },
            },
          }) as unknown as Observable<FetchResponse<any>>
        );

        if (response.data?.frames?.[0]) {
          const backendFrame = response.data.frames[0];

          // Build fields array with values
          const fields = (backendFrame.schema?.fields || []).map((field: any, i: number) => ({
            name: field.name,
            type: this.mapFieldType(field.type),
            values: backendFrame.data?.values?.[i] || [],
          }));

          const frame = new MutableDataFrame({
            refId: target.refId,
            fields: fields,
          });

          return frame;
        }

        return new MutableDataFrame({ refId: target.refId, fields: [] });
      } catch (error: any) {
        throw new Error(`Query failed: ${error.message || 'Unknown error'}`);
      }
    });

    const data = await Promise.all(promises);
    return { data };
  }

  /**
   * Map backend field types to Grafana field types
   */
  private mapFieldType(backendType: string): FieldType {
    const typeMap: Record<string, FieldType> = {
      'time': FieldType.time,
      'number': FieldType.number,
      'string': FieldType.string,
      'boolean': FieldType.boolean,
      'other': FieldType.other,
    };
    return typeMap[backendType] || FieldType.string;
  }

  /**
   * Test datasource connectivity
   * For backend plugins, this delegates to the backend's CheckHealth method
   */
  async testDatasource(): Promise<any> {
    try {
      const response = await lastValueFrom(
        getBackendSrv().fetch<{ status: string; message: string }>({
          url: `/api/datasources/uid/${this.uid}/health`,
          method: 'GET',
        }) as unknown as Observable<FetchResponse<{ status: string; message: string }>>
      );

      if (response.data?.status === 'OK') {
        return {
          status: 'success',
          message: response.data?.message || 'Successfully connected to S3.',
        };
      }

      return {
        status: 'error',
        message: response.data?.message || 'Connection failed',
      };
    } catch (error: any) {
      return {
        status: 'error',
        message: `Connection failed: ${error.data?.message || error.message || 'Unknown error'}`,
      };
    }
  }

  /**
   * Get schema for a specific file or bucket
   */
  async getSchema(bucket: string, key: string): Promise<ParquetSchema> {
    const response = await lastValueFrom(
      getBackendSrv().fetch<ParquetSchema>({
        url: this.getResourceUrl('/schema'),
        method: 'POST',
        data: { bucket, key },
      }) as unknown as Observable<FetchResponse<ParquetSchema>>
    );
    return response.data;
  }

  /**
   * List files in a bucket
   */
  async listFiles(bucket: string, prefix?: string): Promise<S3FileListResponse> {
    const response = await lastValueFrom(
      getBackendSrv().fetch<S3FileListResponse>({
        url: this.getResourceUrl('/files'),
        method: 'POST',
        data: { bucket, prefix },
      }) as unknown as Observable<FetchResponse<S3FileListResponse>>
    );
    return response.data;
  }

  /**
   * List available buckets
   */
  async listBuckets(): Promise<string[]> {
    const response = await lastValueFrom(
      getBackendSrv().fetch<{ buckets: string[] }>({
        url: this.getResourceUrl('/buckets'),
        method: 'GET',
      }) as unknown as Observable<FetchResponse<{ buckets: string[] }>>
    );
    return response.data.buckets;
  }

  /**
   * Support for template variable queries
   */
  async metricFindQuery(query: string, options?: any): Promise<MetricFindValue[]> {
    const templateSrv = getTemplateSrv();
    const interpolatedQuery = templateSrv.replace(query, options?.scopedVars);

    // Parse the query - supports: buckets(), files(bucket, prefix), columns(bucket, key)
    const bucketsMatch = interpolatedQuery.match(/^buckets\(\)$/);
    if (bucketsMatch) {
      const buckets = await this.listBuckets();
      return buckets.map((b) => ({ text: b, value: b }));
    }

    const filesMatch = interpolatedQuery.match(/^files\(([^,]+)(?:,\s*(.+))?\)$/);
    if (filesMatch) {
      const bucket = filesMatch[1].trim();
      const prefix = filesMatch[2]?.trim() || '';
      const result = await this.listFiles(bucket, prefix);
      return result.files.map((f) => ({ text: f.key, value: f.key }));
    }

    const columnsMatch = interpolatedQuery.match(/^columns\(([^,]+),\s*(.+)\)$/);
    if (columnsMatch) {
      const bucket = columnsMatch[1].trim();
      const key = columnsMatch[2].trim();
      const schema = await this.getSchema(bucket, key);
      return schema.columns.map((c) => ({ text: c.name, value: c.name }));
    }

    return [];
  }
}
