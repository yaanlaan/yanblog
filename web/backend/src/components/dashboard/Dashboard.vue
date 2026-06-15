<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <StatCard type="users" :number="stats.users" label="用户数" />
      </el-col>
      <el-col :span="6">
        <StatCard type="articles" :number="stats.articles" label="文章数" />
      </el-col>
      <el-col :span="6">
        <StatCard type="categories" :number="stats.categories" label="分类数" />
      </el-col>
      <el-col :span="6">
        <StatCard type="tags" :number="stats.tags" label="标签数" />
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-row :gutter="20" class="quick-actions-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/article/add')" :icon="Edit">
              新建文章
            </el-button>
            <el-button type="success" plain :icon="Upload">
              ZIP 发布
            </el-button>
            <el-button type="info" plain @click="$router.push('/media')" :icon="Picture">
              媒体库
            </el-button>
            <el-button type="warning" plain @click="openFrontend" :icon="View">
              查看前台
            </el-button>
            <el-button type="danger" plain @click="$router.push('/system/status')" :icon="Monitor">
              系统监控
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表和最近文章 -->
    <el-row :gutter="20" class="content-row">
      <el-col :span="14">
        <DataChart ref="chartRef" />
      </el-col>
      <el-col :span="10">
        <el-card class="recent-articles-card">
          <template #header>
            <div class="card-header">
              <span>最近文章</span>
              <el-button type="primary" link @click="$router.push('/article')">查看全部</el-button>
            </div>
          </template>
          <div v-loading="loadingArticles" class="recent-list">
            <div
              v-for="article in recentArticles"
              :key="article.id"
              class="recent-item"
              @click="$router.push(`/article/edit/${article.id}`)"
            >
              <div class="item-info">
                <span class="item-title">{{ article.title }}</span>
                <span class="item-meta">
                  <el-tag size="small" type="info">{{ article.categoryName }}</el-tag>
                  <span class="item-date">{{ formatDate(article.createdAt) }}</span>
                </span>
              </div>
              <el-icon class="item-arrow"><ArrowRight /></el-icon>
            </div>
            <el-empty v-if="!loadingArticles && recentArticles.length === 0" description="暂无文章" :image-size="60" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 系统概览 -->
    <el-row :gutter="20" class="system-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>系统概览</span>
              <el-button :icon="Refresh" circle size="small" @click="fetchSystemStatus" :loading="loadingSystem" />
            </div>
          </template>
          <el-row :gutter="20" v-loading="loadingSystem && !sysStatus.status">
            <el-col :span="6">
              <div class="sys-item">
                <span class="sys-label">服务状态</span>
                <el-tag :type="sysStatus.status === 'online' ? 'success' : 'danger'" size="small">
                  {{ sysStatus.status === 'online' ? '运行中' : '离线' }}
                </el-tag>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="sys-item">
                <span class="sys-label">运行时长</span>
                <span class="sys-value">{{ sysStatus.uptime || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="3">
              <div class="sys-item">
                <span class="sys-label">CPU</span>
                <span class="sys-value" :style="{ color: getUsageColor(sysStatus.cpu_usage) }">
                  {{ sysStatus.cpu_usage ? sysStatus.cpu_usage.toFixed(1) + '%' : '-' }}
                </span>
              </div>
            </el-col>
            <el-col :span="3">
              <div class="sys-item">
                <span class="sys-label">内存</span>
                <span class="sys-value" :style="{ color: getUsageColor(sysStatus.memory_usage) }">
                  {{ sysStatus.memory_usage ? sysStatus.memory_usage.toFixed(1) + '%' : '-' }}
                </span>
              </div>
            </el-col>
            <el-col :span="3">
              <div class="sys-item">
                <span class="sys-label">磁盘</span>
                <span class="sys-value" :style="{ color: getUsageColor(sysStatus.disk_usage) }">
                  {{ sysStatus.disk_usage ? sysStatus.disk_usage.toFixed(1) + '%' : '-' }}
                </span>
              </div>
            </el-col>
            <el-col :span="3">
              <div class="sys-item">
                <span class="sys-label">Goroutines</span>
                <span class="sys-value">{{ sysStatus.goroutines ?? '-' }}</span>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { userApi, articleApi, categoryApi, tagApi, systemApi } from '@/services/api'
import { ElMessage } from 'element-plus'
import { Edit, Upload, Picture, View, Monitor, Refresh, ArrowRight } from '@element-plus/icons-vue'
import StatCard from './StatCard.vue'
import DataChart from './DataChart.vue'

// 统计数据
const stats = reactive({
  users: 0,
  articles: 0,
  categories: 0,
  tags: 0
})

// 图表引用
const chartRef = ref<InstanceType<typeof DataChart> | null>(null)

// 最近文章
const recentArticles = ref<any[]>([])
const loadingArticles = ref(false)

// 系统状态
const sysStatus = reactive({
  status: '',
  uptime: '',
  memory_usage: 0,
  cpu_usage: 0,
  disk_usage: 0,
  goroutines: 0
})
const loadingSystem = ref(false)

// 获取统计数据
const getStats = async () => {
  try {
    const [userRes, articleRes, categoryRes, tagRes] = await Promise.all([
      userApi.getUsers({ pagesize: -1, pagenum: -1 }),
      articleApi.getArticles({ pagesize: -1, pagenum: -1 }),
      categoryApi.getCategories({ pagesize: -1, pagenum: -1 }),
      tagApi.getTags({ pagesize: -1, pagenum: -1 })
    ])
    stats.users = userRes.data.total || 0
    stats.articles = articleRes.data.total || 0
    stats.categories = categoryRes.data.total || 0
    stats.tags = tagRes.data.total || 0

    // 处理最近文章
    const articles = articleRes.data.data || []
    recentArticles.value = articles.slice(0, 5).map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryName: item.Category?.name || item.category?.name || '未分类',
      createdAt: item.CreatedAt || item.created_at
    }))

    // 更新图表
    if (chartRef.value) {
      chartRef.value.updateChart()
    }
  } catch (error) {
    ElMessage.error('获取统计数据失败')
    console.error(error)
  }
}

