<template>
  <div class="sidebar-card tag-cloud">
    <div class="card-header">
      <h3><i class="iconfont icon-tags tag-icon"></i> 标签云</h3>
      <div class="header-actions">
        <div class="zoom-indicator" v-if="is3DView && Math.abs(zoomLevel - 1.0) > 0.01">
          <span>{{ Math.round(zoomLevel * 100) }}%</span>
        </div>
        <div class="view-switch" @click="toggleView" :title="is3DView ? '切换到列表视图' : '切换到3D视图'">
          <span class="switch-label">{{ is3DView ? '3D' : '列表' }}</span>
          <div class="switch-track" :class="{ 'active': is3DView }">
            <div class="switch-thumb"></div>
          </div>
        </div>
      </div>
    </div>
    <div class="card-content" :class="{ 'content-3d': is3DView }">
      <div v-if="loading" class="skeleton-loader">
        <div class="skeleton-header"></div>
        <div class="skeleton-body">
          <div class="skeleton-tag"></div>
          <div class="skeleton-tag"></div>
          <div class="skeleton-tag"></div>
        </div>
      </div>
      
      <!-- 列表视图 -->
      <div class="tags" v-else-if="!is3DView && tags.length > 0">
        <router-link
          v-for="tag in displayedTags" 
          :key="tag.name" 
          :to="`/articles?keyword=${tag.name}`"
          class="tag"
          :style="{ fontSize: calculateFontSize(tag.count) }"
        >
          {{ tag.name }} <span class="tag-count">({{ tag.count }})</span>
        </router-link>
        
        <div class="see-more-container">
          <button 
            v-if="tags.length > TAG_LIMIT" 
            @click="toggleShowAll"
            class="see-more-button"
          >
            <i class="iconfont icon-seemore"></i>
            <span>{{ showAll ? '收起' : '查看更多' }}</span>
          </button>
        </div>
      </div>

      <!-- 3D 视图 -->
      <div 
        v-else-if="is3DView && tags.length > 0" 
        class="tags-3d-container"
        ref="containerRef"
        @mousedown="handleMouseDown"
        @mousemove="handleMouseMove3D"
        @mouseup="handleMouseUp"
        @mouseleave="handleMouseLeave"
        @touchstart="handleTouchStart"
        @touchmove="handleTouchMove"
        @touchend="handleTouchEnd"
        @wheel="handleWheel"
      >
        <canvas class="connections-canvas" ref="canvasRef"></canvas>
        <div class="tags-3d-wrapper" ref="wrapperRef">
          <router-link
            v-for="(tag, index) in tags3D" 
            :key="tag.name" 
            :to="`/articles?keyword=${tag.name}`"
            class="tag-3d"
            :style="tag.style"
            :data-index="index"
          >
            {{ tag.name }} <span class="tag-count">({{ tag.count }})</span>
          </router-link>
        </div>
      </div>

      <div class="error-message" v-else-if="error">
        <p>{{ error }}</p>
        <button @click="onRetry" class="retry-button">重试</button>
      </div>
      <div class="empty-state" v-else>
        <p>暂无标签</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { tagApi } from '@/services/api'

// 定义标签接口
interface Tag {
  name: string
  count: number
}

interface Tag3D extends Tag {
  x: number
  y: number
  z: number
  style: any
}

// 标签显示上限 (列表视图)
const TAG_LIMIT = 8

const tags = ref<Tag[]>([])
const tags3D = ref<Tag3D[]>([])
const loading = ref(false)
const error = ref('')
const showAll = ref(false)
const is3DView = ref(true) // 默认3D视图

// 3D 动画相关
const containerRef = ref<HTMLElement | null>(null)
const wrapperRef = ref<HTMLElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationId: number | null = null
let resizeObserver: ResizeObserver | null = null
const baseRadius = 100 // 基础球体半径
let radius = 100 // 当前球体半径（可动态调整）
const baseSpeed = 0.005 // 基础旋转速度
let angleX = baseSpeed
let angleY = baseSpeed
let mouseX = 0
let mouseY = 0
let isDragging = false
let lastX = 0
let lastY = 0
let zoomLevel = 1.0 // 缩放级别（0.5 - 2.0）

// 缓存主题色，避免每帧都读取 DOM
let cachedAccentColor = ''
let cachedTextColor = ''
const refreshThemeColors = () => {
  const style = getComputedStyle(document.documentElement)
  cachedAccentColor = style.getPropertyValue('--color-accent').trim()
  cachedTextColor = style.getPropertyValue('--color-text-secondary').trim()
}

