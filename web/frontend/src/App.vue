<template>
  <div class="app">
    <div class="global-bg"></div>
    <header class="header">
      <NavBar />
    </header>

    <main class="main">
      <router-view />
    </main>

    <footer class="footer">
      <Footer />
    </footer>
    
  </div>
</template>

<script setup lang="ts">
// App根组件
import NavBar from '@/components/NavBar.vue'
import Footer from '@/components/Footer.vue'
</script>

<style scoped>
/* Global Background */
.global-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: -1;
  background-image: url('@/assets/img/lonelycat.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.global-bg::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.85); /* 白色半透明遮罩 */
  backdrop-filter: blur(20px); /* 磨砂模糊效果 */
}

/* App layout styles */
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  flex-shrink: 0; /* 防止头部被压缩 */
  height: 60px; /* 占位高度，与 NavBar 高度一致 */
}

.main {
  flex: 1; /* 占据剩余空间 */
  width: 100%;
}

/* Footer styles */
.footer {
  margin-top: auto;
  width: 100%;
  flex-shrink: 0; /* 防止页脚被压缩 */
}

/* Element Plus loading */
[v-loading] {
  position: relative;
}
</style>