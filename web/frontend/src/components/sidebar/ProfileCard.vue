<template>
  <div class="sidebar-card profile-card">
    <div class="profile-header">
      <div 
        class="avatar-wrapper" 
        @mouseenter="refreshHoverQuote"
        @mouseleave="showBubble = false"
      >
        <img :src="siteInfo.author_avatar || '/assets/avatar.jpg'" alt="Avatar" class="avatar">
        <!-- 状态指示点 -->
        <div class="status-dot"></div>
        <!-- 气泡对话框 -->
        <Transition name="pop">
          <div class="speech-bubble" v-if="showBubble">
            {{ hoverQuote }}
          </div>
        </Transition>
      </div>
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
const hoverQuote = ref('')
const showBubble = ref(false)
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

const refreshHoverQuote = () => {
  if (activeQuotes.value.length > 0) {
    const randomIndex = Math.floor(Math.random() * activeQuotes.value.length)
    hoverQuote.value = activeQuotes.value[randomIndex]
    showBubble.value = true
  }
}

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

.avatar-wrapper {
  width: 80px;
  height: 80px;
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  cursor: pointer;
}

.avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 4px solid var(--color-background-soft);
  object-fit: cover;
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
  background-color: var(--color-background-soft);
  position: relative;
  z-index: 1;
}

.avatar-wrapper:hover .avatar {
  transform: scale(0.9);
}

.status-dot {
  position: absolute;
  bottom: 5px;
  right: 5px;
  width: 16px;
  height: 16px;
  background-color: #2ecc71; /* Online green */
  border: 3px solid var(--color-background-soft);
  border-radius: 50%;
  z-index: 2;
  box-shadow: 0 0 0 1px rgba(0,0,0,0.05);
}

.speech-bubble {
  position: absolute;
  bottom: 110%; /* Above avatar */
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--color-background-soft);
  color: var(--color-text);
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  width: max-content;
  max-width: 220px;
  box-shadow: 0 4px 12px var(--color-shadow);
  border: 1px solid var(--color-border);
  z-index: 20;
  text-align: center;
  pointer-events: none; /* Let clicks pass through if needed */
}

/* Bubble Arrow */
.speech-bubble::before {
  content: '';
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  width: 10px;
  height: 10px;
  background-color: var(--color-background-soft);
  border-right: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
  transform: translateX(-50%) rotate(45deg);
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

/* 气泡弹出动画 */
.pop-enter-active,
.pop-leave-active {
  transition: all 0.3s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.pop-enter-from,
.pop-leave-to {
  opacity: 0;
  transform: translate(-50%, 10px) scale(0.8);
}
</style>
