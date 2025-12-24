<template>
  <div class="sidebar-card featured-articles">
    <div class="card-header">
      <h3><i class="iconfont icon-pushpin" style="color: #333; margin-right: 5px;"></i> 置顶博客</h3>
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
            <i class="iconfont icon-rss article-icon"></i>
            <span class="article-title">{{ article.title }}</span>
          </router-link>
        </div>
      </div>
      <div class="error-message" v-else-if="error">
        <p>❌ {{ error }}</p>
        <button @click="onRetry" class="retry-button">重试</button>
      </div>
      <div class="empty-state" v-else>
        <p>暂无置顶文章</p>
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

// 获取置顶文章
const fetchArticles = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    const response = await articleApi.getTopArticles({ num: 5 })
    const { data, status } = response.data
    
    // 检查API返回状态
    if (status !== 200) {
      error.value = response.data.message || '获取置顶文章失败'
      console.error('获取置顶文章失败:', response.data.message)
      return
    }
    
    // 设置文章数据
    articles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '',
      desc: item.desc,
      content: item.content,
      img: item.img,
      createdAt: item.CreatedAt || item.created_at || '',
      updatedAt: item.UpdatedAt || item.updated_at || ''
    }))
  } catch (err: any) {
    error.value = err.message || '获取置顶文章失败'
    console.error('获取置顶文章失败:', err)
  } finally {
    loading.value = false
    emit('loading', false)
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 重试函数
const onRetry = () => {
  fetchArticles()
}

// 暴露方法给父组件
defineExpose({
  fetchArticles
})

// 组件挂载时获取数据
onMounted(() => {
  fetchArticles()
})
</script>

<style scoped>
.sidebar-card.featured-articles {
  background: transparent !important;
  box-shadow: none !important;
  border: none !important;
  padding: 0;
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
}

.article-item:hover {
  background-color: transparent;
  box-shadow: none;
}

.article-link {
  text-decoration: none;
  color: #333;
  display: flex;
  align-items: flex-start;
  gap: 8px;
  transition: color 0.3s ease;
}

.article-link:hover {
  color: #42b883;
}

.article-icon {
  font-size: 18px;
  color: #333;
  margin-top: 2px;
}

.article-title {
  font-size: 15px;
  font-weight: 400;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Skeleton & Error styles remain similar but simplified */
.empty-state, .error-message {
  text-align: center;
  padding: 20px 0;
  color: #888;
}
</style>