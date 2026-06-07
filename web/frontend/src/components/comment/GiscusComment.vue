<template>
  <div v-if="shouldShow" class="comments-section">
    <div class="comment-header">
      <i class="iconfont icon-comment"></i>
      <span>评论</span>
    </div>
    <div ref="giscusContainer" class="giscus-wrapper"></div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch, nextTick } from 'vue'
import { useSiteInfoStore } from '@/stores/siteInfo'
import { useThemeStore } from '@/stores/theme'

const siteInfoStore = useSiteInfoStore()
const themeStore = useThemeStore()
const giscusContainer = ref<HTMLElement | null>(null)

// Check if comments are enabled
const shouldShow = computed(() => {
  return siteInfoStore.siteInfo?.comment?.enable && 
         siteInfoStore.siteInfo?.comment?.type === 'giscus'
})

// Build script element
const loadGiscus = () => {
    if (!shouldShow.value || !giscusContainer.value) return
    
    // Clear existing
    giscusContainer.value.innerHTML = ''
    
    const config = siteInfoStore.siteInfo.comment.giscus
    
    // Check if repo is configured (user might have left default)
    if (config.repo === 'your-username/your-repo' || !config.repo_id) {
        giscusContainer.value.innerHTML = '<div class="giscus-warning">请在 config.yaml 中配置 Giscus 评论系统参数</div>'
        return
    }

    const script = document.createElement('script')
    script.src = 'https://giscus.app/client.js'
    script.setAttribute('data-repo', config.repo)
    script.setAttribute('data-repo-id', config.repo_id)
    script.setAttribute('data-category', config.category)
    script.setAttribute('data-category-id', config.category_id)
    script.setAttribute('data-mapping', config.mapping)
    script.setAttribute('data-reactions-enabled', config.reactions_enabled)
    script.setAttribute('data-emit-metadata', config.emit_metadata)
    script.setAttribute('data-input-position', config.input_position)
    
    // Dynamic theme
    const theme = themeStore.theme === 'dark' ? 'dark' : 'light'
    script.setAttribute('data-theme', theme)
    
    script.setAttribute('data-lang', config.lang)
    script.setAttribute('data-loading', config.loading)
    script.setAttribute('crossorigin', 'anonymous')
    script.async = true
    
    giscusContainer.value.appendChild(script)
}

// Watch theme changes
watch(() => themeStore.theme, (newTheme) => {
    const iframe = document.querySelector<HTMLIFrameElement>('iframe.giscus-frame')
    if (!iframe) return
    const theme = newTheme === 'dark' ? 'dark' : 'light'
    iframe.contentWindow?.postMessage(
        { giscus: { setConfig: { theme } } },
        'https://giscus.app'
    )
})

onMounted(() => {
    // Wait for site info to be loaded
    if (siteInfoStore.siteInfo.comment) {
        loadGiscus()
    } else {
        // Watch for store changes just in case
        const unwatch = watch(() => siteInfoStore.siteInfo, (newVal) => {
            if (newVal.comment) {
                nextTick(() => {
                   loadGiscus()
                   unwatch()
                })
            }
        }, { deep: true })
    }
})

// Optional: Dynamic theme update if we implement dark mode switch later
// This would involve sending a postMessage to the iframe
</script>

<style scoped>
.comments-section {
  margin-top: 60px;
  border-top: 1px solid var(--color-border);
  padding-top: 30px;
  animation: fadeIn 0.5s ease;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--color-heading);
  margin-bottom: 20px;
}

.comment-header i {
  font-size: 20px;
  color: var(--color-primary, #42b883);
}

.giscus-wrapper {
    min-height: 200px;
}

.giscus-warning {
    padding: 20px;
    background: var(--color-background-soft);
    color: var(--color-text);
    border-radius: 4px;
    text-align: center;
    border: 1px solid var(--color-border);
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>
