<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <div class="mx-auto h-12 w-12 flex items-center justify-center rounded-full bg-indigo-100 dark:bg-indigo-900">
          <svg class="h-8 w-8 text-indigo-600 dark:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
          </svg>
        </div>
        <h2 class="mt-6 text-center text-3xl font-bold text-gray-900 dark:text-gray-100">
          {{ t('admin.login.title') }}
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          {{ t('admin.login.description') }}
        </p>
      </div>

      <Card class="border-0 shadow-lg">
        <CardContent class="pt-6">
          <form @submit.prevent="handleLogin" class="space-y-6">
            <div>
              <Label for="password" class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('admin.login.passwordLabel') }}
              </Label>
              <div class="mt-1 relative">
                <Input
                  id="password"
                  name="password"
                  :type="showPassword ? 'text' : 'password'"
                  v-model="password"
                  autocomplete="current-password"
                  :placeholder="t('admin.login.passwordPlaceholder')"
                  class="pr-10"
                  required
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  <svg v-if="showPassword" class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L6.879 6.879a3 3 0 00-4.243 4.243M9.878 9.878a3 3 0 014.242 4.243M15.121 14.121L18.121 17.121a3 3 0 01-4.243 4.243M12 3c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                  </svg>
                  <svg v-else class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                </button>
              </div>
            </div>

            <div>
              <Button
                type="submit"
                :disabled="!password || loading"
                class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
                  <svg class="animate-spin h-5 w-5 text-indigo-300" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </span>
                <svg v-else class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
                </svg>
                {{ loading ? t('admin.login.loggingIn') : t('admin.login.loginButton') }}
              </Button>
            </div>

            <!-- 错误信息显示 -->
            <div v-if="error" class="rounded-md bg-red-50 p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </div>
                <div class="ml-3">
                  <p class="text-sm font-medium text-red-800">
                    {{ error }}
                  </p>
                </div>
                <div class="ml-auto pl-3">
                  <div class="-mx-1.5 -my-1.5">
                    <button
                      @click="error = ''"
                      class="inline-flex bg-red-50 rounded-md p-1.5 text-red-500 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-red-50 focus:ring-red-600"
                    >
                      <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      <!-- 返回首页链接 -->
      <div class="text-center">
        <Button variant="link" @click="router.push('/')" class="text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
          {{ t('nav.backToHome') }}
        </Button>
      </div>
    </div>

    <!-- Toast 组件 -->
    <Toaster />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useI18n } from '@/composables/useI18n'

// 组件导入
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import { Toaster } from '@/components/ui/sonner'

// API 导入
import { adminApi } from '@/lib/api'

const router = useRouter()
const { t } = useI18n()

// 响应式数据
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

// 方法
const handleLogin = async () => {
  if (!password.value.trim()) {
    error.value = t('admin.login.passwordRequired')
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await adminApi.login(password.value.trim())

    if (response.data.code === 200) {
      // 保存 token
      localStorage.setItem('admin_token', response.data.data.token)
      
      toast.success(t('admin.login.loginSuccess'))
      
      // 跳转到管理后台
      router.push('/admin')
    } else {
      throw new Error(response.data.msg || t('admin.login.loginFailed'))
    }
  } catch (err: any) {
    console.error('登录失败:', err)
    
    if (err.response?.status === 401) {
      error.value = t('admin.login.passwordError')
    } else {
      error.value = err.response?.data?.msg || t('admin.login.loginFailed')
    }
    
    // 清空密码
    password.value = ''
  } finally {
    loading.value = false
  }
}

// 检查是否已登录
const checkAuth = () => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    // 如果已有 token，直接跳转到管理后台
    router.replace('/admin')
  }
}

// 组件挂载
onMounted(() => {
  checkAuth()
})
</script>

<style scoped>
/* 自定义样式 */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 表单焦点效果 */
.group:hover .group-hover\:text-indigo-400 {
  color: #818cf8;
}

/* 输入框样式增强 */
input[type="password"]::-ms-reveal,
input[type="password"]::-ms-clear {
  display: none;
}

/* 卡片阴影 */
.shadow-lg {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

/* 按钮悬停效果 */
.hover\:bg-indigo-700:hover {
  background-color: #4338ca;
  transition: background-color 0.2s ease;
}

/* 错误消息动画 */
.bg-red-50 {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 聚焦环效果 */
.focus\:ring-2:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(99, 102, 241, 0.2);
}

.focus\:ring-offset-2:focus {
  outline: none;
  box-shadow: 0 0 0 2px #fff, 0 0 0 4px rgba(99, 102, 241, 0.2);
}
</style>