// 监听主题变化
const themeObserver = new MutationObserver(() => {
  refreshThemeColors()
})

// 计算显示的标签 (列表视图)
const displayedTags = computed(() => {
  if (showAll.value) {
    return tags.value
  }
  return tags.value.slice(0, TAG_LIMIT)
})

// 定义事件
const emit = defineEmits<{
  (e: 'loading', value: boolean): void
}>()

// 切换视图
const toggleView = () => {
  is3DView.value = !is3DView.value
  if (is3DView.value) {
    nextTick(() => initAfterLayout())
  } else {
    stopAnimation()
    if (resizeObserver) {
      resizeObserver.disconnect()
      resizeObserver = null
    }
  }
}

// 等待浏览器完成布局后初始化 3D（双重 rAF 确保 layout 已计算）
const initAfterLayout = () => {
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      init3D()
      initCanvas()
      startAnimation()
      refreshThemeColors()
      setupResizeObserver()
    })
  })
}

// 初始化 Canvas
const initCanvas = () => {
  if (!canvasRef.value || !containerRef.value) return
  
  const container = containerRef.value
  canvasRef.value.width = container.offsetWidth
  canvasRef.value.height = container.offsetHeight
}

// 监听容器尺寸变化，重设 Canvas
const setupResizeObserver = () => {
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
  if (!containerRef.value) return
  
  resizeObserver = new ResizeObserver(() => {
    initCanvas()
  })
  resizeObserver.observe(containerRef.value)
}

// 初始化 3D 坐标 (斐波那契球分布)
const init3D = () => {
  const len = tags.value.length
  
  // 根据标签数量动态调整半径
  if (len <= 10) {
    radius = baseRadius * 0.8 // 标签少时缩小
  } else if (len <= 30) {
    radius = baseRadius // 标准大小
  } else if (len <= 60) {
    radius = baseRadius * 1.3 // 标签多时扩大
  } else {
    radius = baseRadius * 1.6 // 标签很多时更大
  }
  
  // 应用用户缩放级别
  radius = radius * zoomLevel
  
  tags3D.value = tags.value.map((tag, i) => {
    const phi = Math.acos(-1 + (2 * i + 1) / len)
    const theta = Math.sqrt(len * Math.PI) * phi
    return {
      ...tag,
      x: radius * Math.cos(theta) * Math.sin(phi),
      y: radius * Math.sin(theta) * Math.sin(phi),
      z: radius * Math.cos(phi),
      style: {}
    }
  })
}

// 绘制连线
const drawConnections = () => {
  if (!canvasRef.value || !wrapperRef.value || !containerRef.value) return
  
  const ctx = canvasRef.value.getContext('2d')
  if (!ctx) return
  
  ctx.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height)
  
  const containerRect = containerRef.value.getBoundingClientRect()
  const centerX = containerRect.width / 2
  const centerY = containerRect.height / 2
  
  const tagElements = wrapperRef.value.querySelectorAll('.tag-3d')
  
  // 使用主题色绘制连线
  ctx.strokeStyle = cachedAccentColor || 'rgb(66, 184, 131)'
  ctx.lineWidth = 1.2
  ctx.lineCap = 'round'
  
  tagElements.forEach((el, index) => {
    const rect = el.getBoundingClientRect()
    const x = rect.left - containerRect.left + rect.width / 2
    const y = rect.top - containerRect.top + rect.height / 2
    
    const zDepth = tags3D.value[index].z
    const depthFactor = Math.max(0.3, (zDepth + radius) / (2 * radius))
    ctx.globalAlpha = 0.3 * depthFactor
    
    ctx.beginPath()
    ctx.moveTo(centerX, centerY)
    ctx.lineTo(x, y)
    ctx.stroke()
  })
  
  ctx.globalAlpha = 1.0
}

// 3D 动画循环
const animate = () => {
  tags3D.value.forEach(tag => {
    const cosX = Math.cos(angleX)
    const sinX = Math.sin(angleX)
    const y1 = tag.y * cosX - tag.z * sinX
    const z1 = tag.z * cosX + tag.y * sinX
    tag.y = y1
    tag.z = z1

    const cosY = Math.cos(angleY)
    const sinY = Math.sin(angleY)
    const x2 = tag.x * cosY - tag.z * sinY
    const z2 = tag.z * cosY + tag.x * sinY
    tag.x = x2
    tag.z = z2

    const scale = (2 * radius + tag.z) / (2 * radius)
    const alpha = (tag.z + radius) / (2 * radius)
    
    const safeScale = Math.max(0.5, Math.min(1.5, scale))
    const safeAlpha = Math.max(0.3, Math.min(1, alpha + 0.3))

    tag.style = {
      transform: `translate3d(${tag.x}px, ${tag.y}px, 0) scale(${safeScale})`,
      opacity: safeAlpha,
      zIndex: Math.floor(tag.z),
      fontSize: calculateFontSize(tag.count)
    }
  })
  
  drawConnections()
  animationId = requestAnimationFrame(animate)
}

