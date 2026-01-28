<template>
  <nav class="bg-white/75 dark:bg-gray-900/75 backdrop-blur-lg border-b border-gray-200 dark:border-gray-700 shadow-sm transition-colors">
    <div class="max-w-7xl mx-auto px-2 sm:px-4 lg:px-8">
      <div class="flex justify-between h-14 sm:h-16">
        <!-- 左侧 Logo 和菜单 -->
        <div class="flex items-center min-w-0 flex-1">
          <div class="flex-shrink-0 flex items-center space-x-1.5 sm:space-x-2">
            <!-- Logo -->
            <div v-if="config.site?.logo" class="w-7 h-7 sm:w-8 sm:h-8 rounded-lg overflow-hidden flex items-center justify-center flex-shrink-0">
              <img :src="config.site.logo" :alt="config.site?.name || t('site.title')" class="w-full h-full object-contain" />
            </div>
            <div v-else class="w-7 h-7 sm:w-8 sm:h-8 bg-indigo-600 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4 sm:w-5 sm:h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
            </div>
            <!-- 站点名称和徽章 -->
            <div class="flex items-center space-x-1.5 min-w-0">
              <h1 class="hidden sm:block text-lg lg:text-xl font-bold text-gray-900 dark:text-gray-100 truncate">
                {{ config.site?.name || t('site.title') }}
              </h1>
              <Badge variant="secondary" class="text-[10px] sm:text-xs px-1.5 py-0.5 flex-shrink-0">ADMIN</Badge>
            </div>
          </div>
          
          <!-- 桌面端导航菜单 -->
          <div class="hidden lg:ml-6 lg:flex lg:space-x-4 xl:space-x-6 min-w-0">
            <router-link
              to="/admin"
              class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium whitespace-nowrap"
              :class="$route.path === '/admin' ? 'border-indigo-500 text-gray-900 dark:text-gray-100' : 'border-transparent text-gray-500 dark:text-gray-400 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-300'"
            >
              {{ t('admin.nav.overview') }}
            </router-link>
            <router-link
              to="/admin/batches"
              class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium whitespace-nowrap"
              :class="$route.path.includes('/admin/batches') ? 'border-indigo-500 text-gray-900 dark:text-gray-100' : 'border-transparent text-gray-500 dark:text-gray-400 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-300'"
            >
              {{ t('admin.nav.fileManagement') }}
            </router-link>
            <router-link
              to="/admin/tokens"
              class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium whitespace-nowrap"
              :class="$route.path.includes('/admin/tokens') ? 'border-indigo-500 text-gray-900 dark:text-gray-100' : 'border-transparent text-gray-500 dark:text-gray-400 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-300'"
            >
              {{ t('admin.nav.apiManagement') }}
            </router-link>
            <router-link
              to="/admin/config"
              class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium whitespace-nowrap"
              :class="$route.path.includes('/admin/config') ? 'border-indigo-500 text-gray-900 dark:text-gray-100' : 'border-transparent text-gray-500 dark:text-gray-400 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-300'"
            >
              {{ t('admin.nav.systemConfig') }}
            </router-link>
          </div>
        </div>
        
        <!-- 右侧操作按钮 -->
        <div class="flex items-center space-x-0.5 sm:space-x-1 flex-shrink-0">
          <!-- 主题切换 -->
          <ThemeSwitcher class="flex-shrink-0" />
          
          <!-- 语言切换 -->
          <LanguageSwitcher class="flex-shrink-0" />
          
          <!-- 移动端菜单按钮 -->
          <Button 
            variant="ghost" 
            @click="showMobileMenu = !showMobileMenu" 
            class="lg:hidden flex-shrink-0" 
            size="sm"
            :aria-label="t('admin.nav.toggleMenu')"
          >
            <svg 
              class="w-5 h-5 transition-transform duration-200" 
              :class="{ 'rotate-90': showMobileMenu }"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path 
                v-if="!showMobileMenu"
                stroke-linecap="round" 
                stroke-linejoin="round" 
                stroke-width="2" 
                d="M4 6h16M4 12h16M4 18h16"
              ></path>
              <path 
                v-else
                stroke-linecap="round" 
                stroke-linejoin="round" 
                stroke-width="2" 
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </Button>
          
          <!-- 桌面端操作按钮 -->
          <div class="hidden lg:flex lg:items-center lg:space-x-1">
            <Button variant="outline" @click="router.push('/')" size="sm" class="flex-shrink-0">
              <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-1a1 1 0 011-1h2a1 1 0 011 1v1a1 1 0 001 1m-6 0h6"></path>
              </svg>
              <span class="text-sm">{{ t('admin.nav.goToFront') }}</span>
            </Button>
            <Button variant="destructive" @click="handleLogout" size="sm" class="flex-shrink-0">
              <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
              </svg>
              <span class="text-sm">{{ t('admin.nav.logout') }}</span>
            </Button>
          </div>
        </div>
      </div>
      
      <!-- 移动端菜单 -->
      <transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 -translate-y-1"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-1"
      >
        <div class="lg:hidden" v-if="showMobileMenu">
          <div class="pt-2 pb-3 space-y-1 bg-gray-50/50 dark:bg-gray-800/50 backdrop-blur-sm">
            <router-link
              to="/admin"
              class="flex items-center pl-3 pr-4 py-2.5 text-sm sm:text-base font-medium border-l-4 transition-all"
              :class="$route.path === '/admin' ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-700 dark:text-indigo-400' : 'border-transparent text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-800 dark:hover:text-gray-200'"
              @click="showMobileMenu = false"
            >
              <svg class="w-5 h-5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
              </svg>
              <span class="truncate">{{ t('admin.nav.overview') }}</span>
            </router-link>
            <router-link
              to="/admin/batches"
              class="flex items-center pl-3 pr-4 py-2.5 text-sm sm:text-base font-medium border-l-4 transition-all"
              :class="$route.path.includes('/admin/batches') ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-700 dark:text-indigo-400' : 'border-transparent text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-800 dark:hover:text-gray-200'"
              @click="showMobileMenu = false"
            >
              <svg class="w-5 h-5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
              <span class="truncate">{{ t('admin.nav.fileManagement') }}</span>
            </router-link>
            <router-link
              to="/admin/tokens"
              class="flex items-center pl-3 pr-4 py-2.5 text-sm sm:text-base font-medium border-l-4 transition-all"
              :class="$route.path.includes('/admin/tokens') ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-700 dark:text-indigo-400' : 'border-transparent text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-800 dark:hover:text-gray-200'"
              @click="showMobileMenu = false"
            >
              <svg class="w-5 h-5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
              </svg>
              <span class="truncate">{{ t('admin.nav.apiManagement') }}</span>
            </router-link>
            <router-link
              to="/admin/config"
              class="flex items-center pl-3 pr-4 py-2.5 text-sm sm:text-base font-medium border-l-4 transition-all"
              :class="$route.path.includes('/admin/config') ? 'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-500 text-indigo-700 dark:text-indigo-400' : 'border-transparent text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-800 dark:hover:text-gray-200'"
              @click="showMobileMenu = false"
            >
              <svg class="w-5 h-5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
              <span class="truncate">{{ t('admin.nav.systemConfig') }}</span>
            </router-link>
            
            <!-- 移动端快捷操作按钮 -->
            <div class="pt-2 pb-2 px-3 space-y-2">
              <Button variant="outline" @click="router.push('/'); showMobileMenu = false" class="w-full justify-start" size="sm">
                <svg class="w-4 h-4 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-1a1 1 0 011-1h2a1 1 0 011 1v1a1 1 0 001 1m-6 0h6"></path>
                </svg>
                <span class="truncate">{{ t('admin.nav.goToFront') }}</span>
              </Button>
              <Button variant="destructive" @click="handleLogout" class="w-full justify-start" size="sm">
                <svg class="w-4 h-4 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path>
                </svg>
                <span class="truncate">{{ t('admin.nav.logout') }}</span>
              </Button>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { usePublicConfig } from '@/composables/usePublicConfig'
import ThemeSwitcher from '@/components/ui/ThemeSwitcher.vue'
import LanguageSwitcher from '@/components/ui/LanguageSwitcher.vue'
import { useI18n } from '@/composables/useI18n'

const router = useRouter()
const showMobileMenu = ref(false)
const { config } = usePublicConfig()
const { t } = useI18n()

const handleLogout = () => {
  localStorage.removeItem('admin_token')
  router.push('/admin/login')
}
</script>