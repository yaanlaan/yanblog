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
      <div class="articles-grid" v-if="displayArticles.length > 0">
        <ArticleCard 
          v-for="article in displayArticles" 
          :key="article.id" 
          :article="article"
        />
      </div>
      
      <!-- seemore按钮 -->
      <div class="see-more-container" v-if="articles.length > displayCount">
        <button 
          @click="loadMore"
          class="see-more-button"
        >
          <i class="iconfont icon-seemore"></i>
          <span> seemore</span>
        </button>
      </div>
      
      <div class="empty-state" v-else>
        <p>暂无文章</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
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

const props = defineProps<Props>()

// 显示数量控制
const displayCount = ref(5)

// 计算属性，根据displayCount显示文章
const displayArticles = computed(() => {
  return props.articles.slice(0, displayCount.value)
})

// 加载更多文章
const loadMore = () => {
  displayCount.value += 5
}
</script>

<style scoped>
.iconfont {
  font-size: 10px;
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

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #888;
}

.see-more-container {
  display: flex;
  justify-content: center;
  width: 100%;
  margin-top: 30px;
}

.see-more-button {
  width: 40%;
  padding: 10px;
  background-color: #f8f9fa;
  color: #007bff;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.see-more-button:hover {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
  border-radius: 20px;
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