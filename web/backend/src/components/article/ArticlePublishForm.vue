<template>
  <el-card class="publish-card">
    <template #header>
      <div class="card-header">
        <span>发布</span>
      </div>
    </template>
    
    <div class="publish-content">
      <el-form 
        ref="formRef" 
        :model="publishData" 
        :rules="publishRules"
      >
        <el-form-item label="分类" prop="categoryId">
          <el-select 
            v-model="publishData.categoryId" 
            placeholder="请选择分类" 
            clearable
            filterable
            style="width: 100%"
            @change="handleFormChange"
          >
            <el-option 
              v-for="category in categories" 
              :key="category.id" 
              :label="category.name" 
              :value="category.id" 
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="标签">
          <el-input 
            v-model="publishData.tags" 
            placeholder="请输入标签，多个标签用逗号分隔" 
            @input="handleFormChange"
          />
        </el-form-item>

        <el-form-item label="摘要">
          <el-input 
            v-model="publishData.desc" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入文章摘要"
            @input="handleFormChange"
          />
        </el-form-item>
        
        <el-form-item label="置顶等级">
          <el-select 
            v-model="publishData.top" 
            placeholder="请选择置顶等级"
            clearable
            style="width: 100%"
            @change="handleFormChange"
          >
            <el-option label="不置顶" :value="0" />
            <el-option label="等级1 (最高)" :value="1" />
            <el-option label="等级2" :value="2" />
            <el-option label="等级3" :value="3" />
            <el-option label="等级4" :value="4" />
            <el-option label="等级5" :value="5" />
            <el-option label="等级6 (最低)" :value="6" />
          </el-select>
          <div class="top-tip">数字越小等级越高，0表示不置顶</div>
        </el-form-item>
        
        <el-form-item label="封面图">
          <el-input 
            v-model="publishData.img" 
            placeholder="请输入图片URL或上传图片" 
            @input="handleFormChange"
          />
          <el-upload
            class="avatar-uploader"
            action="/api/v1/upload"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :on-error="handleAvatarError"
            :before-upload="beforeAvatarUpload"
            :headers="uploadHeaders"
            :data="{ type: 'cover' }"
          >
            <img v-if="publishData.img" :src="publishData.img" class="avatar" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">点击上传文章封面图片或直接输入图片URL</div>
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleSubmit"
            :loading="submitLoading"
            style="width: 100%"
          >
            {{ isEdit ? '更新文章' : '发布文章' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance, FormRules, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 定义组件属性
const props = defineProps<{
  modelValue: {
    categoryId: number | undefined
    desc: string
    img: string
    top: number
    tags: string
  }
  categories: {id: number, name: string}[]
  isEdit: boolean
  submitLoading: boolean
  title?: string
  articleId?: number
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: {categoryId: number | undefined, desc: string, img: string, top: number, tags: string}): void
  (e: 'submit'): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 发布数据
const publishData = reactive({
  categoryId: props.modelValue.categoryId,
  desc: props.modelValue.desc,
  img: props.modelValue.img,
  top: props.modelValue.top || 0,
  tags: props.modelValue.tags || ''
})

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  publishData.categoryId = newVal.categoryId
  publishData.desc = newVal.desc
  publishData.img = newVal.img
  publishData.top = newVal.top || 0
  publishData.tags = newVal.tags || ''
}, { deep: true })

// 表单验证规则
const publishRules = reactive<FormRules>({
  categoryId: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ]
})

// 上传请求头
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 上传额外参数
import { computed } from 'vue'
const uploadData = computed(() => {
  return {
    type: 'article',
    key: props.title || 'default',
    id: props.articleId || 0
  }
})

// 处理表单变化
const handleFormChange = () => {
  emit('update:modelValue', {
    categoryId: publishData.categoryId,
    desc: publishData.desc,
    img: publishData.img,
    top: publishData.top,
    tags: publishData.tags
  })
}

// 处理上传成功
const handleAvatarSuccess: UploadProps['onSuccess'] = (
  response,
  uploadFile
) => {
  console.log('Upload response:', response); // 调试信息
  
  // 检查后端返回的状态码
  if (response && response.status !== 200) {
    ElMessage.error(response.message || '上传失败')
    return
  }

  // 从响应中获取图片URL
  if (response && response.url) {
    // 后端直接返回url字段
    publishData.img = response.url
  } else {
    // 如果响应格式不正确，显示错误信息
    ElMessage.error('上传失败：服务器响应格式不正确')
    return
  }
  handleFormChange()
}

// 处理上传错误
const handleAvatarError: UploadProps['onError'] = (
  error
) => {
  console.error('Upload error:', error); // 调试信息
  ElMessage.error('上传失败：' + (error.message || '未知错误'))
}

// 上传前检查
const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
  if (rawFile.type !== 'image/jpeg' && rawFile.type !== 'image/png') {
    ElMessage.error('头像图片必须是 JPG 或 PNG 格式!')
    return false
  } else if (rawFile.size / 1024 / 1024 > 10) {
    ElMessage.error('图片大小不能超过 10MB!')
    return false
  }
  return true
}

// 处理提交
const handleSubmit = () => {
  if (!formRef.value) return
  formRef.value.validate((valid) => {
    if (valid) {
      emit('submit')
      handleFormChange()
    }
  })
}
</script>

<style scoped>
.publish-card {
  height: 100%;
}

.card-header {
  font-weight: bold;
}

.avatar-uploader .avatar {
  width: 178px;
  height: 178px;
  display: block;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}

.top-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.upload-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>