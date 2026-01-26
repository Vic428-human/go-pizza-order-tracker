import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
// 例如執行後端這個 endpoint http://localhost:8080/api/admin/dashboard
// 定義在 cmd\routes.go，可以看到後端回傳的訊息
export default defineConfig({
  plugins: [react()],
  server: { proxy: { "/api": "http://localhost:8080" } }
})
