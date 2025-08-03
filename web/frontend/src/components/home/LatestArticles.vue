<template>
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
        <ArticleCard 
          v-for="article in articles" 
          :key="article.id" 
          :article="article"
        />
      </div>
      <div class="empty-state" v-else>
        <p>暂无文章</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ArticleCard from './ArticleCard.vue'

// 定义Props
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

interface Props {
  articles: Article[]
  loading: boolean
}

defineProps<Props>()
</script>

<style scoped>
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

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #888;
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
  }
  
  .section-title {
    font-size: 28px;
  }
}

@media (max-width: 480px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>