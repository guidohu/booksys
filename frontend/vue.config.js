const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  // assetsDir: './src/assets',
  outputDir: 'dist',
  publicPath: '/',
  // indexPath: 'index.html',
  // pages: {
  //   //page's entry file
  //   entry: 'src/main.js',
  //   //template file
  //   template: 'public/index.html',
  //   //Output file in dist/index.html
  //   filename: 'index.html',
  //   // When using the page title option,
  //   // The title tag in template needs to be <title><%= htmlWebpackPlugin.options.title %></title>
  //   title: 'Booking System',
  //   // The blocks contained in this page, by default, contain
  //   // Extracted general chunks and vendor chunk s.
  //   chunks: ['chunk-vendos', 'chunk-common', 'index'],
  // },
  // configureWebpack: {
  //   entry: "./src/main.js",
  //   output: {
  //     filename: 'bundle.js'
  //   },
  //   plugins: [new HtmlWebpackPlugin({
  //     template: './public/index.html',
  //     title: "Booking System"
  //   })]
  // }
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].template = './public/index.html';
        args[0].title = "Booking System";
        return args
      });
    config
      .plugin('pwa')
      .tap(args => {
        args[0].workboxPluginMode = "GenerateSW";
        args[0].manifestOptions = {};
        args[0].name = "Booking System";
        args[0].themeColor = "aliceblue";
        args[0].iconPaths = {
          favicon32: 'img/icons/favicon-32x32.png',
          favicon16: 'img/icons/favicon-16x16.png',
          appleTouchIcon: 'img/icons/apple-touch-icon-512x512.png',
          maskIcon: 'img/icons/safari-pinned-tab.svg',
          msTileIMage: 'img/icons/tile150x150.png'
        };
        return args
      });
  }
}