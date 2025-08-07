<template>
  <nav class="navbar">
    <div class="container">
      <div class="navbar-content">
        <div class="navbar-brand">
          <router-link to="/">言盐盐的博客</router-link>
        </div>
        
        <ul class="navbar-nav">
          <li class="nav-item">
            <router-link to="/" class="nav-link" :class="{ active: $route.name === 'home' }">
              首页
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/articles" class="nav-link" :class="{ active: $route.name === 'articles' }">
              文章列表
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/categories" class="nav-link" :class="{ active: $route.name === 'categories' }">
              分类
            </router-link>
          </li>
          <li class="nav-item">
            <router-link to="/about" class="nav-link" :class="{ active: $route.name === 'about' }">
              关于
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
// 导航栏组件
</script>

<style scoped>
.navbar {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.navbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
}

.navbar-brand a {
  font-size: 1.5rem;
  font-weight: bold;
  color: #333;
  text-decoration: none;
}

.navbar-nav {
  display: flex;
  list-style: none;
}

.nav-item {
  margin-left: 2rem;
}

.nav-link {
  color: #333;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s ease;
  padding: 0.5rem 1rem;
  border-radius: 4px;
}

.nav-link:hover,
.nav-link.active {
  color: #b352ff;
  background-color: rgba(179, 82, 255, 0.1);
}

@media (max-width: 768px) {
  .navbar-content {
    flex-direction: column;
  }
  
  .navbar-nav {
    margin-top: 1rem;
  }
  
  .nav-item {
    margin: 0 0.5rem;
  }
}
</style>