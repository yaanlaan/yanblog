
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
import hljs from 'highlight.js'
import mermaid from 'mermaid'
import 'katex/dist/katex.min.css'
import 'highlight.js/styles/atom-one-dark.css'

// 定义Props
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  // Extend for PDF
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


// 定义事件
const emit = defineEmits<{
  (e: 'imageClick', imageSrc: string, imageAlt: string, images: string[], alts: string[]): void
  (e: 'shareSelection', text: string): void
}>()

// 划词分享相关
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
  if (text.length < 5) { // 太短不显示
    showShareTip.value = false
    return
  }

  selectedText.value = text
  
  // 计算位置
  const range = selection.getRangeAt(0)
  const rect = range.getBoundingClientRect()
  
  // 更新位置 (相对于视口，需要转为相对于页面或绝对定位)
  // 这里我们使用 fixed 定位在 App 级别或者计算相对于 container 的位置
  // 为了简单，我们让 share-tip-btn 使用 fixed 定位
  // 或者计算 scroll
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop
  const scrollLeft = window.pageXOffset || document.documentElement.scrollLeft
  
  // 获取 contentRef 的偏移量，以便转换为相对于 content 的位置（如果是 absolute）
  // 或者直接用 fixed (简单且效果好)
  // 这里假设 share-tip-btn 是 fixed 的
  
  // 实际上 .article-main-content 通常是 static 或 relative。
  // 我们在 handleMouseUp 里计算的是相对于浏览器视口的位置，加上 scroll 就是相对于文档。
  // 但 tip 是在 .article-main-content 里的。
  // 简单起见，我们把 CSS 设置为 fixed，直接用 rect.top/left
  
  // 将按钮放在选区上方居中
  tipPosition.value = {
    x: rect.left + rect.width / 2,
    y: rect.top - 40 // 向上偏移
  }
  
  showShareTip.value = true
}

const handleShareClick = () => {
  emit('shareSelection', selectedText.value)
  showShareTip.value = false
  // 清除选区
  window.getSelection()?.removeAllRanges()
}

// 监听滚动关闭提示
const handleScroll = () => {
  if (showShareTip.value) {
    showShareTip.value = false
  }
}

onMounted(() => {
  renderPostProcess()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})

// 初始化mermaid
mermaid.initialize({ 
  startOnLoad: false,
  theme: 'default',
  securityLevel: 'loose'
})

