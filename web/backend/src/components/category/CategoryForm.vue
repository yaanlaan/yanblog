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
import type { FormInstance, FormRules } from 'element-plus'

// 定义组件属性
const props = defineProps<{
  modelValue: boolean
  title: string
  category: {
    id: number
    name: string
  }
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', category: {id: number, name: string}): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: props.category.id,
  name: props.category.name
})

// 表单验证规则
const formRules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 20, message: '分类名称长度为2-20位', trigger: 'blur' }
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
}, { deep: true })

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
      emit('submit', {...formData})
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
</style>