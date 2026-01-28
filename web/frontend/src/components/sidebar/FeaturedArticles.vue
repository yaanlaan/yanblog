<template>
  <div class="sidebar-card featured-articles">
    <div class="card-header">
      <h3>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 5px;">
          <path d="M4 11a9 9 0 0 1 9 9"></path>
          <path d="M4 4a16 16 0 0 1 16 16"></path>
          <circle cx="5" cy="19" r="1"></circle>
        </svg>
        热门博客
      </h3>
    </div>
    <div class="card-content">
      <div v-if="loading" class="skeleton-loader">
        <div class="skeleton-header"></div>
        <div class="skeleton-body">
          <div class="skeleton-line"></div>
          <div class="skeleton-line"></div>
        </div>
      </div>
      <div class="article-list" v-else-if="articles.length > 0">
        <div 
          v-for="article in articles" 
          :key="article.id" 
          class="article-item"
        >
          <router-link :to="`/article/${article.id}`" class="article-link">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="title-icon">
              <path d="M4 11a9 9 0 0 1 9 9"></path>
              <path d="M4 4a16 16 0 0 1 16 16"></path>
              <circle cx="5" cy="19" r="1"></circle>
            </svg>
            <span class="article-title">{{ article.title }}</span>
            <span class="hover-line"></span>
          </router-link>
        </div>
      </div>
      <div class="error-message" v-else-if="error">
        <p>❌ {{ error }}</p>
        <button @click="onRetry" class="retry-button">重试</button>
      </div>
      <div class="empty-state" v-else>
        <p>暂无热门文章</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { articleApi } from '@/services/api'

// 定义文章接口
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  views: number
  createdAt: string
  updatedAt: string
}

const articles = ref<Article[]>([])
const loading = ref(false)
const error = ref('')

// 定义事件
const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

// 获取热门文章
const fetchArticles = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    const response = await articleApi.getHotArticles({ num: 5 })
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      error.value = response.data.message || '获取热门文章失败'
      return
    }
    
    articles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      views: item.views || 0,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
  } catch (err: any) {
    console.error('获取热门文章错误:', err)
    error.value = '网络请求失败'
  } finally {
    loading.value = false
    emit('loading', false)
  }
}

onMounted(() => {
  fetchArticles()
})

defineExpose({
  fetchArticles
})
</script>

<style scoped>
.sidebar-card.featured-articles {
  background: transparent !important;
  box-shadow: none !important;
  border: none !important;
  padding: 0;
  transform: none !important;
  transition: none !important;
}

.sidebar-card.featured-articles:hover {
  transform: none !important;
  box-shadow: none !important;
}

.card-header {
  padding: 15px 0;
  border-bottom: none;
  background: transparent;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
  font-weight: 600;
  display: flex;
  align-items: center;
}

.card-content {
  padding: 0;
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.article-item {
  border-bottom: none;
  padding: 0;
  margin: 0;
  width: 100%;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.article-item:hover {
  transform: translateY(-2px);
  z-index: 1;
  position: relative;
}

.article-link {
  text-decoration: none;
  color: var(--color-text);
  display: flex;
  align-items: flex-start;
  gap: 8px;
  transition: color 0.3s ease;
  position: relative;
  padding: 5px 0;
}

.article-link:hover {
  color: var(--color-accent);
}

.title-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.article-title {
  font-size: 15px;
  font-weight: 500;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.hover-line {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 0;
  height: 2px;
  background-color: var(--color-accent);
  transition: width 0.3s ease;
  content: "";
}

.article-link:hover .hover-line {
  width: 100%;
}

.retry-button {
  margin-top: 10px;
  padding: 5px 15px;
  background-color: var(--color-accent);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.empty-state, .error-message {
  text-align: center;
  padding: 20px 0;
  color: var(--color-text-secondary);
}

/* Skeleton Loading */
.skeleton-loader {
  padding: 10px;
}

.skeleton-header {
  height: 20px;
  background-color: var(--color-background-mute);
  margin-bottom: 15px;
  border-radius: 4px;
  width: 60%;
}

.skeleton-line {
  height: 16px;
  background-color: var(--color-background-mute);
  margin-bottom: 10px;
  border-radius: 4px;
}

.skeleton-line:last-child {
  width: 80%;
}
</style>