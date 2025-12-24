<template>
  <div class="sidebar-card profile-card">
    <div class="profile-header">
      <img src="@/assets/avatar.jpg" alt="Avatar" class="avatar">
    </div>
    <div class="profile-content">
      <h3 class="name">Yaan</h3>
      <p class="bio">欢迎光临 (￣▽￣)~*</p>
      
      <div class="social-links">
        <a href="https://github.com/yaanlaan" target="_blank" class="social-item" title="GitHub">
          <i class="iconfont icon-github-fill"></i>
        </a>
        <a href="https://space.bilibili.com/3461574693360526" target="_blank" class="social-item" title="Bilibili">
          <i class="iconfont icon-bilibili-line"></i>
        </a>
        <a href="mailto:yanxia2425@foxmail.com" class="social-item" title="Email">
          <i class="iconfont icon-email"></i>
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
        <div class="stat-item">
          <span class="stat-val">{{ wordCount }}</span>
          <span class="stat-label">字数</span>
        </div>
      </div>
      
      <div class="quote-area">
        <p class="quote-text">"{{ currentQuote }}"</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { articleApi, categoryApi } from '@/services/api'

const articleCount = ref(0)
const categoryCount = ref(0)
const wordCount = ref('50k')
const currentQuote = ref('Talk is cheap. Show me the code.')

const quotes = [
  "月亮想着我的心事，一只猫吃了我的奶酪",
  "草木山石，日月星辰",
  "雾霭山岚，风光雨霁"
]

onMounted(async () => {
  currentQuote.value = quotes[Math.floor(Math.random() * quotes.length)]
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
</script>

<style scoped>
.profile-card {
  text-align: center;
  overflow: visible; /* Allow avatar to overlap */
  margin-top: 30px; /* Space for overlapping avatar */
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
  border: 4px solid white;
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
  color: #333;
  margin-bottom: 5px;
}

.bio {
  font-size: 14px;
  color: #666;
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
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  text-decoration: none;
  transition: all 0.3s;
}

.social-item:hover {
  background: #42b883;
  color: white;
  transform: translateY(-3px);
}

.stats-row {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
  padding: 15px 0;
  border-top: 1px solid #eee;
  border-bottom: 1px solid #eee;
}

.stat-item {
  display: flex;
  flex-direction: column;
}

.stat-val {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 12px;
  color: #999;
}

.quote-area {
  margin-top: 15px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 8px;
  border-left: 3px solid #42b883;
}

.quote-text {
  font-size: 13px;
  color: #666;
  font-style: italic;
  margin: 0;
  line-height: 1.5;
}
</style>
