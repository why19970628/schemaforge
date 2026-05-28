import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      '/api': 'http://127.0.0.1:8989',
    },
  },
  build: {
    outDir: '../internal/server/web/dist',
    emptyOutDir: true,
  },
})
