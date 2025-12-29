import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource/datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { ParquetS3Query, ParquetS3DataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, ParquetS3Query, ParquetS3DataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
