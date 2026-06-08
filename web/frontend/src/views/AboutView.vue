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

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { marked } from 'marked'
import axios from 'axios'
import MainLayout from '@/components/layout/MainLayout.vue'
import ProfileCard from '@/components/sidebar/ProfileCard.vue'
import { useSiteInfoStore } from '@/stores/siteInfo'
import { storeToRefs } from 'pinia'
import { useDefaultAvatar } from '@/utils/defaults'

const siteInfoStore = useSiteInfoStore()
const defaultAvatarImg = useDefaultAvatar()
const { siteInfo } = storeToRefs(siteInfoStore)

// 渲染 Markdown 内容
const renderedContent = ref('')

onMounted(async () => {
  try {
    // 添加时间戳防止缓存
    const response = await axios.get(`/static/about.md?t=${new Date().getTime()}`)
    renderedContent.value = marked.parse(response.data) as string
  } catch (error) {
    console.error('Failed to load about.md:', error)
    renderedContent.value = '<h1>加载失败</h1><p>请检查网络或配置文件。</p>'
  }
})

// 打开链接
const openLink = (url: string) => {
  window.open(url, '_blank')
}

</script>

<style scoped>
.about-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.about-content {
  background: var(--color-background-soft);
  border-radius: 8px;
  box-shadow: 0 2px 8px var(--color-shadow);
  padding: 40px;
}

.about-content h1 {
  font-size: 32px;
  margin-bottom: 20px;
  color: var(--color-heading);
}

/* Markdown 样式 */
.markdown-body {
  color: var(--color-text);
  line-height: 1.6;
}

.markdown-body :deep(h1) {
  font-size: 36px;
  font-weight: bold;
  margin: 0 0 30px;
  padding-bottom: 15px;
  border-bottom: 2px solid var(--color-border);
  color: var(--color-heading);
}

.markdown-body :deep(h2) {
  font-size: 24px;
  margin: 30px 0 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border);
  color: var(--color-heading);
}

.markdown-body :deep(h3) {
  font-size: 20px;
  margin: 20px 0 10px;
  color: var(--color-heading);
}

.markdown-body :deep(p) {
  margin-bottom: 16px;
  color: var(--color-text);
}

.markdown-body :deep(ul) {
  padding-left: 20px;
  margin-bottom: 16px;
}

.markdown-body :deep(li) {
  margin-bottom: 8px;
  color: var(--color-text);
}

.markdown-body :deep(a) {
  color: var(--color-accent);
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.about-content h2 {
  font-size: 24px;
  margin: 30px 0 15px;
  color: var(--color-heading);
}

.about-content p {
  font-size: 16px;
  line-height: 1.6;
  color: var(--color-text);
  margin-bottom: 15px;
}

.about-content ul {
  padding-left: 20px;
  margin-bottom: 20px;
}

.about-content li {
  margin-bottom: 10px;
  color: var(--color-text);
}

.sidebar {
  width: 300px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sidebar-card {
  background: var(--color-background-soft);
  border-radius: 12px;
  box-shadow: 0 2px 8px var(--color-shadow);
  overflow: hidden;
  transition: all 0.3s ease;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-background-mute);
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--color-heading);
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
  border: 3px solid var(--color-border);
}

.status-indicator {
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid var(--color-background-soft);
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
  color: var(--color-heading);
}

.signature {
  font-size: 14px;
  color: var(--color-text-secondary);
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
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 4px;
}

.contact-item:hover {
  background-color: var(--color-background-mute);
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
  color: var(--color-accent);
}

.contact-item span {
  font-size: 14px;
  color: var(--color-text);
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
  color: var(--color-heading);
}

.artist {
  font-size: 14px;
  color: var(--color-text-secondary);
  margin: 0;
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

:global([data-theme="dark"]) .contact-item .iconfont {
  color: var(--color-accent) !important;
}
</style>