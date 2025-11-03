<template>
  <div class="hero-section">
    <div class="hero-content">
      <h1 class="hero-title">欢迎来到我的博客</h1>
      <p class="hero-description">
        这里记录了我的技术学习、生活感悟和各种经验分享。
        希望我的文章能对你有所帮助。
      </p>
      <router-link to="/articles" class="hero-button">
        浏览文章
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
// Hero部分组件
</script>

<style scoped>
.hero-section {
  background: linear-gradient(275deg, #80acdfff 0%, #2a5eceff 100%);
  color: white;
  padding: 80px 0;
  margin-bottom: 40px;
  border-radius: 8px;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
  padding: 0 20px;
}

.hero-title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 20px;
}

.hero-description {
  font-size: 20px;
  margin-bottom: 30px;
  opacity: 0.9;
}

.hero-button {
  display: inline-block;
  background: white;
  color: #667eea;
  padding: 15px 30px;
  border-radius: 50px;
  text-decoration: none;
  font-weight: 600;
  font-size: 18px;
  transition: transform 0.3s, box-shadow 0.3s;
}

.hero-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 36px;
  }
  
  .hero-description {
    font-size: 18px;
  }
  
  .hero-section {
    padding: 60px 0;
  }
}
</style>