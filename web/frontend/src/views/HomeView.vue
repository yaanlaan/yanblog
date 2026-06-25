<template>
  <div class="home-page">
    <MainLayout>
      <template #main>
        <HeroSection />

        <div class="articles-section">
          <TopArticles :articles="topArticles" :loading="topLoading" />

          <LatestArticles
            :articles="latestArticles"
            :loading="latestLoading"
            :loading-more="loadingMore"
            :has-more="hasMore"
            @load-more="loadMoreArticles"
          />
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
import { mapArticle } from '@/utils/dataMapper'
import { PAGINATION } from '@/utils/constants'
import type { Article } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import HeroSection from '@/components/home/HeroSection.vue'
import TopArticles from '@/components/home/TopArticles.vue'
import LatestArticles from '@/components/home/LatestArticles.vue'
import Sidebar from '@/components/Sidebar.vue'

const PAGE_SIZE = PAGINATION.HOME_PAGE_SIZE

const topArticles = ref<Article[]>([])
const latestArticles = ref<Article[]>([])
const topLoading = ref(false)
const latestLoading = ref(false)
const loadingMore = ref(false)
const currentPage = ref(1)
const hasMore = ref(true)

const getTopArticles = async () => {
  topLoading.value = true
  try {
    const response = await articleApi.getTopArticles({ num: -1 })
    const { data } = response.data
    topArticles.value = data.map(mapArticle).sort((a, b) => a.top - b.top)
  } catch (error) {
    console.error('获取置顶文章失败:', error)
  } finally {
    topLoading.value = false
  }
}

const getLatestArticles = async () => {
  latestLoading.value = true
  try {
    const response = await articleApi.getArticles({
      pagesize: PAGE_SIZE,
      pagenum: 1,
      excludeTop: true
    })
    const { data, total } = response.data
    latestArticles.value = data.map(mapArticle)
    hasMore.value = latestArticles.value.length < total
  } catch (error) {
    console.error('获取最新文章失败:', error)
  } finally {
    latestLoading.value = false
  }
}

const loadMoreArticles = async () => {
  if (loadingMore.value || !hasMore.value) return

  loadingMore.value = true
  try {
    const nextPage = currentPage.value + 1
    const response = await articleApi.getArticles({
      pagesize: PAGE_SIZE,
      pagenum: nextPage,
      excludeTop: true
    })
    const { data, total } = response.data
    const newArticles = data.map(mapArticle)
    latestArticles.value.push(...newArticles)
    currentPage.value = nextPage
    hasMore.value = latestArticles.value.length < total
  } catch (error) {
    console.error('加载更多文章失败:', error)
  } finally {
    loadingMore.value = false
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