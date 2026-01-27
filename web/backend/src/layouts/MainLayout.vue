<template>
  <div class="layout-container">
    <!-- 侧边栏 -->
    <el-aside width="200px" class="sidebar">
      <div class="logo">
        <h2>博客管理系统</h2>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>仪表板</span>
        </el-menu-item>
        
        <el-sub-menu index="/user">
          <template #title>
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </template>
          <el-menu-item index="/user/list">用户列表</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/category">
          <template #title>
            <el-icon><Folder /></el-icon>
            <span>分类管理</span>
          </template>
          <el-menu-item index="/category/list">分类列表</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/article">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>文章管理</span>
          </template>
          <el-menu-item index="/article/list">文章列表</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/media">
          <el-icon><Picture /></el-icon>
          <span>媒体库</span>
        </el-menu-item>

        <el-sub-menu index="/system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="/system/config">前台配置</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <!-- 主体内容 -->
    <el-container>
      <!-- 顶部栏 -->
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item 
              v-for="item in breadcrumbItems" 
              :key="item.path"
            >
              {{ item.name }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <span class="username">{{ username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 内容区域 -->
      <el-main class="main">
        <router-view v-slot="{ Component, route }">
          <keep-alive>
             <component :is="Component" :key="route.fullPath" />
          </keep-alive>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Odometer, User, Folder, Document, Picture, Setting } from '@element-plus/icons-vue'

// 获取路由实例
const route = useRoute()
const router = useRouter()

// 用户名
const username = ref('')

// 解析JWT token获取用户信息
const parseJwt = (token: string) => {
  try {
    const base64Url = token.split('.')[1]
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(jsonPayload)
  } catch (error) {
    console.error('解析token失败:', error)
    return null
  }
}

// 从token获取用户名
const getUsernameFromToken = () => {
  const token = localStorage.getItem('token')
  if (token) {
    const payload = parseJwt(token)
    if (payload && payload.username) {
      return payload.username
    }
  }
  return '未知用户'
}

// 面包屑导航项
const breadcrumbItems = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    path: item.path,
    name: item.meta?.title as string
  }))
})

// 激活菜单项
const activeMenu = computed(() => {
  const { meta, path } = route
  // 如果当前路由设置了activeMenu，则使用它
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  if (command === 'logout') {
    // 清除本地存储的token
    localStorage.removeItem('token')
    // 跳转到登录页
    router.push('/login')
    ElMessage.success('已退出登录')
  }
}

// 组件挂载时获取用户信息
onMounted(() => {
  // 从token解析获取用户名
  username.value = getUsernameFromToken()
})
</script>

<style scoped>
.layout-container {
  height: 100vh;
  display: flex;
}

.sidebar {
  background-color: #304156;
  color: #fff;
  transition: width 0.28s;
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #253342;
}

.logo h2 {
  color: #fff;
  font-size: 18px;
  margin: 0;
}

.sidebar-menu {
  border: none;
  height: calc(100% - 60px);
}

.header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.12);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 10px;
  font-size: 14px;
  color: #666;
}

.main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>