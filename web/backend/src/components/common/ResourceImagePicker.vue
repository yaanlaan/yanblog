<template>
  <div class="resource-image-picker">
    <div class="picker-main">
      <!-- 预览与直接上传触发框 -->
      <div 
        class="preview-box" 
        :class="[aspectRatio, { 'has-image': !!modelValue }]"
      >
        <template v-if="modelValue">
          <img :src="modelValue" :alt="label || '图片预览'" class="preview-img" />
          <div class="preview-overlay">
            <el-tooltip content="在新窗口查看原图" placement="top">
              <el-button 
                type="primary" 
                circle 
                size="small" 
                :icon="ZoomIn" 
                @click="openPreview" 
              />
            </el-tooltip>
            <el-tooltip content="复制图片链接" placement="top">
              <el-button 
                type="info" 
                circle 
                size="small" 
                :icon="Link" 
                @click="copyUrl" 
              />
            </el-tooltip>
            <el-tooltip content="移除/清空图片" placement="top">
              <el-button 
                type="danger" 
                circle 
                size="small" 
                :icon="Delete" 
                @click="clearImage" 
              />
            </el-tooltip>
          </div>
        </template>
        
        <el-upload
          v-else
          class="upload-trigger"
          action="/api/v1/upload"
          :data="{ type: uploadType }"
          :headers="uploadHeaders"
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          accept="image/*"
        >
          <div class="upload-placeholder">
            <el-icon class="upload-icon"><Plus /></el-icon>
            <span class="upload-text">点击或拖拽上传</span>
          </div>
        </el-upload>
      </div>

      <!-- 右侧控制区 -->
      <div class="picker-controls">
        <!-- 核心操作按钮组 -->
        <div class="btn-group">
          <el-upload
            action="/api/v1/upload"
            :data="{ type: uploadType }"
            :headers="uploadHeaders"
            :show-file-list="false"
            :before-upload="beforeUpload"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            accept="image/*"
            style="display: inline-block;"
          >
            <el-button type="primary" :icon="Upload" :loading="uploading">
              {{ modelValue ? '替换上传' : '上传图片' }}
            </el-button>
          </el-upload>

          <el-button :icon="Picture" @click="openMediaDialog">
            从媒体库选择
          </el-button>

          <el-button v-if="modelValue" :icon="Delete" type="danger" plain @click="clearImage">
            清空
          </el-button>
        </div>

        <!-- 手动输入/外链 URL -->
        <div class="url-input-wrap">
          <el-input
            v-model="inputUrl"
            placeholder="支持粘贴外部网络图片或图床 URL (以 http/https 或 / 开头)"
            clearable
            @change="handleUrlChange"
            @blur="handleUrlChange"
          >
            <template #prepend>图片地址</template>
          </el-input>
        </div>

        <!-- 快捷预设选项 -->
        <div v-if="presetOptions && presetOptions.length > 0" class="presets-wrap">
          <span class="preset-label">快捷选择：</span>
          <el-tag
            v-for="item in presetOptions"
            :key="item.value"
            class="preset-tag"
            :effect="modelValue === item.value ? 'dark' : 'plain'"
            type="info"
            @click="selectPreset(item.value)"
          >
            {{ item.label }}
          </el-tag>
        </div>

        <div v-if="tip" class="picker-tip">{{ tip }}</div>
      </div>
    </div>

    <!-- 媒体库图片选择弹窗 -->
    <el-dialog
      v-model="mediaDialogVisible"
      title="从媒体库选择图片"
      width="820px"
      append-to-body
      destroy-on-close
      class="media-select-dialog"
    >
      <div class="media-dialog-content" v-loading="mediaLoading">
        <!-- 顶部工具栏与面包屑 -->
        <div class="media-nav-bar">
          <div class="breadcrumbs">
            <span class="nav-crumb" @click="loadMediaFiles('')">uploads</span>
            <template v-for="(part, idx) in pathParts" :key="idx">
              <span class="crumb-sep">/</span>
              <span class="nav-crumb" @click="navigateToPart(idx)">{{ part }}</span>
            </template>
          </div>
          <div class="nav-actions">
            <el-button 
              size="small" 
              :icon="Back" 
              :disabled="!currentMediaDir" 
              @click="goUpDir"
            >
              返回上级
            </el-button>
            <el-button size="small" :icon="Refresh" @click="loadMediaFiles(currentMediaDir)">
              刷新
            </el-button>
          </div>
        </div>

        <!-- 文件网格 -->
        <div class="media-grid" v-if="mediaFiles.length > 0">
          <div
            v-for="item in mediaFiles"
            :key="item.path"
            class="media-item"
            :class="{
              'is-dir': item.isDir,
              'is-image': item.isImage,
              'is-selected': selectedMediaPath === getFileUrl(item)
            }"
            @click="handleMediaItemClick(item)"
          >
            <!-- 文件夹视图 -->
            <template v-if="item.isDir">
              <div class="dir-icon-wrap">
                <el-icon :size="36"><Folder /></el-icon>
              </div>
              <div class="media-name" :title="item.name">{{ item.name }}</div>
            </template>

            <!-- 图片视图 -->
            <template v-else-if="item.isImage">
              <div class="img-thumb-wrap">
                <img :src="getFileUrl(item)" :alt="item.name" loading="lazy" />
                <div v-if="selectedMediaPath === getFileUrl(item)" class="selected-badge">
                  <el-icon><Check /></el-icon>
                </div>
              </div>
              <div class="media-name" :title="item.name">{{ item.name }}</div>
              <div class="media-size">{{ formatSize(item.size) }}</div>
            </template>

            <!-- 其他文件 -->
            <template v-else>
              <div class="other-icon-wrap">
                <el-icon :size="32"><Document /></el-icon>
              </div>
              <div class="media-name" :title="item.name">{{ item.name }}</div>
            </template>
          </div>
        </div>

        <el-empty v-else description="当前目录下无图片或文件" :image-size="80" />
      </div>

      <template #footer>
        <div class="dialog-footer">
          <div class="selected-preview-hint">
            <span v-if="selectedMediaPath">已选中：{{ selectedMediaPath }}</span>
            <span v-else class="text-muted">请点击选择一张图片</span>
          </div>
          <div class="footer-btns">
            <el-button @click="mediaDialogVisible = false">取消</el-button>
            <el-button 
              type="primary" 
              :disabled="!selectedMediaPath" 
              @click="confirmMediaSelect"
            >
              确定选择
            </el-button>
          </div>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  Upload, Picture, Delete, ZoomIn, Link, Refresh,
  Folder, Document, Back, Check, Plus
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { UploadProps } from 'element-plus'
import { fileApi } from '@/services/api'

