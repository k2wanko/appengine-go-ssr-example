const http = require('http')
const fs = require('fs')

const { createBundleRenderer } = require('vue-server-renderer')
const serverBundle = require('./vue-ssr-server-bundle')
const clientManifest = require('./app/vue-ssr-client-manifest.json')
const template = fs.readFileSync('./index.html', 'utf-8')

const renderer = createBundleRenderer(serverBundle, {
  runInNewContext: false,
  template,
  clientManifest,
  basedir: './app',
})


const server = http.createServer((req, res) => {
  res.setHeader("Content-Type", "text/html")

  const context = {
    url: req.url
  }

  renderer.renderToString(context, (err, html) => {
    const handleError = err => {
      res.status(500).end('500 | Internal Server Error')
      console.error(`error during render : ${req.url}`)
      console.error(err.stack)
    }
    if (err) {
      return handleError(err)
    }
    res.end(html)
  })
})
server.on('clientError', (err, socket) => {
  socket.end('HTTP/1.1 400 Bad Request\r\n\r\n')
})
server.listen(8000)