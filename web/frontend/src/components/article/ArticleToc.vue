<template>
  <div class="article-toc" v-if="tocItems.length > 0">
    <h3 class="toc-title">目录</h3>
    <ul class="toc-list" ref="tocListRef">
      <li 
        v-for="item in tocItems" 
        :key="item.id"
        :class="['toc-item', `toc-item-level-${item.level}`, { 'active': activeId === item.id }]"
        @click.prevent="scrollToHeading(item.id)"
      >
        <a :href="`#${item.id}`" class="toc-link" :title="item.text">{{ item.text }}</a>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'

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
const tocListRef = ref<HTMLElement | null>(null)

// 从Markdown内容中提取标题 (仅H1和H2)
const extractHeaders = () => {
  const headers: TocItem[] = []
  const headingRegex = /^(#{1,6})\s+(.+)$/gm
  let match
  let counter = 0
  
  while ((match = headingRegex.exec(props.content)) !== null) {
    counter++
    const level = match[1].length
    const text = match[2].trim()
    const id = `heading-${counter}`
    
    // 只保留 H1 和 H2
    if (level <= 2) {
      headers.push({
        id,
        text,
        level
      })
    }
  }
  
  tocItems.value = headers
}

// 滚动到指定标题
const scrollToHeading = (id: string) => {
  const element = document.getElementById(id)
  if (element) {
    const headerOffset = 80
    const elementPosition = element.getBoundingClientRect().top
    const offsetPosition = elementPosition + window.pageYOffset - headerOffset

    window.scrollTo({
      top: offsetPosition,
      behavior: "smooth"
    })
    
    activeId.value = id
  }
}

// 监听滚动事件，更新当前活动的标题
const handleScroll = () => {
  if (tocItems.value.length === 0) return
  
  const scrollPosition = window.scrollY + 100
  let currentId = ''
  
  for (const item of tocItems.value) {
    const element = document.getElementById(item.id)
    if (element && element.offsetTop <= scrollPosition) {
      currentId = item.id
    }
  }
  
  if (currentId && currentId !== activeId.value) {
    activeId.value = currentId
    scrollToActiveTocItem()
  }
}

// 滚动目录列表以保持当前项可见
const scrollToActiveTocItem = () => {
  if (!tocListRef.value) return
  
  const activeItem = tocListRef.value.querySelector('.toc-item.active') as HTMLElement
  if (activeItem) {
    const listHeight = tocListRef.value.clientHeight
    const itemTop = activeItem.offsetTop
    const itemHeight = activeItem.clientHeight
    const scrollTop = tocListRef.value.scrollTop
    
    // 如果激活项不在可视区域内，则滚动
    if (itemTop < scrollTop || itemTop + itemHeight > scrollTop + listHeight) {
       tocListRef.value.scrollTo({
         top: itemTop - listHeight / 2 + itemHeight / 2,
         behavior: 'smooth'
       })
    }
  }
}

watch(() => props.content, () => {
  extractHeaders()
})

onMounted(() => {
  extractHeaders()
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.article-toc {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  max-height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;
}

.toc-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex: 1;
  scrollbar-width: none; /* Firefox */
}

.toc-list::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.toc-item {
  margin-bottom: 4px;
  line-height: 1.5;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition: all 0.2s;
  border-radius: 0 4px 4px 0;
}

.toc-item-level-1 {
  padding-left: 10px;
  font-weight: 600;
}

.toc-item-level-2 {
  padding-left: 25px;
  font-size: 0.95em;
}

.toc-link {
  color: #666;
  text-decoration: none;
  display: block;
  padding: 6px 0;
  transition: all 0.2s;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toc-item:hover .toc-link {
  color: #42b883;
}

.toc-item.active {
  border-left-color: #42b883;
  background-color: rgba(66, 184, 131, 0.08);
}

.toc-item.active .toc-link {
  color: #42b883;
  font-weight: 500;
}
</style>