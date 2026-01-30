<template>
  <div class="article-toc" v-if="tocItems.length > 0">
    <div class="toc-header" title="收起目录栏">
      <div class="toc-title-wrapper">
        <!-- List Icon SVG -->
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="toc-icon"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
        <h3 class="toc-title">目录</h3>
      </div>
      <button class="toc-close-btn" @click.stop="$emit('close')" title="关闭">
        <!-- Close Icon SVG -->
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
      </button>
    </div>
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

const emit = defineEmits<{
  (e: 'close'): void
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
  background: var(--color-background-soft);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 var(--color-shadow);
  max-height: calc(100vh - 100px);
  display: flex;
  flex-direction: column;
  border: 1px solid var(--color-border);
}

.toc-header {
  display: flex;
  justify-content: space-between; /* Space out title and button */
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 0;
  border-bottom: none;
  cursor: default;
  user-select: none;
}

.toc-title-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--color-accent);
}

.toc-close-btn {
  background: transparent;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.toc-close-btn:hover {
  background-color: var(--color-background-soft);
  color: var(--color-text);
}

.toc-close-btn i {
  font-size: 16px;
}

.toc-icon {
  font-size: 18px;
  font-weight: bold;
}

.toc-title {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  padding: 0;
  border: none;
  color: var(--color-heading);
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
  flex: 1;
  scrollbar-width: none; /* Firefox */
  position: relative;
}

.toc-list::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: var(--color-border);
  z-index: 0;
}

.toc-list::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}

.toc-item {
  margin-bottom: 0;
  line-height: 1.5;
  cursor: pointer;
  border-left: 2px solid transparent;
  transition: all 0.2s;
  border-radius: 0;
  position: relative;
  z-index: 1;
  padding-left: 15px; /* Base padding */
  margin-left: 0;
}

.toc-item-level-1 {
  font-weight: 600;
  font-size: 15px;
  margin-top: 8px;
}

.toc-item-level-2 {
  padding-left: 28px; /* Indent for level 2 */
  font-size: 14px;
  color: var(--color-text-secondary);
}

.toc-link {
  color: var(--color-text-secondary);
  text-decoration: none;
  display: block;
  padding: 6px 0;
  transition: all 0.2s;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toc-item:hover .toc-link {
  color: var(--color-heading);
}

.toc-item.active {
  border-left-color: var(--color-accent);
  background-color: transparent; 
}

.toc-item.active .toc-link {
  color: var(--color-accent);
  font-weight: 600;
}
</style>