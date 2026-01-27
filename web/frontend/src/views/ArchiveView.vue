<template>
  <MainLayout>
    <template #main>
      <div class="archive-page">
        <div class="page-header">
          <h1><i class="iconfont icon-archive"></i> 文章归档</h1>
          <p class="subtitle">共 {{ totalArticles }} 篇文章，继续加油！</p>
        </div>

        <div class="tag-filter-section" v-if="tags.length > 0">
          <div class="tag-list">
             <span 
              class="tag-item" 
              :class="{ active: selectedTag === '' }"
              @click="selectTag('')"
            >
              全部
            </span>
            <span 
              v-for="tag in visibleTags" 
              :key="tag.name" 
              class="tag-item" 
              :class="{ active: selectedTag === tag.name }"
              @click="selectTag(tag.name)"
            >
              {{ tag.name }} <span class="count">({{ tag.count }})</span>
            </span>
            <span v-if="tags.length > 15" class="tag-toggle" @click="toggleTags">
               {{ isTagsExpanded ? '收起' : '更多' }}
            </span>
          </div>
        </div>

        <!-- 贡献热力图 -->
        <div class="contribution-section" v-if="!selectedTag && !loading && articles.length > 0">
           <!-- 左侧统计栏 -->
           <div class="contrib-stats">
              <div class="stat-item">
                  <span class="stat-label">总文章数</span>
                  <span class="stat-value">{{ articles.length }} 篇</span>
              </div>
              <div class="stat-item">
                  <span class="stat-label">最长连续</span>
                  <span class="stat-value">{{ longestStreak }} 天</span>
              </div>
              <div class="stat-item">
                  <span class="stat-label">最近更新</span>
                  <span class="stat-value">{{ lastUpdate }}</span>
              </div>
           </div>

           <!-- 右侧日历 -->
           <div class="calendar-wrapper">
              <!-- Month Labels -->
              <div class="months-row">
                 <span 
                  v-for="label in monthLabels" 
                  :key="label.text + label.index" 
                  class="month-label" 
                  :style="{ left: `${30 + label.index * 15}px` }" 
                 >
                   {{ label.text }}
                 </span>
              </div>

              <div class="graph-row">
                 <!-- Weekday Labels -->
                 <div class="weekdays-col">
                     <span></span>
                     <span>Mon</span>
                     <span></span>
                     <span>Wed</span>
                     <span></span>
                     <span>Fri</span>
                     <span></span>
                 </div>
                 <!-- The Grid -->
                 <div class="columns-container">
                    <div v-for="(week, wIndex) in calendarWeeks" :key="wIndex" class="week-column">
                        <div v-for="day in week" 
                             :key="day.date" 
                             class="day-cell"
                             :class="`level-${day.level}`"
                             :title="`${day.date}: ${day.count} 篇`"
                        ></div>
                    </div>
                 </div>
              </div>
              
              <!-- Legend -->
              <div class="calendar-footer">
                  <span class="legend-text">Less</span>
                  <div class="legend-scale">
                      <div class="day-cell level-0"></div>
                      <div class="day-cell level-1"></div>
                      <div class="day-cell level-2"></div>
                      <div class="day-cell level-3"></div>
                      <div class="day-cell level-4"></div>
                  </div>
                  <span class="legend-text">More</span>
              </div>
           </div>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载归档中...</p>
        </div>

        <template v-else>
          <!-- 空状态 -->
          <div v-if="filteredArticles.length === 0" class="empty-state">
             <i class="iconfont icon-coffee" style="font-size: 48px; opacity: 0.5; margin-bottom: 10px;"></i>
             <p>这里暂时没有文章哦~</p>
          </div>

          <div class="archive-content" v-else>
            <div class="timeline">
              <TransitionGroup name="list">
                <div v-for="(group, year) in groupedArticles" :key="year" :id="`year-${year}`" class="timeline-year">
            <h2 class="year-title">{{ year }}</h2>
            
            <div v-for="(articles, month) in group" :key="month" class="timeline-month">
              <h3 class="month-title">{{ month }}月</h3>
              
              <div class="timeline-items">
                <div v-for="article in articles" :key="article.id || article.ID" class="timeline-item">
                  <span class="date">{{ formatDateDay(article.createdAt || article.CreatedAt) }}</span>
                  <router-link :to="`/article/${article.id || article.ID}`" class="title">
                    {{ article.title }}
                  </router-link>
                </div>
              </div>
                </div>
                </div>
              </TransitionGroup>
            </div>
            
            <div class="year-nav">
               <div 
                v-for="(group, year) in groupedArticles" 
                :key="year"
                class="year-nav-item"
                @click="scrollToYear(year as string)"
               >
                 {{ year }}
               </div>
            </div>
          </div>
        </template>
      </div>
    </template>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/components/layout/MainLayout.vue'
