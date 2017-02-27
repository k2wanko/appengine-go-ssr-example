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

function run(context) {
    router.push(context.url)
    const matchedComponents = router.getMatchedComponents()
    if (!matchedComponents.length) {
        context.code = 404
    }

    return Promise.all(matchedComponents.map(component => {
        if (component.preFetch) {
            return component.preFetch(store)
        }
    })).then(() => {
        context.initialState = store.state
        return app
    })
}

module.exports = function (context) {
    const {url, res} = context

    let code = 200

    router.push(url)
    const matchedComponents = router.getMatchedComponents()
    if (!matchedComponents.length) {
        code = 404
    }

    const stream = renderer.renderToStream(app)
    stream.once('data', () => {
        res.write(html.head)
    })
    stream.on('data', chunk => {
        res.write(chunk)
    })
    stream.on('end', () => {
        context.code = code
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