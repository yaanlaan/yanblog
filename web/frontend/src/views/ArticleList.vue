<template>
  <div class="article-list-page">
    <div class="container">
      <div class="content-wrapper">
        <div class="main-content">
          <div class="page-header">
            <h1>{{ pageTitle }}</h1>
          </div>
          
          <div class="filters">
            <!-- 分类筛选 -->
            <div class="filter-group">
              <label>分类:</label>
              <select v-model="selectedCategory" @change="handleCategoryChange">
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
          
          <!-- 文章列表 -->
          <div class="articles" v-loading="loading">
            <div 
              v-for="article in articles" 
              :key="article.id" 
              class="article-item"
            >
              <div class="article-header">
                <h2 class="article-title">
                  <router-link :to="`/article/${article.id}`">
                    {{ article.title }}
                  </router-link>
                </h2>
                <div class="article-meta">
                  <span class="category">
                    分类: {{ article.categoryName }}
                  </span>
                  <span class="date">
                    发布时间: {{ formatDate(article.createdAt) }}
                  </span>
                </div>
              </div>
              
              <div class="article-summary">
                <p>{{ article.desc || '暂无简介' }}</p>
              </div>
              
              <div class="article-footer">
                <router-link :to="`/article/${article.id}`" class="read-more">
                  阅读全文 »
                </router-link>
              </div>
            </div>
            
            <!-- 分页 -->
            <div class="pagination" v-if="total > 0">
              <el-pagination
                v-model:current-page="pagination.currentPage"
                v-model:page-size="pagination.pageSize"
                :page-sizes="[5, 10, 20, 50]"
                :total="total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
            
            <!-- 空状态 -->
            <div class="empty-state" v-if="!loading && articles.length === 0">
              <p>暂无文章</p>
            </div>
          </div>
        </div>
        
        <!-- 右侧边栏 -->
        <Sidebar />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElPagination } from 'element-plus'
import { articleApi, categoryApi } from '@/services/api'
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
const router = useRouter()

// 响应式数据
const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const total = ref(0)
const selectedCategory = ref('')

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

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 获取文章列表
const getArticles = async () => {
  loading.value = true
  try {
    let response
    if (selectedCategory.value) {
      // 获取指定分类下的文章
      response = await articleApi.getCategoryArticles(Number(selectedCategory.value), {
        pagesize: pagination.pageSize,
        pagenum: pagination.currentPage
      })
    } else if (activeCategoryId.value) {
      // 获取路由参数指定的分类下的文章
      response = await articleApi.getCategoryArticles(activeCategoryId.value, {
        pagesize: pagination.pageSize,
        pagenum: pagination.currentPage
      })
    } else {
      // 获取所有文章
      response = await articleApi.getArticles({
        pagesize: pagination.pageSize,
        pagenum: pagination.currentPage
      })
    }
    
    const { data, total: totalCount } = response.data
    articles.value = data.map((item: any) => ({
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
    total.value = totalCount
    
    // 更新选中的分类
    if (activeCategoryId.value) {
      selectedCategory.value = String(activeCategoryId.value)
    }
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

// 处理分类变化
const handleCategoryChange = () => {
  pagination.currentPage = 1
  getArticles()
}

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
  getArticles()
}

// 处理当前页变化
const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
  getArticles()
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
  getArticles()
}, { deep: true })
</script>

<style scoped>
.article-list-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.content-wrapper {
  display: flex;
  gap: 30px;
  min-height: calc(100vh - 280px);
}

.main-content {
  flex: 1;
  min-width: 0;
}

.page-header h1 {
  font-size: 24px;
  color: #333;
  margin-bottom: 20px;
}

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

.article-item {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 25px;
  margin-bottom: 25px;
  transition: box-shadow 0.3s;
}

.article-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.article-header {
  margin-bottom: 15px;
}

.article-title {
  font-size: 22px;
  margin: 0 0 10px 0;
}

.article-title a {
  text-decoration: none;
  color: #333;
  transition: color 0.3s;
}

.article-title a:hover {
  color: #007bff;
}

.article-meta {
  display: flex;
  gap: 15px;
  font-size: 14px;
  color: #888;
}

.article-summary p {
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.article-footer {
  margin-top: 20px;
}

.read-more {
  color: #007bff;
  text-decoration: none;
  font-weight: 500;
}

.read-more:hover {
  text-decoration: underline;
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #888;
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
}

@media (max-width: 768px) {
  .filter-group {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>