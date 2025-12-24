<template>
  <div class="sidebar-card tag-cloud">
    <div class="card-header">
      <h3><i class="iconfont icon-tags" style="color: #2aa073ff; margin-right: 5px;"></i> 标签云</h3>
      <div class="view-switch" @click="toggleView" :title="is3DView ? '切换到列表视图' : '切换到3D视图'">
        <span class="switch-label">{{ is3DView ? '3D' : '列表' }}</span>
        <div class="switch-track" :class="{ 'active': is3DView }">
          <div class="switch-thumb"></div>
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
        @mousemove="handleMouseMove"
        @mouseleave="handleMouseLeave"
      >
        <div class="tags-3d-wrapper">
          <router-link
            v-for="(tag, index) in tags3D" 
            :key="tag.name" 
            :to="`/articles?keyword=${tag.name}`"
            class="tag-3d"
            :style="tag.style"
          >
            {{ tag.name }} <span class="tag-count">({{ tag.count }})</span>
          </router-link>
        </div>
      </div>

      <div class="error-message" v-else-if="error">
        <p>❌ {{ error }}</p>
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
import { articleApi } from '@/services/api'

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
const is3DView = ref(false)

// 3D 动画相关
const containerRef = ref<HTMLElement | null>(null)
let animationId: number | null = null
const radius = 100 // 球体半径
const baseSpeed = 0.005 // 基础旋转速度
let angleX = baseSpeed
let angleY = baseSpeed
let mouseX = 0
let mouseY = 0

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
    nextTick(() => {
      init3D()
      startAnimation()
    })
  } else {
    stopAnimation()
  }
}

// 初始化 3D 坐标 (斐波那契球分布)
const init3D = () => {
  const len = tags.value.length
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

// 3D 动画循环
const animate = () => {
  tags3D.value.forEach(tag => {
    // 绕X轴旋转
    const cosX = Math.cos(angleX)
    const sinX = Math.sin(angleX)
    const y1 = tag.y * cosX - tag.z * sinX
    const z1 = tag.z * cosX + tag.y * sinX
    tag.y = y1
    tag.z = z1

    // 绕Y轴旋转
    const cosY = Math.cos(angleY)
    const sinY = Math.sin(angleY)
    const x2 = tag.x * cosY - tag.z * sinY
    const z2 = tag.z * cosY + tag.x * sinY
    tag.x = x2
    tag.z = z2

    // 更新样式
    const scale = (2 * radius + tag.z) / (2 * radius) // 简单的透视
    const alpha = (tag.z + radius) / (2 * radius)
    
    // 限制 scale 和 opacity 范围，防止过大或过小
    const safeScale = Math.max(0.5, Math.min(1.5, scale))
    const safeAlpha = Math.max(0.3, Math.min(1, alpha + 0.3))

    tag.style = {
      transform: `translate3d(${tag.x}px, ${tag.y}px, 0) scale(${safeScale})`,
      opacity: safeAlpha,
      zIndex: Math.floor(tag.z),
      fontSize: calculateFontSize(tag.count) // 保持原有的大小逻辑
    }
  })
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

// 鼠标交互
const handleMouseMove = (e: MouseEvent) => {
  if (!containerRef.value) return
  const rect = containerRef.value.getBoundingClientRect()
  // 计算鼠标相对于容器中心的坐标 (-1 到 1)
  mouseX = (e.clientX - rect.left - rect.width / 2) / (rect.width / 2)
  mouseY = (e.clientY - rect.top - rect.height / 2) / (rect.height / 2)
  
  // 根据鼠标位置调整旋转速度和方向
  angleY = mouseX * 0.02
  angleX = -mouseY * 0.02
}

const handleMouseLeave = () => {
  // 恢复默认旋转
  angleX = baseSpeed
  angleY = baseSpeed
}

// 获取标签列表（通过获取文章列表聚合）
const fetchTags = async () => {
  try {
    loading.value = true
    error.value = ''
    emit('loading', true)
    
    // 获取最新100篇文章来生成标签云
    const response = await articleApi.getArticles({
      pagesize: 100,
      pagenum: 1
    })
    
    // 检查响应状态
    if (response.status !== 200) {
      error.value = '网络请求失败'
      return
    }
    
    const { data, status, message } = response.data
    
    if (status !== 200) {
      error.value = message || '获取文章列表失败'
      return
    }
    
    const articles = data
    const tagMap = new Map<string, number>()
    
    articles.forEach((article: any) => {
      if (article.tags) {
        // 支持中文逗号和英文逗号
        const articleTags = article.tags.replace(/，/g, ',').split(',')
        articleTags.forEach((tag: string) => {
          const trimmedTag = tag.trim()
          if (trimmedTag) {
            tagMap.set(trimmedTag, (tagMap.get(trimmedTag) || 0) + 1)
          }
        })
      }
    })
    
    tags.value = Array.from(tagMap.entries()).map(([name, count]) => ({
      name,
      count
    })).sort((a, b) => b.count - a.count)
    
    // 如果当前是3D视图，重新初始化
    if (is3DView.value) {
      init3D()
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
})

onBeforeUnmount(() => {
  stopAnimation()
})
</script>

<style scoped>
.iconfont {
  font-size: 10px;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  background: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #333;
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
  background: #3d96fc;
  border-radius: 2px;
}

.view-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.switch-label {
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.switch-track {
  width: 36px;
  height: 20px;
  background-color: #e0e0e0;
  border-radius: 10px;
  position: relative;
  transition: background-color 0.3s ease;
}

.switch-track.active {
  background-color: #42b883;
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
  min-height: 200px;
  position: relative;
}

.content-3d {
  height: 300px; /* 3D 视图需要固定高度 */
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
  background: rgba(61, 150, 252, 0.08);
  border-radius: 20px;
  color: #555;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  text-decoration: none;
  font-weight: 500;
  line-height: 1.2;
  border: 1px solid transparent;
}

.tag:hover {
  background: #42b883;
  color: white;
  transform: translateY(-3px);
  box-shadow: 0 4px 10px rgba(66, 184, 131, 0.3);
}

/* 3D 视图样式 */
.tags-3d-container {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  perspective: 800px;
  cursor: move;
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
  color: #42b883; /* 默认颜色 */
  text-decoration: none;
  font-weight: bold;
  white-space: nowrap;
  transform-origin: center center;
  will-change: transform, opacity;
  /* 移除背景色，只显示文字 */
  text-shadow: 0 1px 2px rgba(255,255,255,0.8);
}

.tag-3d:hover {
  color: #3aa876;
  z-index: 1000 !important; /* 确保 hover 时在最上层 */
  text-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.empty-state {
  text-align: center;
  padding: 30px 10px;
  color: #888;
}

.skeleton-loader {
  animation: skeleton-loading 1s linear infinite alternate;
  padding: 15px;
}

@keyframes skeleton-loading {
  0% {
    background-color: hsl(200, 20%, 80%);
  }
  100% {
    background-color: hsl(200, 20%, 95%);
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
  background-color: #42b883;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.retry-button:hover {
  background-color: #3aa876;
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
  color: #42b883;
  border: 1px solid #42b883;
  border-radius: 20px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s ease;
}

.see-more-button:hover {
  background-color: #42b883;
  color: white;
}

.tag-count {
  font-size: 0.8em;
  opacity: 0.8;
  margin-left: 2px;
}
</style>