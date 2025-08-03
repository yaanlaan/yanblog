<template>
  <div class="article-list-page">
    <MainLayout>
      <template #main>
        <div class="page-header">
          <h1>{{ pageTitle }}</h1>
        </div>
        
        <ArticleFilter 
          :categories="categories"
          :selected-category="selectedCategory"
          :search-keyword="searchKeyword"
          @category-change="handleCategoryChange"
          @search="handleSearch"
          @reset="handleReset"
        />
        
        <ArticleListContent
          :articles="articles"
          :loading="loading"
          :total="total"
          :current-page="pagination.currentPage"
          :page-size="pagination.pageSize"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </template>
      <template #sidebar>
        <Sidebar />
      </template>
    </MainLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { articleApi, categoryApi } from '@/services/api'
import MainLayout from '@/components/layout/MainLayout.vue'
import ArticleFilter from '@/components/article/ArticleFilter.vue'
import ArticleListContent from '@/components/article/ArticleListContent.vue'
import Sidebar from '@/components/Sidebar.vue'

// 类型定义
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  createdAt: string
  updatedAt: string
}

interface Category {
  id: number
  name: string
}

// 路由信息
const route = useRoute()

// 响应式数据
const allArticles = ref<Article[]>([]) // 存储所有文章数据
const articles = ref<Article[]>([]) // 存储当前显示的文章数据
const categories = ref<Category[]>([])
const loading = ref(false)
const total = ref(0)
const selectedCategory = ref('')
const searchKeyword = ref('')

const pagination = reactive({
  currentPage: 1,
  pageSize: 10
})

// 计算属性
const activeCategoryId = computed(() => {
  return route.params.id ? Number(route.params.id) : null
})

const pageTitle = computed(() => {
  if (activeCategoryId.value) {
    const category = categories.value.find(c => c.id === activeCategoryId.value)
    return category ? `${category.name} - 文章列表` : '文章列表'
  }
  return '所有文章'
})

// 获取文章列表
const getArticles = async () => {
  loading.value = true
  try {
    // 一次性获取所有文章数据
    const response = await articleApi.getArticles({
      pagesize: -1,
      pagenum: -1
    })
    
    const { data } = response.data
    allArticles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
    
    // 初始化显示数据
    updateDisplayedArticles()
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const getCategories = async () => {
  try {
    const response = await categoryApi.getCategories({
      pagesize: -1,
      pagenum: -1
    })
    
    const { data } = response.data
    categories.value = data.map((item: any) => ({
      id: item.ID,
      name: item.name
    }))
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 更新显示的文章列表（前端筛选和分页）
const updateDisplayedArticles = () => {
  // 应用筛选条件
  let filteredArticles = [...allArticles.value]
  
  // 应用分类筛选
  const categoryId = selectedCategory.value ? Number(selectedCategory.value) : null
  if (categoryId) {
    filteredArticles = filteredArticles.filter(article => article.categoryId === categoryId)
  } else if (activeCategoryId.value) {
    filteredArticles = filteredArticles.filter(article => article.categoryId === activeCategoryId.value)
  }
  
  // 应用搜索关键词筛选
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filteredArticles = filteredArticles.filter(article => 
      article.title.toLowerCase().includes(keyword) || 
      article.desc.toLowerCase().includes(keyword)
    )
  }
  
  // 更新总数
  total.value = filteredArticles.length
  
  // 应用分页
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  articles.value = filteredArticles.slice(start, end)
}

// 处理分类变化
const handleCategoryChange = (value: string) => {
  selectedCategory.value = value
  pagination.currentPage = 1
  updateDisplayedArticles()
}

// 处理搜索
const handleSearch = (keyword: string) => {
  searchKeyword.value = keyword
  pagination.currentPage = 1
  updateDisplayedArticles()
}

// 处理重置
const handleReset = () => {
  searchKeyword.value = ''
  selectedCategory.value = ''
  pagination.currentPage = 1
  updateDisplayedArticles()
}

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
  updateDisplayedArticles()
}

// 处理当前页变化
const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
  updateDisplayedArticles()
}

// 组件挂载时获取数据
onMounted(() => {
  getCategories()
  getArticles()
})

// 监听路由参数变化
import { watch } from 'vue'
watch(() => route.params, () => {
  pagination.currentPage = 1
  updateDisplayedArticles()
}, { deep: true })
</script>

<style scoped>
.article-list-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.page-header h1 {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
}
</style>