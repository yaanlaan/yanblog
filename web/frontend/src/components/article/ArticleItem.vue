<template>
  <div class="article-card" :class="{ 'is-list-mode': viewMode === 'list' }">
    <!-- 封面图 -->
    <div class="card-cover">
      <router-link :to="`/article/${article.id}`">
        <img :src="article.img || defaultImage" :alt="article.title" loading="lazy" />
      </router-link>
    </div>

    <!-- 内容区域 -->
    <div class="card-content">
      <!-- 标题 -->
      <h3 class="card-title">
        <el-tooltip
          effect="dark"
          :content="article.title"
          placement="top-start"
          :show-after="200"
        >
          <router-link class="article-link" :to="`/article/${article.id}`" v-html="highlightText(article.title)">
          </router-link>
        </el-tooltip>
      </h3>

      <!-- 描述 (仅列表模式显示) -->
      <p class="card-desc" v-if="viewMode === 'list'" v-html="highlightText(article.desc || '暂无简介')"></p>

      <!-- Meta信息 -->
      <div class="card-meta-row">
        <span class="meta-item category">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
          {{ article.categoryName }}
        </span>
        <div class="meta-item tags" v-if="article.tags">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
          <span v-for="(tag, index) in splitTags(article.tags)" :key="tag">
            {{ tag }}{{ index < splitTags(article.tags).length - 1 ? ' ' : '' }}
          </span>
        </div>
      </div>

      <!-- 底部信息 -->
      <div class="card-footer">
        <span class="date-text">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
          {{ formatDate(article.updatedAt || article.createdAt) }}
        </span>
        <span class="views-text">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
          {{ article.views || 0 }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { ElTooltip } from 'element-plus'

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
  views: number
  createdAt: string
  updatedAt: string
}

interface Props {
  article: Article
  viewMode?: 'grid' | 'list'
}

const props = withDefaults(defineProps<Props>(), {
  viewMode: 'grid'
})

const route = useRoute()

// 默认图片
const defaultImage = new URL('../../assets/img/无封面.jpg', import.meta.url).href

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }).replace(/\//g, '年').replace(/,/, '日')
}

// 分割标签
const splitTags = (tags: string) => {
  if (!tags) return []
  return tags.replace(/，/g, ',').split(',').map(t => t.trim()).filter(t => t).slice(0, 3) // 最多显示3个标签
}

// 搜索高亮
const highlightText = (text: string) => {
  const keyword = route.query.search as string
  if (!keyword || !text) return text
  
  // 转义特殊字符，防止正则错误
  const safeKeyword = keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(`(${safeKeyword})`, 'gi')
  return text.replace(regex, '<span class="search-highlight">$1</span>')
}
</script>

<style scoped>
.article-card {
  background: var(--color-background-soft);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px var(--color-shadow);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  height: 100%;
  border: 1px solid var(--color-border);
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12); /* Slightly stronger shadow on hover is okay, or use var if available */
}

.card-cover {
  height: 180px;
  overflow: hidden;
  position: relative;
}

.card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.article-card:hover .card-cover img {
  transform: scale(1.05);
}

.card-content {
  padding: 20px;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.card-title {
  margin: 0 0 15px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

:deep(.article-link) {
  color: var(--color-heading);
  text-decoration: none;
  transition: color 0.2s;
  font-size: 20px !important;
  font-weight: 900 !important;
  line-height: 1.4;
  letter-spacing: -0.5px;
  display: block; /* 确保占满父容器 */
}

:deep(.article-link:hover) {
  color: var(--color-accent);
}

.card-meta-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 20px;
  flex-grow: 1;
  font-size: 12px;
  color: var(--color-text-light);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  transition: color 0.3s ease;
  cursor: pointer;
}

.meta-item:hover {
  color: #42b883;
}

.meta-item svg {
  width: 14px;
  height: 14px;
}

.tags span {
  margin-right: 4px;
}

.tags span:last-child {
  margin-right: 0;
}

.card-footer {
  border-top: 1px solid var(--color-border);
  padding-top: 15px;
  display: flex;
  justify-content: center;
}

.date-text {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-family: monospace;
  display: flex;
  align-items: center;
  gap: 5px;
}

/* 列表模式样式 */
.article-card.is-list-mode {
  flex-direction: row;
  height: 200px;
}

.article-card.is-list-mode .card-cover {
  width: 320px;
  height: 100%;
  flex-shrink: 0;
}

.article-card.is-list-mode .card-content {
  padding: 25px;
  justify-content: center;
}

.article-card.is-list-mode .card-title {
  margin-bottom: 10px;
  height: auto;
  -webkit-line-clamp: 1;
}

.article-card.is-list-mode :deep(.article-link) {
  font-size: 28px !important;
}

.card-desc {
  color: var(--color-text-secondary);
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-card.is-list-mode .card-meta-row {
  margin-bottom: 15px;
  flex-grow: 0;
}

.article-card.is-list-mode .card-footer {
  padding-top: 0;
  border-top: none;
  justify-content: flex-start;
  gap: 20px;
}

.views-text {
  display: flex;
  align-items: center;
  gap: 4px;
}

.views-text svg {
  width: 14px;
  height: 14px;
}

@media (max-width: 768px) {
  .article-card.is-list-mode {
    flex-direction: column;
    height: auto;
  }
  
  .article-card.is-list-mode .card-cover {
    width: 100%;
    height: 180px;
  }
  
  .article-card.is-list-mode .card-title {
    font-size: 16px;
  }
}

/* 全局样式 */
:deep(.search-highlight) {
  color: #e74c3c;
  font-weight: bold;
  background-color: rgba(231, 76, 60, 0.1);
  padding: 0 2px;
  border-radius: 2px;
}
</style>