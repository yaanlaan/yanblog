<template>
  <div class="article-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑文章' : '新增文章' }}</span>
          <div>
            <template v-if="articleType === 1">
              <el-button @click="toggleMarkdownEditor" v-if="!showMarkdownEditor">
                Markdown编辑
              </el-button>
              <el-button @click="toggleMarkdownEditor" v-else>
                普通编辑
              </el-button>
              <el-button @click="togglePreviewOnly" v-if="showMarkdownEditor && !previewOnly">
                仅预览
              </el-button>
              <el-button @click="togglePreviewOnly" v-if="showMarkdownEditor && previewOnly">
                双栏显示
              </el-button>
            </template>
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="18">
          <el-form 
            ref="articleFormRef" 
            :model="articleForm" 
            :rules="articleFormRules" 
            label-width="80px"
          >
            <el-form-item label="标题" prop="title">
              <el-input 
                v-model="articleForm.title" 
                placeholder="请输入文章标题" 
                size="large"
              />
            </el-form-item>

            <el-form-item label="类型">
              <el-radio-group v-model="articleType">
                <el-radio :label="1" :value="1">文本/Markdown</el-radio>
                <el-radio :label="2" :value="2">PDF上传</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="内容" prop="content" v-if="articleType === 1">
              <!-- 普通文本编辑模式 -->
              <el-input 
                v-model="articleForm.content" 
                type="textarea" 
                :rows="20" 
                placeholder="请输入文章内容"
                v-if="!showMarkdownEditor"
                ref="textareaRef"
                @scroll="handleTextareaScroll"
                class="content-textarea"
              />
              
              <!-- Markdown编辑模式 -->
              <ArticleEditor 
                v-else
                v-model="articleForm.content"
                :preview-only="previewOnly"
                :title="articleForm.title"
                :article-id="articleId"
                @save="submitArticle"
                ref="markdownEditorRef"
                class="content-editor"
              />
            </el-form-item>

            <el-form-item label="PDF文件" required v-if="articleType === 2">
              <div class="pdf-edit-container">
                <div class="upload-bar">
                  <el-upload
                    class="pdf-uploader-bar"
                    drag
                    action="/api/v1/upload"
                    :data="{ type: 'pdf' }" 
                    :headers="uploadHeaders"
                    :on-success="handlePdfSuccess"
                    :on-error="handlePdfError"
                    :before-upload="beforePdfUpload"
                    :show-file-list="false"
                    accept=".pdf"
                  >
                    <div class="upload-bar-content">
                       <el-icon class="upload-icon"><Document /></el-icon>
                       <span v-if="pdfUrl" class="file-text">{{ pdfUrl }}</span>
                       <span class="upload-hint">{{ pdfUrl ? ' (拖拽或点击替换)' : '拖拽 PDF 文件到此处或点击上传' }}</span>
                    </div>
                  </el-upload>
                </div>
                
                <div v-if="pdfUrl" class="preview-section-full">
                   <div class="preview-frame">
                     <iframe 
                        :src="pdfUrl" 
                        width="100%" 
                        height="100%" 
                        frameborder="0"
                     ></iframe>
                   </div>
                </div>
              </div>
            </el-form-item>
          </el-form>
        </el-col>
        
        <el-col :span="6">
          <ArticlePublishForm
            v-model="publishForm"
            :categories="categories"
            :is-edit="isEdit"
            :submit-loading="submitLoading"
            :title="articleForm.title"
            :article-id="articleId"
            @submit="submitArticle"
            @update:modelValue="handlePublishFormUpdate"
          />
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>


<script setup lang="ts">
import { ref, reactive, onMounted, computed, onActivated } from 'vue'

defineOptions({
  name: 'ArticleEditor'
})

