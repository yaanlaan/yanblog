<template>
  <div class="main-wrapper">
    <div class="container" :class="{ 'full-width': fullWidth }">
      <div class="content-wrapper">
        <div class="sidebar left-sidebar" v-if="$slots.leftSidebar">
          <slot name="leftSidebar"></slot>
        </div>
        <div class="main-content">
          <slot name="main"></slot>
        </div>
        <div class="sidebar right-sidebar" v-if="$slots.sidebar">
          <slot name="sidebar"></slot>
        </div>
      </div>
      <BackToTop />
    </div>
  </div>
</template>

<script setup lang="ts">
// 主布局组件
import BackToTop from '@/components/BackToTop.vue'

defineProps<{
  fullWidth?: boolean
}>()
</script>

<style scoped>
.main-wrapper {
  width: 100%;
  margin-TOP: 10px;
}

.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.container.full-width {
  max-width: 100%;
  padding: 0 40px;
}

@media (max-width: 768px) {
  .container, .container.full-width {
    padding: 0 16px; /* Decrease padding on mobile */
  }
}

.content-wrapper {
  display: flex;
  gap: 20px;
  min-height: calc(100vh - 280px);
  align-items: flex-start;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.left-sidebar {
  width: 250px; /* 目录栏宽度 */
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: sticky;
  top: 80px;
}

.right-sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: sticky;
  top: 80px;
}

@media (max-width: 1200px) {
  .left-sidebar {
    width: 200px;
  }
  
  .right-sidebar {
    width: 250px;
  }
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .left-sidebar,
  .right-sidebar {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
    position: static;
    max-height: none;
    overflow-y: visible;
  }
  
  .sidebar-card {
    flex: 1 1 calc(50% - 10px);
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 0 15px;
  }
  
  .left-sidebar,
  .right-sidebar {
    flex-direction: column;
  }
  
  .sidebar-card {
    flex: 1 1 100%;
  }
}
</style>