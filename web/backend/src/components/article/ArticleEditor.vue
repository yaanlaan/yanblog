<template>
  <div class="markdown-editor-container">
    <div class="editor-header">
      <div class="header-left">
        <el-button @click="togglePreview" :icon="previewOnly ? Edit : View">
          {{ previewOnly ? '编辑模式' : '预览模式' }}
        </el-button>
        <el-button @click="toggleFullscreen" :icon="isFullscreen ? Fold : Expand">
          {{ isFullscreen ? '退出全屏' : '全屏编辑' }}
        </el-button>
      </div>
      <div class="header-right">
        <span v-if="saveStatus" class="save-status">
          <el-icon v-if="saveStatus === 'saving'" class="status-icon"><Loading /></el-icon>
          <el-icon v-else-if="saveStatus === 'saved'" class="status-icon success"><Check /></el-icon>
          {{ saveStatus === 'saving' ? '保存中...' : '已保存' }}
        </span>
        <el-button @click="saveDraft" type="warning" :icon="Coin">保存草稿</el-button>
        <el-button @click="handleSave" type="primary" :icon="Upload">发布文章</el-button>
      </div>
    </div>
    
    <MdEditor
      v-model="content"
      :previewOnly="previewOnly"
      :toolbarsExclude="previewOnly ? ['save', 'fullscreen', 'pageFullscreen', 'htmlPreview', 'catalog', 'github'] : []"
      :toolbars="customToolbars"
      @onUploadImg="handleUploadImg"
      @onSave="handleSave"
      @onChange="handleContentChange"
      style="height: calc(100% - 56px); min-height: 600px;"
    />
    
    <!-- 字数统计 -->
    <div class="editor-footer">
      <span class="word-count">字数：{{ wordCount }}</span>
      <span class="char-count">字符数：{{ charCount }}</span>
      <span class="last-save">最后保存：{{ lastSaveTime }}</span>
    </div>
    
    <!-- 快速插入菜单 -->
    <div class="quick-insert" v-if="showQuickInsert">
      <el-menu mode="horizontal" @select="handleQuickInsert">
        <el-menu-item index="h2">标题</el-menu-item>
        <el-menu-item index="bold">粗体</el-menu-item>
        <el-menu-item index="code">代码</el-menu-item>
        <el-menu-item index="quote">引用</el-menu-item>
        <el-menu-item index="list">列表</el-menu-item>
        <el-menu-item index="table">表格</el-menu-item>
        <el-menu-item index="link">链接</el-menu-item>
        <el-menu-item index="img">图片</el-menu-item>
      </el-menu>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElNotification } from 'element-plus'
import { MdEditor } from 'md-editor-v3'
import { Edit, View, Expand, Fold, Upload, Loading, Check, Coin } from '@element-plus/icons-vue'
import 'md-editor-v3/lib/style.css'

const props = defineProps<{
  modelValue: string
  previewOnly?: boolean
  title?: string
  articleId?: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'save', value: string): void
  (e: 'saveDraft', value: string): void
}>()

const content = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const previewOnly = ref(props.previewOnly || false)
const isFullscreen = ref(false)
const saveStatus = ref<'saved' | 'saving' | ''>('')
const lastSaveTime = ref('')
const showQuickInsert = ref(false)
const autoSaveTimer = ref<number | null>(null)
const lastContent = ref(props.modelValue)

const customToolbars = [
  'bold', 'del', 'italic', 'quote', 'mark',
  'h1', 'h2', 'h3', 'h4',
  'list', 'ordered-list', 'task',
  'code', 'code-block',
  'table', 'link', 'image',
  'hr', 'br',
  'undo', 'redo',
  'fullscreen', 'pageFullscreen', 'preview', 'htmlPreview',
  'github'
]

const wordCount = computed(() => {
  const text = content.value.replace(/[^\u4e00-\u9fa5a-zA-Z0-9]/g, '')
  return text.length
})

const charCount = computed(() => content.value.length)

const togglePreview = () => {
  previewOnly.value = !previewOnly.value
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
  if (isFullscreen.value) {
    document.documentElement.requestFullscreen?.()
  } else {
    document.exitFullscreen?.()
  }
}

