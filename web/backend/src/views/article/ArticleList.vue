<template>
  <div class="article-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文章管理</span>
          <div>
            <el-button type="danger" plain @click="handleBatchDelete" :disabled="selectedRows.length === 0">
              批量删除 ({{ selectedRows.length }})
            </el-button>
            <el-button type="success" @click="zipDialogVisible = true">ZIP 发布</el-button>
            <el-button type="primary" @click="handleAdd">新增文章</el-button>
            <el-button type="info" plain @click="showHelp = true">Markdown 帮助</el-button>
          </div>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <ArticleSearchForm
        :model-value="searchForm"
        :categories="categories"
        @update:modelValue="handleSearchFormUpdate"
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
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="类型" width="80">
          <template #default="scope">
            <el-tag v-if="scope.row.type === 2" type="warning" size="small">PDF</el-tag>
            <el-tag v-else type="success" size="small">文章</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="封面" width="120">
          <template #default="scope">
            <el-image
              :src="scope.row.img || defaultImage"
              class="article-cover"
              fit="cover"
              :preview-src-list="[scope.row.img || defaultImage]"
              preview-teleported
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" show-overflow-tooltip />
        <el-table-column prop="desc" label="简介" show-overflow-tooltip />
        <el-table-column prop="categoryName" label="分类" width="120" />
        <el-table-column prop="tags" label="标签" show-overflow-tooltip>
          <template #default="scope">
            <el-tag 
              v-for="(tag, index) in (scope.row.tags ? scope.row.tags.split(',') : [])" 
              :key="index"
              size="small"
              style="margin-right: 5px; margin-bottom: 5px;"
            >
              {{ tag }}
            </el-tag>
          </template>
        </el-table-column>
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

      <!-- Markdown 帮助对话框 -->
      <el-dialog v-model="showHelp" title="Markdown 语法 & ZIP 发布说明" width="700px">
        <div class="markdown-help">
          <h3>YAML Front Matter 字段（ZIP 上传）</h3>
          <table>
            <tr><th>字段</th><th>类型</th><th>必填</th><th>说明</th></tr>
            <tr><td><code>title</code></td><td>字符串</td><td>是</td><td>文章标题。为空时使用文件名</td></tr>
            <tr><td><code>date</code></td><td>日期</td><td>否</td><td>格式 YYYY-MM-DD</td></tr>
            <tr><td><code>tags</code></td><td>数组</td><td>否</td><td>如 [Go, Vue]</td></tr>
            <tr><td><code>category</code></td><td>字符串</td><td>否</td><td>分类名。不存在自动创建</td></tr>
            <tr><td><code>desc</code></td><td>字符串</td><td>否</td><td>文章摘要</td></tr>
            <tr><td><code>cover</code></td><td>字符串</td><td>否</td><td>封面图路径（相对于 zip）</td></tr>
          </table>
          <h4>示例</h4>
          <pre><code>---
title: "文章标题"
date: 2024-06-08
tags: [Go, Vue]
category: "技术"
desc: "摘要"
cover: "images/cover.jpg"
---</code></pre>
          <h3>支持的 Markdown 功能</h3>
          <p>代码块（语法高亮+行号+复制）、KaTeX 数学公式 <code>$x^2$</code> <code>$$\int$$</code>、Mermaid 流程图、表格、链接卡片（单独一行 URL 自动渲染为卡片）</p>
        </div>
      </el-dialog>

      <!-- ZIP上传对话框 -->
      <el-dialog v-model="zipDialogVisible" title="批量发布文章 (ZIP)" width="550px" @closed="zipFileList = []; uploadResults = []">
        <el-upload
          class="upload-demo"
          drag
          action="#"
          multiple
          :auto-upload="false"
          :file-list="zipFileList"
          :before-upload="beforeZipUpload"
          :on-remove="handleZipRemove"
          :on-change="handleZipChange"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            拖拽文件到这里 或 <em>点击选择</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持同时选择多个 .zip 文件，每个文件为包含 .md 和 images/ 的压缩包
            </div>
          </template>
        </el-upload>
        <div v-if="uploadResults.length > 0" style="margin-top:16px; max-height:200px; overflow-y:auto;">
          <div v-for="r in uploadResults" :key="r.file_name" style="padding:4px 0; border-bottom:1px solid #eee;">
            <span :style="{color: r.status === 200 ? '#67c23a' : '#f56c6c'}">
              {{ r.status === 200 ? '✓' : '✗' }}
            </span>
            <span style="margin-left:8px;">{{ r.file_name }}</span>
            <span v-if="r.title" style="margin-left:8px; color:#909399;">→ {{ r.title }}</span>
            <span v-if="r.status !== 200" style="margin-left:8px; color:#f56c6c; font-size:12px;">{{ r.message }}</span>
          </div>
        </div>
        <template #footer>
          <el-button @click="zipDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleBatchUpload" :loading="uploading" :disabled="zipFileList.length === 0">
            上传全部 ({{ zipFileList.length }})
          </el-button>
        </template>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, onActivated } from 'vue'

defineOptions({
  name: 'ArticleList'
})

import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import { Picture, UploadFilled } from '@element-plus/icons-vue'
import { articleApi, categoryApi } from '@/services/api'
import ArticleSearchForm from '@/components/article/ArticleSearchForm.vue'
import ArticleActions from '@/components/article/ArticleActions.vue'

// 路由
const router = useRouter()
const route = useRoute()

// 默认图片
const defaultImage = new URL('../../assets/img/无封面.jpg', import.meta.url).href

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
  tags: string
  type?: number
  pdf_url?: string
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
  
  // 更新URL参数
  updateUrlParams()
}

