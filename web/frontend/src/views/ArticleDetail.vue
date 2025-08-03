<template>
  <div class="article-detail-page">
    <div class="container">
      <div class="content-wrapper">
        <div class="main-content">
          <div class="page-header">
            <router-link to="/articles" class="back-link">
              « 返回文章列表
            </router-link>
          </div>
          
          <article class="article-content">
            <!-- 加载状态 -->
            <div v-if="loading" class="loading-state">
              <div class="spinner"></div>
              <p>加载中...</p>
            </div>
            
            <!-- 文章内容 -->
            <div v-else>
              <div class="article-header" v-if="article">
                <h1 class="article-title">{{ article.title }}</h1>
                <div class="article-meta">
                  <span class="category">
                    分类: {{ article.categoryName }}
                  </span>
                  <span class="date">
                    发布时间: {{ formatDate(article.createdAt) }}
                  </span>
                  <span class="date">
                    更新时间: {{ formatDate(article.updatedAt) }}
                  </span>
                </div>
              </div>
              
              <div class="article-body" v-if="article">
                <div class="article-description" v-if="article.desc">
                  <blockquote>{{ article.desc }}</blockquote>
                </div>
                
                <div class="article-main-content">
                  <div class="content" v-html="renderMarkdown(article.content)"></div>
                </div>
              </div>
              
              <div class="empty-state" v-if="!article">
                <p>文章不存在或已被删除</p>
              </div>
            </div>
          </article>
        </div>
        
        <!-- 右侧边栏 -->
        <Sidebar />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { articleApi } from '@/services/api'
import { marked } from 'marked'
import Sidebar from '@/components/Sidebar.vue'

// 类型定义
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

// 路由信息
const route = useRoute()

// 响应式数据
const article = ref<Article | null>(null)
const loading = ref(false)

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 渲染Markdown内容
const renderMarkdown = (content: string) => {
  if (!content) return ''
  return marked.parse(content)
}

// 获取文章详情
const getArticleDetail = async (id: number) => {
  loading.value = true
  try {
    const response = await articleApi.getArticle(id)
    const data = response.data.data
    
    article.value = {
      id: data.ID,
      title: data.title,
      categoryId: data.cid,
      categoryName: data.Category?.name || '未分类',
      desc: data.desc,
      content: data.content,
      img: data.img,
      createdAt: data.CreatedAt || data.created_at,
      updatedAt: data.UpdatedAt || data.updated_at
    }
  } catch (error) {
    console.error('获取文章详情失败:', error)
    article.value = null
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  const articleId = Number(route.params.id)
  if (articleId) {
    getArticleDetail(articleId)
  }
})
</script>

<style scoped>
.article-detail-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.content-wrapper {
  display: flex;
  gap: 30px;
  min-height: calc(100vh - 280px);
}

.main-content {
  flex: 1;
  min-width: 0;
}

.page-header {
  margin-bottom: 20px;
}

.back-link {
  color: #007bff;
  text-decoration: none;
  font-size: 14px;
}

.back-link:hover {
  text-decoration: underline;
}

.article-content {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 40px;
  min-height: 400px;
  position: relative;
  z-index: 1;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.article-header {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.article-title {
  font-size: 28px;
  color: #333;
  margin: 0 0 15px 0;
  line-height: 1.3;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  font-size: 14px;
  color: #888;
}

.article-description blockquote {
  font-size: 16px;
  color: #666;
  border-left: 4px solid #007bff;
  padding: 10px 20px;
  margin: 0 0 30px 0;
  background: #f8f9fa;
}

.article-main-content {
  font-size: 16px;
  line-height: 1.8;
}

.content :deep(h1) {
  font-size: 24px;
  margin: 24px 0 16px;
  color: #333;
}

.content :deep(h2) {
  font-size: 22px;
  margin: 22px 0 14px;
  color: #333;
}

.content :deep(h3) {
  font-size: 20px;
  margin: 20px 0 12px;
  color: #333;
}

.content :deep(p) {
  margin: 16px 0;
  color: #444;
}

.content :deep(ul),
.content :deep(ol) {
  padding-left: 30px;
  margin: 16px 0;
}

.content :deep(li) {
  margin-bottom: 8px;
}

.content :deep(code) {
  background: #f1f1f1;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
  font-size: 14px;
}

.content :deep(pre) {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
  margin: 20px 0;
}

.content :deep(pre code) {
  background: none;
  padding: 0;
}

.content :deep(blockquote) {
  border-left: 4px solid #ddd;
  padding: 10px 20px;
  margin: 20px 0;
  background: #f8f9fa;
  color: #666;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 20px 0;
}

.content :deep(a) {
  color: #007bff;
  text-decoration: none;
}

.content :deep(a:hover) {
  text-decoration: underline;
}

.content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: left;
}

.content :deep(th) {
  background: #f8f9fa;
  font-weight: bold;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #888;
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .article-content {
    padding: 25px;
    min-height: auto;
  }
  
  .article-title {
    font-size: 24px;
  }
  
  .article-meta {
    flex-direction: column;
    gap: 10px;
  }
}
</style>