const handleContentChange = () => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
  autoSaveTimer.value = window.setTimeout(() => {
    if (content.value !== lastContent.value) {
      autoSave()
    }
  }, 3000)
}

const autoSave = async () => {
  saveStatus.value = 'saving'
  try {
    emit('saveDraft', content.value)
    lastContent.value = content.value
    saveStatus.value = 'saved'
    lastSaveTime.value = new Date().toLocaleTimeString('zh-CN')
    setTimeout(() => {
      saveStatus.value = ''
    }, 2000)
  } catch (error) {
    saveStatus.value = ''
    console.error('自动保存失败:', error)
  }
}

const saveDraft = async () => {
  saveStatus.value = 'saving'
  try {
    emit('saveDraft', content.value)
    lastContent.value = content.value
    saveStatus.value = 'saved'
    lastSaveTime.value = new Date().toLocaleTimeString('zh-CN')
    ElMessage.success('草稿保存成功')
    setTimeout(() => {
      saveStatus.value = ''
    }, 2000)
  } catch (error) {
    saveStatus.value = ''
    ElMessage.error('保存失败')
  }
}

const handleSave = (value?: string) => {
  const finalContent = value || content.value
  emit('save', finalContent)
}

const handleUploadImg = async (files: File[], callback: (urls: string[]) => void) => {
  const res = await Promise.all(
    files.map(async (file) => {
      const formdata = new FormData()
      formdata.append('file', file)
      formdata.append('type', 'markdown')

      try {
        const response = await fetch('/api/v1/upload', {
          method: 'POST',
          body: formdata,
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        const data = await response.json()

        if (data.status === 200) {
          return data.url
        } else {
          ElMessage.error(data.message || '上传失败')
          return null
        }
      } catch (error) {
        console.error(error)
        ElMessage.error('上传失败')
        return null
      }
    })
  )

  const urls = res.filter((url) => url !== null) as string[]
  if (urls.length > 0) {
    callback(urls)
    ElNotification.success({
      title: '上传成功',
      message: `成功上传 ${urls.length} 张图片`
    })
  }
}

const handleQuickInsert = (index: string) => {
  const insertMap: Record<string, string> = {
    h2: '\n## 标题\n',
    bold: '**粗体文本**',
    code: '`代码`',
    quote: '\n> 引用文本\n',
    list: '\n- 列表项\n',
    table: '\n| 列1 | 列2 |\n| --- | --- |\n| 内容 | 内容 |\n',
    link: '[链接文本](URL)',
    img: '![描述](图片URL)'
  }
  
  if (insertMap[index]) {
    emit('update:modelValue', content.value + insertMap[index])
  }
  showQuickInsert.value = false
}

watch(() => props.previewOnly, (val) => {
  previewOnly.value = val
})

watch(() => props.modelValue, (val) => {
  lastContent.value = val
})

onMounted(() => {
  lastSaveTime.value = new Date().toLocaleTimeString('zh-CN')
})

onUnmounted(() => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
  }
})
</script>

<style scoped>
.markdown-editor-container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0,0,0,0.08);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: #fafafa;
  border-bottom: 1px solid #e8e8e8;
  flex-shrink: 0;
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.save-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #67c23a;
}

.status-icon.success {
  color: #67c23a;
}

.status-icon {
  color: #409eff;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 20px;
  padding: 10px 20px;
  background: #fafafa;
  border-top: 1px solid #e8e8e8;
  font-size: 13px;
  color: #999;
  flex-shrink: 0;
}

.word-count,
.char-count,
.last-save {
  display: flex;
  align-items: center;
}

.quick-insert {
  position: absolute;
  top: 60px;
  left: 20px;
  z-index: 100;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  padding: 8px;
}

@media (max-width: 768px) {
  .editor-header {
    padding: 10px 12px;
  }
  
  .editor-footer {
    padding: 8px 12px;
    gap: 12px;
    font-size: 12px;
  }
  
  .header-right {
    gap: 6px;
  }
}
</style>