<template>
  <div class="w-full h-full">
    <div class="container mx-auto px-3 sm:px-4 pb-6 sm:pb-12 h-full flex items-center">
        <!-- 主要内容区域 -->
        <div class="max-w-2xl mx-auto w-full">
        <!-- 取件区域 -->
        <Card v-if="!batchData" class="shadow-2xl border-0 overflow-hidden" @click="handleCardClick">
          <CardContent class="pt-8 sm:pt-12 pb-8 sm:pb-10 px-4 sm:px-6">
            <div class="text-center mb-8 sm:mb-10">
              <div class="inline-flex items-center justify-center w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 mb-4 sm:mb-6 shadow-lg">
                <svg class="w-8 h-8 sm:w-10 sm:h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
                </svg>
              </div>
              <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2 sm:mb-3">{{ t('pickup.title') }}</h2>
              <p class="text-sm sm:text-base text-gray-600 dark:text-gray-400">{{ t('pickup.description') }}</p>
            </div>
            
            <div class="flex flex-col items-center space-y-4 sm:space-y-6">
              <!-- InputOTP 组件 -->
              <div ref="otpContainerRef" class="flex justify-center p-4 sm:p-6 bg-gradient-to-br from-blue-50 to-purple-50 dark:from-gray-800 dark:to-gray-800 rounded-2xl shadow-inner w-full max-w-sm">
                <InputOTP
                  v-model="pickupCode"
                  :maxlength="pickupCodeLength"
                  :disabled="loading"
                >
                  <InputOTPGroup>
                    <InputOTPSlot
                      v-for="index in pickupCodeLength"
                      :key="index"
                      :index="index - 1"
                    />
                  </InputOTPGroup>
                </InputOTP>
              </div>
              
              <!-- 加载提示 -->
              <div v-if="loading" class="text-center text-sm sm:text-base text-blue-600 flex items-center animate-pulse">
                <svg class="animate-spin -ml-1 mr-2 sm:mr-3 h-4 w-4 sm:h-5 sm:w-5" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ t('pickup.fetching') }}
              </div>
              
              <div v-else class="text-center">
                <div class="inline-flex items-center px-4 sm:px-5 py-2 sm:py-2.5 rounded-full transition-all" :class="pickupCode === '' ? 'bg-gray-100 dark:bg-gray-800' : 'bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/30 dark:to-purple-900/30'">
                  <svg v-if="pickupCode === ''" class="w-4 h-4 sm:w-5 sm:h-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  <svg v-else class="w-4 h-4 sm:w-5 sm:h-5 mr-2 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
                  </svg>
                  <span class="text-xs sm:text-sm font-medium" :class="pickupCode === '' ? 'text-gray-600 dark:text-gray-400' : 'text-blue-700 dark:text-blue-300'">
                    <template v-if="pickupCode === ''">
                      {{ t('pickup.inputPlaceholder', { length: pickupCodeLength }) }}
                    </template>
                    <template v-else>
                      {{ t('pickup.inputProgress', { current: pickupCode.length, total: pickupCodeLength }) }}
                    </template>
                  </span>
                </div>
              </div>
              
              <!-- 粘贴和清空按钮 -->
              <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-center gap-2 sm:gap-3 w-full sm:w-auto px-4 sm:px-0">
                <Button 
                  variant="ghost" 
                  @click="pasteFromClipboard" 
                  size="default"
                  class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-950/30 transition-all text-sm sm:text-base"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
                  </svg>
                  {{ t('pickup.pasteFromClipboard') }}
                </Button>
                
                <!-- 清空按钮 -->
                <Button 
                  v-if="pickupCode" 
                  variant="ghost" 
                  @click="clearAndFocus" 
                  size="default"
                  class="text-gray-600 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-950/30 transition-all text-sm sm:text-base"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                  {{ t('pickup.clearAndRetry') }}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- 文件展示区域 -->
        <div v-if="batchData" class="space-y-3 sm:space-y-4">
          <!-- 批次信息 -->
          <Card class="shadow-md">
            <CardHeader class="pb-3 sm:pb-4 px-4 sm:px-6">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-3">
                <div class="flex-1 min-w-0">
                  <CardTitle class="flex items-center text-sm sm:text-base mb-1.5">
                    <svg class="w-4 h-4 mr-2 text-green-600 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                    {{ t('pickup.fetchSuccess') }}
                  </CardTitle>
                  <div class="flex items-center gap-2 flex-wrap">
                    <CardDescription class="text-xs">
                      {{ t('pickup.pickupCode') }}: {{ currentCode }}
                    </CardDescription>
                    <Badge :variant="batchData.download_count < (batchData.max_downloads || Infinity) ? 'default' : 'secondary'" class="text-xs">
                      {{ batchData.type === 'text' ? t('pickup.textContent') : `${batchData.files?.length || 0}${t('pickup.fileCount')}` }}
                    </Badge>
                  </div>
                </div>
                <Button variant="outline" size="sm" @click="reset" class="h-8 text-xs w-full sm:w-auto flex-shrink-0">
                  <svg class="w-3.5 h-3.5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                  </svg>
                  {{ t('common.reset') }}
                </Button>
              </div>
            </CardHeader>
            <CardContent class="pt-0 px-4 sm:px-6">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
                <div class="flex items-center gap-2">
                  <span class="text-gray-600 dark:text-gray-400 font-medium min-w-[60px]">{{ t('common.type') }}</span>
                  <span class="text-gray-900 dark:text-gray-100">{{ batchData.type === 'text' ? t('admin.batches.types.text') : t('admin.batches.types.file') }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-gray-600 dark:text-gray-400 font-medium min-w-[60px]">{{ t('pickup.downloadCount') }}</span>
                  <span class="text-gray-900 dark:text-gray-100">{{ batchData.download_count }}{{ batchData.max_downloads ? ` / ${batchData.max_downloads}` : '' }}</span>
                </div>
                <div v-if="batchData.expire_at" class="flex items-center gap-2">
                  <span class="text-gray-600 dark:text-gray-400 font-medium min-w-[60px]">{{ t('common.expireAt') }}</span>
                  <span class="text-gray-900 dark:text-gray-100 text-xs break-all">{{ formatDate(batchData.expire_at) }}</span>
                </div>
                <div v-if="batchData.remark" class="flex items-start gap-2 sm:col-span-2">
                  <span class="text-gray-600 dark:text-gray-400 font-medium min-w-[60px] flex-shrink-0">{{ t('admin.batches.detail.remark', '备注') }}</span>
                  <span class="text-gray-900 dark:text-gray-100 break-words">{{ batchData.remark }}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          <!-- 文本内容 -->
          <Card v-if="batchData.type === 'text'" class="shadow-md">
            <CardHeader class="pb-3 px-4 sm:px-6">
              <div class="flex items-center justify-between">
                <CardTitle class="text-sm sm:text-base">{{ t('pickup.textContent') }}</CardTitle>
                <Button @click="copyText" size="sm" class="h-8 text-xs">
                  <svg class="w-3.5 h-3.5 mr-1 sm:mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                  </svg>
                  <span class="hidden sm:inline">{{ t('pickup.copyText') }}</span>
                </Button>
              </div>
            </CardHeader>
            <CardContent class="pt-3 px-4 sm:px-6">
              <div class="bg-gray-50 dark:bg-gray-800 p-3 rounded-lg border border-gray-200 dark:border-gray-700 max-h-60 sm:max-h-80 overflow-y-auto transition-colors">
                <pre class="whitespace-pre-wrap text-xs sm:text-sm text-gray-800 dark:text-gray-100 break-words">{{ batchData.content }}</pre>
              </div>
            </CardContent>
          </Card>

          <!-- 文件列表 -->
          <Card v-if="batchData.type === 'file' && batchData.files" class="shadow-md">
            <CardHeader class="pb-3 px-4 sm:px-6">
              <div class="flex items-center justify-between gap-2">
                <div class="flex-1 min-w-0">
                  <CardTitle class="text-sm sm:text-base">{{ t('pickup.fileList') }}</CardTitle>
                  <CardDescription class="text-xs">
                    {{ t('common.total', '共') }} {{ batchData.files.length }} {{ t('admin.batches.detail.fileList', '个文件') }}，{{ totalFileSize }}
                  </CardDescription>
                </div>
                <Button 
                  @click="downloadAll" 
                  :disabled="downloading"
                  size="sm"
                  class="h-8 text-xs flex-shrink-0"
                >
                  <svg v-if="downloading" class="animate-spin -ml-1 mr-1 sm:mr-1.5 h-3.5 w-3.5" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="w-3.5 h-3.5 mr-1 sm:mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3M7 7h10a2 2 0 012 2v6a2 2 0 01-2 2H7a2 2 0 01-2-2V9a2 2 0 012-2z"></path>
                  </svg>
                  <span class="hidden sm:inline">{{ downloading ? t('pickup.downloading') : t('pickup.downloadAll') }}</span>
                  <span class="sm:hidden">{{ downloading ? t('common.loading') : t('common.download', '下载') }}</span>
                </Button>
              </div>
            </CardHeader>
            <CardContent class="pt-3 px-4 sm:px-6">
              <div class="space-y-2">
                <div
                  v-for="file in batchData.files"
                  :key="file.id"
                  class="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                >
                  <div class="flex items-center space-x-3 flex-1">
                    <div class="w-9 h-9 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
                      <component :is="getFileIconComponent(file.original_name)" class="w-5 h-5 text-blue-600 dark:text-blue-400" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <h4 class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                        {{ file.original_name }}
                      </h4>
                      <div class="flex items-center space-x-3 text-xs text-gray-500 mt-0.5">
                        <span>{{ formatFileSize(file.size) }}</span>
                        <span>{{ file.mime_type }}</span>
                      </div>
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <Button 
                      @click="downloadFile(file)" 
                      :disabled="downloadingFiles.has(file.id)"
                      size="sm"
                      class="h-8 text-xs"
                    >
                      <svg v-if="downloadingFiles.has(file.id)" class="animate-spin h-3.5 w-3.5" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                      </svg>
                    </Button>
                    
                    <Popover>
                      <PopoverTrigger as-child>
                        <Button 
                          variant="outline"
                          size="sm"
                          class="h-8 w-8 p-0"
                          title="复制下载命令"
                        >
                          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                          </svg>
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent class="w-80" align="end">
                        <div class="space-y-3">
                          <div>
                            <h4 class="font-medium text-sm mb-2">{{ t('pickup.copyDownloadCommand', '复制下载命令') }}</h4>
                            <p class="text-xs text-gray-500 mb-3">{{ t('pickup.selectTool', '选择一种命令行工具复制下载命令') }}</p>
                          </div>
                          
                          <div class="space-y-2">
                            <Button 
                              variant="outline" 
                              size="sm"
                              class="w-full justify-start text-xs font-mono"
                              @click="copyToClipboard(getDownloadCommands(file).url, t('pickup.downloadLink', '下载链接'))"
                            >
                              <svg class="w-3.5 h-3.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path>
                              </svg>
                              复制 URL
                            </Button>
                            
                            <Button 
                              variant="outline" 
                              size="sm"
                              class="w-full justify-start text-xs font-mono"
                              @click="copyToClipboard(getDownloadCommands(file).wget, 'wget 命令')"
                            >
                              <svg class="w-3.5 h-3.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                              </svg>
                              wget
                            </Button>
                            
                            <Button 
                              variant="outline" 
                              size="sm"
                              class="w-full justify-start text-xs font-mono"
                              @click="copyToClipboard(getDownloadCommands(file).curl, 'curl 命令')"
                            >
                              <svg class="w-3.5 h-3.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                              </svg>
                              curl
                            </Button>
                            
                            <Button 
                              variant="outline" 
                              size="sm"
                              class="w-full justify-start text-xs font-mono"
                              @click="copyToClipboard(getDownloadCommands(file).powershell, 'PowerShell 命令')"
                            >
                              <svg class="w-3.5 h-3.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                              </svg>
                              PowerShell
                            </Button>
                            
                            <Button 
                              variant="outline" 
                              size="sm"
                              class="w-full justify-start text-xs font-mono"
                              @click="copyToClipboard(getDownloadCommands(file).aria2c, 'aria2c 命令')"
                            >
                              <svg class="w-3.5 h-3.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"></path>
                              </svg>
                              aria2c
                            </Button>
                          </div>
                        </div>
                      </PopoverContent>
                    </Popover>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { toast } from 'vue-sonner'
import { useI18n } from '@/composables/useI18n'

// Props 定义 - 支持通过路由传递取件码
interface Props {
  code?: string
}
const props = defineProps<Props>()

// 路由和i18n
const route = useRoute()
const { t } = useI18n()

// 组件导入
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { InputOTP, InputOTPGroup, InputOTPSlot } from '@/components/ui/input-otp'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
// API 和工具导入
import { publicApi, utils, type PickupResponse, type FileItem } from '@/lib/api'
import { usePublicConfig } from '@/composables/usePublicConfig'

const { config } = usePublicConfig()

// 响应式数据
const pickupCode = ref('')
const currentCode = ref('')
const batchData = ref<PickupResponse | null>(null)
const loading = ref(false)
const downloading = ref(false)
const downloadingFiles = ref(new Set<string>())
const otpContainerRef = ref<HTMLElement | null>(null)

// 检测是否为移动设备
const isMobile = computed(() => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) ||
    window.innerWidth <= 768
})