import type { FormInstance, FormRules, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'
import { UploadFilled, Document } from '@element-plus/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { articleApi, categoryApi } from '@/services/api'
import ArticleEditor from '@/components/article/ArticleEditor.vue'
import ArticlePublishForm from '@/components/article/ArticlePublishForm.vue'

// 路由和参数
const route = useRoute()
const router = useRouter()

// 表单引用
const articleFormRef = ref<FormInstance>()
const markdownEditorRef = ref()

// 是否为编辑模式
const isEdit = ref(false)
const articleId = ref(0)

// 是否显示Markdown编辑器
const showMarkdownEditor = ref(false)

// 是否仅显示预览
const previewOnly = ref(false)

// 提交状态
const submitLoading = ref(false)

// 文章类型
const articleType = ref(1) // 1: Markdown, 2: PDF
const pdfUrl = ref('')

// 上传Header
const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${localStorage.getItem('token')}`
  }
})

// 文章表单
const articleForm = reactive({
  title: '',
  content: ''
})

// 发布表单
const publishForm = reactive({
  categoryId: undefined as number | undefined,
  desc: '',
  img: '',
  top: 0,
  tags: ''
})

// 分类列表
const categories = ref<{id: number, name: string}[]>([])

// 表单验证规则
const articleFormRules = reactive<FormRules>({
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度为2-100位', trigger: 'blur' }
  ],
  content: [
    { 
      validator: (rule: any, value: any, callback: any) => {
        if (articleType.value === 1 && !value) {
          callback(new Error('请输入文章内容'))
        } else if (articleType.value === 2 && !pdfUrl.value) {
          callback(new Error('请上传PDF文件'))
        } else {
          callback()
        }
      }, 
      trigger: 'blur' 
    }
  ]
})

// 切换Markdown编辑器
const toggleMarkdownEditor = () => {
  showMarkdownEditor.value = !showMarkdownEditor.value
  previewOnly.value = false
}

// 切换仅预览模式
const togglePreviewOnly = () => {
  previewOnly.value = !previewOnly.value
}

// 处理文本区域滚动
const handleTextareaScroll = () => {
  // 普通文本区域滚动处理逻辑
}

// PDF上传成功
const handlePdfSuccess: UploadProps['onSuccess'] = (response) => {
  if (response.status === 200) {
    pdfUrl.value = response.url
    articleForm.content = 'PDF文章: ' + response.url // 填充一点内容通过必填校验，实际使用pdf_url
    ElMessage.success('PDF上传成功')
  } else {
    ElMessage.error(response.message || 'PDF上传失败')
  }
}

// PDF上传失败
const handlePdfError: UploadProps['onError'] = (error) => {
  ElMessage.error('PDF上传失败')
  console.error(error)
}

// PDF上传前检查
const beforePdfUpload: UploadProps['beforeUpload'] = (rawFile) => {
  if (rawFile.type !== 'application/pdf') {
    ElMessage.error('必须上传PDF文件!')
    return false
  }
  if (rawFile.size / 1024 / 1024 > 50) { // 50MB限制
    ElMessage.error('PDF文件大小不能超过50MB!')
    return false
  }
  return true
}

// 获取分类列表
const getCategoryList = async () => {
  try {
    const response = await categoryApi.getCategories({
      pagesize: -1,
      pagenum: -1
    })
    
    // 解析后端返回的数据
    const { data } = response.data
    categories.value = data.map((item: any) => ({
      id: item.ID !== undefined ? parseInt(item.ID, 10) : parseInt(item.id, 10),
      name: item.name
    }))
  } catch (error) {
    ElMessage.error('获取分类列表失败')
    console.error(error)
  }
}

// 处理发布表单更新
const handlePublishFormUpdate = (value: {categoryId: number | undefined, desc: string, img: string, top: number, tags: string}) => {
  publishForm.categoryId = value.categoryId
  publishForm.desc = value.desc
  publishForm.img = value.img
  publishForm.top = value.top
  publishForm.tags = value.tags
}

// 提交文章
const submitArticle = async () => {
  // 先验证文章内容表单
  if (!articleFormRef.value) return
  
  try {
    // 验证文章内容
    await articleFormRef.value.validate()
    
    // 提交数据
    submitLoading.value = true
    
    // PDF文章内容处理
    const contentToSend = articleType.value === 2 ? (articleForm.content || 'PDF Article') : articleForm.content
    
    let res;
    const articleData = {
      title: articleForm.title,
      cid: publishForm.categoryId!,
      desc: publishForm.desc,
      content: contentToSend,
      img: publishForm.img,
      top: publishForm.top,
      tags: publishForm.tags,
      type: articleType.value,
      pdf_url: pdfUrl.value
    }

    if (isEdit.value) {
      // 编辑文章
      res = await articleApi.updateArticle(articleId.value, articleData)
    } else {
      // 新增文章
      res = await articleApi.createArticle(articleData)
    }

    // 检查后端返回的状态码
    if (res.data.status !== 200) {
      ElMessage.error(res.data.message || (isEdit.value ? '文章更新失败' : '文章发布失败'))
      return
    }
    
    ElMessage.success(isEdit.value ? '文章更新成功' : '文章发布成功')
    
    // 返回文章列表
    goBack()
  } catch (error: any) {
    console.error('提交文章失败:', error)
    ElMessage.error(isEdit.value ? '文章更新失败' : '文章发布失败')
  } finally {
    submitLoading.value = false
  }
}

// 返回文章列表
const goBack = () => {
  router.push('/article/list')
}

// 获取文章详情（编辑模式）
const getArticleDetail = async (id: number) => {
  try {
    const response = await articleApi.getArticle(id)
    const article = response.data.data
    
    // 填充表单数据
    articleForm.title = article.title
    articleForm.content = article.content
    publishForm.categoryId = parseInt(article.cid, 10)
    publishForm.desc = article.desc
    publishForm.tags = article.tags || ''
    publishForm.img = article.img
    publishForm.top = article.top || 0
    articleType.value = article.type || 1
    pdfUrl.value = article.pdf_url || ''
  } catch (error) {
    ElMessage.error('获取文章详情失败')
    console.error(error)
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getCategoryList()
  
  // 检查是否为编辑模式
  if (route.name === 'ArticleEdit' && route.params.id) {
    isEdit.value = true
    articleId.value = parseInt(route.params.id as string, 10)
    getArticleDetail(articleId.value)
  }
})

// 组件激活时只更新分类列表，保留已输入的表单内容
onActivated(() => {
  getCategoryList() 
})
</script>

<style scoped>
.article-editor {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.content-textarea,
.content-editor {
  width: 100%;
}

.pdf-edit-container {
  display: flex;
  flex-direction: column;
  gap: 15px;
  width: 100%;
}

.upload-bar {
  width: 100%;
}

.pdf-uploader-bar :deep(.el-upload) {
  width: 100%;
}

.pdf-uploader-bar :deep(.el-upload-dragger) {
  width: 100% !important;
  height: 40px;
  padding: 0 15px;
  border: 1px solid #dcdfe6; /* 类似 Input 边框 */
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  background-color: #fff;
  transition: border-color 0.2s;
}

.pdf-uploader-bar :deep(.el-upload-dragger:hover) {
  border-color: #409eff;
}

.pdf-uploader-bar :deep(.el-upload-dragger.is-dragover) {
  border-color: #409eff;
  background-color: rgba(64, 158, 255, 0.1);
}

.upload-bar-content {
  display: flex;
  align-items: center;
  width: 100%;
  overflow: hidden;
}

.upload-icon {
  margin-right: 8px;
  font-size: 16px;
  color: #909399;
}

.file-text {
  color: #606266;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 10px;
  max-width: 70%;
}

.upload-hint {
  color: #c0c4cc;
  font-size: 13px;
  /* 如果没有文件，提示文字颜色可以深一点 */
  white-space: nowrap;
}

.upload-bar:not(:has(.file-text)) .upload-hint {
    color: #909399;
}

.preview-section-full {
  width: 100%;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #f8f9fa;
  padding: 15px;
}

.preview-frame {
  height: 700px;
  background: white;
  border: 1px solid #e4e7ed;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

</style>