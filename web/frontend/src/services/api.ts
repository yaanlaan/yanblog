import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 文章相关API
export const articleApi = {
  // 获取文章列表
  getArticles: (params: { pagesize: number; pagenum: number }) => 
    apiClient.get('/article', { params }),
  
  // 搜索文章
  searchArticles: (params: { pagesize: number; pagenum: number; keyword?: string; cid?: number }) => 
    apiClient.get('/article/search', { params }),
  
  // 获取分类下的文章
  getCategoryArticles: (id: number, params: { pagesize: number; pagenum: number }) => 
    apiClient.get(`/article/list/${id}`, { params }),
  
  // 获取文章详情
  getArticle: (id: number) => 
    apiClient.get(`/article/info/${id}`),
    
  // 获取置顶文章
  getTopArticles: (params?: { num: number }) => 
    apiClient.get('/article/top', { params })
}

// 分类相关API
export const categoryApi = {
  // 获取分类列表
  getCategories: (params: { pagesize: number; pagenum: number }) => 
    apiClient.get('/category', { params })
}

export default apiClient