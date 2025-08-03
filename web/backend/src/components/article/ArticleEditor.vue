<template>
  <div class="markdown-editor-container">
    <!-- 仅预览模式 -->
    <div v-if="previewOnly" class="markdown-preview-only">
      <div class="markdown-preview" v-html="renderedMarkdown" ref="previewRef"></div>
    </div>
    
    <!-- 双栏编辑模式 -->
    <div class="split-pane" ref="splitPaneRef" v-else>
      <!-- 左侧编辑区 -->
      <div 
        class="split-pane-panel" 
        :style="{ width: leftPanelWidth }"
      >
        <textarea 
          v-model="content" 
          class="markdown-textarea"
          placeholder="请输入Markdown格式的内容"
          ref="textareaRef"
          @scroll="handleTextareaScroll"
        ></textarea>
      </div>
      
      <!-- 拖动条 -->
      <div 
        class="split-pane-divider" 
        @mousedown="startDrag"
      ></div>
      
      <!-- 右侧预览区 -->
      <div 
        class="split-pane-panel preview-panel" 
        :style="{ width: rightPanelWidth }"
      >
        <div 
          class="markdown-preview" 
          v-html="renderedMarkdown"
          ref="previewRef"
          @scroll="handlePreviewScroll"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount, watch } from 'vue'
import { marked } from 'marked'
import katex from 'katex'
import 'katex/dist/katex.min.css'

// 定义组件属性
const props = defineProps<{
  modelValue: string
  previewOnly?: boolean
}>()

// 定义事件
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

// 内容值
const content = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// 仅预览模式
const previewOnly = computed(() => props.previewOnly || false)

// 预览相关
const previewRef = ref<HTMLDivElement | null>(null)

const renderedMarkdown = computed(() => {
  const markdown = content.value || ''
  // 首先使用marked解析Markdown
  let html: string = marked.parse(markdown) as string
  
  // 然后处理数学公式
  html = renderMath(html)
  
  return html
})

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

// 渲染完成后的数学公式处理
const renderMathInElement = () => {
  if (previewRef.value) {
    // 处理块级公式（$$...$$）
    const blockFormulas = previewRef.value.querySelectorAll('p')
    blockFormulas.forEach(element => {
      const text = element.textContent || ''
      if (text.startsWith('$$') && text.endsWith('$$')) {
        try {
          const formula = text.substring(2, text.length - 2)
          const katexHtml = katex.renderToString(formula.trim(), { displayMode: true })
          element.innerHTML = katexHtml
        } catch (error) {
          console.error('KaTeX渲染错误:', error)
        }
      }
    })
  }
}

// 分割面板相关
const splitPaneRef = ref<HTMLDivElement | null>(null)
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const previewRef = ref<HTMLDivElement | null>(null)

// 面板宽度状态
const leftPanelWidth = ref('50%')
const rightPanelWidth = ref('50%')
const isDragging = ref(false)
const dragStartX = ref(0)
const startLeftWidth = ref(0)

// 开始拖动
const startDrag = (e: MouseEvent) => {
  isDragging.value = true
  dragStartX.value = e.clientX
  startLeftWidth.value = parseFloat(leftPanelWidth.value)
  
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
  
  const doDrag = (e: MouseEvent) => {
    if (!isDragging.value || !splitPaneRef.value) return
    
    const containerWidth = splitPaneRef.value.clientWidth
    const deltaX = e.clientX - dragStartX.value
    let newLeftWidth = startLeftWidth.value + (deltaX / containerWidth) * 100
    
    // 限制面板宽度在20%到80%之间
    newLeftWidth = Math.max(20, Math.min(80, newLeftWidth))
    
    leftPanelWidth.value = `${newLeftWidth}%`
    rightPanelWidth.value = `${100 - newLeftWidth}%`
  }
  
  const stopDrag = () => {
    isDragging.value = false
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
    document.removeEventListener('mousemove', doDrag)
    document.removeEventListener('mouseup', stopDrag)
  }
  
  document.addEventListener('mousemove', doDrag)
  document.addEventListener('mouseup', stopDrag)
}

