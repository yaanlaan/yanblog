<template>
  <nav class="navbar" :class="{ hidden: !isNavVisible, visible: isNavVisible }">
    <div class="container">
      <div class="navbar-content">

        <div class="logo-container">
          <router-link to="/">
            <div class="logo">
              <img :src="siteInfo.logo_image || '/assets/yan_icon/yaan.png'" class="avatar" alt="博客头像">
              <span class="blog-name">{{ siteInfo.blog_name }}</span>
            </div>
          </router-link>
        </div>
        
        <div class="navbar-center">
          <Transition name="fade-slide" mode="out-in">
            <div class="site-title-centered" v-if="isScrolled" key="title">
              <span>{{ siteInfo.blog_name }}</span>
            </div>

            <ul class="navbar-nav" v-else key="nav">

              <li class="nav-item">
                <router-link to="/" class="nav-link" :class="{ active: $route.name === 'home' }">
                  <i class="iconfont icon-Homehomepagemenu"></i>
                  <span>首页</span>
                </router-link>
              </li>

              <li class="nav-item">
                <router-link to="/articles" class="nav-link" :class="{ active: $route.name === 'articles' }">
                  <i class="iconfont icon-newspaper"></i>
                  <span>文章</span>
                </router-link>
              </li>

              <li class="nav-item">
                <router-link to="/categories" class="nav-link" :class="{ active: $route.name === 'categories' }">
                  <i class="iconfont icon-categories"></i>
                  <span>分类</span>
                </router-link>
              </li>

              <li class="nav-item">
                <router-link to="/archive" class="nav-link" :class="{ active: $route.name === 'archive' }">
                  <i class="iconfont icon-archive"></i>
                  <span>归档</span>
                </router-link>
              </li>

              <li class="nav-item dropdown-container">
                <router-link to="/about" class="nav-link" :class="{ active: $route.name === 'about' }">
                  <i class="iconfont icon-about"></i>
                  <span>关于</span>
                </router-link>
                
                <!-- 下拉菜单 -->
                <div class="dropdown-menu" v-if="siteInfo.socials && siteInfo.socials.length > 0">
                  <a 
                    v-for="(contact, index) in siteInfo.socials" 
                    :key="index" 
                    :href="contact.url" 
                    target="_blank" 
                    class="dropdown-item"
                  >
                    <i class="iconfont" :class="contact.icon" :style="{ color: contact.color }"></i>
                    <span>{{ contact.name }}</span>
                  </a>
                </div>
              </li>

            </ul>
          </Transition>
        </div>

        <div class="navbar-right">
          <!-- 主题切换按钮 -->
          <button class="theme-toggle-btn" @click="themeStore.toggleTheme" :title="themeStore.theme === 'dark' ? '切换亮色' : '切换暗色'">
            <!-- Sun Icon -->
            <svg v-if="themeStore.theme === 'light'" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"></circle><line x1="12" y1="1" x2="12" y2="3"></line><line x1="12" y1="21" x2="12" y2="23"></line><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line><line x1="1" y1="12" x2="3" y2="12"></line><line x1="21" y1="12" x2="23" y2="12"></line><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line></svg>
            <!-- Moon Icon -->
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
          </button>

          <div class="search-container">
            <input 
              type="text" 
              class="search-box" 
              placeholder="搜索文章..." 
              v-model="searchQuery"
              @keyup.enter="handleSearch"
            >
            <button class="search-icon" @click="handleSearch">
              <i class="iconfont icon-search"></i>
            </button>
          </div>

        </div>
      </div>
    </div>

    <a v-if="siteInfo.admin_url" :href="siteInfo.admin_url" target="_blank" class="login-btn admin-btn-fixed">
      登录
    </a>
    
    <!-- 底部线性进度条 -->
    <div class="progress-bar" :style="{ width: scrollProgress + '%' }"></div>
  </nav>

  <!-- 悬浮阅读进度 -->
  <div class="floating-progress" :class="{ visible: isScrolled }">
    <svg viewBox="0 0 36 36" class="circular-chart">
      <path class="circle-bg"
        d="M18 2.0845
          a 15.9155 15.9155 0 0 1 0 31.831
          a 15.9155 15.9155 0 0 1 0 -31.831"
      />
      <path class="circle"
        :stroke-dasharray="`${scrollProgress}, 100`"
        d="M18 2.0845
          a 15.9155 15.9155 0 0 1 0 31.831
          a 15.9155 15.9155 0 0 1 0 -31.831"
      />
      <text x="18" y="20.35" class="percentage">{{ Math.round(scrollProgress) }}</text>
    </svg>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { useSiteInfoStore } from '@/stores/siteInfo'
