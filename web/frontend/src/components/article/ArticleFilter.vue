<template>
  <div class="filters">
    <!-- 搜索框 -->
    <div class="filter-group">
      <label>搜索:</label>
      <input 
        type="text" 
        :value="searchKeyword" 
        @input="handleSearchInput"
        @keyup.enter="handleSearch"
        placeholder="输入文章标题或摘要关键词"
        class="search-input"
      />
      <button @click="handleSearch" class="search-button">搜索</button>
      <button @click="handleReset" class="reset-button">重置</button>
    </div>
    
    <!-- 分类筛选 -->
    <div class="filter-group">
      <label>分类:</label>
      <select :value="selectedCategory" @change="handleCategoryChange">
        <option value="">全部分类</option>
        <option 
          v-for="category in categories" 
          :key="category.id" 
          :value="category.id"
        >
          {{ category.name }}
        </option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

// 定义Props
interface Category {
  id: number
  name: string
}

interface Props {
  categories: Category[]
  selectedCategory: string
  searchKeyword: string
}

// 定义Emits
const emit = defineEmits<{
  (e: 'category-change', value: string): void
  (e: 'search', keyword: string): void
  (e: 'reset'): void
}>()

const props = defineProps<Props>()

// 防抖定时器
let searchTimer: number | null = null

// 处理分类变化
const handleCategoryChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  emit('category-change', target.value)
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

// 处理搜索（点击搜索按钮或按回车）
const handleSearch = (event: Event) => {
  // 清除防抖定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
  
  const target = event.target as HTMLInputElement
  const keyword = (target.previousElementSibling as HTMLInputElement)?.value || props.searchKeyword
  emit('search', keyword)
}

// 处理重置
const handleReset = () => {
  // 清除防抖定时器
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
  
  emit('reset')
}

// 组件卸载时清除定时器
onUnmounted(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
  }
})
</script>

<style scoped>
.filters {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 20px;
  margin-bottom: 25px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.filter-group:last-child {
  margin-bottom: 0;
}

.filter-group label {
  font-weight: 500;
  color: #333;
  min-width: 40px;
}

.filter-group select,
.filter-group input {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
}

.search-input {
  min-width: 200px;
}

.search-button,
.reset-button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
}

.search-button {
  background-color: #409eff;
  color: white;
}

.reset-button {
  background-color: #f5f5f5;
  color: #666;
}

.search-button:hover {
  background-color: #337ecc;
}

.reset-button:hover {
  background-color: #e0e0e0;
}

@media (max-width: 768px) {
  .filter-group {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .search-input {
    min-width: 150px;
  }
}
</style>