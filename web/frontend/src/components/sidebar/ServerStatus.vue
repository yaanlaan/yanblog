<template>
  <div class="sidebar-card server-status">
    <div class="card-header">
      <h3>服务器状态</h3>
    </div>
    <div class="card-content">
      <div v-if="loading" class="loading-placeholder">
        <p>状态加载中...</p>
      </div>
      <div v-else-if="!error">
        <div class="status-item">
          <span class="label">状态:</span>
          <span class="value" :class="serverStatus.status">
            {{ serverStatus.status === 'online' ? '在线' : '离线' }}
          </span>
        </div>
        <div class="status-item">
          <span class="label">运行时间:</span>
          <span class="value">{{ serverStatus.uptime }}</span>
        </div>
        <div class="status-item">
          <span class="label">内存使用:</span>
          <span class="value">{{ serverStatus.memoryUsage }}%</span>
        </div>
        <div class="status-item">
          <span class="label">CPU使用:</span>
          <span class="value">{{ serverStatus.cpuUsage }}%</span>
        </div>
      </div>
      <div class="error-message" v-else>
        <p>❌ {{ error }}</p>
        <button @click="onRetry" class="retry-button">重试</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { systemApi } from '@/services/api'

// 定义服务器状态接口
interface ServerStatus {
  status: 'online' | 'offline'
  uptime: string
  memoryUsage: number
  cpuUsage: number
  startTime: number
}

const serverStatus = ref<ServerStatus>({
  status: 'offline',
  uptime: '未知',
  memoryUsage: 0,
  cpuUsage: 0,
  startTime: 0
})

const loading = ref(false)
const error = ref('')

// 定义事件
const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

// 定时器引用
let serverStatusTimer: number | null = null
let uptimeTimer: number | null = null

// 获取服务器状态
const fetchServerStatus = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    const response = await systemApi.getSystemStatus()
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      error.value = response.data.message || '获取服务器状态失败'
      console.error('获取服务器状态失败:', response.data.message)
      return
    }
    
    // 设置服务器状态数据
    serverStatus.value = {
      status: data.status,
      uptime: data.uptime,
      memoryUsage: Math.round(data.memory_usage * 100) / 100, // 保留两位小数
      cpuUsage: Math.round(data.cpu_usage * 100) / 100, // 保留两位小数
      startTime: data.start_time // 使用后端返回的启动时间戳
    }
  } catch (err: any) {
    error.value = err.message || '获取服务器状态失败'
    console.error('获取服务器状态失败:', err)
    // 即使获取失败，也保持在线状态
    serverStatus.value.status = 'online'
  } finally {
    loading.value = false
    emit('loading', false)
  }
}

// 计算实时运行时间
const calculateRealTimeUptime = () => {
  if (serverStatus.value.startTime <= 0) return '未知'
  
  const elapsed = Date.now() - serverStatus.value.startTime
  return formatUptime(elapsed)
}

// 格式化运行时间
const formatUptime = (milliseconds: number): string => {
  const seconds = Math.floor(milliseconds / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  const remainingSeconds = seconds % 60
  const remainingMinutes = minutes % 60
  const remainingHours = hours % 24
  
  if (days > 0) {
    return `${days}天${remainingHours}小时${remainingMinutes}分钟${remainingSeconds}秒`
  } else if (hours > 0) {
    return `${remainingHours}小时${remainingMinutes}分钟${remainingSeconds}秒`
  } else if (minutes > 0) {
    return `${remainingMinutes}分钟${remainingSeconds}秒`
  } else {
    return `${remainingSeconds}秒`
  }
}

// 重试函数
const onRetry = () => {
  fetchServerStatus()
}

// 暴露方法给父组件
defineExpose({
  fetchServerStatus
})

// 组件挂载时获取数据
onMounted(() => {
  fetchServerStatus()
  
  // 定期更新服务器状态（CPU、内存等）
  serverStatusTimer = window.setInterval(fetchServerStatus, 30000)
  
  // 每秒更新运行时间显示
  uptimeTimer = window.setInterval(() => {
    if (serverStatus.value.startTime > 0) {
      serverStatus.value.uptime = calculateRealTimeUptime()
    }
  }, 1000)
})

// 组件卸载时清理定时器
onBeforeUnmount(() => {
  if (serverStatusTimer) {
    clearInterval(serverStatusTimer)
    serverStatusTimer = null
  }
  if (uptimeTimer) {
    clearInterval(uptimeTimer)
    uptimeTimer = null
  }
})
</script>

<style scoped>
.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  background: #f8f9fa;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.status-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  padding: 15px;
  border-radius: 10px;
  background-color: #f8f9fa;
  transition: all 0.3s ease;
}

.status-item:hover {
  background-color: #e9ecef;
  border-radius: 12px;
}

.status-item:last-child {
  margin-bottom: 0;
}

.status-item .label {
  font-size: 14px;
  color: #888;
}

.status-item .value {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.status-item .value.online {
  color: #28a745;
}

.status-item .value.offline {
  color: #dc3545;
}

.loading-placeholder {
  text-align: center;
  padding: 30px 10px;
  color: #888;
}

.error-message {
  text-align: center;
  padding: 30px 10px;
  color: #dc3545;
}

.retry-button {
  margin-top: 15px;
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.retry-button:hover {
  background-color: #0056b3;
  border-radius: 8px;
}
</style>