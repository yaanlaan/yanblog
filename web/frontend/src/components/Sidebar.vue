<!-- Sidebar 组件：展示侧边栏内容，包含天气信息、置顶文章、标签云和服务器状态四个模块 -->
<template>
  <div class="sidebar">
    <!-- 天气卡片 -->
    <div class="sidebar-card weather-card">
      <div class="card-header">
        <h3>天气信息</h3>
      </div>
      <div class="card-content">
        <div v-if="loading.weather" class="skeleton-loader">
          <div class="skeleton-header"></div>
          <div class="skeleton-body">
            <div class="skeleton-line"></div>
            <div class="skeleton-line"></div>
          </div>
        </div>
        <div class="weather-info" v-else-if="weather">
          <div class="weather-main">
            <div class="temperature">{{ weather.temperature.toFixed(1) }}°C</div>
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
        <div class="error-message" v-else-if="errors.weather">
          <p>❌ {{ errors.weather }}</p>
          <button @click="fetchWeather" class="retry-button">重试</button>
        </div>
        <div class="weather-placeholder" v-else>
          <p>暂无天气信息</p>
        </div>
      </div>
    </div>

    <!-- 置顶文章 -->
    <div class="sidebar-card featured-articles">
      <div class="card-header">
        <h3>置顶文章</h3>
      </div>
      <div class="card-content">
        <div v-if="loading.articles" class="skeleton-loader">
          <div class="skeleton-header"></div>
          <div class="skeleton-body">
            <div class="skeleton-line"></div>
            <div class="skeleton-line"></div>
          </div>
        </div>
        <div class="article-list" v-else-if="featuredArticles.length > 0">
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
        <div class="error-message" v-else-if="errors.articles">
          <p>❌ {{ errors.articles }}</p>
          <button @click="fetchFeaturedArticles" class="retry-button">重试</button>
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
        <div v-if="loading.categories" class="skeleton-loader">
          <div class="skeleton-header"></div>
          <div class="skeleton-body">
            <div class="skeleton-tag"></div>
            <div class="skeleton-tag"></div>
            <div class="skeleton-tag"></div>
          </div>
        </div>
        <div class="tags" v-else-if="categories.length > 0">
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
        <div class="error-message" v-else-if="errors.categories">
          <p>❌ {{ errors.categories }}</p>
          <button @click="fetchCategories" class="retry-button">重试</button>
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
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { articleApi, categoryApi, weatherApi } from '@/services/api'

// 类型定义
interface Weather {
  city: string
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

// 错误状态
const errors = ref({
  weather: '',
  articles: '',
  categories: ''
})

// 加载状态
const loading = ref({
  weather: false,
  articles: false,
  categories: false
})

// 计算字体大小（基于文章数量）
const calculateFontSize = (count: number) => {
  // 基础字体大小12px，最大字体大小24px
  const minSize = 12
  const maxSize = 24
  // 假设最大文章数为50篇
  const maxCount = 50
  const size = minSize + (maxSize - minSize) * Math.min(count / maxCount, 1)
  return `${size}px`
}

// 获取天气信息
const fetchWeather = async () => {
  try {
    loading.value.weather = true
    errors.value.weather = ''
    console.log('开始获取天气信息...') // 调试日志
    
    const response = await weatherApi.getWeather()
    console.log('天气API响应:', response) // 调试日志
    
    // 统一处理响应数据结构
    if (!response?.data) {
      const errorMessage = '无效的API响应'
      errors.value.weather = errorMessage
      console.error(errorMessage)
      return
    }
    
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      const errorMessage = data?.message || '获取天气信息失败'
      errors.value.weather = errorMessage
      console.error('获取天气信息失败:', errorMessage)
      return
    }
    
    // 正确处理返回的数据结构
    weather.value = {
      city: data.city,
      temperature: data.temperature,
      description: data.description,
      humidity: data.humidity,
      windSpeed: data.wind_speed
    }
    
    console.log('天气数据加载成功:', weather.value) // 调试日志
  } catch (error: any) {
    const errorMessage = error.message || '获取天气信息失败'
    errors.value.weather = errorMessage
    console.error('获取天气信息失败:', error)
  } finally {
    loading.value.weather = false
  }
}

