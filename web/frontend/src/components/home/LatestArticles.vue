<template>
  <div class="latest-articles">
    <h2 class="section-title">最新文章</h2>

    <!-- 初始加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 文章列表 -->
    <div v-else>
      <div class="articles-list" v-if="articles.length > 0">
        <template v-for="(article, index) in articles" :key="article.id">
          <ArticleCard :article="article" />
          <div v-if="index < articles.length - 1" class="article-divider"></div>
        </template>

        <!-- 加载更多按钮 / 加载中提示 / 没有更多提示 -->
        <div class="see-more-container" v-if="articles.length > 0">
          <button
            v-if="hasMore && !loadingMore"
            @click="$emit('loadMore')"
            class="see-more-button"
          >
            <i class="iconfont icon-seemore"></i>
            <span>See More</span>
          </button>
          <div v-else-if="loadingMore" class="loading-hint">
            <div class="mini-spinner"></div>
            <span>加载中...</span>
          </div>
          <div v-else-if="!hasMore" class="no-more-hint">— 已经到底啦 —</div>
        </div>
      </div>

      <div class="empty-state" v-else>
        <p>暂无文章</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ArticleCard from './ArticleCard.vue'

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
  loadingMore?: boolean
  hasMore?: boolean
}

defineProps<Props>()
defineEmits<{
  (e: 'loadMore'): void
}>()
</script>

<style scoped>
.iconfont {
  font-size: 14px;
}
.section-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: var(--color-heading);
  text-align: left;
  padding-left: 10px;
  border-left: 4px solid var(--color-accent);
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
  margin-bottom: 30px;
}

.article-divider {
  height: 1px;
  background-color: #e0e0e0;
  margin: 0 20px;
  width: calc(100% - 40px);
  align-self: center;
}

.see-more-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.see-more-button {
  padding: 10px 30px;
  background-color: var(--color-accent);
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  color: white;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 5px;
  box-shadow: 0 4px 15px color-mix(in srgb, var(--color-accent) 40%, transparent);
}

.see-more-button:hover {
  background-color: color-mix(in srgb, var(--color-accent) 82%, black);
  box-shadow: 0 6px 20px color-mix(in srgb, var(--color-accent) 60%, transparent);
  transform: translateY(-2px);
}

.loading-hint {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-text-secondary, #888);
  font-size: 14px;
}

.mini-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid #f3f3f3;
  border-top: 2px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
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