<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <StatCard type="users" :number="stats.users" label="用户数" />
      </el-col>
      
      <el-col :span="6">
        <StatCard type="articles" :number="stats.articles" label="文章数" />
      </el-col>
      
      <el-col :span="6">
        <StatCard type="categories" :number="stats.categories" label="分类数" />
      </el-col>
      
      <el-col :span="6">
        <StatCard type="views" :number="stats.views" label="浏览量" />
      </el-col>
    </el-row>
    
    <el-row :gutter="20" class="chart-row">
      <el-col :span="24">
        <DataChart ref="chartRef" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { userApi, articleApi, categoryApi } from '@/services/api'
import { ElMessage } from 'element-plus'
import StatCard from './StatCard.vue'
import DataChart from './DataChart.vue'

// 统计数据
const stats = reactive({
  users: 0,
  articles: 0,
  categories: 0,
  views: 0
})

// 图表引用
const chartRef = ref<InstanceType<typeof DataChart> | null>(null)

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
    
    // 更新图表
    if (chartRef.value) {
      chartRef.value.updateChart()
    }
  } catch (error) {
    ElMessage.error('获取统计数据失败')
    console.error(error)
  }
}

// 窗口大小变化时重置图表
const handleResize = () => {
  if (chartRef.value) {
    chartRef.value.resizeChart()
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getStats()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.chart-row {
  margin-top: 20px;
}
</style>