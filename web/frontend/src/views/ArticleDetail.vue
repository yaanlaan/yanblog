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
            
            <!-- 上一篇/下一篇导航 -->
            <div class="article-navigation" v-if="article">
              <div class="nav-item previous" v-if="previousArticle">
                <router-link :to="`/article/${previousArticle.id}`" class="nav-link">
                  <div class="nav-cover">
                    <img :src="previousArticle.img || defaultImage" :alt="previousArticle.title" />
                  </div>
                  <span class="nav-label">上一篇</span>
                  <span class="nav-title">{{ previousArticle.title }}</span>
                </router-link>
              </div>
              <div class="nav-item next" v-if="nextArticle">
                <router-link :to="`/article/${nextArticle.id}`" class="nav-link">
                  <div class="nav-cover">
                    <img :src="nextArticle.img || defaultImage" :alt="nextArticle.title" />
                  </div>
                  <span class="nav-label">下一篇</span>
                  <span class="nav-title">{{ nextArticle.title }}</span>
                </router-link>
              </div>
            </div>
            
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

// 默认图片
const defaultImage = new URL('@/assets/img/无封面.jpg', import.meta.url).href

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
const previousArticle = ref<Article | null>(null)
const nextArticle = ref<Article | null>(null)

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
    
    // 获取相邻文章
    await getAdjacentArticles(id)
    
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

// 获取相邻文章
const getAdjacentArticles = async (currentId: number) => {
  try {
    // 由于后端没有提供专门的API，我们通过获取文章列表来模拟相邻文章功能
    // 这里只是简单实现，实际项目中应该由后端提供专门的API
    
    // 获取所有文章列表
    const response = await articleApi.getArticles({
      pagesize: -1,
      pagenum: -1
    })
    
    const articles = response.data.data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
    
    // 按ID排序（实际项目中应该按创建时间排序）
    articles.sort((a: Article, b: Article) => a.id - b.id)
    
    // 查找当前文章的索引
    const currentIndex = articles.findIndex((article: Article) => article.id === currentId)
    
    // 设置上一篇和下一篇文章
    if (currentIndex > 0) {
      previousArticle.value = articles[currentIndex - 1]
    } else {
      previousArticle.value = null
    }
    
    if (currentIndex < articles.length - 1) {
      nextArticle.value = articles[currentIndex + 1]
    } else {
      nextArticle.value = null
    }
  } catch (error) {
    console.error('获取相邻文章失败:', error)
    previousArticle.value = null
    nextArticle.value = null
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

/* 上一篇/下一篇导航样式 */
.article-navigation {
  display: flex;
  justify-content: space-between;
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.nav-item {
  flex: 1;
  min-width: 0;
}

.nav-item.previous {
  padding-right: 10px;
}

.nav-item.next {
  padding-left: 10px;
  text-align: right;
}

.nav-link {
  display: block;
  text-decoration: none;
  color: #333;
  transition: all 0.3s ease;
}

.nav-link:hover {
  color: #007bff;
}

.nav-label {
  display: block;
  font-size: 12px;
  color: #888;
  margin-bottom: 5px;
}

.nav-title {
  display: block;
  font-size: 16px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 封面预览样式 */
.nav-cover {
  width: 100%;
  height: 120px;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 10px;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

@media (max-width: 992px) {
  .article-content {
    padding: 20px;
  }
  
  .article-navigation {
    flex-direction: column;
    gap: 20px;
  }
  
  .nav-item.previous,
  .nav-item.next {
    padding: 0;
    text-align: left;
  }
  
  .nav-item.next {
    text-align: left;
  }
  
  .nav-cover {
    height: 80px;
  }
}
</style>