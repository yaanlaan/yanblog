<template>
  <div class="tag-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>标签管理</span>
          <el-button type="primary" @click="handleAdd">新增标签</el-button>
        </div>
      </template>
      
      <!-- 标签表格 -->
      <el-table 
        :data="tagList" 
        border 
        style="width: 100%" 
        v-loading="loading"
        :empty-text="error ? '数据加载失败，请检查网络连接' : '暂无数据'"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="标签名称" />
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="handleDelete(scope.row)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.pagenum"
        v-model:page-size="pagination.pagesize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        class="pagination"
      />
    </el-card>
    
    <!-- 标签编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="30%"
    >
      <el-form
        ref="tagFormRef"
        :model="tagForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="标签名称" prop="name">
          <el-form-item label-width="0">
             <el-input v-model="tagForm.name" placeholder="请输入标签名称" />
          </el-form-item>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitTagForm">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { tagApi } from '@/services/api'

interface Tag {
  id: number
  name: string
}

const loading = ref(false)
const error = ref(false)
const tagList = ref<Tag[]>([])
const pagination = reactive({
  pagenum: 1,
  pagesize: 10,
  total: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const tagFormRef = ref<FormInstance>()
const tagForm = reactive({
  id: 0,
  name: ''
})

const rules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 20, message: '长度在 1 到 20 个字符', trigger: 'blur' }
  ]
}

// 获取标签列表
const fetchTags = async () => {
  loading.value = true
  error.value = false
  try {
    const res = await tagApi.getTags({
      pagenum: pagination.pagenum,
      pagesize: pagination.pagesize
    })
    console.log('fetchTags res:', res)
    if (res.status === 200) {
      if (res.data.status === 200) {
        tagList.value = res.data.data
        pagination.total = res.data.total
      } else {
        ElMessage.error(res.data.message || '获取标签列表失败')
      }
    }
  } catch (err) {
    console.error(err)
    error.value = true
    ElMessage.error('获取列表失败')
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (val: number) => {
  pagination.pagesize = val
  fetchTags()
}

const handleCurrentChange = (val: number) => {
  pagination.pagenum = val
  fetchTags()
}

// 新增标签
const handleAdd = () => {
  dialogTitle.value = '新增标签'
  tagForm.id = 0
  tagForm.name = ''
  dialogVisible.value = true
}

// 编辑标签
const handleEdit = (row: Tag) => {
  dialogTitle.value = '编辑标签'
  tagForm.id = row.id
  tagForm.name = row.name
  dialogVisible.value = true
}

// 处理表单提交
const submitTagForm = async () => {
  if (!tagFormRef.value) return
  
  await tagFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        let res
        if (tagForm.id === 0) {
          // 新增
          res = await tagApi.addTag({ name: tagForm.name })
        } else {
          // 编辑
          res = await tagApi.updateTag(tagForm.id, { name: tagForm.name })
        }

        if (res.status === 200 && res.data.status === 200) {
            ElMessage.success(tagForm.id === 0 ? '新增成功' : '更新成功')
            dialogVisible.value = false
            fetchTags()
        } else {
            ElMessage.error(res.data.message || '操作失败')
        }
      } catch (err) {
        console.error(err)
        ElMessage.error('操作失败')
      }
    }
  })
}

// 删除标签
const handleDelete = (row: Tag) => {
  ElMessageBox.confirm(
    '确定要删除该标签吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const res = await tagApi.deleteTag(row.id)
      if (res.status === 200 && res.data.status === 200) {
        ElMessage.success('删除成功')
        fetchTags()
      } else {
        ElMessage.error(res.data.message || '删除失败')
      }
    } catch (err) {
      console.error(err)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

onMounted(() => {
  fetchTags()
})
</script>

<style scoped>
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
