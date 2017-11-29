var ExtractTextPlugin = require('extract-text-webpack-plugin');
var ForkTsCheckerWebpackPlugin = require('fork-ts-checker-webpack-plugin');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var UglifyJsPlugin = require('uglifyjs-webpack-plugin');
var webpack = require('webpack');

var path = require('path');

module.exports = {
    entry: {
        index: './src/index.tsx',
        app: './src/app.ts'
    },

    output: {
        filename: '[name]-bundle.min.js',
        path: path.resolve(__dirname, 'dist')
    },

    context: __dirname, // to automatically find tsconfig.json

    devtool: 'source-map',

    watch: true,

    resolve: {
        extensions: ['.js', '.json', '.ts', '.tsx', '.scss', '.css'],
        modules: ['node_modules', 'src']
    },

    module: {
        rules: [
            {
                test: /\.tsx?$/,
                include: /src/,
                use: [
                    { loader: 'cache-loader' },
                    {
                        loader: 'thread-loader',
                        options: {
                            // There should be 1 cpu for the fork-ts-checker-webpack-plugin.
                            workers: require('os').cpus().length - 1
                        }
                    },
                    {
                        loader: 'babel-loader',
                        options: {
                            presets: ['es2015', 'react'],
                            sourceMap: true
                        }
                    },
                    {
                        loader: 'ts-loader',
                        options: {
                            // IMPORTANT! Use happyPackMode mode to speed-up
                            // compilation and reduce errors reported to webpack.
                            happyPackMode: true
                        }
                    },
                    {
                        loader: 'webpack-module-hot-accept'
                    }
                ]
            },
            {
                test: /\.css$|\.scss$/,
                include: /src/,
                use: ExtractTextPlugin.extract({
                    use: [
                        { loader: 'cache-loader' },
                        {
                            loader: 'thread-loader',
                            options: {
                                workers: require('os').cpus().length - 1
                            }
                        },
                        {
                            loader: 'css-loader',
                            options: {
                                sourceMap: true,
                                minimize: true
                            }
                        },
                        {
                            loader: 'sass-loader',
                            options: {
                                includePaths: [
                                    path.resolve(__dirname, './src/sass'),
                                    path.resolve(__dirname, './node_modules/compass-mixins/lib')
                                ],
                                sourceMap: true
                            }
                        }
                    ],
                    fallback: 'style-loader'
                })
            }
        ]
    },

    devServer: {
        port: 9000,
        host: 'localhost',
        index: 'index.html',
        contentBase: path.join(__dirname, 'dist'),
        // Don't refresh if hot loading fails. Good while
        // implementing the client interface.
        hotOnly: true,
        // Refresh on errors too.
        // hot: true,
        inline: true,
        open: true,
        historyApiFallback: true,
        compress: true
    },

    plugins: [
        // Enable this plugin to let webpack communicate changes
        // to WDS. --hot sets this automatically.
        new webpack.HotModuleReplacementPlugin(),
        new webpack.NamedModulesPlugin(),
        new webpack.optimize.CommonsChunkPlugin({
            name: 'commons',
            filename: 'commons.min.js'
        }),
        new webpack.DefinePlugin({
            'process.env': {
                NODE_ENV: JSON.stringify('production')
            }
        }),
        // Webpack built-in UglifyJsPlugin doesn't work with webpack-dev-server version 2.8.0+
        new UglifyJsPlugin({
            sourceMap: true
        }),
        new ForkTsCheckerWebpackPlugin({
            checkSyntacticErrors: true,
            watch: ['./src']
        }),
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: './src/templates/index.ejs',
            chunks: ['commons', 'index']
        }),
        new HtmlWebpackPlugin({
            filename: 'app.html',
            template: './src/templates/app.ejs',
            chunks: ['commons', 'app']
        }),
        new ExtractTextPlugin({
            filename: '[name]-style.min.css',
            disable: process.env.NODE_ENV === 'development',
            allChunks: true
        })
    ]
};
