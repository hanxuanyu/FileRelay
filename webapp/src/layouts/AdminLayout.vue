<template>
  <div class="h-screen flex flex-col bg-gray-50 dark:bg-gray-950 transition-colors">
    <!-- 固定导航栏 -->
    <div class="fixed top-0 left-0 right-0 z-50">
      <AdminNavBar />
    </div>
    
    <!-- 内容区域 - 从顶部开始滚动 -->
    <div class="flex-1 overflow-auto pt-16">
      <div class="max-w-7xl mx-auto py-4 sm:py-6 px-3 sm:px-6 lg:px-8">
        <div class="sm:px-0">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import AdminNavBar from '@/components/ui/AdminNavBar.vue'
</script>