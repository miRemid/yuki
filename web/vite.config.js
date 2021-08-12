import { defineConfig } from 'vite'
import reactRefresh from '@vitejs/plugin-react-refresh'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [reactRefresh()],
  resolve: {
   alias: {
     "@": path.resolve(__dirname, "./src"),
     "~": path.resolve(__dirname, "./"),
   } 
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://192.168.1.106:8080/api',
        changeOrigin: true,
        rewrite: path => path.replace(/^\/api/, '')
      }
    }
  },
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true,
      }
    }
  }
})
