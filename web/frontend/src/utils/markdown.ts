import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

// ============================================================
// Marked v16 全局配置 — Mac 风格代码块渲染器
// 在 main.ts 中导入，确保所有页面的 marked.parse() 都使用此样式
// ============================================================

marked.use({
  renderer: {
    code(token: { text: string; lang?: string; type: string; raw: string }): string | false {
      const codeText = token.text || ''
      const codeLang = token.lang || 'plaintext'

      if (!codeText || codeText.trim().length === 0) {
        return '<pre><code>代码内容为空</code></pre>'
      }

      // Mermaid 图表代码块保持原样，由前端 onMounted 中的 mermaid.render() 处理
      if (codeLang === 'mermaid') {
        const mermaidId = 'mermaid-' + Math.random().toString(36).substring(2, 9)
        return `<pre class="mermaid-chart" id="${mermaidId}">${codeText}</pre>`
      }

      const cleanCode = codeText.replace(/\n$/, '')
      const encodedCode = cleanCode
        .replace(/&/g, '&amp;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;')
      const lines = cleanCode.split('\n')
      const lineNumbers = lines.map((_, i) => `<span>${i + 1}</span>`).join('')

      // 代码高亮
      let highlighted = ''
      try {
        if (codeLang && hljs.getLanguage(codeLang)) {
          highlighted = hljs.highlight(cleanCode, { language: codeLang }).value
        } else {
          highlighted = hljs.highlightAuto(cleanCode).value
        }
      } catch {
        highlighted = cleanCode
      }

      return `
<div class="mac-code-block">
  <div class="mac-header">
    <div class="mac-dots">
      <span class="dot-red"></span>
      <span class="dot-yellow"></span>
      <span class="dot-green"></span>
    </div>
    <span class="mac-lang">${codeLang}</span>
  </div>
  <div class="mac-code-body">
    <div class="mac-line-numbers">${lineNumbers}</div>
    <div class="mac-code-content">
      <button class="mac-copy-btn" data-code="${encodedCode}">
        <svg class="copy-icon" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
        <svg class="copied-icon" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2"><polyline points="20 6 9 17 4 12"/></svg>
        <span>复制</span>
      </button>
      <pre><code class="hljs language-${codeLang}">${highlighted}</code></pre>
    </div>
  </div>
</div>`
    }
  }
})
