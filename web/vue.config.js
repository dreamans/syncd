module.exports = {
    devServer: {
        proxy: {
            '/api': {
                target: 'http://localhost:8878/',
                changeOrigin: true,
            }
        }
    },
    publicPath: '/static/'
}