// 计算取件码长度
const pickupCodeLength = computed(() => config.value.security?.pickup_code_length || 6)

// 计算属性
const totalFileSize = computed(() => {
  if (!batchData.value?.files) return '0 B'
  const total = batchData.value.files.reduce((sum, file) => sum + file.size, 0)
  return utils.formatFileSize(total)
})

// 方法
const handlePickup = async () => {
  if (!pickupCode.value.trim()) {
    toast.warning(t('pickup.inputPlaceholder', { length: pickupCodeLength.value }))
    return
  }

  loading.value = true
  try {
    const response = await publicApi.getBatch(pickupCode.value.trim())
    
    if (response.data.code === 200) {
      batchData.value = response.data.data
      currentCode.value = pickupCode.value.trim()
      toast.success(t('pickup.fetchSuccess'))
    } else {
      throw new Error(response.data.msg || t('pickup.fetchFailed'))
    }
  } catch (error: any) {
    // 查询失败后清空输入框，方便用户重新输入
    const errorCode = pickupCode.value
    pickupCode.value = ''
    
    // 不在控制台显示业务异常（如404等）
    if (error.response?.status >= 400 && error.response?.status < 500) {
      // 客户端错误，只显示toast，不打印到控制台
      if (error.response?.status === 404) {
        toast.error(t('pickup.invalidCode'), {
          duration: 4000,
          description: h('span', { class: 'text-gray-700 dark:text-gray-200' }, `${t('pickup.pickupCode')}：${errorCode}`)
        })
      } else if (error.response?.status === 410) {
        toast.error(t('pickup.expired'), {
          duration: 4000,
          description: h('span', { class: 'text-gray-700 dark:text-gray-200' }, `${t('pickup.pickupCode')}：${errorCode}`)
        })
      } else {
        toast.error(error.response?.data?.msg || t('pickup.fetchFailed'), {
          duration: 3000
        })
      }
    } else {
      // 服务器错误或网络错误，记录到控制台用于调试
      console.error('获取批次失败:', error)
      toast.error(error.response?.data?.msg || t('common.serverError', '服务器异常，请稍后重试'), {
        duration: 3000
      })
    }
    
    // 错误后自动聚焦输入框，方便用户重新输入
    focusOtpInput()
  } finally {
    loading.value = false
  }
}

