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
      <el-form-item label="用户名" prop="username">
        <el-input v-model="formData.username" placeholder="请输入用户名" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input 
          v-model="formData.password" 
          type="password" 
          :placeholder="isAdd ? '请输入密码' : '不修改请留空'" 
          show-password
        />
      </el-form-item>
      <el-form-item label="角色" prop="role">
        <el-select v-model="formData.role" placeholder="请选择角色" :disabled="!canEditRole">
          <el-option label="超级管理员" :value="1" disabled />
          <el-option label="管理员" :value="2" />
          <el-option label="普通用户" :value="3" />
        </el-select>
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
import { ref, reactive, watch, computed } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

// 定义组件属性
const props = defineProps<{
  modelValue: boolean
  title: string
  isAdd: boolean
  user: {
    id: number
    username: string
    password: string
    role: number
  }
}>()

// 获取当前用户信息（假设存储在localStorage中，实际项目中建议使用Pinia）
const currentUserRole = computed(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      return user.role || 3 // 默认为普通用户
    } catch (e) {
      return 3
    }
  }
  return 3
})

// 是否可以编辑角色
const canEditRole = computed(() => {
  // 1. 只有超级管理员可以修改角色
  if (currentUserRole.value !== 1) return false
  
  // 2. 如果正在编辑的是超级管理员（即自己），则不允许修改角色
  // 注意：这里使用 props.user.role 来判断原始角色
  if (!props.isAdd && props.user.role === 1) return false
  
  return true
})

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', user: {id: number, username: string, password: string, role: number}): void
}>()

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive({
  id: props.user.id,
  username: props.user.username,
  password: props.user.password,
  role: props.user.role
})

// 表单验证规则
const formRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 4, max: 12, message: '用户名长度为4-12位', trigger: 'blur' }
  ],
  password: [
    { required: false, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为6-20位', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
})

// 对话框可见性
const visible = ref(props.modelValue)

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
})

watch(() => props.user, (newVal) => {
  formData.id = newVal.id
  formData.username = newVal.username
  formData.password = newVal.password
  formData.role = newVal.role
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