// 获取系统状态
const fetchSystemStatus = async () => {
  loadingSystem.value = true
  try {
    const res = await systemApi.getSystemStatus()
    if (res.data.status === 200) {
      Object.assign(sysStatus, res.data.data)
    }
  } catch (e) {
    console.error('获取系统状态失败', e)
  } finally {
    loadingSystem.value = false
  }
}

// 打开前台
const openFrontend = () => {
  window.open('/', '_blank')
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
}

// 使用率颜色
const getUsageColor = (val: number) => {
  if (!val) return '#909399'
  if (val < 50) return '#67c23a'
  if (val < 80) return '#e6a23c'
  return '#f56c6c'
}

// 窗口大小变化时重置图表
const handleResize = () => {
  if (chartRef.value) {
    chartRef.value.resizeChart()
  }
}

// 组件挂载时获取数据
onMounted(() => {
  loadingArticles.value = true
  getStats().finally(() => { loadingArticles.value = false })
  fetchSystemStatus()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.quick-actions-row {
  margin-bottom: 20px;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.content-row {
  margin-bottom: 20px;
}

.system-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.recent-articles-card {
  height: 100%;
}

.recent-list {
  min-height: 200px;
}

.recent-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 8px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background 0.2s;
  border-radius: 4px;
}

.recent-item:last-child {
  border-bottom: none;
}

.recent-item:hover {
  background: #f5f7fa;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.item-title {
  font-size: 14px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-date {
  font-size: 12px;
  color: #909399;
}

.item-arrow {
  color: #c0c4cc;
  flex-shrink: 0;
}

.sys-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px;
}

.sys-label {
  font-size: 13px;
  color: #909399;
}

.sys-value {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}
</style>