const startAnimation = () => {
  if (!animationId) {
    animate()
  }
}

const stopAnimation = () => {
  if (animationId) {
    cancelAnimationFrame(animationId)
    animationId = null
  }
}

// 鼠标/触摸事件处理
const handleMouseDown = (e: MouseEvent) => {
  isDragging = true
  lastX = e.clientX
  lastY = e.clientY
  containerRef.value!.style.cursor = 'grabbing'
}

const handleMouseMove3D = (e: MouseEvent) => {
  if (!isDragging) {
    if (!is3DView.value) return
    if (!containerRef.value) return
    const rect = containerRef.value.getBoundingClientRect()
    mouseX = (e.clientX - rect.left - rect.width / 2) / (rect.width / 2)
    mouseY = (e.clientY - rect.top - rect.height / 2) / (rect.height / 2)
    angleY = mouseX * 0.02
    angleX = -mouseY * 0.02
    return
  }
  
  const deltaX = e.clientX - lastX
  const deltaY = e.clientY - lastY
  angleY = deltaX * 0.005
  angleX = -deltaY * 0.005
  lastX = e.clientX
  lastY = e.clientY
}

const handleMouseUp = () => {
  if (isDragging) {
    isDragging = false
    containerRef.value!.style.cursor = 'grab'
    setTimeout(() => {
      if (!isDragging) {
        angleX = baseSpeed
        angleY = baseSpeed
      }
    }, 2000)
  }
}

const handleMouseLeave = () => {
  if (!isDragging) {
    angleX = baseSpeed
    angleY = baseSpeed
  }
  containerRef.value!.style.cursor = 'default'
}

// 触摸事件处理
const handleTouchStart = (e: TouchEvent) => {
  isDragging = true
  lastX = e.touches[0].clientX
  lastY = e.touches[0].clientY
  containerRef.value!.style.cursor = 'grabbing'
  e.preventDefault()
}

const handleTouchMove = (e: TouchEvent) => {
  if (!isDragging) return
  const deltaX = e.touches[0].clientX - lastX
  const deltaY = e.touches[0].clientY - lastY
  angleY = deltaX * 0.005
  angleX = -deltaY * 0.005
  lastX = e.touches[0].clientX
  lastY = e.touches[0].clientY
  e.preventDefault()
}

const handleTouchEnd = () => {
  if (isDragging) {
    isDragging = false
    containerRef.value!.style.cursor = 'grab'
    setTimeout(() => {
      if (!isDragging) {
        angleX = baseSpeed
        angleY = baseSpeed
      }
    }, 2000)
  }
}

// 鼠标滚轮缩放
const handleWheel = (e: WheelEvent) => {
  e.preventDefault()
  
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  zoomLevel = Math.max(0.5, Math.min(2.0, zoomLevel + delta))
  
  // 重新初始化 3D 布局
  init3D()
}

// 获取标签列表
const fetchTags = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    const response = await tagApi.getTags({
      pagesize: 100,
      pagenum: 1
    })
    
    if (response.status !== 200) {
      error.value = '网络请求失败'
      return
    }
    
    const { data, status, message } = response.data
    
    if (status !== 200) {
      error.value = message || '获取标签列表失败'
      return
    }
    
    tags.value = data.map((tag: any) => ({
      name: tag.name,
      count: tag.count
    })).sort((a: any, b: any) => b.count - a.count)
    
    // 标签加载完成后，如果当前是3D视图则初始化
    if (is3DView.value) {
      nextTick(() => initAfterLayout())
    }
    
  } catch (err: any) {
    error.value = err.message || '获取标签失败'
  } finally {
    loading.value = false
    emit('loading', false)
  }
}

// 计算标签字体大小
const calculateFontSize = (count: number) => {
  if (count <= 1) return '12px'
  if (count <= 3) return '14px'
  if (count <= 5) return '16px'
  if (count <= 10) return '18px'
  return '20px'
}

// 切换显示全部/部分标签
const toggleShowAll = () => {
  showAll.value = !showAll.value
}

