<template>
  <nav class="navbar" :class="{ hidden: !isNavVisible, visible: isNavVisible }">
    <div class="container">
      <div class="navbar-content">

        <div class="logo-container">
          <router-link to="/">
            <div class="logo">
              <img src="../assets/yan_icon/yaan.png" class="avatar" alt="博客头像">
              <span class="blog-name">言盐盐的博客</span>
            </div>
          </router-link>
        </div>
        
        <ul class="navbar-nav">

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
            <router-link to="/about" class="nav-link" :class="{ active: $route.name === 'about' }">
              <i class="iconfont icon-about"></i>
              <span>关于</span>
            </router-link>
          </li>

        </ul>

        <div class="navbar-right">
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

          <!-- <button class="login-btn">未登录</button> -->
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'

// 导航栏可见性状态
const isNavVisible = ref(true)
const lastScrollY = ref(0)
const searchQuery = ref('')

const router = useRouter()

// 处理滚动事件
const handleScroll = () => {
  const currentScrollY = window.scrollY
  // 当向下滚动超过50px时隐藏导航栏，向上滚动时显示
  isNavVisible.value = currentScrollY < 50 || currentScrollY < lastScrollY.value
  lastScrollY.value = currentScrollY
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
  background-color: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  width: 100%;
  z-index: 1000;
  transition: opacity 0.3s, transform 0.3s;
  backdrop-filter: blur(10px);
  box-sizing: border-box;

  position: 0;
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
  padding: 10px 0;
  width: 100%;
  box-sizing: border-box;
}

.logo-container {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: left; /* 水平居中 */
  width: 20%;

}

.logo {
  position: relative;
  width: 180px;
  height: 60px;
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
  background: linear-gradient(45deg, #0000ff, #0033ff, #0066ff, #0099ff, #00ccff, #00ffff, #33ffff, #66ffff, #99ffff, #ccffff, #ffffff);
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
  width: 45px;
  height: 45px;
  object-fit: cover;
  border-radius: 10%;
  margin-right: 10px;
}

.blog-name {
  font-size: 1.3em;
  font-weight: bold;
  color: #333;
}

.navbar-nav {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: center; /* 水平居中 */
  width: 50%;
  list-style: none;
  margin: 0;
  padding: 0;
}

.nav-item {
  margin: 0 5px;
}

.nav-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  color: #2c3e50;
  padding: 8px 15px;
  border-radius: 8px;
  transition: all 0.3s ease;
  min-width: 60px;
}

.nav-link:hover {
  color: #3994fc;
  transform: translateY(-2px);
  background-color: rgba(61, 150, 252, 0.1);
}

.nav-link.active {
  color: #7cb7fa;
  background-color: rgba(61, 150, 252, 0.1);
}

.nav-link i {
  font-size: 1.2em;
  margin-bottom: 4px;
}

.navbar-right {
  display: flex;
  align-items: center; /* 垂直居中 */
  justify-content: right; /* 水平居中 */
  width: 30%;
  align-items: center;
  gap: 20px;
}

.search-container {
  display: flex;
  align-items: center;
  border: 1px solid #ddd;
  border-radius: 25px;
  overflow: hidden;
  width: 250px;
  transition: width 0.4s ease-in-out, box-shadow 0.3s ease;
  position: relative;
}

.search-container:focus-within {
  width: 300px;
  box-shadow: 0 0 0 2px #3d96fc;
  border-color: #3d96fc;
}

.search-box {
  height: 35px;
  flex-grow: 1;
  padding: 0 15px;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
}

.search-icon {
  padding: 8px 15px;
  cursor: pointer;
  background-color: transparent;
  border: none;
  outline: none;
  color: #666;
  transition: all 0.3s ease;
  border-radius: 0 25px 25px 0;
}

.search-container:focus-within .search-icon {
  color: #3d96fc;
  background-color: rgba(61, 150, 252, 0.1);
}

.login-btn {
  background-color: #3d96fc;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.login-btn:hover {
  background-color: #2a85f5;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(61, 150, 252, 0.3);
}

@media (max-width: 768px) {
  .navbar-content {
    flex-wrap: wrap;
  }
  
  .navbar-nav {
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