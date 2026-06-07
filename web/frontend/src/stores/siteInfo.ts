import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import yaml from 'js-yaml'

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
  dev_admin_port?: number
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
    tagline: string
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
    wechat_qr: string
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
    dev_admin_port: 3011,
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
      tagline: '',
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
      wechat_qr: '',
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
      let configContent: string | null = null

      try {
        const apiResponse = await axios.get('/api/v1/frontend/config')
        if (apiResponse.data && apiResponse.data.status === 200 && apiResponse.data.data) {
          configContent = apiResponse.data.data
          // console.log('Loaded config from API')
        }
      } catch (apiError) {
        // console.warn('Failed to fetch config from API, falling back to static file:', apiError)
      }

      if (!configContent) {
        const fileResponse = await axios.get(`/config.yaml?t=${new Date().getTime()}`)
        if (fileResponse.data) {
          configContent = fileResponse.data
          // console.log('Loaded config from static file')
       }
      }

      if (configContent) {
        const parsedConfig = yaml.load(configContent) as SiteInfo
        applyLocalEnvironmentOverride(parsedConfig)
        siteInfo.value = parsedConfig
      }
    } catch (error) {
      console.error('Failed to fetch site info:', error)
    }
  }

  const applyLocalEnvironmentOverride = (config: SiteInfo) => {
    if (typeof window === 'undefined') return

    const hostname = window.location.hostname
    const isLocal = hostname === 'localhost' || 
                    hostname === '127.0.0.1' || 
                    hostname.startsWith('192.168.') || 
                    hostname.startsWith('10.') || 
                    hostname.endsWith('.local')
    
    if (isLocal) {
      const port = config.dev_admin_port || 3011
      config.admin_url = `http://${hostname}:${port}`
      // console.log('Detected local environment, overriding admin_url to:', config.admin_url)
    }
  }

  const updateConfig = async (newConfig: SiteInfo): Promise<boolean> => {
    try {
      const yamlContent = yaml.dump(newConfig)
      const response = await axios.put('/api/v1/frontend/config', { content: yamlContent })
      if (response.data && response.data.status === 200) {
        siteInfo.value = newConfig
        return true
      }
      return false
    } catch (error) {
      console.error('Failed to update config:', error)
      return false
    }
  }

  const refreshConfig = async () => {
    await fetchSiteInfo()
  }

  return {
    siteInfo,
    fetchSiteInfo,
    updateConfig,
    refreshConfig
  }
})