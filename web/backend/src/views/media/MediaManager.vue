<template>
  <div class="media-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>媒体库管理</span>
          <div class="header-actions">
            <el-button @click="refreshFiles" :icon="Refresh">刷新</el-button>
            <el-button @click="goUp" :disabled="currentPath === ''" :icon="Back">返回上级</el-button>
          </div>
        </div>
      </template>

      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item><a @click="navigateTo('')">uploads</a></el-breadcrumb-item>
          <el-breadcrumb-item v-for="(part, index) in pathParts" :key="index">
            <a @click="navigateToPart(index)">{{ part }}</a>
          </el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <el-table :data="files" style="width: 100%" v-loading="loading">
        <el-table-column label="名称" min-width="200">
          <template #default="scope">
            <div class="file-name" @click="handleItemClick(scope.row)">
              <el-icon v-if="scope.row.isDir" class="file-icon"><Folder /></el-icon>
              <el-icon v-else class="file-icon"><Document /></el-icon>
              <span>{{ scope.row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120" :formatter="formatSize" />
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button 
              v-if="!scope.row.isDir" 
              type="primary" 
              link 
              @click="previewFile(scope.row)"
            >
              预览
            </el-button>
            <el-popconfirm 
              title="确定要删除吗？" 
              @confirm="deleteItem(scope.row)"
            >
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 图片预览弹窗 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="50%">
      <div class="preview-container">
        <img :src="previewUrl" alt="Preview" style="max-width: 100%;" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Folder, Document, Refresh, Back } from '@element-plus/icons-vue'
import { fileApi } from '@/services/api'

interface FileInfo {
  name: string
  isDir: boolean
  path: string
  size: int
}

const loading = ref(false)
const files = ref<FileInfo[]>([])
const currentPath = ref('')
const previewVisible = ref(false)
const previewUrl = ref('')

const pathParts = computed(() => {
  return currentPath.value ? currentPath.value.split('/').filter(p => p) : []
})

const formatSize = (row: FileInfo) => {
  if (row.isDir) return '-'
  const size = row.size
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(2) + ' KB'
  return (size / (1024 * 1024)).toFixed(2) + ' MB'
}

const refreshFiles = async () => {
  loading.value = true
  try {
    const res = await fileApi.getFiles(currentPath.value)
    if (res.data.status === 200) {
      // 排序：文件夹在前，文件在后
      files.value = res.data.data.sort((a: FileInfo, b: FileInfo) => {
        if (a.isDir === b.isDir) return a.name.localeCompare(b.name)
        return a.isDir ? -1 : 1
      })
    } else {
      ElMessage.error(res.data.message)
    }
  } catch (error) {
    ElMessage.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

const handleItemClick = (item: FileInfo) => {
  if (item.isDir) {
    // 如果是目录，进入目录
    // 注意：后端返回的 path 是相对 uploads 的路径，例如 "articles/title"
    // 但我们前端维护 currentPath 也是相对 uploads 的
    // item.path 已经是相对路径了，直接用
    currentPath.value = item.path
    refreshFiles()
  } else {
    previewFile(item)
  }
}

const goUp = () => {
  if (!currentPath.value) return
  const parts = currentPath.value.split('/')
  parts.pop()
  currentPath.value = parts.join('/')
  refreshFiles()
}

const navigateTo = (path: string) => {
  currentPath.value = path
  refreshFiles()
}

const navigateToPart = (index: number) => {
  const parts = pathParts.value.slice(0, index + 1)
  currentPath.value = parts.join('/')
  refreshFiles()
}

const previewFile = (item: FileInfo) => {
  // 假设是图片，直接预览
  // 路径需要加上 /uploads/ 前缀
  previewUrl.value = '/uploads/' + item.path
  previewVisible.value = true
}

const deleteItem = async (item: FileInfo) => {
  try {
    const res = await fileApi.deleteFile(item.path)
    if (res.data.status === 200) {
      ElMessage.success('删除成功')
      refreshFiles()
    } else {
      ElMessage.error(res.data.message)
    }
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  refreshFiles()
})
</script>

<style scoped>
.media-manager {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.breadcrumb {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.breadcrumb a {
  cursor: pointer;
  color: #606266;
}

.breadcrumb a:hover {
  color: #409eff;
}

.file-name {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.file-icon {
  margin-right: 8px;
  font-size: 18px;
}

.file-name:hover {
  color: #409eff;
}

.preview-container {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