// 更新URL参数
const updateUrlParams = () => {
  const query: Record<string, string | undefined> = {}
  
  if (searchForm.title) {
    query.title = searchForm.title
  }
  
  if (searchForm.categoryId !== undefined) {
    query.categoryId = searchForm.categoryId.toString()
  }
  
  if (pagination.currentPage > 1) {
    query.page = pagination.currentPage.toString()
  }
  
  if (pagination.pageSize !== 5) {
    query.pageSize = pagination.pageSize.toString()
  }
  
  // 只有当查询参数发生变化时才更新路由
  const currentQuery = route.query
  let needUpdate = false
  
  // 检查参数是否发生变化
  const paramKeys = ['title', 'categoryId', 'page', 'pageSize']
  for (const key of paramKeys) {
    if (query[key] !== currentQuery[key]) {
      needUpdate = true
      break
    }
  }
  
  // 检查是否有额外的参数需要移除
  for (const key in currentQuery) {
    if (!['title', 'categoryId', 'page', 'pageSize'].includes(key) && query[key] === undefined) {
      needUpdate = true
      break
    }
  }
  
  if (needUpdate) {
    router.replace({ query })
  }
}

// 从URL参数初始化搜索表单和分页
const initFromUrlParams = () => {
  const query = route.query
  
  // 初始化搜索表单
  searchForm.title = (query.title as string) || ''
  searchForm.categoryId = query.categoryId ? Number(query.categoryId) : undefined
  
  // 初始化分页
  pagination.currentPage = query.page ? Number(query.page) : 1
  pagination.pageSize = query.pageSize ? Number(query.pageSize) : 5
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
      tags: item.tags || '',
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

// 处理搜索表单更新
const handleSearchFormUpdate = (value: {title: string, categoryId: number | undefined}) => {
  Object.assign(searchForm, value)
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

// 批量选择
const selectedRows = ref<Article[]>([])
const handleSelectionChange = (rows: Article[]) => {
  selectedRows.value = rows
}

// 批量删除
const handleBatchDelete = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先选择要删除的文章')
    return
  }
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedRows.value.length} 篇文章吗？此操作不可恢复！`,
      '批量删除确认',
      { confirmButtonText: '确定删除', cancelButtonText: '取消', type: 'error' }
    )
    const ids = selectedRows.value.map(row => row.id)
    const res = await articleApi.batchDeleteArticles(ids)
    if (res.data.status === 200) {
      ElMessage.success(`成功删除 ${res.data.deleted} 篇`)
      selectedRows.value = []
      getArticleList()
    } else {
      ElMessage.error(res.data.message || '删除失败')
    }
  } catch (err) {
    // 取消操作
  }
}

// ZIP上传对话框
const zipDialogVisible = ref(false)
const zipFileList = ref<any[]>([])
const uploadResults = ref<any[]>([])
const uploading = ref(false)

// Markdown 帮助对话框
const showHelp = ref(false)

// 文件校验
const beforeZipUpload = (file: File) => {
  const isZip = file.type === 'application/zip' || file.name.endsWith('.zip')
  const isLt50M = file.size / 1024 / 1024 < 50
  if (!isZip) { ElMessage.error('只能上传 ZIP 文件!'); return false }
  if (!isLt50M) { ElMessage.error('文件大小不能超过 50MB!'); return false }
  return true
}

const handleZipChange = () => {
  uploadResults.value = []
}

const handleZipRemove = (_file: any, fileList: any[]) => {
  zipFileList.value = fileList
}

// 批量上传
const handleBatchUpload = async () => {
  if (zipFileList.value.length === 0) {
    ElMessage.warning('请先选择 ZIP 文件')
    return
  }
  uploading.value = true
  uploadResults.value = []
  try {
    const formData = new FormData()
    zipFileList.value.forEach((f: any) => {
      formData.append('files', f.raw || f)
    })
    const res = await articleApi.uploadZipBatch(formData)
    if (res.data.results) {
      uploadResults.value = res.data.results
      const successCount = res.data.success || 0
      if (successCount > 0) {
        ElMessage.success(`成功上传 ${successCount}/${res.data.total} 篇文章`)
        getArticleList()
      }
      if (successCount < res.data.total) {
        ElMessage.warning(`${res.data.total - successCount} 篇失败，详见下方列表`)
      }
    }
  } catch (err) {
    ElMessage.error('上传出错')
  } finally {
    uploading.value = false
  }
}

// 监听路由参数变化
watch(
  () => route.query,
  () => {
    initFromUrlParams()
    updateCurrentPageData()
  }
)

// 组件挂载时获取数据
onMounted(() => {
  initFromUrlParams()
  getArticleList()
  getCategoryList()
})

// 组件激活时更新数据
onActivated(() => {
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
.markdown-help h3 { margin: 16px 0 8px; border-bottom: 1px solid #ebeef5; padding-bottom: 4px; }
.markdown-help h4 { margin: 12px 0 6px; font-size: 14px; }
.markdown-help table { width: 100%; border-collapse: collapse; margin: 8px 0 16px; }
.markdown-help th, .markdown-help td { border: 1px solid #ebeef5; padding: 6px 10px; text-align: left; font-size: 13px; }
.markdown-help th { background: #f5f7fa; }
.markdown-help code { background: #f5f7fa; padding: 1px 5px; border-radius: 3px; font-size: 13px; color: #e96900; }
.markdown-help pre { background: #f5f7fa; padding: 10px; border-radius: 4px; overflow-x: auto; font-size: 13px; margin: 8px 0 16px; }
.markdown-help p { font-size: 13px; color: #606266; margin: 4px 0; }
</style>