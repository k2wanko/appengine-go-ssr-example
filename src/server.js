import './polyfill.js'
import CreateApp from './main'
import { createRenderer } from 'vue-server-renderer'

import layout from 'html!./index.html'

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

export function renderStream(res) {
    const vm = CreateApp()
    const stream = renderer.renderToStream(vm)
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