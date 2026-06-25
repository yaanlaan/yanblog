import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/articles',
      name: 'articles',
      component: () => import('../views/ArticleList.vue'),
    },
    {
      path: '/article/:id',
      name: 'article-detail',
      component: () => import('../views/ArticleDetail.vue'),
      props: true
    },
    {
      path: '/category/:id',
      name: 'category-articles',
      component: () => import('../views/ArticleList.vue'),
      props: true
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('../views/CategoryView.vue'),
    },
    {
      path: '/archive',
      name: 'archive',
      component: () => import('../views/ArchiveView.vue'),
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFound.vue')
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    // 如果浏览器支持且有保存的位置，使用保存的滚动位置（如浏览器后退）
    if (savedPosition) {
      return savedPosition
    }
    // 否则滚动到页面顶部
    return { top: 0 }
  }
})

// 全局前置守卫：参数校验 + 权限控制
router.beforeEach((to, _from, next) => {
  // 校验文章ID等动态参数是否为合法数字
  if (to.params.id && !/^\d+$/.test(to.params.id as string)) {
    // 非法ID，重定向到404或首页
    next({ name: 'not-found', params: [to.path] })
    return
  }

  next()
})

export default router
