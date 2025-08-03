<template>
  <div class="article-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文章管理</span>
          <el-button type="primary" @click="handleAdd">新增文章</el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <ArticleSearchForm
        v-model="searchForm"
        :categories="categories"
        @search="handleSearch"
        @reset="handleReset"
      />
      
      <!-- 文章表格 -->
      <el-table 
        :data="articleList" 
        border 
        style="width: 100%" 
        v-loading="loading"
        :empty-text="error ? '数据加载失败，请检查网络连接' : '暂无数据'"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="封面" width="120">
          <template #default="scope">
            <el-image
              v-if="scope.row.img"
              :src="scope.row.img"
              class="article-cover"
              fit="cover"
              :preview-src-list="[scope.row.img]"
              preview-teleported
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div v-else class="no-cover">
              <el-icon><Picture /></el-icon>
              <span>无封面</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" show-overflow-tooltip />
        <el-table-column prop="desc" label="简介" show-overflow-tooltip />
        <el-table-column prop="categoryName" label="分类" width="120" />
        <el-table-column prop="top" label="置顶等级" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.top > 0" type="danger">等级{{ scope.row.top }}</el-tag>
            <el-tag v-else type="info">未置顶</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170">
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="170">
          <template #default="scope">
            {{ formatDateTime(scope.row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <ArticleActions
              :article="scope.row"
              @edit="handleEdit"
              @delete="handleDelete"
            />
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.currentPage"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[5, 10, 20, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { Picture } from '@element-plus/icons-vue'
import { articleApi, categoryApi } from '@/services/api'
import ArticleSearchForm from '@/components/article/ArticleSearchForm.vue'
import ArticleActions from '@/components/article/ArticleActions.vue'

// 路由
const router = useRouter()

// 文章数据类型
interface Article {
  id: number
  title: string
  categoryId: number | undefined
  categoryName: string
  desc: string
  content: string
  img: string
  top: number
  createdAt: string
  updatedAt: string
}

// 分类数据类型
interface Category {
  id: number
  name: string
}

// 所有文章数据（用于前端分页）
const allArticles = ref<Article[]>([])

// 搜索表单
const searchForm = reactive({
  title: '',
  categoryId: undefined as number | undefined
})

// 更新当前页数据（前端分页和筛选）
const updateCurrentPageData = () => {
  // 应用搜索和筛选
  let filteredArticles = [...allArticles.value]
  
  // 标题搜索
  if (searchForm.title) {
    filteredArticles = filteredArticles.filter(article => 
      article.title.toLowerCase().includes(searchForm.title.toLowerCase())
    )
  }
  
  // 分类筛选
  if (searchForm.categoryId !== undefined && searchForm.categoryId !== null) {
    const categoryId = Number(searchForm.categoryId)
    filteredArticles = filteredArticles.filter(article => 
      article.categoryId === categoryId
    )
  }
  
  // 更新总数
  pagination.total = filteredArticles.length
  
  // 计算当前页数据
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  articleList.value = filteredArticles.slice(start, end)
}

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 5,
  total: 0
})

// 文章列表（当前页数据）
const articleList = ref<Article[]>([])

// 分类列表
const categories = ref<Category[]>([])

// 加载状态
const loading = ref(false)

// 错误状态
const error = ref(false)

// 格式化日期时间
const formatDateTime = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).replace(/\//g, '-')
}

// 获取文章列表
const getArticleList = async () => {
  loading.value = true
  error.value = false
  try {
    const response = await articleApi.getArticles({
      pagesize: -1,  // 获取所有数据
      pagenum: -1    // 获取所有数据
    })
    
    // 解析后端返回的数据
    const { data } = response.data
    allArticles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid !== undefined ? parseInt(item.cid, 10) : undefined,
      categoryName: (item.Category?.name || item.category?.name || '未分类') as string,
      desc: item.desc,
      content: item.content,
      img: item.img,
      top: item.top || 0,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
    
    // 更新当前页数据
    updateCurrentPageData()
  } catch (err) {
    error.value = true
    ElMessage.error('获取文章列表失败')
    console.error('获取文章列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  pagination.currentPage = 1
  updateCurrentPageData()
}

// 处理当前页变化
const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
  updateCurrentPageData()
}

// 处理搜索
const handleSearch = () => {
  pagination.currentPage = 1
  updateCurrentPageData()
}

// 处理重置
const handleReset = () => {
  searchForm.title = ''
  searchForm.categoryId = undefined
  pagination.currentPage = 1
  updateCurrentPageData()
}

// 获取分类列表
const getCategoryList = async () => {
  try {
    const response = await categoryApi.getCategories({
      pagesize: -1,
      pagenum: -1
    })
    
    // 解析后端返回的数据
    const { data } = response.data
    console.log('从后端获取的分类数据:', data)
    
    categories.value = data.map((item: any) => ({
      id: item.ID !== undefined ? parseInt(item.ID, 10) : parseInt(item.id, 10),
      name: item.name
    }))
    
    console.log('处理后的分类数据:', categories.value)
  } catch (error) {
    ElMessage.error('获取分类列表失败')
    console.error(error)
  }
}

// 处理新增文章
const handleAdd = () => {
  router.push('/article/add')
}

// 处理编辑文章
const handleEdit = (article: Article) => {
  router.push(`/article/edit/${article.id}`)
}

// 处理删除文章
const handleDelete = (article: Article) => {
  ElMessageBox.confirm(
    `确定要删除文章《${article.title}》吗？此操作不可恢复！`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await articleApi.deleteArticle(article.id)
      ElMessage.success('文章删除成功')
      // 重新加载数据
      await getArticleList()
    } catch (error) {
      ElMessage.error('文章删除失败')
      console.error(error)
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 组件挂载时获取数据
onMounted(() => {
  getArticleList()
  getCategoryList()
})
</script>

<style scoped>
.article-list {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.article-cover {
  width: 100%;
  height: 60px;
  border-radius: 4px;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #999;
}

.no-cover {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 60px;
  color: #999;
  font-size: 12px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>