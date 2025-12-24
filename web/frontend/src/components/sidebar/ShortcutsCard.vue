<template>
  <div class="sidebar-card shortcuts-card">
    <div class="card-header">
      <h3><i class="iconfont icon-link"></i> 常用链接</h3>
    </div>
    <div class="card-content">
      <div class="shortcut-list">
        <a v-for="(item, index) in siteInfo.shortcuts" :key="index" :href="item.url" target="_blank" class="shortcut-item">
          <div class="icon-wrapper" :style="{ background: item.color }">
            <i :class="['iconfont', item.icon]"></i>
          </div>
          <span class="name">{{ item.name }}</span>
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSiteInfoStore } from '@/stores/siteInfo'
import { storeToRefs } from 'pinia'

const siteInfoStore = useSiteInfoStore()
const { siteInfo } = storeToRefs(siteInfoStore)
</script>

<style scoped>
.shortcuts-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  margin-bottom: 20px;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
  background: white;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-content {
  padding: 15px;
}

.shortcut-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.shortcut-item {
  display: flex;
  align-items: center;
  padding: 10px;
  background: #f9f9f9;
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.3s ease;
}

.shortcut-item:hover {
  background: #f0f0f0;
  transform: translateX(5px);
}

.icon-wrapper {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  color: white;
  font-size: 18px;
  flex-shrink: 0;
}

.name {
  color: #333;
  font-size: 14px;
  font-weight: 500;
}
</style>
