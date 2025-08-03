<!-- Sidebar 组件：展示侧边栏内容，包含天气信息、置顶文章、标签云和服务器状态四个模块 -->
<template>
  <div class="sidebar">
    <!-- 天气卡片 -->
    <div class="sidebar-card weather-card">
      <div class="card-header">
        <h3>天气信息</h3>
      </div>
      <div class="card-content">
        <div class="weather-info" v-if="weather">
          <div class="weather-main">
            <div class="temperature">{{ weather.temperature }}°C</div>
            <div class="weather-description">{{ weather.description }}</div>
          </div>
          <div class="weather-details">
            <div class="detail-item">
              <span class="label">湿度:</span>
              <span class="value">{{ weather.humidity }}%</span>
            </div>
            <div class="detail-item">
              <span class="label">风速:</span>
              <span class="value">{{ weather.windSpeed }} m/s</span>
            </div>
          </div>
        </div>
        <div class="weather-placeholder" v-else>
          <p>天气信息加载中...</p>
        </div>
      </div>
    </div>

    <!-- 置顶文章 -->
    <div class="sidebar-card featured-articles">
      <div class="card-header">
        <h3>置顶文章</h3>
      </div>
      <div class="card-content">
        <div class="article-list" v-if="featuredArticles.length > 0">
          <div 
            v-for="article in featuredArticles" 
            :key="article.id" 
            class="article-item"
          >
            <router-link :to="`/article/${article.id}`" class="article-link">
              <div class="article-title">{{ article.title }}</div>
              <div class="article-date">{{ formatDate(article.createdAt) }}</div>
            </router-link>
          </div>
        </div>
        <div class="empty-state" v-else>
          <p>暂无置顶文章</p>
        </div>
      </div>
    </div>

    <!-- 标签云 -->
    <div class="sidebar-card tag-cloud">
      <div class="card-header">
        <h3>标签云</h3>
      </div>
      <div class="card-content">
        <div class="tags" v-if="tags.length > 0">
          <span 
            v-for="tag in tags" 
            :key="tag.id" 
            class="tag"
            :style="{ fontSize: calculateFontSize(tag.count) }"
          >
            {{ tag.name }}
          </span>
        </div>
        <div class="empty-state" v-else>
          <p>暂无标签</p>
        </div>
      </div>
    </div>

    <!-- 服务器状态 -->
    <div class="sidebar-card server-status">
      <div class="card-header">
        <h3>服务器状态</h3>
      </div>
      <div class="card-content">
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

// 定义类型
interface Weather {
  temperature: number
  description: string
  humidity: number
  windSpeed: number
}

interface Article {
  id: number
  title: string
  createdAt: string
}

interface Tag {
  id: number
  name: string
  count: number
}

interface ServerStatus {
  status: 'online' | 'offline'
  uptime: string
  memoryUsage: number
  cpuUsage: number
}

// 响应式数据
const weather = ref<Weather | null>(null)
const featuredArticles = ref<Article[]>([])
const tags = ref<Tag[]>([])
const serverStatus = ref<ServerStatus>({
  status: 'online',
  uptime: '0天0小时0分钟',
  memoryUsage: 0,
  cpuUsage: 0
})

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 计算标签字体大小
const calculateFontSize = (count: number) => {
  // 假设标签计数在1-100之间
  const minSize = 12
  const maxSize = 24
  const size = minSize + (count / 100) * (maxSize - minSize)
  return `${size}px`
}

// 模拟获取天气信息
const fetchWeather = () => {
  // 模拟API调用
  setTimeout(() => {
    weather.value = {
      temperature: 22,
      description: '晴',
      humidity: 65,
      windSpeed: 3.5
    }
  }, 500)
}

// 模拟获取置顶文章
const fetchFeaturedArticles = () => {
  // 模拟API调用
  setTimeout(() => {
    featuredArticles.value = [
      { id: 1, title: 'Go语言并发编程指南', createdAt: '2025-07-20T10:00:00Z' },
      { id: 2, title: 'Vue 3状态管理最佳实践', createdAt: '2025-07-15T14:30:00Z' },
      { id: 3, title: '数据库设计优化技巧', createdAt: '2025-07-10T09:15:00Z' }
    ]
  }, 800)
}

// 模拟获取标签数据
const fetchTags = () => {
  // 模拟API调用
  setTimeout(() => {
    tags.value = [
      { id: 1, name: 'Go', count: 80 },
      { id: 2, name: 'Vue', count: 70 },
      { id: 3, name: '数据库', count: 60 },
      { id: 4, name: 'Docker', count: 50 },
      { id: 5, name: 'Kubernetes', count: 40 },
      { id: 6, name: '微服务', count: 45 },
      { id: 7, name: '测试', count: 30 },
      { id: 8, name: '部署', count: 35 }
    ]
  }, 600)
}

// 模拟获取服务器状态
const fetchServerStatus = () => {
  // 模拟API调用
  setTimeout(() => {
    serverStatus.value = {
      status: 'online',
      uptime: '15天6小时32分钟',
      memoryUsage: 65,
      cpuUsage: 28
    }
  }, 300)
}

// 组件挂载时获取数据
onMounted(() => {
  fetchWeather()
  fetchFeaturedArticles()
  fetchTags()
  fetchServerStatus()
  
  // 定期更新服务器状态
  setInterval(fetchServerStatus, 30000)
})
</script>

<style scoped>
.sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

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

.card-content {
  padding: 20px;
}

/* ==================== 天气卡片样式 ==================== */
.weather-info {
  text-align: center;
}

.weather-main {
  margin-bottom: 15px;
}

.temperature {
  font-size: 32px;
  font-weight: bold;
  color: #333;
}

.weather-description {
  font-size: 16px;
  color: #666;
  margin-top: 5px;
}

.weather-details {
  display: flex;
  justify-content: space-around;
}

.detail-item {
  text-align: center;
}

.detail-item .label {
  display: block;
  font-size: 14px;
  color: #888;
}

.detail-item .value {
  display: block;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.weather-placeholder {
  text-align: center;
  padding: 20px 0;
  color: #888;
}

/* ==================== 置顶文章样式 ==================== */
.article-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.article-item {
  border-bottom: 1px solid #eee;
  padding-bottom: 15px;
}

.article-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.article-link {
  text-decoration: none;
  color: inherit;
  display: block;
  transition: color 0.3s;
}

.article-link:hover {
  color: #007bff;
}

.article-title {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 5px;
  line-height: 1.4;
}

.article-date {
  font-size: 13px;
  color: #888;
}

/* ==================== 标签云样式 ==================== */
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag {
  display: inline-block;
  padding: 5px 10px;
  background: #e9ecef;
  border-radius: 20px;
  color: #495057;
  cursor: pointer;
  transition: all 0.3s;
}

.tag:hover {
  background: #007bff;
  color: white;
  transform: scale(1.05);
}

/* ==================== 服务器状态样式 ==================== */
.status-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
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

.empty-state {
  text-align: center;
  padding: 20px 0;
  color: #888;
}

/* ==================== 响应式样式 ==================== */
@media (max-width: 992px) {
  .sidebar {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .sidebar-card {
    flex: 1 1 calc(50% - 10px);
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    flex-direction: column;
  }
  
  .sidebar-card {
    flex: 1 1 100%;
  }
}
</style>