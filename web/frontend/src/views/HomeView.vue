<template>
  <div class="home-page">
    <!-- 全局加载动画 -->
    <div v-if="pageLoading" class="loading-overlay">
      <div class="arc"></div>
      <h1 class="loading-text"><span>LOADING</span></h1>
    </div>
    
    <!-- 主页内容 -->
    <MainLayout v-else>
      <template #main>
        <HeroSection />
        <TopArticles :articles="topArticles" :loading="topLoading" />
        <LatestArticles :articles="latestArticles" :loading="latestLoading" />
      </template>
      <template #sidebar>
        <Sidebar />
      </template>
    </MainLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { articleApi } from '@/services/api'
import MainLayout from '@/components/layout/MainLayout.vue'
import HeroSection from '@/components/home/HeroSection.vue'
import TopArticles from '@/components/home/TopArticles.vue'
import LatestArticles from '@/components/home/LatestArticles.vue'
import Sidebar from '@/components/Sidebar.vue'

// 类型定义
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  top: number
  createdAt: string
  updatedAt: string
}

// 响应式数据
const topArticles = ref<Article[]>([])
const latestArticles = ref<Article[]>([])
const topLoading = ref(false)
const latestLoading = ref(false)
const pageLoading = ref(true) // 页面全局加载状态

// 获取置顶文章
const getTopArticles = async () => {
  topLoading.value = true
  try {
    const response = await articleApi.getTopArticles({
      num: -1 // 获取所有置顶文章
    })
    
    const { data } = response.data
    topArticles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      top: item.top,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
    // 按置顶等级排序，数字越小等级越高
    .sort((a, b) => a.top - b.top)
  } catch (error) {
    console.error('获取置顶文章失败:', error)
  } finally {
    topLoading.value = false
    checkPageLoading()
  }
}

// 获取最新文章
const getLatestArticles = async () => {
  latestLoading.value = true
  try {
    const response = await articleApi.getArticles({
      pagesize: 10, // 获取更多文章，然后在组件中限制显示数量
      pagenum: 1
    })
    
    const { data } = response.data
    // 过滤掉置顶文章，只显示非置顶文章
    const nonTopArticles = data
      .filter((item: any) => !item.top || item.top === 0)
      .map((item: any) => ({
        id: item.ID,
        title: item.title,
        categoryId: item.cid,
        categoryName: item.Category?.name || '未分类',
        desc: item.desc,
        content: item.content,
        img: item.img,
        top: item.top,
        createdAt: item.CreatedAt || item.created_at,
        updatedAt: item.UpdatedAt || item.updated_at
      }))
      // 按创建时间倒序排列
      .sort((a: Article, b: Article) => {
        const dateA = new Date(a.createdAt).getTime()
        const dateB = new Date(b.createdAt).getTime()
        return dateB - dateA
      })
    
    latestArticles.value = nonTopArticles
  } catch (error) {
    console.error('获取最新文章失败:', error)
  } finally {
    latestLoading.value = false
    checkPageLoading()
  }
}

// 检查页面加载状态
const checkPageLoading = () => {
  // 当所有数据加载完成时，隐藏全局加载动画
  if (!topLoading.value && !latestLoading.value) {
    // 添加一点延迟以确保用户体验
    setTimeout(() => {
      pageLoading.value = false
    }, 300)
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getTopArticles()
  getLatestArticles()
})
</script>

<style scoped>
.home-page {
  width: 100%;
  min-height: calc(100vh - 200px);
  position: relative;
}

.loading-overlay {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  background-color: var(--color-background);
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
}

.arc {
  position: absolute;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  border-top: 3px solid var(--color-accent);
  border-left: 1px solid transparent;
  border-right: 1px solid transparent;
  animation: rt 2s infinite linear;
}

.arc::before {
  content: "";
  position: absolute;
  width: 70px;
  height: 70px;
  border-radius: 50%;
  border-top: 2px solid #c85540;
  border-left: 1px solid transparent;
  border-right: 1px solid transparent;
  animation: rt 4s infinite linear reverse;
  margin: auto;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}

.arc::after {
  content: "";
  position: absolute;
  width: 0;
  height: 0;
  border-radius: 50%;
  border-top: initial;
  border-left: initial;
  border-right: initial;
  animation: cw 1s infinite;
  margin: auto;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  background: var(--color-accent);
  opacity: 0.5;
}

.loading-text {
  position: absolute;
  height: 40px;
  margin: auto;
  top: 200px;
  left: 0;
  right: 0;
  bottom: 0;
  text-transform: uppercase;
  text-align: center;
  letter-spacing: 0.1em;
  font-size: 14px;
  font-weight: lighter;
  color: var(--color-accent);
}

.loading-text span {
  display: none;
}

.loading-text::after {
  content: "";
  animation: txt 5s infinite;
}

@keyframes rt {
  100% {
    transform: rotate(360deg);
  }
}

@keyframes cw {
  0% {
    width: 0;
    height: 0;
  }

  75% {
    width: 40px;
    height: 40px;
  }

  100% {
    width: 0;
    height: 0;
  }
}

@keyframes txt {
  0% {
    content: "LOADING.";
  }

  50% {
    content: "LOADING..";
  }

  100% {
    content: "LOADING...";
  }
}
</style>