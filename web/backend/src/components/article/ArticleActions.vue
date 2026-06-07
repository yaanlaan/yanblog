<template>
  <div class="article-actions">
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
  article: any
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'edit', article: any): void
  (e: 'delete', article: any): void
}>()

// 处理编辑
const handleEdit = () => {
  emit('edit', props.article)
}

// 处理删除
const handleDelete = () => {
  emit('delete', props.article)
}
</script>

<style scoped>
.article-actions {
  display: flex;
  gap: 8px;
}
</style>