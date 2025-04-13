const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

const isProduction = process.env.NODE_ENV === 'production';

const config = {
    entry:  {
        vendor: './src/scripts/vendor.js',
        main: './src/scripts/main.js',
        auth: './src/scripts/auth.js',
    },
    output: {
        filename: '[name].bundle.js',
        path: path.resolve(__dirname, 'dist'),
        clean: true,
    },
    devServer: {
        static: path.resolve(__dirname, '.dist'),
        compress: true,
        open: true,
        host: 'localhost',
        port: 8080,
        headers: {
            'Cache-Control': 'no-store',
        },
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: '[name].css'
        }),
        new HtmlWebpackPlugin({
            filename: 'main.html',
            template: './src/pages/main.pug',
            chunks: ['vendor', 'main'],
            minify: false,
        }),
        new HtmlWebpackPlugin({
            filename: 'auth.html',
            template: './src/pages/auth/auth.pug',
            chunks: ['vendor', 'auth'],
            minify: false,
        }),
    ],
    module: {
        rules: [
            {
                test: /\.m?js$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader",
                    options: {
                    presets: ['@babel/preset-env']
                    }
                }
            },
            {
                test: /\.css$/i,
                use: [ MiniCssExtractPlugin.loader, {
                    loader: 'css-loader',
                    options: {
                      importLoaders: 1
                    }
                  }, 'postcss-loader'],
            },
            {
                test: /\.pug$/,
                loader: 'pug-loader',
                options: {
                  self: true // Доступ к миксинам
                }
            },
            {
                test: /\.(png|svg|jpg|jpeg|gif)$/i,
                type: 'asset/resource',
            },
            {
                test: /\.(woff|woff2|eot|ttf|otf)$/i,
                type: 'asset/resource',
            },
            {
                test: /\.html$/i,
                type: 'html-loader',
            }
        ],
    },
    optimization: { // ПРоблема с devServer т.к. много страниц & предотвращение дублирования кода
        runtimeChunk: 'single',
    },
};

module.exports = () => {
    if (isProduction) {
        config.mode = 'production';
        config.plugins.push(new CleanWebpackPlugin());
    } else {
        config.mode = 'development';
    }
    return config;
};
