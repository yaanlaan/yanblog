<template>
  <div class="top-article-card">
    <div class="article-content">
      <div class="article-image" v-if="article.img || defaultImage">
        <img :src="article.img || defaultImage" :alt="article.title" />
      </div>
      <div class="article-info">
        <h3 class="article-title">
          <router-link :to="`/article/${article.id}`">
            {{ article.title }}
          </router-link>
        </h3>
        <div class="article-meta">
          <span class="category">{{ article.categoryName }}</span>
          <span class="date">{{ formatDate(article.createdAt) }}</span>
          <span class="top-tag">置顶</span>
        </div>
        <p class="article-excerpt">{{ article.desc || '暂无简介' }}</p>
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

const props = defineProps<Props>()

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
.top-article-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  margin-bottom: 20px;
}

.top-article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.article-content {
  display: flex;
  height: 200px;
}

.article-image {
  width: 35%;
  height: 100%;
  overflow: hidden;
}

.article-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.top-article-card:hover .article-image img {
  transform: scale(1.05);
}

.article-info {
  width: 65%;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.article-title {
  font-size: 20px;
  margin: 0 0 10px 0;
  flex: 1;
}

.article-title a {
  color: #333;
  text-decoration: none;
}

.article-title a:hover {
  color: #007bff;
}

.article-meta {
  display: flex;
  gap: 15px;
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.top-tag {
  background-color: #dc3545;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.article-excerpt {
  color: #666;
  line-height: 1.6;
  margin: 0 0 15px 0;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.read-more {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
  align-self: flex-start;
}

.read-more:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .article-content {
    flex-direction: column;
    height: auto;
  }
  
  .article-image {
    width: 100%;
    height: 200px;
  }
  
  .article-info {
    width: 100%;
  }
}
</style>