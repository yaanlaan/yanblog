<template>
  <div class="user-actions">
    <el-button size="small" @click="handleEdit">编辑</el-button>
    <el-button 
      size="small" 
      type="danger" 
      @click="handleDelete"
      :disabled="user.role === 2"
    >
      删除
    </el-button>
  </div>
</template>

<script setup lang="ts">
// 定义组件属性
const props = defineProps<{
  user: any
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'edit', user: any): void
  (e: 'delete', user: any): void
}>()

// 处理编辑
const handleEdit = () => {
  emit('edit', props.user)
}

// 处理删除
const handleDelete = () => {
  emit('delete', props.user)
}
</script>

<style scoped>
.user-actions {
  display: flex;
  gap: 8px;
}
</style>