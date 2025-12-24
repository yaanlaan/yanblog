<template>
  <div class="category-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>分类管理</span>
          <el-button type="primary" @click="handleAdd">新增分类</el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <CategorySearchForm
        :model-value="searchForm"
        @update:modelValue="handleSearchFormUpdate"
        @search="handleSearch"
        @reset="handleReset"
      />
      
      <!-- 分类表格 -->
      <el-table 
        :data="categoryList" 
        border 
        style="width: 100%" 
        v-loading="loading"
        :empty-text="error ? '数据加载失败，请检查网络连接' : '暂无数据'"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" />
        <el-table-column label="封面图片" width="120">
          <template #default="scope">
            <el-image
              :src="scope.row.img || defaultImage"
              class="category-cover"
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
        <el-table-column prop="top" label="置顶排序" width="100" />
        <el-table-column prop="createdAt" label="创建时间" />
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <CategoryActions
              :category="scope.row"
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
    
    <!-- 分类编辑表单 -->
    <CategoryForm
      v-model="dialogVisible"
      :title="dialogTitle"
      :category="categoryForm"
      @submit="submitCategoryForm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import { Picture } from '@element-plus/icons-vue'
import { categoryApi } from '@/services/api'
import CategorySearchForm from '@/components/category/CategorySearchForm.vue'
import CategoryActions from '@/components/category/CategoryActions.vue'
import CategoryForm from '@/components/category/CategoryForm.vue'

// 路由
const router = useRouter()
const route = useRoute()

// 默认图片
const defaultImage = new URL('../../assets/img/无封面.jpg', import.meta.url).href

// 分类数据类型
interface Category {
  id: number
  name: string
  img: string
  top: number
  createdAt: string
}

// 所有分类数据（用于前端分页）
const allCategories = ref<Category[]>([])

// 当前页分类数据
const categoryList = ref<Category[]>([])

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 5,
  total: 0
})

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAdd = ref(true)

// 加载状态
const loading = ref(false)
const error = ref(false)

// 分类表单
const categoryForm = reactive({
  id: 0,
  name: '',
  img: '',
  top: 0
})

// 计算当前页数据
const updateCurrentPageData = () => {
  // 应用搜索过滤
  let resultCategories = [...filteredCategories.value]
  
  // 更新总数
  pagination.total = resultCategories.length
  
  // 计算当前页数据
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  categoryList.value = resultCategories.slice(start, end)
  
  // 更新URL参数
  updateUrlParams()
}

// filteredCategories 应该是一个 ref
const filteredCategories = ref<Category[]>([])

