<template>
  <div class="article-main-content">
    <div class="article-description" v-if="article.desc">
      <blockquote>{{ article.desc }}</blockquote>
    </div>
    
    <div class="content" v-html="renderedContent"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'

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

// 渲染Markdown内容
const renderedContent = computed(() => {
  if (!props.article.content) return ''
  return marked.parse(props.article.content)
})
</script>

<style scoped>
.article-description blockquote {
  font-size: 16px;
  color: #666;
  border-left: 4px solid #007bff;
  padding: 10px 20px;
  margin: 0 0 30px 0;
  background: #f8f9fa;
}

.content {
  font-size: 16px;
  line-height: 1.8;
}

.content :deep(h1) {
  font-size: 24px;
  margin: 24px 0 16px;
  color: #333;
}

.content :deep(h2) {
  font-size: 22px;
  margin: 22px 0 14px;
  color: #333;
}

.content :deep(h3) {
  font-size: 20px;
  margin: 20px 0 12px;
  color: #333;
}

.content :deep(p) {
  margin: 16px 0;
  color: #444;
}

.content :deep(ul),
.content :deep(ol) {
  padding-left: 30px;
  margin: 16px 0;
}

.content :deep(li) {
  margin-bottom: 8px;
}

.content :deep(code) {
  background: #f1f1f1;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
  font-size: 14px;
}

.content :deep(pre) {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
  margin: 20px 0;
}

.content :deep(pre code) {
  background: none;
  padding: 0;
}

.content :deep(blockquote) {
  border-left: 4px solid #ddd;
  padding: 10px 20px;
  margin: 20px 0;
  background: #f8f9fa;
  color: #666;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 20px 0;
}

.content :deep(a) {
  color: #007bff;
  text-decoration: none;
}

.content :deep(a:hover) {
  text-decoration: underline;
}

.content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: left;
}

.content :deep(th) {
  background: #f8f9fa;
  font-weight: bold;
}
</style>