// 获取置顶文章
const fetchFeaturedArticles = async () => {
  try {
    loading.value.articles = true
    errors.value.articles = ''
    console.log('开始获取置顶文章...') // 调试日志
    
    const response = await articleApi.getTopArticles({ num: 3 })
    console.log('置顶文章API响应:', response) // 调试日志
    
    // 统一处理响应数据结构
    if (!response?.data) {
      const errorMessage = '无效的API响应'
      errors.value.articles = errorMessage
      console.error(errorMessage)
      return
    }
    
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      const errorMessage = data?.message || '获取置顶文章失败'
      errors.value.articles = errorMessage
      console.error('获取置顶文章失败:', errorMessage)
      return
    }
    
    // 设置置顶文章数据
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
    
    console.log('置顶文章加载成功:', featuredArticles.value) // 调试日志
  } catch (error: any) {
    const errorMessage = error.message || '获取置顶文章失败'
    errors.value.articles = errorMessage
    console.error('获取置顶文章失败:', error)
  } finally {
    loading.value.articles = false
  }
}

// 获取分类列表（用于标签云）
const fetchCategories = async () => {
  try {
    loading.value.categories = true
    errors.value.categories = ''
    console.log('开始获取分类列表...') // 调试日志
    
    const response = await categoryApi.getCategories({
      pagesize: -1, // 获取所有分类
      pagenum: -1
    })
    console.log('分类API响应:', response) // 调试日志
    
    // 统一处理响应数据结构
    if (!response?.data) {
      const errorMessage = '无效的API响应'
      errors.value.categories = errorMessage
      console.error(errorMessage)
      return
    }
    
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      const errorMessage = data?.message || '获取分类列表失败'
      errors.value.categories = errorMessage
      console.error('获取分类列表失败:', errorMessage)
      return
    }
    
    // 设置分类数据
    categories.value = data.map((item: any) => ({
      id: item.ID,
      name: item.name,
      articleCount: item.article_count || 0
    }))
    
    console.log('分类数据加载成功:', categories.value) // 调试日志
  } catch (error: any) {
    const errorMessage = error.message || '获取分类列表失败'
    errors.value.categories = errorMessage
    console.error('获取分类列表失败:', error)
  } finally {
    loading.value.categories = false
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

// 定时器引用
let serverStatusTimer: number | null = null

// 组件挂载时获取数据（并行执行，不阻塞）
onMounted(() => {
  console.log('Sidebar组件挂载完成，开始加载数据...') // 调试日志
  
  // 并行执行所有API调用，避免阻塞
  Promise.allSettled([
    fetchWeather(),
    fetchFeaturedArticles(),
    fetchCategories()
  ]).then(() => {
    console.log('所有Sidebar数据加载完成')
  })
  
  // 初始加载服务器状态
  fetchServerStatus()
  
  // 定期更新服务器状态
  serverStatusTimer = window.setInterval(fetchServerStatus, 30000)
})

// 组件卸载时清理定时器
onBeforeUnmount(() => {
  if (serverStatusTimer) {
    clearInterval(serverStatusTimer)
    serverStatusTimer = null
  }
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
  padding: 1rem;
  color: var(--el-text-color-secondary);
  font-size: 0.9rem;
}

/* 骨架屏样式 */
.skeleton-loader {
  padding: 1rem;
}

.skeleton-header {
  width: 40%;
  height: 1.2rem;
  margin-bottom: 1.5rem;
  background: var(--el-fill-color-light);
  border-radius: 4px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.skeleton-line {
  width: 100%;
  height: 1rem;
  background: var(--el-fill-color-light);
  border-radius: 4px;
}

.skeleton-tag {
  width: 30%;
  height: 1.2rem;
  background: var(--el-fill-color-light);
  border-radius: 20px;
  display: inline-block;
  margin: 0.2rem;
}

/* 错误信息样式 */
.error-message {
  text-align: center;
  padding: 1rem;
  color: var(--el-color-danger);
}

.error-message p {
  margin: 0 0 1rem 0;
}

.retry-button {
  background: var(--el-color-primary);
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.retry-button:hover {
  background: var(--el-color-primary-light-3);
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