import { articleApi, tagApi } from '@/services/api'

interface Article {
  id?: number
  ID?: number
  title: string
  createdAt?: string
  CreatedAt?: string
  tags?: string // "tag1, tag2"
  tag_models?: { id: number; name: string }[]
}

interface Tag {
  name: string
  count: number
}

const loading = ref(false)
const articles = ref<Article[]>([])
const tags = ref<Tag[]>([])
const totalArticles = ref(0)
const selectedTag = ref('')
const isTagsExpanded = ref(false)

// 贡献图数据
const contributionData = computed(() => {
  const map = new Map<string, number>()
  articles.value.forEach(art => {
    const d = art.createdAt || art.CreatedAt
    if (d) {
      const day = d.split('T')[0] // 假设 ISO 格式 YYYY-MM-DD
      map.set(day, (map.get(day) || 0) + 1)
    }
  })
  return map
})

// 生成类似 Github 的周视图数据 (52列 x 7行)
const calendarWeeks = computed(() => {
  const weeks: { date: string; count: number; level: number }[][] = []
  const end = new Date()
  const start = new Date()
  start.setFullYear(start.getFullYear() - 1)
  
  // 调整到上一年的那个周日，确保第一列从周日开始
  while (start.getDay() !== 0) {
    start.setDate(start.getDate() - 1)
  }

  let current = new Date(start)
  
  // 生成53周的数据，填满一行
  for (let w = 0; w < 53; w++) {
    const week: { date: string; count: number; level: number }[] = []
    for (let i = 0; i < 7; i++) {
        const dateStr = current.toISOString().split('T')[0]
        const count = contributionData.value.get(dateStr) || 0
        
        let level = 0
        if (count > 0) level = 1
        if (count > 1) level = 2
        if (count > 2) level = 3
        if (count > 4) level = 4
        
        // 只有当日期在end之前（或者就是end）才显示有效，否则是未来的日期
        // GitHub如果这一行没满（今天周三），后面的格子是不显示的还是空的？
        // 观察GitHub，最后一周如果不满7天，后面的格子是不渲染的或者不可见的。
        // 这里我们可以保留格子但是level-0，或者判断日期
        
        if (current > end) {
           // 未来的日子，不渲染或者特殊处理
           // 这里还是push进去，但是可以通过样式控制或者就留白
        }

        week.push({ date: dateStr, count, level })
        current.setDate(current.getDate() + 1)
    }
    weeks.push(week)
  }
  return weeks
})

// 生成月份标签
const monthLabels = computed(() => {
  const labels: { text: string; index: number }[] = []
  
  calendarWeeks.value.forEach((week, index) => {
    const firstDay = new Date(week[0].date)
    if (firstDay.getDate() <= 7) {
       const monthText = firstDay.toLocaleString('default', { month: 'short' })
       if (labels.length === 0 || labels[labels.length - 1].text !== monthText) {
          labels.push({ text: monthText, index })
       }
    }
  })
  return labels
})

const lastUpdate = computed(() => {
    if (articles.value.length === 0) return '暂无'
    // 假设文章已按时间倒序排列，否则需自行排序
    const sorted = [...articles.value].sort((a, b) => {
        const dA = new Date(a.createdAt || a.CreatedAt || 0).getTime()
        const dB = new Date(b.createdAt || b.CreatedAt || 0).getTime()
        return dB - dA
    })
    const latest = sorted[0].createdAt || sorted[0].CreatedAt
    if (!latest) return '未知'
    return latest.split('T')[0]
})

const longestStreak = computed(() => {
    const dates = Array.from(contributionData.value.keys()).sort()
    if (dates.length === 0) return 0
    
    let maxStreak = 1
    let currentStreak = 1
    
    for (let i = 1; i < dates.length; i++) {
        const prev = new Date(dates[i-1])
        const curr = new Date(dates[i])
        const diffTime = Math.abs(curr.getTime() - prev.getTime())
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
        
        if (diffDays === 1) {
            currentStreak++
        } else {
            maxStreak = Math.max(maxStreak, currentStreak)
            currentStreak = 1
        }
    }
    return Math.max(maxStreak, currentStreak)
})


// 展示的标签
const visibleTags = computed(() => {
  if (isTagsExpanded.value) {
    return tags.value
  }
  return tags.value.slice(0, 15) //只显示前15个
})

// 获取所有标签
const fetchTags = async () => {
  try {
    const res = await tagApi.getTags({ pagesize: 100, pagenum: 1 })
    if (res.status === 200 && res.data.status === 200) {
      tags.value = res.data.data
    }
  } catch (error) {
    console.error('Failed to fetch tags:', error)
  }
}

