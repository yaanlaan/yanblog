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
        v-model="searchForm"
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
        <el-table-column prop="role" label="角色">
          <template #default="scope">
            <el-tag v-if="scope.row.role === 2" type="danger">管理员</el-tag>
            <el-tag v-else type="success">普通用户</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" />
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
      :is-add="isAdd"
      :user="userForm"
      @submit="submitUserForm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi } from '@/services/api'
import UserSearchForm from '@/components/user/UserSearchForm.vue'
import UserActions from '@/components/user/UserActions.vue'
import UserForm from '@/components/user/UserForm.vue'

// 用户数据类型
interface User {
  id: number
  username: string
  role: number
  createdAt: string
}

// 搜索表单
const searchForm = reactive({
  username: '',
  role: undefined as number | undefined
})

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 用户列表
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

// 获取用户列表
const getUserList = async () => {
  loading.value = true
  error.value = false
  try {
    const response = await userApi.getUsers({
      pagesize: pagination.pageSize,
      pagenum: pagination.currentPage
    })
    
    // 解析后端返回的数据
    const { data, total } = response.data
    console.log('用户列表数据:', data, total) // 调试信息
    userList.value = data.map((item: any) => ({
      id: item.ID,
      username: item.username,
      role: item.role,
      createdAt: item.CreatedAt
    }))
    pagination.total = total
  } catch (err) {
    error.value = true
    ElMessage.error('获取用户列表失败')
    console.error('获取用户列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = () => {
  pagination.currentPage = 1
  getUserList()
}

// 处理重置
const handleReset = () => {
  pagination.currentPage = 1
  getUserList()
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
const submitUserForm = async (formData: {id: number, username: string, password: string, role: number}) => {
  try {
    if (isAdd.value) {
      // 新增用户
      await userApi.createUser({
        username: formData.username,
        password: formData.password,
        role: formData.role
      })
      ElMessage.success('新增成功')
    } else {
      // 编辑用户
      await userApi.updateUser(formData.id, {
        username: formData.username,
        role: formData.role
      })
      ElMessage.success('修改成功')
    }
    dialogVisible.value = false
    getUserList()
  } catch (error) {
    ElMessage.error(isAdd.value ? '新增失败' : '修改失败')
    console.error(error)
  }
}

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  pagination.pageSize = val
  getUserList()
}

// 处理当前页变化
const handleCurrentChange = (val: number) => {
  pagination.currentPage = val
  getUserList()
}

// 组件挂载时获取数据
onMounted(() => {
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