import type { Configuration } from 'webpack';
import { DefinePlugin } from 'webpack';
import CopyWebpackPlugin from 'copy-webpack-plugin';
import ESLintPlugin from 'eslint-webpack-plugin';
import ForkTsCheckerWebpackPlugin from 'fork-ts-checker-webpack-plugin';
import path from 'path';

const PLUGIN_ID = 'tobiasworkstech-parquet-s3-datasource';

const config = async (env: { production?: boolean }): Promise<Configuration> => {
  const isProduction = env.production === true;

  return {
    mode: isProduction ? 'production' : 'development',
    devtool: isProduction ? 'source-map' : 'eval-source-map',
    target: 'web',
    context: path.join(process.cwd()),
    entry: {
      module: './src/module.ts',
    },
    output: {
      clean: true,
      filename: '[name].js',
      library: {
        type: 'amd',
      },
      path: path.resolve(process.cwd(), 'dist'),
      publicPath: `public/plugins/${PLUGIN_ID}/`,
    },
    externals: [
      'lodash',
      'jquery',
      'moment',
      'slate',
      'emotion',
      '@emotion/react',
      '@emotion/css',
      'prismjs',
      'slate-plain-serializer',
      '@grafana/slate-react',
      'react',
      'react-dom',
      'react-redux',
      'redux',
      'rxjs',
      'd3',
      'angular',
      '@grafana/ui',
      '@grafana/runtime',
      '@grafana/data',
    ],
    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          { from: 'src/plugin.json', to: '.' },
          { from: '../README.md', to: '.', noErrorOnMissing: true },
          { from: 'CHANGELOG.md', to: '.', noErrorOnMissing: true },
          { from: 'src/img/', to: 'img/', noErrorOnMissing: true },
        ],
      }),
      new ForkTsCheckerWebpackPlugin({
        async: Boolean(!isProduction),
        issue: {
          include: [{ file: '**/*.{ts,tsx}' }],
        },
        typescript: { configFile: path.join(process.cwd(), 'tsconfig.json') },
      }),
      new ESLintPlugin({
        extensions: ['.ts', '.tsx'],
        lintDirtyModulesOnly: Boolean(!isProduction),
      }),
      new DefinePlugin({
        'process.env.NODE_ENV': JSON.stringify(isProduction ? 'production' : 'development'),
      }),
    ],
    resolve: {
      extensions: ['.js', '.jsx', '.ts', '.tsx'],
      unsafeCache: true,
    },
    module: {
      rules: [
        {
          test: /\.[tj]sx?$/,
          exclude: /node_modules/,
          use: {
            loader: 'swc-loader',
            options: {
              jsc: {
                parser: {
                  syntax: 'typescript',
                  tsx: true,
                  decorators: false,
                  dynamicImport: true,
                },
                target: 'es2018',
                loose: false,
              },
            },
          },
        },
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader'],
        },
        {
          test: /\.s[ac]ss$/,
          use: ['style-loader', 'css-loader', 'sass-loader'],
        },
        {
          test: /\.(png|jpe?g|gif|svg)$/,
          type: 'asset/resource',
          generator: {
            filename: 'img/[hash][ext]',
          },
        },
        {
          test: /\.(woff|woff2|eot|ttf|otf)(\?v=\d+\.\d+\.\d+)?$/,
          type: 'asset/resource',
          generator: {
            filename: 'fonts/[hash][ext]',
          },
        },
      ],
    },
  };
};

export default config;
