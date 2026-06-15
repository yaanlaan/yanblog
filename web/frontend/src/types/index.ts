/**
 * 前端统一的 TypeScript 类型定义
 * 消除各组件中重复的类型定义，保持一致性
 */

// 文章类型
export interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  top: number
  tags: string
  views: number
  type?: number
  pdf_url?: string
  createdAt: string
  updatedAt: string
}

// 分类类型
export interface Category {
  id: number
  name: string
}

// 标签类型
export interface Tag {
  id: number
  name: string
}

// 用户类型（后台用）
export interface User {
  id: number
  username: string
  role: number
  createdAt: string
  updatedAt: string
}

// 分页响应类型
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  status: number
  message: string
}

// API 响应类型
export interface ApiResponse<T = any> {
  status: number
  data: T
  message: string
}

// 分页参数
export interface PaginationParams {
  pagesize: number
  pagenum: number
}

// 搜索参数
export interface SearchParams extends PaginationParams {
  keyword?: string
  cid?: number
}
