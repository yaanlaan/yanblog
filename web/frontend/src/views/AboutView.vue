<template>
  <div class="about-page">
    <MainLayout>
      <template #main>
        <div class="about-content">
          <h1>关于我们</h1>

          <div class="section">
            <h2>项目简介</h2>
            <p>这是一个使用 Vue 3、Go 语言和现代 Web 技术构建的全栈博客系统。项目采用前后端分离架构，后端使用 Gin 框架提供 RESTful API，前端使用 Vue 3 构建响应式界面。</p>
          </div>

          <div class="section">
            <h2>技术栈</h2>
            <div class="tech-stack">
              <div class="tech-category">
                <h3>后端技术</h3>
                <ul>
                  <li>Go 语言</li>
                  <li>Gin Web 框架</li>
                  <li>GORM (数据库ORM)</li>
                  <li>JWT (身份认证)</li>
                  <li>MySQL (数据库)</li>
                </ul>
              </div>

              <div class="tech-category">
                <h3>前端技术</h3>
                <ul>
                  <li>Vue 3 (Composition API)</li>
                  <li>TypeScript</li>
                  <li>Vite (构建工具)</li>
                  <li>Element Plus (UI组件库)</li>
                  <li>Vue Router (路由管理)</li>
                </ul>
              </div>

              <div class="tech-category">
                <h3>部署与工具</h3>
                <ul>
                  <li>Git (版本控制)</li>
                  <li>RESTful API 设计</li>
                  <li>响应式设计</li>
                </ul>
              </div>
            </div>
          </div>

          <div class="section">
            <h2>功能特性</h2>
            <div class="features">
              <div class="feature-card">
                <i class="iconfont icon-article"></i>
                <h3>文章管理</h3>
                <p>完整的文章发布、编辑、删除功能，支持Markdown格式</p>
              </div>

              <div class="feature-card">
                <i class="iconfont icon-category"></i>
                <h3>分类系统</h3>
                <p>灵活的文章分类管理，便于内容组织</p>
              </div>

              <div class="feature-card">
                <i class="iconfont icon-user"></i>
                <h3>用户权限</h3>
                <p>基于JWT的用户认证和权限控制系统</p>
              </div>

              <div class="feature-card">
                <i class="iconfont icon-upload"></i>
                <h3>文件上传</h3>
                <p>支持图片等文件上传到云存储服务</p>
              </div>

              <div class="feature-card">
                <i class="iconfont icon-weather"></i>
                <h3>天气信息</h3>
                <p>集成第三方天气API，实时显示天气状况</p>
              </div>

              <div class="feature-card">
                <i class="iconfont icon-monitor"></i>
                <h3>系统监控</h3>
                <p>实时监控服务器状态，包括内存、CPU使用率等</p>
              </div>
            </div>
          </div>

          <div class="section">
            <h2>视频插入</h2>
            <div class="demo-content">
              <h3>视频示例（博客配图可以上传七牛云并复制外链添加，视频就自己从bilibili引外链吧，下面是一个例子）</h3>
              <div class="video-container">
                <iframe
                  src="//player.bilibili.com/player.html?isOutside=true&aid=638975717&bvid=BV1HY4y1r731&cid=717992692&page=1&high_quality=1&danmaku=0"
                  allowfullscreen="allowfullscreen" width="100%" height="500" scrolling="no" frameborder="0"
                  sandbox="allow-top-navigation allow-same-origin allow-forms allow-scripts"></iframe>
              </div>
              <p class="video-description">莫扎特 D大调嬉游曲 K. 136丨Pekka Kuusisto丨挪威室内乐团</p>
            </div>
          </div>
          <div class="section"> 
            <h2>参考说明</h2>
          <p>本博客主要参考了ginblog项目的后端，和知乎恢锦的博客元素设计</p>
          </div>
        </div>
      </template>
      <template #sidebar>
        <div class="sidebar right-sidebar">
          <!-- 个人卡片 -->
          <div class="sidebar-card profile-card">
            <div class="card-header">
              <h3>个人资料</h3>
            </div>
            <div class="card-content">
              <div class="profile-wrapper">
                <div class="avatar-wrapper">
                  <img src="@/assets/avatar.jpg" alt="头像" class="avatar" @error="handleAvatarError" />
                  <!-- 状态指示器 -->
                  <div class="status-indicator online"></div>
                </div>
                <div class="profile-info">
                  <h4 class="username">yaanlaan</h4>
                  <p class="signature">月亮想着我的心事，猫儿吃了我的奶酪</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 联系方式卡片 -->
          <div class="sidebar-card contact-card" ref="contactCard">
            <div class="card-header">
              <h3>联系方式</h3>
            </div>
            <div class="card-content">
              <ul class="contact-list">
                <li class="contact-item" @click="openLink('mailto:example@email.com')">
                  <i class="iconfont icon-email"></i>
                  <span>example@email.com</span>
                </li>
                <li class="contact-item" @click="openQQChat('123456789')">
                  <i class="iconfont icon-QQ"></i>
                  <span>123456789</span>
                </li>
                <li class="contact-item" @click="showQRCode('wechat')">
                  <i class="iconfont icon-wechat-fill"></i>
                  <span>wechat_id</span>
                </li>
                <li class="contact-item" @click="openLink('https://github.com/username')">
                  <i class="iconfont icon-github-fill"></i>
                  <span>github_username</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- 音乐播放器卡片 -->
          <div class="sidebar-card music-card">
            <div class="card-header">
              <h3>音乐播放器</h3>
            </div>
            <div class="card-content">
              <div class="music-player">
                <!-- 音乐播放器 -->
                <iframe frameborder="no" border="0" marginwidth="0" marginheight="0" width="100%" height="86"
                  src="//music.163.com/outchain/player?type=2&id=555503430&auto=0&height=66">
                </iframe>
              </div>
            </div>
          </div>
        </div>
      </template>
    </MainLayout>

    <!-- 微信二维码弹窗 -->
    <div v-if="showWechatQR" class="qr-modal" @click="showWechatQR = false">
      <div class="qr-content" @click.stop>
        <h3>微信二维码</h3>
        <img src="@/assets/wechat-qr.png" alt="微信二维码" @error="handleQRError">
        <p>扫描二维码添加微信</p>
        <button class="close-btn" @click="showWechatQR = false">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MainLayout from '@/components/layout/MainLayout.vue'

