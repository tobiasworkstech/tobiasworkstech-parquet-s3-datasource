import type { Configuration } from 'webpack';
import { merge } from 'webpack-merge';
import grafanaConfig from './.grafana/webpack';
import path from 'path';

const config = async (env: any): Promise<Configuration> => {
  const baseConfig = await grafanaConfig(env);

  return merge(baseConfig, {
    resolve: {
      alias: {
        '@': path.resolve(__dirname, '../../src'),
      },
    },
  });
};

export default config;
