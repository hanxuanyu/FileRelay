<template>
  <div class="h-screen flex flex-col bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-950 dark:to-gray-900 transition-colors">
    <!-- 固定导航栏 -->
    <div class="fixed top-0 left-0 right-0 z-50">
      <NavBar :showDescription="true" :activeTab="activeTab" @update:activeTab="activeTab = $event" />
    </div>
    
    <!-- 内容区域 - 从顶部开始滚动 -->
    <div class="flex-1 overflow-auto pt-16">
      <Transition name="fade" mode="out-in">
        <component :is="currentComponent" :key="activeTab" />
      </Transition>
    </div>
    
    <!-- Toast 组件 -->
    <Toaster />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from '@/components/ui/NavBar.vue'
import HomePage from './HomePage.vue'
import UploadPage from './UploadPage.vue'
import { Toaster } from '@/components/ui/sonner'

const route = useRoute()
const router = useRouter()

// 根据路由确定当前激活的标签页
const activeTab = ref<'pickup' | 'upload'>(route.path === '/upload' ? 'upload' : 'pickup')

// 当前显示的组件
const currentComponent = computed(() => {
  return activeTab.value === 'upload' ? UploadPage : HomePage
})

// 监听路由变化
watch(() => route.path, (newPath) => {
  activeTab.value = newPath === '/upload' ? 'upload' : 'pickup'
})

// 监听 activeTab 变化，更新 URL（不刷新页面）
watch(activeTab, (newTab) => {
  const targetPath = newTab === 'upload' ? '/upload' : '/'
  if (route.path !== targetPath) {
    router.replace(targetPath)
  }
})
</script>

<style scoped>
/* 淡入淡出过渡效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
