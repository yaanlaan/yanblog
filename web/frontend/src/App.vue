<template>
  <div class="app">
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
/* App layout styles */
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  flex-shrink: 0; /* 防止头部被压缩 */
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