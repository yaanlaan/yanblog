<template>
  <div class="article-main-content">
    <div class="article-description" v-if="article.desc">
      <blockquote>{{ article.desc }}</blockquote>
    </div>
    
    <!-- PDF 文章 -->
    <div v-if="article.type === 2" class="pdf-container">
      <div v-if="article.pdf_url" class="pdf-wrapper">
         <iframe 
            :src="article.pdf_url" 
            width="100%" 
            height="800px" 
            style="border: none;"
            title="PDF Viewer"
         >
            <p>您的浏览器不支持 PDF 预览，请<a :href="article.pdf_url">下载 PDF</a>查看。</p>
         </iframe>
      </div>
      <div v-else class="pdf-error">
        PDF 文件未找到
      </div>
    </div>

    <!-- Markdwon 文章 -->
    <div v-else class="content" v-html="renderedContent" ref="contentRef" @click="handleContentClick" @mouseup="handleMouseUp"></div>
    
    <!-- 划词分享提示按钮 -->
    <div 
      v-if="showShareTip" 
      class="share-tip-btn" 
      :style="{ top: tipPosition.y + 'px', left: tipPosition.x + 'px' }"
      @mousedown.prevent="handleShareClick"
    >
      <i class="iconfont icon-share"></i> 分享
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUpdated, onUnmounted } from 'vue'
import { marked } from 'marked'
import katex from 'katex'
import mermaid from 'mermaid'
import 'katex/dist/katex.min.css'

const escapeHtml = (unsafe: unknown) => {
  if (typeof unsafe !== 'string') {
    return ''
  }
  return unsafe
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  type?: number
  pdf_url?: string
  createdAt: string
  updatedAt: string
}

interface Props {
  article: Article
}

const props = defineProps<Props>()
const contentRef = ref<HTMLElement | null>(null)

const emit = defineEmits<{
  (e: 'imageClick', imageSrc: string, imageAlt: string, images: string[], alts: string[]): void
  (e: 'shareSelection', text: string): void
}>()

const showShareTip = ref(false)
const tipPosition = ref({ x: 0, y: 0 })
const selectedText = ref('')

const handleMouseUp = () => {
  const selection = window.getSelection()
  if (!selection || selection.isCollapsed) {
    showShareTip.value = false
    return
  }

  const text = selection.toString().trim()
  if (text.length < 5) {
    showShareTip.value = false
    return
  }

  selectedText.value = text
  
  const range = selection.getRangeAt(0)
  const rect = range.getBoundingClientRect()
  
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop
  const scrollLeft = window.pageXOffset || document.documentElement.scrollLeft
  
  tipPosition.value = {
    x: rect.left + rect.width / 2,
    y: rect.top - 40
  }
  
  showShareTip.value = true
}

const handleShareClick = () => {
  emit('shareSelection', selectedText.value)
  showShareTip.value = false
  window.getSelection()?.removeAllRanges()
}

const handleScroll = () => {
  if (showShareTip.value) {
    showShareTip.value = false
  }
}

mermaid.initialize({
  startOnLoad: false,
  theme: 'neutral',
  securityLevel: 'loose'
})

// 独立链接渲染为卡片：将 <p><a href="...">text</a></p> 转为 link-card
const renderLinkCards = (html: string): string => {
  return html.replace(
    /<p>\s*<a href="([^"]+)"[^>]*>([^<]+)<\/a>\s*<\/p>/g,
    (_, href, text) => {
      // 提取域名作为描述
      let domain = ''
      try {
        domain = new URL(href).hostname
      } catch { domain = href }
      return `
<a href="${href}" target="_blank" rel="noopener noreferrer" class="link-card">
  <span class="link-card-title">${text}</span>
  <span class="link-card-desc">${domain}</span>
</a>`
    }
  )
}

const renderedContent = computed(() => {
  if (!props.article.content) {
    return ''
  }

  const contentWithIds = addIdsToHeadings(props.article.content)

  let html: string = marked.parse(contentWithIds) as string

  html = renderMath(html)
  html = renderLinkCards(html)

  return html
})

