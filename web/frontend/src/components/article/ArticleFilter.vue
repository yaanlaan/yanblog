<template>
  <div class="filters">
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
// 定义Props
interface Category {
  id: number
  name: string
}

interface Props {
  categories: Category[]
  selectedCategory: string
}

// 定义Emits
const emit = defineEmits<{
  (e: 'category-change', value: string): void
}>()

const props = defineProps<Props>()

// 处理分类变化
const handleCategoryChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  emit('category-change', target.value)
}
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
}

.filter-group label {
  font-weight: 500;
  color: #333;
}

.filter-group select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
}

@media (max-width: 768px) {
  .filter-group {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>