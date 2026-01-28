<template>
  <div class="sidebar-card profile-card">
    <div class="profile-header">
      <img :src="siteInfo.author_avatar || '/assets/avatar.jpg'" alt="Avatar" class="avatar">
    </div>
    <div class="profile-content">
      <h3 class="name">{{ siteInfo.author_name }}</h3>
      <p class="bio">{{ siteInfo.author_bio }}</p>
      
      <div class="social-links">
        <a 
          v-for="(item, index) in siteInfo.socials" 
          :key="index"
          :href="item.url" 
          target="_blank" 
          class="social-item" 
          :title="item.name"
          :style="item.is_circle ? { '--hover-bg': item.color } : {}"
        >
          <i class="iconfont" :class="item.icon"></i>
        </a>
      </div>
      
      <div class="stats-row">
        <div class="stat-item">
          <span class="stat-val">{{ articleCount }}</span>
          <span class="stat-label">文章</span>
        </div>
        <div class="stat-item">
          <span class="stat-val">{{ categoryCount }}</span>
          <span class="stat-label">分类</span>
        </div>
      </div>
      
      <div class="quote-area">
        <Transition name="fade-slide" mode="out-in">
          <p class="quote-text" :key="currentQuote">"{{ currentQuote }}"</p>
        </Transition>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { articleApi, categoryApi } from '@/services/api'
import { useSiteInfoStore } from '@/stores/siteInfo'
import { storeToRefs } from 'pinia'

const siteInfoStore = useSiteInfoStore()
const { siteInfo } = storeToRefs(siteInfoStore)

const articleCount = ref(0)
const categoryCount = ref(0)
const currentQuote = ref('')
const currentQuoteIndex = ref(0)
let quoteInterval: number | null = null

const defaultQuotes = [
  "月亮想着我的心事，一只猫吃了我的奶酪",
  "草木山石，日月星辰",
  "雾霭山岚，风光雨霁",
]

const activeQuotes = computed(() => {
  return (siteInfo.value.quotes && siteInfo.value.quotes.length > 0) 
    ? siteInfo.value.quotes 
    : defaultQuotes
})

onMounted(async () => {
  // 初始化名言
  currentQuote.value = activeQuotes.value[0]
  
  // 启动轮播
  quoteInterval = setInterval(() => {
    currentQuoteIndex.value = (currentQuoteIndex.value + 1) % activeQuotes.value.length
    currentQuote.value = activeQuotes.value[currentQuoteIndex.value]
  }, 4000) as unknown as number

  try {
    const artRes = await articleApi.getArticles({ pagesize: 1, pagenum: 1 })
    if (artRes.data && artRes.data.total) {
      articleCount.value = artRes.data.total
    }
    
    const catRes = await categoryApi.getCategories({ pagesize: 1, pagenum: 1 })
    if (catRes.data && catRes.data.total) {
      categoryCount.value = catRes.data.total
    }
  } catch (e) {
    console.error('Failed to fetch stats', e)
  }
})

onUnmounted(() => {
  if (quoteInterval) {
    clearInterval(quoteInterval)
  }
})
</script>

<style scoped>
.profile-card {
  text-align: center;
  overflow: visible; /* Allow avatar to overlap */
  margin-top: 0; 
}

.profile-header {
  height: 100px;
  background: linear-gradient(135deg, #42b883 0%, #3aa876 100%);
  border-radius: 12px 12px 0 0;
  position: relative;
  margin-bottom: 40px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 4px solid var(--color-background-soft);
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  object-fit: cover;
  transition: transform 0.5s ease;
}

.avatar:hover {
  transform: translateX(-50%) rotate(360deg);
}

.profile-content {
  padding: 0 20px 20px;
}

.name {
  font-size: 20px;
  font-weight: bold;
  color: var(--color-heading);
  margin-bottom: 5px;
}

.bio {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin-bottom: 15px;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-bottom: 20px;
}

.social-item {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--color-background-mute);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all 0.3s;
}

.social-item:hover {
  background: var(--hover-bg, var(--color-accent));
  color: white;
  transform: translateY(-3px);
}

.stats-row {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
  padding: 15px 0;
  border-top: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
}

.stat-item {
  display: flex;
  flex-direction: column;
}

.stat-val {
  font-size: 18px;
  font-weight: bold;
  color: var(--color-heading);
}

.stat-label {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.quote-area {
  margin-top: 15px;
  padding: 10px;
  background-color: var(--color-background-mute);
  border-radius: 8px;
  border-left: 3px solid var(--color-accent);
  min-height: 60px; /* 预留高度防止跳动 */
  display: flex;
  align-items: center;
  justify-content: center;
}

.quote-text {
  font-size: 13px;
  color: var(--color-text-secondary);
  font-style: italic;
  margin: 0;
  line-height: 1.5;
}

/* 动画效果 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.5s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
