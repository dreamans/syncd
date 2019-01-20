module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:8868/',
                changeOrigin: true,
            }
        }
    }
}
