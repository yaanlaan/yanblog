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
import { onMounted } from 'vue'
import { useSiteInfoStore } from '@/stores/siteInfo'

const siteInfoStore = useSiteInfoStore()

onMounted(async () => {
  await siteInfoStore.fetchSiteInfo()
  
  // 动态加载 iconfont
  if (siteInfoStore.siteInfo.iconfont_url) {
    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = siteInfoStore.siteInfo.iconfont_url
    document.head.appendChild(link)
  } else {
    // 如果没有配置 iconfont_url，使用默认的离线图标样式（简单的圆点作为占位符）
    const style = document.createElement('style')
    style.textContent = `
      .iconfont {
        display: inline-block;
        font-style: normal;
        vertical-align: baseline;
        text-align: center;
        text-transform: none;
        line-height: 1;
        text-rendering: optimizeLegibility;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
      }
      .iconfont::before {
        content: "\\25CF"; /* 实心圆点 */
      }
    `
    document.head.appendChild(style)
  }

  // 设置 Favicon
  if (siteInfoStore.siteInfo.favicon) {
    let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
    if (!link) {
      link = document.createElement('link')
      link.rel = 'icon'
      document.head.appendChild(link)
    }
    link.href = siteInfoStore.siteInfo.favicon
  }

  // 设置页面标题
  if (siteInfoStore.siteInfo.page_title?.default) {
    document.title = siteInfoStore.siteInfo.page_title.default
  }

  // 处理页面标题模糊效果
  let title = document.title
  window.onblur = function () {
    title = document.title
    if (siteInfoStore.siteInfo.page_title?.blur) {
      document.title = siteInfoStore.siteInfo.page_title.blur
    }
  }
  window.onfocus = function () {
    if (title) {
      document.title = title
    }
  }
})
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
  background-color: var(--color-background-overlay); /* 主题适配背景遮罩 */
  backdrop-filter: blur(20px); /* 磨砂模糊效果 */
  transition: background-color 0.5s ease;
}

/* App layout styles */
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: transparent; /* 防止背景色遮挡 global-bg */
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