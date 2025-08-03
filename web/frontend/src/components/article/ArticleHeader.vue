，<template>
  <div class="article-header">
    <h1 class="article-title">{{ article.title }}</h1>
    <div class="article-image" v-if="article.img || defaultImage">
      <img :src="article.img || defaultImage" :alt="article.title" />
    </div>
    <div class="article-meta">
      <span class="category">
        分类: {{ article.categoryName }}
      </span>
      <span class="date">
        发布时间: {{ formatDate(article.createdAt) }}
      </span>
      <span class="date">
        更新时间: {{ formatDate(article.updatedAt) }}
      </span>
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
.article-header {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.article-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 20px;
  color: #333;
  line-height: 1.3;
}

.article-image {
  width: 100%;
  height: 300px;
  overflow: hidden;
  margin-bottom: 20px;
  border-radius: 8px;
}

.article-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-meta {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #666;
}

@media (max-width: 768px) {
  .article-title {
    font-size: 24px;
  }
  
  .article-meta {
    flex-direction: column;
    gap: 5px;
  }
  
  .article-image {
    height: 200px;
  }
}
</style>