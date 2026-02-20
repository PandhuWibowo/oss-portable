import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    fs: {
      allow: ['..'] // allow importing docs from project root
    },
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
