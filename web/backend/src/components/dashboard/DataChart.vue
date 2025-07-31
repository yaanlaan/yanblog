<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>文章发布统计</span>
        <div class="chart-controls">
          <el-radio-group v-model="timeRange" size="small" @change="updateChart">
            <el-radio-button label="7days">近7天</el-radio-button>
            <el-radio-button label="6months">近6个月</el-radio-button>
          </el-radio-group>
          <el-radio-group v-model="chartType" size="small" @change="updateChart" style="margin-left: 10px;">
            <el-radio-button label="bar">柱状图</el-radio-button>
            <el-radio-button label="line">折线图</el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </template>
    <div ref="chartRef" class="chart-container"></div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, defineExpose } from 'vue'
import * as echarts from 'echarts'
import { articleApi } from '@/services/api'
import { ElMessage } from 'element-plus'

// 图表引用
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

// 时间范围
const timeRange = ref<'7days' | '6months'>('7days')

// 图表类型
const chartType = ref<'line' | 'bar'>('bar')

// 生成近7天的日期数组
const generateLast7Days = () => {
  const dates = []
  for (let i = 6; i >= 0; i--) {
    const date = new Date()
    date.setDate(date.getDate() - i)
    dates.push(date.toISOString().split('T')[0])
  }
  return dates
}

// 生成近6个月的月份数组
const generateLast6Months = () => {
  const months = []
  const now = new Date()
  for (let i = 5; i >= 0; i--) {
    const date = new Date(now.getFullYear(), now.getMonth() - i, 1)
    months.push(`${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`)
  }
  return months
}

// 按日期统计文章数量
const countArticlesByDate = (articles: any[], dates: string[]) => {
  const counts: Record<string, number> = {}
  
  // 初始化计数器
  dates.forEach(date => {
    counts[date] = 0
  })
  
  // 统计文章
  articles.forEach(article => {
    const articleDate = article.CreatedAt.split('T')[0]
    if (counts.hasOwnProperty(articleDate)) {
      counts[articleDate]++
    }
  })
  
  return dates.map(date => counts[date])
}

// 按月份统计文章数量
const countArticlesByMonth = (articles: any[], months: string[]) => {
  const counts: Record<string, number> = {}
  
  // 初始化计数器
  months.forEach(month => {
    counts[month] = 0
  })
  
  // 统计文章
  articles.forEach(article => {
    const articleMonth = article.CreatedAt.substring(0, 7)
    if (counts.hasOwnProperty(articleMonth)) {
      counts[articleMonth]++
    }
  })
  
  return months.map(month => counts[month])
}

// 获取文章数据并更新图表
const updateChart = async () => {
  if (!chartRef.value) return
  
  try {
    // 获取所有文章数据
    const response = await articleApi.getArticles({ pagesize: -1, pagenum: -1 })
    const articles = response.data.data
    
    let xAxisData: string[] = []
    let seriesData: number[] = []
    
    if (timeRange.value === '7days') {
      // 近7天数据
      xAxisData = generateLast7Days()
      seriesData = countArticlesByDate(articles, xAxisData)
    } else {
      // 近6个月数据
      xAxisData = generateLast6Months()
      seriesData = countArticlesByMonth(articles, xAxisData)
    }
    
    // 计算Y轴最大值，确保至少显示到5，并且是整数
    const maxValue = Math.max(...seriesData);
    const yAxisMax = Math.max(5, Math.ceil(maxValue));
    
    // 更新图表
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['文章数']
      },
      xAxis: {
        type: 'category',
        data: xAxisData
      },
      yAxis: {
        type: 'value',
        min: 0,
        max: yAxisMax,
        interval: 1, // 设置刻度间隔为1，确保显示整数刻度
        axisLabel: {
          formatter: '{value}' // 确保显示为整数
        }
      },
      series: [
        {
          name: '文章数',
          type: chartType.value,
          data: seriesData,
          smooth: chartType.value === 'line' // 只有折线图才启用平滑
        }
      ]
    }
    
    if (!chart) {
      chart = echarts.init(chartRef.value)
    }
    
    chart.setOption(option, true)
  } catch (error) {
    ElMessage.error('获取文章统计数据失败')
    console.error(error)
  }
}

// 窗口大小变化时重置图表
const resizeChart = () => {
  if (chart) {
    chart.resize()
  }
}

// 组件挂载时初始化图表
onMounted(() => {
  updateChart()
})

// 组件卸载时清理资源
onUnmounted(() => {
  if (chart) {
    chart.dispose()
  }
})

// 暴露方法给父组件使用
defineExpose({
  resizeChart,
  updateChart
})
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-controls {
  display: flex;
  align-items: center;
}
</style>