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
    items: Array<{
      name: string
      url: string
      icon: string
      color: string
      is_circle: boolean
    }>
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
      items: []
    }
  })

  const fetchSiteInfo = async () => {
    try {
      // 直接请求 public 目录下的 config.yaml
      const response = await axios.get('/config.yaml')
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
