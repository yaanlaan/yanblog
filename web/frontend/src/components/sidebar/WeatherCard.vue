<template>
  <div class="sidebar-card weather-card">
    <div class="weather-content" v-if="weather">
      <!-- Location Header -->
      <div class="location-row">
        <i class="iconfont icon-location"></i>
        <span class="city">{{ weather.city }} hefei</span>
      </div>
      
      <!-- Main Weather Info -->
      <div class="main-weather">
        <span class="temperature">{{ weather.temperature.toFixed(1) }}°</span>
        <span class="weather-text">{{ weather.description }}</span>
      </div>

      <!-- Weather Details List -->
      <div class="weather-details-list">
        <div class="detail-row">
          <i class="iconfont icon-wind"></i>
          <span>西北风 {{ weather.windSpeed }}级</span>
        </div>
        <div class="detail-row">
          <i class="iconfont icon-humidity"></i>
          <span>{{ weather.humidity }}%</span>
        </div>
        <div class="detail-row">
          <i class="iconfont icon-leaf"></i>
          <span>97 良</span>
        </div>
      </div>
    </div>
    
    <!-- Loading State -->
    <div v-else-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载天气中...</p>
    </div>
    
    <!-- Error State -->
    <div v-else class="error-state">
      <p>{{ error || '无法获取天气' }}</p>
      <button @click="fetchWeather" class="retry-btn">重试</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { weatherApi } from '@/services/api'

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

const fetchWeather = async () => {
  try {
    loading.value = true
    error.value = ''
    // Try to get real weather, but fallback to mock if it fails or returns error
    try {
        const response = await weatherApi.getWeather('Shanghai') 
        if (response.data.status === 200) {
            weather.value = response.data.data
            return // Success
        }
    } catch (e) {
        console.log("Weather API failed, using mock data")
    }
    
    // Fallback/Mock data
    weather.value = {
      city: '合肥',
      temperature: 5.3,
      description: '霾',
      humidity: 81,
      windSpeed: 2
    }
  } catch (err: any) {
    error.value = '网络错误'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchWeather()
})
</script>

<style>
/* 
  Global style block to handle dark mode overrides reliably.
*/
html[data-theme="dark"] .sidebar-card.weather-card {
  background: linear-gradient(135deg, #141e30 0%, #243b55 100%) !important;
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3) !important;
}

html[data-theme="dark"] .sidebar-card.weather-card .location-row,
html[data-theme="dark"] .sidebar-card.weather-card .temperature,
html[data-theme="dark"] .sidebar-card.weather-card .weather-text {
  color: #ffffff !important;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
}

html[data-theme="dark"] .sidebar-card.weather-card .detail-row {
  color: rgba(255, 255, 255, 0.95) !important;
}

html[data-theme="dark"] .sidebar-card.weather-card .detail-row .iconfont {
  color: rgba(255, 255, 255, 0.8) !important;
}
</style>

<style scoped>
.sidebar-card.weather-card {
  /* Standard card background (Light Mode) - Soft Blue Gradient */
  background: linear-gradient(135deg, #e0f7fa 0%, #ffffff 100%);
  border-radius: 8px;
  color: var(--color-text); /* Ensure text is dark in light mode */
  padding: 25px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  border: 1px solid var(--color-border);
  overflow: hidden;
  position: relative;
  min-height: 240px;
  display: flex;
  flex-direction: column;
  margin-bottom: 20px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  transition: all 0.3s ease;
}

.location-row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  margin-bottom: 30px;
  opacity: 0.95;
  font-weight: 400;
  color: var(--color-heading);
}

.location-row .iconfont {
  font-size: 20px;
  color: var(--color-accent);
}

.main-weather {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 35px;
}

.temperature {
  font-size: 48px;
  font-weight: 400;
  line-height: 1;
  font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
  color: var(--color-heading);
}

.weather-text {
  font-size: 36px;
  /* Serif font for the Chinese character "霾" */
  font-family: "Songti SC", "SimSun", "STSong", "Times New Roman", serif;
  font-weight: 400;
  color: var(--color-text);
}

.weather-details-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
  opacity: 0.9;
  font-weight: 300;
  color: var(--color-text);
}

.detail-row .iconfont {
  font-size: 18px;
  width: 20px;
  text-align: center;
  color: var(--color-text-light);
}

/* Loading & Error States */
.loading-state, .error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
  color: var(--color-text-light);
}

.retry-btn {
  background: var(--color-accent);
  border: 1px solid var(--color-accent);
  color: white;
  padding: 6px 18px;
  border-radius: 20px;
  cursor: pointer;
  margin-top: 15px;
  transition: background 0.3s;
}

.retry-btn:hover {
  background: var(--color-accent-hover);
}

.spinner {
  width: 30px;
  height: 30px;
  border: 3px solid var(--color-border);
  border-radius: 50%;
  border-top-color: var(--color-accent);
  animation: spin 1s ease-in-out infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
