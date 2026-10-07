<template>
  <div class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
      <div class="logo">
        <h2 v-if="!isCollapse">博客管理系统</h2>
        <el-icon v-else :size="22" color="#fff"><Management /></el-icon>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        :collapse="isCollapse"
        :collapse-transition="false"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>仪表板</span>
        </el-menu-item>
        
        <el-menu-item index="/user">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
        
        <el-menu-item index="/category">
          <el-icon><Folder /></el-icon>
          <span>分类管理</span>
        </el-menu-item>

        <el-menu-item index="/tag">
          <el-icon><Collection /></el-icon>
          <span>标签管理</span>
        </el-menu-item>
        
        <el-menu-item index="/article">
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </el-menu-item>

        <el-menu-item index="/media">
          <el-icon><Picture /></el-icon>
          <span>媒体库</span>
        </el-menu-item>

        <el-sub-menu index="/system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </template>
          <el-menu-item index="/system/status">系统监控</el-menu-item>
          <el-menu-item index="/system/config">前台配置</el-menu-item>
          <el-menu-item index="/system/backend">后端配置</el-menu-item>
          <el-menu-item index="/system/about">关于页管理</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    
    <!-- 主体内容 -->
    <el-container>
      <!-- 顶部栏 -->
      <el-header class="header">
        <div class="header-left">
          <el-button 
            link 
            class="collapse-btn" 
            @click="isCollapse = !isCollapse"
            :title="isCollapse ? '展开侧边栏' : '折叠侧边栏'"
          >
            <el-icon :size="20">
              <component :is="isCollapse ? Expand : Fold" />
            </el-icon>
          </el-button>

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
          <!-- 快捷浏览前台入口 -->
          <el-button 
            type="primary" 
            plain 
            size="small" 
            :icon="View" 
            @click="openBlog" 
            class="visit-blog-btn"
          >
            浏览前台 ↗
          </el-button>

          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="26" :icon="User" class="user-avatar" />
              <span class="username">{{ username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="toConfig">前台配置</el-dropdown-item>
                <el-dropdown-item command="toStatus">系统监控</el-dropdown-item>
                <el-dropdown-item divided command="logout" style="color: #f56c6c;">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 内容区域 -->
      <el-main class="main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Odometer, User, Folder, Document, Picture, 
  Setting, Collection, Fold, Expand, View, Management 
} from '@element-plus/icons-vue'

// 路由实例
const route = useRoute()
const router = useRouter()

// 侧边栏折叠状态
const isCollapse = ref(false)

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
  return '管理员'
}

// 打开前台
const openBlog = () => {
  window.open('/', '_blank')
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
  if (meta?.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

// 处理下拉菜单命令
const handleCommand = (command: string) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/login')
    ElMessage.success('已安全退出登录')
  } else if (command === 'toConfig') {
    router.push('/system/config')
  } else if (command === 'toStatus') {
    router.push('/system/status')
  }
}

// 组件挂载时获取用户信息
onMounted(() => {
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
  transition: width 0.25s ease;
  box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);
  overflow-x: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #253342;
  transition: all 0.25s;
}

.logo h2 {
  color: #fff;
  font-size: 16px;
  margin: 0;
  white-space: nowrap;
}

.sidebar-menu {
  border: none;
  height: calc(100% - 60px);
}

.header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.collapse-btn {
  color: #606266;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.collapse-btn:hover {
  background-color: #f2f3f5;
  color: #409eff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.visit-blog-btn {
  font-size: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background 0.2s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-avatar {
  background-color: #409eff;
  color: #fff;
}

.username {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

/* 页面切换动画 */
.fade-transform-leave-active,
.fade-transform-enter-active {
  transition: all 0.2s ease;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(10px);
}
</style>
