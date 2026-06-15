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
                  <span class="stat-number">{{ articles.length }}</span>
                  <span class="stat-label">总文章数</span>
              </div>
              <div class="stat-item">
                  <span class="stat-number">{{ longestStreak }}</span>
                  <span class="stat-label">最长连续(天)</span>
              </div>
              <div class="stat-item">
                  <span class="stat-number">{{ thisYearCount }}</span>
                  <span class="stat-label">今年发布</span>
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
                  <h2 class="year-title">
                    <span class="year-badge">{{ year }}</span>
                    <span class="year-count">{{ getYearCount(group) }} 篇</span>
                  </h2>
                  
                  <div v-for="(monthArticles, month) in group" :key="month" class="timeline-month">
                    <h3 class="month-title">{{ month }}月 <span class="month-count">{{ monthArticles.length }}</span></h3>
                    
                    <div class="timeline-items">
                      <div v-for="article in monthArticles" :key="article.id || article.ID" class="timeline-item">
                        <span class="date-dot"></span>
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
                :class="{ active: activeYear === year }"
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
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
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
const activeYear = ref('')

// Helper: 格式化本地日期 YYYY-MM-DD
const toDateKey = (date: Date) => {
  const y = date.getFullYear()
  const m = (date.getMonth() + 1).toString().padStart(2, '0')
  const d = date.getDate().toString().padStart(2, '0')
  return `${y}-${m}-${d}`
}

// 贡献图数据
const contributionData = computed(() => {
  const map = new Map<string, number>()
  articles.value.forEach(art => {
    const dStr = art.createdAt || art.CreatedAt
    if (dStr) {
      const date = new Date(dStr)
      if (!isNaN(date.getTime())) {
         const key = toDateKey(date)
         map.set(key, (map.get(key) || 0) + 1)
      }
    }
  })
  return map
})

// 今年文章数
const thisYearCount = computed(() => {
  const year = new Date().getFullYear()
  return articles.value.filter(art => {
    const dStr = art.createdAt || art.CreatedAt
    if (!dStr) return false
    return new Date(dStr).getFullYear() === year
  }).length
})

// 生成类似 Github 的周视图数据 (52列 x 7行)
const calendarWeeks = computed(() => {
  const weeks: { date: string; count: number; level: number }[][] = []
  const end = new Date()
  
  const start = new Date()
  start.setDate(start.getDate() - 52 * 7)
  
  // 调整到该周的周日 (Start of the week)
  while (start.getDay() !== 0) {
    start.setDate(start.getDate() - 1)
  }

  const current = new Date(start)
  
  for (let w = 0; w < 53; w++) {
    const week: { date: string; count: number; level: number }[] = []
    for (let i = 0; i < 7; i++) {
        const dateStr = toDateKey(current)
        const count = contributionData.value.get(dateStr) || 0
        
        let level = 0
        if (count > 0) level = 1
        if (count > 2) level = 2
        if (count > 4) level = 3
        if (count > 6) level = 4
        
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
  return tags.value.slice(0, 15)
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
    const res = await articleApi.getArticles({ pagesize: -1, pagenum: -1 })
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
    if (article.tag_models && article.tag_models.length > 0) {
      return article.tag_models.some(t => t.name === selectedTag.value)
    }
    if (article.tags) {
      const tagsList = article.tags.split(/,|，/).map(t => t.trim())
      return tagsList.includes(selectedTag.value)
    }
    return false
  })
})

// 按年月分组
const groupedArticles = computed(() => {
  const groups: Record<string, Record<string, Article[]>> = {}
  
  filteredArticles.value.forEach(article => {
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
  
  const sortedGroups: Record<string, Record<string, Article[]>> = {}
  Object.keys(groups).sort((a, b) => Number(b) - Number(a)).forEach(year => {
    sortedGroups[year] = {}
    Object.keys(groups[year]).sort((a, b) => Number(b) - Number(a)).forEach(month => {
      sortedGroups[year][month] = groups[year][month]
    })
  })
  
  return sortedGroups
})

const getYearCount = (group: Record<string, Article[]>) => {
  return Object.values(group).reduce((sum, arts) => sum + arts.length, 0)
}

const formatDateDay = (dateStr?: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''
  return `${date.getDate()}日`
}

const selectTag = (tagName: string) => {
  if (selectedTag.value === tagName) {
    selectedTag.value = ''
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

// IntersectionObserver 用于年份导航高亮
let observer: IntersectionObserver | null = null

onMounted(() => {
  fetchArchive()
  fetchTags()

  // 监听年份进入视口，高亮对应的导航项
  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const id = entry.target.getAttribute('id')
        if (id) {
          activeYear.value = id.replace('year-', '')
        }
      }
    })
  }, { rootMargin: '-100px 0px -60% 0px' })
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
  }
})

