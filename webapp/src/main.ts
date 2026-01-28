import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import { initPublicConfig } from '@/composables/usePublicConfig'
import { initDarkMode } from '@/composables/useDarkMode'

async function initApp() {
  // 初始化暗黑模式（先执行，避免闪烁）
  initDarkMode()
  
  // 初始化公共配置
  await initPublicConfig()
  
  // 创建应用
  createApp(App).use(router).use(i18n).mount('#app')
}

initApp()
