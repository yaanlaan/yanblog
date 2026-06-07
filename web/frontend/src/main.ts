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

app.mount('#app')