// 获取归档数据
const fetchArchive = async () => {
  loading.value = true
  try {
    // 获取足够多的文章用于归档
    const res = await articleApi.getArticles({ pagesize: 1000, pagenum: 1 })
    if (res.status === 200 && res.data.status === 200) {
      articles.value = res.data.data
      totalArticles.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to fetch archive:', error)
  } finally {
    loading.value = false
  }
}

// 过滤后的文章
const filteredArticles = computed(() => {
  if (!selectedTag.value) {
    return articles.value
  }
  return articles.value.filter(article => {
    // 优先检查 tag_models
    if (article.tag_models && article.tag_models.length > 0) {
      return article.tag_models.some(t => t.name === selectedTag.value)
    }
    // 降级检查 tags 字符串
    if (article.tags) {
      const tagsList = article.tags.split(/,|，/).map(t => t.trim())
      return tagsList.includes(selectedTag.value)
    }
    return false
  })
})

// 按年月分组 (基于筛选后的文章)
const groupedArticles = computed(() => {
  const groups: Record<string, Record<string, Article[]>> = {}
  
  filteredArticles.value.forEach(article => {
    // 兼容后端可能返回 CreatedAt 或 createdAt
    const dateStr = article.createdAt || article.CreatedAt
    if (!dateStr) return

    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return

    const year = date.getFullYear().toString()
    const month = (date.getMonth() + 1).toString()
    
    if (!groups[year]) {
      groups[year] = {}
    }
    if (!groups[year][month]) {
      groups[year][month] = []
    }
    groups[year][month].push(article)
  })
  
  // 排序：年份倒序，月份倒序
  const sortedGroups: Record<string, Record<string, Article[]>> = {}
  Object.keys(groups).sort((a, b) => Number(b) - Number(a)).forEach(year => {
    sortedGroups[year] = {}
    Object.keys(groups[year]).sort((a, b) => Number(b) - Number(a)).forEach(month => {
      sortedGroups[year][month] = groups[year][month]
    })
  })
  
  return sortedGroups
})

const formatDateDay = (dateStr?: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  return `${date.getDate()}日`
}

const selectTag = (tagName: string) => {
  if (selectedTag.value === tagName) {
    selectedTag.value = '' // 取消选择
  } else {
    selectedTag.value = tagName
  }
}

const toggleTags = () => {
  isTagsExpanded.value = !isTagsExpanded.value
}

const scrollToYear = (year: string) => {
  const el = document.getElementById(`year-${year}`)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
  }
}

onMounted(() => {
  fetchArchive()
  fetchTags()
})
</script>

<style scoped>
.archive-page {
  background: var(--color-background);
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  position: relative;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--color-border);
}

.page-header h1 {
  font-size: 28px;
  color: var(--color-heading);
  margin-bottom: 10px;
}

.subtitle {
  color: var(--color-text);
  opacity: 0.8;
}

.timeline {
  flex: 1;
  position: relative;
  padding-left: 20px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: var(--color-border);
}

.timeline-year {
  margin-bottom: 40px;
}

.year-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--color-heading);
  margin-bottom: 20px;
  position: relative;
}

.year-title::before {
  content: '';
  position: absolute;
  left: -25px;
  top: 50%;
  transform: translateY(-50%);
  width: 12px;
  height: 12px;
  background: #06bac7ff;
  border-radius: 50%;
  border: 2px solid var(--color-background);
}

.timeline-month {
  margin-left: 20px;
  margin-bottom: 20px;
}

.month-title {
  font-size: 18px;
  color: var(--color-text);
  margin-bottom: 15px;
  opacity: 0.9;
}

.timeline-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.timeline-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 10px;
  border-radius: 6px;
  transition: all 0.3s;
}

.timeline-item:hover {
  background: var(--color-background-soft);
  transform: translateX(5px);
}

.date {
  color: var(--color-text);
  opacity: 0.7;
  font-size: 14px;
  min-width: 40px;
}

.title {
  color: var(--color-heading);
  text-decoration: none;
  font-size: 16px;
  transition: color 0.3s;
}

.title:hover {
  color: #21e7ddff;
}

.tag-filter-section {
  margin-bottom: 30px;
  padding: 0 20px;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
}

.tag-item {
  padding: 6px 14px;
  background: var(--color-background-soft);
  border-radius: 20px;
  font-size: 14px;
  color: var(--color-text);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  border: 1px solid transparent;
  display: flex;
  align-items: center;
  gap: 6px;
}

.tag-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  border-color: #21e7ddff;
  color: #21e7ddff;
}

