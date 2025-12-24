<template>
  <div class="user-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleAdd">新增用户</el-button>
        </div>
      </template>
      
      <!-- 搜索表单 -->
      <UserSearchForm
        :model-value="searchForm"
        @update:modelValue="handleSearchFormUpdate"
        @search="handleSearch"
        @reset="handleReset"
      />
      
      <!-- 用户表格 -->
      <el-table 
        :data="userList" 
        border 
        style="width: 100%" 
        v-loading="loading"
        :empty-text="error ? '数据加载失败，请检查网络连接' : '暂无数据'"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column label="角色" width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.role === 1" type="danger">超级管理员</el-tag>
            <el-tag v-else-if="scope.row.role === 2" type="warning">管理员</el-tag>
            <el-tag v-else type="info">普通用户</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <UserActions
              :user="scope.row"
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
    
    <!-- 用户编辑表单 -->
    <UserForm
      v-model="dialogVisible"
      :title="dialogTitle"
      :user="userForm"
      :is-add="isAdd"
      @submit="submitUserForm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import { userApi } from '@/services/api'
import UserSearchForm from '@/components/user/UserSearchForm.vue'
import UserActions from '@/components/user/UserActions.vue'
import UserForm from '@/components/user/UserForm.vue'

// 路由
const router = useRouter()
const route = useRoute()

// 用户数据类型
interface User {
  id: number
  username: string
  role: number
  createdAt: string
}

// 所有用户数据（用于前端分页）
const allUsers = ref<User[]>([])

// 搜索表单
const searchForm = reactive({
  username: '',
  role: undefined as number | undefined
})

// 更新当前页数据（前端分页和筛选）
const updateCurrentPageData = () => {
  // 应用搜索和筛选
  let filteredUsers = [...allUsers.value]
  
  // 用户名搜索
  if (searchForm.username) {
    filteredUsers = filteredUsers.filter(user => 
      user.username.toLowerCase().includes(searchForm.username.toLowerCase())
    )
  }
  
  // 角色筛选
  if (searchForm.role !== undefined && searchForm.role !== null) {
    const role = Number(searchForm.role)
    filteredUsers = filteredUsers.filter(user => 
      user.role === role
    )
  }
  
  // 更新总数
  pagination.total = filteredUsers.length
  
  // 计算当前页数据
  const start = (pagination.currentPage - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  userList.value = filteredUsers.slice(start, end)
  
  // 更新URL参数
  updateUrlParams()
}

// 更新URL参数
const updateUrlParams = () => {
  const query: Record<string, string | undefined> = {}
  
  if (searchForm.username) {
    query.username = searchForm.username
  }
  
  if (searchForm.role !== undefined) {
    query.role = searchForm.role.toString()
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
  const paramKeys = ['username', 'role', 'page', 'pageSize']
  for (const key of paramKeys) {
    if (query[key] !== currentQuery[key]) {
      needUpdate = true
      break
    }
  }
  
  // 检查是否有额外的参数需要移除
  for (const key in currentQuery) {
    if (!['username', 'role', 'page', 'pageSize'].includes(key) && query[key] === undefined) {
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
  searchForm.username = (query.username as string) || ''
  searchForm.role = query.role ? Number(query.role) : undefined
  
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

// 用户列表（当前页数据）
const userList = ref<User[]>([])

// 加载状态
const loading = ref(false)

// 错误状态
const error = ref(false)

// 对话框相关
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isAdd = ref(true)

// 用户表单
const userForm = reactive({
  id: 0,
  username: '',
  password: '',
  role: 3 // 默认为普通用户
})

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

// 获取用户列表（获取所有数据）
const getUserList = async () => {
  loading.value = true
  error.value = false
  try {
    const response = await userApi.getUsers({
      pagesize: -1,  // 获取所有数据
      pagenum: -1    // 获取所有数据
    })
    
    // 解析后端返回的数据
    const { data } = response.data
    allUsers.value = data.map((item: any) => ({
      id: item.ID,
      username: item.username,
      role: item.role,
      createdAt: item.CreatedAt
    }))
    
    // 更新当前页数据
    updateCurrentPageData()
  } catch (err) {
    error.value = true
    ElMessage.error('获取用户列表失败')
    console.error('获取用户列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 处理搜索表单更新
const handleSearchFormUpdate = (value: {username: string, role: number | undefined}) => {
  Object.assign(searchForm, value)
}

// 处理搜索
const handleSearch = () => {
  pagination.currentPage = 1
  updateCurrentPageData()
}

// 处理重置
const handleReset = () => {
  searchForm.username = ''
  searchForm.role = undefined
  pagination.currentPage = 1
  updateCurrentPageData()
}

// 处理新增
const handleAdd = () => {
  dialogTitle.value = '新增用户'
  isAdd.value = true
  dialogVisible.value = true
  // 重置表单
  Object.assign(userForm, {
    id: 0,
    username: '',
    password: '',
    role: 3 // 默认为普通用户
  })
}

// 处理编辑
const handleEdit = (row: User) => {
  dialogTitle.value = '编辑用户'
  isAdd.value = false
  dialogVisible.value = true
  // 填充表单数据
  Object.assign(userForm, {
    id: row.id,
    username: row.username,
    password: '',
    role: row.role
  })
}

// 处理删除
const handleDelete = (row: User) => {
  if (row.role === 1) {
    ElMessage.warning('不能删除超级管理员')
    return
  }
  
  ElMessageBox.confirm(
    `确定要删除用户 "${row.username}" 吗？`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await userApi.deleteUser(row.id)
      ElMessage.success('删除成功')
      getUserList()
    } catch (error) {
      ElMessage.error('删除失败')
      console.error(error)
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 提交用户表单
const submitUserForm = async (formData: { username: string; role: number; password?: string }) => {
  try {
    let res;
    if (isAdd.value) {
      // 新增用户
      res = await userApi.createUser({
        username: formData.username,
        password: formData.password || '123456', // 默认密码
        role: formData.role
      })
    } else {
      // 编辑用户
      res = await userApi.updateUser(userForm.id, {
        username: formData.username,
        role: formData.role,
        password: formData.password // 传递密码
      })
    }

    // 检查后端返回的状态码
    if (res.data.status !== 200) {
      ElMessage.error(res.data.message || (isAdd.value ? '用户创建失败' : '用户更新失败'))
      return
    }
    
    ElMessage.success(isAdd.value ? '用户创建成功' : '用户更新成功')
    
    // 关闭对话框
    dialogVisible.value = false
    
    // 重新加载数据
    await getUserList()
  } catch (error: any) {
    console.error('提交用户表单失败:', error)
    ElMessage.error(isAdd.value ? '用户创建失败' : '用户更新失败')
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
  getUserList()
})
</script>

<style scoped>
.user-list {
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