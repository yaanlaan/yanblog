import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ArticleList from '../views/ArticleList.vue'
import ArticleDetail from '../views/ArticleDetail.vue'
import CategoryView from '../views/CategoryView.vue'

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
      component: ArticleList,
    },
    {
      path: '/article/:id',
      name: 'article-detail',
      component: ArticleDetail,
      props: true
    },
    {
      path: '/category/:id',
      name: 'category-articles',
      component: ArticleList,
      props: true
    },
    {
      path: '/categories',
      name: 'categories',
      component: CategoryView,
    },
    {
      path: '/archive',
      name: 'archive',
      component: () => import('../views/ArchiveView.vue'),
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

export default router