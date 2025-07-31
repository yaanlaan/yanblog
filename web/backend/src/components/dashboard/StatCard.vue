<template>
  <el-card class="stat-card">
    <div class="stat-content">
      <div class="stat-icon" :class="type">
        <el-icon :size="24" color="white">
          <component :is="icon" />
        </el-icon>
      </div>
      <div class="stat-info">
        <div class="stat-number">{{ number }}</div>
        <div class="stat-label">{{ label }}</div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { User, Document, Folder, View } from '@element-plus/icons-vue'

// 定义组件属性
interface Props {
  type: 'users' | 'articles' | 'categories' | 'views'
  number: number
  label: string
}

const props = defineProps<Props>()

// 图标映射
const iconMap = {
  users: User,
  articles: Document,
  categories: Folder,
  views: View
}

// 获取图标
const icon = iconMap[props.type]
</script>

<style scoped>
.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  font-size: 24px;
  color: white;
}

.stat-icon.users {
  background-color: #409eff;
}

.stat-icon.articles {
  background-color: #67c23a;
}

.stat-icon.categories {
  background-color: #e6a23c;
}

.stat-icon.views {
  background-color: #f56c6c;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #666;
}
</style>