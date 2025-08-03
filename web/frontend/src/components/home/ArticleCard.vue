<template>
  <div class="article-card">
    <div class="article-image" v-if="article.img">
      <img :src="article.img" :alt="article.title" />
    </div>
    <div class="article-content">
      <h3 class="article-title">
        <router-link :to="`/article/${article.id}`">
          {{ article.title }}
        </router-link>
      </h3>
      <div class="article-meta">
        <span class="category">{{ article.categoryName }}</span>
        <span class="date">{{ formatDate(article.createdAt) }}</span>
      </div>
      <p class="article-excerpt">{{ article.desc || '暂无简介' }}</p>
      <router-link :to="`/article/${article.id}`" class="read-more">
        阅读全文 »
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
  createdAt: string
  updatedAt: string
}

interface Props {
  article: Article
}

const props = defineProps<Props>()

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.article-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.article-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
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

.article-card:hover .article-image img {
  transform: scale(1.05);
}

.article-content {
  padding: 25px;
  display: flex;
  flex-direction: column;
  height: 100%;
  flex: 1;
}

.article-title {
  font-size: 20px;
  margin: 0 0 15px 0;
}

.article-title a {
  color: #333;
  text-decoration: none;
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
  margin-bottom: 15px;
}

.article-excerpt {
  color: #666;
  line-height: 1.6;
  margin-bottom: 20px;
  flex-grow: 1;
}

.read-more {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
  margin-top: auto;
}

.read-more:hover {
  text-decoration: underline;
}
</style>