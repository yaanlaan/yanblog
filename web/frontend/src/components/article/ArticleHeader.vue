<template>
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
      <span class="tags" v-if="article.tags">
        标签: {{ article.tags }}
      </span>
      <span class="word-count">
        字数: {{ formatNumber(wordCount) }}
      </span>
      <span class="reading-time">
        预计阅读: {{ readingTime }} 分钟
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// 定义Props
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  tags: string
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

// 计算字数
const wordCount = computed(() => {
  if (!props.article.content) return 0
  // 去除HTML标签后计算字数
  const text = props.article.content.replace(/<[^>]*>/g, '')
  // 计算中文字符和英文单词
  const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || []
  const englishWords = text.match(/[a-zA-Z]+/g) || []
  return chineseChars.length + englishWords.length
})

// 计算预计阅读时间（中文约200字/分钟，英文约150词/分钟）
const readingTime = computed(() => {
  if (!props.article.content) return 0
  // 去除HTML标签后计算字数
  const text = props.article.content.replace(/<[^>]*>/g, '')
  // 分别计算中英文字符数
  const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || []
  const englishWords = text.match(/[a-zA-Z]+/g) || []
  
  // 中文每分钟200字，英文每分钟150词
  const chineseTime = chineseChars.length / 200
  const englishTime = englishWords.length / 150
  
  // 总时间，向上取整
  const totalMinutes = Math.ceil(chineseTime + englishTime)
  return Math.max(1, totalMinutes) // 至少1分钟
})

// 格式化数字，添加千位分隔符
const formatNumber = (num: number) => {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
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

.article-meta .word-count,
.article-meta .reading-time {
  background-color: #f0f8ff;
  padding: 2px 8px;
  border-radius: 4px;
  color: #409eff;
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