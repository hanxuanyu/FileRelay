import { ref, watch, onMounted } from 'vue'

// 暗黑模式类型：'light' | 'dark' | 'system'
export type ThemeMode = 'light' | 'dark' | 'system'

// 当前主题模式（用户选择的）
const themeMode = ref<ThemeMode>('system')
// 实际应用的主题（解析后的）
const isDark = ref(false)

// 从 localStorage 加载主题设置
const loadTheme = (): ThemeMode => {
  const saved = localStorage.getItem('theme-mode')
  if (saved === 'light' || saved === 'dark' || saved === 'system') {
    return saved
  }
  return 'system'
}

// 保存主题设置到 localStorage
const saveTheme = (mode: ThemeMode) => {
  localStorage.setItem('theme-mode', mode)
}

// 获取系统主题偏好
const getSystemTheme = (): boolean => {
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

// 应用暗黑模式到 DOM
const applyTheme = (dark: boolean) => {
  isDark.value = dark
  if (dark) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// 根据当前模式计算实际主题
const resolveTheme = (mode: ThemeMode): boolean => {
  if (mode === 'system') {
    return getSystemTheme()
  }
  return mode === 'dark'
}

// 更新主题
const updateTheme = () => {
  const dark = resolveTheme(themeMode.value)
  applyTheme(dark)
}

// 监听系统主题变化
let mediaQuery: MediaQueryList | null = null

export function useDarkMode() {
  // 初始化
  onMounted(() => {
    // 加载保存的主题设置
    themeMode.value = loadTheme()
    
    // 应用主题
    updateTheme()
    
    // 监听系统主题变化
    mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    const handleChange = () => {
      if (themeMode.value === 'system') {
        updateTheme()
      }
    }
    mediaQuery.addEventListener('change', handleChange)
  })

  // 监听主题模式变化
  watch(themeMode, (newMode) => {
    saveTheme(newMode)
    updateTheme()
  })

  // 切换到下一个模式（循环切换）
  const cycleTheme = () => {
    const modes: ThemeMode[] = ['light', 'dark', 'system']
    const currentIndex = modes.indexOf(themeMode.value)
    const nextIndex = (currentIndex + 1) % modes.length
    themeMode.value = modes[nextIndex] as ThemeMode
  }

  // 设置特定主题模式
  const setTheme = (mode: ThemeMode) => {
    themeMode.value = mode
  }

  return {
    themeMode,
    isDark,
    cycleTheme,
    setTheme
  }
}

// 初始化主题（在 app 启动时调用）
export function initDarkMode() {
  themeMode.value = loadTheme()
  updateTheme()
  
  // 监听系统主题变化
  mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const handleChange = () => {
    if (themeMode.value === 'system') {
      updateTheme()
    }
  }
  mediaQuery.addEventListener('change', handleChange)
}
