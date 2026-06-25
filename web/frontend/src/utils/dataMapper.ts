/**
 * 数据映射工具函数
 * 将后端 API 返回的数据映射为前端统一的类型格式
 * 消除各组件中重复的映射代码
 */

import type { Article, Category } from '@/types'

/**
 * 映射文章数据
 * 将后端返回的原始文章数据转换为前端 Article 类型
 */
export const mapArticle = (item: any): Article => ({
  id: item.ID,
  title: item.title,
  categoryId: item.cid,
  categoryName: item.Category?.name || '未分类',
  desc: item.desc,
  content: item.content,
  img: item.img,
  top: item.top || 0,
  tags: item.tags || '',
  views: item.views || 0,
  type: item.type,
  pdf_url: item.pdf_url,
  createdAt: item.CreatedAt || item.created_at,
  updatedAt: item.UpdatedAt || item.updated_at
})

/**
 * 映射分类数据
 * 将后端返回的原始分类数据转换为前端 Category 类型
 */
export const mapCategory = (item: any): Category => ({
  id: item.ID,
  name: item.name
})

/**
 * 批量映射文章列表
 */
export const mapArticleList = (items: any[]): Article[] =>
  items.map(mapArticle)

/**
 * 批量映射分类列表
 */
export const mapCategoryList = (items: any[]): Category[] =>
  items.map(mapCategory)
