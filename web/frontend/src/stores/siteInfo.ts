import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import yaml from 'js-yaml'

// 定义站点信息接口
export interface SiteInfo {
  blog_name: string
  author_name: string
  author_avatar: string
  author_bio: string
  hero: {
    title: string
    subtitle: string
    welcome: string
    welcome_image: string
  }
  quotes: string[]
  logo_text: string
  logo_image: string
  favicon: string
  admin_url: string
  page_title: {
    default: string
    blur: string
  }
  iconfont_url: string
  music_player: {
    show: boolean
    url: string
  }
  shortcuts: Array<{
    name: string
    url: string
    icon: string
    color: string
  }>
  footer: {
    copyright: string
    email: string
    icp: {
      show: boolean
      text: string
      link: string
    }
    powered_by: string
    powered_by_link: string
    portfolio: {
      show: boolean
      title: string
      icon: string
      items: Array<{
        name: string
        url: string
        icon: string
      }>
    }
    related_links: {
      show: boolean
      title: string
      items: Array<{
        name: string
        url: string
        icon: string
      }>
    }
  }
  socials: Array<{
    name: string
    url: string
    icon: string
    color: string
    is_circle: boolean
  }>
  contacts: {
    show: boolean
    wechat_qr: string // 新增微信二维码配置
    items: Array<{
      name: string
      url: string
      icon: string
      color: string
      is_circle: boolean
    }>
  }
  comment: {
    enable: boolean
    type: 'giscus' | 'other'
    giscus: {
      repo: string
      repo_id: string
      category: string
      category_id: string
      mapping: string
      reactions_enabled: string
      emit_metadata: string
      input_position: string
      theme: string
      lang: string
      loading: string
    }
  }
}

export const useSiteInfoStore = defineStore('siteInfo', () => {
  const siteInfo = ref<SiteInfo>({
    blog_name: '言盐盐的博客',
    author_name: 'Yaan',
    author_avatar: '',
    author_bio: '',
    hero: {
      title: '',
      subtitle: '',
      welcome: '',
      welcome_image: ''
    },
    quotes: [],
    logo_text: '言盐盐的博客',
    logo_image: '',
    favicon: '',
    admin_url: '',
    page_title: {
      default: '',
      blur: ''
    },
    iconfont_url: '',
    music_player: {
      show: false,
      url: ''
    },
    shortcuts: [],
    footer: {
      copyright: '',
      email: '',
      icp: {
        show: false,
        text: '',
        link: ''
      },
      powered_by: '',
      powered_by_link: '',
      portfolio: {
        show: false,
        title: '',
        icon: '',
        items: []
      },
      related_links: {
        show: false,
        title: '',
        items: []
      }
    },
    socials: [],
    contacts: {
      show: false,
      wechat_qr: '', // 初始化
      items: []
    },
    comment: {
      enable: false,
      type: 'giscus',
      giscus: {
        repo: '',
        repo_id: '',
        category: '',
        category_id: '',
        mapping: 'pathname',
        reactions_enabled: '1',
        emit_metadata: '0',
        input_position: 'top',
        theme: 'light',
        lang: 'zh-CN',
        loading: 'lazy'
      }
    }
  })

  const fetchSiteInfo = async () => {
    try {
      // 直接请求 public 目录下的 config.yaml，添加时间戳防止缓存
      const response = await axios.get(`/config.yaml?t=${new Date().getTime()}`)
      if (response.data) {
        // 解析 YAML 内容
        const parsedConfig = yaml.load(response.data) as SiteInfo
        siteInfo.value = parsedConfig
      }
    } catch (error) {
      console.error('Failed to fetch site info:', error)
    }
  }

  return {
    siteInfo,
    fetchSiteInfo
  }
})
