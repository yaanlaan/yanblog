<template>
  <div class="media-manager">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span>媒体库管理</span>
            <span class="stats-info">
              共 {{ stats.totalFiles }} 个文件，{{ stats.totalDirs }} 个文件夹，{{ formatSize(stats.totalSize) }}
            </span>
          </div>
          <div class="header-actions">
            <el-button @click="toggleViewMode" :icon="viewMode === 'grid' ? List : Grid">
              {{ viewMode === 'grid' ? '列表视图' : '网格视图' }}
            </el-button>
            <el-button @click="uploadFiles" :icon="Upload" type="primary">上传文件</el-button>
            <el-button 
              v-if="selectedFiles.length > 0" 
              @click="batchDeleteSelected" 
              :icon="Delete" 
              type="danger"
            >
              批量删除 ({{ selectedFiles.length }})
            </el-button>
            <el-button @click="createFolderBtn" :icon="FolderAdd">新建文件夹</el-button>
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

      <!-- 拖拽上传区域 -->
      <div 
        class="upload-drop-zone" 
        @drop.prevent="handleDrop"
        @dragover.prevent
        @dragenter.prevent
        v-if="!loading"
      >
        <el-icon class="drop-icon"><Upload /></el-icon>
        <p>拖拽文件到此处上传，或点击上方按钮上传</p>
      </div>

      <!-- 列表视图 -->
      <el-table 
        v-if="viewMode === 'list'" 
        :data="files" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
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
        <el-table-column prop="ext" label="类型" width="100" />
        <el-table-column prop="modTime" label="修改时间" width="180" :formatter="formatTime" />
        <el-table-column label="操作" width="280">
          <template #default="scope">
            <el-button 
              v-if="!scope.row.isDir && scope.row.isImage" 
              type="primary" 
              link 
              @click="previewFile(scope.row)"
            >
              预览
            </el-button>
            <el-button 
              v-if="!scope.row.isDir" 
              type="success" 
              link 
              @click="copyFile(scope.row)"
            >
              复制
            </el-button>
            <el-button 
              type="warning" 
              link 
              @click="renameItem(scope.row)"
            >
              重命名
            </el-button>
            <el-button 
              type="info" 
              link 
              @click="moveItem(scope.row)"
            >
              移动
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

      <!-- 网格视图 -->
      <div v-else class="grid-view" v-loading="loading">
        <div 
          v-for="file in files" 
          :key="file.path"
          class="grid-item"
          :class="{ selected: selectedFiles.includes(file.path) }"
          @click="handleGridItemClick(file)"
          @contextmenu.prevent="showContextMenu($event, file)"
        >
          <div class="grid-item-checkbox" v-if="selectedFiles.length > 0">
            <el-checkbox 
              v-model="selectedFiles" 
              :value="file.path" 
              @change.stop
            />
          </div>
          <div class="grid-item-icon">
            <img 
              v-if="file.isImage && !file.isDir" 
              :src="file.thumbnail" 
              :alt="file.name"
              class="grid-thumb"
            />
            <el-icon v-else-if="file.isDir" class="dir-icon"><FolderOpened /></el-icon>
            <el-icon v-else class="file-icon-small"><Document /></el-icon>
          </div>
          <div class="grid-item-name">{{ file.name }}</div>
          <div class="grid-item-size">{{ file.isDir ? '文件夹' : formatSize(file.size) }}</div>
        </div>
      </div>

      <!-- 上下文菜单 -->
      <el-menu 
        v-if="contextMenuVisible" 
        :default-active="contextMenuActive"
        class="context-menu"
        :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
        @select="handleContextMenuSelect"
      >
        <el-menu-item v-if="!contextMenuFile?.isDir" index="preview">预览</el-menu-item>
        <el-menu-item v-if="!contextMenuFile?.isDir" index="copy">复制</el-menu-item>
        <el-menu-item index="rename">重命名</el-menu-item>
        <el-menu-item index="move">移动</el-menu-item>
        <el-menu-item index="delete" style="color: #f56c6c">删除</el-menu-item>
      </el-menu>
    </el-card>

    <!-- 图片预览弹窗 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="80%">
      <div class="preview-container">
        <img :src="previewUrl" alt="Preview" class="preview-image" />
        <div class="preview-info">
          <p>文件名: {{ previewFileName }}</p>
          <p>大小: {{ previewFileSize }}</p>
        </div>
      </div>
    </el-dialog>

    <!-- 移动文件弹窗 -->
    <el-dialog v-model="moveVisible" title="移动文件" width="400px">
      <el-form :model="moveForm" label-width="80px">
        <el-form-item label="目标位置">
          <el-select 
            v-model="moveForm.targetPath" 
            placeholder="选择目标目录"
            style="width: 100%"
          >
            <el-option label="根目录" value="" />
            <el-option 
              v-for="dir in availableDirs" 
              :key="dir.path" 
              :label="dir.name" 
              :value="dir.path" 
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="moveVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmMove">确定移动</el-button>
      </template>
    </el-dialog>

    <!-- 批量上传弹窗 -->
    <el-dialog v-model="uploadVisible" title="批量上传" width="500px">
      <div class="upload-dialog">
        <input 
          ref="fileInput"
          type="file" 
          multiple 
          class="upload-input"
          @change="handleFileSelect"
        />
        <div class="upload-area" @click="triggerFileInput">
          <el-icon class="upload-icon"><Upload /></el-icon>
          <p>点击或拖拽文件到此处</p>
          <p class="upload-tip">支持 JPG, PNG, GIF, PDF 等格式</p>
        </div>
        <div v-if="uploadFilesList.length > 0" class="upload-list">
          <div 
            v-for="(file, index) in uploadFilesList" 
            :key="index" 
            class="upload-item"
          >
            <span class="upload-filename">{{ file.name }}</span>
            <span class="upload-size">{{ formatSize(file.size) }}</span>
            <el-button size="small" @click="removeUploadFile(index)">移除</el-button>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="uploadVisible = false">取消</el-button>
        <el-button type="primary" @click="startUpload">开始上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onActivated, onDeactivated } from 'vue'