// 重试函数
const onRetry = () => {
  fetchTags()
}

// 暴露方法给父组件
defineExpose({
  fetchTags
})

// 组件挂载时获取数据
onMounted(() => {
  fetchTags()
  // 监听主题切换（data-theme 属性变化）
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })
})

onBeforeUnmount(() => {
  stopAnimation()
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  themeObserver.disconnect()
})
</script>

<style scoped>
.iconfont {
  font-size: 10px;
}

.tag-icon {
  color: var(--color-accent);
  margin-right: 5px;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid var(--color-border);
  background: transparent;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--color-heading);
  position: relative;
  padding-left: 12px;
}

.card-header h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 16px;
  background: var(--color-accent);
  border-radius: 2px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.zoom-indicator {
  padding: 2px 8px;
  background: color-mix(in srgb, var(--color-accent) 15%, transparent);
  border-radius: 12px;
  font-size: 11px;
  color: var(--color-accent);
  font-weight: 600;
  transition: all 0.3s ease;
}

.view-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.switch-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.switch-track {
  width: 36px;
  height: 20px;
  background-color: var(--color-border);
  border-radius: 10px;
  position: relative;
  transition: background-color 0.3s ease;
}

.switch-track.active {
  background-color: var(--color-accent);
}

.switch-thumb {
  width: 16px;
  height: 16px;
  background-color: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.switch-track.active .switch-thumb {
  transform: translateX(16px);
}

.card-content {
  height: 300px;
  position: relative;
  overflow-y: auto;
  scrollbar-width: thin;
  transition: all 0.3s ease;
}

.card-content::-webkit-scrollbar {
  width: 4px;
}

.card-content::-webkit-scrollbar-thumb {
  background-color: var(--color-border-hover);
  border-radius: 2px;
}

.content-3d {
  overflow: hidden;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 15px;
}

.tag {
  display: inline-block;
  padding: 6px 14px;
  background: color-mix(in srgb, var(--color-accent) 10%, transparent);
  border-radius: 20px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  text-decoration: none;
  font-weight: 500;
  line-height: 1.2;
  border: 1px solid transparent;
}

.tag:hover {
  background: var(--color-accent);
  color: white;
  transform: translateY(-3px);
  box-shadow: 0 4px 10px color-mix(in srgb, var(--color-accent) 30%, transparent);
}

/* 3D 视图样式 */
.tags-3d-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  perspective: 800px;
  cursor: grab;
  position: relative;
}

.tags-3d-container:active {
  cursor: grabbing;
}

.connections-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.tags-3d-wrapper {
  position: relative;
  width: 0;
  height: 0;
  transform-style: preserve-3d;
}

.tag-3d {
  position: absolute;
  left: 0;
  top: 0;
  color: var(--color-text-secondary);
  text-decoration: none;
  font-weight: bold;
  white-space: nowrap;
  transform-origin: center center;
  will-change: transform, opacity;
  text-shadow: 0 1px 2px var(--color-background);
  transition: color 0.3s ease;
}

.tag-3d:hover {
  color: var(--color-accent) !important;
  z-index: 1000 !important;
  text-shadow: 0 2px 4px var(--color-shadow);
}

.empty-state {
  text-align: center;
  padding: 30px 10px;
  color: var(--color-text-secondary);
}

.skeleton-loader {
  animation: skeleton-loading 1s linear infinite alternate;
  padding: 15px;
}

@keyframes skeleton-loading {
  0% {
    background-color: var(--color-background-soft);
  }
  100% {
    background-color: var(--color-background-mute);
  }
}

.skeleton-header {
  height: 20px;
  width: 60%;
  margin-bottom: 15px;
  border-radius: 4px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.skeleton-tag {
  height: 24px;
  width: 60px;
  border-radius: 12px;
}

.error-message {
  text-align: center;
  padding: 30px 10px;
  color: #dc3545;
}

.retry-button {
  margin-top: 15px;
  padding: 8px 16px;
  background-color: var(--color-accent);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.retry-button:hover {
  opacity: 0.85;
  border-radius: 8px;
}

.see-more-container {
  display: flex;
  justify-content: center;
  width: 100%;
  margin-top: 10px;
}

.see-more-button {
  padding: 8px 15px;
  background-color: transparent;
  color: var(--color-accent);
  border: 1px solid var(--color-accent);
  border-radius: 20px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s ease;
}

.see-more-button:hover {
  background-color: var(--color-accent);
  color: white;
}

.tag-count {
  font-size: 0.8em;
  opacity: 0.8;
  margin-left: 2px;
}
</style>
