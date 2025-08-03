<template>
  <div class="home-page">
    <div class="container">
      <div class="content-wrapper">
        <div class="main-content">
          <div class="hero-section">
            <div class="hero-content">
              <h1 class="hero-title">欢迎来到我的博客</h1>
              <p class="hero-description">
                这里记录了我的技术学习、生活感悟和各种经验分享。
                希望我的文章能对你有所帮助。
              </p>
              <router-link to="/articles" class="hero-button">
                浏览文章
              </router-link>
            </div>
          </div>

          <div class="latest-articles">
            <h2 class="section-title">最新文章</h2>
            
            <!-- 加载状态 -->
            <div v-if="loading" class="loading-state">
              <div class="spinner"></div>
              <p>加载中...</p>
            </div>
            
            <!-- 文章列表 -->
            <div v-else>
              <div class="articles-grid" v-if="articles.length > 0">
                <div 
                  v-for="article in articles" 
                  :key="article.id" 
                  class="article-card"
                >
                  <div class="article-content">
                    <h3 class="article-title">
                      <router-link :to="`/article/${article.id}`">
                        {{ article.title }}
                      </router-link>
                    </h3>
                    <div class="article-meta">
                      <span class="category">{{ article.categoryName }}</span>
                      <span class="date">{{ formatDate(article.createdAt) }}</span>
                    </div>
                    <p class="article-excerpt">{{ article.desc || '暂无简介' }}</p>
                    <router-link :to="`/article/${article.id}`" class="read-more">
                      阅读全文 »
                    </router-link>
                  </div>
                </div>
              </div>
              <div class="empty-state" v-else>
                <p>暂无文章</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 右侧边栏 -->
        <Sidebar />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { articleApi } from '@/services/api'
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

// 响应式数据
const articles = ref<Article[]>([])
const loading = ref(false)

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 获取最新文章
const getLatestArticles = async () => {
  loading.value = true
  try {
    const response = await articleApi.getArticles({
      pagesize: 6,
      pagenum: 1
    })
    
    const { data } = response.data
    articles.value = data.map((item: any) => ({
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
  } catch (error) {
    console.error('获取最新文章失败:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getLatestArticles()
})
</script>

<style scoped>
.home-page {
  width: 100%;
  min-height: calc(100vh - 200px);
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

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 80px 0;
  margin-bottom: 40px;
  border-radius: 8px;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
  padding: 0 20px;
}

.hero-title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 20px;
}

.hero-description {
  font-size: 20px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.hero-button {
  display: inline-block;
  background: white;
  color: #667eea;
  padding: 15px 30px;
  border-radius: 30px;
  text-decoration: none;
  font-weight: 600;
  font-size: 18px;
  transition: all 0.3s;
  border: 2px solid white;
}

.hero-button:hover {
  background: transparent;
  color: white;
}

.section-title {
  font-size: 32px;
  margin-bottom: 40px;
  color: #333;
  text-align: center;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
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

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 30px;
}

.article-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  height: 100%;
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.article-content {
  padding: 25px;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.article-title {
  font-size: 20px;
  margin: 0 0 15px 0;
}

.article-title a {
  color: #333;
  text-decoration: none;
  transition: color 0.3s;
}

.article-title a:hover {
  color: #007bff;
}

.article-meta {
  display: flex;
  gap: 15px;
  font-size: 14px;
  color: #888;
  margin-bottom: 15px;
}

.article-excerpt {
  color: #666;
  line-height: 1.6;
  margin-bottom: 20px;
  flex-grow: 1;
}

.read-more {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
  margin-top: auto;
}

.read-more:hover {
  text-decoration: underline;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #888;
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .hero-title {
    font-size: 36px;
  }
  
  .hero-description {
    font-size: 18px;
  }
  
  .articles-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}

@media (max-width: 768px) {
  .container {
    padding: 0 15px;
  }
  
  .articles-grid {
    grid-template-columns: 1fr;
  }
}
</style>