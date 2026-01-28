<template>
  <AdminLayout>
    <!-- 页面标题 -->
    <div class="mb-6 sm:mb-8">
      <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-gray-100">{{ t('admin.dashboard.title') }}</h2>
      <p class="mt-1 sm:mt-2 text-sm sm:text-base text-gray-600 dark:text-gray-400">{{ t('admin.dashboard.description') }}</p>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 md:gap-6 mb-6 sm:mb-8 relative z-0">
          <!-- 总批次数 -->
          <Card>
            <CardContent class="p-4 sm:p-6">
              <div class="flex flex-col sm:flex-row sm:items-center">
                <div class="flex-shrink-0 mb-2 sm:mb-0">
                  <div class="w-8 h-8 sm:w-10 sm:h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2h4l2 2h4a2 2 0 012 2v2M7 7V3a2 2 0 012-2h6a2 2 0 012 2v4"></path>
                    </svg>
                  </div>
                </div>
                <div class="sm:ml-3 md:ml-4">
                  <p class="text-xs sm:text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.stats.totalBatches') }}</p>
                  <p class="text-lg sm:text-xl md:text-2xl font-semibold text-gray-900 dark:text-gray-100">{{ stats.totalBatches || 0 }}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 活跃批次 -->
          <Card>
            <CardContent class="p-4 sm:p-6">
              <div class="flex flex-col sm:flex-row sm:items-center">
                <div class="flex-shrink-0 mb-2 sm:mb-0">
                  <div class="w-8 h-8 sm:w-10 sm:h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                </div>
                <div class="sm:ml-3 md:ml-4">
                  <p class="text-xs sm:text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.stats.activeBatches') }}</p>
                  <p class="text-lg sm:text-xl md:text-2xl font-semibold text-gray-900 dark:text-gray-100">{{ stats.activeBatches || 0 }}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 已过期批次 -->
          <Card>
            <CardContent class="p-4 sm:p-6">
              <div class="flex flex-col sm:flex-row sm:items-center">
                <div class="flex-shrink-0 mb-2 sm:mb-0">
                  <div class="w-8 h-8 sm:w-10 sm:h-10 bg-yellow-100 dark:bg-yellow-900/30 rounded-lg flex items-center justify-center">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 text-yellow-600 dark:text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                </div>
                <div class="sm:ml-3 md:ml-4">
                  <p class="text-xs sm:text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.stats.expiredBatches') }}</p>
                  <p class="text-lg sm:text-xl md:text-2xl font-semibold text-gray-900 dark:text-gray-100">{{ stats.expiredBatches || 0 }}</p>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 总文件数 -->
          <Card>
            <CardContent class="p-4 sm:p-6">
              <div class="flex flex-col sm:flex-row sm:items-center">
                <div class="flex-shrink-0 mb-2 sm:mb-0">
                  <div class="w-8 h-8 sm:w-10 sm:h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
                    <svg class="w-4 h-4 sm:w-5 sm:h-5 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                    </svg>
                  </div>
                </div>
                <div class="sm:ml-3 md:ml-4">
                  <p class="text-xs sm:text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('admin.dashboard.stats.totalFiles') }}</p>
                  <p class="text-lg sm:text-xl md:text-2xl font-semibold text-gray-900 dark:text-gray-100">{{ stats.totalFiles || 0 }}</p>
                </div>
              </div>
            </CardContent>
      </Card>
    </div>

    <!-- 最近批次 -->
        <Card class="mb-6 sm:mb-8">
          <CardHeader class="px-4 sm:px-6">
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
              <div>
                <CardTitle class="text-base sm:text-lg">{{ t('admin.dashboard.recentBatches.title') }}</CardTitle>
                <CardDescription class="text-xs sm:text-sm">{{ t('admin.dashboard.recentBatches.description') }}</CardDescription>
              </div>
              <Button variant="outline" @click="router.push('/admin/batches')" size="sm" class="w-full sm:w-auto">
                <span class="text-xs sm:text-sm">{{ t('common.viewAll') }}</span>
                <svg class="w-3 h-3 sm:w-4 sm:h-4 ml-1 sm:ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                </svg>
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="loading" class="space-y-4">
              <Skeleton class="h-4 w-full" />
              <Skeleton class="h-4 w-3/4" />
              <Skeleton class="h-4 w-1/2" />
            </div>
            <div v-else-if="recentBatches.length === 0" class="text-center py-8">
              <svg class="w-12 h-12 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
              <p class="text-gray-500">{{ t('admin.dashboard.recentBatches.empty') }}</p>
            </div>
            <div v-else class="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{{ t('admin.batches.list.columns.pickupCode') }}</TableHead>
                    <TableHead>{{ t('common.type') }}</TableHead>
                    <TableHead>{{ t('common.status') }}</TableHead>
                    <TableHead>{{ t('common.createdAt') }}</TableHead>
                    <TableHead>{{ t('common.expireAt') }}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="batch in recentBatches" :key="batch.id" class="dark:hover:bg-black/20">
                    <TableCell class="font-medium font-mono">{{ batch.pickup_code }}</TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {{ batch.type === 'text' ? t('admin.batches.types.text') : t('admin.batches.types.file') }}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge :variant="getStatusVariant(batch.status)">
                        {{ getStatusText(batch.status) }}
                      </Badge>
                    </TableCell>
                    <TableCell>{{ formatDate(batch.created_at) }}</TableCell>
                    <TableCell>
                      <span v-if="batch.expire_at">{{ formatDate(batch.expire_at) }}</span>
                      <span v-else class="text-gray-500">{{ t('pickup.neverExpire') }}</span>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
    </Card>

    <!-- 快速操作 -->
        <Card class="mb-6 sm:mb-8">
          <CardHeader>
            <CardTitle>{{ t('admin.dashboard.quickActions.title') }}</CardTitle>
            <CardDescription>{{ t('admin.dashboard.quickActions.description') }}</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
              <Button @click="router.push('/admin/batches')" variant="outline" class="h-auto p-4 flex-col items-start">
                <svg class="w-8 h-8 text-blue-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2h4l2 2h4a2 2 0 012 2v2M7 7V3a2 2 0 012-2h6a2 2 0 012 2v4"></path>
                </svg>
                <div class="text-left">
                  <div class="font-medium">{{ t('admin.dashboard.quickActions.batchManagement') }}</div>
                  <div class="text-sm text-gray-500">{{ t('admin.dashboard.quickActions.batchManagementDesc') }}</div>
                </div>
              </Button>

              <Button @click="router.push('/admin/tokens')" variant="outline" class="h-auto p-4 flex-col items-start">
                <svg class="w-8 h-8 text-green-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
                </svg>
                <div class="text-left">
                  <div class="font-medium">{{ t('admin.dashboard.quickActions.tokenManagement') }}</div>
                  <div class="text-sm text-gray-500">{{ t('admin.dashboard.quickActions.tokenManagementDesc') }}</div>
                </div>
              </Button>

              <Button @click="router.push('/admin/config')" variant="outline" class="h-auto p-4 flex-col items-start">
                <svg class="w-8 h-8 text-orange-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
                <div class="text-left">
                  <div class="font-medium">{{ t('admin.dashboard.quickActions.systemConfig') }}</div>
                  <div class="text-sm text-gray-500">{{ t('admin.dashboard.quickActions.systemConfigDesc') }}</div>
                </div>
              </Button>

              <Button @click="refreshData" variant="outline" class="h-auto p-4 flex-col items-start">
                <svg class="w-8 h-8 text-purple-600 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                <div class="text-left">
                  <div class="font-medium">{{ t('admin.dashboard.quickActions.refreshData') }}</div>
                  <div class="text-sm text-gray-500">{{ t('admin.dashboard.quickActions.refreshDataDesc') }}</div>
                </div>
              </Button>
            </div>
          </CardContent>
    </Card>

    <!-- Toast 组件 -->
    <Toaster />
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { toast } from 'vue-sonner'

