// const path = require('path');

module.exports = {
  outputDir: "dist",
  publicPath: "/",
  devServer: {
    proxy: "http://localhost:80",
  },
  configureWebpack: {
    mode: "production",
    optimization: {
      splitChunks: {
        chunks: "all",
        cacheGroups: {
          vue: {
            test: /[\\/]node_modules[\\/]vue/,
            reuseExistingChunk: true,
            chunks: "all",
          },
          raphael: {
            test: /[\\/]node_modules[\\/]raphael/,
            reuseExistingChunk: true,
            minSize: 0,
            chunks: "all",
            priority: 20,
          },
          popper: {
            test: /[\\/]node_modules[\\/]popper/,
            reuseExistingChunk: true,
            minSize: 0,
            chunks: "all",
            priority: 15,
          },
          coreJs: {
            test: /[\\/]node_modules[\\/]core-js/,
            reuseExistingChunk: true,
            minSize: 0,
            chunks: "all",
            priority: 10,
          },
          components: {
            test: /[\\/]src[\\/]components[\\/]/,
            chunks: "all",
            minSize: 0,
            reuseExistingChunk: true,
          },
          default: {
            chunks: "all",
            minChunks: 2,
            reuseExistingChunk: true,
          },
        },
      },
    },
  },
  pwa: {
    manifestOptions: {
      name: "Booking System",
      short_name: "Booking System",
      start_url: "./",
      display: "standalone",
      theme_color: "#1accda",
      icons: [
        {
          src: "./img/icons/favicon-32x32.png",
          sizes: "32x32",
          type: "image/png",
          purpose: "any maskable",
        },
      ],
    },
    themeColor: "#1accda",
    msTileColor: "#1accda",
    appleMobileWebAppCapable: "no",
    appleMobileWebAppStatusBarStyle: "black",
    iconPaths: {
      maskicon: null,
      favicon32: "img/icons/favicon-32x32.png",
      favicon16: "img/icons/favicon-16x16.png",
      appleTouchIcon: "img/icons/apple-touch-icon-512x512.png",
      msTileImage: "img/icons/tile150x150.png",
    },
    workboxPluginMode: "GenerateSW",
  },
};