const addIdsToHeadings = (content: string) => {
  let headingCounter = 0
  return content.replace(/^(#{1,6})\s+(.+)$/gm, (_, hashes, text) => {
    headingCounter++
    const headingId = `heading-${headingCounter}`
    const slug = text.trim().replace(/\s+/g, '-').replace(/[^\w\u4e00-\u9fff-]/g, '')
    return `${hashes} <span id="${headingId}" data-slug="${slug}" class="heading-anchor"></span>${text}`
  })
}

const renderMath = (html: string) => {
  html = html.replace(/\$\$([\s\S]*?)\$\$/g, (_, formula) => {
    const cleanFormula = formula.trim()
    try {
      const rendered = katex.renderToString(cleanFormula, { displayMode: true })
      return `<div class="math-block" data-tex="${encodeURIComponent(cleanFormula)}">${rendered}</div>`
    } catch (error) {
      console.error('KaTeX渲染错误:', error)
      return `<span style="color: red;">公式渲染错误: ${formula}</span>`
    }
  })
  
  html = html.replace(/\$([^\$\n]+?)\$/g, (_, formula) => {
    const cleanFormula = formula.trim()
    try {
      const rendered = katex.renderToString(cleanFormula, { displayMode: false })
      return `<span class="math-inline" data-tex="${encodeURIComponent(cleanFormula)}">${rendered}</span>`
    } catch (error) {
      console.error('KaTeX渲染错误:', error)
      return `<span style="color: red;">公式渲染错误: ${formula}</span>`
    }
  })
  
  return html
}

const renderPostProcess = async () => {
  if (contentRef.value) {
    const tables = contentRef.value.querySelectorAll('table');
    tables.forEach((table) => {
      if (table.parentElement && !table.parentElement.classList.contains('table-wrapper')) {
        const wrapper = document.createElement('div');
        wrapper.className = 'table-wrapper';
        wrapper.style.overflowX = 'auto';
        wrapper.style.marginBottom = '20px';
        wrapper.style.maxWidth = '100%';
        table.parentElement.insertBefore(wrapper, table);
        wrapper.appendChild(table);
      }
    });
    
    const mermaidBlocks = contentRef.value.querySelectorAll('.mermaid-chart');
    for (let i = 0; i < mermaidBlocks.length; i++) {
      const block = mermaidBlocks[i];
      if (block instanceof HTMLElement) {
        try {
          const code = block.textContent || '';
          const { svg } = await mermaid.render(block.id + '-svg', code);
          block.innerHTML = svg;
        } catch (error) {
          console.error('Mermaid渲染错误:', error);
          block.innerHTML = '<p style="color: red;">图表渲染失败</p>';
        }
      }
    }
    
    const blockFormulas = contentRef.value.querySelectorAll('p');
    blockFormulas.forEach(element => {
      const text = element.textContent || '';
      if (text.startsWith('$$') && text.endsWith('$$')) {
        try {
          const formula = text.substring(2, text.length - 2);
          const katexHtml = katex.renderToString(formula.trim(), { displayMode: true });
          element.innerHTML = katexHtml;
        } catch (error) {
          console.error('KaTeX渲染错误:', error);
        }
      }
    });
  }
}