// 组件导入
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Toaster } from '@/components/ui/sonner'
import AdminLayout from '@/layouts/AdminLayout.vue'

// API 和工具导入
import { adminApi, utils, type FileBatch } from '@/lib/api'

const router = useRouter()
const { t } = useI18n()

// 响应式数据
const loading = ref(true)
const recentBatches = ref<FileBatch[]>([])
const stats = ref({
  totalBatches: 0,
  activeBatches: 0,
  expiredBatches: 0,
  totalFiles: 0
})

// 方法
const fetchData = async () => {
  loading.value = true
  
  try {
    // 获取最近的批次数据
    const response = await adminApi.getBatches({ page: 1, page_size: 10 })
    
    if (response.data.code === 200) {
      const batches = response.data.data.data
      recentBatches.value = batches
      
      // 计算统计数据
      stats.value.totalBatches = response.data.data.total
      stats.value.activeBatches = batches.filter(batch => batch.status === 'active').length
      stats.value.expiredBatches = batches.filter(batch => batch.status === 'expired').length
      stats.value.totalFiles = batches.reduce((sum, batch) => sum + (batch.file_items?.length || 0), 0)
    }
  } catch (error: any) {
    console.error('获取数据失败:', error)
    toast.error('获取数据失败，请重试')
  } finally {
    loading.value = false
  }
}

const refreshData = async () => {
  toast.info('正在刷新数据...')
  await fetchData()
  toast.success('数据已刷新')
}

const getStatusVariant = (status: string) => {
  switch (status) {
    case 'active':
      return 'default'
    case 'expired':
      return 'secondary'
    case 'deleted':
      return 'destructive'
    default:
      return 'outline'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'active':
      return t('admin.dashboard.statusTexts.active')
    case 'expired':
      return t('admin.dashboard.statusTexts.expired')
    case 'deleted':
      return t('admin.dashboard.statusTexts.deleted')
    default:
      return t('common.status', '未知')
  }
}

const formatDate = (dateString: string): string => {
  return utils.formatDate(dateString)
}

// 组件挂载
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 自定义样式 */
.border-b-2 {
  border-bottom-width: 2px;
}

/* 导航链接悬停效果 */
.hover\:border-gray-300:hover {
  border-color: #d1d5db;
}

.hover\:text-gray-700:hover {
  color: #374151;
}

/* 卡片悬停效果 */
.hover\:shadow-lg:hover {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.3s ease;
}

/* 按钮组悬停效果 */
.flex-col:hover {
  transform: translateY(-2px);
  transition: transform 0.2s ease;
}

/* 表格行悬停效果 */
.hover\:bg-gray-50:hover {
  background-color: #f9fafb;
  transition: background-color 0.2s ease;
}

/* 统计卡片动画 */
@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.grid > * {
  animation: slideInUp 0.3s ease forwards;
}

.grid > *:nth-child(1) { animation-delay: 0.1s; }
.grid > *:nth-child(2) { animation-delay: 0.2s; }
.grid > *:nth-child(3) { animation-delay: 0.3s; }
.grid > *:nth-child(4) { animation-delay: 0.4s; }
</style>