import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { 
  Folder, FolderOpened, Document, Refresh, Back, FolderAdd, 
  Upload, Delete, List, Grid 
} from '@element-plus/icons-vue'
import { fileApi } from '@/services/api'

interface FileInfo {
  name: string
  isDir: boolean
  path: string
  size: number
  ext: string
  modTime: string
  isImage: boolean
  thumbnail?: string
}

const loading = ref(false)
const files = ref<FileInfo[]>([])
const currentPath = ref('')
const viewMode = ref<'list' | 'grid'>('grid')
const previewVisible = ref(false)
const previewUrl = ref('')
const previewFileName = ref('')
const previewFileSize = ref('')
const selectedFiles = ref<string[]>([])
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuFile = ref<FileInfo | null>(null)
const moveVisible = ref(false)
const moveForm = ref({ targetPath: '' })
const moveSourcePath = ref('')
const uploadVisible = ref(false)
const uploadFilesList = ref<File[]>([])
const fileInput = ref<HTMLInputElement | null>(null)
const stats = ref({ totalFiles: 0, totalDirs: 0, totalSize: 0 })

const pathParts = computed(() => {
  return currentPath.value ? currentPath.value.split('/').filter(p => p) : []
})

const availableDirs = computed(() => {
  return files.value
    .filter(f => f.isDir)
    .map(f => ({ name: f.name, path: f.path }))
})

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

const formatTime = (timeStr: string) => {
  return new Date(timeStr).toLocaleString('zh-CN')
}

const refreshFiles = async () => {
  loading.value = true
  try {
    const res = await fileApi.getFiles(currentPath.value)
    if (res.data.status === 200) {
      files.value = res.data.data.sort((a: FileInfo, b: FileInfo) => {
        if (a.isDir === b.isDir) return a.name.localeCompare(b.name)
        return a.isDir ? -1 : 1
      })
    } else {
      ElMessage.error(res.data.message)
    }
    await fetchStats()
  } catch (error) {
    ElMessage.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res = await fileApi.getStorageStats()
    if (res.data.status === 200) {
      stats.value = {
        totalFiles: res.data.totalFiles,
        totalDirs: res.data.totalDirs,
        totalSize: res.data.totalSize
      }
    }
  } catch (error) {
    console.error('获取存储统计失败', error)
  }
}

