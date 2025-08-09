import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { createHtmlPlugin } from "vite-plugin-html";
import path from "path";
import { visualizer } from "rollup-plugin-visualizer"; // Import the visualizer plugin
import { VitePWA } from "vite-plugin-pwa";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    createHtmlPlugin({
      inject: {
        data: {
          title: "Booksys",
        },
      },
    }),
    // Add the visualizer plugin
    visualizer({
      open: false,
      gzipSize: true,
      brotliSize: true,
      filename: "bundle-analysis.html",
      template: "treemap",
    }),
    VitePWA({
      manifest: {
        name: "Wake and Surf Booksys",
        short_name: "Booksys",
        description: "Wake and Surf Booking System for boat communities.",
        theme_color: "#ffffff",
        background_color: "#ffffff",
        display: "standalone",
        scope: "/",
        start_url: "/",
      },
      devOptions: {
        enabled: false, // Enable PWA in development
      },
    }),
  ],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:9090",
        changeOrigin: true,
      },
      "^/uploads": {
        target: "http://localhost:9090",
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      // Map 'assets' to the absolute path of your assets directory
      booksys: path.resolve(__dirname, "./src"),
      // If you also use an '@' alias for 'src', add it here too:
      // '@': path.resolve(__dirname, './src'),
    },
  },
});
