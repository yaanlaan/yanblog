<template>
  <div class="category-view">
    <div class="container">
      <div class="page-header">
        <h1>文章分类</h1>
        <p>浏览所有文章分类</p>
      </div>
      
      <div class="category-grid" v-loading="loading">
        <div 
          v-for="category in categories" 
          :key="category.id" 
          class="category-card"
          @click="goToCategory(category.id)"
        >
          <div class="category-image">
            <img 
              :src="category.img || defaultImage" 
              :alt="category.name"
              @error="handleImageError"
            />
          </div>
          <div class="category-info">
            <h3>{{ category.name }}</h3>
            <p>{{ category.article_count }} 篇文章</p>
          </div>
        </div>
        
        <div v-if="!loading && categories.length === 0" class="no-data">
          <p>暂无分类数据</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { categoryApi } from '@/services/api'

// 定义分类类型
interface Category {
  id: number
  name: string
  img: string
  article_count: number
  created_at?: string
}

// 获取路由实例
const router = useRouter()

// 响应式数据
const categories = ref<Category[]>([])
const loading = ref(true)

// 默认图片
const defaultImage = new URL('@/assets/img/无封面.jpg', import.meta.url).href

// 处理图片加载错误
const handleImageError = (e: Event) => {
  const target = e.target as HTMLImageElement
  target.src = defaultImage
}

// 获取分类列表
const getCategories = async () => {
  try {
    loading.value = true
    const response = await categoryApi.getCategories({
      pagesize: -1,
      pagenum: -1
    })
    
    // 解析数据
    const { data } = response.data
    categories.value = data.map((item: any) => ({
      id: item.ID || item.id,
      name: item.name,
      img: item.img,
      article_count: item.article_count || 0
    }))
  } catch (error) {
    console.error('获取分类列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 跳转到分类文章页面
const goToCategory = (id: number) => {
  router.push(`/category/${id}`)
}

// 组件挂载时获取数据
onMounted(() => {
  getCategories()
})
</script>

<style scoped>
.category-view {
  padding: 20px 0;
  min-height: calc(100vh - 120px);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.page-header {
  text-align: center;
  margin-bottom: 30px;
}

.page-header h1 {
  font-size: 2rem;
  margin-bottom: 10px;
  color: #333;
}

.page-header p {
  color: #666;
  font-size: 1rem;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
}

.category-card {
  border: 1px solid #eaeaea;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.category-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.category-image {
  height: 160px;
  overflow: hidden;
}

.category-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.category-card:hover .category-image img {
  transform: scale(1.05);
}

.category-info {
  padding: 15px;
  text-align: center;
}

.category-info h3 {
  margin: 0 0 8px 0;
  font-size: 1.2rem;
  color: #333;
}

.category-info p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.no-data {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  color: #666;
}

@media (max-width: 768px) {
  .category-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
  }
  
  .category-image {
    height: 140px;
  }
  
  .page-header h1 {
    font-size: 1.5rem;
  }
}

@media (max-width: 480px) {
  .category-grid {
    grid-template-columns: 1fr;
  }
  
  .category-image {
    height: 180px;
  }
}
</style>