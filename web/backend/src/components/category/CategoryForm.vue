<template>
  <el-dialog 
    v-model="visible" 
    :title="title" 
    width="500px"
    @close="handleClose"
  >
    <el-form 
      ref="formRef" 
      :model="formData" 
      :rules="formRules" 
      label-width="80px"
    >
      <el-form-item label="分类名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入分类名称" />
      </el-form-item>
      <el-form-item label="封面图片" prop="img">
        <el-input v-model="formData.img" placeholder="请输入图片URL或上传图片" />
        <el-upload
          class="avatar-uploader"
          action="/api/v1/upload"
          :show-file-list="false"
          :on-success="handleAvatarSuccess"
          :on-error="handleAvatarError"
          :before-upload="beforeAvatarUpload"
          :headers="uploadHeaders"
        >
          <img v-if="formData.img" :src="formData.img" class="avatar" />
          <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
        </el-upload>
        <div class="upload-tip">点击上传分类封面图片或直接输入图片URL</div>
      </el-form-item>
      <el-form-item label="置顶排序" prop="top">
        <el-input-number 
          v-model="formData.top" 
          :min="0" 
          :max="999" 
          controls-position="right" 
          placeholder="请输入置顶排序值"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance, FormRules, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 定义组件属性
const props = defineProps<{
  modelValue: boolean      // 控制对话框可见性
  title: string           // 对话框标题
  category?: {            // 分类数据，可选属性（新建时可能不需要）
    id?: number | null    // id可空（新建时可能为null）
    name: string
    img: string
    top: number
  }
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', category: {id: number, name: string, img: string, top: number}): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: props.category.id,
  name: props.category.name,
  img: props.category.img,
  top: props.category.top
})

// 上传请求头
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}

// 表单验证规则
const formRules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 20, message: '分类名称长度为2-20位', trigger: 'blur' }
  ],
  img: [
    { required: false, message: '请上传封面图片', trigger: 'blur' }
  ],
  top: [
    { required: false, message: '请输入置顶排序值', trigger: 'blur' }
  ]
})

// 对话框可见性
const visible = ref(props.modelValue)

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
})

watch(() => props.category, (newVal) => {
  formData.id = newVal.id
  formData.name = newVal.name
  formData.img = newVal.img
  formData.top = newVal.top
}, { deep: true })

// 处理上传成功
const handleAvatarSuccess: UploadProps['onSuccess'] = (
  response
) => {
  // 从响应中获取图片URL
  if (response && response.url) {
    formData.img = response.url
  } else {
    ElMessage.error('上传失败：服务器响应格式不正确')
    return
  }
}

// 处理上传错误
const handleAvatarError: UploadProps['onError'] = (
  error
) => {
  ElMessage.error('上传失败：' + (error.message || '未知错误'))
}

// 上传前检查
const beforeAvatarUpload: UploadProps['beforeUpload'] = (rawFile) => {
  if (rawFile.type !== 'image/jpeg' && rawFile.type !== 'image/png') {
    ElMessage.error('封面图片必须是 JPG 或 PNG 格式!')
    return false
  } else if (rawFile.size / 1024 / 1024 > 10) {
    ElMessage.error('图片大小不能超过 10MB!')
    return false
  }
  return true
}

// 处理关闭
const handleClose = () => {
  emit('update:modelValue', false)
  formRef.value?.resetFields()
}

// 处理取消
const handleCancel = () => {
  handleClose()
}

// 处理提交
const handleSubmit = () => {
  if (!formRef.value) return
  formRef.value.validate((valid) => {
    if (valid) {
      // 添加日志以便调试
      console.log('Submitting category form data:', {...formData})
      // 确保正确传递分类数据，包括id、名称、图片和置顶排序值
      emit('submit', {
        id: formData.id as number,
        name: formData.name,
        img: formData.img,
        top: formData.top
      })
      // 提交后重置表单
      formRef.value.resetFields()
      // 关闭对话框
      handleClose()
    } else {
      console.log('Form validation failed')
      ElMessage.error('表单验证失败，请检查输入内容')
      return false
    }
  })
}
</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.avatar-uploader .avatar {
  width: 178px;
  height: 178px;
  display: block;
}

.avatar-uploader .el-upload {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
}

.upload-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}
</style>