interface PresetOption {
  label: string
  value: string
}

interface FileItem {
  name: string
  isDir: boolean
  path: string
  size: number
  ext: string
  isImage: boolean
  thumbnail?: string
}

const props = withDefaults(
  defineProps<{
    modelValue?: string
    label?: string
    tip?: string
    uploadType?: string
    presetOptions?: PresetOption[]
    aspectRatio?: 'square' | 'wide' | 'auto'
  }>(),
  {
    modelValue: '',
    label: '',
    tip: '',
    uploadType: 'system',
    presetOptions: () => [],
    aspectRatio: 'square'
  }
)

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void
  (e: 'change', val: string): void
}>()

const inputUrl = ref(props.modelValue || '')
const uploading = ref(false)

watch(
  () => props.modelValue,
  (val) => {
    inputUrl.value = val || ''
  }
)

const uploadHeaders = computed(() => ({
  Authorization: 'Bearer ' + (localStorage.getItem('token') || '')
}))

const updateValue = (val: string) => {
  emit('update:modelValue', val)
  emit('change', val)
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  if (file.size / 1024 / 1024 > 10) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  uploading.value = true
  return true
}

const handleUploadSuccess: UploadProps['onSuccess'] = (res) => {
  uploading.value = false
  if (res.status === 200 && res.url) {
    updateValue(res.url)
    ElMessage.success('图片上传并应用成功！')
  } else {
    ElMessage.error(res.message || '上传失败')
  }
}

const handleUploadError = () => {
  uploading.value = false
  ElMessage.error('上传失败，请检查网络或后端状态')
}

const handleUrlChange = () => {
  const trimmed = inputUrl.value.trim()
  updateValue(trimmed)
}

const selectPreset = (val: string) => {
  inputUrl.value = val
  updateValue(val)
  ElMessage.success('已切换为预设资源')
}

const clearImage = () => {
  inputUrl.value = ''
  updateValue('')
  ElMessage.info('已移除图片')
}

const openPreview = () => {
  if (props.modelValue) {
    window.open(props.modelValue, '_blank')
  }
}

const copyUrl = async () => {
  if (!props.modelValue) return
  try {
    await navigator.clipboard.writeText(props.modelValue)
    ElMessage.success('链接已复制到剪贴板')
  } catch {
    ElMessage.info('复制失败，链接为: ' + props.modelValue)
  }
}

// --- 媒体库选择器弹窗 ---
const mediaDialogVisible = ref(false)
const mediaLoading = ref(false)
const mediaFiles = ref<FileItem[]>([])
const currentMediaDir = ref('')
const selectedMediaPath = ref('')

const pathParts = computed(() => {
  return currentMediaDir.value ? currentMediaDir.value.split('/').filter(p => p) : []
})

const openMediaDialog = () => {
  selectedMediaPath.value = props.modelValue || ''
  mediaDialogVisible.value = true
  // 默认根据上传类型定位到常用子目录，例如 defaults 或 system
  if (!currentMediaDir.value) {
    currentMediaDir.value = props.uploadType === 'avatar' ? 'avatar' : (props.uploadType === 'system' ? 'system' : '')
  }
  loadMediaFiles(currentMediaDir.value)
}

