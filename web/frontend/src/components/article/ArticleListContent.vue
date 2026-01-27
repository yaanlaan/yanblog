<template>
  <div class="articles" v-loading="loading">
    <div class="articles-grid" :class="{ 'list-view': viewMode === 'list' }">
      <ArticleItem
        v-for="article in articles" 
        :key="article.id" 
        :article="article"
        :view-mode="viewMode"
      />
    </div>
    
    <!-- 滚动加载触发器 -->
    <div ref="loadingTrigger" class="loading-trigger">
      <div v-if="loading" class="loading-more">
        <div class="spinner"></div>
        <span>加载更多...</span>
      </div>
      <div v-else-if="articles.length >= total && total > 0" class="no-more">
        <span>- 也就是这些了 -</span>
      </div>
    </div>
    
    <!-- 空状态 -->
    <div class="empty-state" v-if="!loading && articles.length === 0">
      <p>暂无文章</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import ArticleItem from './ArticleItem.vue'

// 定义Props
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  top: number
  tags: string
  createdAt: string
  updatedAt: string
}

interface Props {
  articles: Article[]
  loading: boolean
  total: number
  viewMode?: 'grid' | 'list'
}

// 定义Emits
const emit = defineEmits<{
  (e: 'load-more'): void
}>()

const props = withDefaults(defineProps<Props>(), {
  viewMode: 'grid'
})

// 滚动监听
const loadingTrigger = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

const setupObserver = () => {
    if (observer) observer.disconnect()
    
    observer = new IntersectionObserver((entries) => {
        const entry = entries[0]
        if (entry.isIntersecting && !props.loading && props.articles.length < props.total) {
            emit('load-more')
        }
    }, {
        rootMargin: '100px'
    })
    
    if (loadingTrigger.value) {
        observer.observe(loadingTrigger.value)
    }
}

onMounted(() => {
    setupObserver()
})

onUnmounted(() => {
    if (observer) observer.disconnect()
})

// 监听 articles 变化，当数据更新后重新 hook observer 或者确保 observer 还在
watch(() => props.articles, () => {
    // DOM 更新后可能需要 check 一下? IntersectionObserver 是持续的，所以不需要每次 destroy
    // 但如果在 empty 状态或者 trigger 被 v-if 移除，则需要小心。
    // 这里 trigger 始终在 DOM 中 (除了内容变化)，所以应该没事。
}, { flush: 'post' })

</script>

<style scoped>
.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 25px;
  margin-bottom: 30px;
}

.articles-grid.list-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.loading-trigger {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;
}

.loading-more {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--color-text-light);
  font-size: 14px;
}

.no-more {
  color: var(--color-text-light);
  font-size: 13px;
  opacity: 0.7;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-primary, #42b883);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #888;
  background: #f9f9f9;
  border-radius: 12px;
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>