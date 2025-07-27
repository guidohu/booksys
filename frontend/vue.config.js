// const path = require('path');

module.exports = {
  outputDir: "dist",
  publicPath: "/",
  devServer: {
    port: 8080,
    proxy: {
      "^/api/v2": {
        target: "http://localhost:9090",
        changeOrigin: true,
      },
      "^/uploads": {
        target: "http://localhost:9090",
        changeOrigin: true,
      },
    },
  },
  pages: {
    index: {
      title: 'Wake and Surf',
      entry: 'src/main.js',
    }
  },
}
