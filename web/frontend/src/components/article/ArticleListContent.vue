<template>
  <div class="articles">
    <!-- 骨架屏加载状态 (仅在无数据且加载中显示) -->
    <div v-if="loading && articles.length === 0" class="articles-grid" :class="{ 'list-view': viewMode === 'list' }">
      <div v-for="i in 6" :key="i" class="skeleton-card" :class="{ 'is-list-mode': viewMode === 'list' }">
        <el-skeleton animated class="skeleton-content">
          <template #template>
            <!-- Grid 模式 -->
            <div v-if="viewMode === 'grid'" class="skeleton-grid-layout">
              <el-skeleton-item variant="image" class="skeleton-cover" />
              <div class="skeleton-body">
                <el-skeleton-item variant="h3" style="width: 60%; margin-bottom: 15px;" />
                <div style="display: flex; gap: 10px; margin-bottom: 20px;">
                  <el-skeleton-item variant="text" style="width: 30%" />
                  <el-skeleton-item variant="text" style="width: 20%" />
                </div>
                 <el-skeleton-item variant="text" style="width: 40%" />
              </div>
            </div>

            <!-- List 模式 -->
            <div v-else class="skeleton-list-layout">
              <el-skeleton-item variant="image" class="skeleton-cover-list" />
              <div class="skeleton-body-list">
                <el-skeleton-item variant="h1" style="width: 40%; margin-bottom: 15px; height: 28px;" />
                <el-skeleton-item variant="p" style="width: 90%; margin-bottom: 10px;" />
                <el-skeleton-item variant="p" style="width: 70%; margin-bottom: 15px;" />
                <div style="display: flex; gap: 15px;">
                  <el-skeleton-item variant="text" style="width: 100px;" />
                  <el-skeleton-item variant="text" style="width: 150px;" />
                </div>
              </div>
            </div>
          </template>
        </el-skeleton>
      </div>
    </div>

    <!-- 真实文章列表 -->
    <div v-else class="articles-grid" :class="{ 'list-view': viewMode === 'list' }">
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
  color: var(--color-text-secondary);
  background: var(--color-background-soft);
  border-radius: 12px;
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}

/* 骨架屏样式 */
.skeleton-card {
  background: var(--color-background-soft);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px var(--color-shadow);
  border: 1px solid var(--color-border);
}

/* Grid 模式骨架屏 */
.skeleton-grid-layout {
  display: flex;
  flex-direction: column;
}

.skeleton-cover {
  width: 100%;
  height: 180px !important; /* 强制覆盖 el-skeleton-item 默认样式 */
}

.skeleton-body {
  padding: 20px;
}

/* List 模式骨架屏 */
.skeleton-card.is-list-mode {
  height: 200px;
}

.skeleton-list-layout {
  display: flex;
  height: 100%;
}

.skeleton-cover-list {
  width: 320px !important;
  height: 100% !important;
  flex-shrink: 0;
}

.skeleton-body-list {
  padding: 25px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

@media (max-width: 768px) {
  .skeleton-card.is-list-mode {
    height: auto;
  }
  
  .skeleton-list-layout {
    flex-direction: column;
  }

  .skeleton-cover-list {
    width: 100% !important;
    height: 180px !important;
  }
}
</style>