const pasteFromClipboard = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text.trim()) {
      const trimmedText = text.trim()
      pickupCode.value = trimmedText
      toast.success(t('common.pasteSuccess', '已粘贴取件码'))
      
      // watch监听会自动触发查询，无需手动调用handlePickup
      // 如果长度匹配，watch会立即触发handlePickup()
    } else {
      toast.warning(t('common.clipboardEmpty', '剪贴板中没有内容'))
    }
  } catch (error) {
    toast.error(t('common.clipboardFailed', '读取剪贴板失败，请手动输入'))
  }
}

const copyText = async () => {
  if (batchData.value?.content) {
    const success = await utils.copyToClipboard(batchData.value.content)
    if (success) {
      toast.success(t('pickup.copyText') + t('common.success', '成功'))
    } else {
      toast.error(t('common.copyFailed'))
    }
  }
}

// 刷新下载次数
const refreshDownloadCount = async () => {
  if (!currentCode.value || !batchData.value) return
  
  try {
    const response = await publicApi.getDownloadCount(currentCode.value)
    if (response.data.code === 200) {
      // 只更新下载次数，不重新加载整个批次信息
      batchData.value.download_count = response.data.data.download_count
      if (response.data.data.max_downloads !== undefined) {
        batchData.value.max_downloads = response.data.data.max_downloads
      }
    }
  } catch (error: any) {
    // 刷新失败不显示错误，避免影响用户体验
    // 只有非业务异常才记录到控制台
    if (!error.response || error.response.status >= 500) {
      console.error('刷新下载次数失败:', error)
    }
  }
}

