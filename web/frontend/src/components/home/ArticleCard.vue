<template>
  <div class="article-item">
    <div class="article-main">
      <div class="article-header">
        <h2 class="article-title">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="title-icon"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line><polyline points="10 9 9 9 8 9"></polyline></svg>
          <router-link :to="`/article/${article.id}`">
            {{ article.title }}
          </router-link>
        </h2>
      </div>
      
      <div class="article-meta-row">
        <div class="tags-container" v-if="article.tags">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tag-icon"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>
          <span v-for="tag in splitTags(article.tags)" :key="tag" class="tag-pill">
            {{ tag }}
          </span>
        </div>
        <span v-else class="category-pill">{{ article.categoryName }}</span>
      </div>

      <div class="article-summary">
        <p>{{ article.desc || '暂无简介' }}</p>
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
  top?: number
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
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

// 分割标签
const splitTags = (tags: string) => {
  if (!tags) return []
  return tags.replace(/，/g, ',').split(',').map(t => t.trim()).filter(t => t).slice(0, 3) // 最多显示3个标签
}
</script>

<style scoped>
.article-item {
  display: flex;
  flex-direction: row;
  padding: 25px 15px;
  margin-bottom: 0;
  background: transparent;
  border-bottom: 1px solid #f0f0f0;
  position: relative;
  transition: all 0.3s ease;
  height: 160px; /* 紧凑高度 */
  
  /* 对角线高亮效果 */
  background-repeat: no-repeat;
  background-image: 
    linear-gradient(#42b883, #42b883), /* 上边框 */
    linear-gradient(#42b883, #42b883), /* 右边框 */
    linear-gradient(#42b883, #42b883), /* 下边框 */
    linear-gradient(#42b883, #42b883); /* 左边框 */
    
  /* 初始大小为0 */
  background-size: 0% 2px, 2px 0%, 0% 2px, 2px 0%;
  
  /* 定位：左上，右下，右下，左上 */
  background-position: 0 0, 100% 100%, 100% 100%, 0 0;
}

.article-item:hover {
  background-size: 100% 2px, 2px 100%, 100% 2px, 2px 100%;
  background-color: rgba(66, 184, 131, 0.01);
}

/* 左侧内容区 */
.article-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  min-width: 0;
  padding-right: 20px;
}

/* 标题 */
.article-header {
  margin-bottom: 8px;
}

.article-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  line-height: 1.4;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #2c3e50;
}

.title-icon {
  color: #42b883; /* 主题色图标 */
  flex-shrink: 0;
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

/* Meta信息行 (标签等) */
.article-meta-row {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.tags-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-icon {
  color: #666;
  width: 14px;
  height: 14px;
}

.tag-pill {
  font-size: 12px;
  color: #42b883;
  border: 1px solid #42b883;
  padding: 1px 8px;
  border-radius: 12px; /* 胶囊形状 */
  background: transparent;
  transition: all 0.3s;
}

.tag-pill:hover {
  background: #42b883;
  color: white;
}

.category-pill {
  font-size: 12px;
  color: #666;
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 4px;
}

/* 摘要 */
.article-summary {
  flex: 1;
  overflow: hidden;
}

.article-summary p {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2; /* 显示2行 */
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 右侧封面图 */
.article-cover {
  width: 200px; /* 稍微缩小图片宽度 */
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.article-cover img {
  width: 100%;
  height: 110px; /* 固定高度，保持整齐 */
  object-fit: cover;
  border-radius: 8px;
  transition: transform 0.5s ease;
}

.article-item:hover .article-cover img {
  transform: scale(1.05);
}

/* 响应式适配 */
@media (max-width: 768px) {
  .article-item {
    flex-direction: column;
    height: auto;
    padding: 20px;
    border-bottom: 1px solid #eee;
  }
  
  /* 移动端取消对角线动画，改用简单背景色 */
  .article-item:hover {
    background-size: 0 0;
    background-color: #f9f9f9;
  }

  .article-cover {
    width: 100%;
    height: 160px;
    margin-top: 15px;
    order: 1; /* 图片在下方 */
  }
  
  .article-cover img {
    height: 100%;
  }

  .article-main {
    padding-right: 0;
  }
}
</style>