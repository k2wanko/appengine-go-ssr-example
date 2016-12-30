var path = require('path')
var webpack = require('webpack')

var ExtractTextPlugin = require("extract-text-webpack-plugin")

const client = {
  entry: './src/client.js',
  output: {
    path: path.resolve(__dirname, '../backend/app'),
    publicPath: '/',
    filename: 'build.js'
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue',
        options: {
          loaders: {
            css: ExtractTextPlugin.extract({
              loader: 'css-loader',
              fallbackLoader: 'vue-style-loader'
            })
          }
        }
      },
      {
        test: /\.js$/,
        loader: 'babel',
        exclude: /node_modules/
      },
      {
        test: /\.(png|jpg|gif|svg)$/,
        loader: 'file',
        options: {
          name: '[name].[ext]?[hash]'
        }
      },
      {
        test: /\.json$/,
        loader: 'json'
      },
      {
        test: /\.html$/,
        loader: 'html',
        query: {
          minimize: process.env.NODE_ENV === 'production',
          removeAttributeQuotes: false,
        }
      }
    ]
  },
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.common.js',
    }
  },
  devServer: {
    historyApiFallback: true,
    noInfo: true
  },
  devtool: '#eval-source-map',
  plugins: [
    new ExtractTextPlugin('style.css'),
  ]
}

var NodeSourcePlugin = require('webpack/lib/node/NodeSourcePlugin')
serverModule = Object.assign({}, client.module)
delete (serverModule.rules[1].exclude)

var binding = process.binding
process.binding = function (name) {
  if (name === 'natives') return {}
  return binding.apply(process, arguments)
}

var server = Object.assign({}, client, {
  entry: './src/server.js',
  target: 'node',
  devtool: false,
  output: Object.assign({}, client.output, {
    filename: 'server-build.js',
    libraryTarget: 'commonjs2'
  }),
  plugins: (client.plugins || []).concat([
    new NodeSourcePlugin(
      {
        console: true,
        process: false,
        global: true,
        Buffer: true,
        setImmediate: true,
        module: 'empty',
        __filename: 'mock',
        __dirname: 'mock'
      }),
  ]),
  module: serverModule,
})

if (process.env.NODE_ENV === 'production') {
  client.devtool = '#source-map'
  // http://vue-loader.vuejs.org/en/workflow/production.html
  client.plugins = (client.plugins || []).concat([
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: '"production"'
      }
    }),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      }
    }),
    new webpack.LoaderOptionsPlugin({
      minimize: true
    })
  ])

  server.plugins = (server.plugins || []).concat([
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: false
      }
    }),
  ])
}

module.exports = [client, server]
