<!-- Sidebar 组件：展示侧边栏内容 -->
<template>
  <WeatherCard ref="weatherCard" />
  <ShortcutsCard />
  <FeaturedArticles ref="featuredArticles" />
  <ServerStatus ref="serverStatus" />
  <TagCloud ref="tagCloud" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import WeatherCard from './sidebar/WeatherCard.vue'
import ShortcutsCard from './sidebar/ShortcutsCard.vue'
import FeaturedArticles from './sidebar/FeaturedArticles.vue'
import ServerStatus from './sidebar/ServerStatus.vue'
import TagCloud from './sidebar/TagCloud.vue'

// 组件引用
const featuredArticles = ref<InstanceType<typeof FeaturedArticles> | null>(null)
const tagCloud = ref<InstanceType<typeof TagCloud> | null>(null)
const weatherCard = ref<InstanceType<typeof WeatherCard> | null>(null)
const serverStatus = ref<InstanceType<typeof ServerStatus> | null>(null)

// 暴露刷新方法
defineExpose({
  refreshAll: () => {
    featuredArticles.value?.fetchArticles()
    tagCloud.value?.fetchTags()
    weatherCard.value?.fetchWeather()
    serverStatus.value?.fetchServerStatus()
  }
})
</script>

<style scoped>
/* ==================== 通用卡片样式 ==================== */
.sidebar-card {
  background: var(--color-background-soft);
  border-radius: 12px;
  box-shadow: 0 4px 20px var(--color-shadow);
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
  border: 1px solid var(--color-border);
}

.sidebar-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 30px var(--color-shadow);
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid var(--color-border);
  background: transparent;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--color-heading);
}

.card-content {
  padding: 20px;
}

/* ==================== 响应式样式 ==================== */
@media (max-width: 992px) {
  .sidebar-card {
    flex: 1 1 calc(50% - 10px);
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .sidebar-card {
    flex: 1 1 100%;
  }
}
</style>