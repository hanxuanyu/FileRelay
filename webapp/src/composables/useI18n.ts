import { useI18n as vueUseI18n } from 'vue-i18n'

export function useI18n() {
  const i18n = vueUseI18n()
  
  const setLocale = (locale: string) => {
    i18n.locale.value = locale
    localStorage.setItem('locale', locale)
  }

  return {
    ...i18n,
    setLocale,
  }
}
