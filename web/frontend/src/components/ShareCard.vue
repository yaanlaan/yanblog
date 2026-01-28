<template>
  <div v-if="visible" class="share-modal-overlay" @click="handleOverlayClick">
    <div class="share-modal" @click.stop>
      <div class="modal-header">
        <h3>分享卡片</h3>
        <button class="close-btn" @click="close">×</button>
      </div>
      
      <div class="modal-body">
        <!-- 预览区域 (将被转为图片) -->
        <div class="card-preview-container">
          <div 
            ref="cardRef" 
            class="share-card" 
            :class="theme"
            :style="{ backgroundColor: currentTheme.bg, color: currentTheme.text }"
          >
            <div class="card-content">
              <div class="quote-icon">❝</div>
              <div class="text-content">
                {{ content }}
              </div>
              <div class="author-info">
                <div class="left">
                    <!-- 假设默认头像，实际可替换 -->
                  <img src="@/assets/default-avatar.jpg" alt="avatar" class="avatar" onerror="this.src='/favicon.ico'" />
                  <span class="nickname">YanBlog</span>
                </div>
                <div class="right logo">
                    YanBlog
                </div>
              </div>
            </div>
            
            <div class="card-footer" :style="{ backgroundColor: currentTheme.footerBg }">
              <div class="footer-info">
                <h4 class="article-title">{{ title }}</h4>
                <p class="site-desc">分享自 YanBlog</p>
              </div>
              <div class="qrcode-wrapper">
                <img :src="qrcodeUrl" class="qrcode-img" />
                <span class="scan-tip">扫码阅读全文</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 控制区域 -->
        <div class="controls">
          <div class="theme-selector">
            <span class="label">配色风格：</span>
            <button 
              v-for="t in themes" 
              :key="t.name"
              class="theme-btn"
              :class="{ active: theme === t.name }"
              :style="{ background: t.preview }"
              @click="theme = t.name"
              :title="t.label"
            ></button>
          </div>
          
          <button class="download-btn" @click="generateImage" :disabled="generating">
            <span v-if="generating">生成中...</span>
            <span v-else>下载分享海报</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import QRCode from 'qrcode'
import html2canvas from 'html2canvas'

const props = defineProps<{
  visible: boolean
  title: string
  content: string // 选中的文字或摘要
  url: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const cardRef = ref<HTMLElement | null>(null)
const qrcodeUrl = ref('')
const generating = ref(false)
const theme = ref('default')

// 主题配置
const themes = [
  { name: 'default', label: '知乎蓝', preview: 'linear-gradient(135deg, #fff 50%, #f6f6f6 50%)', bg: '#ffffff', text: '#121212', footerBg: '#f6f6f6' },
  { name: 'dark', label: '暗夜黑', preview: '#1e1e1e', bg: '#1e1e1e', text: '#e0e0e0', footerBg: '#2d2d2d' },
  { name: 'parchment', label: '羊皮纸', preview: '#fbf7ed', bg: '#fbf7ed', text: '#5c4b37', footerBg: '#f2e8d5' },
  { name: 'mint', label: '清新绿', preview: '#e8f5e9', bg: '#ffffff', text: '#2e4c34', footerBg: '#e8f5e9' },
]

const currentTheme = computed(() => {
  return themes.find(t => t.name === theme.value) || themes[0]
})

// 生成二维码
watch(() => props.visible, async (val) => {
  if (val && props.url) {
    try {
      qrcodeUrl.value = await QRCode.toDataURL(props.url, {
        margin: 1,
        width: 100,
        color: {
          dark: '#000000',
          light: '#ffffff00' // 透明背景
        }
      })
    } catch (err) {
      console.error('Generate QRCode failed', err)
    }
  }
})

const close = () => {
  emit('close')
}

const handleOverlayClick = () => {
  close()
}

const generateImage = async () => {
  if (!cardRef.value) return
  generating.value = true
  
  try {
    // 等待图片加载（如果有头像等）
    await nextTick()
    
    const canvas = await html2canvas(cardRef.value, {
      useCORS: true,
      scale: 2, // 提高清晰度
      backgroundColor: currentTheme.value.bg
    })
    
    // 下载图片
    const link = document.createElement('a')
    link.download = `share_${Date.now()}.png`
    link.href = canvas.toDataURL('image/png')
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (err) {
    console.error('Generate image failed', err)
    alert('生成海报失败，请重试')
  } finally {
    generating.value = false
  }
}
</script>

<style scoped>
.share-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.6);
  z-index: 2000;
  display: flex;
  justify-content: center;
  align-items: center;
  backdrop-filter: blur(4px);
}

.share-modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 15px 20px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

/* 卡片样式 - 也是生成图片的源 */
.card-preview-container {
  width: 100%;
  max-width: 375px; /* 模拟手机宽度 */
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
}

.share-card {
  width: 100%;
  min-height: 480px;
  background: white;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
}

.card-content {
  flex: 1;
  padding: 30px 24px;
  display: flex;
  flex-direction: column;
}

.quote-icon {
  font-size: 48px;
  line-height: 1;
  color: rgba(0,0,0,0.1);
  font-family: Georgia, serif;
  margin-bottom: 10px;
}

/* 主题适配引用引号颜色 */
.share-card.dark .quote-icon { color: rgba(255,255,255,0.1); }
.share-card.parchment .quote-icon { color: rgba(92, 75, 55, 0.15); }
.share-card.mint .quote-icon { color: rgba(46, 76, 52, 0.15); }

.text-content {
  font-size: 18px;
  line-height: 1.6;
  font-weight: 400;
  margin-bottom: 30px;
  flex: 1;
  /* 限制最大高度，防止截图太长 */
  max-height: 400px;
  overflow-y: hidden; 
  position: relative;
  text-align: justify;
}

/* 如果内容过长，可以使用 JS 截断或者 CSS 遮罩（截图时遮罩可能不好看，先只做溢出隐藏） */

.author-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: auto;
  padding-top: 20px;
}

.author-info .left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.nickname {
  font-size: 14px;
  font-weight: 500;
  opacity: 0.8;
}

.author-info .right {
  font-weight: bold;
  font-family: 'Times New Roman', serif;
  opacity: 0.3;
  letter-spacing: 1px;
}

.card-footer {
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-info {
  flex: 1;
  padding-right: 15px;
}

.article-title {
  margin: 0 0 8px;
  font-size: 15px;
  font-weight: bold;
  opacity: 0.9;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.site-desc {
  margin: 0;
  font-size: 12px;
  opacity: 0.5;
}

.qrcode-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.qrcode-img {
  width: 60px;
  height: 60px;
  display: block;
  mix-blend-mode: multiply; /* 混合模式，去掉白底(虽然生成时已经是透明底) */
}

/* 暗色模式二维码特殊处理 */
.share-card.dark .qrcode-img {
  background: white; /* 必须有白底才能扫 */
  padding: 2px;
  border-radius: 4px;
  mix-blend-mode: normal;
}

.scan-tip {
  font-size: 10px;
  opacity: 0.5;
  transform: scale(0.9);
}

/* 控制区域 */
.controls {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.theme-selector {
  display: flex;
  align-items: center;
  gap: 10px;
}

.label {
  font-size: 14px;
  color: #666;
}

.theme-btn {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}

.theme-btn.active {
  border-color: #1890ff;
  transform: scale(1.1);
}

.download-btn {
  width: 100%;
  max-width: 300px;
  padding: 12px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 25px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

.download-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  box-shadow: none;
}

.download-btn:hover:not(:disabled) {
  background: #40a9ff;
  transform: translateY(-2px);
}
</style>