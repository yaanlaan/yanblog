<template>
  <div class="article-detail-page">
    <MainLayout>
      <template #leftSidebar>
        <ArticleToc 
          v-if="article" 
          :content="article.content" 
          ref="tocRef"
        />
      </template>
      
      <template #main>
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
            <ArticleHeader v-if="article" :article="article" />
            
            <ArticleContent v-if="article" :article="article" />
            
            <div class="empty-state" v-if="!article">
              <p>文章不存在或已被删除</p>
            </div>
          </div>
        </article>
      </template>
      
      <template #sidebar>
        <Sidebar />
      </template>
    </MainLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { articleApi } from '@/services/api'
import MainLayout from '@/components/layout/MainLayout.vue'
import ArticleHeader from '@/components/article/ArticleHeader.vue'
import ArticleContent from '@/components/article/ArticleContent.vue'
import ArticleToc from '@/components/article/ArticleToc.vue'
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
const tocRef = ref<InstanceType<typeof ArticleToc> | null>(null)

// 响应式数据
const article = ref<Article | null>(null)
const loading = ref(false)

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
    
    // 更新目录
    setTimeout(() => {
      if (tocRef.value) {
        tocRef.value.extractHeaders()
      }
    }, 0)
  } catch (error) {
    console.error('获取文章详情失败:', error)
    article.value = null
  } finally {
    loading.value = false
  }
}

// 监听路由变化
watch(
  () => route.params.id,
  (newId) => {
    const articleId = Number(newId)
    if (articleId) {
      getArticleDetail(articleId)
    }
  }
)

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
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 30px;
  min-height: 400px;
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

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #888;
}

@media (max-width: 992px) {
  .article-content {
    padding: 20px;
  }
}
</style>