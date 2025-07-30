<template>
  <el-form :model="searchData" label-width="80px" class="search-form">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-form-item label="分类名称">
          <el-input v-model="searchData.name" placeholder="请输入分类名称" />
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
    name: string
  }
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: {name: string}): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

// 搜索数据
const searchData = reactive({
  name: props.modelValue.name
})

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  searchData.name = newVal.name
}, { deep: true })

// 处理搜索
const handleSearch = () => {
  emit('update:modelValue', {...searchData})
  emit('search')
}

// 处理重置
const handleReset = () => {
  searchData.name = ''
  emit('update:modelValue', {...searchData})
  emit('reset')
}
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>