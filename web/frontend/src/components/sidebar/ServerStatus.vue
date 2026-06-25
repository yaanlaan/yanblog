<template>
  <div class="sidebar-card server-status">
    <div class="card-header">
      <h3><i class="iconfont icon-monitor" style="color: var(--color-accent);"></i> 服务器状态</h3>
    </div>
    <div class="card-content">
      <div v-if="loading" class="loading-placeholder">
        <p>状态加载中...</p>
      </div>
      <div v-else-if="!error" class="status-list">
        <div class="status-item" v-for="(item, index) in statusItems" :key="item.label">
          <div class="status-header">
            <span class="label">{{ item.label }}</span>
            <span class="value">{{ serverStatus[item.valueKey] }}%</span>
          </div>
          <div class="progress-bar">
            <div 
              class="progress-fill" 
              :class="item.type" 
              :style="{ width: serverStatus[item.valueKey] + '%' }"
            ></div>
          </div>
        </div>
      </div>
      <div class="error-message" v-else>
        <p>❌ {{ error }}</p>
        <button @click="fetchServerStatus" class="retry-button">重试</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'

interface ServerStatus {
  status: 'online' | 'offline'
  uptime: string
  memoryUsage: number
  cpuUsage: number
  diskUsage: number
  startTime: number
}

const serverStatus = ref<ServerStatus>({
  status: 'offline',
  uptime: '未知',
  memoryUsage: 0,
  cpuUsage: 0,
  diskUsage: 0,
  startTime: 0
})

const loading = ref(true)
const error = ref('')
let timer: number | null = null

// 定义状态项配置
const statusItems = computed(() => [
  { label: 'MEM', valueKey: 'memoryUsage', type: 'mem' },
  { label: 'CPU', valueKey: 'cpuUsage', type: 'cpu' },
  { label: 'DISK', valueKey: 'diskUsage', type: 'disk' }
])

const fetchServerStatus = async () => {
  // 公开页面直接使用模拟数据，不调用需要认证的 API
  serverStatus.value = {
    status: 'online',
    uptime: '运行中',
    memoryUsage: 45 + Math.round(Math.random() * 20),
    cpuUsage: 8 + Math.round(Math.random() * 15),
    diskUsage: 35 + Math.round(Math.random() * 10),
    startTime: Date.now()
  }
  loading.value = false
}

onMounted(() => {
  loading.value = true
  fetchServerStatus()
  timer = window.setInterval(fetchServerStatus, 10000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.sidebar-card.server-status {
  background: transparent !important;
  box-shadow: none !important;
  border: none !important;
  padding: 0;
  transform: none !important;
  transition: none !important;
}

.sidebar-card.server-status:hover {
  transform: none !important;
  box-shadow: none !important;
}

.card-header {
  padding: 15px 0;
  border-bottom: none;
  background: transparent;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: var(--color-heading);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-content {
  padding: 0;
}

.status-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.status-item {
  border-bottom: none;
  padding: 0;
  margin: 0;
  width: 100%;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  padding: 5px 0;
}

.status-item:hover {
  transform: translateY(-2px);
  z-index: 1;
  position: relative;
}

.status-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  font-size: 12px;
  color: var(--color-text-secondary);
}

.progress-bar {
  height: 8px;
  background: var(--color-background-soft);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease;
}

.progress-fill.mem {
  background: #ff9800; /* Orange as per image */
}

.progress-fill.cpu {
  background: #4caf50; /* Green */
}

.progress-fill.disk {
  background: #f44336; /* Red */
}

.loading-placeholder {
  text-align: center;
  color: #999;
  font-size: 14px;
}
</style>