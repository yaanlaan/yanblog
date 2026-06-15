<template>
  <div class="system-status">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统监控</span>
          <div class="header-actions">
            <el-tag :type="status.status === 'online' ? 'success' : 'danger'" size="large" effect="dark">
              {{ status.status === 'online' ? '在线' : '离线' }}
            </el-tag>
            <el-button :icon="Refresh" circle @click="fetchStatus" :loading="loading" />
          </div>
        </div>
      </template>

      <div v-loading="loading && !status.status">
        <!-- 基本信息 -->
        <el-descriptions :column="2" border class="mb-4">
          <el-descriptions-item label="服务状态">
            <el-tag :type="status.status === 'online' ? 'success' : 'danger'">
              {{ status.status === 'online' ? '运行中' : '已停止' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="运行时长">
            {{ status.uptime || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="Goroutine 数量">
            {{ status.goroutines ?? '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="上次刷新">
            {{ lastUpdate || '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 资源使用率 -->
        <h3 class="section-title">资源使用率</h3>
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="usage-card">
              <div class="usage-label">CPU 使用率</div>
              <el-progress
                :percentage="Math.round(status.cpu_usage || 0)"
                :color="getProgressColor(status.cpu_usage || 0)"
                :stroke-width="20"
                :text-inside="true"
              />
            </div>
          </el-col>
          <el-col :span="8">
            <div class="usage-card">
              <div class="usage-label">内存使用率</div>
              <el-progress
                :percentage="Math.round(status.memory_usage || 0)"
                :color="getProgressColor(status.memory_usage || 0)"
                :stroke-width="20"
                :text-inside="true"
              />
            </div>
          </el-col>
          <el-col :span="8">
            <div class="usage-card">
              <div class="usage-label">磁盘使用率</div>
              <el-progress
                :percentage="Math.round(status.disk_usage || 0)"
                :color="getProgressColor(status.disk_usage || 0)"
                :stroke-width="20"
                :text-inside="true"
              />
            </div>
          </el-col>
        </el-row>

        <!-- 自动刷新设置 -->
        <div class="auto-refresh">
          <el-switch v-model="autoRefresh" @change="toggleAutoRefresh" />
          <span class="auto-refresh-label">自动刷新 (每 5 秒)</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { systemApi } from '@/services/api'

defineOptions({
  name: 'SystemStatus'
})

const loading = ref(false)
const autoRefresh = ref(true)
let timer: ReturnType<typeof setInterval> | null = null
const lastUpdate = ref('')

const status = reactive({
  status: '',
  uptime: '',
  memory_usage: 0,
  cpu_usage: 0,
  disk_usage: 0,
  goroutines: 0,
  start_time: 0
})

const getProgressColor = (percentage: number): string => {
  if (percentage < 50) return '#67c23a'
  if (percentage < 80) return '#e6a23c'
  return '#f56c6c'
}

const fetchStatus = async () => {
  loading.value = true
  try {
    const res = await systemApi.getSystemStatus()
    if (res.data.status === 200) {
      const data = res.data.data
      Object.assign(status, data)
      lastUpdate.value = new Date().toLocaleTimeString('zh-CN')
    }
  } catch (e) {
    ElMessage.error('获取系统状态失败')
    console.error(e)
  } finally {
    loading.value = false
  }
}

const startAutoRefresh = () => {
  if (timer) clearInterval(timer)
  timer = setInterval(fetchStatus, 5000)
}

const stopAutoRefresh = () => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

const toggleAutoRefresh = (val: boolean) => {
  if (val) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

onMounted(() => {
  fetchStatus()
  if (autoRefresh.value) {
    startAutoRefresh()
  }
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.system-status {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mb-4 {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin: 0 0 16px 0;
}

.usage-card {
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  text-align: center;
}

.usage-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 12px;
}

.auto-refresh {
  margin-top: 24px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.auto-refresh-label {
  font-size: 14px;
  color: #606266;
}
</style>
