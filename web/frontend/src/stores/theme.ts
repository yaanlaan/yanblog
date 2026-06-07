import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
    // 默认优先使用本地存储，否则跟随系统，默认浅色
    const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    
    const theme = ref<'light' | 'dark'>(savedTheme || (prefersDark ? 'dark' : 'light'))

    const applyTheme = (currentTheme: 'light' | 'dark') => {
        if (currentTheme === 'dark') {
            document.documentElement.setAttribute('data-theme', 'dark')
            document.documentElement.classList.add('dark')
        } else {
            document.documentElement.removeAttribute('data-theme')
            document.documentElement.classList.remove('dark')
        }
    }

    // 初始化应用主题
    applyTheme(theme.value)

    const toggleTheme = () => {
        theme.value = theme.value === 'light' ? 'dark' : 'light'
    }

    // 监听变化并保存
    watch(theme, (newTheme) => {
        localStorage.setItem('theme', newTheme)
        applyTheme(newTheme)
    })

    return {
        theme,
        toggleTheme
    }
})
