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
        // Matches --bk-deep, the top of the mobile canvas, so the splash
        // and the task switcher do not flash white before the app paints.
        theme_color: "#0a1c27",
        background_color: "#0a1c27",
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
