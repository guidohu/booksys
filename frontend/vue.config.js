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
        chunks: 'all',
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
            test: /[\\/]node_modules[\\/]core\-js/,
            reuseExistingChunk: true,
            minSize: 0,
            chunks: "all",
            priority: 10,
          },
          components: {
            test: /[\\/]src[\\/]components[\\/]/,
            chunks: "all",
            minSize: 0,
            reuseExistingChunk: true
          },
          default: {
            chunks: "all",
            minChunks: 2,
            reuseExistingChunk: true
          }
        }
      },
      // minimize: true,
      // minimizer: [
      //   new CssMinimizerPlugin()
      // ]
    }
  },
  // chainWebpack: config => {
  //   config
  //     .plugin('html')
  //     .tap(args => {
  //       args[0].template = './public/index.html';
  //       args[0].title = "Booking System";
  //       return args
  //     });
  // },
  pwa: {
    manifestOptions: {
      name: "Booking System",
      short_name: "Booking System",
      start_url: "./",
      display: "standalone",
      theme_color: "#1accda",
      icons: [
        {
          src: "./img/icons/android-chrome-512x512.png",
          sizes: "512x512",
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
      favicon32: 'img/icons/favicon-32x32.png',
      favicon16: 'img/icons/favicon-16x16.png',
      appleTouchIcon: 'img/icons/apple-touch-icon-512x512.png',
      msTileImage: 'img/icons/tile150x150.png'
    },
    workboxPluginMode: "GenerateSW",
  }
}