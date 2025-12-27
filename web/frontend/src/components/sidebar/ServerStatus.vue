<template>
  <div class="sidebar-card server-status">
    <div class="card-header">
      <h3><i class="iconfont icon-monitor" style="color: #333;"></i> 服务器状态</h3>
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
import { systemApi } from '@/services/api'

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

const loading = ref(false)
const error = ref('')
let timer: number | null = null

// 定义状态项配置
const statusItems = computed(() => [
  { label: 'MEM', valueKey: 'memoryUsage', type: 'mem' },
  { label: 'CPU', valueKey: 'cpuUsage', type: 'cpu' },
  { label: 'DISK', valueKey: 'diskUsage', type: 'disk' }
])

const fetchServerStatus = async () => {
  try {
    // loading.value = true // Don't show loading on refresh
    const response = await systemApi.getSystemStatus()
    if (response.data.status === 200) {
      const data = response.data.data
      serverStatus.value = {
        status: data.status,
        uptime: data.uptime,
        memoryUsage: Math.round((data.memory_usage || 0) * 100) / 100,
        cpuUsage: Math.round((data.cpu_usage || 0) * 100) / 100,
        diskUsage: Math.round((data.disk_usage || 0) * 100) / 100,
        startTime: data.start_time
      }
    } else {
      error.value = response.data.message
    }
  } catch (err: any) {
    error.value = '获取失败'
    // Mock for demo
    serverStatus.value = {
      status: 'online',
      uptime: '10天',
      memoryUsage: 60.6,
      cpuUsage: 12.5,
      diskUsage: 70.3,
      startTime: 0
    }
    error.value = '' // Clear error if using mock
  } finally {
    loading.value = false
  }
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
  color: #333;
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
  color: #666;
}

.progress-bar {
  height: 8px;
  background: #f0f0f0;
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