// 当分组数据加载后，开始观察每个年份 section
const observeYears = () => {
  if (!observer) return
  observer.disconnect()
  const yearSections = document.querySelectorAll('.timeline-year')
  yearSections.forEach(el => observer!.observe(el))
}

// 监听 articles 变化，重新设置 Observer
watch(() => Object.keys(groupedArticles.value).length, () => {
  setTimeout(observeYears, 100)
})
</script>

<style scoped>
.archive-page {
  background: var(--color-background);
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 12px var(--color-shadow);
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

.page-header h1 .iconfont {
  color: var(--color-accent);
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 15px;
}

/* Tag Filter */
.tag-filter-section {
  margin-bottom: 30px;
  padding: 0 20px;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.tag-item {
  padding: 5px 12px;
  background: var(--color-background-soft);
  border-radius: 20px;
  font-size: 13px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.25s ease;
  border: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  gap: 5px;
  user-select: none;
}

.tag-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 8px var(--color-shadow);
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.tag-item.active {
  background: var(--color-accent);
  color: white;
  border-color: var(--color-accent);
  box-shadow: 0 3px 10px color-mix(in srgb, var(--color-accent) 30%, transparent);
}

.tag-item .count {
  font-size: 11px;
  opacity: 0.8;
}

.tag-toggle {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 5px 10px;
  background: transparent;
  border-radius: 20px;
  font-size: 12px;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.25s;
  border: 1px dashed var(--color-border);
}
.tag-toggle:hover {
  color: var(--color-accent);
  border-color: var(--color-accent);
}

/* Heatmap */
.contribution-section {
    display: flex;
    gap: 24px;
    margin-bottom: 40px;
    padding: 24px;
    background: var(--color-background-soft);
    border-radius: 8px;
    border: 1px solid var(--color-border);
}

.contrib-stats {
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-width: 110px;
  border-right: 1px solid var(--color-border);
  padding-right: 20px;
  justify-content: center;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.stat-number {
  font-size: 22px;
  font-weight: 700;
  color: var(--color-accent);
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.calendar-wrapper {
  flex: 1;
  min-width: 0;
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
  color: var(--color-text-secondary);
}

.month-label {
  position: absolute;
  top: 0;
  font-size: 10px;
}

.graph-row {
  display: flex;
}

.weekdays-col {
  display: flex;
  flex-direction: column;
  gap: 3px;
  width: 30px;
  flex-shrink: 0;
  padding-right: 8px;
  font-size: 10px; 
  color: var(--color-text-secondary);
  line-height: 12px;
}

.weekdays-col span {
  height: 12px;
  display: block;
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
    border-radius: 2px;
    box-sizing: border-box;
    transition: transform 0.15s ease;
}
.day-cell:hover {
  transform: scale(1.3);
  outline: 1px solid var(--color-accent);
}

/* Light mode heatmap */
.day-cell.level-0 { background: var(--color-background-mute); border: 1px solid var(--color-border); }
.day-cell.level-1 { background: color-mix(in srgb, var(--color-accent) 20%, var(--color-background)); }
.day-cell.level-2 { background: color-mix(in srgb, var(--color-accent) 40%, var(--color-background)); }
.day-cell.level-3 { background: color-mix(in srgb, var(--color-accent) 65%, var(--color-background)); }
.day-cell.level-4 { background: var(--color-accent); }

.calendar-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  width: 100%;
  margin-top: 8px;
  font-size: 10px;
  color: var(--color-text-secondary);
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

/* Timeline */
.timeline {
  flex: 1;
  position: relative;
  padding-left: 30px;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 6px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(to bottom, var(--color-accent), var(--color-border));
}

.timeline-year {
  margin-bottom: 40px;
}

.year-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  position: relative;
}

.year-title::before {
  content: '';
  position: absolute;
  left: -29px;
  top: 50%;
  transform: translateY(-50%);
  width: 14px;
  height: 14px;
  background: var(--color-accent);
  border-radius: 50%;
  border: 3px solid var(--color-background);
  box-shadow: 0 0 0 2px var(--color-accent);
}

.year-badge {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-heading);
}

.year-count {
  font-size: 13px;
  color: var(--color-text-secondary);
  background: var(--color-background-soft);
  padding: 2px 10px;
  border-radius: 12px;
  border: 1px solid var(--color-border);
}

.timeline-month {
  margin-left: 10px;
  margin-bottom: 24px;
}

.month-title {
  font-size: 16px;
  color: var(--color-text);
  margin-bottom: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.month-count {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-weight: 400;
}

.timeline-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.timeline-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 6px;
  transition: all 0.25s ease;
}

.timeline-item:hover {
  background: var(--color-background-soft);
  transform: translateX(4px);
}

.date-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--color-border);
  flex-shrink: 0;
  transition: background 0.25s;
}

