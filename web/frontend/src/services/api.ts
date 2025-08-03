import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 15000, // 设置15秒超时
  headers: {
    'Content-Type': 'application/json'
  }
})

// 添加请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 在发送请求之前做些什么
    console.log('API Request:', config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    // 对请求错误做些什么
    console.error('API Request Error:', error);
    return Promise.reject(error);
  }
);

// 添加响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    // 对响应数据做点什么
    console.log('API Response:', response.status, response.config.url);
    return response;
  },
  (error) => {
    // 对响应错误做点什么
    if (error.code === 'ECONNABORTED') {
      console.error('API Timeout:', error.config.url);
      return Promise.reject(new Error('请求超时，请稍后重试'));
    }
    
    if (!error.response) {
      console.error('API Network Error:', error.config?.url);
      return Promise.reject(new Error('网络连接错误，请检查后端服务是否启动'));
    }
    
    console.error('API Error:', error.response.status, error.response.config.url);
    return Promise.reject(error);
  }
);

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
    apiClient.get('/category', { params }),
    
  // 获取分类信息
  getCategoryInfo: (id: number) => 
    apiClient.get(`/category/info/${id}`)
}

// 天气相关API
export const weatherApi = {
  // 获取天气信息
  getWeather: (params?: { city?: string }) => 
    apiClient.get('/weather', { params })
}

export default apiClient