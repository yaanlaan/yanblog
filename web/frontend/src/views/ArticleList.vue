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
          :loading="loading || loadingMore"
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
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { articleApi, categoryApi } from '@/services/api'
import { mapArticle, mapCategory, mapCategoryList } from '@/utils/dataMapper'
import { PAGINATION } from '@/utils/constants'
import type { Article, Category } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import ArticleFilter from '@/components/article/ArticleFilter.vue'
import ArticleListContent from '@/components/article/ArticleListContent.vue'
import Sidebar from '@/components/Sidebar.vue'

// 路由信息
const route = useRoute()

const PAGE_SIZE = PAGINATION.ARTICLE_LIST_PAGE_SIZE

// 响应式数据
const articles = ref<Article[]>([])
const categories = ref<Category[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const total = ref(0)
const selectedCategory = ref('')
const searchKeyword = ref('')
const viewMode = ref<'grid' | 'list'>('grid')
const currentPage = ref(1)
const hasMore = ref(true)

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

// 重置状态并加载第一页
const resetAndFetch = async () => {
  articles.value = []
  currentPage.value = 1
  hasMore.value = true
  await fetchArticles()
}

// 获取文章列表（后端分页 + 筛选）
const fetchArticles = async (isLoadMore = false) => {
  if (!isLoadMore) {
    loading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    let res: any

    // 判断当前筛选模式：搜索优先 > 分类筛选 > 全部文章
    if (searchKeyword.value) {
      // 使用后端搜索接口
      const cid = selectedCategory.value ? Number(selectedCategory.value) : (activeCategoryId.value || undefined)
      res = await articleApi.searchArticles({
        keyword: searchKeyword.value,
        cid: cid || undefined,
        pagesize: PAGE_SIZE,
        pagenum: currentPage.value
      })
    } else if (selectedCategory.value) {
      // 前端选择了某个分类（非路由参数的分类）
      res = await articleApi.getCategoryArticles(Number(selectedCategory.value), {
        pagesize: PAGE_SIZE,
        pagenum: currentPage.value
      })
    } else if (activeCategoryId.value) {
      // 路由参数中的分类
      res = await articleApi.getCategoryArticles(activeCategoryId.value, {
        pagesize: PAGE_SIZE,
        pagenum: currentPage.value
      })
    } else {
      // 全部文章
      res = await articleApi.getArticles({
        pagesize: PAGE_SIZE,
        pagenum: currentPage.value
      })
    }

    const { data, status } = res.data

    if (status !== 200) {
      console.error('获取文章列表失败:', res.data.message)
      return
    }

    const newArticles = data.map(mapArticle)

    if (isLoadMore) {
      articles.value.push(...newArticles)
    } else {
      articles.value = newArticles
    }

    total.value = res.data.total
    hasMore.value = articles.value.length < total.value
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

// 加载更多
const loadMore = async () => {
  if (loadingMore.value || !hasMore.value) return
  currentPage.value++
  await fetchArticles(true)
}

// 获取分类列表
const getCategories = async () => {
  try {
    const response = await categoryApi.getCategories({
      pagesize: -1,
      pagenum: -1
    })
    const { data } = response.data
    categories.value = mapCategoryList(data)
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 处理分类变化 → 重新从第1页加载
const handleCategoryChange = async (value: string) => {
  selectedCategory.value = value
  await resetAndFetch()
}

// 处理搜索 → 重新从第1页加载
const handleSearch = async (keyword: string) => {
  searchKeyword.value = keyword
  await resetAndFetch()
}

// 处理重置
const handleReset = async () => {
  searchKeyword.value = ''
  selectedCategory.value = ''
  await resetAndFetch()
}

onMounted(async () => {
  if (route.query.keyword) {
    searchKeyword.value = route.query.keyword as string
  }
  getCategories()
  await fetchArticles()
})

watch(() => route.params, () => {
  // 分类路由切换时重置筛选并重新加载
  selectedCategory.value = ''
  searchKeyword.value = ''
  resetAndFetch()
}, { deep: true })

watch(() => route.query, () => {
  if (route.query.keyword) {
    searchKeyword.value = route.query.keyword as string
    resetAndFetch()
  }
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