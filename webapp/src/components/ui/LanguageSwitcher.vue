<template>
  <div class="relative z-50">
    <button
      @click="toggleDropdown"
      class="flex items-center space-x-1 px-2 py-1.5 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors flex-shrink-0"
      :title="t('language.switch', '切换语言')"
    >
      <svg class="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"></path>
      </svg>
      <span class="text-sm font-medium hidden sm:inline">{{ currentLanguageName }}</span>
      <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
      </svg>
    </button>

    <!-- 下拉菜单 - 使用最高的 z-index 确保不被任何元素覆盖 -->
    <Transition
      enter-active-class="transition ease-out duration-100"
      enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100"
      leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100"
      leave-to-class="transform opacity-0 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute right-0 mt-2 w-40 bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 py-1"
        style="z-index: 99999;"
        @click="isOpen = false"
      >
        <button
          v-for="lang in languages"
          :key="lang.code"
          @click="switchLanguage(lang.code)"
          class="w-full px-4 py-2 text-left text-sm hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors flex items-center justify-between"
          :class="{
            'text-blue-600 dark:text-blue-400 font-medium': locale === lang.code,
            'text-gray-700 dark:text-gray-300': locale !== lang.code
          }"
        >
          <span class="flex items-center">
            <span class="mr-2">{{ lang.flag }}</span>
            <span>{{ lang.name }}</span>
          </span>
          <svg
            v-if="locale === lang.code"
            class="w-4 h-4"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"></path>
          </svg>
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { languages, getLanguageName } from '@/i18n/languages'

const { locale, setLocale, t } = useI18n()
const isOpen = ref(false)

const currentLanguageName = computed(() => {
  return getLanguageName(locale.value)
})

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

const switchLanguage = (langCode: string) => {
  setLocale(langCode)
  isOpen.value = false
}

// 点击外部关闭下拉菜单
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (isOpen.value && !target.closest('.relative')) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
/* 确保下拉菜单在其他元素之上 */
.z-50 {
  z-index: 50;
}
</style>
