<template>
  <div class="about-page">
    <MainLayout>
      <template #main>
        <div class="about-content">
          <!-- Markdown 内容渲染区域 -->
          <div class="markdown-body" v-html="renderedContent"></div>
        </div>
      </template>
      <template #sidebar>
          <!-- 个人卡片 -->
          <ProfileCard />

          <!-- 联系方式卡片 -->
          <div class="sidebar-card contact-card" ref="contactCard" v-if="siteInfo.contacts?.show">
            <div class="card-header">
              <h3>联系方式</h3>
            </div>
            <div class="card-content">
              <ul class="contact-list">
                <li v-for="(item, index) in siteInfo.contacts.items" :key="index" class="contact-item" @click="openLink(item.url)">
                  <i class="iconfont" :class="item.icon" :style="{ color: item.color }"></i>
                  <span>{{ item.name }}</span>
                </li>
              </ul>
            </div>
          </div>

          <!-- 音乐播放器卡片 -->
          <div class="sidebar-card music-card" v-if="siteInfo.music_player?.show && siteInfo.music_player?.url">
            <div class="card-header">
              <h3>音乐播放器</h3>
            </div>
            <div class="card-content">
              <div class="music-player">
                <!-- 音乐播放器 -->
                <iframe frameborder="no" border="0" marginwidth="0" marginheight="0" width="100%" height="86"
                  :src="siteInfo.music_player.url">
                </iframe>
              </div>
            </div>
          </div>
      </template>
    </MainLayout>

    <!-- 微信二维码弹窗 -->
    <div v-if="showWechatQR" class="qr-modal" @click="showWechatQR = false">
      <div class="qr-content" @click.stop>
        <h3>微信二维码</h3>
        <img :src="siteInfo.contacts?.wechat_qr || '/assets/default-qr.png'" alt="微信二维码" @error="handleQRError">
        <p>扫描二维码添加微信</p>
        <button class="close-btn" @click="showWechatQR = false">关闭</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { marked } from 'marked'
import MainLayout from '@/components/layout/MainLayout.vue'
import ProfileCard from '@/components/sidebar/ProfileCard.vue'
import aboutContent from '@/assets/about/about.md?raw'
import { useSiteInfoStore } from '@/stores/siteInfo'
import { storeToRefs } from 'pinia'

const siteInfoStore = useSiteInfoStore()
const { siteInfo } = storeToRefs(siteInfoStore)

// 渲染 Markdown 内容
const renderedContent = ref('')

onMounted(() => {
  renderedContent.value = marked.parse(aboutContent) as string
})

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

/* Markdown 样式 */
.markdown-body {
  color: #333;
  line-height: 1.6;
}

.markdown-body :deep(h2) {
  font-size: 24px;
  margin: 30px 0 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
  color: #333;
}

.markdown-body :deep(h3) {
  font-size: 20px;
  margin: 20px 0 10px;
  color: #444;
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
  color: #666;
}

.markdown-body :deep(ul) {
  padding-left: 20px;
  margin-bottom: 16px;
}

.markdown-body :deep(li) {
  margin-bottom: 8px;
  color: #666;
}

.markdown-body :deep(a) {
  color: #42b883;
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
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
  color: #42b883;
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
  background: #42b883;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 15px;
}

.close-btn:hover {
  background: #3aa876;
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