<template>
  <div class="article-filter-bar">
    <!-- 左侧分类Tab -->
    <div class="category-tabs">
      <button 
        class="tab-item" 
        :class="{ active: !selectedCategory }"
        @click="selectCategory('')"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tab-icon"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
        全部文章
      </button>
      <button 
        v-for="category in categories" 
        :key="category.id" 
        class="tab-item"
        :class="{ active: selectedCategory == String(category.id) }"
        @click="selectCategory(String(category.id))"
      >
        {{ category.name }}
      </button>
    </div>

    <!-- 右侧搜索与视图切换 -->
    <div class="right-actions">
      <div class="view-toggles">
        <button 
          class="view-btn" 
          :class="{ active: viewMode === 'grid' }"
          @click="$emit('view-change', 'grid')"
          title="网格视图"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
        </button>
        <button 
          class="view-btn" 
          :class="{ active: viewMode === 'list' }"
          @click="$emit('view-change', 'list')"
          title="列表视图"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
        </button>
      </div>

      <div class="search-wrapper">
        <div class="search-box">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="search-icon"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          <input 
            type="text" 
            :value="searchKeyword" 
            @input="handleSearchInput"
            @keyup.enter="handleSearch"
            placeholder="搜索文章..."
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onUnmounted } from 'vue'

// 定义Props
interface Category {
  id: number
  name: string
}

interface Props {
  categories: Category[]
  selectedCategory: string
  searchKeyword: string
  viewMode?: 'grid' | 'list'
}

// 定义Emits
const emit = defineEmits<{
  (e: 'category-change', value: string): void
  (e: 'search', keyword: string): void
  (e: 'reset'): void
  (e: 'view-change', mode: 'grid' | 'list'): void
}>()

const props = defineProps<Props>()

// 防抖定时器
let searchTimer: number | null = null

// 处理分类选择
const selectCategory = (id: string) => {
  emit('category-change', id)
}

// 处理搜索输入
const handleSearchInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const keyword = target.value
  
  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
  
  // 设置新的防抖定时器
  searchTimer = window.setTimeout(() => {
    emit('search', keyword)
  }, 300) // 300ms防抖延迟
}

// 处理搜索（按回车）
const handleSearch = (event: Event) => {
  // 清除防抖定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
  
  const target = event.target as HTMLInputElement
  emit('search', target.value)
}

// 组件卸载时清除定时器
onUnmounted(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})
</script>

<style scoped>
.article-filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  background: var(--color-background-soft);
  border: 1px solid var(--color-border);
  padding: 15px;
  border-radius: 8px;
  flex-wrap: wrap;
  gap: 15px;
  box-shadow: 0 2px 8px var(--color-shadow);
}

.category-tabs {
  display: flex;
  gap: 20px;
  overflow-x: auto;
  padding-bottom: 5px; /* 滚动条空间 */
}

.tab-item {
  background: none;
  border: none;
  font-size: 15px;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  font-weight: 500;
}

.tab-item:hover {
  color: var(--color-accent);
  background-color: rgba(66, 184, 131, 0.05); /* Usually accent color with opacity is fine */
}

.tab-item.active {
  color: var(--color-accent);
  background-color: rgba(66, 184, 131, 0.1);
  font-weight: 600;
}

.tab-icon {
  opacity: 0.8;
}

.right-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

.view-toggles {
  display: flex;
  background: var(--color-background-soft);
  padding: 3px;
  border-radius: 8px;
  gap: 2px;
}

.view-btn {
  border: none;
  background: none;
  padding: 6px;
  border-radius: 6px;
  color: var(--color-text-light);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.view-btn:hover {
  color: var(--color-text);
  background: var(--color-background-mute);
}

.view-btn.active {
  background: var(--color-background-mute);
  color: var(--color-accent);
  box-shadow: 0 2px 5px var(--color-shadow);
}

.search-wrapper {
  display: flex;
  align-items: center;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  color: var(--color-text-light);
  pointer-events: none;
}

.search-box input {
  padding: 8px 12px 8px 36px;
  border: 1px solid var(--color-border);
  background: var(--color-background-soft);
  border-radius: 20px;
  font-size: 14px;
  width: 200px;
  transition: all 0.3s;
  outline: none;
  color: var(--color-text);
}

.search-box input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.1);
  width: 240px;
}

@media (max-width: 768px) {
  .article-filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  
  .search-box input {
    width: 100%;
  }
  
  .search-box input:focus {
    width: 100%;
  }
}
</style>