this.global = global || this
global.process = {}
global.process.env = {}
global.process.nextTick = function (fn) {
    setTimeout(fn, 0)
}