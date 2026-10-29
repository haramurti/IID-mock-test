import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    // In dev, the browser talks only to Vite (:5173). Vite forwards /api/* to
    // the Go backend, so both look like the same origin and no CORS is needed.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
