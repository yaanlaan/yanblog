<template>
  <div id="app">
    <header class="header">
      <div class="navbar">
        <div class="container">
          <div class="navbar-content">
            <div class="navbar-brand">
              <router-link to="/">我的博客</router-link>
            </div>
            <ul class="navbar-nav">
              <li class="nav-item">
                <router-link to="/" class="nav-link" exact-active-class="active">
                  首页
                </router-link>
              </li>
              <li class="nav-item">
                <router-link to="/articles" class="nav-link" exact-active-class="active">
                  文章
                </router-link>
              </li>
              <li class="nav-item">
                <router-link to="/about" class="nav-link" exact-active-class="active">
                  关于
                </router-link>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </header>

    <main class="main">
      <router-view />
    </main>

    <footer class="footer">
      <div class="container">
        <p>&copy; 2025 我的博客. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
// App根组件
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
}

body {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  color: #333;
  background-color: #f5f7fa;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header styles */
.header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative; /* 修改为relative避免内容遮挡 */
  width: 100%;
  flex-shrink: 0; /* 防止头部被压缩 */
}

.navbar {
  width: 100%;
  background: white;
}

.container {
  width: 70%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.navbar-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
}

.navbar-brand a {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  text-decoration: none;
}

.navbar-nav {
  display: flex;
  list-style: none;
  gap: 30px;
}

.nav-link {
  text-decoration: none;
  color: #666;
  font-weight: 500;
  transition: color 0.3s;
  padding: 5px 10px;
  border-radius: 4px;
}

.nav-link:hover,
.nav-link.active {
  color: #007bff;
  background-color: rgba(0, 123, 255, 0.1);
}

/* Main content */
.main {
  flex: 1;
  width: 100%;
  padding: 20px 0;
}

/* Footer styles */
.footer {
  background: #333;
  color: white;
  padding: 20px 0;
  text-align: center;
  margin-top: auto;
  width: 100%;
  flex-shrink: 0; /* 防止页脚被压缩 */
}

/* Element Plus loading */
[v-loading] {
  position: relative;
}

/* Responsive design */
@media (max-width: 768px) {
  .navbar-content {
    flex-direction: column;
    gap: 15px;
  }
  
  .navbar-nav {
    gap: 15px;
  }
  
  .container {
    width: 90%;
    padding: 0 15px;
  }
}
</style>