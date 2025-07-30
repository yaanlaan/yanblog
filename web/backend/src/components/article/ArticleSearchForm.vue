<template>
  <el-form :model="searchData" label-width="80px" class="search-form">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-form-item label="文章标题">
          <el-input v-model="searchData.title" placeholder="请输入文章标题" />
        </el-form-item>
      </el-col>
      <el-col :span="6">
        <el-form-item label="分类">
          <el-select v-model="searchData.categoryId" placeholder="请选择分类" clearable>
            <el-option 
              v-for="category in categories" 
              :key="category.id" 
              :label="category.name" 
              :value="category.id" 
            />
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
    title: string
    categoryId: number | undefined
  }
  categories: {id: number, name: string}[]
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: {title: string, categoryId: number | undefined}): void
  (e: 'search'): void
  (e: 'reset'): void
}>()

// 搜索数据
const searchData = reactive({
  title: props.modelValue.title,
  categoryId: props.modelValue.categoryId
})

// 监听属性变化
watch(() => props.modelValue, (newVal) => {
  searchData.title = newVal.title
  searchData.categoryId = newVal.categoryId
}, { deep: true })

// 处理搜索
const handleSearch = () => {
  emit('update:modelValue', {...searchData})
  emit('search')
}

// 处理重置
const handleReset = () => {
  searchData.title = ''
  searchData.categoryId = undefined
  emit('update:modelValue', {...searchData})
  emit('reset')
}
</script>

<style scoped>
.search-form {
  margin-bottom: 20px;
}
</style>