const handleContentClick = async (event: MouseEvent) => {
  const target = event.target as HTMLElement
  
  if (target.tagName === 'A' && target.getAttribute('href')?.startsWith('#')) {
    event.preventDefault()
    const hash = target.getAttribute('href')?.substring(1)
    
    if (hash) {
      try {
        const decodedHash = decodeURIComponent(hash)
        const slug = decodedHash.replace(/\s+/g, '-').replace(/[^\w\u4e00-\u9fff-]/g, '')
        
        const anchor = contentRef.value?.querySelector(`[data-slug="${slug}"]`)
        
        if (anchor) {
          anchor.scrollIntoView({ behavior: 'smooth', block: 'start' })
        }
      } catch (error) {
        console.error('Error decoding URI:', error)
      }
    }
    return
  }
  
  if (target.tagName === 'IMG') {
    event.preventDefault()
    
    const images = contentRef.value?.querySelectorAll('img') || []
    const imageSources: string[] = []
    const imageAlts: string[] = []
    
    images.forEach(img => {
      imageSources.push(img.getAttribute('src') || '')
      imageAlts.push(img.getAttribute('alt') || '')
    })
    
    const currentSrc = target.getAttribute('src') || ''
    const currentAlt = target.getAttribute('alt') || ''
    
    emit('imageClick', currentSrc, currentAlt, imageSources, imageAlts)
  }

  const copyBtn = target.closest('.mac-copy-btn') as HTMLElement
  if (copyBtn) {
    const code = copyBtn.getAttribute('data-code') || ''
    const decodedCode = code.replace(/&quot;/g, '"').replace(/&#39;/g, "'")
    if (decodedCode) {
      try {
        await navigator.clipboard.writeText(decodedCode)
        copyBtn.classList.add('copied')
        setTimeout(() => {
          copyBtn.classList.remove('copied')
        }, 2000)
      } catch (err) {
        console.error('复制失败:', err)
      }
    }
  }

  const mathElement = target.closest('.math-block, .math-inline') as HTMLElement
  if (mathElement) {
    event.preventDefault()
    const isShowingCode = mathElement.classList.contains('show-code')
    const tex = decodeURIComponent(mathElement.getAttribute('data-tex') || '')
    
    if (isShowingCode) {
      mathElement.classList.remove('show-code')
      try {
         const displayMode = mathElement.classList.contains('math-block')
         mathElement.innerHTML = katex.renderToString(tex, { displayMode })
      } catch(e) {
        console.error('KaTeX重渲染错误:', e)
      }
    } else {
      mathElement.classList.add('show-code')
      mathElement.textContent = tex
    }
  }
}

onMounted(() => {
  renderPostProcess()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

onUpdated(() => {
  renderPostProcess()
})
</script>

<style scoped>
.share-tip-btn {
  position: fixed;
  z-index: 1000;
  background: #333;
  color: #fff;
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  user-select: none;
}

.share-tip-btn::after {
  content: '';
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-top: 6px solid #333;
}

.article-description blockquote {
  font-size: 16px;
  color: var(--color-text-secondary);
  border-left: 4px solid var(--color-accent);
  padding: 15px 20px;
  margin: 0 0 30px 0;
  background: var(--color-background-mute);
  border-radius: 0 6px 6px 0;
  font-style: italic;
}

.content {
  font-size: 16px;
  line-height: 1.9;
  counter-reset: h1;
}

.content :deep(h1) {
  font-size: 28px;
  margin: 30px 0 20px;
  color: var(--color-heading);
  counter-reset: h2;
  font-weight: 700;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--color-border);
}

.content :deep(h2) {
  font-size: 24px;
  margin: 26px 0 16px;
  color: var(--color-heading);
  counter-reset: h3;
  font-weight: 600;
}

.content :deep(h3) {
  font-size: 20px;
  margin: 22px 0 14px;
  color: var(--color-heading);
  font-weight: 600;
}

.content :deep(h4),
.content :deep(h5),
.content :deep(h6) {
  font-size: 18px;
  margin: 18px 0 12px;
  color: var(--color-heading);
  font-weight: 600;
}

.content :deep(p) {
  margin: 18px 0;
  color: var(--color-text);
  text-align: justify;
}

.content :deep(ul),
.content :deep(ol) {
  padding-left: 30px;
  margin: 18px 0;
  color: var(--color-text);
}

.content :deep(li) {
  margin-bottom: 10px;
  position: relative;
}

.content :deep(ul) li::marker {
  color: var(--color-accent);
}

.content :deep(code) {
  background: var(--color-background-mute);
  color: var(--color-accent);
  padding: 3px 8px;
  border-radius: 4px;
  font-family: 'Fira Code', Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  border: 1px solid var(--color-border);
  word-break: break-all;
}

.content :deep(.mac-code-block) {
  margin: 20px 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0,0,0,0.15);
  background: var(--color-code-bg, #1e1e1e);
  border: 1px solid var(--color-code-border, #3c3c3c);
}

.content :deep(.mac-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: var(--color-code-header-bg, #252526);
  border-bottom: 1px solid var(--color-code-border, #3c3c3c);
}

.content :deep(.mac-dots) {
  display: flex;
  gap: 8px;
}

.content :deep(.mac-dots span) {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  transition: opacity 0.2s;
}

.content :deep(.mac-code-block:hover .mac-dots span) {
  opacity: 0.8;
}

.content :deep(.dot-red) {
  background: #ff5f56;
  box-shadow: 0 0 2px rgba(255, 95, 86, 0.3);
}

.content :deep(.dot-yellow) {
  background: #ffbd2e;
  box-shadow: 0 0 2px rgba(255, 189, 46, 0.3);
}

.content :deep(.dot-green) {
  background: #27c93f;
  box-shadow: 0 0 2px rgba(39, 201, 63, 0.3);
}

.content :deep(.mac-lang) {
  font-size: 12px;
  color: var(--color-code-lang, #858585);
  font-family: 'SF Mono', 'Fira Code', Consolas, monospace;
  text-transform: lowercase;
  letter-spacing: 0.5px;
}

.content :deep(.mac-code-body) {
  display: flex;
  overflow-x: auto;
  background: var(--color-code-body-bg, #1e1e1e);
}

.content :deep(.mac-line-numbers) {
  padding: 20px 6px;
  background: var(--color-code-line-num-bg, #252526);
  text-align: right;
  user-select: none;
  min-width: 38px;
  border-right: 1px solid var(--color-code-border, #3c3c3c);
}

.content :deep(.mac-line-numbers span) {
  display: block;
  line-height: 1.7;
  font-size: 12px;
  color: var(--color-code-line-num, #6b6b6b);
  font-family: 'Fira Code', Consolas, Monaco, monospace;
  user-select: none;
  opacity: 0.6;
}

.content :deep(.mac-code-content) {
  flex: 1;
  position: relative;
  overflow-x: auto;
  min-width: 0;
}

.content :deep(.mac-copy-btn) {
  position: absolute;
  top: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--color-code-copy-btn-bg, rgba(255,255,255,0.08));
  border: 1px solid var(--color-code-copy-btn-border, rgba(255,255,255,0.1));
  border-radius: 6px;
  color: var(--color-code-copy-btn-text, #858585);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  opacity: 0;
  visibility: hidden;
  backdrop-filter: blur(8px);
  font-weight: 500;
}

.content :deep(.mac-code-block:hover .mac-copy-btn) {
  opacity: 1;
  visibility: visible;
}

.content :deep(.mac-copy-btn:hover) {
  background: var(--color-code-copy-btn-hover-bg, rgba(255,255,255,0.12));
  border-color: var(--color-code-copy-btn-hover-border, rgba(255,255,255,0.2));
  color: var(--color-code-copy-btn-hover-text, #a0a0a0);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.2);
}

.content :deep(.mac-copy-btn.copied) {
  background: #27c93f;
  border-color: #27c93f;
  color: #fff;
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(39, 201, 63, 0.3);
}

.content :deep(.mac-copy-btn .copy-icon),
.content :deep(.mac-copy-btn .copied-icon) {
  width: 14px;
  height: 14px;
  transition: transform 0.2s;
}

.content :deep(.mac-copy-btn:hover .copy-icon) {
  transform: scale(1.1);
}

.content :deep(.mac-copy-btn.copied .copy-icon) {
  display: none;
}

.content :deep(.mac-copy-btn:not(.copied) .copied-icon) {
  display: none;
}

.content :deep(.mac-code-content pre) {
  margin: 0;
  padding: 20px 12px;
  background: transparent;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.content :deep(.mac-code-content code) {
  font-family: 'Fira Code', 'SF Mono', Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.7;
  background: transparent !important;
  padding: 0 !important;
  border: none !important;
  white-space: pre;
  color: var(--color-code-text, #abb2bf);
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
}

.content :deep(.mac-code-body::-webkit-scrollbar) {
  height: 8px;
}

.content :deep(.mac-code-body::-webkit-scrollbar-track) {
  background: var(--color-code-scrollbar-track, #1e1e1e);
}

.content :deep(.mac-code-body::-webkit-scrollbar-thumb) {
  background: var(--color-code-scrollbar-thumb, #4a4a4a);
  border-radius: 4px;
}

.content :deep(.mac-code-body::-webkit-scrollbar-thumb:hover) {
  background: var(--color-code-scrollbar-thumb-hover, #5a5a5a);
}

.line-numbers span {
  display: block;
  line-height: 1.6;
  font-size: 13px;
  color: #5c6370;
  font-family: monospace;
}

.code-body pre {
  margin: 0;
  padding: 16px;
  overflow: visible;
  background: transparent;
}

.code-body code {
  font-family: 'Fira Code', Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  background: transparent;
  padding: 0;
  border: none;
  white-space: pre;
  color: #abb2bf;
}

.content :deep(blockquote) {
  border-left: 4px solid var(--color-accent);
  padding: 15px 20px;
  margin: 20px 0;
  background: var(--color-background-mute);
  color: var(--color-text-secondary);
  border-radius: 0 6px 6px 0;
  font-style: italic;
}

.content :deep(.article-image) {
  max-width: 100%;
  height: auto;
  border-radius: 10px;
  margin: 20px 0;
  cursor: zoom-in;
  transition: all 0.3s;
  box-shadow: 0 4px 16px var(--color-shadow);
  display: block;
}

.content :deep(.article-image:hover) {
  opacity: 0.95;
  box-shadow: 0 8px 24px var(--color-shadow);
}

.content :deep(.article-link) {
  color: var(--color-accent);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.3s;
  font-weight: 500;
}

.content :deep(.article-link:hover) {
  border-bottom-color: var(--color-accent);
  background: rgba(66, 184, 131, 0.1);
  padding: 2px 4px;
  margin: -2px -4px;
  border-radius: 3px;
}

.content :deep(.table-wrapper) {
  margin: 20px 0;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  overflow-x: auto;
}

.content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  min-width: 600px;
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid var(--color-border);
  padding: 14px 16px;
  text-align: left;
  color: var(--color-text);
}

.content :deep(th) {
  background: var(--color-background-mute);
  color: var(--color-heading);
  font-weight: 600;
  white-space: nowrap;
}

.content :deep(tr:nth-child(even)) {
  background: var(--color-background-soft);
}

.content :deep(tr:hover) {
  background: var(--color-background-mute);
}

.content :deep(.katex-display) {
  margin: 25px 0;
  overflow-x: auto;
  overflow-y: hidden;
  max-width: 100%;
}

.content :deep(.katex) {
  white-space: normal;
  max-width: 100%;
}

.content :deep(.katex-html) {
  overflow-x: auto;
  overflow-y: hidden;
  max-width: 100%;
}

.content :deep(.mermaid-chart) {
  text-align: center;
  margin: 25px 0;
  overflow-x: auto;
  max-width: 100%;
  background: var(--mermaid-bg, transparent);
  border-radius: 8px;
  padding: 12px;
}

.content :deep(.mermaid-chart svg) {
  max-width: 100%;
  height: auto;
}

/* Mermaid 图表文字在暗色模式下可见 */
.content :deep(.mermaid-chart svg .label),
.content :deep(.mermaid-chart svg text),
.content :deep(.mermaid-chart svg .edgeLabel) {
  fill: var(--mermaid-text, currentColor) !important;
}

.content :deep(.mermaid-chart svg .edgePath .path) {
  stroke: var(--mermaid-line, currentColor) !important;
}

.content :deep(.heading-anchor) {
  position: absolute;
  top: -80px;
  left: 0;
  visibility: hidden;
}

.content :deep(h1),
.content :deep(h2),
.content :deep(h3),
.content :deep(h4),
.content :deep(h5),
.content :deep(h6) {
  position: relative;
}

.content :deep(.math-block),
.content :deep(.math-inline) {
  cursor: pointer;
  transition: opacity 0.2s;
}

.content :deep(.math-block:hover),
.content :deep(.math-inline:hover) {
  opacity: 0.8;
}

.content :deep(.math-block.show-code) {
  background: var(--color-background-soft);
  padding: 15px;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  white-space: pre-wrap;
  font-family: monospace;
  color: var(--color-text);
  display: block;
  text-align: left;
  margin: 20px 0;
}

.article-main-content {
  background: var(--color-background-soft);
  padding: 45px 50px;
  border-radius: 16px;
  box-shadow: 0 4px 20px var(--color-shadow);
  color: var(--color-text);
  line-height: 1.9;
  font-size: 16px;
  position: relative;
  min-height: 200px;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

@media (max-width: 992px) {
  .article-main-content {
    padding: 30px 25px;
    font-size: 15px;
  }
  
  .content :deep(h1) {
    font-size: 24px;
  }
  
  .content :deep(h2) {
    font-size: 22px;
  }
  
  .content :deep(h3) {
    font-size: 19px;
  }
}

@media (max-width: 768px) {
  .article-main-content {
    padding: 20px 16px;
    font-size: 15px;
    border-radius: 10px;
  }
  
  .content :deep(h1) {
    font-size: 22px;
    padding-bottom: 8px;
  }
  
  .content :deep(h2) {
    font-size: 20px;
  }
  
  .content :deep(h3) {
    font-size: 18px;
  }
  
  .content :deep(p) {
    margin: 15px 0;
  }
  
  .content :deep(ul),
  .content :deep(ol) {
    padding-left: 20px;
  }
  
  .mac-header {
    padding: 8px 10px;
  }
  
  .mac-lang {
    font-size: 11px;
  }
  
  .mac-dots span {
    width: 10px;
    height: 10px;
  }
  
  .line-numbers {
    min-width: 35px;
    padding: 12px 6px;
  }
  
  .line-numbers span {
    font-size: 11px;
  }
  
  .mac-code-content pre {
    padding: 14px 12px;
  }
  
  .mac-code-content code {
    font-size: 12px;
    line-height: 1.6;
  }
  
  .mac-copy-btn {
    padding: 4px 8px;
    font-size: 11px;
  }
  
  .mac-copy-btn svg {
    width: 12px;
    height: 12px;
  }
}

.share-tip-btn {
  position: fixed;
  background: var(--color-heading);
  color: var(--color-background-soft);
  padding: 10px 18px;
  border-radius: 8px;
  cursor: pointer;
  z-index: 1000;
  transform: translate(-50%, -100%);
  box-shadow: 0 4px 16px rgba(0,0,0,0.25);
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  transition: all 0.2s;
}

.share-tip-btn::after {
  content: '';
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  border-width: 6px 6px 0;
  border-style: solid;
  border-color: var(--color-heading) transparent transparent transparent;
}

.share-tip-btn:hover {
  transform: translate(-50%, -110%);
  box-shadow: 0 6px 20px rgba(0,0,0,0.3);
}

.pdf-container {
  width: 100%;
  margin-top: 20px;
}

.pdf-wrapper {
  width: 100%;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  border-radius: 10px;
  overflow: hidden;
}

.pdf-error {
  padding: 60px 40px;
  text-align: center;
  background: var(--color-background-soft);
  color: var(--color-text-secondary);
  border-radius: 10px;
}

.content :deep(.link-card) {
  display: block;
  max-width: 400px;
  width: 100%;
  margin: 25px auto;
  padding: 20px;
  background-color: var(--color-background-mute);
  border-radius: 10px;
  text-decoration: none;
  transition: all 0.3s;
  border: 1px solid var(--color-border);
  box-sizing: border-box;
}

.content :deep(.link-card:hover) {
  background-color: var(--color-background-soft);
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  transform: translateY(-2px);
}

.content :deep(.link-card-title) {
  display: block;
  font-weight: 600;
  font-size: 16px;
  color: var(--color-heading);
  margin-bottom: 8px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.content :deep(.link-card-desc) {
  display: block;
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 768px) {
  .content :deep(.link-card) {
    max-width: 100%;
    padding: 16px;
  }
}
</style>