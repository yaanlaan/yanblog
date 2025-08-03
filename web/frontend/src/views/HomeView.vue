<template>
  <div class="home-page">
    <MainLayout>
      <template #main>
        <HeroSection />
        <LatestArticles :articles="articles" :loading="loading" />
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
  createdAt: string
  updatedAt: string
}

// 响应式数据
const articles = ref<Article[]>([])
const loading = ref(false)

// 获取最新文章
const getLatestArticles = async () => {
  loading.value = true
  try {
    const response = await articleApi.getArticles({
      pagesize: 6,
      pagenum: 1
    })
    
    const { data } = response.data
    articles.value = data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
  } catch (error) {
    console.error('获取最新文章失败:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getLatestArticles()
})
</script>

<style scoped>
.home-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}
</style>