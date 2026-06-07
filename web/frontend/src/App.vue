<template>
  <div class="app">
    <LoadingSpinner :loading="isLoading" />
    
    <div class="global-bg"></div>
    
    <header class="header" v-show="!isLoading">
      <NavBar />
    </header>

    <main class="main" v-show="!isLoading">
      <Transition name="fade" mode="out-in">
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </Transition>
    </main>

    <footer class="footer" v-show="!isLoading">
      <Footer />
    </footer>
    
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '@/components/NavBar.vue'
import Footer from '@/components/Footer.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import { useSiteInfoStore } from '@/stores/siteInfo'

const isLoading = ref(true)
const siteInfoStore = useSiteInfoStore()

onMounted(async () => {
  await siteInfoStore.fetchSiteInfo()
  
  if (siteInfoStore.siteInfo.iconfont_url) {
    const link = document.createElement('link')
    link.rel = 'stylesheet'
    link.href = siteInfoStore.siteInfo.iconfont_url
    document.head.appendChild(link)
  } else {
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
        content: "\\25CF";
      }
    `
    document.head.appendChild(style)
  }

  if (siteInfoStore.siteInfo.favicon) {
    let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
    if (!link) {
      link = document.createElement('link')
      link.rel = 'icon'
      document.head.appendChild(link)
    }
    link.href = siteInfoStore.siteInfo.favicon
  }

  if (siteInfoStore.siteInfo.page_title?.default) {
    document.title = siteInfoStore.siteInfo.page_title.default
  }

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

  setTimeout(() => {
    isLoading.value = false
  }, 500)
})
</script>

<style scoped>
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
  background-color: var(--color-background-overlay);
  backdrop-filter: blur(20px);
  transition: background-color 0.5s ease;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: transparent;
}

.header {
  flex-shrink: 0;
  height: 60px;
  animation: slideDown 0.5s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.main {
  flex: 1;
  width: 100%;
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.footer {
  margin-top: auto;
  width: 100%;
  flex-shrink: 0;
  animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

[v-loading] {
  position: relative;
}
</style>