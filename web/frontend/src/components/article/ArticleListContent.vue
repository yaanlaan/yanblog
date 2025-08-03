<template>
  <div class="articles" v-loading="loading">
    <div class="articles-grid">
      <ArticleItem
        v-for="article in articles" 
        :key="article.id" 
        :article="article"
      />
    </div>
    
    <!-- 分页 -->
    <div class="pagination" v-if="total > 0">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[5, 10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
    
    <!-- 空状态 -->
    <div class="empty-state" v-if="!loading && articles.length === 0">
      <p>暂无文章</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElPagination } from 'element-plus'
import ArticleItem from './ArticleItem.vue'

// 定义Props
interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  createdAt: string
  updatedAt: string
}

interface Props {
  articles: Article[]
  loading: boolean
  total: number
  currentPage: number
  pageSize: number
}

// 定义Emits
const emit = defineEmits<{
  (e: 'size-change', value: number): void
  (e: 'current-change', value: number): void
}>()

const props = defineProps<Props>()

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  emit('size-change', val)
}

// 处理当前页变化
const handleCurrentChange = (val: number) => {
  emit('current-change', val)
}
</script>

<style scoped>
.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 30px;
  margin-bottom: 30px;
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
  color: #888;
}

@media (max-width: 768px) {
  .articles-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}
</style>