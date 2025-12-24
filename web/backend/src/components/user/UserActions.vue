<template>
  <div class="user-actions">
    <el-button 
      size="small" 
      @click="handleEdit"
      v-if="canEdit"
    >
      编辑
    </el-button>
    <el-button 
      size="small" 
      type="danger" 
      @click="handleDelete"
      v-if="canDelete"
    >
      删除
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// 定义组件属性
const props = defineProps<{
  user: any
}>()

// 获取当前用户信息
const currentUser = computed(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch (e) {
      return { role: 3, username: '' }
    }
  }
  return { role: 3, username: '' }
})

// 是否可以编辑
const canEdit = computed(() => {
  const myRole = currentUser.value.role
  const targetRole = props.user.role
  const myUsername = currentUser.value.username
  const targetUsername = props.user.username

  // 1. 超级管理员(1)可以编辑任何人
  if (myRole === 1) return true

  // 2. 管理员(2)
  if (myRole === 2) {
    // 可以编辑自己
    if (myUsername === targetUsername) return true
    // 可以编辑普通用户(3)
    if (targetRole === 3) return true
    // 不能编辑超级管理员(1)和其他管理员(2)
    return false
  }

  // 3. 普通用户(3)只能编辑自己
  if (myRole === 3) {
    return myUsername === targetUsername
  }

  return false
})

// 是否可以删除
const canDelete = computed(() => {
  const myRole = currentUser.value.role
  const targetRole = props.user.role
  
  // 1. 超级管理员(1)可以删除任何人（除了自己，通常前端做限制）
  if (myRole === 1) {
      // 简单起见，不让删除自己，防止误删导致无管理员
      return props.user.username !== currentUser.value.username
  }

  // 2. 管理员(2)只能删除普通用户(3)
  if (myRole === 2) {
    return targetRole === 3
  }

  // 3. 普通用户(3)不能删除任何人
  return false
})

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