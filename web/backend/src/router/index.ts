import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 布局组件
import MainLayout from '@/layouts/MainLayout.vue'

// 页面组件
const Login = () => import('@/views/Login.vue')
const Dashboard = () => import('@/views/dashboard/Dashboard.vue')
const UserList = () => import('@/views/user/UserList.vue')
const CategoryList = () => import('@/views/category/CategoryList.vue')
const ArticleList = () => import('@/views/article/ArticleList.vue')

// 定义路由
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: {
      title: '登录'
    }
  },
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: {
          title: '仪表板',
          icon: 'Odometer'
        }
      },
      {
        path: '/user',
        name: 'User',
        redirect: '/user/list',
        meta: {
          title: '用户管理',
          icon: 'User'
        }
      },
      {
        path: '/user/list',
        name: 'UserList',
        component: UserList,
        meta: {
          title: '用户列表',
          activeMenu: '/user'
        }
      },
      {
        path: '/category',
        name: 'Category',
        redirect: '/category/list',
        meta: {
          title: '分类管理',
          icon: 'Folder'
        }
      },
      {
        path: '/category/list',
        name: 'CategoryList',
        component: CategoryList,
        meta: {
          title: '分类列表',
          activeMenu: '/category'
        }
      },
      {
        path: '/article',
        name: 'Article',
        redirect: '/article/list',
        meta: {
          title: '文章管理',
          icon: 'Document'
        }
      },
      {
        path: '/article/list',
        name: 'ArticleList',
        component: ArticleList,
        meta: {
          title: '文章列表',
          activeMenu: '/article'
        }
      },
      {
        path: '/article/add',
        name: 'ArticleAdd',
        component: () => import('@/views/article/ArticleEditor.vue'),
        meta: {
          title: '新增文章'
        }
      },
      {
        path: '/article/edit/:id',
        name: 'ArticleEdit',
        component: () => import('@/views/article/ArticleEditor.vue'),
        meta: {
          title: '编辑文章'
        }
      }
    ]
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - 博客后台管理系统`
  }
  
  // 检查是否需要登录
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    // 未登录且不是去登录页，则重定向到登录页
    next('/login')
  } else if (to.path === '/login' && token) {
    // 已登录且去登录页，则重定向到首页
    next('/')
  } else {
    // 其他情况正常放行
    next()
  }
})

export default router