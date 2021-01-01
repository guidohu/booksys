const path = require('path');

module.exports = {
  mode: 'production',
  plugins: [
    new BundleAnalyzerPlugin()
  ]
}