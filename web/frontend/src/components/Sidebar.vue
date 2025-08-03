<!-- Sidebar 组件：展示侧边栏内容，包含天气信息、置顶文章、标签云和服务器状态四个模块 -->
<template>
  <div class="sidebar">
    <WeatherCard ref="weatherCard" />
    <FeaturedArticles ref="featuredArticles" />
    <TagCloud ref="tagCloud" />
    <ServerStatus ref="serverStatus" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import WeatherCard from './sidebar/WeatherCard.vue'
import FeaturedArticles from './sidebar/FeaturedArticles.vue'
import TagCloud from './sidebar/TagCloud.vue'
import ServerStatus from './sidebar/ServerStatus.vue'

// 组件引用
const weatherCard = ref<InstanceType<typeof WeatherCard> | null>(null)
const featuredArticles = ref<InstanceType<typeof FeaturedArticles> | null>(null)
const tagCloud = ref<InstanceType<typeof TagCloud> | null>(null)
const serverStatus = ref<InstanceType<typeof ServerStatus> | null>(null)

// 暴露刷新方法
defineExpose({
  refreshAll: () => {
    weatherCard.value?.fetchWeather()
    featuredArticles.value?.fetchArticles()
    tagCloud.value?.fetchCategories()
    serverStatus.value?.fetchServerStatus()
  }
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

/* ==================== 通用卡片样式 ==================== */
.sidebar-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 50% 50%;
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