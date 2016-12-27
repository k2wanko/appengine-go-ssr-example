import { app, router } from './app'
import { createRenderer } from 'vue-server-renderer'

import layout from './index.html'

const html = (() => {
    const target = '<div id="app"></div>'
    const i = layout.indexOf(target)
    return {
        head: layout.slice(0, i),
        tail: layout.slice(i + target.length)
    }
})()

const renderer = createRenderer({
    cache: global.ComponentCache,
})

export function renderStream({ url, res }) {
    router.push(url)

    const stream = renderer.renderToStream(app)
    stream.once('data', () => {
        res.write(html.head)
    })
    stream.on('data', chunk => {
        res.write(chunk)
    })
    stream.on('end', () => {
        res.end(html.tail)
    })
    stream.on('error', err => {
        res.error(err)
    })
}