// 微信二维码显示状态
const showWechatQR = ref(false)

// 联系方式卡片引用
const contactCard = ref<HTMLElement | null>(null)

// 处理头像加载错误
const handleAvatarError = (e: Event) => {
  const target = e.target as HTMLImageElement
  target.src = '@/assets/default-avatar.jpg'
}

// 处理二维码图片加载错误
const handleQRError = (e: Event) => {
  const target = e.target as HTMLImageElement
  target.src = '@/assets/default-qr.png'
}

// 打开链接
const openLink = (url: string) => {
  window.open(url, '_blank')
}

// 打开QQ聊天
const openQQChat = (qqNumber: string) => {
  window.open(`http://wpa.qq.com/msgrd?v=3&uin=${qqNumber}&site=qq&menu=yes`, '_blank')
}

// 显示二维码
const showQRCode = (type: string) => {
  if (type === 'wechat') {
    showWechatQR.value = true
  }
}

// 组件挂载时添加滚动监听
// 注意：由于这是一个局部组件，我们不需要在这里添加滚动监听
// 吸附效果更适合在整个应用级别实现
</script>

<style scoped>
.about-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.about-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.about-content h1 {
  font-size: 32px;
  margin-bottom: 20px;
  color: #333;
}

.about-content h2 {
  font-size: 24px;
  margin: 30px 0 15px;
  color: #333;
}

.about-content p {
  font-size: 16px;
  line-height: 1.6;
  color: #666;
  margin-bottom: 15px;
}

.about-content ul {
  padding-left: 20px;
  margin-bottom: 20px;
}

.about-content li {
  margin-bottom: 10px;
  color: #666;
}

.sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  background: #f8f9fa;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.card-content {
  padding: 20px;
}

/* 个人卡片样式 */
.profile-wrapper {
  text-align: center;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
  margin-bottom: 15px;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #f0f0f0;
}

.status-indicator {
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid white;
  bottom: 5px;
  right: 5px;
}

.status-indicator.online {
  background-color: #4caf50;
}

.status-indicator.offline {
  background-color: #9e9e9e;
}

.status-indicator.busy {
  background-color: #f44336;
}

.status-indicator.away {
  background-color: #ff9800;
}

.username {
  margin: 10px 0 5px;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.signature {
  font-size: 14px;
  color: #666;
  margin: 0;
}

/* 联系方式卡片样式 */
.contact-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.contact-item {
  display: flex;
  align-items: center;
  padding: 12px 15px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.contact-item:hover {
  background-color: #f5f7fa;
  transform: translateX(5px);
}

.contact-item:last-child {
  border-bottom: none;
}

.contact-item i {
  font-size: 18px;
  margin-right: 10px;
  width: 24px;
  text-align: center;
  color: #409eff;
}

.contact-item span {
  font-size: 14px;
  color: #333;
}

/* 音乐播放器卡片样式 */
.music-player {
  text-align: center;
}

.music-info {
  margin-top: 15px;
}

.song-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 5px;
  color: #333;
}

.artist {
  font-size: 14px;
  color: #666;
  margin: 0;
}

/* 二维码弹窗样式 */
.qr-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.qr-content {
  background: white;
  padding: 30px;
  border-radius: 8px;
  text-align: center;
  max-width: 300px;
  width: 90%;
}

.qr-content h3 {
  margin-top: 0;
  color: #333;
}

.qr-content img {
  width: 200px;
  height: 200px;
  object-fit: contain;
  margin: 15px 0;
}

.qr-content p {
  margin: 10px 0;
  color: #666;
}

.close-btn {
  background: #409eff;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 15px;
}

.close-btn:hover {
  background: #337ecc;
}

@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
  }

  .sidebar-card {
    flex: 1 1 calc(50% - 10px);
    min-width: 200px;
  }

  .about-content {
    padding: 20px;
  }

  .about-content h1 {
    font-size: 28px;
  }

  .about-content h2 {
    font-size: 22px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    flex-direction: column;
  }

  .sidebar-card {
    flex: 1 1 100%;
  }
}
</style>