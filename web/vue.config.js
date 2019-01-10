let assetsDir = (() => {
    let t = Math.ceil(new Date().getTime()/1000)
    return 'static/' + t
})()

module.exports = {
    baseUrl: process.env.BASE_URL,
    devServer: {
        host: 'localhost',
        port: 8801,
        compress: true
    },
    assetsDir: assetsDir,
    productionSourceMap: false
}
