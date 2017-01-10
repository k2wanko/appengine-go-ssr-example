import { app, router, store } from './app'
import Vue from 'vue'
import { createRenderer } from 'vue-server-renderer'
import serialize from 'serialize-javascript'
import layout from './index.html'

Vue.config.devtools = false

const html = (() => {
    const target = '<div id="app"></div>'
    const i = layout.indexOf(target)
    return {
        head: layout.slice(0, i),
        tail: layout.slice(i + target.length)
    }
})()

const renderer = createRenderer({
    cache: global['__ComponentCache__'],
})

module.exports = function (context) {
    const {url, res} = context
    router.push(url)

    const stream = renderer.renderToStream(app)
    stream.once('data', () => {
        res.write(html.head)
    })
    stream.on('data', chunk => {
        res.write(chunk)
    })
    stream.on('end', () => {
        res.write(
            `<script>window.__INITIAL_STATE__=${
            serialize(store.state, { isJSON: true })
            }</script>`
        )
        res.end(html.tail)
    })
    stream.on('error', err => {
        res.error(err)
    })
}