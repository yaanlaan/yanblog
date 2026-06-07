<template>
  <div class="loading-overlay" :class="{ hidden: !loading }">
    <div class="loading-container">
      <div class="loading-spinner">
        <div class="spinner-ring"></div>
        <div class="spinner-ring delay-1"></div>
        <div class="spinner-ring delay-2"></div>
        <div class="spinner-center">
          <span class="logo-text">Y</span>
        </div>
      </div>
      <div class="loading-text">
        <span class="loading-words">{{ loadingText }}</span>
        <span class="loading-dots">
          <span class="dot">.</span>
          <span class="dot delay-1">.</span>
          <span class="dot delay-2">.</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

defineProps<{
  loading: boolean
}>()

const loadingText = ref('正在加载')

const loadingTexts = [
  '正在加载',
  '努力加载中',
  '马上就好',
  '稍等片刻'
]

let textInterval: ReturnType<typeof setInterval>

onMounted(() => {
  let textIndex = 0
  textInterval = setInterval(() => {
    textIndex = (textIndex + 1) % loadingTexts.length
    loadingText.value = loadingTexts[textIndex]
  }, 2000)
})

onUnmounted(() => {
  clearInterval(textInterval)
})
</script>

<style scoped>
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, var(--color-background) 0%, var(--color-background-soft) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  transition: opacity 0.5s ease, visibility 0.5s ease;
}

.loading-overlay.hidden {
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 30px;
}

.loading-spinner {
  position: relative;
  width: 120px;
  height: 120px;
}

.spinner-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100%;
  height: 100%;
  border: 3px solid transparent;
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: spin 2s linear infinite;
}

.spinner-ring:nth-child(1) {
  border-top-color: var(--color-accent);
}

.spinner-ring.delay-1 {
  animation-delay: 0.3s;
  border-top-color: var(--color-accent-light);
  width: 85%;
  height: 85%;
}

.spinner-ring.delay-2 {
  animation-delay: 0.6s;
  border-top-color: var(--color-accent-dark);
  width: 70%;
  height: 70%;
}

@keyframes spin {
  0% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  100% {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

.spinner-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 20px rgba(66, 184, 131, 0.4);
}

.logo-text {
  font-size: 24px;
  font-weight: bold;
  color: white;
}

.loading-text {
  display: flex;
  align-items: center;
  font-size: 18px;
  color: var(--color-text-secondary);
}

.loading-words {
  font-weight: 500;
}

.loading-dots {
  display: flex;
  margin-left: 2px;
}

.dot {
  animation: dotBounce 1s infinite;
}

.dot.delay-1 {
  animation-delay: 0.2s;
}

.dot.delay-2 {
  animation-delay: 0.4s;
}

@keyframes dotBounce {
  0%, 100% {
    opacity: 0.3;
    transform: translateY(0);
  }
  50% {
    opacity: 1;
    transform: translateY(-4px);
  }
}
</style>