const handleItemClick = (item: FileInfo) => {
  if (item.isDir) {
    currentPath.value = item.path
    refreshFiles()
  } else {
    previewFile(item)
  }
}

const handleGridItemClick = (item: FileInfo) => {
  if (selectedFiles.value.length > 0) {
    const index = selectedFiles.value.indexOf(item.path)
    if (index > -1) {
      selectedFiles.value.splice(index, 1)
    } else {
      selectedFiles.value.push(item.path)
    }
  } else {
    handleItemClick(item)
  }
}

const handleSelectionChange = (val: string[]) => {
  selectedFiles.value = val
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

const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid'
}

const createFolderBtn = () => {
  ElMessageBox.prompt('请输入文件夹名称', '新建文件夹', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^[^\\/:*?"<>|]+$/,
    inputErrorMessage: '包含非法字符'
  }).then(async ({ value }) => {
    try {
      const res = await fileApi.createFolder(currentPath.value, value)
      if (res.data.status === 200) {
        ElMessage.success('创建成功')
        refreshFiles()
      } else {
        ElMessage.error(res.data.message || '创建失败')
      }
    } catch (e) {
      ElMessage.error('网络错误')
    }
  }).catch(() => {})
}

const renameItem = (row: FileInfo) => {
  ElMessageBox.prompt('请输入新名称', '重命名', {
    inputValue: row.name,
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /^[^\\/:*?"<>|]+$/,
    inputErrorMessage: '包含非法字符'
  }).then(async ({ value }) => {
    if (value === row.name) return
    try {
      const res = await fileApi.renameFile(row.path, value)
      if (res.data.status === 200) {
        ElMessage.success('重命名成功')
        refreshFiles()
      } else {
        ElMessage.error(res.data.message || '重命名失败')
      }
    } catch (e) {
      ElMessage.error('网络错误')
    }
  }).catch(() => {})
}

const copyFile = (row: FileInfo) => {
  ElMessageBox.prompt('请输入目标目录（相对于 uploads）', '复制文件', {
    inputValue: currentPath.value,
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(async ({ value }) => {
    try {
      const res = await fileApi.copyFile(row.path, value)
      if (res.data.status === 200) {
        ElMessage.success('复制成功')
        refreshFiles()
      } else {
        ElMessage.error(res.data.message || '复制失败')
      }
    } catch (e) {
      ElMessage.error('网络错误')
    }
  }).catch(() => {})
}

const moveItem = (row: FileInfo) => {
  moveSourcePath.value = row.path
  moveForm.value.targetPath = ''
  moveVisible.value = true
}

const confirmMove = async () => {
  try {
    const res = await fileApi.moveFile(moveSourcePath.value, moveForm.value.targetPath)
    if (res.data.status === 200) {
      ElMessage.success('移动成功')
      moveVisible.value = false
      refreshFiles()
    } else {
      ElMessage.error(res.data.message || '移动失败')
    }
  } catch (e) {
    ElMessage.error('网络错误')
  }
}

const previewFile = (item: FileInfo) => {
  previewUrl.value = '/uploads/' + item.path
  previewFileName.value = item.name
  previewFileSize.value = formatSize(item.size)
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

const batchDeleteSelected = () => {
  ElMessageBox.confirm(
    `确定要删除选中的 ${selectedFiles.value.length} 个文件/文件夹吗？`,
    '批量删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const res = await fileApi.batchDeleteFiles(selectedFiles.value)
      if (res.data.status === 200) {
        ElNotification.success({
          title: '批量删除完成',
          message: `成功删除 ${res.data.successCount} 个，失败 ${res.data.failCount} 个`
        })
        selectedFiles.value = []
        refreshFiles()
      } else {
        ElMessage.error(res.data.message)
      }
    } catch (error) {
      ElMessage.error('批量删除失败')
    }
  }).catch(() => {})
}

const showContextMenu = (event: MouseEvent, file: FileInfo) => {
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  contextMenuFile.value = file
  contextMenuVisible.value = true
}

const handleContextMenuSelect = (index: string) => {
  if (!contextMenuFile.value) return
  
  switch (index) {
    case 'preview':
      previewFile(contextMenuFile.value)
      break
    case 'copy':
      copyFile(contextMenuFile.value)
      break
    case 'rename':
      renameItem(contextMenuFile.value)
      break
    case 'move':
      moveItem(contextMenuFile.value)
      break
    case 'delete':
      deleteItem(contextMenuFile.value)
      break
  }
  contextMenuVisible.value = false
}

const closeContextMenu = (event: MouseEvent) => {
  if ((event.target as HTMLElement).closest('.context-menu')) return
  contextMenuVisible.value = false
}

const uploadFiles = () => {
  uploadVisible.value = true
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    uploadFilesList.value = [...uploadFilesList.value, ...Array.from(target.files)]
  }
  target.value = ''
}

const removeUploadFile = (index: number) => {
  uploadFilesList.value.splice(index, 1)
}

const handleDrop = (event: DragEvent) => {
  const files = event.dataTransfer?.files
  if (files) {
    uploadFilesList.value = [...uploadFilesList.value, ...Array.from(files)]
  }
}

const startUpload = async () => {
  if (uploadFilesList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  const formData = new FormData()
  uploadFilesList.value.forEach(file => {
    formData.append('files', file)
  })
  formData.append('dir', currentPath.value)

  try {
    const res = await fileApi.batchUploadFiles(formData)
    if (res.data.status === 200) {
      ElNotification.success({
        title: '批量上传完成',
        message: `成功上传 ${res.data.successCount} 个，失败 ${res.data.failCount} 个`
      })
      uploadFilesList.value = []
      uploadVisible.value = false
      refreshFiles()
    } else {
      ElMessage.error(res.data.message)
    }
  } catch (error) {
    ElMessage.error('上传失败')
  }
}

onMounted(() => {
  refreshFiles()
  document.addEventListener('click', closeContextMenu)
})

onActivated(() => {
  refreshFiles()
})

onDeactivated(() => {
  document.removeEventListener('click', closeContextMenu)
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

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stats-info {
  font-size: 14px;
  color: #909399;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.breadcrumb {
  margin-bottom: 15px;
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

.upload-drop-zone {
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  padding: 30px;
  text-align: center;
  margin-bottom: 15px;
  transition: all 0.3s;
}

.upload-drop-zone:hover {
  border-color: #409eff;
  background-color: #f0f5ff;
}

.drop-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 10px;
}

.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 15px;
  padding: 10px 0;
}

.grid-item {
  position: relative;
  background: #fafafa;
  border-radius: 8px;
  padding: 15px;
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
}

.grid-item:hover {
  background: #f0f5ff;
  transform: translateY(-2px);
}

.grid-item.selected {
  background: #e6f7ff;
  border: 2px solid #409eff;
}

.grid-item-checkbox {
  position: absolute;
  top: 8px;
  left: 8px;
}

.grid-item-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.grid-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.dir-icon {
  font-size: 48px;
  color: #e6a23c;
}

.file-icon-small {
  font-size: 40px;
  color: #909399;
}

.grid-item-name {
  font-size: 13px;
  color: #606266;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.grid-item-size {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.context-menu {
  position: fixed;
  z-index: 9999;
  min-width: 120px;
}

.preview-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.preview-image {
  max-width: 100%;
  max-height: 600px;
  border-radius: 8px;
}

.preview-info {
  margin-top: 15px;
  text-align: center;
  color: #606266;
}

.upload-dialog {
  padding: 20px;
}

.upload-input {
  display: none;
}

.upload-area {
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.upload-area:hover {
  border-color: #409eff;
  background-color: #f0f5ff;
}

.upload-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 10px;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.upload-list {
  margin-top: 20px;
  max-height: 200px;
  overflow-y: auto;
}

.upload-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 8px;
}

.upload-filename {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.upload-size {
  margin-right: 10px;
  color: #909399;
}
</style>