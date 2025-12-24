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

      <!-- 标签 -->
      <div class="card-tags">
        <template v-if="article.tags">
          <span v-for="tag in splitTags(article.tags)" :key="tag" class="tag-pill">
            <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tag-icon"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
            {{ tag }}
          </span>
        </template>
        <span v-else class="category-pill">{{ article.categoryName }}</span>
      </div>

      <!-- 底部信息 -->
      <div class="card-footer">
        <span class="date-text">Modify {{ formatDate(article.updatedAt || article.createdAt) }}</span>
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
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  height: 100%;
  border: 1px solid rgba(0,0,0,0.05);
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
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
  color: #2c3e50;
  text-decoration: none;
  transition: color 0.2s;
  font-size: 20px !important;
  font-weight: 900 !important;
  line-height: 1.4;
  letter-spacing: -0.5px;
  display: block; /* 确保占满父容器 */
}

:deep(.article-link:hover) {
  color: #42b883;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 20px;
  flex-grow: 1;
}

.tag-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background-color: #e0f2f1;
  color: #00695c;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.tag-icon {
  opacity: 0.7;
}

.category-pill {
  padding: 4px 10px;
  background-color: #f3e5f5;
  color: #7b1fa2;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.card-footer {
  border-top: 1px solid #f0f0f0;
  padding-top: 15px;
  display: flex;
  justify-content: center;
}

.date-text {
  font-size: 12px;
  color: #999;
  font-family: monospace;
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
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 15px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-card.is-list-mode .card-tags {
  margin-bottom: 15px;
  flex-grow: 0;
}

.article-card.is-list-mode .card-footer {
  padding-top: 0;
  border-top: none;
  justify-content: flex-start;
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