

module.exports = function(req, res) {
    res.setHeader("Content-Type", "text/html")
    res.writeHead(200)
    res.write(req.url)
}