const downloadFile = async (file: FileItem) => {
  downloadingFiles.value.add(file.id)
  try {
    // 如果后端返回了 download_url，使用该 URL 直接下载
    if (file.download_url) {
      // 使用 fetch 下载文件
      const response = await fetch(file.download_url)
      if (!response.ok) {
        // 构造一个类似 axios 的错误对象，以便统一处理
        const error: any = new Error('Download failed')
        error.response = {
          status: response.status,
          data: await response.json().catch(() => ({ msg: 'Download failed' }))
        }
        throw error
      }
      const blob = await response.blob()
      utils.downloadBlob(blob, file.original_name)
    } else {
      // 否则使用原有的 API 方式下载
      const response = await publicApi.downloadFile(file.id)
      utils.downloadBlob(response.data, file.original_name)
    }
    toast.success(`${t('common.download', '下载')} ${file.original_name} ${t('common.success', '成功')}`)
    
    // 仅文件类型的批次下载后刷新，避免文本类型增加查询次数
    // 文件下载成功后使用专用接口刷新下载次数
    if (batchData.value?.type === 'file') {
      // 稍微延迟一下再刷新，确保后端已经更新了下载次数
      await new Promise(resolve => setTimeout(resolve, 100))
      await refreshDownloadCount()
    }
  } catch (error: any) {
    // 不在控制台显示业务异常（如404等）
    if (error.response?.status >= 400 && error.response?.status < 500) {
      // 客户端错误，只显示toast，不打印到控制台
      if (error.response?.status === 404) {
        toast.error(t('admin.batches.errors.batchNotFound', '批次不存在或文件已被删除'))
      } else if (error.response?.status === 410) {
        toast.error(t('pickup.expired'))
      } else {
        toast.error(error.response?.data?.msg || t('common.downloadFailed', '下载失败'))
      }
    } else {
      // 服务器错误或网络错误，记录到控制台用于调试
      console.error('下载文件失败:', error)
      toast.error(t('common.serverError', '服务器异常或网络错误，请重试'))
    }
  } finally {
    downloadingFiles.value.delete(file.id)
  }
}

