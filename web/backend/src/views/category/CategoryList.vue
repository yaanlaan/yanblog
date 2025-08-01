，<template>
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
        v-model="searchForm"
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { categoryApi } from '@/services/api'
import CategorySearchForm from '@/components/category/CategorySearchForm.vue'
import CategoryActions from '@/components/category/CategoryActions.vue'
import CategoryForm from '@/components/category/CategoryForm.vue'

// 分类数据类型
interface Category {
  id: number
  name: string
  createdAt: string
}

// 所有分类数据（用于前端分页）
const allCategories = ref<Category[]>([])

// 搜索表单
const searchForm = reactive({
  name: ''
})

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 分类列表（当前页数据）
const categoryList = ref<Category[]>([])

// 加载状态
const loading = ref(false)

// 错误状态
const error = ref(false)

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAdd = ref(true)

// 分类表单
const categoryForm = reactive({
  id: 0,
  name: ''
})

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
    const { data, total } = response.data
    console.log('分类列表数据:', data, total) // 调试信息
    allCategories.value = data.map((item: any) => ({
      id: item.ID,
      name: item.name,
      createdAt: item.CreatedAt || item.created_at || ''  // 修复创建时间字段
    }))
    pagination.total = total
    
    // 更新当前页数据
    updateCurrentPageData()
  } catch (err) {
    error.value = true
    ElMessage.error('获取分类列表失败')
    console.error('获取分类列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 更新当前页数据（前端分页）
const updateCurrentPageData = () => {
  // 应用搜索过滤
  let filteredCategories = allCategories.value
  if (searchForm.name) {
    filteredCategories = filteredCategories.filter(category => 
      category.name.includes(searchForm.name)
    )
  }
  
  // 更新总数
  pagination.total = filteredCategories.length
  
  // 计算当前页数据
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  categoryList.value = filteredCategories.slice(start, end)
}

// 处理搜索
const handleSearch = () => {
  pagination.currentPage = 1
  updateCurrentPageData()
}

// 处理重置
const handleReset = () => {
  searchForm.name = ''
  pagination.currentPage = 1
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
    name: ''
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
    name: row.name
  })
}

// 处理删除
const handleDelete = (row: Category) => {
  ElMessageBox.confirm(
    `确定要删除分类 "${row.name}" 吗？`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await categoryApi.deleteCategory(row.id)
      ElMessage.success('删除成功')
      getCategoryList()
    } catch (error) {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 提交分类表单
const submitCategoryForm = async (formData: {id: number, name: string}) => {
  try {
    if (isAdd.value) {
      // 新增分类
      await categoryApi.createCategory({
        name: formData.name
      })
      ElMessage.success('新增成功')
    } else {
      // 编辑分类
      await categoryApi.updateCategory(formData.id, {
        name: formData.name
      })
      ElMessage.success('修改成功')
    }
    dialogVisible.value = false
    getCategoryList()
  } catch (error) {
    ElMessage.error(isAdd.value ? '新增失败' : '修改失败')
    console.error(error)
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

// 组件挂载时获取数据
onMounted(() => {
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
</style>