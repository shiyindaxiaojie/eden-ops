import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        cookieDomainRewrite: 'localhost',
        configure: (proxy, options) => {
          // 代理事件处理器，打印请求和响应
          proxy.on('proxyReq', function(proxyReq, req, res) {
            console.log('发送到服务器的请求:', {
              path: req.url,
              headers: proxyReq.getHeaders()
            });
          });
          
          proxy.on('proxyRes', function(proxyRes, req, res) {
            console.log('收到服务器的响应:', {
              path: req.url,
              statusCode: proxyRes.statusCode,
              headers: proxyRes.headers
            });
          });
        }
      }
    },
  },
  define: {
    'import.meta.env.VITE_APP_BASE_API': '""',
  },
})