// 生成下载命令
const getDownloadCommands = (file: FileItem) => {
  // 优先使用后端返回的 download_url，如果不存在则使用前端拼接的 URL
  const baseUrl = window.location.origin
  const downloadUrl = file.download_url || `${baseUrl}/api/files/${file.id}/download`
  const filename = file.original_name
  
  return {
    url: downloadUrl,
    wget: `wget -O "${filename}" "${downloadUrl}"`,
    curl: `curl -o "${filename}" "${downloadUrl}"`,
    powershell: `Invoke-WebRequest -Uri "${downloadUrl}" -OutFile "${filename}"`,
    aria2c: `aria2c -o "${filename}" "${downloadUrl}"`
  }
}

// 复制到剪贴板
const copyToClipboard = async (text: string, label: string) => {
  try {
    await navigator.clipboard.writeText(text)
    toast.success(`${label}已复制到剪贴板`)
  } catch (error) {
    toast.error('复制失败，请手动复制')
  }
}

const downloadAll = async () => {
  if (!currentCode.value) return

  downloading.value = true
  try {
    const response = await publicApi.downloadBatch(currentCode.value)
    const filename = `files_${currentCode.value}.zip`
    utils.downloadBlob(response.data, filename)
    toast.success(t('pickup.downloadAll') + t('common.success', '成功'))
    
    // 仅文件类型的批次下载后刷新，避免文本类型增加查询次数
    // 打包下载成功后使用专用接口刷新下载次数
    if (batchData.value?.type === 'file') {
      // 稍微延迟一下再刷新，确保后端已经更新了下载次数
      await new Promise(resolve => setTimeout(resolve, 100))
      await refreshDownloadCount()
    }
  } catch (error: any) {
    // 不在控制台显示业务异常（如404等）
    if (error.response?.status >= 400 && error.response?.status < 500) {
      // 客户端错误，只显示toast，不打印到控制台
      if (error.response?.status === 404) {
        toast.error(t('admin.batches.errors.batchNotFound', '批次不存在或已被删除'))
      } else if (error.response?.status === 410) {
        toast.error(t('pickup.expired'))
      } else {
        toast.error(error.response?.data?.msg || t('common.downloadFailed', '下载失败'))
      }
    } else {
      // 服务器错误或网络错误，记录到控制台用于调试
      console.error('打包下载失败:', error)
      toast.error(t('common.serverError', '服务器异常或网络错误，请重试'))
    }
  } finally {
    downloading.value = false
  }
}