import { useThemeStore } from '@/stores/theme'
import { storeToRefs } from 'pinia'

const siteInfoStore = useSiteInfoStore()
const { siteInfo } = storeToRefs(siteInfoStore)
const themeStore = useThemeStore()

// 导航栏可见性状态
const isNavVisible = ref(true)
const isScrolled = ref(false)
const lastScrollY = ref(0)
const searchQuery = ref('')
const scrollProgress = ref(0)

const router = useRouter()

// 处理滚动事件
const handleScroll = () => {
  const currentScrollY = window.scrollY
  // 始终保持导航栏可见，但根据滚动距离切换显示内容
  isNavVisible.value = true
  isScrolled.value = currentScrollY > 100
  lastScrollY.value = currentScrollY
  
  // 计算阅读进度
  const docHeight = document.documentElement.scrollHeight - document.documentElement.clientHeight
  if (docHeight > 0) {
    scrollProgress.value = Math.min(100, Math.max(0, (currentScrollY / docHeight) * 100))
  } else {
    scrollProgress.value = 0
  }
}

// 处理搜索
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      name: 'articles',
      query: { search: searchQuery.value.trim() }
    })
  }
}

// 添加滚动事件监听器
onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

// 移除滚动事件监听器
onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.navbar {
  background-color: var(--color-background-soft);
  box-shadow: 0 2px 4px var(--color-shadow);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  width: 100%;
  z-index: 1000;
  transition: opacity 0.3s, transform 0.3s;
  backdrop-filter: blur(10px);
  box-sizing: border-box;

  margin: 0;
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: center; /* 水平居中 */
}

.navbar.hidden {
  opacity: 0;
  transform: translateY(-100%);
}

.navbar.visible {
  opacity: 1;
  transform: translateY(0);
}

.container {
  width: 100%;
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: center; /* 水平居中 */
  max-width: 1200px;
  margin: 0;
  padding: 0 20px;
  box-sizing: border-box;
}

.navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px; /* 固定高度 */
  width: 100%;
  box-sizing: border-box;
}

.logo-container {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: left; /* 水平居中 */
  width: 20%;
  height: 100%;
}

.logo {
  position: relative;
  padding: 0 10px;
  height: 40px; /* 减小高度 */
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.3s ease-in-out;
}

.logo:hover {
  transform: scale(1.05);
}

.logo::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  border-radius: 12px;
  background: linear-gradient(45deg, hsla(197, 100%, 50%, 1.00), hsla(177, 100%, 50%, 1.00));
  z-index: -1;
  animation: glow 2s linear infinite;
  opacity: 0;
  transition: opacity 0.3s ease-in-out;
}

.logo:hover::before {
  opacity: 1;
}

@keyframes glow {
  0% {
    filter: hue-rotate(0deg);
  }
  100% {
    filter: hue-rotate(360deg);
  }
}

.avatar {
  width: 32px; /* 减小头像 */
  height: 32px;
  object-fit: cover;
  border-radius: 10%;
  margin-right: 8px;
}

.blog-name {
  font-size: 1.1em; /* 稍微减小字体 */
  font-weight: bold;
  color: var(--color-heading);
}

.navbar-center {
  width: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.navbar-nav {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: center; /* 水平居中 */
  width: 100%;
  list-style: none;
  margin: 0;
  padding: 0;
}

.nav-item {
  margin: 0 5px;
}

.nav-link {
  display: flex;
  flex-direction: row; /* 改为水平排列 */
  align-items: center;
  justify-content: center;
  text-decoration: none;
  color: var(--color-text);
  padding: 6px 12px; /* 减小 padding */
  border-radius: 8px;
  transition: all 0.3s ease;
  min-width: auto; /* 移除最小宽度限制 */
  gap: 6px; /* 图标和文字间距 */
}

.nav-link:hover {
  color: var(--color-accent);
  transform: translateY(-1px); /* 减小位移 */
  background-color: var(--color-border-hover);
}

.nav-link.active {
  color: var(--color-accent);
  background-color: var(--color-border-hover);
}

.nav-link i {
  font-size: 1.1em;
  margin-bottom: 0; /* 移除底部间距 */
}

/* 下拉菜单样式 - Mac风格 */
.dropdown-container {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(10px);
  background: var(--color-background-soft);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 12px;
  padding: 6px;
  min-width: 140px;
  box-shadow: 
    0 4px 6px -1px var(--color-shadow),
    0 2px 4px -1px var(--color-shadow),
    0 0 0 1px var(--color-border) inset;
  border: 1px solid var(--color-border);
  opacity: 0;
  visibility: hidden;
  transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
  z-index: 1000;
  margin-top: 5px;
}


.dropdown-container:hover .dropdown-menu {
  opacity: 1;
  visibility: visible;
  transform: translateX(-50%) translateY(0);
}

/* 小三角箭头 */
.dropdown-menu::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-bottom: 6px solid var(--color-background-soft);
}

