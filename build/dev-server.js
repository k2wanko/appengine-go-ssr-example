const webpack = require('webpack')
const Server = require('webpack-dev-server')

const wpConfig = require('./webpack.config')
const clientConfig = wpConfig[0]
const serverConfig = wpConfig[1]

const options = {
    host: 'localhost',
    port: '8001'
}

clientConfig.entry = [
    `webpack-dev-server/client?http://${options.host}:${options.port}`,
    'webpack/hot/dev-server',
    clientConfig.entry
]

clientConfig.plugins.push(
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoErrorsPlugin()
)

const compiler = webpack([clientConfig, serverConfig])
compiler.apply(new webpack.ProgressPlugin({
    profile: true,
}))

const server = new Server(compiler, {})
server.listen(options.port, options.host, err => {
    if (err) throw err
    console.log(`listen: ${options.host}:${options.port}`)
})

const NodeOutputFileSystem = require('webpack/lib/node/NodeOutputFileSystem')
compiler.outputFileSystem = new NodeOutputFileSystem()
compiler.plugin("done", stats => {
    stats = stats.toJson()
    stats.errors.forEach(err => console.error(err))
    stats.warnings.forEach(err => console.warn(err))
})
