/**
 * 数据映射工具函数
 * 将后端 API 返回的数据映射为前端统一的类型格式
 * 消除各组件中重复的映射代码
 */

import type { Article, Category, RawArticle, RawCategory } from '@/types'

export const mapArticle = (item: RawArticle): Article => ({
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
  createdAt: item.CreatedAt || item.created_at || '',
  updatedAt: item.UpdatedAt || item.updated_at || ''
})

export const mapCategory = (item: RawCategory): Category => ({
  id: item.ID,
  name: item.name
})

export const mapArticleList = (items: RawArticle[]): Article[] =>
  items.map(mapArticle)

export const mapCategoryList = (items: RawCategory[]): Category[] =>
  items.map(mapCategory)
