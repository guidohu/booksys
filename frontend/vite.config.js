
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { createHtmlPlugin } from 'vite-plugin-html'
import path from 'path';
import { visualizer } from 'rollup-plugin-visualizer'; // Import the visualizer plugin

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    createHtmlPlugin({
      inject: {
        data: {
          title: 'Booksys',
        },
      },
    }),
    // Add the visualizer plugin
    visualizer({
      open: true,
      gzipSize: true,
      brotliSize: true,
      filename: 'bundle-analysis.html',
      template: 'treemap',
    }),

  ],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:9090',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      // Map 'assets' to the absolute path of your assets directory
      'api': path.resolve(__dirname, './src/api'),
      'assets': path.resolve(__dirname, './src/assets'),
      'libs': path.resolve(__dirname, './src/libs'),
      'store': path.resolve(__dirname, './src/store'),
      'views': path.resolve(__dirname, './src/views'),
      'booksys': path.resolve(__dirname, './src'),
      // If you also use an '@' alias for 'src', add it here too:
      // '@': path.resolve(__dirname, './src'),
    },
  },
})
