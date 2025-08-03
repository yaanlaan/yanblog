<template>
  <div class="container">
    <div class="content-wrapper">
      <div class="main-content">
        <slot name="main"></slot>
      </div>
      <div class="sidebar">
        <slot name="sidebar"></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 主布局组件
</script>

<style scoped>
.container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.content-wrapper {
  display: flex;
  gap: 30px;
  min-height: calc(100vh - 280px);
}

.main-content {
  flex: 1;
  min-width: 0;
}

.sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
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
  
  .sidebar {
    flex-direction: column;
  }
  
  .sidebar-card {
    flex: 1 1 100%;
  }
}
</style>