const loadMediaFiles = async (dir: string) => {
  mediaLoading.value = true
  currentMediaDir.value = dir
  try {
    const res = await fileApi.getFiles(dir)
    if (res.data?.status === 200 && Array.isArray(res.data.data)) {
      mediaFiles.value = res.data.data.sort((a: FileItem, b: FileItem) => {
        if (a.isDir === b.isDir) return a.name.localeCompare(b.name)
        return a.isDir ? -1 : 1
      })
    } else {
      mediaFiles.value = []
    }
  } catch (err) {
    console.error('Failed to load media files:', err)
    ElMessage.error('读取媒体库失败')
    mediaFiles.value = []
  } finally {
    mediaLoading.value = false
  }
}

const getFileUrl = (item: FileItem) => {
  if (item.thumbnail) return item.thumbnail
  return `/uploads/${item.path}`
}

const handleMediaItemClick = (item: FileItem) => {
  if (item.isDir) {
    loadMediaFiles(item.path)
  } else if (item.isImage) {
    selectedMediaPath.value = getFileUrl(item)
  }
}

const navigateToPart = (idx: number) => {
  const target = pathParts.value.slice(0, idx + 1).join('/')
  loadMediaFiles(target)
}

const goUpDir = () => {
  if (!currentMediaDir.value) return
  const parts = currentMediaDir.value.split('/').filter(p => p)
  parts.pop()
  loadMediaFiles(parts.join('/'))
}

const confirmMediaSelect = () => {
  if (selectedMediaPath.value) {
    updateValue(selectedMediaPath.value)
    mediaDialogVisible.value = false
    ElMessage.success('已选用媒体库图片')
  }
}

const formatSize = (bytes: number) => {
  if (!bytes) return '0 B'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}
</script>

<style scoped>
.resource-image-picker {
  width: 100%;
}

.picker-main {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

/* 预览框 */
.preview-box {
  position: relative;
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  overflow: hidden;
  background-color: #fcfcfc;
  /* 棋盘格衬底，方便看透透明 PNG/SVG */
  background-image: 
    linear-gradient(45deg, #f0f0f0 25%, transparent 25%), 
    linear-gradient(-45deg, #f0f0f0 25%, transparent 25%), 
    linear-gradient(45deg, transparent 75%, #f0f0f0 75%), 
    linear-gradient(-45deg, transparent 75%, #f0f0f0 75%);
  background-size: 16px 16px;
  background-position: 0 0, 0 8px, 8px -8px, -8px 0px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.25s ease;
}

.preview-box.square {
  width: 130px;
  height: 130px;
}

.preview-box.wide {
  width: 220px;
  height: 124px; /* 16:9 比例 */
}

.preview-box.auto {
  min-width: 130px;
  max-width: 220px;
  height: 130px;
}

.preview-box:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
}

/* 悬停工具条 */
.preview-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.preview-box:hover .preview-overlay {
  opacity: 1;
}

.upload-trigger {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.upload-trigger :deep(.el-upload) {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  color: var(--el-text-color-secondary);
}

.upload-icon {
  font-size: 26px;
  color: #909399;
}

.upload-text {
  font-size: 12px;
  color: #909399;
}

/* 控制区 */
.picker-controls {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.btn-group {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.url-input-wrap {
  max-width: 520px;
}

.presets-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: 2px;
}

.preset-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.preset-tag {
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.preset-tag:hover {
  transform: translateY(-1px);
  color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

.picker-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.4;
}

/* 媒体库选择弹窗 */
.media-nav-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #f8f9fa;
  border-radius: 6px;
  margin-bottom: 15px;
}

.breadcrumbs {
  font-size: 13px;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav-crumb {
  color: var(--el-color-primary);
  cursor: pointer;
  font-weight: 500;
}

.nav-crumb:hover {
  text-decoration: underline;
}

.crumb-sep {
  color: #909399;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
  gap: 12px;
  max-height: 420px;
  overflow-y: auto;
  padding: 4px;
}

.media-item {
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #fff;
  position: relative;
}

.media-item:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.media-item.is-selected {
  border-color: var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-5);
}

.dir-icon-wrap, .other-icon-wrap {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #e6a23c;
}

.other-icon-wrap {
  color: #909399;
}

.img-thumb-wrap {
  width: 90px;
  height: 90px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  border-radius: 4px;
  overflow: hidden;
}

.img-thumb-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.selected-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--el-color-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.media-name {
  font-size: 12px;
  color: var(--el-text-color-primary);
  margin-top: 6px;
  width: 100%;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.media-size {
  font-size: 11px;
  color: #909399;
  margin-top: 2px;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.selected-preview-hint {
  font-size: 13px;
  color: var(--el-color-primary);
  max-width: 500px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-muted {
  color: #909399;
}
</style>
