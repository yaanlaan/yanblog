import { fileURLToPath, URL } from 'node:url'
import fs from 'node:fs'
import yaml from 'js-yaml'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// 读取前端目录下的 public/config.yaml 中的 allowed_hosts
let allowedHosts: string[] = []
try {
  // 假设后端和前端在同一级目录下 (web/backend 和 web/frontend)
  const configPath = fileURLToPath(new URL('../frontend/public/config.yaml', import.meta.url))
  if (fs.existsSync(configPath)) {
    const configFile = fs.readFileSync(configPath, 'utf8')
    const config = yaml.load(configFile) as any
    if (config && config.allowed_hosts && Array.isArray(config.allowed_hosts)) {
       allowedHosts = [...allowedHosts, ...config.allowed_hosts]
    }
    // 同时也尝试添加 admin_url 中的域名
    if (config && config.admin_url) {
       try {
         const adminHostname = new URL(config.admin_url).hostname
         allowedHosts.push(adminHostname)
       } catch (e) {}
    }
    allowedHosts = [...new Set(allowedHosts)]
  }
} catch (e) {
  console.warn('Failed to load allowed_hosts from frontend config:', e)
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    allowedHosts: allowedHosts,
    host: '0.0.0.0', // 强制监听 IPv4
    port: 3001,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        // 不再重写路径，保持/api前缀
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/assets': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/static': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/iconfont': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },    }
  }
})