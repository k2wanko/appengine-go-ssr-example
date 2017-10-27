var path = require('path')
var webpack = require('webpack')
var merge = require('webpack-merge')
var VueSSRServerPlugin = require('vue-server-renderer/server-plugin')
var HtmlWebpackPlugin = require('html-webpack-plugin')
var NodeExternals = require('webpack-node-externals')
var ExtractTextPlugin = require("extract-text-webpack-plugin")

var isProd = process.env.NODE_ENV === 'production'

var clientConfig = Object.assign({}, require('./webpack.config'))
clientConfig.plugins = []

module.exports = merge(clientConfig, {
  target: 'node',
  devtool: 'none',
  entry: './src/server.js',
  output: {
    path: path.resolve(__dirname, './backend'),
    filename: 'server-bundle.js',
    libraryTarget: 'commonjs2',
  },
  // externals: NodeExternals({
  //   whitelist: /\.css$/
  // }),
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development'),
      'process.env.VUE_ENV': '"server"'
    }),
    new ExtractTextPlugin({
      filename: '[name].css?[contenthash]',
    }),
    new VueSSRServerPlugin(),
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: 'src/index.html',
      inject: false,
      minify: isProd ? {
        html5: true,
        collapseWhitespace: true,
      } : false,
    }),
  ],
})