// 获取分类列表（获取所有数据）
const getCategoryList = async () => {
  loading.value = true
  error.value = false
  try {
    const response = await categoryApi.getCategories({
      pagesize: -1,  // 获取所有数据
      pagenum: -1    // 获取所有数据
    })
    
    // 解析后端返回的数据
    const { data } = response.data
    
    // 转换数据格式
    allCategories.value = data.map((item: any) => ({
      id: item.ID || item.id,
      name: item.name,
      img: item.img,
      top: item.top,
      createdAt: item.CreatedAt || item.created_at || ''
    }))
    
    // 应用筛选条件
    applyFilters()
  } catch (err) {
    error.value = true
    ElMessage.error('获取分类列表失败')
    console.error('获取分类列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 应用筛选条件
const applyFilters = (resetPage: boolean = true) => {
  // 应用搜索过滤
  let filtered = allCategories.value
  if (searchForm.name) {
    filtered = filtered.filter(category => 
      category.name.toLowerCase().includes(searchForm.name.toLowerCase())
    )
  }
  
  filteredCategories.value = filtered
  
  // 仅在需要时重置到第一页（如搜索时）
  if (resetPage) {
    pagination.currentPage = 1
  }
  
  // 更新分页信息和当前页数据
  updateCurrentPageData()
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
const handleSearchFormUpdate = (value: { name: string }) => {
  Object.assign(searchForm, value)
}

// 处理搜索
const handleSearch = () => {
  applyFilters(true) // 搜索时重置页码
}

// 处理重置
const handleReset = () => {
  // 重置搜索表单
  searchForm.name = ''
  
  // 恢复所有数据
  filteredCategories.value = [...allCategories.value]
  
  // 重置到第一页
  pagination.currentPage = 1
  
  // 更新分页信息和当前页数据
  updateCurrentPageData()
}

// 处理新增
const handleAdd = () => {
  dialogTitle.value = '新增分类'
  isAdd.value = true
  dialogVisible.value = true
  // 重置表单
  Object.assign(categoryForm, {
    id: 0,
    name: '',
    img: '',
    top: 0
  })
}

// 处理编辑
const handleEdit = (row: Category) => {
  dialogTitle.value = '编辑分类'
  isAdd.value = false
  dialogVisible.value = true
  // 填充表单数据
  Object.assign(categoryForm, {
    id: row.id,
    name: row.name,
    img: row.img || '',  // 确保img字段为空字符串而不是undefined
    top: row.top
  })
}

// 处理删除
const handleDelete = (row: Category) => {
  ElMessageBox.confirm(
    `确定要删除分类"${row.name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await categoryApi.deleteCategory(row.id)
      ElMessage.success('分类删除成功')
      // 重新获取分类列表
      await getCategoryList()
    } catch (error) {
      ElMessage.error('分类删除失败')
      console.error('分类删除失败:', error)
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 提交分类表单
const submitCategoryForm = async (formData: { id: number; name: string; img: string; top: number }) => {
  try {
    console.log('submitCategoryForm called with:', formData);
    let res;
    if (isAdd.value) {
      // 新增分类
      console.log('Creating new category');
      res = await categoryApi.createCategory({
        name: formData.name,
        img: formData.img, // 不再设置默认值，保持为空或用户输入的值
        top: formData.top
      })
    } else {
      // 编辑分类
      console.log('Updating category with id:', formData.id);
      console.log('Data being sent:', {
        name: formData.name,
        img: formData.img,
        top: formData.top
      });
      
      // 添加更详细的调试信息
      const requestData = {
        name: formData.name,
        img: formData.img,
        top: formData.top
      };
      
      console.log('Sending PUT request to:', `/api/v1/category/${formData.id}`);
      console.log('Request data:', requestData);
      
      res = await categoryApi.updateCategory(formData.id, requestData);
    }

    // 检查后端返回的状态码
    if (res.data.status !== 200) {
      ElMessage.error(res.data.message || (isAdd.value ? '分类创建失败' : '分类更新失败'))
      return
    }
    
    ElMessage.success(isAdd.value ? '分类创建成功' : '分类更新成功')
    
    // 关闭对话框
    dialogVisible.value = false
    
    // 重新获取分类列表
    console.log('Refreshing category list');
    await getCategoryList()
  } catch (error: any) {
    console.error('提交分类表单失败:', error)
    console.error('Error details:', {
      message: error.message,
      response: error.response,
      request: error.request
    });
    ElMessage.error(isAdd.value ? '分类创建失败' : '分类更新失败')
  }
}

// 更新URL参数
const updateUrlParams = () => {
  const query: Record<string, string | undefined> = {}
  
  if (searchForm.name) {
    query.name = searchForm.name
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
  const paramKeys = ['name', 'page', 'pageSize']
  for (const key of paramKeys) {
    if (query[key] !== currentQuery[key]) {
      needUpdate = true
      break
    }
  }
  
  // 检查是否有额外的参数需要移除
  for (const key in currentQuery) {
    if (!['name', 'page', 'pageSize'].includes(key) && query[key] === undefined) {
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
  searchForm.name = (query.name as string) || ''
  
  // 初始化分页
  pagination.currentPage = parseInt(query.page as string) || 1
  pagination.pageSize = parseInt(query.pageSize as string) || 5
}

// 监听路由参数变化
watch(
  () => route.query,
  () => {
    initFromUrlParams()
    // 应用筛选条件，但不重置页码（因为页码已经在initFromUrlParams中设置）
    applyFilters(false)
  }
)

// 组件挂载时获取数据
onMounted(() => {
  initFromUrlParams()
  getCategoryList()
})
</script>

<style scoped>
.category-list {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.category-cover {
  width: 100%;
  height: 100%;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 30px;
}
</style>
