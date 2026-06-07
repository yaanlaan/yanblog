<template>
  <div class="home-page">
    <MainLayout>
      <template #main>
        <HeroSection />
        
        <div class="articles-section">
          <TopArticles :articles="topArticles" :loading="topLoading" />
          
          <LatestArticles :articles="latestArticles" :loading="latestLoading" />
        </div>
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
import ArticleSkeleton from '@/components/skeleton/ArticleSkeleton.vue'

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

const topArticles = ref<Article[]>([])
const latestArticles = ref<Article[]>([])
const topLoading = ref(false)
const latestLoading = ref(false)

const getTopArticles = async () => {
  topLoading.value = true
  try {
    const response = await articleApi.getTopArticles({ num: -1 })
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
    })).sort((a, b) => a.top - b.top)
  } catch (error) {
    console.error('获取置顶文章失败:', error)
  } finally {
    topLoading.value = false
  }
}

const getLatestArticles = async () => {
  latestLoading.value = true
  try {
    const response = await articleApi.getArticles({ pagesize: 10, pagenum: 1 })
    const { data } = response.data
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
  }
}

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

.articles-section {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.skeleton-wrapper {
  padding: 20px 0;
}
</style>