.tag-item.active {
  background: linear-gradient(135deg, #21e7ddff 0%, #06bac7ff 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(33, 231, 221, 0.3);
  border-color: transparent;
}

.tag-item .count {
  font-size: 12px;
  opacity: 0.8;
}

.tag-toggle {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: transparent;
  border-radius: 20px;
  font-size: 13px;
  color: var(--color-text);
  opacity: 0.7;
  cursor: pointer;
  transition: all 0.3s;
}
.tag-toggle:hover {
  color: #21e7ddff;
  background: var(--color-background-soft);
  opacity: 1;
}

/* Heatmap - GitHub Style */
.contribution-section {
    display: flex;
    gap: 20px;
    margin-bottom: 40px;
    padding: 24px;
    background: var(--color-background-soft);
    border-radius: 8px;
    border: 1px solid var(--color-border);
}

/* New Stats Sidebar */
.contrib-stats {
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 120px;
  border-right: 1px solid var(--color-border);
  padding-right: 20px;
  justify-content: center;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--color-text);
  opacity: 0.7;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-heading);
}

.calendar-wrapper {
  flex: 1;
  min-width: 0; /* Important for scroll within flex item */
  overflow-x: auto;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
}

.months-row {
  position: relative;
  height: 20px;
  width: 100%;
  margin-bottom: 4px;
  font-size: 12px;
  color: var(--color-text);
}

.month-label {
  position: absolute;
  top: 0;
  font-size: 10px; /* Reduced font size to match GitHub */
}

/* Container for labels + grid */
.graph-row {
  display: flex;
}

.weekdays-col {
  display: flex;
  flex-direction: column;
  gap: 3px; /* Match grid gap */
  width: 30px; /* Fixed width for alignment */
  flex-shrink: 0;
  padding-top: 0; 
  padding-right: 8px;
  /* margin-top removed */
  font-size: 10px; 
  color: var(--color-text);
  opacity: 0.7;
  line-height: 12px; /* Match cell height */
}

.weekdays-col span {
  height: 12px;
  display: block; /* Ensure height applies */
}

.columns-container {
  display: flex;
  gap: 3px;
}
.week-column {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.day-cell {
    width: 12px;
    height: 12px;
    background: var(--color-border); /* level 0 */
    border-radius: 2px;
    box-sizing: border-box;
    border: 1px solid rgba(27, 31, 35, 0.06);
    transition: all 0.1s;
}
.day-cell:hover {
  border-color: rgba(0,0,0,0.3);
  transform: scale(1.2);
}

.day-cell.level-0 { background: #ebedf0; border-color: rgba(27,31,35,0.06); }
[data-theme='dark'] .day-cell.level-0 { background: #161b22; border-color: rgba(240,246,252,0.1); }

.day-cell.level-1 { background: #9be9a8; border-color: rgba(27,31,35,0.06); }
.day-cell.level-2 { background: #40c463; border-color: rgba(27,31,35,0.06); }
.day-cell.level-3 { background: #30a14e; border-color: rgba(27,31,35,0.06); }
.day-cell.level-4 { background: #216e39; border-color: rgba(27,31,35,0.06); }

/* Theme color override (Teal) */
.day-cell.level-1 { background: #b1f1ee; }
.day-cell.level-2 { background: #5eead4; }
.day-cell.level-3 { background: #2dd4bf; }
.day-cell.level-4 { background: #0f766e; }


.calendar-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  width: 100%;
  margin-top: 8px;
  font-size: 10px;
  color: var(--color-text);
  opacity: 0.8;
}
.legend-scale {
  display: flex;
  gap: 3px;
}

/* Archive Content Layout */
.archive-content {
    display: flex;
    gap: 40px;
    position: relative;
    align-items: flex-start;
}

/* Sticky Nav */
.year-nav {
    width: 60px;
    position: sticky;
    top: 100px; /* Adjust based on your header height */
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding-left: 20px;
    border-left: 2px solid var(--color-border-soft, #eee);
}

.year-nav-item {
    cursor: pointer;
    font-size: 14px;
    color: var(--color-text-light, #999);
    transition: all 0.2s;
    font-family: inherit;
    font-weight: 500;
}
.year-nav-item:hover {
    color: #06bac7ff;
    transform: translateX(5px) scale(1.1);
}

/* Empty State */
.empty-state {
    text-align: center;
    padding: 80px 0;
    color: var(--color-text);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}
.empty-state p {
    margin-top: 15px;
    font-size: 16px;
    opacity: 0.8;
}

/* List Transitions */
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

.list-leave-active {
  position: absolute; 
  width: 100%;
}

.loading-state {
  text-align: center;
  padding: 40px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: #04aa94ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 15px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
