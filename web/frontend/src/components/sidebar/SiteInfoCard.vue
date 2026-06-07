<template>
  <div class="sidebar-card site-info-card">
    <div class="card-header">
      <h3><i class="iconfont icon-chart" style="color: #42b883; margin-right: 5px;"></i> 站点信息</h3>
    </div>
    <div class="card-content">
      <div class="info-item">
        <span class="label">文章总数</span>
        <span class="value">36</span>
      </div>
      <div class="info-item">
        <span class="label">建站天数</span>
        <span class="value">{{ runDays }} 天</span>
      </div>
      <div class="info-item">
        <span class="label">全站字数</span>
        <span class="value">50k</span>
      </div>
      <div class="info-item">
        <span class="label">最后更新</span>
        <span class="value">2025-12-24</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const runDays = ref(0)

onMounted(() => {
  const startDate = new Date('2022-01-01') // Replace with actual start date
  const now = new Date()
  const diff = now.getTime() - startDate.getTime()
  runDays.value = Math.floor(diff / (1000 * 60 * 60 * 24))
})
</script>

<style scoped>
.info-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  font-size: 14px;
  border-bottom: 1px dashed #eee;
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  color: #666;
}

.value {
  color: #333;
  font-weight: bold;
}
</style>
