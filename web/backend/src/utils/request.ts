import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: '/api', // 基础URL，与Vite代理配置匹配
  timeout: 10000, // 请求超时时间
  headers: {
    'Content-Type': 'application/json;charset=UTF-8'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 在发送请求之前做些什么
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    // 对请求错误做些什么
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    // 对响应数据做些什么
    const res = response.data
    
    // 根据后端约定的状态码进行处理
    if (res.status === 200) {
      return res
    } else {
      // 处理其他状态码
      console.error('响应错误:', res.message)
      return Promise.reject(new Error(res.message || 'Error'))
    }
  },
  (error) => {
    // 对响应错误做些什么
    console.error('响应错误:', error)
    return Promise.reject(error)
  }
)

export default service