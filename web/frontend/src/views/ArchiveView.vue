<template>
  <MainLayout>
    <template #main>
      <div class="archive-page">
        <div class="page-header">
          <h1><i class="iconfont icon-archive"></i> 文章归档</h1>
          <p class="subtitle">共 {{ totalArticles }} 篇文章，继续加油！</p>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载归档中...</p>
        </div>

        <div v-else class="timeline">
          <div v-for="(group, year) in groupedArticles" :key="year" class="timeline-year">
            <h2 class="year-title">{{ year }}</h2>
            
            <div v-for="(articles, month) in group" :key="month" class="timeline-month">
              <h3 class="month-title">{{ month }}月</h3>
              
              <div class="timeline-items">
                <div v-for="article in articles" :key="article.id" class="timeline-item">
                  <span class="date">{{ formatDateDay(article.createdAt || article.CreatedAt) }}</span>
                  <router-link :to="`/article/${article.id}`" class="title">
                    {{ article.title }}
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/components/layout/MainLayout.vue'
import { articleApi } from '@/services/api'

interface Article {
  id: number
  title: string
  createdAt?: string
  CreatedAt?: string
}

const loading = ref(false)
const articles = ref<Article[]>([])
const totalArticles = ref(0)

// 获取归档数据
const fetchArchive = async () => {
  loading.value = true
  try {
    // 获取足够多的文章用于归档
    const res = await articleApi.getArticles({ pagesize: 1000, pagenum: 1 })
    if (res.status === 200 && res.data.status === 200) {
      articles.value = res.data.data
      totalArticles.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to fetch archive:', error)
  } finally {
    loading.value = false
  }
}

// 按年月分组
const groupedArticles = computed(() => {
  const groups: Record<string, Record<string, Article[]>> = {}
  
  articles.value.forEach(article => {
    // 兼容后端可能返回 CreatedAt 或 createdAt
    const dateStr = article.createdAt || article.CreatedAt
    if (!dateStr) return

    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return

    const year = date.getFullYear().toString()
    const month = (date.getMonth() + 1).toString()
    
    if (!groups[year]) {
      groups[year] = {}
    }
    if (!groups[year][month]) {
      groups[year][month] = []
    }
    groups[year][month].push(article)
  })
  
  // 排序：年份倒序，月份倒序
  const sortedGroups: Record<string, Record<string, Article[]>> = {}
  Object.keys(groups).sort((a, b) => Number(b) - Number(a)).forEach(year => {
    sortedGroups[year] = {}
    Object.keys(groups[year]).sort((a, b) => Number(b) - Number(a)).forEach(month => {
      sortedGroups[year][month] = groups[year][month]
    })
  })
  
  return sortedGroups
})

const formatDateDay = (dateStr?: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  return `${date.getDate()}日`
}

onMounted(() => {
  fetchArchive()
})
</script>

<style scoped>
.archive-page {
  background: var(--color-background);
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--color-border);
}

.page-header h1 {
  font-size: 28px;
  color: var(--color-heading);
  margin-bottom: 10px;
}

.subtitle {
  color: var(--color-text);
  opacity: 0.8;
}

.timeline {
  position: relative;
  padding-left: 20px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: var(--color-border);
}

.timeline-year {
  margin-bottom: 40px;
}

.year-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--color-heading);
  margin-bottom: 20px;
  position: relative;
}

.year-title::before {
  content: '';
  position: absolute;
  left: -25px;
  top: 50%;
  transform: translateY(-50%);
  width: 12px;
  height: 12px;
  background: #06bac7ff;
  border-radius: 50%;
  border: 2px solid var(--color-background);
}

.timeline-month {
  margin-left: 20px;
  margin-bottom: 20px;
}

.month-title {
  font-size: 18px;
  color: var(--color-text);
  margin-bottom: 15px;
  opacity: 0.9;
}

.timeline-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.timeline-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px;
  border-radius: 6px;
  transition: all 0.3s;
}

.timeline-item:hover {
  background: var(--color-background-soft);
  transform: translateX(5px);
}

.date {
  color: var(--color-text);
  opacity: 0.7;
  font-size: 14px;
  min-width: 40px;
}

.title {
  color: var(--color-heading);
  text-decoration: none;
  font-size: 16px;
  transition: color 0.3s;
}

.title:hover {
  color: #21e7ddff;
}

.loading-state {
  text-align: center;
  padding: 40px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: #04aa94ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 15px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
