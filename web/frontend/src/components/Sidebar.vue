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
        <div class="tags" v-if="categories.length > 0">
          <router-link
            v-for="category in categories" 
            :key="category.id" 
            :to="`/category/${category.id}`"
            class="tag"
            :style="{ fontSize: calculateFontSize(category.articleCount) }"
          >
            {{ category.name }}
          </router-link>
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
import { articleApi, categoryApi } from '@/services/api'

// 类型定义
interface Weather {
  temperature: number
  description: string
  humidity: number
  windSpeed: number
}

interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  createdAt: string
  updatedAt: string
}

interface Tag {
  id: number
  name: string
  count: number
}

interface Category {
  id: number
  name: string
  articleCount: number
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
const categories = ref<Category[]>([])
const serverStatus = ref<ServerStatus>({
  status: 'offline',
  uptime: '未知',
  memoryUsage: 0,
  cpuUsage: 0
})

// 计算字体大小（基于文章数量）
const calculateFontSize = (count: number) => {
  // 基础字体大小12px，最大字体大小24px
  const minSize = 12
  const maxSize = 24
  // 使用对数函数使字体大小变化更平滑
  const size = minSize + (maxSize - minSize) * (Math.log(count + 1) / Math.log(100))
  return `${size}px`
}

// 获取天气信息
const fetchWeather = async () => {
  try {
    // 模拟API调用
    setTimeout(() => {
      weather.value = {
        temperature: 22,
        description: '晴',
        humidity: 65,
        windSpeed: 3.5
      }
    }, 500)
  } catch (error) {
    console.error('获取天气信息失败:', error)
  }
}

// 获取置顶文章
const fetchFeaturedArticles = async () => {
  try {
    const response = await articleApi.getTopArticles({ num: 3 })
    const { data } = response.data
    featuredArticles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
  } catch (error) {
    console.error('获取置顶文章失败:', error)
    // 如果API调用失败，使用模拟数据
    setTimeout(() => {
      featuredArticles.value = [
        { id: 1, title: 'Go语言并发编程指南', categoryId: 1, categoryName: '技术', desc: '', content: '', img: '', createdAt: '2025-07-20T10:00:00Z', updatedAt: '2025-07-20T10:00:00Z' },
        { id: 2, title: 'Vue 3状态管理最佳实践', categoryId: 1, categoryName: '技术', desc: '', content: '', img: '', createdAt: '2025-07-15T14:30:00Z', updatedAt: '2025-07-15T14:30:00Z' },
        { id: 3, title: '数据库设计优化技巧', categoryId: 1, categoryName: '技术', desc: '', content: '', img: '', createdAt: '2025-07-10T09:15:00Z', updatedAt: '2025-07-10T09:15:00Z' }
      ]
    }, 800)
  }
}

// 获取分类列表（用于标签云）
const fetchCategories = async () => {
  try {
    const response = await categoryApi.getCategories({
      pagesize: -1, // 获取所有分类
      pagenum: -1
    })
    
    const { data } = response.data
    // 获取每个分类的文章数量
    const categoriesWithCount = await Promise.all(
      data.map(async (item: any) => {
        return {
          id: item.ID,
          name: item.name,
          articleCount: item.article_count || 0
        }
      })
    )
    
    categories.value = categoriesWithCount
  } catch (error) {
    console.error('获取分类列表失败:', error)
    // 如果API调用失败，使用模拟数据
    setTimeout(() => {
      categories.value = [
        { id: 1, name: '技术', articleCount: 80 },
        { id: 2, name: '生活', articleCount: 70 },
        { id: 3, name: '旅行', articleCount: 60 },
        { id: 4, name: '美食', articleCount: 50 },
        { id: 5, name: '摄影', articleCount: 40 },
        { id: 6, name: '读书', articleCount: 45 },
        { id: 7, name: '电影', articleCount: 30 },
        { id: 8, name: '音乐', articleCount: 35 }
      ]
    }, 600)
  }
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
  fetchCategories()
  fetchServerStatus()
  
  // 定期更新服务器状态
  setInterval(fetchServerStatus, 30000)
})

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}
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