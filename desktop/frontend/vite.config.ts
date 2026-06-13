import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    wails(path.resolve(__dirname, "./bindings")),
    tailwindcss(),
  ],
  server: {
    host: "127.0.0.1",
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@bindings": path.resolve(__dirname, "./bindings"),
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('apexcharts') || id.includes('vue3-apexcharts')) return 'charts'
          if (id.includes('vuetify')) return 'vuetify'
          if (id.includes('vue') || id.includes('pinia')) return 'vue-vendor'
        }
      }
    }
  }
});
