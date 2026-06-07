<template>
  <div class="article-item top-article">
    <div class="article-main">
      <div class="article-header">
        <h2 class="article-title">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="title-icon"><path d="M4 11a9 9 0 0 1 9 9"></path><path d="M4 4a16 16 0 0 1 16 16"></path><circle cx="5" cy="19" r="1"></circle></svg>
          <router-link :to="`/article/${article.id}`">
            {{ article.title }}
          </router-link>
        </h2>
      </div>
      
      <div class="article-summary">
        <p>{{ article.desc || '暂无简介' }}</p>
      </div>
      
      <div class="article-footer">
        <div class="footer-left">
          <span class="meta-item date">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            {{ formatDate(article.createdAt) }}
          </span>
          <span class="meta-item category">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>
            {{ article.categoryName }}
          </span>
          <span class="meta-item tags" v-if="article.tags">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
            <span v-for="(tag, index) in splitTags(article.tags)" :key="tag">
              {{ tag }}{{ index < splitTags(article.tags).length - 1 ? ', ' : '' }}
            </span>
          </span>
        </div>
        <div class="footer-right">
          <span class="top-badge">置顶</span>
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
  tags?: string
  createdAt: string
  updatedAt: string
}

interface Props {
  article: Article
}

defineProps<Props>()

// 默认图片
const defaultImage = new URL('../../assets/img/无封面.jpg', import.meta.url).href

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }).replace(/\//g, '-')
}

// 分割标签
const splitTags = (tags: string) => {
  if (!tags) return []
  return tags.replace(/，/g, ',').split(',').map(t => t.trim()).filter(t => t).slice(0, 3)
}
</script>

<style scoped>
/* ... existing styles ... */
.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 15px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.footer-left {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
  font-size: 0.85rem;
  color: #999;
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

.top-badge {
  background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  box-shadow: 0 2px 6px rgba(255, 77, 79, 0.2);
}
</style>

<style scoped>
.article-item {
  display: flex;
  flex-direction: row;
  background: var(--color-background-soft);
  border-radius: 0 12px 12px 0; /* 左侧无圆角 */
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.03);
  overflow: hidden;
  transition: all 0.3s ease;
  margin-bottom: 15px;
  height: 180px; /* 固定高度 */
  position: relative;
  border: 1px solid var(--color-border);
}

.article-item::before {
  content: "";
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: linear-gradient(180deg, var(--color-accent) 0%, #35495e 100%); /* 置顶文章使用绿色渐变 */
  z-index: 1;
}

.article-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.06);
}

/* 左侧内容区 */
.article-main {
  flex: 1;
  padding: 15px 20px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0; /* 防止flex子项溢出 */
}

/* 标题 */
.article-header {
  margin-bottom: 8px;
}

.article-title {
  font-size: 20px;
  font-weight: 900;
  margin: 0;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-heading);
}

.title-icon {
  color: var(--color-accent); /* 置顶图标颜色 */
  width: 18px;
  height: 18px;
}

.article-title a {
  text-decoration: none;
  color: inherit;
  transition: color 0.3s;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  font-weight: 900 !important;
  font-size: 20px !important;
}

.article-title a:hover {
  color: var(--color-accent);
}

/* 摘要 */
.article-summary {
  flex: 1;
  margin-bottom: 10px;
  overflow: hidden;
}

.article-summary p {
  color: var(--color-text);
  font-size: 14px;
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
  font-size: 13px;
  color: var(--color-text-light);
  border-top: 1px solid var(--color-border);
  padding-top: 10px;
}

.footer-left {
  display: flex;
  align-items: center;
}

.date {
  display: flex;
  align-items: center;
  gap: 5px;
}

.footer-right {
  display: flex;
  gap: 8px;
}

.category-badge {
  padding: 2px 10px;
  border: 1px solid #42b883;
  color: #42b883;
  border-radius: 4px;
  font-size: 12px;
  transition: all 0.3s;
  cursor: pointer;
}

.category-badge:hover {
  background-color: #42b883;
  color: white;
}

.top-badge {
  padding: 2px 10px;
  background: linear-gradient(45deg, #42b883, #66c798);
  color: white;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  box-shadow: 0 2px 6px rgba(66, 184, 131, 0.3);
}

/* 右侧封面图 */
.article-cover {
  width: 240px;
  height: 100%;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.article-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 6px;
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
    border-radius: 12px;
  }
  
  .article-item::before {
    width: 100%;
    height: 4px;
    bottom: auto;
    background: linear-gradient(90deg, #42b883 0%, #35495e 100%);
  }

  .article-cover {
    width: 100%;
    height: 160px;
    padding: 0;
    order: -1; /* 图片在上方 */
    margin-top: 4px;
  }
  
  .article-cover img {
    border-radius: 0;
  }

  .article-main {
    padding: 15px;
  }

  .article-title {
    font-size: 16px;
  }
  
  .article-summary p {
    -webkit-line-clamp: 2;
  }
}
</style>