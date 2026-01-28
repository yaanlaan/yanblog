<template>
  <div class="category-view">
    <div class="container">
      <div class="page-header">
        <h1>文章分类</h1>
        <p>浏览所有文章分类</p>
      </div>
      
      <div class="category-grid" v-loading="loading">
        <div 
          v-for="category in displayedCategories" 
          :key="category.id" 
          class="category-card"
          @click="goToCategory(category.id)"
        >
          <div class="category-top-tag" v-if="category.top > 0">
            置顶
          </div>
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
        
        <div v-if="!loading && displayedCategories.length === 0" class="no-data">
          <p>暂无分类数据</p>
        </div>
      </div>
      
      <!-- 滚动加载触发器 -->
      <div ref="loadingTrigger" class="loading-trigger">
        <div v-if="loadingCategories" class="loading-more">
           <!-- loading state is usually for initial load, but here we can reuse or use a separate loadingMore flag if async -->
           <!-- Since we fetch ALL at start, 'loading' is only for initial fetch. 
                Scrolling is instant, so we might not need a spinner unless we simulate delay. -->
           <span v-if="displayedCategories.length < total">- 向下滑动加载更多 -</span>
           <span v-else>- 也就是这些了 -</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { categoryApi } from '@/services/api'
// import { ElPagination } from 'element-plus'

// 定义分类类型
interface Category {
  id: number
  name: string
  img: string
  article_count: number
  top: number
  created_at?: string
}

// 获取路由实例
const route = useRoute()
const router = useRouter()

// 响应式数据
const allCategories = ref<Category[]>([]) // 存储所有分类数据
const displayedCategories = ref<Category[]>([]) // 存储当前显示的分类数据
const loading = ref(true)
const loadingCategories = ref(false) // Just for uniformity if needed
const total = ref(0)
const displayCount = ref(12)

// 滚动监听
const loadingTrigger = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

// 默认图片
const defaultImage = new URL('@/assets/img/无封面.jpg', import.meta.url).href

// 处理图片加载错误
const handleImageError = (e: Event) => {
  const target = e.target as HTMLImageElement
  target.src = defaultImage
}

// 获取所有分类列表
const getCategories = async () => {
  try {
    loading.value = true
    const response = await categoryApi.getCategories({
      pagesize: -1,
      pagenum: -1
    })
    
    // 解析数据并按置顶等级排序
    const { data, total: totalCount } = response.data
    total.value = totalCount
    
    allCategories.value = data.map((item: any) => ({
      id: item.ID || item.id,
      name: item.name,
      img: item.img,
      top: item.top || 0,
      article_count: item.article_count || 0
    })).sort((a: Category, b: Category) => {
      // 首先按置顶等级排序（0表示不置顶，其他值越小等级越高）
      if (a.top !== b.top) {
        // 如果其中一个为0（不置顶），则不置顶的排在后面
        if (a.top === 0) return 1;
        if (b.top === 0) return -1;
        // 都置顶的情况下，top值小的排在前面
        return a.top - b.top;
      }
      // 如果置顶等级相同，保持原有顺序
      return 0;
    })
    
    // 初始化显示数据
    updateDisplayedCategories()
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

// 更新显示的分类列表（无限滚动）
const updateDisplayedCategories = () => {
  // 应用切片
  displayedCategories.value = allCategories.value.slice(0, displayCount.value)
}

// 加载更多
const loadMore = () => {
    if (displayedCategories.value.length >= total.value) return
    displayCount.value += 8
    updateDisplayedCategories()
}

// 监听器设置
const setupObserver = () => {
    if (observer) observer.disconnect()
    
    observer = new IntersectionObserver((entries) => {
        const entry = entries[0]
        if (entry.isIntersecting && !loading.value && displayedCategories.value.length < total.value) {
            loadMore()
        }
    }, {
        rootMargin: '100px'
    })
    
    if (loadingTrigger.value) {
        observer.observe(loadingTrigger.value)
    }
}

// 组件挂载时获取数据
onMounted(async () => {
  await getCategories()
  nextTick(() => {
      setupObserver()
  })
})

onUnmounted(() => {
    if (observer) observer.disconnect()
})

// 监听数据显示变化，确保 observer 有效
watch(() => displayedCategories.value, () => {
    // DOM更新后可能位置变化，但 observer 通常不需要重置
})

// 在组件卸载前清理watcher（可选）
// onUnmounted(() => {
//   // 清理逻辑（如果需要）
// })
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
  color: var(--color-heading);
}

.page-header p {
  color: var(--color-text-secondary);
  font-size: 1rem;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.category-card {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
  background: var(--color-background-soft);
  box-shadow: 0 2px 8px var(--color-shadow);
  position: relative;
}

.category-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 16px var(--color-shadow);
}

.category-top-tag {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: #ff6b6b;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
  z-index: 1;
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
  color: var(--color-heading);
}

.category-info p {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.no-data {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  color: var(--color-text-secondary);
}

.loading-trigger {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 30px;
}

.loading-more {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--color-text-secondary);
  font-size: 14px;
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