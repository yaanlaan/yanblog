<template>
  <div class="article-toc" v-if="tocItems.length > 0">
    <h3 class="toc-title">目录</h3>
    <ul class="toc-list">
      <li 
        v-for="item in tocItems" 
        :key="item.id"
        :class="['toc-item', `toc-item-level-${item.level}`, { 'active': activeId === item.id }]"
        @click="scrollToHeading(item.id)"
      >
        <a href="javascript:void(0)" class="toc-link">{{ item.text }}</a>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

interface TocItem {
  id: string
  text: string
  level: number
}

const props = defineProps<{
  content: string
}>()

const tocItems = ref<TocItem[]>([])
const activeId = ref<string>('')

// 从Markdown内容中提取标题
const extractHeaders = () => {
  const headers: TocItem[] = []
  const headingRegex = /^(#{1,6})\s+(.+)$/gm
  let match
  
  while ((match = headingRegex.exec(props.content)) !== null) {
    const level = match[1].length
    const text = match[2].trim()
    const id = `heading-${headers.length + 1}`
    
    headers.push({
      id,
      text,
      level
    })
  }
  
  tocItems.value = headers
}

// 滚动到指定标题
const scrollToHeading = (id: string) => {
  const element = document.getElementById(id)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

// 监听滚动事件，更新当前活动的标题
const handleScroll = () => {
  if (tocItems.value.length === 0) return
  
  const headings = tocItems.value.map(item => {
    return document.getElementById(item.id)
  }).filter(Boolean) as HTMLElement[]
  
  if (headings.length === 0) return
  
  const scrollPosition = window.scrollY + 100
  
  for (let i = headings.length - 1; i >= 0; i--) {
    const heading = headings[i]
    if (heading.offsetTop <= scrollPosition) {
      activeId.value = tocItems.value[i].id
      break
    }
  }
}

// 为文章内容添加ID
const addIdsToHeadings = () => {
  // 这个功能将在ArticleContent组件中实现
}

onMounted(() => {
  extractHeaders()
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 监听内容变化
defineExpose({
  extractHeaders
})
</script>

<style scoped>
.article-toc {
  position: sticky;
  top: 20px;
  padding: 20px;
  background: var(--color-background-soft);
  border-radius: 8px;
  max-height: calc(100vh - 40px);
  overflow-y: auto;
}

.toc-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0 0 15px 0;
  color: var(--color-heading);
  padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border);
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.toc-item {
  margin-bottom: 8px;
  padding: 5px 0;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.toc-item:hover {
  background-color: var(--color-background-mute);
}

.toc-item.active {
  background-color: #3d96fc;
}

.toc-item.active .toc-link {
  color: white;
}

.toc-link {
  text-decoration: none;
  color: var(--color-text);
  font-size: 14px;
  display: block;
  padding: 3px 10px;
  border-radius: 4px;
  transition: color 0.2s;
}

.toc-item-level-1 .toc-link {
  padding-left: 10px;
  font-weight: bold;
}

.toc-item-level-2 .toc-link {
  padding-left: 25px;
}

.toc-item-level-3 .toc-link {
  padding-left: 40px;
}

.toc-item-level-4 .toc-link {
  padding-left: 55px;
}

.toc-item-level-5 .toc-link {
  padding-left: 70px;
}

.toc-item-level-6 .toc-link {
  padding-left: 85px;
}
</style>