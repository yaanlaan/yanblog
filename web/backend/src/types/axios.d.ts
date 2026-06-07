// axios类型定义扩展
import 'axios'

declare module 'axios' {
  interface AxiosRequestConfig {
    // 扩展AxiosRequestConfig类型
  }

  interface AxiosResponse<T = any> {
    // 扩展AxiosResponse类型
    data: T
    code: number
    message: string
  }
}