<template>
  <div class="article-detail-page">
    <MainLayout>
      <template #leftSidebar v-if="article && article.type !== 2 && isTocOpen">
        <!-- 仅非 PDF 文章显示目录 -->
        <ArticleToc 
          :content="article.content" 
          ref="tocRef"
          @close="isTocOpen = false"
        />
      </template>

      <template #main>
        <!-- 目录展开按钮 (当目录收起时显示) -->
        <div 
          v-if="!isTocOpen && article && article.type !== 2" 
          class="toc-fab" 
          @click="isTocOpen = true"
          title="展开目录"
        >
          <span class="fab-icon">☰</span>
        </div>

        <div class="page-header">
          <router-link to="/articles" class="back-link">
            « 返回文章列表
          </router-link>
        </div>
        
        <article class="article-content">
          <!-- 加载状态 -->
          <div v-if="loading" class="loading-state">
            <div class="spinner"></div>
            <p>加载中...</p>
          </div>
          
          <!-- 文章内容 -->
          <div v-else>
            <ArticleHeader v-if="article" :article="article" />
            
            <ArticleContent 
              v-if="article" 
              :article="article" 
              @image-click="handleImageClick"
            />
            
            <!-- 上一篇/下一篇导航 -->
            <div class="article-navigation" v-if="article">
              <div class="nav-item previous" v-if="previousArticle">
                <router-link :to="`/article/${previousArticle.id}`" class="nav-link">
                  <div class="nav-cover">
                    <img :src="previousArticle.img || defaultImage" :alt="previousArticle.title" />
                  </div>
                  <span class="nav-label">上一篇</span>
                  <span class="nav-title">{{ previousArticle.title }}</span>
                </router-link>
              </div>
              <div class="nav-item next" v-if="nextArticle">
                <router-link :to="`/article/${nextArticle.id}`" class="nav-link">
                  <div class="nav-cover">
                    <img :src="nextArticle.img || defaultImage" :alt="nextArticle.title" />
                  </div>
                  <span class="nav-label">下一篇</span>
                  <span class="nav-title">{{ nextArticle.title }}</span>
                </router-link>
              </div>
            </div>

            <!-- 相关文章推荐 -->
            <div class="related-articles-section" v-if="relatedArticles.length > 0">
              <h3 class="section-title">✨ 相关推荐</h3>
              <div class="related-grid">
                <router-link 
                  v-for="item in relatedArticles" 
                  :key="item.id" 
                  :to="`/article/${item.id}`" 
                  class="related-card"
                >
                  <div class="related-cover">
                    <img :src="item.img || defaultImage" :alt="item.title" loading="lazy" />
                  </div>
                  <div class="related-info">
                    <h4 class="related-item-title">{{ item.title }}</h4>
                    <span class="related-date">{{ formatDate(item.createdAt) }}</span>
                  </div>
                </router-link>
              </div>
            </div>

            <!-- 评论区 -->
            <GiscusComment />
            
            <div class="empty-state" v-if="!article">
              <p>文章不存在或已被删除</p>
            </div>
          </div>
        </article>
      </template>
    </MainLayout>
    
    <!-- 图片查看器 -->
    <div v-if="showImageViewer" class="image-viewer" @click="closeImageViewer">
      <div class="image-viewer-content" @click.stop>
        <button class="close-btn" @click="closeImageViewer">×</button>
        <button class="nav-btn prev-btn" @click="prevImage" v-if="imageList.length > 1">‹</button>
        <button class="nav-btn next-btn" @click="nextImage" v-if="imageList.length > 1">›</button>
        
        <div class="image-container" @wheel.prevent="handleWheel">
          <img 
            :src="currentImage" 
            :alt="`图片 ${currentImageIndex + 1}`"
            class="viewer-image"
            :style="{ transform: `scale(${imageScale})`, transition: 'transform 0.1s ease' }"
            @load="onImageLoad"
            @error="onImageError"
          />
        </div>
        
        <div class="image-info">
          <span class="image-counter">{{ currentImageIndex + 1 }} / {{ imageList.length }}</span>
          <span class="image-alt">{{ imageAltList[currentImageIndex] || `图片 ${currentImageIndex + 1}` }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { articleApi } from '@/services/api'
import MainLayout from '@/components/layout/MainLayout.vue'
import ArticleHeader from '@/components/article/ArticleHeader.vue'
import ArticleContent from '@/components/article/ArticleContent.vue'
import ArticleToc from '@/components/article/ArticleToc.vue'
import GiscusComment from '@/components/comment/GiscusComment.vue'

// 默认图片
const defaultImage = new URL('@/assets/img/无封面.jpg', import.meta.url).href

// 类型定义
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  tags: string
  views: number
  type?: number
  pdf_url?: string
  createdAt: string
  updatedAt: string
}

// 路由信息
const route = useRoute()
const tocRef = ref<InstanceType<typeof ArticleToc> | null>(null)

