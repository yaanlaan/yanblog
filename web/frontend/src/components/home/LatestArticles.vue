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
      <div class="articles-list" v-if="articles.length > 0">
        <ArticleCard 
          v-for="article in displayArticles" 
          :key="article.id" 
          :article="article"
        />
        
        <!-- seemore按钮 -->
        <div class="see-more-container">
          <button 
            v-if="displayCount < articles.length"
            @click="loadMore"
            class="see-more-button"
          >
            <i class="iconfont icon-seemore"></i>
            <span> 查看更多</span>
          </button>
          <p v-else class="no-more-hint">没有更多文章了</p>
        </div>
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
  font-size: 14px;
}
.section-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #333;
  text-align: left;
  padding-left: 10px;
  border-left: 4px solid #42b883;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 150px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 30px;
}

.see-more-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.see-more-button {
  padding: 10px 30px;
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 5px;
}

.see-more-button:hover {
  background-color: #f8f9fa;
  color: #333;
  border-color: #ccc;
}

.no-more-hint {
  color: #999;
  font-size: 14px;
  padding: 10px 0;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #888;
}
</style>