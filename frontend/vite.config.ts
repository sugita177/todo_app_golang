import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // Dockerなどの外部からの接続を許可
    port: 5173,
    watch: {
      usePolling: true, // Windows環境などでホットリロードが効かない場合に有効
    }
  }
})