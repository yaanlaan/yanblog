<template>
  <div class="sidebar-card featured-articles">
    <div class="card-header">
      <h3>置顶文章</h3>
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
            <div class="article-title">{{ article.title }}</div>
            <div class="article-date">{{ formatDate(article.createdAt) }}</div>
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

.article-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
  padding: 5px;
}

.article-item {
  border-bottom: 1px solid #eee;
  padding: 15px;
  border-radius: 8px;
  transition: all 0.3s ease;
  margin: 0 -10px; /* 补偿.card-content的padding */
  width: calc(100% + 20px);
}

.article-item:hover {
  background-color: #f8f9fa;
  border-radius: 10px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.article-item:last-child {
  border-bottom: none;
  padding-bottom: 15px;
}

.article-link {
  text-decoration: none;
  color: inherit;
  display: block;
}

.article-link:hover {
  color: #007bff; /* 绿色 */
}

.article-title {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 8px;
  line-height: 1.4;
}

.article-date {
  font-size: 13px;
  color: #888;
}

.empty-state {
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