<template>
  <el-form :model="searchData" label-width="80px" class="search-form">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-form-item label="用户名">
          <el-input v-model="searchData.username" placeholder="请输入用户名" />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item label="角色">
          <el-select v-model="searchData.role" placeholder="请选择角色" clearable>
            <el-option label="管理员" :value="2" />
            <el-option label="普通用户" :value="3" />
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-col>
    </el-row>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'

// 定义组件属性
const props = defineProps<{
  modelValue: {
    username: string
    role: number | undefined
  }
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: {username: string, role: number | undefined}): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

// 搜索数据
const searchData = reactive({
  username: props.modelValue.username,
  role: props.modelValue.role
})

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  searchData.username = newVal.username
  searchData.role = newVal.role
}, { deep: true })

// 处理搜索
const handleSearch = () => {
  emit('update:modelValue', {...searchData})
  emit('search')
}

// 处理重置
const handleReset = () => {
  searchData.username = ''
  searchData.role = undefined
  emit('update:modelValue', {...searchData})
  emit('reset')
}
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>