.timeline-item:hover .date-dot {
  background: var(--color-accent);
}

.date {
  color: var(--color-text-secondary);
  font-size: 13px;
  min-width: 36px;
}

.title {
  color: var(--color-heading);
  text-decoration: none;
  font-size: 15px;
  transition: color 0.25s;
  line-height: 1.5;
}

.title:hover {
  color: var(--color-accent);
}

/* Sticky Year Nav */
.year-nav {
    width: 56px;
    position: sticky;
    top: 100px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding-left: 16px;
    border-left: 2px solid var(--color-border);
}

.year-nav-item {
    cursor: pointer;
    font-size: 13px;
    color: var(--color-text-secondary);
    transition: all 0.25s;
    font-weight: 500;
    padding: 4px 0;
    position: relative;
}

.year-nav-item::before {
  content: '';
  position: absolute;
  left: -19px;
  top: 50%;
  transform: translateY(-50%);
  width: 0;
  height: 0;
  border-radius: 50%;
  background: var(--color-accent);
  transition: all 0.25s;
}

.year-nav-item:hover {
    color: var(--color-accent);
    transform: translateX(3px);
}

.year-nav-item.active {
    color: var(--color-accent);
    font-weight: 700;
}

.year-nav-item.active::before {
  width: 8px;
  height: 8px;
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
    color: var(--color-text-secondary);
}

/* Loading */
.loading-state {
  text-align: center;
  padding: 60px 40px;
}

.spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* List Transitions */
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.4s cubic-bezier(0.55, 0, 0.1, 1);
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.list-leave-active {
  position: absolute; 
  width: 100%;
}

/* Responsive */
@media (max-width: 768px) {
  .archive-page {
    padding: 20px 15px;
  }

  .page-header h1 {
    font-size: 22px;
  }
  
  .tag-list {
    justify-content: flex-start;
    overflow-x: auto;
    padding-bottom: 8px;
    flex-wrap: nowrap;
    white-space: nowrap;
    -webkit-overflow-scrolling: touch;
  }
  
  .tag-item {
    flex-shrink: 0;
  }

  .tag-toggle {
    display: none; 
  }

  .contribution-section {
    flex-direction: column;
    padding: 16px;
    gap: 16px;
  }
  
  .contrib-stats {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--color-border);
    padding-right: 0;
    padding-bottom: 14px;
    flex-direction: row;
    justify-content: space-around;
    min-width: unset;
  }

  .stat-number {
    font-size: 18px;
  }

  .timeline {
    padding-left: 20px;
  }

  .timeline::before {
    left: 4px;
  }

  .year-title::before {
    left: -22px;
    width: 10px;
    height: 10px;
  }

  .year-nav {
    display: none;
  }

  .archive-content {
    gap: 0;
  }
}
</style>
