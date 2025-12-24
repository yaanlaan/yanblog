<template>
  <div class="hero-section">
    <!-- 上部分：两个大卡片 -->
    <div class="hero-top-grid">
      <!-- 左侧：介绍与技术栈 -->
      <div class="hero-card intro-card">
        <div class="intro-content">
          <h1 class="blog-title">草木山石<br>日月星辰</h1>
          <p class="blog-subtitle">yaan's blog</p>
        </div>
        <div class="tech-stack-visual">
          <!-- 装饰性图标背景 -->
          <div class="tech-icon icon-js">JS</div>
          <div class="tech-icon icon-vue">Vue</div>
          <div class="tech-icon icon-react">React</div>
          <div class="tech-icon icon-html">HTML5</div>
          <div class="tech-icon icon-css">CSS3</div>
          <div class="tech-icon icon-ts">TS</div>
        </div>
      </div>

      <!-- 右侧：欢迎图与推荐 -->
      <div class="hero-card welcome-card">
        <div class="welcome-overlay"></div>
        <div class="welcome-content">
          <h2 class="welcome-title">Welcome to<br>Yaan's Blog</h2>
          <button class="recommend-btn" @click="handleRandomVisit">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"/></svg>
            随便逛逛
          </button>
        </div>
      </div>
    </div>

    <!-- 下部分：三个导航卡片 -->
    <div class="hero-bottom-grid">
      <div class="nav-card blue-card">
        <div class="nav-content">
          <h3>黎明破晓</h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="nav-icon"><path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/></svg>
        </div>
      </div>
      <div class="nav-card orange-card">
        <div class="nav-content">
          <h3>日落黄昏</h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="nav-icon"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
        </div>
      </div>
      <div class="nav-card green-card">
        <div class="nav-content">
          <h3>午夜降临</h3>
          <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="nav-icon"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { articleApi } from '@/services/api'

const router = useRouter()

const handleRandomVisit = async () => {
  try {
    const res = await articleApi.getArticles({ pagesize: -1, pagenum: -1 })
    const { data, status } = res.data
    if (status === 200 && data && data.length > 0) {
      const randomIndex = Math.floor(Math.random() * data.length)
      const randomArticle = data[randomIndex]
      router.push({ name: 'article-detail', params: { id: randomArticle.ID } })
    }
  } catch (error) {
    console.error('Failed to fetch articles for random visit:', error)
  }
}
</script>

<style scoped>
.hero-section {
  margin-bottom: 30px;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

/* 上部分网格 */
.hero-top-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  height: 280px;
}

.hero-card {
  border-radius: 16px;
  position: relative;
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.hero-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.15);
}

/* 左侧介绍卡片 */
.intro-card {
  background: #1a1a1a;
  color: white;
  display: flex;
  align-items: center;
  padding: 40px;
  position: relative;
}

.intro-content {
  z-index: 2;
}

.blog-title {
  font-size: 36px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 10px;
  letter-spacing: 2px;
}

.blog-subtitle {
  font-size: 16px;
  opacity: 0.6;
  font-family: monospace;
}

/* 技术栈背景装饰 */
.tech-stack-visual {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 60%;
  overflow: hidden;
  z-index: 1;
  opacity: 0.8;
}

@keyframes float {
  0% { transform: translateY(0px) rotate(-15deg); }
  50% { transform: translateY(-15px) rotate(-15deg); }
  100% { transform: translateY(0px) rotate(-15deg); }
}

.tech-icon {
  position: absolute;
  width: 60px;
  height: 60px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
  color: white;
  box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  transform: rotate(-15deg);
  animation: float 6s ease-in-out infinite;
}

/* 模拟图标位置和颜色 */
.icon-js { background: #f7df1e; color: #323330; top: 20%; right: 40%; width: 70px; height: 70px; z-index: 2; animation-delay: 0s; }
.icon-vue { background: #42b883; top: 50%; right: 20%; width: 80px; height: 80px; z-index: 3; animation-delay: 1s; }
.icon-react { background: #61dafb; color: #20232a; top: 10%; right: 10%; animation-delay: 2s; }
.icon-html { background: #e34f26; bottom: 10%; right: 45%; animation-delay: 3s; }
.icon-css { background: #1572b6; bottom: 20%; right: 5%; animation-delay: 4s; }
.icon-ts { background: #3178c6; top: 40%; right: 55%; width: 50px; height: 50px; animation-delay: 1.5s; }

/* 右侧欢迎卡片 */
.welcome-card {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 30px;
  color: white;
  overflow: hidden;
}

.welcome-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: url('/src/assets/img/1412.jpg');
  background-size: cover;
  background-position: center;
  transition: all 0.5s ease;
  z-index: 0;
}

.welcome-card:hover::before {
  filter: blur(4px);
  transform: scale(1.1);
}

.welcome-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0) 60%);
  z-index: 1;
}

.welcome-content {
  position: relative;
  z-index: 2;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
}

.welcome-title {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
}

.recommend-btn {
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(5px);
  border: 1px solid rgba(255,255,255,0.4);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 5px;
}

.recommend-btn:hover {
  background: white;
  color: #333;
}

/* 下部分网格 */
.hero-bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 15px;
  height: 100px;
}

.nav-card {
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  color: white;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  overflow: hidden;
}

.nav-card:hover {
  transform: translateY(-2px);
  filter: brightness(1.1);
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  z-index: 2;
}

.nav-content h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.nav-icon {
  opacity: 0.8;
}

/* 颜色变体 */
.blue-card {
  background: linear-gradient(135deg, #2196f3, #00bcd4);
}

.orange-card {
  background: linear-gradient(135deg, #ff5722, #ff9800);
}

.green-card {
  background: linear-gradient(135deg, #4caf50, #8bc34a);
}

/* 响应式 */
@media (max-width: 768px) {
  .hero-top-grid {
    grid-template-columns: 1fr;
    height: auto;
  }
  
  .hero-card {
    height: 200px;
  }
  
  .hero-bottom-grid {
    grid-template-columns: 1fr;
    height: auto;
  }
  
  .nav-card {
    height: 80px;
  }
}
</style>