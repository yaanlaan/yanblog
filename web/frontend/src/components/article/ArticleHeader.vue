，<template>
  <div class="article-header">
    <h1 class="article-title">{{ article.title }}</h1>
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
  font-size: 28px;
  color: #333;
  margin: 0 0 15px 0;
  line-height: 1.3;
}

.article-meta {
  display: flex;
  gap: 20px;
  font-size: 14px;
  color: #888;
}

@media (max-width: 768px) {
  .article-title {
    font-size: 28px;
  }
  
  .article-meta {
    flex-direction: column;
    gap: 5px;
  }
}
</style>