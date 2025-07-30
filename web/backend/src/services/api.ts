// API服务封装
import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

// 创建axios实例
const apiClient: AxiosInstance = axios.create({
  baseURL: '/api', // 基础URL，实际项目中应配置为后端API地址
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // 在发送请求之前做些什么
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    // 对请求错误做些什么
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // 对响应数据做点什么
    return response
  },
  (error) => {
    // 对响应错误做点什么
    if (error.response?.status === 401) {
      // token过期或无效，清除本地存储并跳转到登录页
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 用户相关API
export const userApi = {
  // 登录
  login: (data: { username: string; password: string }) => 
    apiClient.post('/v1/login', data),
  
  // 获取用户列表
  getUsers: (params: { pagesize: number; pagenum: number }) => 
    apiClient.get('/v1/users', { params }),
  
  // 创建用户
  createUser: (data: { username: string; password: string; role: number }) => 
    apiClient.post('/v1/user/add', data),
  
  // 更新用户
  updateUser: (id: number, data: { username: string; role: number }) => 
    apiClient.put(`/v1/user/${id}`, data),
  
  // 删除用户
  deleteUser: (id: number) => 
    apiClient.delete(`/v1/user/${id}`)
}

// 分类相关API
export const categoryApi = {
  // 获取分类列表
  getCategories: (params: { pagesize: number; pagenum: number }) => 
    apiClient.get('/v1/category', { params }),
  
  // 创建分类
  createCategory: (data: { name: string }) => 
    apiClient.post('/v1/category/add', data),
  
  // 更新分类
  updateCategory: (id: number, data: { name: string }) => 
    apiClient.put(`/v1/category/${id}`, data),
  
  // 删除分类
  deleteCategory: (id: number) => 
    apiClient.delete(`/v1/category/${id}`)
}

// 文章相关API
export const articleApi = {
  // 获取文章列表
  getArticles: (params: { pagesize: number; pagenum: number }) => 
    apiClient.get('/v1/article', { params }),
  
  // 获取分类下的文章
  getCategoryArticles: (id: number, params: { pagesize: number; pagenum: number }) => 
    apiClient.get(`/v1/article/list/${id}`, { params }),
  
  // 获取文章详情
  getArticle: (id: number) => 
    apiClient.get(`/v1/article/info/${id}`),
  
  // 创建文章
  createArticle: (data: { 
    title: string; 
    cid: number; 
    desc: string; 
    content: string; 
    img: string 
  }) => 
    apiClient.post('/v1/article/add', data),
  
  // 更新文章
  updateArticle: (id: number, data: { 
    title: string; 
    cid: number; 
    desc: string; 
    content: string; 
    img: string 
  }) => 
    apiClient.put(`/v1/article/${id}`, data),
  
  // 删除文章
  deleteArticle: (id: number) => 
    apiClient.delete(`/v1/article/${id}`)
}

// 文件上传API
export const uploadApi = {
  // 上传文件
  uploadFile: (formData: FormData) => 
    apiClient.post('/v1/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
}

export default apiClient