const path = require('path');
// const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  // assetsDir: 'assets',
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
  configureWebpack: {
    mode: "production",
  //   entry: "./src/main.js",
  //   output: {
  //     filename: 'bundle.js'
  //   },
  //   plugins: [new HtmlWebpackPlugin({
  //     template: './public/index.html',
  //     title: "Booking System"
  //   })]
    // plugins: [
    //   new BundleAnalyzerPlugin()
    // ],
    // plugins: [
    //   new MiniCssExtractPlugin()
    // ],
    // module: {
    //   rules: [
    //     {
    //         test: /\.css$/,
    //         use: ExtractTextPlugin.extract({
    //             use: ['css-loader'],
    //             fallback: 'style-loader'
    //         })
    //     },
    //     {
    //         test: /\.scss$/,
    //         use: ExtractTextPlugin.extract({
    //             use: ['css-loader', 'sass-loader' ],
    //             fallback: 'style-loader'
    //         })
    //     }
    //   ]
    // },
    // plugins: [
    //     new ExtractTextPlugin(cssOutput),
    //     new PurgecssPlugin({
    //         paths: glob.sync('./Views/**/*.cshtml', { nodir: true }),
    //         whitelistPatterns: [ /selectize-.*/ ]
    //     })
    // ],
    // module: {
    //   rules: [
    //     {
    //       test: /\.css$/i,
    //       use: [MiniCssExtractPlugin.loader, 'css-loader']
    //     }
    //   ]
    // },
    optimization: {
      splitChunks: {
        chunks: 'all'
  //     minSize: 10000,
  //     maxSize: 250000
        // cacheGroups: {
        //   vendor: {
        //     name: "node_vendors", // part of the bundle name and
        //     // can be used in chunks array of HtmlWebpackPlugin
        //     test: /[\\/]node_modules[\\/]/,
        //     chunks: "initial",
        //   },
        //   common: {
        //     test: /[\\/]src[\\/]components[\\/]/,
        //     chunks: "all",
        //     minSize: 0,
        //   },
        // }
      },
      // minimize: true,
      // minimizer: [
      //   new CssMinimizerPlugin()
      // ]
    }
  },
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