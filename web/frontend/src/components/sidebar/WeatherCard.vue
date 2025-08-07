<template>
  <div class="sidebar-card weather-card">
    <div class="card-header">
      <h3>天气信息</h3>
    </div>
    <div class="card-content">
      <div v-if="loading" class="skeleton-loader">
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
            <span v-for="(icon, index) in weatherIcons" :key="index" class="icon">{{ icon }}</span>
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
      <div class="error-message" v-else-if="error">
        <p>❌ {{ error }}</p>
        <button @click="onRetry" class="retry-button">重试</button>
      </div>
      <div class="weather-placeholder" v-else>
        <p>暂无天气信息</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { weatherApi } from '@/services/api'

// 定义组件属性
interface Weather {
  city: string
  temperature: number
  description: string
  humidity: number
  windSpeed: number
}

const weather = ref<Weather | null>(null)
const loading = ref(false)
const error = ref('')

// 定义事件
const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

// 计算属性：根据天气描述获取所有对应的图标
const weatherIcons = computed(() => {
  if (!weather.value) return [];
  
  const description = weather.value.description;
  const icons = [];
  
  // 如果包含多种天气，拆分处理
  if (description.includes('，')) {
    const types = description.split('，');
    types.forEach(type => {
      const icon = getSingleWeatherIcon(type.trim());
      if (icon) {
        icons.push(icon);
      }
    });
  } else {
    // 单一天气
    const icon = getSingleWeatherIcon(description);
    if (icon) {
      icons.push(icon);
    }
  }
  
  return icons;
});

// 获取单个天气图标
const getSingleWeatherIcon = (description: string) => {
  // 根据天气描述返回对应的图标
  switch (description) {
    case '晴':
      return '☀️'
    case '多云':
      return '☁️'
    case '阴':
      return '⛅'
    case '阵雨':
      return '🌦️'
    case '雷阵雨':
      return '⛈️'
    case '小雨':
      return '🌧️'
    case '中雨':
      return '🌧️'
    case '大雨':
      return '🌧️'
    case '暴雨':
      return '🌧️'
    case '小雪':
      return '🌨️'
    case '中雪':
      return '🌨️'
    case '大雪':
      return '🌨️'
    case '暴雪':
      return '🌨️'
    case '雾':
      return '🌫️'
    case '霾':
      return '🌫️'
    default:
      return '🌈' // 默认图标
  }
}

// 获取天气信息
const fetchWeather = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    const response = await weatherApi.getWeather()
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      error.value = response.data.message || '获取天气信息失败'
      console.error('获取天气信息失败:', response.data.message)
      return
    }
    
    // 设置天气数据
    weather.value = {
      city: data.city,
      temperature: data.temperature,
      description: data.description,
      humidity: data.humidity,
      windSpeed: data.wind_speed
    }
  } catch (err: any) {
    error.value = err.message || '获取天气信息失败'
    console.error('获取天气信息失败:', err)
  } finally {
    loading.value = false
    emit('loading', false)
  }
}

// 重试函数
const onRetry = () => {
  fetchWeather()
}

// 暴露方法给父组件
defineExpose({
  fetchWeather
})

// 组件挂载时获取数据
onMounted(() => {
  fetchWeather()
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

.weather-info {
  text-align: center;
  padding: 10px;
}

.weather-main {
  margin-bottom: 20px;
}

.city {
  font-size: 20px;
  font-weight: bold;
  color: #333;
  margin-bottom: 15px;
}

.weather-icon {
  font-size: 48px;
  margin: 15px 0;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.temperature {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  margin: 10px 0;
}

.weather-description {
  font-size: 16px;
  color: #666;
  margin: 10px 0;
}

.weather-details {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
  padding: 15px;
  border-radius: 10px;
  background-color: #f8f9fa;
}

.detail-item {
  text-align: center;
  padding: 5px 10px;
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
  padding: 30px 10px;
  color: #888;
}

.skeleton-loader {
  animation: skeleton-loading 1s linear infinite alternate;
  padding: 10px;
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
  margin-bottom: 15px;
  border-radius: 4px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
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