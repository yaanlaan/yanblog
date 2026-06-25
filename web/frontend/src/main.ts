import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

import './assets/main.css'
import './utils/markdown' // Marked v16 全局配置（Mac 风格代码块）

import App from './App.vue'
import router from './router'
import lazy from './directives/v-lazy'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.directive('lazy', lazy)

// 全局错误处理：防止子组件异常导致白屏
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误捕获:', err)
  console.error('错误组件:', instance?.$options?.name || '未知')
  console.error('错误信息:', info)
}

app.mount('#app')