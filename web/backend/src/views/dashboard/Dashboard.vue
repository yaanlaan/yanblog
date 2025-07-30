<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon users">
              <i class="el-icon-user"></i>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.users }}</div>
              <div class="stat-label">用户数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon articles">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.articles }}</div>
              <div class="stat-label">文章数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon categories">
              <i class="el-icon-folder"></i>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.categories }}</div>
              <div class="stat-label">分类数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon views">
              <i class="el-icon-view"></i>
            </div>
            <div class="stat-info">
              <div class="stat-number">{{ stats.views }}</div>
              <div class="stat-label">浏览量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>数据统计</span>
            </div>
          </template>
          <div ref="chartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { userApi, articleApi, categoryApi } from '@/services/api'
import { ElMessage } from 'element-plus'

// 统计数据
const stats = reactive({
  users: 0,
  articles: 0,
  categories: 0,
  views: 0
})

// 图表引用
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

// 获取统计数据
const getStats = async () => {
  try {
    // 获取用户总数
    const userResponse = await userApi.getUsers({ pagesize: -1, pagenum: -1 })
    stats.users = userResponse.data.total || 0
    
    // 获取文章总数
    const articleResponse = await articleApi.getArticles({ pagesize: -1, pagenum: -1 })
    stats.articles = articleResponse.data.total || 0
    
    // 获取分类总数
    const categoryResponse = await categoryApi.getCategories({ pagesize: -1, pagenum: -1 })
    stats.categories = categoryResponse.data.total || 0
    
    // 浏览量暂时使用模拟数据
    stats.views = 0
  } catch (error) {
    ElMessage.error('获取统计数据失败')
    console.error(error)
  }
}

// 初始化图表
const initChart = () => {
  if (!chartRef.value) return
  
  chart = echarts.init(chartRef.value)
  
  // 模拟图表数据
  const option = {
    title: {
      text: '文章发布统计'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['文章数']
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '文章数',
        type: 'line',
        data: [12, 23, 15, 20, 18, 25, 30],
        smooth: true
      }
    ]
  }
  
  chart.setOption(option)
}

// 窗口大小变化时重置图表
const handleResize = () => {
  if (chart) {
    chart.resize()
  }
}

// 组件挂载时获取数据并初始化图表
onMounted(() => {
  getStats()
  initChart()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (chart) {
    chart.dispose()
  }
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  font-size: 24px;
  color: white;
}

.stat-icon.users {
  background-color: #409eff;
}

.stat-icon.articles {
  background-color: #67c23a;
}

.stat-icon.categories {
  background-color: #e6a23c;
}

.stat-icon.views {
  background-color: #f56c6c;
}

.stat-info {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.chart-row {
  margin-top: 20px;
}

.chart-container {
  width: 100%;
  height: 400px;
}
</style>