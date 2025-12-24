<template>
  <div class="sidebar-card category-card">
    <div class="card-header">
      <h3><i class="iconfont icon-categories" style="color: #e6a23c; margin-right: 5px;"></i> 分类</h3>
      <router-link to="/categories" class="more-link">更多 ></router-link>
    </div>
    <div class="card-content">
      <div class="category-grid">
        <router-link 
          v-for="cat in categories" 
          :key="cat.id" 
          :to="`/categories/${cat.id}`"
          class="category-item"
        >
          <span class="cat-name">{{ cat.name }}</span>
          <span class="cat-count">{{ cat.count || 0 }}</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { categoryApi } from '@/services/api'

interface Category {
  id: number
  name: string
  count?: number 
}

const categories = ref<Category[]>([])

const fetchCategories = async () => {
  try {
    const res = await categoryApi.getCategories({ pagesize: 6, pagenum: 1 })
    if (res.data.status === 200) {
      categories.value = res.data.data
    }
  } catch (error) {
    console.error('Failed to fetch categories', error)
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.more-link {
  font-size: 12px;
  color: #999;
  text-decoration: none;
}

.more-link:hover {
  color: #42b883;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f9f9f9;
  border-radius: 6px;
  text-decoration: none;
  color: #666;
  transition: all 0.3s;
  border: 1px solid transparent;
}

.category-item:hover {
  background: white;
  border-color: #42b883;
  color: #42b883;
  box-shadow: 0 2px 8px rgba(66, 184, 131, 0.1);
}

.cat-name {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 80px; /* Adjust based on layout */
}

.cat-count {
  font-size: 12px;
  background: #eee;
  padding: 2px 6px;
  border-radius: 10px;
  color: #999;
}

.category-item:hover .cat-count {
  background: #42b883;
  color: white;
}
</style>
