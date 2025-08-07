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
import 'highlight.js/styles/github.css'

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
  marked.setOptions({
    highlight: function(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext';
      return hljs.highlight(code, { language }).value;
    },
    langPrefix: 'hljs language-',
  });
  
  // 首先使用marked解析Markdown
  let html: string = marked.parse(contentWithIds) as string
  
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
    try {
      return katex.renderToString(formula.trim(), { displayMode: true })
    } catch (error) {
      console.error('KaTeX渲染错误:', error)
      return `<span style="color: red;">公式渲染错误: ${formula}</span>`
    }
  })
  
  // 处理行内公式（$...$）
  html = html.replace(/\$([^\$\n]+?)\$/g, (_, formula) => {
    try {
      return katex.renderToString(formula.trim(), { displayMode: false })
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
const handleContentClick = (event: MouseEvent) => {
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
  border-left: 4px solid #007bff;
  padding: 10px 20px;
  margin: 0 0 30px 0;
  background: #f8f9fa;
}

.content {
  font-size: 16px;
  line-height: 1.8;
}

.content :deep(h1) {
  font-size: 24px;
  margin: 24px 0 16px;
  color: #333;
}

.content :deep(h2) {
  font-size: 22px;
  margin: 22px 0 14px;
  color: #333;
}

.content :deep(h3) {
  font-size: 20px;
  margin: 20px 0 12px;
  color: #333;
}

.content :deep(p) {
  margin: 16px 0;
  color: #444;
}

.content :deep(ul),
.content :deep(ol) {
  padding-left: 30px;
  margin: 16px 0;
}

.content :deep(li) {
  margin-bottom: 8px;
}

.content :deep(code) {
  background: #f1f1f1;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
  font-size: 14px;
}

.content :deep(pre) {
  background: #f8f9fa;
  padding: 16px;
  border-radius: 6px;
  overflow: auto;
  margin: 20px 0;
}

.content :deep(pre code) {
  background: none;
  padding: 0;
}

.content :deep(blockquote) {
  border-left: 4px solid #ddd;
  padding: 10px 20px;
  margin: 20px 0;
  background: #f8f9fa;
  color: #666;
}

.content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 20px 0;
  cursor: pointer;
  transition: opacity 0.3s;
}

.content :deep(img:hover) {
  opacity: 0.8;
}

.content :deep(a) {
  color: #007bff;
  text-decoration: none;
}

.content :deep(a:hover) {
  text-decoration: underline;
}

.content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: left;
}

.content :deep(th) {
  background: #f8f9fa;
  font-weight: bold;
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
  position: relative;
  top: -80px; /* 修正锚点位置 */
  display: block;
  height: 0;
  visibility: hidden;
}
</style>