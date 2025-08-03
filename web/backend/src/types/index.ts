// 公共类型定义

// 用户类型
export interface User {
  id: number
  username: string
  role: number // 1: 管理员, 2: 普通用户
  createdAt: string
}

// 分类类型
export interface Category {
  id: number
  name: string
  createdAt: string
}

// 文章类型
export interface Article {
  id: number
  title: string
  categoryId: number
  categoryName: string
  desc: string
  content: string
  img: string
  top: number // 置顶等级，0表示不置顶，1-6表示置顶等级，数字越小等级越高
  createdAt: string
}

// 分页参数类型
export interface PaginationParams {
  pagesize: number
  pagenum: number
}

// 分页响应类型
export interface PaginationResponse<T> {
  status: number
  data: T[]
  total: number
  message: string
}