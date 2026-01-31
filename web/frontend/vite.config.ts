import { fileURLToPath, URL } from 'node:url'
import fs from 'node:fs'
import yaml from 'js-yaml'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// 读取 public/config.yaml 中的 allowed_hosts
let allowedHosts: string[] = []
try {
  const configPath = fileURLToPath(new URL('./public/config.yaml', import.meta.url))
  const configFile = fs.readFileSync(configPath, 'utf8')
  const config = yaml.load(configFile) as any
  if (config && config.allowed_hosts && Array.isArray(config.allowed_hosts)) {
    allowedHosts = [...allowedHosts, ...config.allowed_hosts]
    // 去重
    allowedHosts = [...new Set(allowedHosts)]
  }
} catch (e) {
  console.warn('Failed to load allowed_hosts from config.yaml:', e)
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
    host: '127.0.0.1',
    port: 3002,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, '/api')
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    },
    // 增加超时时间
    // timeout: 30000
  }
})