// 渲染Markdown内容
const renderedContent = computed(() => {
  if (!props.article.content) return ''
  
  // 添加ID到标题
  let contentWithIds = addIdsToHeadings(props.article.content)
  
  // 配置marked
  const renderer = new marked.Renderer()
  
  // 自定义代码块渲染
  renderer.code = ({ text, lang }: { text: string, lang?: string }) => {
    const language = lang || 'plaintext'
    
    // Mermaid 特殊处理
    if (language === 'mermaid') {
      return `<div class="mermaid-chart">${text}</div>`
    }

    const validLang = hljs.getLanguage(language) ? language : 'plaintext'
    const highlighted = hljs.highlight(text, { language: validLang }).value
    
    // 生成行号
    const lines = text.replace(/\n$/, '').split('\n')
    const lineNumbers = lines.map((_, i) => `<span>${i + 1}</span>`).join('')

    return `<div class="code-wrapper">
              <div class="code-header">
                <div class="mac-buttons">
                  <span class="mac-button red"></span>
                  <span class="mac-button yellow"></span>
                  <span class="mac-button green"></span>
                </div>
                <span class="lang-name">${validLang}</span>
                <button class="copy-btn" data-code="${encodeURIComponent(text)}">
                  <span class="copy-text">复制</span>
                </button>
              </div>
              <div class="code-body">
                <div class="line-numbers">${lineNumbers}</div>
                <pre><code class="hljs language-${validLang}">${highlighted}</code></pre>
              </div>
            </div>`
  }

  // 自定义链接渲染 (Link Card Support)
  // Syntax: [card: Title | Description](URL)
  renderer.link = ({ href, title, text }: { href: string, title?: string | null, text: string }) => {
    if (text.startsWith('card:')) {
      const content = text.replace(/^card:\s*/, '')
      const parts = content.split('|')
      const cardTitle = parts[0].trim()
      const cardDesc = parts[1] ? parts[1].trim() : href.replace(/^https?:\/\//, '')
      
      return `
        <a class="link-card" href="${href}" target="_blank" rel="noopener noreferrer">
          <span class="link-card-content">
            <span class="link-card-title">${cardTitle}</span>
            <span class="link-card-desc">${cardDesc}</span>
          </span>
        </a>
      `
    }
    
    return `<a href="${href}" target="_blank" rel="noopener noreferrer" title="${title || ''}">${text}</a>`
  }

  // 首先使用marked解析Markdown
  let html: string = marked.parse(contentWithIds, { renderer }) as string
  
  // 处理mermaid图表 (Modified: renderer.code now handles mermaid divs)
  // html = renderMermaid(html)
  
  // 然后处理数学公式
  html = renderMath(html)
  
  return html
})

// 为标题添加ID
const addIdsToHeadings = (content: string) => {
  let headingCounter = 0
  return content.replace(/^(#{1,6})\s+(.+)$/gm, (_, hashes, text) => {
    headingCounter++
    const id = `heading-${headingCounter}`
    return `${hashes} <span id="${id}" class="heading-anchor"></span>${text}`
  })
}

// 渲染mermaid图表 - Deprecated since we handle it in renderer.code
// const renderMermaid = (html: string) => { ... }

// (Wait, I should define renderer.code properly instead of post-regex)



// 渲染数学公式
const renderMath = (html: string) => {
  // 处理块级公式（$$...$$）
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
  
  // 处理行内公式（$...$）
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

// 渲染完成后的处理
const renderPostProcess = async () => {
  if (contentRef.value) {
    // 处理代码高亮
    const codeBlocks = contentRef.value.querySelectorAll('pre code');
    codeBlocks.forEach((block) => {
      if (block instanceof HTMLElement) {
        hljs.highlightElement(block);
      }
    });
    
    // 处理mermaid图表
    const mermaidBlocks = contentRef.value.querySelectorAll('.mermaid-chart');
    for (let i = 0; i < mermaidBlocks.length; i++) {
      const block = mermaidBlocks[i];
      if (block instanceof HTMLElement) {
        try {
          const code = block.textContent || '';
          // 正确处理mermaid.render的返回值
          const { svg } = await mermaid.render(block.id + '-svg', code);
          block.innerHTML = svg;
        } catch (error) {
          console.error('Mermaid渲染错误:', error);
          block.innerHTML = '<p style="color: red;">图表渲染失败</p>';
        }
      }
    }
    
    // 处理块级公式（$$...$$）
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

// 处理内容区域点击事件
const handleContentClick = async (event: MouseEvent) => {
  const target = event.target as HTMLElement
  
  // 检查是否点击了图片
  if (target.tagName === 'IMG') {
    event.preventDefault()
    
    // 获取所有图片
    const images = contentRef.value?.querySelectorAll('img') || []
    const imageSources: string[] = []
    const imageAlts: string[] = []
    
    // 提取所有图片的src和alt属性
    images.forEach(img => {
      imageSources.push(img.getAttribute('src') || '')
      imageAlts.push(img.getAttribute('alt') || '')
    })
    
    // 获取当前点击图片的src
    const currentSrc = target.getAttribute('src') || ''
    const currentAlt = target.getAttribute('alt') || ''
    
    // 触发imageClick事件
    emit('imageClick', currentSrc, currentAlt, imageSources, imageAlts)
  }

  // 检查是否点击了复制按钮
  const copyBtn = target.closest('.copy-btn') as HTMLElement
  if (copyBtn) {
    const code = decodeURIComponent(copyBtn.getAttribute('data-code') || '')
    if (code) {
      try {
        await navigator.clipboard.writeText(code)
        const textSpan = copyBtn.querySelector('.copy-text')
        if (textSpan) {
          const originalText = textSpan.textContent
          textSpan.textContent = '已复制!'
          setTimeout(() => {
            textSpan.textContent = originalText
          }, 2000)
        }
      } catch (err) {
        console.error('复制失败:', err)
      }
    }
  }

  // 检查是否点击了数学公式
  const mathElement = target.closest('.math-block, .math-inline') as HTMLElement
  if (mathElement) {
    event.preventDefault()
    const isShowingCode = mathElement.classList.contains('show-code')
    const tex = decodeURIComponent(mathElement.getAttribute('data-tex') || '')
    
    if (isShowingCode) {
      // 切换回渲染视图
      mathElement.classList.remove('show-code')
      try {
         const displayMode = mathElement.classList.contains('math-block')
         mathElement.innerHTML = katex.renderToString(tex, { displayMode })
      } catch(e) {
        console.error('KaTeX重渲染错误:', e)
      }
    } else {
      // 显示代码
      mathElement.classList.add('show-code')
      mathElement.textContent = tex
    }
  }
}

// 在组件挂载和更新时进行后处理
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
  border-left: 4px solid var(--color-border);
  padding: 10px 20px;
  margin: 0 0 30px 0;
  background: var(--color-background-mute);
  border-radius: 6px;
}

.content {
  font-size: 16px;
  line-height: 1.8;
  counter-reset: h1; /* 初始化计数器 */
}

/* 标题自动编号 */
.content :deep(h1) {
  font-size: 24px;
  margin: 24px 0 16px;
  color: var(--color-heading);
  counter-reset: h2;
}

.content :deep(h1)::before {
  counter-increment: h1;
  content: counter(h1) ". ";
  color: var(--color-text-secondary);
  font-weight: normal;
  margin-right: 8px;
}

.content :deep(h2) {
  font-size: 22px;
  margin: 22px 0 14px;
  color: var(--color-heading);
  counter-reset: h3;
}

.content :deep(h2)::before {
  counter-increment: h2;
  content: counter(h1) "." counter(h2) " ";
  color: var(--color-text-secondary);
  font-weight: normal;
  margin-right: 8px;
}

.content :deep(h3) {
  font-size: 20px;
  margin: 20px 0 12px;
  color: var(--color-heading);
}

.content :deep(h3)::before {
  counter-increment: h3;
  content: counter(h1) "." counter(h2) "." counter(h3) " ";
  color: var(--color-text-secondary);
  font-weight: normal;
  margin-right: 8px;
}

.content :deep(p) {
  margin: 16px 0;
  color: var(--color-text);
}

.content :deep(ul),
.content :deep(ol) {
  padding-left: 30px;
  margin: 16px 0;
  color: var(--color-text);
}

.content :deep(li) {
  margin-bottom: 8px;
}

/* 行内代码 */
.content :deep(code) {
  background: var(--color-background-mute);
  color: var(--color-text);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  border: 1px solid var(--color-border);
}

/* 代码块容器 */
.content :deep(.code-wrapper) {
  margin: 20px 0;
  border-radius: 8px;
  overflow: hidden;
  background: #282c34; /* Keep dark for code blocks usually */
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

/* 代码块头部 (Mac 风格) */
.content :deep(.code-header) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: #21252b;
  border-bottom: 1px solid #181a1f;
}

.content :deep(.mac-buttons) {
  display: flex;
  gap: 8px;
}

.content :deep(.mac-button) {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.content :deep(.mac-button.red) { background: #ff5f56; }
.content :deep(.mac-button.yellow) { background: #ffbd2e; }
.content :deep(.mac-button.green) { background: #27c93f; }

.content :deep(.lang-name) {
  color: #abb2bf;
  font-size: 12px;
  text-transform: uppercase;
  font-family: sans-serif;
}

.content :deep(.copy-btn) {
  background: transparent;
  border: none;
  color: #abb2bf;
  cursor: pointer;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
}

.content :deep(.copy-btn:hover) {
  background: rgba(255,255,255,0.1);
  color: #fff;
}

.content :deep(.code-body) {
  display: flex;
  background: #282c34;
  overflow: auto;
  position: relative;
}

.content :deep(.line-numbers) {
  padding: 16px 0 16px 10px;
  text-align: right;
  color: #495162;
  border-right: 1px solid #353b45;
  background: #282c34;
  user-select: none;
  display: flex;
  flex-direction: column;
  position: sticky;
  left: 0;
  z-index: 1;
  min-width: 40px;
}

.content :deep(.line-numbers span) {
  padding-right: 10px;
  line-height: 1.6;
  font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  height: 22.4px; /* 14px * 1.6 */
}

.content :deep(pre) {
  background: transparent;
  padding: 16px;
  margin: 0;
  overflow: visible;
  color: #abb2bf;
  flex: 1;
}

.content :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
  border: none;
  font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.6 !important;
  display: block;
}

/* 引用块 (灰白黑风格) */
.content :deep(blockquote) {
  border-left: 4px solid var(--color-border);
  padding: 10px 20px;
  margin: 20px 0;
  background: var(--color-background-mute);
  color: var(--color-text-secondary);
  border-radius: 0 4px 4px 0;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 20px 0;
  cursor: zoom-in;
  transition: all 0.3s;
  box-shadow: 0 4px 12px var(--color-shadow);
}

.content :deep(img:hover) {
  opacity: 0.95;
  box-shadow: 0 8px 24px var(--color-shadow);
}

/* 链接 (灰白黑风格) */
.content :deep(a) {
  color: var(--color-text);
  text-decoration: none;
  border-bottom: 1px solid var(--color-text-secondary);
  transition: all 0.3s;
  font-weight: 500;
}

.content :deep(a:hover) {
  color: var(--color-heading);
  border-bottom-color: var(--color-heading);
  background: var(--color-background-mute);
}

/* 表格 (灰白黑风格) */
.content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid var(--color-border);
  padding: 12px 15px;
  text-align: left;
  color: var(--color-text);
}

.content :deep(th) {
  background: var(--color-background-mute);
  color: var(--color-heading);
  font-weight: 600;
}

/* KaTeX 样式修复 */
.content :deep(.katex-display) {
  margin: 20px 0;
  overflow-x: auto;
  overflow-y: hidden;
}

.content :deep(.katex) {
  white-space: nowrap;
}

/* Mermaid 图表样式 */
.content :deep(.mermaid-chart) {
  text-align: center;
  margin: 20px 0;
  overflow-x: auto;
}

.content :deep(.mermaid-chart svg) {
  max-width: 100%;
  height: auto;
}

/* 标题锚点样式 */
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

/* 数学公式交互样式 */
.content :deep(.math-block),
.content :deep(.math-inline) {
  cursor: pointer;
  transition: opacity 0.2s;
}

.content :deep(.math-block:hover),
.content :deep(.math-inline:hover) {
  opacity: 0.8;
}

.content :deep(.math-code) {
  font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  background: var(--color-background-mute);
  padding: 2px 6px;
  border-radius: 3px;
  color: #d63384;
  font-size: 0.9em;
}

.content :deep(.math-block.show-code) {
  background: var(--color-background-soft);
  padding: 10px;
  border-radius: 4px;
  border: 1px solid var(--color-border);
  white-space: pre-wrap;
  font-family: monospace;
  color: var(--color-text);
  display: block;
  text-align: left;
}

.article-main-content {
  background: var(--color-background-soft);
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px var(--color-shadow);
  color: var(--color-text);
  line-height: 1.8;
  font-size: 16px;
  position: relative;
  min-height: 200px; /* 防止内容过少时太扁 */
}

/* Share Tip Button */
.share-tip-btn {
  position: fixed;
  background: var(--color-heading);
  color: var(--color-background-soft);
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  z-index: 1000;
  transform: translate(-50%, -100%);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 500;
  transition: all 0.2s;
  /* Add arrow */
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
}

/* PDF 样式 */
.pdf-container {
  width: 100%;
  margin-top: 20px;
}

.pdf-wrapper {
  width: 100%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  overflow: hidden;
}

.pdf-error {
  padding: 40px;
  text-align: center;
  background: var(--color-background-soft);
  color: var(--color-text-secondary);
  border-radius: 8px;
}

/* Link Card Styles */
.content :deep(.link-card) {
  display: block; /* 块级显示 */
  max-width: 390px; /* 限制宽度 */
  width: 100%;
  margin: 20px auto; /* 上下20px，左右自动(居中) */
  padding: 16px 20px; /* 增加内边距 */
  background-color: var(--color-background-mute); /* 知乎卡片背景色 */
  border-radius: 8px;
  text-decoration: none;
  transition: background-color 0.2s;
  border: none;
}

.content :deep(.link-card:hover) {
  background-color: var(--color-border-hover); /* 悬停颜色稍深 */
}

.content :deep(.link-card-content) {
  display: flex;
  flex-direction: column; /* 上下布局 */
  justify-content: center;
}

.content :deep(.link-card-title) {
  display: block; /* 覆盖默认inline */
  font-weight: 600;
  font-size: 16px;
  color: #121212;
  margin-bottom: 4px; /* 标题和描述的间距 */
  line-height: 1.4;
  /* 允许换行，最多两行 */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.content :deep(.link-card-desc) {
  display: block; /* 覆盖默认inline */
  font-size: 13px;
  color: #999;
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.content :deep(.link-card-icon) {
  /* 隐藏之前的图标样式 */
  display: none;
}

/* 适配移动端 */
@media (max-width: 768px) {
  .content :deep(.link-card) {
    max-width: 100%; /* 移动端占满 */
  }
}
</style>