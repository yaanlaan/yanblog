<template>
  <div class="not-found-page">
    <div class="container">
      <div class="content">
        <div class="error-code-wrapper">
          <h1 class="error-code">404</h1>
          <div class="cat-container">
            <svg class="cat-svg" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 4C8.5 4 6 6.5 6 9C6 11.5 8.5 13 11 13C11 15 12 17 12 17C12 17 13 15 13 13C15.5 13 18 11.5 18 9C18 6.5 15.5 4 12 4Z" fill="var(--color-heading)" stroke="var(--color-heading)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M9 8L9 9" stroke="var(--color-background)" stroke-width="2" stroke-linecap="round"/>
              <path d="M15 8L15 9" stroke="var(--color-background)" stroke-width="2" stroke-linecap="round"/>
              <path d="M12 11L12 11.5" stroke="var(--color-background)" stroke-width="2" stroke-linecap="round"/>
              <path d="M7 4L5 2" stroke="var(--color-heading)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M17 4L19 2" stroke="var(--color-heading)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M12 17V21M12 21L10 19M12 21L14 19" stroke="var(--color-heading)" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
        </div>
        
        <h2 class="title">哎呀，页面被猫咪吃掉了！</h2>
        <p class="description">
          你访问的页面似乎去流浪了，或者原本就不存在。<br>
          不妨回到首页看看其他有趣的内容吧。
        </p>
        
        <div class="actions">
          <router-link to="/" class="btn primary-btn">
            <i class="iconfont icon-home"></i> 返回首页
          </router-link>
          <button @click="goBack" class="btn secondary-btn">
            <i class="iconfont icon-back"></i> 返回上一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.not-found-page {
  min-height: calc(100vh - 120px); /* 减去header footer高度的大致估算 */
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 20px;
}

.container {
  max-width: 600px;
  width: 100%;
}

.content {
  background: var(--color-background-soft);
  padding: 60px 40px;
  border-radius: 20px;
  box-shadow: 0 4px 20px var(--color-shadow);
  border: 1px solid var(--color-border);
  transition: all 0.3s ease;
}

.error-code-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 20px;
}

.error-code {
  font-size: 120px;
  font-weight: 900;
  color: var(--color-accent);
  line-height: 1;
  margin: 0;
  text-shadow: 4px 4px 0px var(--color-background-mute);
  position: relative;
  z-index: 1;
}

.cat-container {
  position: absolute;
  top: -40px;
  right: -20px;
  z-index: 2;
  animation: float 3s ease-in-out infinite;
}

.cat-svg {
  width: 80px;
  height: 80px;
  filter: drop-shadow(0 4px 6px var(--color-shadow));
}

.title {
  font-size: 24px;
  color: var(--color-heading);
  margin-bottom: 15px;
  font-weight: bold;
}

.description {
  font-size: 16px;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin-bottom: 30px;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  border: none;
  font-weight: 500;
}

.primary-btn {
  background: var(--color-accent);
  color: white;
  text-decoration: none;
}

.primary-btn:hover {
  background: #3aa876; /* Slightly darker accent */
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(66, 184, 131, 0.3);
}

.secondary-btn {
  background: var(--color-background-mute);
  color: var(--color-text);
  border: 1px solid var(--color-border);
}

.secondary-btn:hover {
  background: var(--color-border);
  transform: translateY(-2px);
}

@keyframes float {
  0% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-10px) rotate(5deg); }
  100% { transform: translateY(0px) rotate(0deg); }
}

@media (max-width: 480px) {
  .error-code {
    font-size: 80px;
  }
  
  .content {
    padding: 40px 20px;
  }
  
  .actions {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>