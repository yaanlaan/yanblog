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
          >
            <el-option 
              v-for="category in categories" 
              :key="category.id" 
              :label="category.name" 
              :value="category.id" 
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="摘要">
          <el-input 
            v-model="publishData.desc" 
            type="textarea" 
            :rows="4" 
            placeholder="请输入文章摘要"
          />
        </el-form-item>
        
        <el-form-item label="封面图">
          <el-upload
            class="avatar-uploader"
            action="/api/v1/upload"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
            :before-upload="beforeAvatarUpload"
            :headers="uploadHeaders"
          >
            <img v-if="publishData.img" :src="publishData.img" class="avatar" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
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
  }
  categories: {id: number, name: string}[]
  isEdit: boolean
  submitLoading: boolean
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: {categoryId: number | undefined, desc: string, img: string}): void
  (e: 'submit'): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 发布数据
const publishData = reactive({
  categoryId: props.modelValue.categoryId,
  desc: props.modelValue.desc,
  img: props.modelValue.img
})

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  publishData.categoryId = newVal.categoryId
  publishData.desc = newVal.desc
  publishData.img = newVal.img
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

// 处理上传成功
const handleAvatarSuccess: UploadProps['onSuccess'] = (
  response,
  uploadFile
) => {
  // 从响应中获取图片URL
  if (response && response.data && response.data.url) {
    publishData.img = response.data.url
  } else {
    // 如果响应格式不正确，使用本地URL
    publishData.img = URL.createObjectURL(uploadFile.raw!)
  }
  emit('update:modelValue', {...publishData})
}

// 上传前检查
const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
  if (rawFile.type !== 'image/jpeg' && rawFile.type !== 'image/png') {
    ElMessage.error('头像图片必须是 JPG 或 PNG 格式!')
    return false
  } else if (rawFile.size / 1024 / 1024 > 2) {
    ElMessage.error('图片大小不能超过 2MB!')
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
      emit('update:modelValue', {...publishData})
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
</style>