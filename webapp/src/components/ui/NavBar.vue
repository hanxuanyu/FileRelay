<template>
  <nav class="bg-white/75 dark:bg-gray-900/75 backdrop-blur-lg border-b border-gray-200 dark:border-gray-700 transition-colors">
    <div class="container mx-auto px-3 sm:px-4">
      <div class="flex justify-between items-center h-16">
        <!-- 左侧：站点信息 -->
        <div class="flex items-center space-x-2 sm:space-x-4">
          <a href="/" class="flex items-center space-x-2 sm:space-x-3">
            <!-- Logo -->
            <div v-if="config.site?.logo" class="w-8 h-8 sm:w-9 sm:h-9 rounded-lg overflow-hidden flex items-center justify-center flex-shrink-0">
              <img :src="config.site.logo" :alt="config.site?.name || t('site.title', '文件中转站')" class="w-full h-full object-contain" />
            </div>
            <div v-else class="w-8 h-8 sm:w-9 sm:h-9 bg-blue-600 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
            </div>
            <!-- 移动端隐藏文字，只保留图标 -->
            <div class="hidden sm:block">
              <h1 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ config.site?.name || t('site.title', '文件中转站') }}</h1>
              <p v-if="showDescription" class="text-xs text-gray-500 dark:text-gray-400">{{ config.site?.description || t('site.description', '安全、便捷的文件暂存服务') }}</p>
            </div>
          </a>
        </div>

        <!-- 右侧：导航链接 -->
        <div class="flex items-center space-x-1 sm:space-x-2">
          <!-- 主题切换 - 放在最左侧 -->
          <ThemeSwitcher />
          
          <!-- 取件/发送切换按钮组 -->
          <div ref="buttonGroup" class="relative inline-flex items-center rounded-full bg-gray-100 dark:bg-gray-800 p-0.5 transition-colors">
            <!-- 滑动指示器 -->
            <div 
              :class="[
                'absolute inset-y-0.5 bg-white dark:bg-gray-700 rounded-full shadow-sm',
                isInitialized && 'transition-all duration-300 ease-in-out'
              ]"
              :style="{ 
                left: indicatorPosition,
                width: indicatorWidth
              }"
            ></div>
            
            <button @click="switchTab('pickup')" ref="pickupButton">
              <Button 
                variant="ghost" 
                size="sm" 
                class="relative z-10 h-7 px-2 sm:px-4 rounded-full hover:bg-transparent"
              >
                <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5l4-4 4 4"></path>
                </svg>
                <span class="hidden sm:inline">{{ t('nav.pickup') }}</span>
              </Button>
            </button>
            
            <button @click="switchTab('upload')" ref="uploadButton">
              <Button 
                variant="ghost" 
                size="sm" 
                class="relative z-10 h-7 px-2 sm:px-4 rounded-full hover:bg-transparent"
              >
                <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
                </svg>
                <span class="hidden sm:inline">{{ t('nav.upload') }}</span>
              </Button>
            </button>
          </div>
          
          <!-- 语言切换 - 放在最右侧 -->
          <LanguageSwitcher />
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue'
import { usePublicConfig } from '@/composables/usePublicConfig'
import { useI18n } from '@/composables/useI18n'
import { Button } from '@/components/ui/button'
import ThemeSwitcher from '@/components/ui/ThemeSwitcher.vue'
import LanguageSwitcher from '@/components/ui/LanguageSwitcher.vue'

const props = defineProps<{
  showDescription?: boolean
  activeTab?: 'pickup' | 'upload'
}>()

const emit = defineEmits<{
  'update:activeTab': [value: 'pickup' | 'upload']
}>()

const { config } = usePublicConfig()
const { t, locale } = useI18n()

// 用于动态计算按钮位置和宽度
const buttonGroup = ref<HTMLElement | null>(null)
const pickupButton = ref<HTMLElement | null>(null)
const uploadButton = ref<HTMLElement | null>(null)
const indicatorWidth = ref('0px')
const indicatorPosition = ref('0px')
const isInitialized = ref(false)

// 切换标签页
const switchTab = (tab: 'pickup' | 'upload') => {
  emit('update:activeTab', tab)
}

// 更新指示器位置和宽度
const updateIndicator = async () => {
  // 等待 DOM 更新完成
  await nextTick()
  
  if (!buttonGroup.value || !pickupButton.value || !uploadButton.value) return
  
  const containerRect = buttonGroup.value.getBoundingClientRect()
  let activeButton: HTMLElement | null = null
  
  if (props.activeTab === 'pickup') {
    activeButton = pickupButton.value
  } else if (props.activeTab === 'upload') {
    activeButton = uploadButton.value
  }
  
  if (activeButton) {
    const buttonRect = activeButton.getBoundingClientRect()
    indicatorWidth.value = `${buttonRect.width}px`
    indicatorPosition.value = `${buttonRect.left - containerRect.left}px`
  }
}

// 监听 activeTab 变化
watch(() => props.activeTab, async () => {
  await updateIndicator()
})

// 监听语言变化，重新计算指示器宽度
watch(locale, async () => {
  // 等待翻译更新后再计算指示器
  await nextTick()
  await updateIndicator()
})

// 在组件挂载后初始化
onMounted(async () => {
  // 先设置位置，不显示动画
  await updateIndicator()
  // 使用 requestAnimationFrame 确保位置已经应用到 DOM
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      // 启用动画效果
      isInitialized.value = true
    })
  })
  // 监听窗口大小变化，重新计算位置
  window.addEventListener('resize', updateIndicator)
})

// 组件卸载时清理事件监听
onUnmounted(() => {
  window.removeEventListener('resize', updateIndicator)
})
</script>