const reset = async () => {
  batchData.value = null
  currentCode.value = ''
  pickupCode.value = ''
  
  // 重置后自动聚焦到输入框
  focusOtpInput()
}

const formatFileSize = (bytes: number): string => {
  return utils.formatFileSize(bytes)
}

const formatDate = (dateString: string): string => {
  return utils.formatDate(dateString)
}

// 图标组件映射
// 监听取件码变化，自动触发提取
watch(pickupCode, (newCode) => {
  if (newCode.length === pickupCodeLength.value) {
    handlePickup()
  }
})

// 自动聚焦函数
const focusOtpInput = async () => {
  await nextTick()
  // 添加一个小延迟确保 DOM 完全渲染
  setTimeout(() => {
    // 方式1: 通过容器查找任何 input（不限定类型）
    if (otpContainerRef.value) {
      const input = otpContainerRef.value.querySelector('input') as HTMLInputElement
      if (input) {
        input.focus()
        return
      }
    }
    
    // 方式2: 通过 data-slot 属性查找
    const otpContainer = document.querySelector('[data-slot="input-otp"]')
    if (otpContainer) {
      const input = otpContainer.querySelector('input') as HTMLInputElement
      if (input) {
        input.focus()
        return
      }
    }
    
    // 方式3: 查找 OTP 输入槽位内的 input
    const otpSlots = document.querySelector('[data-slot="input-otp-slot"]')
    if (otpSlots) {
      const parent = otpSlots.parentElement?.parentElement
      if (parent) {
        const input = parent.querySelector('input') as HTMLInputElement
        if (input) {
          input.focus()
          return
        }
      }
    }
    
    // 方式4: 全局查找所有 input（不限定类型）
    const allInputs = document.querySelectorAll('input')
    for (const input of allInputs) {
      const htmlInput = input as HTMLInputElement
      // 检查是否是可见且未禁用的输入框
      if (htmlInput.offsetParent !== null && !htmlInput.disabled) {
        htmlInput.focus()
        return
      }
    }
  }, 150)
}