.dropdown-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  color: var(--color-text);
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.2s ease;
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
}

.dropdown-item:hover {
  background-color: var(--color-accent); /* Mac Blue */
  color: white;
}

.dropdown-item:hover i {
  color: white !important; /* 强制图标变白 */
}

.dropdown-item i {
  margin-right: 10px;
  font-size: 16px;
  width: 20px;
  text-align: center;
  transition: color 0.2s ease;
}

.navbar-right {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: right; /* 水平居中 */
  width: 30%;
  align-items: center;
  gap: 20px;
}

.theme-toggle-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--color-text);
  padding: 8px;
  border-radius: 50%;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.theme-toggle-btn:hover {
  background-color: var(--color-border-hover);
  color: var(--color-accent);
}

.search-container {
  display: flex;
  align-items: center;
  border: 1px solid var(--color-border);
  border-radius: 25px;
  overflow: hidden;
  width: 250px;
  transition: width 0.4s ease-in-out, box-shadow 0.3s ease;
  position: relative;
}

.search-container:focus-within {
  width: 300px;
  box-shadow: 0 0 0 2px var(--color-accent);
  border-color: var(--color-accent);
}

.search-box {
  height: 35px;
  flex-grow: 1;
  padding: 0 15px;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
  color: var(--color-text);
}

.search-box::placeholder {
  color: var(--color-text-secondary);
  opacity: 0.8;
}

.search-icon {
  padding: 8px 15px;
  cursor: pointer;
  background-color: transparent;
  border: none;
  outline: none;
  color: var(--color-text-secondary);
  transition: all 0.3s ease;
  border-radius: 0 25px 25px 0;
}

.search-container:focus-within .search-icon {
  color: var(--color-accent);
  background-color: var(--color-border-hover);
}

.login-btn {
  background-color: var(--color-accent);
  color: white;
  text-decoration: none;
  width: 36px;
  height: 36px;
  padding: 0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.login-btn:hover {
  background-color: #3aa876;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(66, 184, 131, 0.3);
}

.login-btn i {
  font-size: 20px;
}

.admin-btn-fixed {
  position: absolute;
  right: 20px;
  top: 12px; /* (60px header - 36px button) / 2 */
  z-index: 1002;
}

/* 悬浮阅读进度条样式 */
.floating-progress {
  position: fixed;
  top: 80px;
  right: 20px;
  width: 40px;
  height: 40px;
  z-index: 999;
  background: var(--color-background-soft);
  border-radius: 50%;
  box-shadow: 0 2px 10px var(--color-shadow);
  padding: 2px;
  opacity: 0;
  transform: translateX(20px);
  transition: all 0.3s ease;
  pointer-events: none;
  backdrop-filter: blur(5px);
}

.floating-progress.visible {
  opacity: 1;
  transform: translateX(0);
  pointer-events: auto;
}

.circular-chart {
  display: block;
  margin: 0 auto;
  max-width: 100%;
  max-height: 100%;
}

.circle-bg {
  fill: none;
  stroke: #eee;
  stroke-width: 3.8;
}

.circle {
  fill: none;
  stroke-width: 2.8;
  stroke-linecap: round;
  stroke: #42b883;
  animation: progress 1s ease-out forwards;
  transition: stroke-dasharray 0.1s;
}

.percentage {
  fill: #666;
  font-family: sans-serif;
  font-weight: bold;
  font-size: 10px;
  text-anchor: middle;
}

.progress-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  height: 3px;
  background: linear-gradient(to right, #42b883, #3aa876);
  transition: width 0.1s;
  z-index: 1001;
}

.site-title-centered {
  width: 100%;
  height: 100%; /* 撑满高度 */
  text-align: center;
  font-size: 1.3em; /* 稍微减小字体 */
  font-weight: bold;
  color: var(--color-heading);
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 动画效果 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

@media (max-width: 768px) {
  .navbar-content {
    flex-wrap: wrap;
  }
  
  .navbar-nav,
  .site-title-centered,
  .navbar-center {
    order: 3;
    width: 100%;
    justify-content: center;
    margin-top: 10px;
  }
  
  .navbar-right {
    order: 2;
    width: 100%;
    margin: 10px 0;
    justify-content: center;
  }
  
  .search-container {
    margin: 0;
    width: 100%;
    max-width: 300px;
  }
  
  .logo-container {
    order: 1;
  }
  
  .login-btn {
    order: 4;
    margin: 10px auto;
  }
}
</style>