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
          :view-mode="viewMode"
          @category-change="handleCategoryChange"
          @search="handleSearch"
          @reset="handleReset"
          @view-change="(mode) => viewMode = mode"
        />
        
        <ArticleListContent
          :articles="articles"
          :loading="loading"
          :total="total"
          :view-mode="viewMode"
          @load-more="loadMore"
        />
      </template>
      <template #sidebar>
        <Sidebar />
      </template>
    </MainLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
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
  top: number // 添加top字段
  tags: string
  views: number // 添加views字段
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
const viewMode = ref<'grid' | 'list'>('grid')

// 无限滚动显示数量
const displayCount = ref(10)

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
    let response;
    
    // 如果是分类文章页面，使用分类文章API
    if (activeCategoryId.value) {
      response = await articleApi.getCategoryArticles(activeCategoryId.value, {
        pagesize: -1,
        pagenum: -1
      });
    } else {
      // 否则使用普通文章列表API
      response = await articleApi.getArticles({
        pagesize: -1,
        pagenum: -1
      });
    }
    
    const { data, status } = response.data
    
    if (status !== 200) {
      console.error('获取文章列表失败:', response.data.message)
      return
    }

    allArticles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      top: item.top || 0, // 添加top字段
      tags: item.tags || '', // 添加tags字段
      views: item.views || 0, // 添加views字段
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
    // 按创建时间倒序排列（最新的在最前面）
    .sort((a: Article, b: Article) => {
      const dateA = new Date(a.createdAt).getTime()
      const dateB = new Date(b.createdAt).getTime()
      return dateB - dateA
    })
    
    // 初始化显示数据
    displayCount.value = 10
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
      article.desc.toLowerCase().includes(keyword) ||
      (article.tags && article.tags.toLowerCase().includes(keyword))
    )
  }
  
  // 更新总数
  total.value = filteredArticles.length
  
  // 应用切片 (无限滚动)
  articles.value = filteredArticles.slice(0, displayCount.value)
}

// 加载更多文章
const loadMore = () => {
    // 增加显示数量
    displayCount.value += 10
    updateDisplayedArticles()
}

// 处理分类变化
const handleCategoryChange = (value: string) => {
  selectedCategory.value = value
  displayCount.value = 10
  updateDisplayedArticles()
}

// 处理搜索
const handleSearch = (keyword: string) => {
  searchKeyword.value = keyword
  displayCount.value = 10
  updateDisplayedArticles()
}

// 处理重置
const handleReset = () => {
  searchKeyword.value = ''
  selectedCategory.value = ''
  displayCount.value = 10
  updateDisplayedArticles()
}

// 组件挂载时获取数据
onMounted(() => {
  // 初始化搜索关键词
  if (route.query.keyword) {
    searchKeyword.value = route.query.keyword as string
  }
  getCategories()
  getArticles()
})

// 监听路由参数变化
watch(() => route.params, (newParams, oldParams) => {
  // 当路由参数变化时，重新获取文章数据
  // displayCount 在 getArticles 内部会重置
  getArticles()
}, { deep: true })

// 监听路由查询参数变化（如果有的话）
watch(() => route.query, () => {
    displayCount.value = 10
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
  color: var(--color-heading);
  margin-bottom: 20px;
}
</style>