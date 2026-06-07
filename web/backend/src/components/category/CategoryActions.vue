<template>
  <div class="category-actions">
    <el-button size="small" @click="handleEdit">编辑</el-button>
    <el-button 
      size="small" 
      type="danger" 
      @click="handleDelete"
    >
      删除
    </el-button>
  </div>
</template>

<script setup lang="ts">
// 定义组件属性
const props = defineProps<{
  category: any
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'edit', category: any): void
  (e: 'delete', category: any): void
}>()

// 处理编辑
const handleEdit = () => {
  emit('edit', props.category)
}

// 处理删除
const handleDelete = () => {
  emit('delete', props.category)
}
</script>

<style scoped>
.category-actions {
  display: flex;
  gap: 8px;
}
</style>