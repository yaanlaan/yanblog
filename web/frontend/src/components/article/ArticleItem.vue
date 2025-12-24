<template>
  <div class="article-item">
    <div class="article-main">
      <div class="article-header">
        <h2 class="article-title">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="title-icon"><path d="M4 11a9 9 0 0 1 9 9"></path><path d="M4 4a16 16 0 0 1 16 16"></path><circle cx="5" cy="19" r="1"></circle></svg>
          <router-link :to="`/article/${article.id}`" v-html="highlightText(article.title)">
          </router-link>
        </h2>
      </div>
      
      <div class="article-summary">
        <p v-html="highlightText(article.desc || '暂无简介')"></p>
      </div>
      
      <div class="article-footer">
        <div class="footer-left">
          <span class="date">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            {{ formatDate(article.createdAt) }}
          </span>
        </div>
        <div class="footer-right">
          <div class="tags-container" v-if="article.tags">
            <span v-for="tag in splitTags(article.tags)" :key="tag" class="tag-badge">
              {{ tag }}
            </span>
          </div>
          <span v-else class="category-badge">{{ article.categoryName }}</span>
        </div>
      </div>
    </div>

    <div class="article-cover">
      <router-link :to="`/article/${article.id}`">
        <img :src="article.img || defaultImage" :alt="article.title" />
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

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
}

defineProps<Props>()

const route = useRoute()

// 默认图片
const defaultImage = new URL('../../assets/img/无封面.jpg', import.meta.url).href

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
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

<style>
/* 全局样式 */
.search-highlight {
  color: #e74c3c;
  font-weight: bold;
  background-color: rgba(231, 76, 60, 0.1);
  padding: 0 2px;
  border-radius: 2px;
}
</style>

<style scoped>
.article-item {
  display: flex;
  flex-direction: row;
  background: #fdfdfd; /* 极淡的背景色 */
  border-radius: 0 12px 12px 0; /* 左侧无圆角 */
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
  overflow: hidden;
  transition: all 0.3s ease;
  margin-bottom: 20px;
  height: 220px; /* 固定高度 */
  position: relative;
}

.article-item::before {
  content: "";
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 5px;
  background: linear-gradient(180deg, #42b883 0%, #35495e 100%); /* 渐变色 */
  z-index: 1;
}

.article-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

/* 左侧内容区 */
.article-main {
  flex: 1;
  padding: 25px 30px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0; /* 防止flex子项溢出 */
}

/* 标题 */
.article-header {
  margin-bottom: 10px;
}

.article-title {
  font-size: 22px;
  font-weight: 600;
  margin: 0;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #2c3e50;
}

.title-icon {
  color: #e74c3c; /* 图标颜色 */
  font-size: 24px;
}

.article-title a {
  text-decoration: none;
  color: inherit;
  transition: color 0.3s;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-title a:hover {
  color: #42b883;
}

/* 摘要 */
.article-summary {
  flex: 1;
  margin-bottom: 15px;
  overflow: hidden;
}

.article-summary p {
  color: #666;
  font-size: 15px;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3; /* 显示3行 */
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 底部信息 */
.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #999;
  border-top: 1px solid rgba(0,0,0,0.03);
  padding-top: 15px;
}

.footer-left {
  display: flex;
  align-items: center;
}

.date {
  display: flex;
  align-items: center;
  gap: 6px;
}

.date i {
  font-size: 16px;
  color: #42b883;
}

.footer-right {
  display: flex;
  gap: 10px;
}

.tag-badge, .category-badge {
  padding: 4px 12px;
  border: 1px solid #42b883;
  color: #42b883;
  border-radius: 6px;
  font-size: 12px;
  transition: all 0.3s;
  cursor: pointer;
}

.tag-badge:hover, .category-badge:hover {
  background-color: #42b883;
  color: white;
}

/* 右侧封面图 */
.article-cover {
  width: 280px;
  height: 100%;
  padding: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  transition: transform 0.5s ease;
}

.article-item:hover .article-cover img {
  transform: scale(1.03);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .article-item {
    flex-direction: column;
    height: auto;
    border-left: none;
    border-top: 4px solid #42b883;
  }

  .article-cover {
    width: 100%;
    height: 180px;
    padding: 0;
    order: -1; /* 图片在上方 */
  }
  
  .article-cover img {
    border-radius: 0;
  }

  .article-main {
    padding: 20px;
  }

  .article-title {
    font-size: 18px;
  }
  
  .article-summary p {
    -webkit-line-clamp: 2;
  }
}
</style>