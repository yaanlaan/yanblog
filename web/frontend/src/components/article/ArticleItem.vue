<template>
  <div class="article-item">
    <div class="article-image">
      <img :src="article.img || defaultImage" :alt="article.title" />
    </div>
    <div class="article-content">
      <div class="article-header">
        <h2 class="article-title">
          <router-link :to="`/article/${article.id}`">
            {{ article.title }}
          </router-link>
        </h2>
        <div class="article-meta">
          <span class="category">
            分类: {{ article.categoryName }}
          </span>
          <span class="date">
            发布时间: {{ formatDate(article.createdAt) }}
          </span>
        </div>
      </div>
      
      <div class="article-summary">
        <p>{{ article.desc || '暂无简介' }}</p>
      </div>
      
      <div class="article-footer">
        <router-link :to="`/article/${article.id}`" class="read-more">
          阅读全文 »
        </router-link>
      </div>
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
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.article-item {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: box-shadow 0.3s;
  display: flex;
  flex-direction: column;
  position: relative; /* 保留相对定位，以防其他功能需要 */
}

.article-top-tag {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: #ff6b6b;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  z-index: 1;
}

.article-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.article-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.article-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.article-item:hover .article-image img {
  transform: scale(1.05);
}

.article-content {
  padding: 25px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.article-header {
  margin-bottom: 15px;
}

.article-title {
  font-size: 22px;
  margin: 0 0 10px 0;
}

.article-title a {
  text-decoration: none;
  color: #333;
  transition: color 0.3s;
}

.article-title a:hover {
  color: #007bff;
}

.article-meta {
  display: flex;
  gap: 15px;
  font-size: 14px;
  color: #888;
}

.article-summary p {
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.article-footer {
  margin-top: 20px;
}

.read-more {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
}

.read-more:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .article-title {
    font-size: 20px;
  }
  
  .article-image {
    height: 150px;
  }
}
</style>