// 处理文本区域滚动
const handleTextareaScroll = () => {
  if (!textareaRef.value || !previewRef.value) return
  
  const textarea = textareaRef.value
  const preview = previewRef.value
  
  // 计算滚动比例
  const scrollTop = textarea.scrollTop
  const scrollHeight = textarea.scrollHeight
  const clientHeight = textarea.clientHeight
  
  // 避免除以零
  if (scrollHeight <= clientHeight) return
  
  const scrollPercentage = scrollTop / (scrollHeight - clientHeight)
  preview.scrollTop = scrollPercentage * (preview.scrollHeight - preview.clientHeight)
}

// 处理预览区域滚动
const handlePreviewScroll = () => {
  if (!textareaRef.value || !previewRef.value) return
  
  const textarea = textareaRef.value
  const preview = previewRef.value
  
  // 计算滚动比例
  const scrollTop = preview.scrollTop
  const scrollHeight = preview.scrollHeight
  const clientHeight = preview.clientHeight
  
  // 避免除以零
  if (scrollHeight <= clientHeight) return
  
  const scrollPercentage = scrollTop / (scrollHeight - clientHeight)
  textarea.scrollTop = scrollPercentage * (textarea.scrollHeight - textarea.clientHeight)
}

// 切换仅预览模式
const togglePreviewOnly = () => {
  previewOnly.value = !previewOnly.value
}

// 监听内容变化，重新渲染数学公式
watch(renderedMarkdown, () => {
  // 使用nextTick确保DOM更新后再处理数学公式
  setTimeout(() => {
    renderMathInElement()
  }, 0)
})

// 组件销毁前清理事件监听器
onBeforeUnmount(() => {
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
})
</script>

<style scoped>
/* Markdown编辑器样式 */
.markdown-editor-container {
  height: 500px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.markdown-preview-only {
  width: 100%;
  height: 100%;
  overflow: auto;
}

.split-pane {
  display: flex;
  height: 100%;
}

.split-pane-panel {
  height: 100%;
  overflow: auto;
}

.split-pane-panel textarea {
  width: 100%;
  height: 100%;
  border: none;
  outline: none;
  resize: none;
  padding: 20px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  background-color: #fafafa;
}

.split-pane-divider {
  width: 6px;
  background-color: #ebeef5;
  cursor: col-resize;
  position: relative;
  flex-shrink: 0;
}

.split-pane-divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 2px;
  width: 2px;
  height: 30px;
  background-color: #c0c4cc;
  transform: translateY(-50%);
}

.preview-panel {
  background-color: #fff;
}

.markdown-preview {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
}

/* Markdown样式 */
.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3),
.markdown-preview :deep(h4),
.markdown-preview :deep(h5),
.markdown-preview :deep(h6) {
  margin: 10px 0;
  font-weight: bold;
}

.markdown-preview :deep(h1) {
  font-size: 24px;
  border-bottom: 1px solid #dcdfe6;
  padding-bottom: 10px;
}

.markdown-preview :deep(h2) {
  font-size: 20px;
  border-bottom: 1px solid #dcdfe6;
  padding-bottom: 8px;
}

.markdown-preview :deep(h3) {
  font-size: 18px;
}

.markdown-preview :deep(h4) {
  font-size: 16px;
}

.markdown-preview :deep(p) {
  margin: 10px 0;
  line-height: 1.7;
}

.markdown-preview :deep(code) {
  background-color: #f5f7fa;
  padding: 2px 4px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
}

.markdown-preview :deep(pre) {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  overflow: auto;
  margin: 10px 0;
}

.markdown-preview :deep(pre code) {
  background-color: transparent;
  padding: 0;
}

.markdown-preview :deep(blockquote) {
  border-left: 4px solid #dcdfe6;
  padding-left: 10px;
  margin: 10px 0;
  color: #666;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  padding-left: 20px;
  margin: 10px 0;
}

.markdown-preview :deep(li) {
  margin: 5px 0;
}

.markdown-preview :deep(a) {
  color: #409eff;
  text-decoration: none;
}

.markdown-preview :deep(a:hover) {
  text-decoration: underline;
}

.markdown-preview :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  border: 1px solid #dcdfe6;
  padding: 8px 12px;
  text-align: left;
}

.markdown-preview :deep(th) {
  background-color: #f5f7fa;
  font-weight: bold;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  height: auto;
}
</style>