import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import zhTW from './locales/zh-TW'
import enUS from './locales/en-US'
import jaJP from './locales/ja-JP'
import koKR from './locales/ko-KR'
import ruRU from './locales/ru-RU'

// 获取存储的语言或默认为中文
const getDefaultLocale = () => {
  const stored = localStorage.getItem('locale')
  if (stored && ['zh-CN', 'zh-TW', 'en-US', 'ja-JP', 'ko-KR', 'ru-RU'].includes(stored)) {
    return stored
  }
  return 'zh-CN'
}

const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: getDefaultLocale(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'zh-TW': zhTW,
    'en-US': enUS,
    'ja-JP': jaJP,
    'ko-KR': koKR,
    'ru-RU': ruRU,
  },
})

export default i18n