// 组件挂载时处理 URL 参数
onMounted(async () => {
  // 配置已在应用启动时加载，无需重复加载
  
  // 如果 URL 中有取件码参数，自动填充并获取
  const urlCode = props.code || (route.params.code as string)
  if (urlCode) {
    pickupCode.value = urlCode
    // watch 会自动触发 handlePickup
  } else {
    // 如果没有 URL 取件码，自动为输入框获取焦点
    // PC端：立即聚焦
    // 移动端：延迟聚焦，避免页面加载时就弹出键盘影响用户浏览
    if (isMobile.value) {
      // 移动端延迟500ms后聚焦，给用户一个查看页面的时间
      setTimeout(() => {
        // 再次检查是否仍然符合聚焦条件
        if (!batchData.value && !pickupCode.value) {
          focusOtpInput()
        }
      }, 500)
    } else {
      focusOtpInput()
    }
  }
  
  // 添加全局点击事件监听器（仅PC端），点击空白区域时聚焦输入框
  // 移动端不添加此监听器，避免频繁弹出键盘
  if (!isMobile.value) {
    document.addEventListener('click', handleDocumentClick)
  }
  
  // 添加全局焦点监听器（仅PC端），当没有元素获取焦点时自动聚焦输入框
  // 移动端不需要这个功能，因为会导致键盘频繁弹出
  if (!isMobile.value) {
    document.addEventListener('focusin', handleFocusIn)
    document.addEventListener('focusout', handleFocusOut)
  }
})

// 组件卸载时清理事件监听器
onUnmounted(() => {
  if (!isMobile.value) {
    document.removeEventListener('click', handleDocumentClick)
    document.removeEventListener('focusin', handleFocusIn)
    document.removeEventListener('focusout', handleFocusOut)
  }
})

// 处理卡片点击事件
const handleCardClick = () => {
  // 移动端：只在输入框为空时才聚焦，避免用户浏览时意外弹出键盘
  // PC端：任何时候点击卡片都聚焦
  if (!isMobile.value || !pickupCode.value) {
    focusOtpInput()
  }
}

// 处理文档点击事件（仅PC端使用）
const handleDocumentClick = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  
  // 如果已经有批次数据，不处理
  if (batchData.value) return
  
  // 如果点击的是按钮或输入框本身，不处理
  if (target.tagName === 'BUTTON' || target.tagName === 'INPUT' || target.closest('button')) {
    return
  }
  
  // 如果输入框不为空，不自动聚焦
  if (pickupCode.value) return
  
  // 其他情况下自动聚焦输入框
  focusOtpInput()
}

// 处理焦点进入事件（仅PC端使用）
let focusOutTimer: number | null = null
const handleFocusIn = () => {
  // 清除失去焦点的定时器
  if (focusOutTimer) {
    clearTimeout(focusOutTimer)
    focusOutTimer = null
  }
}

// 处理焦点离开事件（仅PC端使用）
const handleFocusOut = () => {
  // 如果已经有批次数据或输入框有内容，不处理
  if (batchData.value || pickupCode.value) return
  
  // 延迟检查，因为焦点可能转移到另一个元素
  focusOutTimer = window.setTimeout(() => {
    // 检查当前是否有元素获取焦点
    const activeElement = document.activeElement
    
    // 如果没有元素获取焦点（activeElement 是 body），且输入框为空，则自动聚焦
    if (activeElement === document.body && !pickupCode.value && !batchData.value) {
      focusOtpInput()
    }
  }, 100)
}

// 清空输入框并自动聚焦
const clearAndFocus = () => {
  pickupCode.value = ''
  focusOtpInput()
}

const getFileIconComponent = (filename: string) => {
  const icon = utils.getFileTypeIcon(filename)
  
  const iconMap: Record<string, any> = {
    image: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z' })
    ]),
    video: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z' })
    ]),
    audio: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3' })
    ]),
    archive: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 8l4 4 4-4m6-4v12a2 2 0 01-2 2H7a2 2 0 01-2-2V4a2 2 0 012-2h10a2 2 0 012 2z' })
    ]),
    file: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' })
    ])
  }

  return iconMap[icon] || iconMap.file
}
</script>

<style scoped>
.container {
  max-width: 1200px;
}

/* 输入框样式 */
.tracking-widest {
  letter-spacing: 0.1em;
}

/* 动画效果 */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 卡片阴影 */
.shadow-lg {
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

/* 悬停效果 */
.hover\:bg-gray-100:hover {
  transition: background-color 0.2s ease;
}

/* 滚动条 */
.overflow-y-auto {
  scrollbar-width: thin;
  scrollbar-color: rgba(156, 163, 175, 0.5) transparent;
}

.overflow-y-auto::-webkit-scrollbar {
  width: 8px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.5);
  border-radius: 4px;
}
</style>