// 响应式数据
const article = ref<Article | null>(null)
const loading = ref(false)
const previousArticle = ref<Article | null>(null)
const nextArticle = ref<Article | null>(null)
const relatedArticles = ref<Article[]>([])
const isTocOpen = ref(true)

// 图片查看器相关
const showImageViewer = ref(false)
const imageList = ref<string[]>([])
const imageAltList = ref<string[]>([])
const currentImageIndex = ref(0)
const currentImage = ref('')
const imageScale = ref(1)

// 处理滚轮缩放
const handleWheel = (e: WheelEvent) => {
  const delta = e.deltaY > 0 ? -0.1 : 0.1
  const newScale = Math.max(0.1, Math.min(5, imageScale.value + delta))
  imageScale.value = Number(newScale.toFixed(1))
}

// 获取文章详情
const getArticleDetail = async (id: number) => {
  loading.value = true
  try {
    const response = await articleApi.getArticle(id)
    
    if (response.data.status !== 200) {
      console.error('获取文章详情失败:', response.data.message)
      article.value = null
      return
    }

    const data = response.data.data
    
    article.value = {
      id: data.ID,
      title: data.title,
      categoryId: data.cid,
      categoryName: data.Category?.name || '未分类',
      desc: data.desc,
      content: data.content,
      img: data.img,
      tags: data.tags || '',
      views: data.views || 0,
      type: data.type,
      pdf_url: data.pdf_url,
      createdAt: data.CreatedAt || data.created_at,
      updatedAt: data.UpdatedAt || data.updated_at
    }
    
    // 获取相邻文章
    await getAdjacentArticles(id)

    // 获取相关文章
    await getRelatedArticles(id)
    
    // 更新目录
    setTimeout(() => {
      if (tocRef.value) {
        tocRef.value.extractHeaders()
      }
    }, 0)
  } catch (error) {
    console.error('获取文章详情失败:', error)
    article.value = null
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 获取相关文章
const getRelatedArticles = async (id: number) => {
  try {
    const response = await articleApi.getRelatedArticles(id)
    if (response.data.status === 200) {
       relatedArticles.value = response.data.data.map((item: any) => ({
        id: item.ID,
        title: item.title,
        img: item.img,
        createdAt: item.CreatedAt || item.created_at,
       }))
    }
  } catch (error) {
    console.error('获取相关文章失败:', error)
  }
}

// 获取相邻文章
const getAdjacentArticles = async (currentId: number) => {
  try {
    // 由于后端没有提供专门的API，我们通过获取文章列表来模拟相邻文章功能
    // 这里只是简单实现，实际项目中应该由后端提供专门的API
    
    // 获取所有文章列表
    const response = await articleApi.getArticles({
      pagesize: -1,
      pagenum: -1
    })
    
    if (response.data.status !== 200) {
      return
    }

    const articles = response.data.data.map((item: any) => ({
      id: item.ID,
      title: item.title,
      categoryId: item.cid,
      categoryName: item.Category?.name || '未分类',
      desc: item.desc,
      content: item.content,
      img: item.img,
      createdAt: item.CreatedAt || item.created_at,
      updatedAt: item.UpdatedAt || item.updated_at
    }))
    
    // 按ID排序（实际项目中应该按创建时间排序）
    articles.sort((a: Article, b: Article) => a.id - b.id)
    
    // 查找当前文章的索引
    const currentIndex = articles.findIndex((article: Article) => article.id === currentId)
    
    // 设置上一篇和下一篇文章
    if (currentIndex > 0) {
      previousArticle.value = articles[currentIndex - 1]
    } else {
      previousArticle.value = null
    }
    
    if (currentIndex < articles.length - 1) {
      nextArticle.value = articles[currentIndex + 1]
    } else {
      nextArticle.value = null
    }
  } catch (error) {
    console.error('获取相邻文章失败:', error)
    previousArticle.value = null
    nextArticle.value = null
  }
}

// 处理图片点击事件
const handleImageClick = (imageSrc: string, imageAlt: string, images: string[], alts: string[]) => {
  imageList.value = images
  imageAltList.value = alts
  currentImageIndex.value = images.indexOf(imageSrc)
  currentImage.value = imageSrc
  showImageViewer.value = true
  
  // 禁止body滚动
  document.body.style.overflow = 'hidden'
}

// 关闭图片查看器
const closeImageViewer = () => {
  showImageViewer.value = false
  imageScale.value = 1
  // 恢复body滚动
  document.body.style.overflow = ''
}

// 上一张图片
const prevImage = () => {
  if (imageList.value.length <= 1) return
  currentImageIndex.value = (currentImageIndex.value - 1 + imageList.value.length) % imageList.value.length
  currentImage.value = imageList.value[currentImageIndex.value]
  imageScale.value = 1
}

// 下一张图片
const nextImage = () => {
  if (imageList.value.length <= 1) return
  currentImageIndex.value = (currentImageIndex.value + 1) % imageList.value.length
  currentImage.value = imageList.value[currentImageIndex.value]
  imageScale.value = 1
}

// 图片加载完成
const onImageLoad = () => {
  // 图片加载完成后的处理
}

// 图片加载错误
const onImageError = (e: Event) => {
  console.error('图片加载失败:', e)
}

// 监听键盘事件
const handleKeydown = (e: KeyboardEvent) => {
  if (!showImageViewer.value) return
  
  switch (e.key) {
    case 'Escape':
      closeImageViewer()
      break
    case 'ArrowLeft':
      prevImage()
      break
    case 'ArrowRight':
      nextImage()
      break
  }
}

// 监听路由变化
watch(
  () => route.params.id,
  (newId) => {
    const articleId = Number(newId)
    if (articleId) {
      getArticleDetail(articleId)
    }
  }
)

// 组件挂载时获取数据
onMounted(() => {
  const articleId = Number(route.params.id)
  if (articleId) {
    getArticleDetail(articleId)
  }
  
  // 添加键盘事件监听
  window.addEventListener('keydown', handleKeydown)
})

// 组件卸载时移除事件监听
// 注意：在Vue 3的组合式API中，需要手动清理事件监听器
</script>

<style scoped>
.article-detail-page {
  width: 100%;
  min-height: calc(100vh - 200px);
}

.page-header {
  margin-bottom: 20px;
}

.back-link {
  color: #00a5a5ff;
  text-decoration: none;
  font-size: 14px;
}

.back-link:hover {
  text-decoration: underline;
}

.article-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 30px;
  min-height: 400px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 400px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #888;
}

/* 上一篇/下一篇导航样式 */
.article-navigation {
  display: flex;
  justify-content: space-between;
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.nav-item {
  flex: 1;
  min-width: 0;
}

.nav-item.previous {
  padding-right: 10px;
}

.nav-item.next {
  padding-left: 10px;
  text-align: right;
}

.nav-link {
  display: block;
  text-decoration: none;
  color: #333;
  transition: all 0.3s ease;
}

.nav-link:hover {
  color: #007bff;
}

.nav-label {
  display: block;
  font-size: 12px;
  color: #888;
  margin-bottom: 5px;
}

.nav-title {
  display: block;
  font-size: 16px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 封面预览样式 */
.nav-cover {
  width: 100%;
  height: 120px;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 10px;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

@media (max-width: 992px) {
  .article-content {
    padding: 20px;
  }
  
  .article-navigation {
    flex-direction: column;
    gap: 20px;
  }
  
  .nav-item.previous,
  .nav-item.next {
    padding: 0;
    text-align: left;
  }
  
  .nav-item.next {
    text-align: left;
  }
  
  .nav-cover {
    height: 80px;
  }
  
  .image-viewer-content {
    max-width: 95%;
  }
  
  .nav-btn {
    width: 40px;
    height: 40px;
    font-size: 24px;
  }
  
  .prev-btn {
    left: 10px;
  }
  
  .next-btn {
    right: 10px;
  }
}

/* 图片查看器样式 */
.image-viewer {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
}

.image-viewer-content {
  position: relative;
  max-width: 90%;
  max-height: 90%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.close-btn {
  position: absolute;
  top: -40px;
  right: 0;
  background: none;
  border: none;
  color: white;
  font-size: 32px;
  cursor: pointer;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.close-btn:hover {
  color: #ccc;
}

.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  font-size: 32px;
  cursor: pointer;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: background 0.3s;
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.prev-btn {
  left: 20px;
}

.next-btn {
  right: 20px;
}

.image-container {
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 100%;
  max-height: 80vh;
}

.viewer-image {
  max-width: 100%;
  max-height: 80vh;
  object-fit: contain;
}

.image-info {
  color: white;
  margin-top: 15px;
  text-align: center;
}

.image-counter {
  margin-right: 20px;
  font-size: 14px;
}

.image-alt {
  font-size: 16px;
}

.toc-fab {
  position: fixed;
  bottom: 30px;
  left: 30px;
  width: 50px;
  height: 50px;
  background-color: white;
  border-radius: 50%;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 99;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  color: #555;
  border: 1px solid #ebeef5;
}

.toc-fab:hover {
  transform: translateY(-5px) scale(1.05);
  background-color: #42b883;
  color: white;
  box-shadow: 0 8px 24px rgba(66, 184, 131, 0.3);
  border-color: #42b883;
}

.toc-fab .fab-icon {
  font-size: 20px;
  font-weight: bold;
}

/* 相关文章样式 */
.related-articles-section {
  margin: 40px 0;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.section-title {
  font-size: 20px;
  margin-bottom: 20px;
  color: #333;
  font-weight: 600;
}

.related-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.related-card {
  display: block;
  text-decoration: none;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  transition: all 0.3s ease;
  height: 100%;
}

.related-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
}

.related-cover {
  height: 120px;
  overflow: hidden;
}

.related-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.related-card:hover .related-cover img {
  transform: scale(1.1);
}

.related-info {
  padding: 12px;
}

.related-item-title {
  margin: 0 0 8px;
  font-size: 15px;
  color: #333;
  line-height: 1.4;
  height: 42px; /* 2 lines */
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.related-date {
  font-size: 12px;
  color: #999;
}
</style>