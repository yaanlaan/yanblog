/**
 * 前端公共常量定义
 * 集中管理所有常量，消除硬编码的魔法数字
 */

// 响应式断点（像素）
export const BREAKPOINTS = {
  MOBILE: 768,
  TABLET: 992,
  DESKTOP: 1450,
  WIDE: 1600
} as const

// 分页配置
export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 10,
  ARTICLE_LIST_PAGE_SIZE: 12,
  HOME_PAGE_SIZE: 10
} as const

// 文章相关
export const ARTICLE = {
  HOT_COUNT: 5,           // 热门文章数量
  TOP_COUNT: 6,           // 置顶文章数量
  RELATED_COUNT: 5,       // 相关文章数量
  DEFAULT_COVER: '/assets/img/default-cover.jpg'
} as const

// 动画时长（毫秒）
export const ANIMATION = {
  FAST: 200,
  NORMAL: 300,
  SLOW: 500,
  LOADING_DELAY: 500
} as const

// 路由名称
export const ROUTES = {
  HOME: 'home',
  ARTICLES: 'articles',
  ARTICLE_DETAIL: 'article-detail',
  CATEGORY_ARTICLES: 'category-articles',
  CATEGORIES: 'categories',
  ARCHIVE: 'archive',
  ABOUT: 'about',
  NOT_FOUND: 'not-found'
} as const
