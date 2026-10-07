import axios from 'axios'
import type { AxiosPromise, AxiosRequestConfig } from 'axios'

export const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

type CancelablePromise<T> = AxiosPromise<T> & { abort: () => void }

const requestWithCancel = <T>(config: AxiosRequestConfig): CancelablePromise<T> => {
  const controller = new AbortController()
  const promise = apiClient({ ...config, signal: controller.signal }) as CancelablePromise<T>
  promise.abort = () => controller.abort()
  return promise
}

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new Error('请求超时，请稍后重试'))
    }
    if (!error.response) {
      return Promise.reject(new Error('网络连接错误，请检查后端服务是否启动'))
    }
    return Promise.reject(error)
  }
)

// 文章相关API
export const articleApi = {
  // 获取文章列表
  getArticles: (params: { pagesize: number; pagenum: number; excludeTop?: boolean }) =>
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
    apiClient.get('/article/top', { params }),

  // 获取热门文章
  getHotArticles: (params?: { num: number }) =>
    apiClient.get('/article/hot', { params }),

  // 获取相关文章
  getRelatedArticles: (id: number) =>
    apiClient.get(`/article/related/${id}`),

  // 随机获取一篇文章
  getRandomArticle: () =>
    apiClient.get('/article/random'),

  // 获取相邻文章（上一篇/下一篇）
  getAdjacentArticle: (id: number) =>
    apiClient.get(`/article/adjacent/${id}`)
}

// 分类相关API
export const categoryApi = {
  // 获取分类列表
  getCategories: (params: { pagesize: number; pagenum: number }) =>
    apiClient.get('/category', { params }),

  // 获取分类信息
  getCategoryInfo: (id: number) =>
    apiClient.get(`/category/info/${id}`)
}

// 标签相关API
export const tagApi = {
  // 获取标签列表
  getTags: (params: { pagesize: number; pagenum: number }) =>
    apiClient.get('/tags', { params })
}

// 天气相关API
export const weatherApi = {
  // 获取天气信息
  getWeather: (params?: { city?: string }) =>
    apiClient.get('/weather', { params })
}

// 系统状态相关API
export const systemApi = {
  // 获取系统状态信息
  getSystemStatus: () =>
    apiClient.get('/system/status')
}

export default apiClient