<template>
  <div class="article-editor">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑文章' : '新增文章' }}</span>
          <div>
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
            
            <el-form-item label="内容" prop="content">
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
                ref="markdownEditorRef"
                class="content-editor"
              />
            </el-form-item>
          </el-form>
        </el-col>
        
        <el-col :span="6">
          <ArticlePublishForm
            v-model="publishForm"
            :categories="categories"
            :is-edit="isEdit"
            :submit-loading="submitLoading"
            @submit="submitArticle"
            @update:modelValue="handlePublishFormUpdate"
          />
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
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
  top: 0
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
    { required: true, message: '请输入文章内容', trigger: 'blur' }
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
const handlePublishFormUpdate = (value: {categoryId: number | undefined, desc: string, img: string, top: number}) => {
  publishForm.categoryId = value.categoryId
  publishForm.desc = value.desc
  publishForm.img = value.img
  publishForm.top = value.top
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
    
    if (isEdit.value) {
      // 编辑文章
      await articleApi.updateArticle(articleId.value, {
        title: articleForm.title,
        cid: publishForm.categoryId!,
        desc: publishForm.desc,
        content: articleForm.content,
        img: publishForm.img,
        top: publishForm.top
      })
      ElMessage.success('文章更新成功')
    } else {
      // 新增文章
      await articleApi.createArticle({
        title: articleForm.title,
        cid: publishForm.categoryId!,
        desc: publishForm.desc,
        content: articleForm.content,
        img: publishForm.img,
        top: publishForm.top
      })
      ElMessage.success('文章发布成功')
    }
    
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
    publishForm.img = article.img
    publishForm.top = article.top || 0
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
</style>