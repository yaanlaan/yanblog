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
            <div class="city">{{ weather.city }}</div>
            <div class="weather-icon">
              <span class="icon">{{ getWeatherIcon(weather.description) }}</span>
            </div>
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
        <div v-if="loading.serverStatus" class="loading-placeholder">
          <p>状态加载中...</p>
        </div>
        <div v-else-if="!errors.serverStatus">
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
          <p>❌ {{ errors.serverStatus }}</p>
          <button @click="fetchServerStatus" class="retry-button">重试</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { articleApi, categoryApi, weatherApi, systemApi } from '@/services/api'

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
  startTime: number // 添加服务器启动时间戳
}

// 响应式数据
const weather = ref<Weather | null>(null)
const featuredArticles = ref<Article[]>([])
const categories = ref<Category[]>([])
const serverStatus = ref<ServerStatus>({
  status: 'offline',
  uptime: '未知',
  memoryUsage: 0,
  cpuUsage: 0,
  startTime: 0 // 初始化启动时间戳
})

// 错误状态
const errors = ref({
  weather: '',
  articles: '',
  categories: '',
  serverStatus: ''
})

// 加载状态
const loading = ref({
  weather: false,
  articles: false,
  categories: false,
  serverStatus: false
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

// 获取天气图标
const getWeatherIcon = (description: string) => {
  // 根据天气描述返回对应的图标
  switch (description) {
    case '晴':
      return '☀️';
    case '多云':
      return '☁️';
    case '阴':
      return '⛅';
    case '阵雨':
      return '🌦️';
    case '雷阵雨':
      return '⛈️';
    case '小雨':
      return '🌧️';
    case '中雨':
      return '🌧️';
    case '大雨':
      return '🌧️';
    case '暴雨':
      return '🌧️';
    case '小雪':
      return '🌨️';
    case '中雪':
      return '🌨️';
    case '大雪':
      return '🌨️';
    case '暴雪':
      return '🌨️';
    case '雾':
      return '🌫️';
    case '霾':
      return '🌫️';
    default:
      return '🌈'; // 默认图标
  }
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

// 获取服务器状态
const fetchServerStatus = async () => {
  try {
    loading.value.serverStatus = true
    errors.value.serverStatus = ''
    
    const response = await systemApi.getSystemStatus()
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      errors.value.serverStatus = response.data.message || '获取服务器状态失败'
      console.error('获取服务器状态失败:', response.data.message)
      return
    }
    
    // 设置服务器状态数据
    serverStatus.value = {
      status: data.status,
      uptime: data.uptime,
      memoryUsage: Math.round(data.memory_usage * 100) / 100, // 保留两位小数
      cpuUsage: Math.round(data.cpu_usage * 100) / 100, // 保留两位小数
      startTime: Date.now() - parseUptimeToMilliseconds(data.uptime) // 计算启动时间戳
    }
  } catch (error: any) {
    errors.value.serverStatus = error.message || '获取服务器状态失败'
    console.error('获取服务器状态失败:', error)
    // 即使获取失败，也保持在线状态
    serverStatus.value.status = 'online'
  } finally {
    loading.value.serverStatus = false
  }
}

// 将运行时间字符串解析为毫秒数
const parseUptimeToMilliseconds = (uptime: string): number => {
  // 解析格式如"1天2小时3分钟4秒"或"2小时3分钟4秒"等
  let totalMilliseconds = 0;
  
  // 匹配天数
  const daysMatch = uptime.match(/(\d+)天/);
  if (daysMatch) {
    totalMilliseconds += parseInt(daysMatch[1]) * 24 * 60 * 60 * 1000;
  }
  
  // 匹配小时
  const hoursMatch = uptime.match(/(\d+)小时/);
  if (hoursMatch) {
    totalMilliseconds += parseInt(hoursMatch[1]) * 60 * 60 * 1000;
  }
  
  // 匹配分钟
  const minutesMatch = uptime.match(/(\d+)分钟/);
  if (minutesMatch) {
    totalMilliseconds += parseInt(minutesMatch[1]) * 60 * 1000;
  }
  
  // 匹配秒数
  const secondsMatch = uptime.match(/(\d+)秒/);
  if (secondsMatch) {
    totalMilliseconds += parseInt(secondsMatch[1]) * 1000;
  }
  
  return totalMilliseconds;
}

// 计算实时运行时间
const calculateRealTimeUptime = () => {
  if (serverStatus.value.startTime <= 0) return '未知';
  
  const elapsed = Date.now() - serverStatus.value.startTime;
  return formatUptime(elapsed);
}

// 格式化运行时间
const formatUptime = (milliseconds: number): string => {
  const seconds = Math.floor(milliseconds / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  
  const remainingSeconds = seconds % 60;
  const remainingMinutes = minutes % 60;
  const remainingHours = hours % 24;
  
  if (days > 0) {
    return `${days}天${remainingHours}小时${remainingMinutes}分钟${remainingSeconds}秒`;
  } else if (hours > 0) {
    return `${remainingHours}小时${remainingMinutes}分钟${remainingSeconds}秒`;
  } else if (minutes > 0) {
    return `${remainingMinutes}分钟${remainingSeconds}秒`;
  } else {
    return `${remainingSeconds}秒`;
  }
}

// 定时器引用
let serverStatusTimer: number | null = null
let uptimeTimer: number | null = null

// 组件挂载时获取数据（并行执行，不阻塞）
onMounted(() => {
  console.log('Sidebar组件挂载完成，开始加载数据...') // 调试日志
  
  // 并行执行所有API调用，避免阻塞
  Promise.allSettled([
    fetchWeather(),
    fetchFeaturedArticles(),
    fetchCategories(),
    fetchServerStatus()
  ]).then(() => {
    console.log('所有Sidebar数据加载完成')
  })
  
  // 定期更新服务器状态（CPU、内存等）
  serverStatusTimer = window.setInterval(fetchServerStatus, 30000)
  
  // 每秒更新运行时间显示
  uptimeTimer = window.setInterval(() => {
    if (serverStatus.value.startTime > 0) {
      serverStatus.value.uptime = calculateRealTimeUptime();
    }
  }, 1000);
})

// 组件卸载时清理定时器
onBeforeUnmount(() => {
  if (serverStatusTimer) {
    clearInterval(serverStatusTimer);
    serverStatusTimer = null;
  }
  if (uptimeTimer) {
    clearInterval(uptimeTimer);
    uptimeTimer = null;
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

.city {
  font-size: 20px;
  font-weight: bold;
  color: #333;
  margin-bottom: 10px;
}

.weather-icon {
  font-size: 48px;
  margin: 10px 0;
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

/* ==================== 加载状态样式 ==================== */
.skeleton-loader {
  animation: skeleton-loading 1s linear infinite alternate;
}

@keyframes skeleton-loading {
  0% {
    background-color: hsl(200, 20%, 80%);
  }
  100% {
    background-color: hsl(200, 20%, 95%);
  }
}

.skeleton-header {
  height: 20px;
  width: 60%;
  margin-bottom: 10px;
  border-radius: 4px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.skeleton-line {
  height: 16px;
  border-radius: 4px;
}

.skeleton-line:first-child {
  width: 100%;
}

.skeleton-line:nth-child(2) {
  width: 80%;
}

.skeleton-tag {
  height: 24px;
  width: 60px;
  border-radius: 12px;
}

.error-message {
  text-align: center;
  padding: 20px 0;
  color: #dc3545;
}

.retry-button {
  margin-top: 10px;
  padding: 6px 12px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.retry-button:hover {
  background-color: #0056b3;
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