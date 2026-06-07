<template>
  <div class="skeleton-container">
    <div v-for="n in count" :key="n" class="skeleton-card">
      <div class="skeleton-header">
        <div class="skeleton-avatar"></div>
        <div class="skeleton-info">
          <div class="skeleton-name"></div>
          <div class="skeleton-time"></div>
        </div>
      </div>
      <div class="skeleton-title"></div>
      <div class="skeleton-content">
        <div class="skeleton-line"></div>
        <div class="skeleton-line"></div>
        <div class="skeleton-line short"></div>
      </div>
      <div class="skeleton-footer">
        <div class="skeleton-category"></div>
        <div class="skeleton-tags">
          <div class="skeleton-tag"></div>
          <div class="skeleton-tag"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  count?: number
}>()
</script>

<style scoped>
.skeleton-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  padding: 20px 0;
}

.skeleton-card {
  background: var(--color-background-soft);
  border-radius: 12px;
  padding: 20px;
  animation: skeleton-loading 1.5s ease-in-out infinite;
}

@keyframes skeleton-loading {
  0%, 100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

.skeleton-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.skeleton-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
}

.skeleton-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.skeleton-name {
  height: 16px;
  width: 60%;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 4px;
}

.skeleton-time {
  height: 12px;
  width: 40%;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 4px;
}

.skeleton-title {
  height: 24px;
  width: 80%;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 4px;
  margin-bottom: 12px;
}

.skeleton-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.skeleton-line {
  height: 14px;
  width: 100%;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 4px;
}

.skeleton-line.short {
  width: 60%;
}

.skeleton-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.skeleton-category {
  height: 24px;
  width: 60px;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 12px;
}

.skeleton-tags {
  display: flex;
  gap: 8px;
}

.skeleton-tag {
  height: 20px;
  width: 40px;
  background: linear-gradient(90deg, var(--color-border) 25%, var(--color-border-hover) 50%, var(--color-border) 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
  border-radius: 10px;
}

@keyframes skeleton-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

@media (max-width: 768px) {
  .skeleton-container {
    grid-template-columns: 1fr;
  }
}
</style>