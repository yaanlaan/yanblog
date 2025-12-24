<template>
  <div class="markdown-editor-container">
    <MdEditor
      v-model="content"
      :previewOnly="previewOnly"
      :toolbarsExclude="previewOnly ? ['save', 'fullscreen', 'pageFullscreen', 'htmlPreview', 'catalog', 'github'] : []"
      @onUploadImg="handleUploadImg"
      @onSave="handleSave"
      style="height: 100%; min-height: 600px;"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

// 定义组件属性
const props = defineProps<{
  modelValue: string
  previewOnly?: boolean
  title?: string
  articleId?: number
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'save', value: string): void
}>()

// 内容值
const content = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 保存事件
const handleSave = (value: string) => {
  emit('save', value)
}

// 图片上传
const handleUploadImg = async (files: File[], callback: (urls: string[]) => void) => {
  const res = await Promise.all(
    files.map(async (file) => {
      const formdata = new FormData()
      formdata.append('file', file)
      formdata.append('type', 'article')
      formdata.append('key', props.title || 'default')
      formdata.append('id', (props.articleId || 0).toString())

      try {
        const response = await axios({
          url: '/api/v1/upload',
          method: 'post',
          data: formdata,
          headers: { 
            'Content-Type': 'multipart/form-data',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })

        if (response.data.status === 200) {
          return response.data.url
        } else {
          ElMessage.error(response.data.message || '上传失败')
          return null
        }
      } catch (error) {
        console.error(error)
        ElMessage.error('上传失败')
        return null
      }
    })
  )

  // 过滤掉失败的上传
  const urls = res.filter((url) => url !== null) as string[]
  if (urls.length > 0) {
    callback(urls)
    ElMessage.success(`成功上传 ${urls.length} 张图片`)
  }
}
</script>

<style scoped>
.markdown-editor-container {
  height: 100%;
  width: 100%;
}
</style>