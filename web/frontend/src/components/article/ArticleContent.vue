<template>
  <div class="article-main-content">
    <div class="article-description" v-if="article.desc">
      <blockquote>{{ article.desc }}</blockquote>
    </div>
    
    <div class="content" v-html="renderedContent" ref="contentRef" @click="handleContentClick"></div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUpdated } from 'vue'
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
}>()

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
    const validLang = hljs.getLanguage(language) ? language : 'plaintext'
    const highlighted = hljs.highlight(text, { language: validLang }).value
    
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
              <pre><code class="hljs language-${validLang}">${highlighted}</code></pre>
            </div>`
  }

  // 首先使用marked解析Markdown
  let html: string = marked.parse(contentWithIds, { renderer }) as string
  
  // 处理mermaid图表
  html = renderMermaid(html)
  
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

// 渲染mermaid图表
const renderMermaid = (html: string) => {
  // 查找mermaid代码块
  const mermaidRegex = /<pre><code class="([^"]*)mermaid([^"]*)">([\s\S]*?)<\/code><\/pre>/g;
  
  let match;
  let newHtml = html;
  let mermaidCounter = 0;
  
  while ((match = mermaidRegex.exec(html)) !== null) {
    const fullMatch = match[0];
    const mermaidCode = match[3];
    const mermaidId = `mermaid-${mermaidCounter++}`;
    
    // 替换为占位符div
    newHtml = newHtml.replace(
      fullMatch,
      `<div class="mermaid-chart" id="${mermaidId}">${mermaidCode}</div>`
    );
  }
  
  return newHtml;
}

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
})

onUpdated(() => {
  renderPostProcess()
})
</script>

<style scoped>
.article-description blockquote {
  font-size: 16px;
  color: #666;
  border-left: 4px solid #ccc;
  padding: 10px 20px;
  margin: 0 0 30px 0;
  background: #f8f9fa;
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
  color: #333;
  counter-reset: h2;
}

.content :deep(h1)::before {
  counter-increment: h1;
  content: counter(h1) ". ";
  color: #666;
  font-weight: normal;
  margin-right: 8px;
}

.content :deep(h2) {
  font-size: 22px;
  margin: 22px 0 14px;
  color: #333;
  counter-reset: h3;
}

.content :deep(h2)::before {
  counter-increment: h2;
  content: counter(h1) "." counter(h2) " ";
  color: #666;
  font-weight: normal;
  margin-right: 8px;
}

.content :deep(h3) {
  font-size: 20px;
  margin: 20px 0 12px;
  color: #333;
}

.content :deep(h3)::before {
  counter-increment: h3;
  content: counter(h1) "." counter(h2) "." counter(h3) " ";
  color: #666;
  font-weight: normal;
  margin-right: 8px;
}

.content :deep(p) {
  margin: 16px 0;
  color: #333;
}

.content :deep(ul),
.content :deep(ol) {
  padding-left: 30px;
  margin: 16px 0;
  color: #333;
}

.content :deep(li) {
  margin-bottom: 8px;
}

/* 行内代码 */
.content :deep(code) {
  background: #f0f0f0;
  color: #333;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
  font-size: 14px;
  border: 1px solid #e0e0e0;
}

/* 代码块容器 */
.content :deep(.code-wrapper) {
  margin: 20px 0;
  border-radius: 8px;
  overflow: hidden;
  background: #282c34;
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

.content :deep(pre) {
  background: #282c34;
  padding: 16px;
  margin: 0;
  overflow: auto;
  color: #abb2bf;
}

.content :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
  border: none;
}

/* 引用块 (灰白黑风格) */
.content :deep(blockquote) {
  border-left: 4px solid #ccc;
  padding: 10px 20px;
  margin: 20px 0;
  background: #f8f9fa;
  color: #666;
  border-radius: 0 4px 4px 0;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 20px 0;
  cursor: zoom-in;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.content :deep(img:hover) {
  opacity: 0.95;
  box-shadow: 0 8px 24px rgba(0,0,0,0.15);
}

/* 链接 (灰白黑风格) */
.content :deep(a) {
  color: #333;
  text-decoration: none;
  border-bottom: 1px solid #999;
  transition: all 0.3s;
  font-weight: 500;
}

.content :deep(a:hover) {
  color: #000;
  border-bottom-color: #000;
  background: rgba(0,0,0,0.05);
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
  border: 1px solid #e0e0e0;
  padding: 12px 15px;
  text-align: left;
  color: #333;
}

.content :deep(th) {
  background: #f5f5f5;
  color: #000;
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
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  color: #d63384;
  font-size: 0.9em;
}

.content :deep(.math-block.show-code) {
  background: #f8f9fa;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #eee;
  white-space: pre-wrap;
  font-family: monospace;
  color: #333;
  display: block;
  text-align: left;
}
</style>