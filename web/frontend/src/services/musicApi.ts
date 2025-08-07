// 网易云音乐API服务
import axios from 'axios'

// 网易云音乐API基础URL (使用第三方代理)
const MUSIC_API_BASE = 'https://api.bzqll.com/music/tencent'

// 创建axios实例
const musicApi = axios.create({
  baseURL: MUSIC_API_BASE,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  }
})

// 获取推荐歌曲
export const getRecommendSongs = async () => {
  try {
    const response = await musicApi.get('/recommendSongs')
    return response.data
  } catch (error) {
    console.error('获取推荐歌曲失败:', error)
    throw error
  }
}

// 获取歌单详情
export const getPlaylistDetail = async (id: string) => {
  try {
    const response = await musicApi.get(`/playlist/detail?id=${id}`)
    return response.data
  } catch (error) {
    console.error('获取歌单详情失败:', error)
    throw error
  }
}

// 获取歌曲详情
export const getSongDetail = async (ids: string) => {
  try {
    const response = await musicApi.get(`/song/detail?ids=${ids}`)
    return response.data
  } catch (error) {
    console.error('获取歌曲详情失败:', error)
    throw error
  }
}

// 获取歌曲播放链接
export const getSongUrl = async (id: string) => {
  try {
    const response = await musicApi.get(`/song/url?id=${id}`)
    return response.data
  } catch (error) {
    console.error('获取歌曲播放链接失败:', error)
    throw error
  }
}

// 搜索歌曲
export const searchSongs = async (keywords: string, limit: number = 30) => {
  try {
    const response = await musicApi.get(`/search?keywords=${keywords}&limit=${limit}`)
    return response.data
  } catch (error) {
    console.error('搜索歌曲失败:', error)
    throw error
  }
}

export default {
  getRecommendSongs,
  getPlaylistDetail,
  getSongDetail,
  getSongUrl,
  searchSongs
}