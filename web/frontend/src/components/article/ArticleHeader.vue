<template>
  <div class="article-header">
    <h1 class="article-title">{{ article.title }}</h1>
    
    <div class="article-meta-row">
      <div class="author-info">
        <img :src="siteInfo.author_avatar || '/assets/avatar.jpg'" class="author-avatar" alt="author">
        <div class="author-details">
          <span class="author-name">{{ siteInfo.author_name || 'yaan' }}</span>
          <span class="publish-date">{{ formatDate(article.createdAt) }}</span>
        </div>
      </div>
      
      <div class="meta-stats">
         <span class="meta-item" title="字数">
           <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
           {{ formatNumber(wordCount) }} 字
         </span>
         <span class="meta-item" title="阅读时间">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            约 {{ readingTime }} 分钟
         </span>
         <span class="meta-item tag-item" title="分类" v-if="article.categoryName">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
            {{ article.categoryName }}
         </span>
      </div>
    </div>

    <!-- 操作按钮行 -->
    <div class="article-actions-bar">
      <!-- 占位，保持右对齐 -->
      <div></div>
      
      <div class="right-actions">
         <button class="action-btn-pill comment-btn" title="评论" @click="$emit('comment')">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"></path></svg>
            评论
         </button>
         <button class="action-btn-pill share-btn" title="分享" @click="$emit('share')">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 12v8a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-8"></path><polyline points="16 6 12 2 8 6"></polyline><line x1="12" y1="2" x2="12" y2="15"></line></svg>
         </button>
      </div>
    </div>

    <div class="article-image" v-if="article.img || defaultImage">
      <img :src="article.img || defaultImage" :alt="article.title" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSiteInfoStore } from '@/stores/siteInfo'

const siteInfoStore = useSiteInfoStore()
const siteInfo = computed(() => siteInfoStore.siteInfo)

// 定义事件
const emit = defineEmits<{
  (e: 'share'): void
  (e: 'comment'): void
  (e: 'like'): void
  (e: 'subscribe'): void
}>()


// 定义Props
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  tags: string
  views: number
  createdAt: string
  updatedAt: string
}

interface Props {
  article: Article
}

const props = defineProps<Props>()

// 默认图片
const defaultImage = new URL('../../assets/img/无封面.jpg', import.meta.url).href

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 计算字数
const wordCount = computed(() => {
  if (!props.article.content) return 0
  // 去除HTML标签后计算字数
  const text = props.article.content.replace(/<[^>]*>/g, '')
  // 计算中文字符和英文单词
  const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || []
  const englishWords = text.match(/[a-zA-Z]+/g) || []
  return chineseChars.length + englishWords.length
})

// 计算预计阅读时间（中文约200字/分钟，英文约150词/分钟）
const readingTime = computed(() => {
  if (!props.article.content) return 0
  // 去除HTML标签后计算字数
  const text = props.article.content.replace(/<[^>]*>/g, '')
  // 分别计算中英文字符数
  const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || []
  const englishWords = text.match(/[a-zA-Z]+/g) || []
  
  // 中文每分钟200字，英文每分钟150词
  const chineseTime = chineseChars.length / 200
  const englishTime = englishWords.length / 150
  
  // 总时间，向上取整
  const totalMinutes = Math.ceil(chineseTime + englishTime)
  return Math.max(1, totalMinutes) // 至少1分钟
})

// 格式化数字，添加千位分隔符
const formatNumber = (num: number) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}
</script>

<style scoped>
.article-header {
  margin-bottom: 40px;
  padding-bottom: 0;
  border-bottom: none;
}

.article-title {
  font-size: 36px;
  font-weight: 800;
  margin-bottom: 25px;
  color: var(--color-heading);
  line-height: 1.4;
  letter-spacing: -0.5px;
}

.article-meta-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-end; /* Align to bottom so date aligns with stats */
  margin-bottom: 25px;
  flex-wrap: wrap;
  gap: 20px;
  padding-bottom: 25px;
  border-bottom: 1px solid var(--color-border);
}

.author-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--color-background);
  box-shadow: 0 0 0 1px var(--color-border);
}

.author-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.author-name {
  font-weight: 600;
  color: var(--color-heading);
  font-size: 15px;
}

.publish-date {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.meta-stats {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: var(--color-text-secondary);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: var(--color-background-soft); /* Light background */
  border: 1px solid var(--color-border);
  border-radius: 12px; /* Pill shape */
  font-size: 12px;
  font-weight: 500;
  transition: all 0.2s;
  cursor: default;
}

.meta-item:hover {
  background: var(--color-background-mute);
  border-color: var(--color-border-hover);
}

.meta-item svg {
  color: var(--color-text-secondary);
  opacity: 0.8;
}

/* Category specific style to make it pop */
.tag-item {
  color: var(--color-accent);
  border-color: var(--color-accent);
  background: transparent;
}

.tag-item:hover {
  background: var(--color-accent);
  color: white;
}

.tag-item svg {
  color: currentColor;
  opacity: 1;
}

.article-actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.left-actions, .right-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.action-btn-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  border-radius: 8px; /* Slightly squarer like screenshot */
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  height: 36px;
  min-width: 60px;
  justify-content: center;
}

.action-btn-pill:hover {
  background: var(--color-background-soft);
  border-color: var(--color-text-secondary);
}

.like-btn svg {
  color: var(--color-text-secondary);
}

.like-btn:hover svg {
  color: #ff6b6b;
  fill: #ff6b6b; /* Optionally fill on hover */
}

.subscribe-btn svg {
  color: var(--color-text-secondary);
}

/* Comment and Share buttons share the same pill style now */
.comment-btn, .share-btn {
  /* Inherits action-btn-pill */
}


.article-image {
  width: 100%;
  height: auto;
  /* Removed max-height to display full image */
  overflow: hidden;
  margin-bottom: 40px;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0,0,0,0.08);
}

.article-image img {
  width: 100%;
  /* Removed fixed height to maintain aspect ratio */
  height: auto;
  display: block;
  object-fit: contain; /* Ensure full image is visible */
  transition: transform 0.5s ease;
}

.article-image:hover img {
  transform: scale(1.02);
}

@media (max-width: 768px) {
  .article-title {
    font-size: 26px;
  }
  
  .article-meta-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
  
  .meta-stats {
    width: 100%;
    justify-content: flex-start;
    padding-left: 60px; /* Align with text not avatar */
    margin-top: -10px;
  }
  
  .article-image {
    max-height: 300px;
  }
}
</style>