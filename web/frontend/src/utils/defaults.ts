import { computed } from 'vue'
import { useSiteInfoStore } from '@/stores/siteInfo'
import { storeToRefs } from 'pinia'

// 内置回退图片（Vite 构建时解析为实际路径）
const fallbackCover = new URL('../assets/img/无封面.jpg', import.meta.url).href
const fallbackAvatar = new URL('../assets/default-avatar.jpg', import.meta.url).href
const fallbackQr = new URL('../assets/default-qr.png', import.meta.url).href

// 默认封面图（优先使用配置，缺失时回退到内置图片）
export function useDefaultCover() {
  const { siteInfo } = storeToRefs(useSiteInfoStore())
  return computed(() => siteInfo.value.default_images?.cover || fallbackCover)
}

// 默认头像
export function useDefaultAvatar() {
  const { siteInfo } = storeToRefs(useSiteInfoStore())
  return computed(() => siteInfo.value.default_images?.avatar || fallbackAvatar)
}

// 默认二维码
export function useDefaultQrCode() {
  const { siteInfo } = storeToRefs(useSiteInfoStore())
  return computed(() => siteInfo.value.default_images?.qr_code || fallbackQr)
}
