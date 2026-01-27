<template>
  <div class="about-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>关于页面内容编辑</span>
          <el-button type="primary" @click="saveContent" :loading="saving">保存更改</el-button>
        </div>
      </template>

      <div class="editor-container" v-loading="loading">
        <MdEditor v-model="content" style="height: 700px" @onUploadImg="handleUploadImg" @onSave="saveContent" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from '@/utils/request'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

const content = ref('')
const loading = ref(false)
const saving = ref(false)

const fetchContent = async () => {
  loading.value = true
  try {
    const res = await axios.get('v1/about')
    if (res.status === 200) {
      content.value = res.data
    } else {
      ElMessage.error(res.message || '获取内容失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('获取内容失败')
  } finally {
    loading.value = false
  }
}

const saveContent = async () => {
  saving.value = true
  try {
    const res = await axios.put('v1/about', { content: content.value })
    if (res.status === 200) {
      ElMessage.success('保存成功')
    } else {
      ElMessage.error(res.message || '保存失败')
    }
  } catch (error) {
    console.error(error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleUploadImg = async (files: File[], callback: (urls: string[]) => void) => {
  const res = await Promise.all(
    files.map(async (file) => {
      const formdata = new FormData()
      formdata.append('file', file)
      const res = await axios.post('/upload', formdata, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      return res.url
    })
  )
  callback(res)
}

onMounted(() => {
  fetchContent()
})
</script>

<style scoped>
.about-editor {
  height: 100%;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.editor-container {
  min-height: 400px;
}
</style>
