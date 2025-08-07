<template>
  <div class="sidebar-card tag-cloud">
    <div class="card-header">
      <h3>标签云</h3>
    </div>
    <div class="card-content">
      <div v-if="loading" class="skeleton-loader">
        <div class="skeleton-header"></div>
        <div class="skeleton-body">
          <div class="skeleton-tag"></div>
          <div class="skeleton-tag"></div>
          <div class="skeleton-tag"></div>
        </div>
      </div>
      <div class="tags" v-else-if="categories.length > 0">
        <router-link
          v-for="category in displayedCategories" 
          :key="category.id" 
          :to="`/category/${category.id}`"
          class="tag"
          :style="{ fontSize: calculateFontSize(category.articleCount) }"
        >
          {{ category.name }}
        </router-link>
        
        <div class="see-more-container">
          <button 
            v-if="categories.length > TAG_LIMIT" 
            @click="toggleShowAll"
            class="see-more-button"
          >
            {{ showAll ? 'tidy-display' : 'seemore' }}
          </button>
        </div>
      </div>
      <div class="error-message" v-else-if="error">
        <p>❌ {{ error }}</p>
        <button @click="onRetry" class="retry-button">重试</button>
      </div>
      <div class="empty-state" v-else>
        <p>暂无标签</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { categoryApi } from '@/services/api'

// 定义分类接口
interface Category {
  id: number
  name: string
  articleCount: number
}

// 标签显示上限
const TAG_LIMIT = 5

const categories = ref<Category[]>([])
const loading = ref(false)
const error = ref('')
const showAll = ref(false)

// 计算显示的标签
const displayedCategories = computed(() => {
  if (showAll.value) {
    return categories.value
  }
  return categories.value.slice(0, TAG_LIMIT)
})

// 定义事件
const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

// 获取分类列表
const fetchCategories = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    const response = await categoryApi.getCategories({
      pagesize: 100,
      pagenum: 1
    })
    
    console.log('Category API Response:', response)
    
    // 检查响应状态
    if (response.status !== 200) {
      error.value = '网络请求失败'
      return
    }
    
    const { data, status, message } = response.data
    
    // 检查API返回状态 - 根据后端API响应，0表示成功
    if (status !== 0) {
      error.value = message || '获取分类列表失败'
      console.error('获取分类列表失败:', message)
      return
    }
    
    // 确保data是数组
    if (!Array.isArray(data)) {
      error.value = '数据格式错误'
      console.error('数据格式错误，期望数组但得到:', typeof data)
      return
    }
    
    // 设置分类数据 - 修复数据结构问题，兼容不同字段名
    categories.value = data.map((item: any) => ({
      id: item.ID !== undefined ? item.ID : item.id,
      name: item.name,
      articleCount: item.article_count !== undefined ? item.article_count : (item.ArticleCount || 0)
    }))
  } catch (err: any) {
    error.value = err.message || '获取分类列表失败'
  } finally {
    loading.value = false
    emit('loading', false)
  }
}

// 计算标签字体大小
const calculateFontSize = (count: number) => {
  if (count <= 0) return '12px'
  if (count <= 5) return '14px'
  if (count <= 10) return '16px'
  if (count <= 20) return '18px'
  return '20px'
}

// 切换显示全部/部分标签
const toggleShowAll = () => {
  showAll.value = !showAll.value
}

// 重试函数
const onRetry = () => {
  fetchCategories()
}

// 暴露方法给父组件
defineExpose({
  fetchCategories
})

// 组件挂载时获取数据
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  background: #f8f9fa;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding: 15px;
  margin: 0 -10px; /* 补偿.card-content的padding */
  width: calc(100% + 20px);
}

.tag {
  display: inline-block;
  padding: 8px 15px;
  background: #e9ecef;
  border-radius: 20px;
  color: #495057;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  font-weight: 500;
}

.tag:hover {
  background: #007bff;
  color: white;
  transform: scale(1.05);
  border-radius: 22px;
}

.empty-state {
  text-align: center;
  padding: 30px 10px;
  color: #888;
}

.skeleton-loader {
  animation: skeleton-loading 1s linear infinite alternate;
  padding: 15px;
}

@keyframes skeleton-loading {
  0% {
    background-color: hsl(200, 20%, 80%);
  }
  100% {
    background-color: hsl(200, 20%, 95%);
  }
}

.skeleton-header {
  height: 20px;
  width: 60%;
  margin-bottom: 15px;
  border-radius: 4px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.skeleton-tag {
  height: 24px;
  width: 60px;
  border-radius: 12px;
}

.error-message {
  text-align: center;
  padding: 30px 10px;
  color: #dc3545;
}

.retry-button {
  margin-top: 15px;
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.retry-button:hover {
  background-color: #0056b3;
  border-radius: 8px;
}

.see-more-container {
  display: flex;
  justify-content: center;
  width: 100%;
}

.see-more-button {
  width: 40%;
  padding: 10px;
  margin-top: 15px;
  background-color: #f8f9fa;
  color: #007bff;
  border: 1px solid #dee2e6;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.see-more-button:hover {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
  border-radius: 20px;
}
</style>
