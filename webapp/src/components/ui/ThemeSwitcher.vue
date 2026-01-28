<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useDarkMode, type ThemeMode } from '@/composables/useDarkMode'

const { themeMode, setTheme } = useDarkMode()
const isExpanded = ref(false)
const isMobile = ref(false)

// 检测是否为移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768 || ('ontouchstart' in window)
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

// 获取当前主题的图标
const currentIcon = computed(() => {
  switch (themeMode.value) {
    case 'light':
      return {
        svg: '<circle cx="12" cy="12" r="4"/><path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/>',
        color: 'text-amber-500',
        index: 0
      }
    case 'system':
      return {
        svg: '<rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/>',
        color: 'text-blue-500 dark:text-blue-400',
        index: 1
      }
    case 'dark':
      return {
        svg: '<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>',
        color: 'text-indigo-500 dark:text-indigo-400',
        index: 2
      }
  }
})

// 计算指示器位置（展开时）
const indicatorPosition = computed(() => {
  if (themeMode.value === 'light') return '0.0625rem'
  if (themeMode.value === 'system') return '1.75rem'
  return '3.4375rem'
})

// 计算容器偏移（收起时让当前按钮居中显示）
const containerOffset = computed(() => {
  if (isExpanded.value) return '0px'
  // 收起时根据当前主题位置调整偏移
  const index = currentIcon.value.index
  return `${-index * 28}px` // 28px = w-7 = 1.75rem
})

// 循环切换主题（移动端）
const cycleTheme = () => {
  const modes: ThemeMode[] = ['light', 'system', 'dark']
  const currentIndex = modes.indexOf(themeMode.value)
  const nextIndex = (currentIndex + 1) % modes.length
  setTheme(modes[nextIndex] as ThemeMode)
}

// 处理鼠标进入（仅桌面端）
const handleMouseEnter = () => {
  if (!isMobile.value) {
    isExpanded.value = true
  }
}

// 处理鼠标离开（仅桌面端）
const handleMouseLeave = () => {
  if (!isMobile.value) {
    isExpanded.value = false
  }
}

// 处理点击
const handleClick = () => {
  if (isMobile.value) {
    cycleTheme()
  }
}

// 处理按钮点击
const handleButtonClick = (mode: ThemeMode, event: Event) => {
  if (isMobile.value) {
    // 移动端让事件冒泡到外层处理
    return
  }
  // 桌面端阻止冒泡并切换到指定主题
  event.stopPropagation()
  setTheme(mode)
}
</script>

<template>
  <div 
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    @click="handleClick"
    class="relative overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700 transition-all duration-300 ease-in-out cursor-pointer"
    :style="{ width: isExpanded ? '5.5rem' : '2rem', padding: '0.125rem' }"
  >
    <!-- 容器，用于横向滑动显示当前按钮 -->
    <div 
      class="relative flex items-center transition-transform duration-300 ease-in-out"
      :style="{ transform: `translateX(${containerOffset})` }"
    >
      <!-- 滑动指示器（仅展开时显示） -->
      <div 
        v-show="isExpanded"
        class="absolute top-0 bottom-0 w-7 bg-white dark:bg-gray-600 rounded-full shadow-sm transition-all duration-300 ease-in-out"
        :style="{ left: indicatorPosition }"
      ></div>
      
      <!-- 浅色模式按钮 -->
      <button
        @click="(e) => handleButtonClick('light', e)"
        :class="[
          'relative z-10 flex items-center justify-center w-7 h-7 rounded-full transition-all duration-200 flex-shrink-0',
          themeMode === 'light' ? 'text-amber-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
        ]"
        :title="isExpanded ? '浅色模式' : ''"
        type="button"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="4"/>
          <path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.32 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/>
        </svg>
      </button>
      
      <!-- 自动模式按钮 -->
      <button
        @click="(e) => handleButtonClick('system', e)"
        :class="[
          'relative z-10 flex items-center justify-center w-7 h-7 rounded-full transition-all duration-200 flex-shrink-0',
          themeMode === 'system' ? 'text-blue-500 dark:text-blue-400' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
        ]"
        :title="isExpanded ? '跟随系统' : ''"
        type="button"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <rect x="2" y="3" width="20" height="14" rx="2"/>
          <path d="M8 21h8M12 17v4"/>
        </svg>
      </button>
      
      <!-- 深色模式按钮 -->
      <button
        @click="(e) => handleButtonClick('dark', e)"
        :class="[
          'relative z-10 flex items-center justify-center w-7 h-7 rounded-full transition-all duration-200 flex-shrink-0',
          themeMode === 'dark' ? 'text-indigo-500 dark:text-indigo-400' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300'
        ]"
        :title="isExpanded ? '深色模式' : ''"
        type="button"
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
        </svg>
      </button>
    </div>
  </div>
</template>
